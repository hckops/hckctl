package lab

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	commonCmd "github.com/hckops/hckctl/internal/command/common"
	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/internal/command/config"
	"github.com/hckops/hckctl/internal/command/version"
	boxModel "github.com/hckops/hckctl/pkg/box/model"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/lab"
	labModel "github.com/hckops/hckctl/pkg/lab/model"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/template"
)

type labCmdOptions struct {
	configRef  *config.ConfigRef
	inputsFlag []string
	parameters commonModel.Parameters
}

func NewLabCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := labCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:     "lab [name]",
		Short:   "Create a managed lab",
		Long:    heredoc.Doc(`TODO long`),
		Example: heredoc.Doc(`TODO example`),
		Args:    cobra.ExactArgs(1),
		PreRunE: opts.validate,
		RunE:    opts.run,
		Hidden:  false,
	}

	// N --inputs
	const (
		inputFlagName  = "input"
		inputFlagUsage = "override defaults"
	)
	command.Flags().StringArrayVarP(&opts.inputsFlag, inputFlagName, commonFlag.NoneFlagShortHand, []string{}, inputFlagUsage)

	return command
}

func (opts *labCmdOptions) validate(cmd *cobra.Command, args []string) error {
	// inputs
	if validParameters, err := commonFlag.ValidateParametersFlag(opts.inputsFlag); err != nil {
		return err
	} else {
		opts.parameters = validParameters
	}
	return nil
}

func (opts *labCmdOptions) run(cmd *cobra.Command, args []string) error {
	name := args[0]
	revision := commonCmd.TemplateSourceRevision
	log.Debug().Msgf("create lab from git template: name=%s revision=%s", name, revision)

	sourceOpts := commonCmd.NewGitSourceOptions(opts.configRef.Config.Template.CacheDir, revision)
	sourceLoader := template.NewGitLoader[labModel.LabV1](sourceOpts, name)
	return startLab(sourceLoader, opts.configRef, opts.parameters)
}

func startLab(sourceLoader template.SourceLoader[labModel.LabV1], configRef *config.ConfigRef, parameters commonModel.Parameters) error {

	info, err := sourceLoader.Read()
	if err != nil || info.Value.Kind != schema.KindLabV1 {
		log.Warn().Err(err).Msg("error reading template")
		return errors.New("invalid template")
	}

	loader := commonCmd.NewLoader()
	loader.Start("loading template %s", info.Value.Data.Name)
	defer loader.Stop()

	log.Info().Msgf("loading template: name=%s\n%s", info.Value.Data.Name, info.Value.Data.Pretty())

	labClient, err := newDefaultLabClient(configRef, loader)
	if err != nil {
		return err
	}

	createOpts := &labModel.CreateOptions{
		LabTemplate:   &info.Value.Data,
		BoxTemplates:  map[string]*boxModel.BoxV1{},  // cloud only
		DumpTemplates: map[string]*labModel.DumpV1{}, // cloud only
		Parameters:    parameters,
		Labels:        commonModel.Labels{}, // cloud only
	}

	if labInfo, err := labClient.Create(createOpts); err != nil {
		return err
	} else {
		loader.Stop()
		fmt.Println(labInfo.Name)
	}
	return nil
}

func newDefaultLabClient(configRef *config.ConfigRef, loader *commonCmd.Loader) (lab.LabClient, error) {
	provider := labModel.Cloud
	labClientOpts := &labModel.LabClientOptions{
		Provider:  provider,
		CloudOpts: configRef.Config.Provider.Cloud.ToCloudOptions(version.ClientVersion()),
	}

	labClient, err := lab.NewLabClient(labClientOpts)
	if err != nil {
		log.Error().Err(err).Msgf("error lab client provider=%s", provider)
		return nil, fmt.Errorf("error %s client", provider)
	}

	labClient.Events().Subscribe(commonCmd.EventCallback(loader))
	return labClient, nil
}
