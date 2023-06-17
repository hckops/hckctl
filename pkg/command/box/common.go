package box

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/event"
	"github.com/hckops/hckctl/pkg/template"
)

type boxClientOpts struct {
	client   box.BoxClient
	template *model.BoxV1
	loader   *common.Loader
}

func runBoxClient(src template.TemplateSource, provider model.BoxProvider, invokeClient func(*boxClientOpts) error) error {

	boxTemplate, err := src.ReadBox()
	if err != nil {
		log.Warn().Err(err).Msg("error reading template")
		return errors.New("invalid template")
	}

	loader := common.NewLoader()
	loader.Start("loading template %s", boxTemplate.Name)
	defer loader.Stop()

	log.Info().Msgf("loading template: provider=%s name=%s\n%s", provider, boxTemplate.Name, boxTemplate.Pretty())

	boxClient, err := box.NewBoxClient(provider)
	if err != nil {
		log.Warn().Err(err).Msgf("error creating client: provider=%v", provider)
		return fmt.Errorf("create %v client error", provider)
	}

	boxClient.Events().Subscribe(func(e event.Event) {
		switch e.Kind() {
		case event.PrintConsole:
			loader.Refresh("loading")
			fmt.Println(e.String())
		case event.LoaderUpdate:
			loader.Refresh(e.String())
		case event.LoaderStop:
			loader.Stop()
		default:
			log.Debug().Msgf("[%v][%s] %s", e.Source(), e.Kind(), e.String())
		}
	})

	opts := &boxClientOpts{
		client:   boxClient,
		template: boxTemplate,
		loader:   loader,
	}
	if err := invokeClient(opts); err != nil {
		log.Warn().Err(err).Msgf("error invoking client: provider=%v", provider)
		return fmt.Errorf("invoke %v client error", provider)
	}
	return nil
}

func runRemoteBoxClient(configRef *config.ConfigRef, boxName string, invokeClient func(box.BoxClient, *model.BoxV1) error) error {

	// best effort approach to resolve remote box template by name with default revision
	// WARNING this might return unexpected results if the container was created with a different revision
	revisionOpts := &template.RevisionOpts{
		SourceCacheDir: configRef.Config.Template.CacheDir,
		SourceUrl:      common.TemplateSourceUrl,
		SourceRevision: common.TemplateSourceRevision,
		Revision:       common.TemplateSourceRevision, // TODO create container with Labels="com.hckops.revision=<REVISION>" to resolve exact template
	}
	templateName := model.ToBoxTemplateName(boxName)
	boxTemplate, err := template.NewRemoteSource(revisionOpts, templateName).ReadBox()
	if err != nil {
		log.Warn().Err(err).Msgf("error reading box template: templateName=%v", templateName)
		return errors.New("invalid template")
	}

	// TODO attempt all providers: model.BoxProviders()
	for _, provider := range []model.BoxProvider{model.Docker} {
		boxClient, err := newDefaultBoxClient(provider)
		if err != nil {
			return err
		}

		if err := invokeClient(boxClient, boxTemplate); err != nil {
			log.Warn().Err(err).Msgf("error invoking client: provider=%v", provider)
			return fmt.Errorf("invoke %v client error", provider)
		}
	}

	return nil
}

func newDefaultBoxClient(provider model.BoxProvider) (box.BoxClient, error) {
	boxClient, err := box.NewBoxClient(provider)
	if err != nil {
		log.Warn().Err(err).Msgf("error creating client: provider=%v", provider)
		return nil, fmt.Errorf("create %v client error", provider)
	}

	boxClient.Events().Subscribe(func(e event.Event) {
		log.Debug().Msgf("[%v][%s] %s", e.Source(), e.Kind(), e.String())
	})
	return boxClient, nil
}
