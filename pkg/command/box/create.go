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
		log.Debug().Msgf("create box from local template: path=%s", path)

		return createBox(source.NewLocalSource(path), opts.configRef)

	} else if len(args) == 1 {
		name := args[0]
		revisionOpts := &source.RevisionOpts{
			SourceCacheDir: opts.configRef.Config.Template.CacheDir,
			SourceUrl:      common.TemplateSourceUrl,
			SourceRevision: common.TemplateSourceRevision,
			Revision:       opts.sourceFlag.Revision,
		}
		log.Debug().Msgf("create box from remote template: name=%s revision=%s", name, opts.sourceFlag.Revision)

		return createBox(source.NewRemoteSource(revisionOpts, name), opts.configRef)

	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func createBox(src source.TemplateSource, configRef *config.ConfigRef) error {

	boxTemplate, err := src.ReadBox()
	if err != nil {
		log.Warn().Err(err).Msg("error reading template")
		return errors.New("invalid template")
	}

	loader := common.NewLoader()
	loader.Start("loading template %s", boxTemplate.Name)

	provider := configRef.Config.Box.Provider
	log.Debug().Msgf("creating box: provider=%s name=%s\n%s", provider, boxTemplate.Name, boxTemplate.Pretty())

	client, err := box.NewBoxClient(provider, boxTemplate)
	if err != nil {
		log.Warn().Err(err).Msg("error creating client")
		loader.Stop()
		return errors.New("client error")
	}

	boxId, err := client.Create()
	if err != nil {
		log.Warn().Err(err).Msg("error creating box")
		loader.Stop()
		return errors.New("create error")
	}

	loader.Stop()
	log.Info().Msgf("new box successfully created: boxId=%s", boxId)
	fmt.Println(boxId)
	return nil
}
