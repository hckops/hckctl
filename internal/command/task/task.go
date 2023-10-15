package task

import (
	"fmt"
	"strings"

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
	configRef *config.ConfigRef
	// flags
	commandFlag        *taskFlag.CommandFlag
	networkVpnFlag     string
	providerFlag       *commonFlag.ProviderFlag
	templateSourceFlag *commonFlag.TemplateSourceFlag
	// internal
	provider   taskModel.TaskProvider
	parameters commonModel.Parameters
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

	// --inline or --command with N --inputs
	opts.commandFlag = taskFlag.AddCommandFlag(command)
	// --network-vpn
	commonFlag.AddNetworkVpnFlag(command, &opts.networkVpnFlag)
	// --provider (enum)
	opts.providerFlag = taskFlag.AddTaskProviderFlag(command)
	// --revision or --local
	opts.templateSourceFlag = commonFlag.AddTemplateSourceFlag(command)

	return command
}

func (opts *taskCmdOptions) validate(cmd *cobra.Command, args []string) error {
	// command
	if validParameters, err := commonFlag.ValidateParametersFlag(opts.commandFlag.Inputs); err != nil {
		return err
	} else {
		opts.parameters = validParameters
	}
	// provider
	if validProvider, err := taskFlag.ValidateTaskProviderFlag(opts.configRef.Config.Task.Provider, opts.providerFlag); err != nil {
		return err
	} else {
		opts.provider = validProvider
	}
	// network-vpn (after provider validation)
	if vpnNetworkInfo, err := commonFlag.ValidateNetworkVpnFlag(opts.networkVpnFlag, opts.configRef.Config.Network.VpnNetworks()); err != nil {
		return err
	} else if vpnNetworkInfo != nil && opts.provider == taskModel.Cloud {
		return fmt.Errorf("%s: use flow", commonFlag.ErrorFlagNotSupported)
	}
	// source
	if err := commonFlag.ValidateTemplateSourceFlag(opts.providerFlag, opts.templateSourceFlag); err != nil {
		log.Warn().Err(err).Msgf(commonFlag.ErrorFlagNotSupported)
		return errors.New(commonFlag.ErrorFlagNotSupported)
	}
	return nil
}

func (opts *taskCmdOptions) run(cmd *cobra.Command, args []string) error {

	if opts.templateSourceFlag.Local {
		path := args[0]
		log.Debug().Msgf("run task from local template: path=%s", path)

		sourceLoader := template.NewLocalCachedLoader[taskModel.TaskV1](path, opts.configRef.Config.Template.CacheDir)
		return opts.runTask(sourceLoader, taskModel.NewTaskLabels().AddDefaultLocal(), args[1:])

	} else {
		name := args[0]
		log.Debug().Msgf("run task from git template: name=%s revision=%s", name, opts.templateSourceFlag.Revision)

		sourceOpts := commonCmd.NewGitSourceOptions(opts.configRef.Config.Template.CacheDir, opts.templateSourceFlag.Revision)
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

	templateName := commonCmd.PrettyName(info, opts.configRef.Config.Template.CacheDir, info.Value.Data.Name)
	loader := commonCmd.NewLoader()
	loader.Start("loading template %s", templateName)
	defer loader.Stop()

	log.Info().Msgf("loading template: provider=%s name=%s\n%s", opts.provider, templateName, info.Value.Data.Pretty())

	taskClient, err := newDefaultTaskClient(opts.provider, opts.configRef, loader)
	if err != nil {
		return err
	}

	var arguments []string
	if opts.commandFlag.Inline {
		log.Info().Msgf("run task inline arguments=[%s]", strings.Join(inlineArguments, ","))

		arguments = inlineArguments
	} else {
		taskCommand, err := info.Value.Data.LoadCommand(opts.commandFlag.Preset)
		if err != nil {
			log.Warn().Err(err).Msg("error loading command")
			return errors.New("invalid command")
		}
		expandedArguments, err := taskCommand.ExpandCommandArguments(opts.parameters)
		if err != nil {
			log.Warn().Err(err).Msg("error expanding command arguments")
			return errors.New("invalid command arguments")
		}
		log.Info().Msgf("run task command=%s arguments=[%s] inputs=%v expanded=[%s]",
			taskCommand.Name, strings.Join(taskCommand.Arguments, ","), opts.parameters, strings.Join(expandedArguments, ","))

		arguments = expandedArguments
	}

	var networkVpn *commonModel.NetworkVpnInfo
	if networkVpnInfo, err := opts.configRef.Config.Network.ToNetworkVpnInfo(opts.networkVpnFlag); err != nil {
		log.Warn().Err(err).Msg("error invalid vpn config")
		return err
	} else if networkVpnInfo != nil {
		log.Info().Msgf("run task connected to vpn network name=%s path=%s", networkVpnInfo.Name, networkVpnInfo.LocalPath)
		networkVpn = networkVpnInfo
	}

	runOpts := &taskModel.RunOptions{
		Template: &info.Value.Data,
		Labels:   commonCmd.AddTemplateLabels[taskModel.TaskV1](info, labels),
		CommonInfo: commonModel.CommonInfo{
			NetworkVpn: networkVpn,
			ShareDir:   opts.configRef.Config.Common.ToShareDirInfo(),
		},
		StreamOpts: commonModel.NewStdStreamOpts(false),
		Arguments:  arguments,
		LogDir:     opts.configRef.Config.Task.LogDir,
	}

	if err := taskClient.Run(runOpts); err != nil {
		log.Warn().Err(err).Msg("error run task")
		return errors.New("error run task")
	}
	return nil
}

func newDefaultTaskClient(provider taskModel.TaskProvider, configRef *config.ConfigRef, loader *commonCmd.Loader) (task.TaskClient, error) {
	taskClientOpts := &taskModel.TaskClientOptions{
		Provider:   provider,
		DockerOpts: configRef.Config.Provider.Docker.ToDockerOptions(),
		KubeOpts:   configRef.Config.Provider.Kube.ToKubeOptions(),
	}

	taskClient, err := task.NewTaskClient(taskClientOpts)
	if err != nil {
		log.Error().Err(err).Msgf("error task client provider=%s", provider)
		return nil, fmt.Errorf("error %s client", provider)
	}

	taskClient.Events().Subscribe(commonCmd.EventCallback(loader))
	return taskClient, nil
}
