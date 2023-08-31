package box

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	boxFlag "github.com/hckops/hckctl/internal/command/box/flag"
	commonCmd "github.com/hckops/hckctl/internal/command/common"
	"github.com/hckops/hckctl/internal/command/config"
	"github.com/hckops/hckctl/internal/command/version"
	"github.com/hckops/hckctl/pkg/box"
	boxModel "github.com/hckops/hckctl/pkg/box/model"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/template"
)

type invokeOptions struct {
	client   box.BoxClient
	template *template.TemplateInfo[boxModel.BoxV1]
	loader   *commonCmd.Loader
}

// start and temporary
func runBoxClient(sourceLoader template.SourceLoader[boxModel.BoxV1], provider boxModel.BoxProvider, configRef *config.ConfigRef, invokeClient func(*invokeOptions) error) error {

	boxTemplate, err := sourceLoader.Read()
	if err != nil || boxTemplate.Value.Kind != schema.KindBoxV1 {
		log.Warn().Err(err).Msg("error reading template")
		return errors.New("invalid template")
	}

	loader := commonCmd.NewLoader()
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

// open, info and stop-one
func attemptRunBoxClients(configRef *config.ConfigRef, boxName string, invokeClient func(*invokeOptions, *boxModel.BoxDetails) error) error {

	loader := commonCmd.NewLoader()
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
		if err != nil || templateInfo.Value.Kind != schema.KindBoxV1 {
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

func newSourceLoader(boxDetails *boxModel.BoxDetails, cacheDir string) template.SourceLoader[boxModel.BoxV1] {
	if boxDetails.TemplateInfo.IsCached() {
		return template.NewLocalLoader[boxModel.BoxV1](boxDetails.TemplateInfo.CachedTemplate.Path)
	} else {
		sourceOpts := commonCmd.NewGitSourceOptions(cacheDir, boxDetails.TemplateInfo.GitTemplate.Commit)
		return template.NewGitLoader[boxModel.BoxV1](sourceOpts, boxDetails.TemplateInfo.GitTemplate.Name)
	}
}

func newBoxClientOpts(provider boxModel.BoxProvider, configRef *config.ConfigRef) *boxModel.BoxClientOptions {
	return &boxModel.BoxClientOptions{
		Provider:   provider,
		DockerOpts: configRef.Config.Provider.Docker.ToDockerOptions(),
		KubeOpts:   configRef.Config.Provider.Kube.ToKubeOptions(),
		CloudOpts:  configRef.Config.Provider.Cloud.ToCloudOptions(version.ClientVersion()),
	}
}

func newDefaultBoxClient(provider boxModel.BoxProvider, configRef *config.ConfigRef, loader *commonCmd.Loader) (box.BoxClient, error) {

	boxClientOpts := newBoxClientOpts(provider, configRef)
	boxClient, err := box.NewBoxClient(boxClientOpts)
	if err != nil {
		log.Error().Err(err).Msgf("error box client provider=%s", provider)
		return nil, fmt.Errorf("error %s client", provider)
	}

	boxClient.Events().Subscribe(commonCmd.EventCallback(loader))
	return boxClient, nil
}

func newCreateOptions(info *template.TemplateInfo[boxModel.BoxV1], labels commonModel.Labels, sizeValue string) (*boxModel.CreateOptions, error) {
	size, err := boxModel.ExistResourceSize(sizeValue)
	if err != nil {
		return nil, err
	}

	var allLabels commonModel.Labels
	switch info.SourceType {
	case template.Local:
		allLabels = boxModel.AddBoxSize(labels, size).AddLocal(info.Path)
	case template.Git:
		allLabels = boxModel.AddBoxSize(labels, size).AddGit(info.Path, info.Revision)
	}

	return &boxModel.CreateOptions{
		Template: &info.Value.Data,
		Size:     size,
		Labels:   allLabels,
	}, nil
}
