package box

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/box/docker"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/event"
	"github.com/hckops/hckctl/pkg/template"
)

// TODO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

func addBoxProviderFlag(command *cobra.Command) {
	const (
		providerFlagName = "provider"
	)
	command.Flags().StringP(providerFlagName, common.NoneFlagShortHand, string(model.Docker),
		fmt.Sprintf("switch box provider, one of %s", strings.Join(model.BoxProviderValues(), "|")))
	viper.BindPFlag(fmt.Sprintf("box.%s", providerFlagName), command.Flags().Lookup(providerFlagName))
}

func NewBoxClient(provider model.BoxProvider) (box.BoxClient, error) {
	opts := model.NewBoxOpts()
	switch provider {
	case model.Docker:
		return docker.NewDockerBox(opts)
	case model.Kubernetes:
		// TODO
		return nil, nil
	case model.Argo:
		// TODO
		return nil, nil
	case model.Cloud:
		// TODO
		return nil, nil
	default:
		return nil, errors.New("invalid provider")
	}
}

// TODO refactor createBox
func openBox(src template.TemplateSource, configRef *config.ConfigRef) error {

	boxTemplate, err := src.ReadBox()
	if err != nil {
		log.Warn().Err(err).Msg("error reading template")
		return errors.New("invalid template")
	}

	loader := common.NewLoader()
	loader.Start("loading template %s", boxTemplate.Name)
	defer loader.Stop()

	provider := configRef.Config.Box.Provider
	log.Debug().Msgf("opening box: provider=%s name=%s\n%s", provider, boxTemplate.Name, boxTemplate.Pretty())

	boxClient, err := NewBoxClient(provider)
	if err != nil {
		log.Warn().Err(err).Msg("error creating client")
		return errors.New("client error")
	}

	handleOpenEvents(boxClient, loader)

	if err := boxClient.Open(boxTemplate); err != nil {
		log.Warn().Err(err).Msg("error opening box")
		return errors.New("open error")
	}
	return nil
}

func handleOpenEvents(boxClient box.BoxClient, loader *common.Loader) {
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
}

// TODO refactor
func createBox(src template.TemplateSource, configRef *config.ConfigRef) error {

	boxTemplate, err := src.ReadBox()
	if err != nil {
		log.Warn().Err(err).Msg("error reading template")
		return errors.New("invalid template")
	}

	loader := common.NewLoader()
	loader.Start("loading template %s", boxTemplate.Name)
	defer loader.Stop()

	provider := configRef.Config.Box.Provider
	log.Debug().Msgf("creating box: provider=%s name=%s\n%s", provider, boxTemplate.Name, boxTemplate.Pretty())

	boxClient, err := NewBoxClient(provider)
	if err != nil {
		log.Warn().Err(err).Msg("error creating client")
		return errors.New("client error")
	}

	handleOpenEvents(boxClient, loader)

	if boxInfo, err := boxClient.Create(boxTemplate); err != nil {
		log.Warn().Err(err).Msg("error creating box")
		return errors.New("create error")
	} else {
		loader.Stop()
		fmt.Println(boxInfo.Name)
	}
	return nil
}

// resolve box by name
func runBoxClient(configRef *config.ConfigRef, boxName string, run func(box.BoxClient, *model.BoxV1) error) error {

	// TODO create container with Labels=revision?
	// best effort mode to resolve default template
	// WARNING this might return unexpected results if the container was created with a different revision
	revision := common.TemplateSourceRevision

	revisionOpts := &template.RevisionOpts{
		SourceCacheDir: configRef.Config.Template.CacheDir,
		SourceUrl:      common.TemplateSourceUrl,
		SourceRevision: common.TemplateSourceRevision,
		Revision:       revision,
	}

	templateName := model.ToBoxTemplateName(boxName)
	boxTemplate, err := template.NewRemoteSource(revisionOpts, templateName).ReadBox()
	if err != nil {
		log.Warn().Err(err).Msg("error reading box template")
		return errors.New("invalid template")
	}

	// TODO how to resolve provider without attempting all of them?
	boxClient, err := NewBoxClient(model.Docker)
	if err != nil {
		log.Warn().Err(err).Msg("error creating client")
		return errors.New("client error")
	}
	boxClient.Events().Subscribe(func(event event.Event) {
		log.Debug().Msg(event.String())
	})

	if err := run(boxClient, boxTemplate); err != nil {
		log.Warn().Err(err).Msg("error invoking client")
		return errors.New("run error")
	}

	return nil
}
