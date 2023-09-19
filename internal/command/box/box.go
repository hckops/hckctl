package box

import (
	"github.com/MakeNowJust/heredoc"
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

type boxCmdOptions struct {
	configRef          *config.ConfigRef
	templateSourceFlag *commonFlag.TemplateSourceFlag
	providerFlag       *commonFlag.ProviderFlag
	provider           boxModel.BoxProvider
	tunnelFlag         *boxFlag.TunnelFlag
}

func NewBoxCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "box [name]",
		Short: "Access and tunnel a box",
		Long: heredoc.Doc(`
			Access and tunnel a box

			  Create and access an ephemeral Box, tunnelling locally all the open ports.
			  All public templates are versioned under the /box/ sub-path on GitHub
			  at https://github.com/hckops/megalopolis

			  Independently from the provider and the template used, it will spawn a shell
			  that when closed will automatically delete and cleanup all the resources.

			  The main purpose of a Box is to provide a ready-to-go and always up-to-date
			  environment with an uniformed experience, abstracting the actual providers
			  e.g. Docker, Kubernetes, etc.
		`),
		Example: heredoc.Doc(`

			# creates and accesses a temporary "box/base/parrot" docker container,
			# spawns a /bin/bash shell and tunnels the following ports:
			# (vnc)			vncviewer localhost:5900
			# (novnc)		http://localhost:6080
			# (tty)			http://localhost:7681
			hckctl box parrot

			# opens a box deployed on kubernetes (docker|kube|cloud)
			hckctl box kali --provider kube

			# opens a box tunneling all ports, without spawning a shell (ignored by docker)
			hckctl box arch --no-exec

			# opens a box spawning a shell, without tunneling the ports (ignored by docker)
			hckctl box alpine --no-tunnel

			# opens a box using a specific version (branch|tag|sha)
			hckctl box vulnerable/dvwa --revision main

			# opens a box defined locally
			hckctl box ../megalopolis/box/base/alpine.yml --local
		`),
		Args:    cobra.ExactArgs(1),
		PreRunE: opts.validate,
		RunE:    opts.run,
	}

	// --revision or --local
	opts.templateSourceFlag = commonFlag.AddTemplateSourceFlag(command)
	// --provider (enum)
	opts.providerFlag = boxFlag.AddBoxProviderFlag(command)
	// --no-exec or --no-tunnel
	opts.tunnelFlag = boxFlag.AddTunnelFlag(command)

	command.AddCommand(NewBoxInfoCmd(configRef))
	command.AddCommand(NewBoxListCmd(configRef))
	command.AddCommand(NewBoxOpenCmd(configRef))
	command.AddCommand(NewBoxStartCmd(configRef))
	command.AddCommand(NewBoxStopCmd(configRef))

	return command
}

func (opts *boxCmdOptions) validate(cmd *cobra.Command, args []string) error {

	validProvider, err := boxFlag.ValidateBoxProvider(opts.configRef.Config.Box.Provider, opts.providerFlag)
	if err != nil {
		return err
	}
	opts.provider = validProvider

	if err := commonFlag.ValidateTemplateSourceFlag(opts.providerFlag, opts.templateSourceFlag); err != nil {
		log.Warn().Err(err).Msgf(commonFlag.ErrorFlagNotSupported)
		return errors.New(commonFlag.ErrorFlagNotSupported)
	}

	if err := boxFlag.ValidateTunnelFlag(opts.tunnelFlag, opts.provider); err != nil {
		log.Warn().Err(err).Msgf("ignore validation %s", commonFlag.ErrorFlagNotSupported)
		// ignore validation
		return nil
	}
	return nil
}

func (opts *boxCmdOptions) run(cmd *cobra.Command, args []string) error {

	if opts.templateSourceFlag.Local {
		path := args[0]
		log.Debug().Msgf("temporary box from local template: path=%s", path)

		sourceLoader := template.NewLocalCachedLoader[boxModel.BoxV1](path, opts.configRef.Config.Template.CacheDir)
		return opts.temporaryBox(sourceLoader, boxModel.NewBoxLabels().AddDefaultLocal())

	} else {
		name := args[0]
		log.Debug().Msgf("temporary box from git template: name=%s revision=%s", name, opts.templateSourceFlag.Revision)

		sourceOpts := commonCmd.NewGitSourceOptions(opts.configRef.Config.Template.CacheDir, opts.templateSourceFlag.Revision)
		sourceLoader := template.NewGitLoader[boxModel.BoxV1](sourceOpts, name)
		labels := boxModel.NewBoxLabels().AddDefaultGit(sourceOpts.RepositoryUrl, sourceOpts.DefaultRevision, sourceOpts.CacheDirName())
		return opts.temporaryBox(sourceLoader, labels)
	}
}

func (opts *boxCmdOptions) temporaryBox(sourceLoader template.SourceLoader[boxModel.BoxV1], labels commonModel.Labels) error {

	temporaryClient := func(invokeOpts *invokeOptions) error {

		createOpts, err := newCreateOptions(invokeOpts.template, labels, opts.configRef.Config.Box.Size)
		if err != nil {
			return err
		}
		boxInfo, err := invokeOpts.client.Create(createOpts)
		if err != nil {
			return err
		}

		connectOpts := opts.tunnelFlag.ToConnectOptions(&invokeOpts.template.Value.Data, boxInfo.Name, true)
		return invokeOpts.client.Connect(connectOpts)
	}
	return runBoxClient(sourceLoader, opts.provider, opts.configRef, temporaryClient)
}
