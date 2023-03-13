package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag/v2"

	"github.com/hckops/hckctl/internal/box"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/template"
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
					box.NewDockerBox(template).Open()
				case KubernetesFlag:
					box.NewKubeBox(template, &config.Kube).Open()
				case CloudFlag:
					box.NewCloudBox(template, &config.Cloud).Open()
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

func requestBoxTemplate(name string, revision string) *schema.BoxV1 {
	log.Info().Msgf("request box template: name=%s revision=%s", name, revision)

	rawTemplate, err := template.RequestTemplate(NewBoxParam(name, revision))
	if err != nil {
		printFatalError(err, "unable to fetch box template")
	}

	boxTemplate, err := schema.ParseValidBoxV1(rawTemplate)
	if err != nil {
		printFatalError(err, "invalid box template")
	}

	return boxTemplate
}

// TODO shared with template cmd
func NewBoxParam(name, revision string) *template.TemplateParam {
	return &template.TemplateParam{
		TemplateKind:  "box/v1", // TODO enum
		TemplateName:  name,
		Revision:      revision,
		ClientVersion: "hckctl-v0.0.0", // TODO sha/tag
	}
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
