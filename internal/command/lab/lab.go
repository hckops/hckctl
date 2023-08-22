package lab

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/internal/command/common"
	"github.com/hckops/hckctl/internal/command/config"
	"github.com/hckops/hckctl/internal/command/version"
	"github.com/hckops/hckctl/pkg/event"
	"github.com/hckops/hckctl/pkg/lab"
	"github.com/hckops/hckctl/pkg/lab/model"
	"github.com/hckops/hckctl/pkg/template"
)

// TODO command create, list, describe, delete

type labCmdOptions struct {
	configRef *config.ConfigRef
}

func NewLabCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := labCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "lab [name]",
		Short: "TODO",
		Args:  cobra.ExactArgs(1),
		RunE:  opts.run,
	}

	return command
}

func (opts *labCmdOptions) run(cmd *cobra.Command, args []string) error {
	name := args[0]
	log.Debug().Msgf("start lab from git template: name=%s", name)

	sourceOpts := newGitSourceOptions(opts.configRef.Config.Template.CacheDir)
	sourceLoader := template.NewGitLoader[model.LabV1](sourceOpts, name)
	return startLab(sourceLoader, model.Cloud, opts.configRef)
}

func startLab(sourceLoader template.SourceLoader[model.LabV1], provider model.LabProvider, configRef *config.ConfigRef) error {

	labTemplate, err := sourceLoader.Read()
	if err != nil {
		log.Warn().Err(err).Msg("error reading template")
		return errors.New("invalid template")
	}

	loader := common.NewLoader()
	loader.Start("loading template %s", labTemplate.Value.Data.Name)
	defer loader.Stop()

	log.Info().Msgf("loading template: provider=%s name=%s\n%s", provider, labTemplate.Value.Data.Name, labTemplate.Value.Data.Pretty())

	labClient, err := newDefaultLabClient(provider, configRef, loader)
	if err != nil {
		return err
	}

	createOpts := &model.CreateOptions{
		Template:   &labTemplate.Value.Data,
		Parameters: map[string]string{},
		Labels:     map[string]string{}, // TODO
	}

	if labInfo, err := labClient.Create(createOpts); err != nil {
		return err
	} else {
		loader.Stop()
		fmt.Println(labInfo.Name)
	}
	return nil
}

func newGitSourceOptions(cacheDir string) *template.GitSourceOptions {
	return &template.GitSourceOptions{
		CacheBaseDir:    cacheDir,
		RepositoryUrl:   common.TemplateSourceUrl,
		DefaultRevision: common.TemplateSourceRevision,
		Revision:        common.TemplateSourceRevision,
		AllowOffline:    true,
	}
}

func newDefaultLabClient(provider model.LabProvider, configRef *config.ConfigRef, loader *common.Loader) (lab.LabClient, error) {
	labClientOpts := &model.LabClientOptions{
		Provider:  provider,
		CloudOpts: configRef.Config.Provider.Cloud.ToCloudOptions(version.ClientVersion()),
	}

	labClient, err := lab.NewLabClient(labClientOpts)
	if err != nil {
		log.Error().Err(err).Msgf("error lab client")
		return nil, errors.New("error client")
	}

	// TODO common/generics
	labClient.Events().Subscribe(func(e event.Event) {
		switch e.Kind() {
		case event.PrintConsole:
			loader.Refresh("loading")
			fmt.Println(e.String())
		case event.LoaderUpdate:
			loader.Refresh(e.String())
		case event.LoaderStop:
			loader.Stop()
		case event.LogInfo:
			log.Info().Msgf("[%v] %s", e.Source(), e.String())
		case event.LogWarning:
			log.Warn().Msgf("[%v] %s", e.Source(), e.String())
		case event.LogError:
			log.Error().Msgf("[%v] %s", e.Source(), e.String())
		default:
			log.Debug().Msgf("[%v] %s", e.Source(), e.String())
		}
	})
	return labClient, nil
}
