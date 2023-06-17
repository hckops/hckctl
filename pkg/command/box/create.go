package box

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/template"
)

type boxCreateCmdOptions struct {
	configRef         *config.ConfigRef
	sourceFlag        *common.SourceFlag
	providerValueFlag *providerFlag
}

func NewBoxCreateCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxCreateCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a detached box",
		RunE:  opts.run,
	}

	// --provider
	opts.providerValueFlag = addBoxProviderFlag(command)
	// --revision or --local
	opts.sourceFlag = common.AddTemplateSourceFlag(command)

	return command
}

func (opts *boxCreateCmdOptions) run(cmd *cobra.Command, args []string) error {

	provider, err := validateProvider(opts.configRef.Config.Box.Provider, opts.providerValueFlag)
	if err != nil {
		return err
	}

	if len(args) == 1 && opts.sourceFlag.Local {
		path := args[0]
		log.Debug().Msgf("create box from local template: path=%s", path)

		return createBox(template.NewLocalSource(path), provider)

	} else if len(args) == 1 {
		name := args[0]
		revisionOpts := &template.RevisionOpts{
			SourceCacheDir: opts.configRef.Config.Template.CacheDir,
			SourceUrl:      common.TemplateSourceUrl,
			SourceRevision: common.TemplateSourceRevision,
			Revision:       opts.sourceFlag.Revision,
		}
		log.Debug().Msgf("create box from remote template: name=%s revision=%s", name, opts.sourceFlag.Revision)

		return createBox(template.NewRemoteSource(revisionOpts, name), provider)

	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func createBox(src template.TemplateSource, provider model.BoxProvider) error {
	createClient := func(opts *boxClientOpts) error {
		if boxInfo, err := opts.client.Create(opts.template); err != nil {
			return err
		} else {
			opts.loader.Stop()
			fmt.Println(boxInfo.Name)
		}
		return nil
	}
	return runBoxClient(src, provider, createClient)
}
