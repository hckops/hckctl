package box

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/client"
	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/template/model"
	"github.com/hckops/hckctl/pkg/template/source"
)

type boxExecCmdOptions struct {
	configRef *config.ConfigRef
}

func NewBoxExecCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxExecCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "exec",
		Short: "exec in a box",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxExecCmdOptions) run(cmd *cobra.Command, args []string) error {

	if len(args) == 1 {
		boxName := args[0]
		log.Debug().Msgf("exec remote box: boxName=%s", boxName)

		execClient := func(boxClient box.BoxClient, boxTemplate *model.BoxV1) error {
			return boxClient.Exec(boxName, boxTemplate.Shell)
		}
		return runBoxClient(opts.configRef, boxName, execClient)

	} else {
		cmd.HelpFunc()(cmd, args)
	}

	return nil
}

// resolve box by name
func runBoxClient(configRef *config.ConfigRef, boxName string, run func(box.BoxClient, *model.BoxV1) error) error {

	// TODO create container with Labels=revision?
	// best effort mode to resolve default template
	// WARNING this might return unexpected results if the container was created with a different revision
	revision := common.TemplateSourceRevision

	revisionOpts := &source.RevisionOpts{
		SourceCacheDir: configRef.Config.Template.CacheDir,
		SourceUrl:      common.TemplateSourceUrl,
		SourceRevision: common.TemplateSourceRevision,
		Revision:       revision,
	}

	templateName := model.ToBoxTemplateName(boxName)
	boxTemplate, err := source.NewRemoteSource(revisionOpts, templateName).ReadBox()
	if err != nil {
		log.Warn().Err(err).Msg("error reading box template")
		return errors.New("invalid template")
	}

	// TODO how to resolve provider without attempting all of them?
	boxClient, err := box.NewBoxClient(box.Docker)
	if err != nil {
		log.Warn().Err(err).Msg("error creating client")
		return errors.New("client error")
	}
	boxClient.Events().Subscribe(func(event client.Event) {
		log.Debug().Msg(event.String())
	})

	if err := run(boxClient, boxTemplate); err != nil {
		log.Warn().Err(err).Msg("error invoking client")
		return errors.New("run error")
	}

	return nil
}
