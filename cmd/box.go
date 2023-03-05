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
				template := requestBoxTemplate(name, config.Revision)

				switch provider {
				case DockerFlag:
					box.NewDockerBox(template).OpenBox()
				case KubernetesFlag:
					box.NewKubeBox(template, &config.Kube).OpenBox()
				case CloudFlag:
					log.Debug().Msg("init cloud box")
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

func requestBoxTemplate(name string, revision string) *model.BoxV1 {
	log.Debug().Msgf("request box template: name=%s revision=%s", name, revision)

	rawTemplate, err := template.FetchTemplate(name, revision)
	if err != nil {
		log.Fatal().Err(err).Msg("fetch box template")
	}

	boxTemplate, err := template.ParseValidBoxV1(rawTemplate)
	if err != nil {
		log.Fatal().Err(err).Msg("validate box template")
	}

	return boxTemplate
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
