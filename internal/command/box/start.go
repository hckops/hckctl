package box

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	boxFlag "github.com/hckops/hckctl/internal/command/box/flag"
	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/internal/command/config"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/template"
)

type boxStartCmdOptions struct {
	configRef    *config.ConfigRef
	sourceFlag   *commonFlag.SourceFlag
	providerFlag *commonFlag.ProviderFlag
}

func NewBoxStartCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxStartCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "start [name]",
		Short: "Start a long running detached box",
		RunE:  opts.run,
	}

	// --provider (enum)
	opts.providerFlag = boxFlag.AddBoxProviderFlag(command)
	// --revision or --local
	opts.sourceFlag = commonFlag.AddTemplateSourceFlag(command)

	return command
}

func (opts *boxStartCmdOptions) run(cmd *cobra.Command, args []string) error {

	provider, err := boxFlag.ValidateBoxProvider(opts.configRef.Config.Box.Provider, opts.providerFlag)
	if err != nil {
		return err
	} else if len(args) == 1 {

		if err := boxFlag.ValidateSourceFlag(provider, opts.sourceFlag); err != nil {
			log.Warn().Err(err).Msgf(commonFlag.ErrorFlagNotSupported)
			return errors.New(commonFlag.ErrorFlagNotSupported)

		} else if opts.sourceFlag.Local {
			path := args[0]
			log.Debug().Msgf("start box from local template: path=%s", path)

			sourceLoader := template.NewLocalCachedLoader[model.BoxV1](path, opts.configRef.Config.Template.CacheDir)
			return startBox(sourceLoader, provider, opts.configRef, model.NewLocalLabels())

		} else {
			name := args[0]
			log.Debug().Msgf("start box from git template: name=%s revision=%s", name, opts.sourceFlag.Revision)

			sourceOpts := newGitSourceOptions(opts.configRef.Config.Template.CacheDir, opts.sourceFlag.Revision)
			sourceLoader := template.NewGitLoader[model.BoxV1](sourceOpts, name)
			labels := model.NewGitLabels(sourceOpts.RepositoryUrl, sourceOpts.DefaultRevision, sourceOpts.CacheDirName())
			return startBox(sourceLoader, provider, opts.configRef, labels)
		}
	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func startBox(sourceLoader template.SourceLoader[model.BoxV1], provider model.BoxProvider, configRef *config.ConfigRef, labels model.BoxLabels) error {

	createClient := func(invokeOpts *invokeOptions) error {

		createOpts, err := newCreateOptions(invokeOpts.template, labels, configRef.Config.Box.Size)
		if err != nil {
			return err
		}
		if boxInfo, err := invokeOpts.client.Create(createOpts); err != nil {
			return err
		} else {
			invokeOpts.loader.Stop()
			fmt.Println(boxInfo.Name)
		}
		return nil
	}
	return runBoxClient(sourceLoader, provider, configRef, createClient)
}
