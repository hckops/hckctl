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
	configRef *config.ConfigRef
	// flags
	networkVpnFlag     string
	providerFlag       *commonFlag.ProviderFlag
	templateSourceFlag *commonFlag.TemplateSourceFlag
	// internal
	provider boxModel.BoxProvider
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

	// --network-vpn
	commonFlag.AddNetworkVpnFlag(command, &opts.networkVpnFlag)
	// --provider (enum)
	opts.providerFlag = boxFlag.AddBoxProviderFlag(command)
	// --revision or --local
	opts.templateSourceFlag = commonFlag.AddTemplateSourceFlag(command)

	return command
}

func (opts *boxStartCmdOptions) validate(cmd *cobra.Command, args []string) error {
	// provider
	if validProvider, err := boxFlag.ValidateBoxProviderFlag(opts.configRef.Config.Box.Provider, opts.providerFlag); err != nil {
		return err
	} else {
		opts.provider = validProvider
	}
	// network-vpn (after provider validation)
	if vpnNetworkInfo, err := commonFlag.ValidateNetworkVpnFlag(opts.networkVpnFlag, opts.configRef.Config.Network.VpnNetworks()); err != nil {
		return err
	} else if vpnNetworkInfo != nil && opts.provider == boxModel.Cloud {
		return fmt.Errorf("%s: use lab", commonFlag.ErrorFlagNotSupported)
	}
	// source
	if err := commonFlag.ValidateTemplateSourceFlag(opts.providerFlag, opts.templateSourceFlag); err != nil {
		log.Warn().Err(err).Msgf(commonFlag.ErrorFlagNotSupported)
		return errors.New(commonFlag.ErrorFlagNotSupported)
	}
	return nil
}

func (opts *boxStartCmdOptions) run(cmd *cobra.Command, args []string) error {

	if opts.templateSourceFlag.Local {
		path := args[0]
		log.Debug().Msgf("start box from local template: path=%s", path)

		sourceLoader := template.NewLocalCachedLoader[boxModel.BoxV1](path, opts.configRef.Config.Template.CacheDir)
		return opts.startBox(sourceLoader, boxModel.NewBoxLabels().AddDefaultLocal())

	} else {
		name := args[0]
		log.Debug().Msgf("start box from git template: name=%s revision=%s", name, opts.templateSourceFlag.Revision)

		sourceOpts := commonCmd.NewGitSourceOptions(opts.configRef.Config.Template.CacheDir, opts.templateSourceFlag.Revision)
		sourceLoader := template.NewGitLoader[boxModel.BoxV1](sourceOpts, name)
		labels := boxModel.NewBoxLabels().AddDefaultGit(sourceOpts.RepositoryUrl, sourceOpts.DefaultRevision, sourceOpts.CacheDirName())
		return opts.startBox(sourceLoader, labels)
	}
}

func (opts *boxStartCmdOptions) startBox(sourceLoader template.SourceLoader[boxModel.BoxV1], labels commonModel.Labels) error {

	createClient := func(invokeOpts *invokeOptions) error {

		createOpts, err := newCreateOptions(invokeOpts.template, labels, opts.configRef, opts.networkVpnFlag)
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
