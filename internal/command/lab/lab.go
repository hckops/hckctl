package lab

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	commonCmd "github.com/hckops/hckctl/internal/command/common"
	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/internal/command/config"
	labFlag "github.com/hckops/hckctl/internal/command/lab/flag"
	"github.com/hckops/hckctl/internal/command/version"
	boxModel "github.com/hckops/hckctl/pkg/box/model"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/lab"
	labModel "github.com/hckops/hckctl/pkg/lab/model"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/template"
)

type labCmdOptions struct {
	configRef    *config.ConfigRef
	sourceFlag   *commonFlag.SourceFlag
	providerFlag *commonFlag.ProviderFlag
	provider     labModel.LabProvider
}

func NewLabCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := labCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:     "lab [name]",
		Short:   "Create a lab",
		Args:    cobra.ExactArgs(1),
		PreRunE: opts.validate,
		RunE:    opts.run,
		Hidden:  false, // TODO WIP
	}

	// --revision or --local
	opts.sourceFlag = commonFlag.AddTemplateSourceFlag(command)
	// --provider (enum)
	opts.providerFlag = labFlag.AddLabProviderFlag(command)

	return command
}

func (opts *labCmdOptions) validate(cmd *cobra.Command, args []string) error {

	validProvider, err := labFlag.ValidateLabProvider(opts.configRef.Config.Lab.Provider, opts.providerFlag)
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

func (opts *labCmdOptions) run(cmd *cobra.Command, args []string) error {

	if opts.sourceFlag.Local {
		path := args[0]
		log.Debug().Msgf("create lab from local template: path=%s", path)

		sourceLoader := template.NewLocalCachedLoader[labModel.LabV1](path, opts.configRef.Config.Template.CacheDir)
		// TODO labels
		return startLab(sourceLoader, opts.provider, opts.configRef)

	} else {
		name := args[0]
		log.Debug().Msgf("create lab from git template: name=%s revision=%s", name, opts.sourceFlag.Revision)

		sourceOpts := commonCmd.NewGitSourceOptions(opts.configRef.Config.Template.CacheDir, opts.sourceFlag.Revision)
		sourceLoader := template.NewGitLoader[labModel.LabV1](sourceOpts, name)
		// TODO labels
		return startLab(sourceLoader, opts.provider, opts.configRef)
	}
}

func startLab(sourceLoader template.SourceLoader[labModel.LabV1], provider labModel.LabProvider, configRef *config.ConfigRef) error {

	labTemplate, err := sourceLoader.Read()
	if err != nil || labTemplate.Value.Kind != schema.KindLabV1 {
		log.Warn().Err(err).Msg("error reading template")
		return errors.New("invalid template")
	}

	loader := commonCmd.NewLoader()
	loader.Start("loading template %s", labTemplate.Value.Data.Name)
	defer loader.Stop()

	log.Info().Msgf("loading template: provider=%s name=%s\n%s", provider, labTemplate.Value.Data.Name, labTemplate.Value.Data.Pretty())

	labClient, err := newDefaultLabClient(provider, configRef, loader)
	if err != nil {
		return err
	}

	createOpts := &labModel.CreateOptions{
		LabTemplate:   &labTemplate.Value.Data,
		BoxTemplates:  map[string]*boxModel.BoxV1{}, // TODO load box template
		DumpTemplates: map[string]*labModel.DumpV1{},
		Parameters:    map[string]string{},  // TODO add override --input alias=parrot --input password=changeme --input vpn=htb-eu
		Labels:        commonModel.Labels{}, // TODO box+lab labels
	}

	if labInfo, err := labClient.Create(createOpts); err != nil {
		return err
	} else {
		loader.Stop()
		fmt.Println(labInfo.Name)
	}
	return nil
}

func newDefaultLabClient(provider labModel.LabProvider, configRef *config.ConfigRef, loader *commonCmd.Loader) (lab.LabClient, error) {
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
