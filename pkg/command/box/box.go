package box

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/template/source"
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
			hckctl box kali --revision kube

			# opens a box using a specific version (branch|tag|sha)
			hckctl box vulnerable/dvwa --revision main

			# opens a box defined locally
			hckctl box ../megalopolis/boxes/official/powershell.yml --local
		`),
		RunE: opts.run,
	}

	const (
		providerFlagName = "provider"
	)
	// --provider
	command.Flags().StringP(providerFlagName, common.NoneFlagShortHand, string(box.Docker),
		fmt.Sprintf("switch box provider, one of %s", strings.Join(box.BoxProviderValues(), "|")))
	viper.BindPFlag(fmt.Sprintf("box.%s", providerFlagName), command.Flags().Lookup(providerFlagName))

	// --revision or --local
	opts.sourceFlag = common.AddTemplateSourceFlag(command)

	//command.AddCommand(NewBoxCopyCmd(opts))
	command.AddCommand(NewBoxCreateCmd(configRef))
	//command.AddCommand(NewBoxDeleteCmd(opts))
	//command.AddCommand(NewBoxExecCmd(opts))
	command.AddCommand(NewBoxListCmd(configRef))
	//command.AddCommand(NewBoxTunnelCmd(opts))

	return command
}

func (opts *boxCmdOptions) run(cmd *cobra.Command, args []string) error {
	provider := opts.configRef.Config.Box.Provider

	if len(args) == 1 && opts.sourceFlag.Local {
		path := args[0]
		log.Debug().Msgf("open local box: %s", path)

		return openBox(source.NewLocalSource(path), provider)

	} else if len(args) == 1 {
		name := args[0]
		revisionOpts := &source.RevisionOpts{
			SourceCacheDir: opts.configRef.Config.Template.CacheDir,
			SourceUrl:      common.TemplateSourceUrl,
			SourceRevision: common.TemplateSourceRevision,
			Revision:       opts.sourceFlag.Revision,
		}
		log.Debug().Msgf("open remote box: %s", name)

		return openBox(source.NewRemoteSource(revisionOpts, name), provider)

	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func openBox(src source.TemplateSource, provider box.BoxProvider) error {
	log.Debug().Msg("TODO")
	loader := common.NewLoader()
	loader.Start("TODO")

	boxTemplate, err := src.ReadBox()
	if err != nil {
		log.Warn().Err(err).Msg("error reading template")
		return errors.New("invalid template")
	}

	loader.Sleep(2)
	loader.Refresh("update")
	loader.Sleep(2)

	if client, err := box.NewBoxClient(provider, boxTemplate); err != nil {
		log.Warn().Err(err).Msg("error creating client")
		return errors.New("client error")
	} else {
		// TODO
		client.Setup()
	}

	loader.Stop()
	return nil
}
