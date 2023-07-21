package box

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/box/model"
	boxFlag "github.com/hckops/hckctl/pkg/command/box/flag"
	"github.com/hckops/hckctl/pkg/command/common"
	commonFlag "github.com/hckops/hckctl/pkg/command/common/flag"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/command/version"
	"github.com/hckops/hckctl/pkg/event"
	"github.com/hckops/hckctl/pkg/template"
)

type boxClientOptions struct {
	template *template.TemplateInfo[model.BoxV1]
	client   box.BoxClient
	loader   *common.Loader
}

// open and create
func runBoxClient(sourceLoader template.SourceLoader[model.BoxV1], provider model.BoxProvider, configRef *config.ConfigRef, invokeClient func(*boxClientOptions) error) error {

	boxTemplate, err := sourceLoader.Read()
	if err != nil {
		log.Warn().Err(err).Msg("error reading template")
		return errors.New("invalid template")
	}

	loader := common.NewLoader()
	loader.Start("loading template %s", boxTemplate.Value.Data.Name)
	defer loader.Stop()

	log.Info().Msgf("loading template: provider=%s name=%s\n%s", provider, boxTemplate.Value.Data.Name, boxTemplate.Value.Data.Pretty())

	boxClientOpts := newBoxClientOpts(provider, configRef)
	boxClient, err := box.NewBoxClient(boxClientOpts)
	if err != nil {
		log.Warn().Err(err).Msgf("error creating client: provider=%v", provider)
		return fmt.Errorf("error %v client create", provider)
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
		case event.LogWarning:
			log.Warn().Msgf("[%v] %s", e.Source(), e.String())
		default:
			log.Debug().Msgf("[%v][%s] %s", e.Source(), e.Kind(), e.String())
		}
	})

	opts := &boxClientOptions{
		template: boxTemplate,
		client:   boxClient,
		loader:   loader,
	}
	if err := invokeClient(opts); err != nil {
		log.Warn().Err(err).Msgf("error invoking client: provider=%v", provider)
		return fmt.Errorf("error %v client invoke", provider)
	}
	return nil
}

// exec, describe and delete
func attemptRunBoxClients(configRef *config.ConfigRef, boxName string, invokeClient func(box.BoxClient, *model.BoxV1) error) error {

	// TODO describe box and depending on the source type instantiate sourceLoader

	// best effort approach to resolve the box template by name with git source and default revision
	// WARNING this might return unexpected results if the box was created with a different revision
	sourceOpts := &template.GitSourceOptions{
		CacheBaseDir:    configRef.Config.Template.CacheDir,
		RepositoryUrl:   common.TemplateSourceUrl,
		DefaultRevision: common.TemplateSourceRevision,
		Revision:        common.TemplateSourceRevision, // TODO always default, read from labels and convert main to sha
		AllowOffline:    true,
	}
	// TODO add name to label and search for all provider
	templateName := model.ToBoxTemplateName(boxName)
	boxTemplate, err := template.NewGitLoader[model.BoxV1](sourceOpts, templateName).Read()
	if err != nil {
		log.Warn().Err(err).Msgf("error reading box template: templateName=%v", templateName)
		return errors.New("invalid template")
	}

	// silently fail attempting all the providers
	for _, providerFlag := range boxFlag.BoxProviders() {
		log.Debug().Msgf("attempt box template: providerFlag=%v", providerFlag)

		boxClient, err := newDefaultBoxClient(providerFlag, configRef)
		if err != nil {
			log.Warn().Err(err).Msgf("ignoring error default client: providerFlag=%v", providerFlag)
			// skip to the next provider
			break
		}

		if err := invokeClient(boxClient, &boxTemplate.Value.Data); err != nil {
			log.Warn().Err(err).Msgf("ignoring error invoking client: providerFlag=%v", providerFlag)
		} else {
			// return as soon as the client is invoked with success
			return nil
		}
	}
	// nothing happened and all the providers failed
	return errors.New("not found")
}

func newBoxClientOpts(provider model.BoxProvider, configRef *config.ConfigRef) *model.BoxClientOptions {

	return &model.BoxClientOptions{
		Provider:   provider,
		CommonOpts: model.NewCommonBoxOpts(),
		DockerOpts: configRef.Config.Provider.Docker.ToDockerBoxOptions(),
		KubeOpts:   configRef.Config.Provider.Kube.ToKubeBoxOptions(),
		CloudOpts:  configRef.Config.Provider.Cloud.ToCloudBoxOptions(version.ClientVersion()),
	}
}

func newDefaultBoxClient(providerFlag commonFlag.ProviderFlag, configRef *config.ConfigRef) (box.BoxClient, error) {

	provider, err := boxFlag.ToBoxProvider(providerFlag)
	if err != nil {
		return nil, err
	}
	opts := newBoxClientOpts(provider, configRef)
	boxClient, err := box.NewBoxClient(opts)
	if err != nil {
		log.Warn().Err(err).Msgf("error creating client: provider=%v", opts.Provider)
		return nil, fmt.Errorf("create %v client error", opts.Provider)
	}

	boxClient.Events().Subscribe(func(e event.Event) {
		switch e.Kind() {
		case event.PrintConsole:
			fmt.Println(e.String())
		case event.LogWarning:
			log.Warn().Msgf("[%v] %s", e.Source(), e.String())
		default:
			log.Debug().Msgf("[%v][%s] %s", e.Source(), e.Kind(), e.String())
		}
	})
	return boxClient, nil
}

func newTemplateOptions(info *template.TemplateInfo[model.BoxV1], labels model.BoxLabels, sizeValue string) (*model.TemplateOptions, error) {
	size, err := model.ExistResourceSize(sizeValue)
	if err != nil {
		return nil, err
	}

	var allLabels model.BoxLabels
	switch info.SourceType {
	case template.Local:
		allLabels = labels.AddLocalLabels(size, info.Path)
	case template.Git:
		allLabels = labels.AddGitLabels(size, info.Path, info.Revision)
	}

	templateOpts := &model.TemplateOptions{
		Template: &info.Value.Data,
		Size:     size,
		Labels:   allLabels,
	}
	return templateOpts, nil
}
