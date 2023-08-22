package box

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	boxFlag "github.com/hckops/hckctl/internal/command/box/flag"
	"github.com/hckops/hckctl/internal/command/common"
	"github.com/hckops/hckctl/internal/command/config"
	"github.com/hckops/hckctl/internal/command/version"
	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/template"
)

type invokeOptions struct {
	client   box.BoxClient
	template *template.TemplateInfo[model.BoxV1]
	loader   *common.Loader
}

// create and temporary
func runBoxClient(sourceLoader template.SourceLoader[model.BoxV1], provider model.BoxProvider, configRef *config.ConfigRef, invokeClient func(*invokeOptions) error) error {

	boxTemplate, err := sourceLoader.Read()
	if err != nil {
		log.Warn().Err(err).Msg("error reading template")
		return errors.New("invalid template")
	}

	loader := common.NewLoader()
	loader.Start("loading template %s", boxTemplate.Value.Data.Name)
	defer loader.Stop()

	log.Info().Msgf("loading template: provider=%s name=%s\n%s", provider, boxTemplate.Value.Data.Name, boxTemplate.Value.Data.Pretty())

	boxClient, err := newDefaultBoxClient(provider, configRef, loader)
	if err != nil {
		return err
	}

	invokeOpts := &invokeOptions{
		client:   boxClient,
		template: boxTemplate,
		loader:   loader,
	}
	if err := invokeClient(invokeOpts); err != nil {
		log.Warn().Err(err).Msgf("error invoking client: provider=%v", provider)
		return fmt.Errorf("error %v client invoke", provider)
	}
	return nil
}

// open, info and delete-one
func attemptRunBoxClients(configRef *config.ConfigRef, boxName string, invokeClient func(*invokeOptions, *model.BoxDetails) error) error {

	loader := common.NewLoader()
	loader.Start("loading %s", boxName)
	defer loader.Stop()

	// silently fail attempting all the providers
	for _, providerFlag := range boxFlag.BoxProviders() {
		log.Debug().Msgf("attempt box template: providerFlag=%s boxName=%s", providerFlag, boxName)

		provider, err := boxFlag.ToBoxProvider(providerFlag)
		if err != nil {
			log.Warn().Err(err).Msgf("ignoring error provider: provider=%s", provider)
			// skip to the next provider
			continue
		}

		boxClient, err := newDefaultBoxClient(provider, configRef, loader)
		if err != nil {
			log.Warn().Err(err).Msgf("ignoring error default client: provider=%s", provider)
			continue
		}

		boxDetails, err := boxClient.Describe(boxName)
		if err != nil {
			log.Warn().Err(err).Msgf("ignoring error describe box: provider=%s boxName=%s", provider, boxName)
			continue
		}

		templateInfo, err := newSourceLoader(boxDetails, configRef.Config.Template.CacheDir).Read()
		if err != nil {
			log.Warn().Err(err).Msgf("ignoring error reading source: provider=%s boxName=%s", provider, boxName)
			continue
		}

		invokeOpts := &invokeOptions{
			client:   boxClient,
			template: templateInfo,
			loader:   loader,
		}
		if err := invokeClient(invokeOpts, boxDetails); err != nil {
			log.Warn().Err(err).Msgf("ignoring error invoking client: provider=%s boxName=%s", provider, boxName)
		} else {
			// return as soon as the client is invoked with success
			return nil
		}
	}
	// nothing happened and all the providers failed
	return errors.New("not found")
}

func newSourceLoader(boxDetails *model.BoxDetails, cacheDir string) template.SourceLoader[model.BoxV1] {
	if boxDetails.TemplateInfo.IsCached() {
		return template.NewLocalLoader[model.BoxV1](boxDetails.TemplateInfo.CachedTemplate.Path)
	} else {
		sourceOpts := newGitSourceOptions(cacheDir, boxDetails.TemplateInfo.GitTemplate.Commit)
		return template.NewGitLoader[model.BoxV1](sourceOpts, boxDetails.TemplateInfo.GitTemplate.Name)
	}
}

func newGitSourceOptions(cacheDir string, revision string) *template.GitSourceOptions {
	return &template.GitSourceOptions{
		CacheBaseDir:    cacheDir,
		RepositoryUrl:   common.TemplateSourceUrl,
		DefaultRevision: common.TemplateSourceRevision,
		Revision:        revision,
		AllowOffline:    true,
	}
}

func newBoxClientOpts(provider model.BoxProvider, configRef *config.ConfigRef) *model.BoxClientOptions {
	return &model.BoxClientOptions{
		Provider:   provider,
		DockerOpts: configRef.Config.Provider.Docker.ToDockerOptions(),
		KubeOpts:   configRef.Config.Provider.Kube.ToKubeOptions(),
		CloudOpts:  configRef.Config.Provider.Cloud.ToCloudOptions(version.ClientVersion()),
	}
}

func newDefaultBoxClient(provider model.BoxProvider, configRef *config.ConfigRef, loader *common.Loader) (box.BoxClient, error) {

	boxClientOpts := newBoxClientOpts(provider, configRef)
	boxClient, err := box.NewBoxClient(boxClientOpts)
	if err != nil {
		log.Error().Err(err).Msgf("error box client provider=%s", provider)
		return nil, fmt.Errorf("error %s client", provider)
	}

	boxClient.Events().Subscribe(common.EventCallback(loader))
	return boxClient, nil
}

func newCreateOptions(info *template.TemplateInfo[model.BoxV1], labels model.BoxLabels, sizeValue string) (*model.CreateOptions, error) {
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

	return &model.CreateOptions{
		Template: &info.Value.Data,
		Size:     size,
		Labels:   allLabels,
	}, nil
}
