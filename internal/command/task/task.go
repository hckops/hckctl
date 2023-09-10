package task

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	commonCmd "github.com/hckops/hckctl/internal/command/common"
	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/internal/command/config"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/task"
	"github.com/hckops/hckctl/pkg/task/model"
	taskModel "github.com/hckops/hckctl/pkg/task/model"
	"github.com/hckops/hckctl/pkg/template"
)

type taskCmdOptions struct {
	configRef    *config.ConfigRef
	sourceFlag   *commonFlag.SourceFlag
	providerFlag *commonFlag.ProviderFlag
	commandFlag  string // e.g. default (nothing), inline (reserved keyword), other values
}

func NewTaskCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := taskCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:     "task [name]",
		Short:   "Run a task",
		Args:    cobra.ExactArgs(1),
		PreRunE: opts.validate,
		RunE:    opts.run,
		Hidden:  false, // TODO WIP
	}

	// --revision or --local
	opts.sourceFlag = commonFlag.AddTemplateSourceFlag(command)
	// --provider (enum)
	// TODO opts.providerFlag = labFlag.AddLabProviderFlag(command)

	return command
}

func (opts *taskCmdOptions) validate(cmd *cobra.Command, args []string) error {

	// TODO validate provider

	if err := commonFlag.ValidateSourceFlag(opts.providerFlag, opts.sourceFlag); err != nil {
		log.Warn().Err(err).Msgf(commonFlag.ErrorFlagNotSupported)
		return errors.New(commonFlag.ErrorFlagNotSupported)
	}
	return nil
}

func (opts *taskCmdOptions) run(cmd *cobra.Command, args []string) error {

	if opts.sourceFlag.Local {
		path := args[0]
		log.Debug().Msgf("run task from local template: path=%s", path)

		sourceLoader := template.NewLocalCachedLoader[model.TaskV1](path, opts.configRef.Config.Template.CacheDir)
		// TODO labels
		// TODO provider
		return runTask(sourceLoader, "TODO opts.provider", opts.configRef)

	} else {
		name := args[0]
		log.Debug().Msgf("run task from git template: name=%s revision=%s", name, opts.sourceFlag.Revision)

		sourceOpts := commonCmd.NewGitSourceOptions(opts.configRef.Config.Template.CacheDir, opts.sourceFlag.Revision)
		sourceLoader := template.NewGitLoader[model.TaskV1](sourceOpts, name)
		// TODO labels
		// TODO provider
		return runTask(sourceLoader, "TODO opts.provider", opts.configRef)
	}
}

func runTask(sourceLoader template.SourceLoader[model.TaskV1], provider taskModel.TaskProvider, configRef *config.ConfigRef) error {

	loader := commonCmd.NewLoader()
	// TODO loader.Start("loading template %s", labTemplate.Value.Data.Name)
	defer loader.Stop()

	taskClient, err := newDefaultTaskClient(provider, configRef, loader)
	if err != nil {
		return err
	}

	createOpts := &taskModel.CreateOptions{
		TaskTemplate: nil,                      // TODO
		Parameters:   commonModel.Parameters{}, // TODO common model
		Labels:       commonModel.Labels{},     // TODO
	}

	return taskClient.Run(createOpts)
}

func newDefaultTaskClient(provider taskModel.TaskProvider, configRef *config.ConfigRef, loader *commonCmd.Loader) (task.TaskClient, error) {
	taskClientOpts := &taskModel.TaskClientOptions{
		Provider:   provider,
		DockerOpts: nil, // TODO
	}

	taskClient, err := task.NewTaskClient(taskClientOpts)
	if err != nil {
		log.Error().Err(err).Msgf("error task client provider=%s", provider)
		return nil, fmt.Errorf("error %s client", provider)
	}

	taskClient.Events().Subscribe(commonCmd.EventCallback(loader))
	return taskClient, nil
}
