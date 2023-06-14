package box

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/event"
	"github.com/hckops/hckctl/pkg/template"
)

func addBoxProviderFlag(command *cobra.Command) {
	const (
		providerFlagName = "provider"
	)
	command.Flags().StringP(providerFlagName, common.NoneFlagShortHand, string(model.Docker),
		fmt.Sprintf("switch box provider, one of %s", strings.Join(model.BoxProviderValues(), "|")))
	viper.BindPFlag(fmt.Sprintf("box.%s", providerFlagName), command.Flags().Lookup(providerFlagName))
}

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

	// TODO attempt all providers
	providers := []model.BoxProvider{model.Docker}
	provider := providers[0]

	boxClient, err := box.NewBoxClient(provider)
	if err != nil {
		log.Warn().Err(err).Msgf("error creating client: provider=%v", provider)
		return fmt.Errorf("create %v client error", provider)
	}
	boxClient.Events().Subscribe(func(e event.Event) {
		log.Debug().Msgf("[%v][%s] %s", e.Source(), e.Kind(), e.String())
	})

	if err := invokeClient(boxClient, boxTemplate); err != nil {
		log.Warn().Err(err).Msgf("error invoking client: provider=%v", provider)
		return fmt.Errorf("invoke %v client error", provider)
	}
	return nil
}
