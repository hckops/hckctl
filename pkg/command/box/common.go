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
	client   box.BoxClient
	template *model.BoxV1
	loader   *common.Loader
}

// open and create
func runBoxClient(src template.TemplateSource, provider model.BoxProvider, configRef *config.ConfigRef, invokeClient func(*boxClientOptions) error) error {

	boxTemplate, err := src.ReadBox()
	if err != nil {
		log.Warn().Err(err).Msg("error reading template")
		return errors.New("invalid template")
	}

	loader := common.NewLoader()
	loader.Start("loading template %s", boxTemplate.Name)
	defer loader.Stop()

	log.Info().Msgf("loading template: provider=%s name=%s\n%s", provider, boxTemplate.Name, boxTemplate.Pretty())

	boxOpts, err := newBoxOpts(provider, configRef)
	if err != nil {
		return err
	}
	boxClient, err := box.NewBoxClient(boxOpts)
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
		default:
			log.Debug().Msgf("[%v][%s] %s", e.Source(), e.Kind(), e.String())
		}
	})

	opts := &boxClientOptions{
		client:   boxClient,
		template: boxTemplate,
		loader:   loader,
	}
	if err := invokeClient(opts); err != nil {
		log.Warn().Err(err).Msgf("error invoking client: provider=%v", provider)
		return fmt.Errorf("error %v client invoke", provider)
	}
	return nil
}

// exec and delete
func attemptRunBoxClients(configRef *config.ConfigRef, boxName string, invokeClient func(box.BoxClient, *model.BoxV1) error) error {

	// best effort approach to resolve remote box template by name with default revision
	// WARNING this might return unexpected results if the box was created with a different revision
	revisionOpts := &template.RevisionOptions{
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

	// silently fail attempting all the providers
	for _, providerFlag := range boxFlag.BoxProviders() {
		log.Debug().Msgf("attempt box template: providerFlag=%v", providerFlag)

		boxClient, err := newDefaultBoxClient(providerFlag, configRef)
		if err != nil {
			log.Warn().Err(err).Msgf("ignoring error default client: providerFlag=%v", providerFlag)
			// skip to the next provider
			break
		}

		if err := invokeClient(boxClient, boxTemplate); err != nil {
			log.Warn().Err(err).Msgf("ignoring error invoking client: providerFlag=%v", providerFlag)
		} else {
			// return as soon as the client is invoked with success
			return nil
		}
	}
	// nothing happened and all the providers failed
	return errors.New("not found")
}

func newBoxOpts(provider model.BoxProvider, configRef *config.ConfigRef) (*model.BoxOptions, error) {

	kubeClientConfig, err := configRef.Config.Provider.Kube.ToKubeClientConfig()
	if err != nil {
		log.Warn().Err(err).Msgf("error kube config: kubeConfig=%v", configRef.Config.Provider.Kube)
		return nil, errors.Wrap(err, "invalid kube config")
	}

	dockerClientConfig := configRef.Config.Provider.Docker.ToDockerClientConfig()
	sshClientConfig := configRef.Config.Provider.Cloud.ToSshClientConfig()
	internalOpts := model.NewBoxInternalOpts(version.ClientVersion())
	boxOpts := &model.BoxOptions{
		Provider:     provider,
		DockerConfig: dockerClientConfig,
		KubeConfig:   kubeClientConfig,
		SshConfig:    sshClientConfig,
		InternalOpts: internalOpts,
	}
	return boxOpts, nil
}

func newDefaultBoxClient(providerFlag commonFlag.ProviderFlag, configRef *config.ConfigRef) (box.BoxClient, error) {

	provider, err := boxFlag.ToBoxProvider(providerFlag)
	if err != nil {
		return nil, err
	}
	opts, err := newBoxOpts(provider, configRef)
	if err != nil {
		return nil, err
	}
	boxClient, err := box.NewBoxClient(opts)
	if err != nil {
		log.Warn().Err(err).Msgf("error creating client: provider=%v", opts.Provider)
		return nil, fmt.Errorf("create %v client error", opts.Provider)
	}

	boxClient.Events().Subscribe(func(e event.Event) {
		switch e.Kind() {
		case event.PrintConsole:
			fmt.Println(e.String())
		default:
			log.Debug().Msgf("[%v][%s] %s", e.Source(), e.Kind(), e.String())
		}
	})
	return boxClient, nil
}
