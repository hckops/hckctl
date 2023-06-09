package box

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/template/source"
)

type boxCreateCmdOptions struct {
	configRef  *config.ConfigRef
	sourceFlag *common.SourceFlag
}

func NewBoxCreateCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxCreateCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "create [name]",
		Short: "create a detached box",
		RunE:  opts.run,
	}

	// --provider
	addBoxProviderFlag(command)
	// --revision or --local
	opts.sourceFlag = common.AddTemplateSourceFlag(command)

	return command
}

func (opts *boxCreateCmdOptions) run(cmd *cobra.Command, args []string) error {

	if len(args) == 1 && opts.sourceFlag.Local {
		path := args[0]
		log.Debug().Msgf("create local box: %s", path)

		return createBox(source.NewLocalSource(path), opts.configRef)

	} else if len(args) == 1 {
		name := args[0]
		revisionOpts := &source.RevisionOpts{
			SourceCacheDir: opts.configRef.Config.Template.CacheDir,
			SourceUrl:      common.TemplateSourceUrl,
			SourceRevision: common.TemplateSourceRevision,
			Revision:       opts.sourceFlag.Revision,
		}
		log.Debug().Msgf("create remote box: %s", name)

		return createBox(source.NewRemoteSource(revisionOpts, name), opts.configRef)

	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

// TODO before return error invoke stop or it will screw the terminal

func createBox(src source.TemplateSource, configRef *config.ConfigRef) error {

	boxTemplate, err := src.ReadBox()
	if err != nil {
		log.Warn().Err(err).Msg("error reading template")
		return errors.New("invalid template")
	}

	loader := common.NewLoader()
	loader.Start("loading template %s", boxTemplate.Name)

	provider := configRef.Config.Box.Provider
	boxId := boxTemplate.GenerateName(common.BoxPrefix)
	log.Debug().Msgf("creating box provider=%s name=%s boxId=%s\n%s", provider, boxTemplate.Name, boxId, boxTemplate.Pretty())

	client, err := box.NewBoxClient(provider, boxTemplate)
	if err != nil {
		log.Warn().Err(err).Msg("error creating client")
		loader.Stop()
		return errors.New("client error")
	}

	// TODO
	client.Setup()
	loader.Stop()
	fmt.Println(boxId)
	return nil
}
