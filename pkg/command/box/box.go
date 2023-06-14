package box

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/template"
)

type boxCmdOptions struct {
	configRef  *config.ConfigRef
	sourceFlag *common.SourceFlag
}

func NewBoxCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "box [name]",
		Short: "attach and tunnel boxes",
		Long: heredoc.Doc(`
			attach and tunnel boxes

			  Create and attach to an ephemeral Box, tunnelling locally all the open ports.
			  All public templates are versioned under the /boxes/ sub-path on GitHub
			  at https://github.com/hckops/megalopolis

			  Independently from the provider and the template used, it will spawn a shell
			  that when closed will automatically remove and cleanup the running instance.

			  The main purpose of Boxes is to provide a ready-to-go and always up-to-date
			  hacking environment with an uniformed experience, abstracting the actual providers
			  e.g. Docker, Kubernetes, etc.

			  Boxes are a simple building block, for more advanced scenarios prefer Labs.
		`),
		Example: heredoc.Doc(`

			# creates and attaches to a "boxes/official/parrot" docker container,
			# spawns a /bin/bash shell and tunnels the following ports:
			# (vnc)			vncviewer localhost:5900
			# (novnc)		http://localhost:6080
			# (tty)			http://localhost:7681
			hckctl box parrot

			# opens a box deployed on kubernetes (docker|kube|argo|cloud)
			hckctl box kali --provider kube

			# opens a box using a specific version (branch|tag|sha)
			hckctl box vulnerable/dvwa --revision main

			# opens a box defined locally
			hckctl box ../megalopolis/boxes/official/powershell.yml --local
		`),
		RunE: opts.run,
	}

	// --provider
	addBoxProviderFlag(command)
	// --revision or --localtemplate
	opts.sourceFlag = common.AddTemplateSourceFlag(command)

	command.AddCommand(NewBoxCreateCmd(configRef))
	command.AddCommand(NewBoxDeleteCmd(configRef))
	command.AddCommand(NewBoxExecCmd(configRef))
	command.AddCommand(NewBoxListCmd(configRef))

	return command
}

func (opts *boxCmdOptions) run(cmd *cobra.Command, args []string) error {

	if len(args) == 1 && opts.sourceFlag.Local {
		path := args[0]
		log.Debug().Msgf("open box from local template: path=%s", path)

		return openBox(template.NewLocalSource(path), opts.configRef)

	} else if len(args) == 1 {
		name := args[0]
		revisionOpts := &template.RevisionOpts{
			SourceCacheDir: opts.configRef.Config.Template.CacheDir,
			SourceUrl:      common.TemplateSourceUrl,
			SourceRevision: common.TemplateSourceRevision,
			Revision:       opts.sourceFlag.Revision,
		}
		log.Debug().Msgf("open box from remote template: name=%s revision=%s", name, opts.sourceFlag.Revision)

		return openBox(template.NewRemoteSource(revisionOpts, name), opts.configRef)

	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func openBox(src template.TemplateSource, configRef *config.ConfigRef) error {
	provider := configRef.Config.Box.Provider

	openClient := func(client box.BoxClient, template *model.BoxV1) error {
		return client.Open(template)
	}
	return runBoxClient(src, provider, openClient)
}
