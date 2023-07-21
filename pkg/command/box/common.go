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

type invokeOptions struct {
	client   box.BoxClient
	template *template.TemplateInfo[model.BoxV1]
}

// open and create
func runBoxClient(sourceLoader template.SourceLoader[model.BoxV1], provider model.BoxProvider, configRef *config.ConfigRef, invokeClient func(*invokeOptions, *common.Loader) error) error {

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

	invokeOpts := &invokeOptions{
		client:   boxClient,
		template: boxTemplate,
	}
	if err := invokeClient(invokeOpts, loader); err != nil {
		log.Warn().Err(err).Msgf("error invoking client: provider=%v", provider)
		return fmt.Errorf("error %v client invoke", provider)
	}
	return nil
}

// exec, describe and delete
func attemptRunBoxClients(configRef *config.ConfigRef, boxName string, invokeClient func(*invokeOptions) error) error {

	// silently fail attempting all the providers
	for _, providerFlag := range boxFlag.BoxProviders() {
		log.Debug().Msgf("attempt box template: providerFlag=%v", providerFlag)

		boxClient, err := newDefaultBoxClient(providerFlag, configRef)
		if err != nil {
			log.Warn().Err(err).Msgf("ignoring error default client: providerFlag=%v", providerFlag)
			// skip to the next provider
			break
		}

		boxDetails, err := boxClient.Describe(boxName)
		if err != nil {
			log.Warn().Err(err).Msgf("ignoring error describe box: providerFlag=%v", providerFlag)
			break
		}

		sourceLoader, err := newSourceLoader(boxDetails, configRef.Config.Template.CacheDir)
		if err != nil {
			log.Warn().Err(err).Msgf("ignoring error source loader: providerFlag=%v", providerFlag)
			break
		}

		templateInfo, err := sourceLoader.Read()
		if err != nil {
			log.Warn().Err(err).Msgf("ignoring error reading source: providerFlag=%v ", providerFlag)
			break
		}

		invokeOpts := &invokeOptions{
			client:   boxClient,
			template: templateInfo,
		}
		if err := invokeClient(invokeOpts); err != nil {
			log.Warn().Err(err).Msgf("ignoring error invoking client: providerFlag=%v", providerFlag)
		} else {
			// return as soon as the client is invoked with success
			return nil
		}
	}
	// nothing happened and all the providers failed
	return errors.New("not found")
}

// TODO issue sourceType

func newSourceLoader(boxDetails *model.BoxDetails, cacheDir string) (template.SourceLoader[model.BoxV1], error) {

	//switch boxDetails.TemplateInfo.SourceType {
	//case template.Local:
	//	return template.NewLocalLoader[model.BoxV1](boxDetails.TemplateInfo.LocalTemplate.Path), nil
	//case template.Git:
	//	sourceOpts := &template.GitSourceOptions{
	//		CacheBaseDir:    cacheDir,
	//		RepositoryUrl:   common.TemplateSourceUrl,
	//		DefaultRevision: common.TemplateSourceRevision,
	//		Revision:        boxDetails.TemplateInfo.GitTemplate.Commit,
	//		AllowOffline:    true,
	//	}
	//	return template.NewGitLoader[model.BoxV1](sourceOpts, boxDetails.Info.Name), nil
	//}
	return nil, errors.New("invalid source type")
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
