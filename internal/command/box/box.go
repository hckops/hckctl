package box

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	boxFlag "github.com/hckops/hckctl/internal/command/box/flag"
	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/internal/command/config"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/template"
)

type boxCmdOptions struct {
	configRef    *config.ConfigRef
	sourceFlag   *commonFlag.SourceFlag
	providerFlag *commonFlag.ProviderFlag
	tunnelFlag   *boxFlag.TunnelFlag
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
		RunE: opts.run,
	}

	// --revision or --local
	opts.sourceFlag = commonFlag.AddTemplateSourceFlag(command)
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

func (opts *boxCmdOptions) run(cmd *cobra.Command, args []string) error {

	provider, err := boxFlag.ValidateBoxProvider(opts.configRef.Config.Box.Provider, opts.providerFlag)
	if err != nil {
		return err
	} else if len(args) == 1 {

		if err := opts.validateFlags(provider); err != nil {
			log.Warn().Err(err).Msgf(commonFlag.ErrorFlagNotSupported)
			return errors.New(commonFlag.ErrorFlagNotSupported)

		} else if opts.sourceFlag.Local {
			path := args[0]
			log.Debug().Msgf("temporary box from local template: path=%s", path)

			sourceLoader := template.NewLocalCachedLoader[model.BoxV1](path, opts.configRef.Config.Template.CacheDir)
			return opts.temporaryBox(sourceLoader, provider, model.NewLocalLabels())

		} else {
			name := args[0]
			log.Debug().Msgf("temporary box from git template: name=%s revision=%s", name, opts.sourceFlag.Revision)

			sourceOpts := newGitSourceOptions(opts.configRef.Config.Template.CacheDir, opts.sourceFlag.Revision)
			sourceLoader := template.NewGitLoader[model.BoxV1](sourceOpts, name)
			labels := model.NewGitLabels(sourceOpts.RepositoryUrl, sourceOpts.DefaultRevision, sourceOpts.CacheDirName())
			return opts.temporaryBox(sourceLoader, provider, labels)
		}

	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func (opts *boxCmdOptions) validateFlags(provider model.BoxProvider) error {
	if err := boxFlag.ValidateSourceFlag(provider, opts.sourceFlag); err != nil {
		log.Warn().Err(err).Msgf(commonFlag.ErrorFlagNotSupported)
		return err
	}
	if err := boxFlag.ValidateTunnelFlag(provider, opts.tunnelFlag); err != nil {
		log.Warn().Err(err).Msgf("ignore validation %s", commonFlag.ErrorFlagNotSupported)
		return nil
	}
	return nil
}

func (opts *boxCmdOptions) temporaryBox(sourceLoader template.SourceLoader[model.BoxV1], provider model.BoxProvider, labels model.BoxLabels) error {

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
	return runBoxClient(sourceLoader, provider, opts.configRef, temporaryClient)
}