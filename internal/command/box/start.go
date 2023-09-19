package box

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	boxFlag "github.com/hckops/hckctl/internal/command/box/flag"
	commonCmd "github.com/hckops/hckctl/internal/command/common"
	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/internal/command/config"
	boxModel "github.com/hckops/hckctl/pkg/box/model"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/template"
)

type boxStartCmdOptions struct {
	configRef    *config.ConfigRef
	sourceFlag   *commonFlag.SourceFlag
	providerFlag *commonFlag.ProviderFlag
	provider     boxModel.BoxProvider
}

func NewBoxStartCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxStartCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:     "start [name]",
		Short:   "Start a long running detached box",
		Args:    cobra.ExactArgs(1),
		PreRunE: opts.validate,
		RunE:    opts.run,
	}

	// --revision or --local
	opts.sourceFlag = commonFlag.AddSourceFlag(command)
	// --provider (enum)
	opts.providerFlag = boxFlag.AddBoxProviderFlag(command)

	return command
}

func (opts *boxStartCmdOptions) validate(cmd *cobra.Command, args []string) error {

	validProvider, err := boxFlag.ValidateBoxProvider(opts.configRef.Config.Box.Provider, opts.providerFlag)
	if err != nil {
		return err
	}
	opts.provider = validProvider

	if err := commonFlag.ValidateSourceFlag(opts.providerFlag, opts.sourceFlag); err != nil {
		log.Warn().Err(err).Msgf(commonFlag.ErrorFlagNotSupported)
		return errors.New(commonFlag.ErrorFlagNotSupported)
	}
	return nil
}

func (opts *boxStartCmdOptions) run(cmd *cobra.Command, args []string) error {

	if opts.sourceFlag.Local {
		path := args[0]
		log.Debug().Msgf("start box from local template: path=%s", path)

		sourceLoader := template.NewLocalCachedLoader[boxModel.BoxV1](path, opts.configRef.Config.Template.CacheDir)
		return opts.startBox(sourceLoader, boxModel.NewBoxLabels().AddDefaultLocal())

	} else {
		name := args[0]
		log.Debug().Msgf("start box from git template: name=%s revision=%s", name, opts.sourceFlag.Revision)

		sourceOpts := commonCmd.NewGitSourceOptions(opts.configRef.Config.Template.CacheDir, opts.sourceFlag.Revision)
		sourceLoader := template.NewGitLoader[boxModel.BoxV1](sourceOpts, name)
		labels := boxModel.NewBoxLabels().AddDefaultGit(sourceOpts.RepositoryUrl, sourceOpts.DefaultRevision, sourceOpts.CacheDirName())
		return opts.startBox(sourceLoader, labels)
	}
}

func (opts *boxStartCmdOptions) startBox(sourceLoader template.SourceLoader[boxModel.BoxV1], labels commonModel.Labels) error {

	createClient := func(invokeOpts *invokeOptions) error {

		createOpts, err := newCreateOptions(invokeOpts.template, labels, opts.configRef.Config.Box.Size)
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
	return runBoxClient(sourceLoader, opts.provider, opts.configRef, createClient)
}
