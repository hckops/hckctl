package box

import (
	"fmt"

	"github.com/pkg/errors"
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
		Short: "Create a long running detached box",
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
	} else if len(args) == 1 {

		if err := boxFlag.ValidateSourceFlag(provider, opts.sourceFlag); err != nil {
			log.Warn().Err(err).Msgf(commonFlag.ErrorFlagNotSupported)
			return errors.New(commonFlag.ErrorFlagNotSupported)

		} else if opts.sourceFlag.Local {
			path := args[0]
			log.Debug().Msgf("create box from local template: path=%s", path)

			return createBox(template.NewLocalSource(path), provider, opts.configRef)

		} else {
			name := args[0]
			sourceOpts := &template.SourceOptions{
				SourceCacheDir: opts.configRef.Config.Template.CacheDir,
				SourceUrl:      common.TemplateSourceUrl,
				SourceRevision: common.TemplateSourceRevision,
				Revision:       opts.sourceFlag.Revision,
			}
			log.Debug().Msgf("create box from git template: name=%s revision=%s", name, opts.sourceFlag.Revision)

			return createBox(template.NewGitSource(sourceOpts, name), provider, opts.configRef)
		}
	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func createBox(src template.TemplateSource, provider model.BoxProvider, configRef *config.ConfigRef) error {
	createClient := func(opts *boxClientOptions) error {
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
