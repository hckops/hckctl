package task

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
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
	configRef      *config.ConfigRef
	sourceFlag     *commonFlag.TemplateSourceFlag
	providerFlag   *commonFlag.ProviderFlag
	commandFlag    *taskFlag.CommandFlag
	networkVpnFlag string
	provider       taskModel.TaskProvider
	parameters     commonModel.Parameters
}

func NewTaskCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := taskCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:     "task [name]",
		Short:   "Run a task",
		Long:    heredoc.Doc(`TODO long`),
		Example: heredoc.Doc(`TODO example`),
		Args:    cobra.MinimumNArgs(1),
		PreRunE: opts.validate,
		RunE:    opts.run,
	}

	// --revision or --local
	opts.sourceFlag = commonFlag.AddTemplateSourceFlag(command)
	// --network-vpn
	commonFlag.AddNetworkVpnFlag(command, &opts.networkVpnFlag)
	// --provider (enum)
	opts.providerFlag = taskFlag.AddTaskProviderFlag(command)
	// --inline or --command with N --inputs
	opts.commandFlag = taskFlag.AddCommandFlag(command)

	return command
}

func (opts *taskCmdOptions) validate(cmd *cobra.Command, args []string) error {

	// source
	if err := commonFlag.ValidateTemplateSourceFlag(opts.providerFlag, opts.sourceFlag); err != nil {
		log.Warn().Err(err).Msgf(commonFlag.ErrorFlagNotSupported)
		return errors.New(commonFlag.ErrorFlagNotSupported)
	}

	// vpn
	if err := commonFlag.ValidateNetworkVpnFlag(opts.networkVpnFlag, opts.configRef.Config.Network.VpnNetworks()); err != nil {
		return err
	}

	// provider
	if validProvider, err := taskFlag.ValidateTaskProviderFlag(opts.configRef.Config.Task.Provider, opts.providerFlag); err != nil {
		return err
	} else {
		opts.provider = validProvider
	}

	// inputs
	if validParameters, err := taskFlag.ValidateCommandInputsFlag(opts.commandFlag.Inputs); err != nil {
		return err
	} else {
		opts.parameters = validParameters
	}
	return nil
}

func (opts *taskCmdOptions) run(cmd *cobra.Command, args []string) error {

	if opts.sourceFlag.Local {
		path := args[0]
		log.Debug().Msgf("run task from local template: path=%s", path)

		sourceLoader := template.NewLocalCachedLoader[taskModel.TaskV1](path, opts.configRef.Config.Template.CacheDir)
		return opts.runTask(sourceLoader, taskModel.NewTaskLabels().AddDefaultLocal(), args[1:])

	} else {
		name := args[0]
		log.Debug().Msgf("run task from git template: name=%s revision=%s", name, opts.sourceFlag.Revision)

		sourceOpts := commonCmd.NewGitSourceOptions(opts.configRef.Config.Template.CacheDir, opts.sourceFlag.Revision)
		sourceLoader := template.NewGitLoader[taskModel.TaskV1](sourceOpts, name)
		labels := taskModel.NewTaskLabels().AddDefaultGit(sourceOpts.RepositoryUrl, sourceOpts.DefaultRevision, sourceOpts.CacheDirName())
		return opts.runTask(sourceLoader, labels, args[1:])
	}
}

func (opts *taskCmdOptions) runTask(sourceLoader template.SourceLoader[taskModel.TaskV1], labels commonModel.Labels, inlineArguments []string) error {

	info, err := sourceLoader.Read()
	if err != nil || info.Value.Kind != schema.KindTaskV1 {
		log.Warn().Err(err).Msg("error reading template")
		return errors.New("invalid template")
	}

	loader := commonCmd.NewLoader()
	loader.Start("loading template %s", info.Value.Data.Name) // TODO review template name e.g task/name (lowercase)
	defer loader.Stop()

	log.Info().Msgf("loading template: provider=%s name=%s\n%s", opts.provider, info.Value.Data.Name, info.Value.Data.Pretty())

	taskClient, err := newDefaultTaskClient(opts.provider, opts.configRef, loader)
	if err != nil {
		return err
	}

	// TODO --command and --input
	// TODO opts.networkVpnFlag

	var arguments []string
	if opts.commandFlag.Inline {
		arguments = inlineArguments
	} else {
		// TODO expand/merge values
		arguments = info.Value.Data.DefaultCommandArguments()
	}

	runOpts := &taskModel.RunOptions{
		Template:   &info.Value.Data,
		Arguments:  arguments,
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
