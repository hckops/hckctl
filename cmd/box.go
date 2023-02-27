package cmd

import (
	"fmt"

	"github.com/hckops/hckctl/internal/box"
	"github.com/hckops/hckctl/internal/model"
	"github.com/hckops/hckctl/internal/template"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag/v2"
)

func NewBoxCmd() *cobra.Command {
	var provider ProviderFlag

	command := &cobra.Command{
		Use:   "box [NAME]",
		Short: "attach and tunnel a box",
		Run: func(cmd *cobra.Command, args []string) {

			config := GetCliConfig().Box
			// use value from merged config
			var provider, err = ProviderToFlag(config.Provider)
			if err != nil {
				cmd.HelpFunc()(cmd, args)
				log.Fatal().Err(err).Msgf("invalid provider: %s", config.Provider)
			}

			if len(args) == 1 {
				name := args[0]

				switch provider {
				case DockerFlag:
					runDockerBoxCmd(name, config)
				case KubernetesFlag:
					runKubeBoxCmd(name, config)
				case CloudFlag:
					runCloudBoxCmd(name, config)
				}

			} else {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	const (
		RevisionFlag = "revision"
		ProviderFlag = "provider"
	)

	// TODO should this be global?
	command.PersistentFlags().StringP(RevisionFlag, "r", "main", "megalopolis git source version i.e. branch|tag|sha")
	viper.BindPFlag("box.revision", command.PersistentFlags().Lookup(RevisionFlag))

	providerValue := enumflag.New(&provider, ProviderFlag, ProviderIds, enumflag.EnumCaseInsensitive)
	command.Flags().Var(providerValue, ProviderFlag, "set the box provider, one of docker|kube|cloud")
	viper.BindPFlag("box.provider", command.Flags().Lookup(ProviderFlag))

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list available boxes",
		Run: func(cmd *cobra.Command, args []string) {
			runBoxListCmd()
		},
	}

	command.AddCommand(listCmd)
	return command
}

func runDockerBoxCmd(name string, config model.BoxConfig) {
	log.Debug().Msgf("request docker box: name=%s revision=%s", name, config.Revision)

	rawTemplate, err := template.FetchTemplate(name, config.Revision)
	if err != nil {
		log.Fatal().Err(err).Msg("fetch box template")
	}

	boxTemplate, err := template.ParseValidBoxV1(rawTemplate)
	if err != nil {
		log.Fatal().Err(err).Msg("validate box template")
	}

	box.NewDockerBox(boxTemplate).InitBox()
}

func runKubeBoxCmd(name string, config model.BoxConfig) {
	log.Debug().Msgf("request kube box: name=%s revision=%s", name, config.Revision)

	rawTemplate, err := template.FetchTemplate(name, config.Revision)
	if err != nil {
		log.Fatal().Err(err).Msg("fetch box template")
	}

	boxTemplate, err := template.ParseValidBoxV1(rawTemplate)
	if err != nil {
		log.Fatal().Err(err).Msg("validate box template")
	}

	box.NewKubeBox(boxTemplate).InitBox(config.Kube)
}

func runCloudBoxCmd(name string, config model.BoxConfig) {
	log.Debug().Msgf("request cloud box: name=%s revision=%s", name, config.Revision)
}

func runBoxListCmd() {
	for _, box := range getBoxes() {
		fmt.Println(box)
	}
}

// TODO api endpoint or fetch from https://github.com/hckops/megalopolis/tree/main/boxes
// TODO see github api
func getBoxes() []string {
	// TODO struct: name, alias e.g. alpine -> official/alpine
	// TODO revision param

	return []string{
		"alpine",
		"parrot",
	}
}
