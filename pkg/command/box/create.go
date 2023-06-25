package box

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/box/model"
	boxFlag "github.com/hckops/hckctl/pkg/command/box/flag"
	"github.com/hckops/hckctl/pkg/command/common"
	commonFlag "github.com/hckops/hckctl/pkg/command/common/flag"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/template"
)

type boxCreateCmdOptions struct {
	configRef    *config.ConfigRef
	sourceFlag   *commonFlag.SourceFlag
	providerFlag *commonFlag.ProviderFlag
}

func NewBoxCreateCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxCreateCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a detached long running box",
		RunE:  opts.run,
	}

	// --provider (enum)
	opts.providerFlag = boxFlag.AddBoxProviderFlag(command)
	// --revision or --local
	opts.sourceFlag = commonFlag.AddTemplateSourceFlag(command)

	return command
}

func (opts *boxCreateCmdOptions) run(cmd *cobra.Command, args []string) error {

	provider, err := boxFlag.ValidateBoxProvider(opts.configRef.Config.Box.Provider, opts.providerFlag)
	if err != nil {
		return err
	} else if len(args) == 1 && opts.sourceFlag.Local {
		path := args[0]
		log.Debug().Msgf("create box from local template: path=%s", path)

		return createBox(template.NewLocalSource(path), provider, opts.configRef)

	} else if len(args) == 1 {
		name := args[0]

		if provider == model.Cloud && opts.sourceFlag.Revision != common.TemplateSourceRevision {
			log.Warn().Msgf("revision flag not supported by cloud provider: revision=%s", opts.sourceFlag.Revision)
			return errors.New("invalid revision")
		}

		revisionOpts := &template.RevisionOpts{
			SourceCacheDir: opts.configRef.Config.Template.CacheDir,
			SourceUrl:      common.TemplateSourceUrl,
			SourceRevision: common.TemplateSourceRevision,
			Revision:       opts.sourceFlag.Revision,
		}
		log.Debug().Msgf("create box from remote template: name=%s revision=%s", name, opts.sourceFlag.Revision)

		return createBox(template.NewRemoteSource(revisionOpts, name), provider, opts.configRef)

	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func createBox(src template.TemplateSource, provider model.BoxProvider, configRef *config.ConfigRef) error {
	createClient := func(opts *boxClientOpts) error {
		if boxInfo, err := opts.client.Create(opts.template); err != nil {
			return err
		} else {
			opts.loader.Stop()
			fmt.Println(boxInfo.Name)
		}
		return nil
	}
	return runBoxClient(src, provider, configRef, createClient)
}
