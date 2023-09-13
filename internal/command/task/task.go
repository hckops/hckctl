package task

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	commonCmd "github.com/hckops/hckctl/internal/command/common"
	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/internal/command/config"
	taskFlag "github.com/hckops/hckctl/internal/command/task/flag"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/task"
	taskModel "github.com/hckops/hckctl/pkg/task/model"
	"github.com/hckops/hckctl/pkg/template"
)

type taskCmdOptions struct {
	configRef    *config.ConfigRef
	sourceFlag   *commonFlag.SourceFlag
	providerFlag *commonFlag.ProviderFlag
	provider     taskModel.TaskProvider
	inlineFlag   bool
}

func NewTaskCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := taskCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:     "task [name]",
		Short:   "Run a task",
		PreRunE: opts.validate,
		RunE:    opts.run,
	}

	// --revision or --local
	opts.sourceFlag = commonFlag.AddTemplateSourceFlag(command)
	// --provider (enum)
	opts.providerFlag = taskFlag.AddTaskProviderFlag(command)

	const (
		inlineFlagName  = "inline"
		inlineFlagUsage = "inline arguments"
	)
	command.Flags().BoolVarP(&opts.inlineFlag, inlineFlagName, commonFlag.NoneFlagShortHand, false, inlineFlagUsage)

	return command
}

func (opts *taskCmdOptions) validate(cmd *cobra.Command, args []string) error {

	validProvider, err := taskFlag.ValidateTaskProvider(opts.configRef.Config.Task.Provider, opts.providerFlag)
	if err != nil {
		return err
	}
	opts.provider = validProvider

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

		sourceLoader := template.NewLocalCachedLoader[taskModel.TaskV1](path, opts.configRef.Config.Template.CacheDir)
		return runTask(sourceLoader, opts.provider, opts.configRef, taskModel.NewTaskLabels().AddDefaultLocal())

	} else {
		name := args[0]
		log.Debug().Msgf("run task from git template: name=%s revision=%s", name, opts.sourceFlag.Revision)

		sourceOpts := commonCmd.NewGitSourceOptions(opts.configRef.Config.Template.CacheDir, opts.sourceFlag.Revision)
		sourceLoader := template.NewGitLoader[taskModel.TaskV1](sourceOpts, name)
		labels := taskModel.NewTaskLabels().AddDefaultGit(sourceOpts.RepositoryUrl, sourceOpts.DefaultRevision, sourceOpts.CacheDirName())
		return runTask(sourceLoader, opts.provider, opts.configRef, labels)
	}
}

func runTask(sourceLoader template.SourceLoader[taskModel.TaskV1], provider taskModel.TaskProvider, configRef *config.ConfigRef, labels commonModel.Labels) error {

	info, err := sourceLoader.Read()
	if err != nil || info.Value.Kind != schema.KindTaskV1 {
		log.Warn().Err(err).Msg("error reading template")
		return errors.New("invalid template")
	}

	loader := commonCmd.NewLoader()
	loader.Start("loading template %s", info.Value.Data.Name) // TODO review template name e.g task/name (lowercase)
	defer loader.Stop()

	taskClient, err := newDefaultTaskClient(provider, configRef, loader)
	if err != nil {
		return err
	}

	runOpts := &taskModel.RunOptions{
		Template:   &info.Value.Data,
		Arguments:  info.Value.Data.DefaultCommandArgs(), // TODO expand/merge values
		Labels:     commonCmd.AddTemplateLabels[taskModel.TaskV1](info, labels),
		StreamOpts: commonModel.NewStdStreamOpts(false),
	}

	return taskClient.Run(runOpts)
}

func newDefaultTaskClient(provider taskModel.TaskProvider, configRef *config.ConfigRef, loader *commonCmd.Loader) (task.TaskClient, error) {
	taskClientOpts := &taskModel.TaskClientOptions{
		Provider:   provider,
		DockerOpts: configRef.Config.Provider.Docker.ToDockerOptions(),
	}

	taskClient, err := task.NewTaskClient(taskClientOpts)
	if err != nil {
		log.Error().Err(err).Msgf("error task client provider=%s", provider)
		return nil, fmt.Errorf("error %s client", provider)
	}

	taskClient.Events().Subscribe(commonCmd.EventCallback(loader))
	return taskClient, nil
}
