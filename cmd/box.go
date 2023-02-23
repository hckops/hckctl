package cmd

import (
	"fmt"

	"github.com/hckops/hckctl/internal/box"
	"github.com/hckops/hckctl/internal/template"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewBoxCmd() *cobra.Command {
	var cloud bool
	var kubernetes bool
	var docker bool

	command := &cobra.Command{
		Use:   "box [NAME]",
		Short: "attach and tunnel a box",
		Run: func(cmd *cobra.Command, args []string) {

			config := GetCliConfig().Box

			if len(args) == 1 {
				name := args[0]

				if kubernetes {
					runKubeBoxCmd(name, config)
				} else if docker {
					runDockerBoxCmd(name, config)
				} else {
					runCloudBoxCmd(name, config)
				}

			} else {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	const (
		RevisionFlag = "revision"
	)

	// TODO should this be global?
	command.PersistentFlags().StringP(RevisionFlag, "r", "main", "megalopolis git source version i.e. branch|tag|sha")
	viper.BindPFlag("box.revision", command.PersistentFlags().Lookup(RevisionFlag))

	command.Flags().BoolVar(&cloud, "cloud", true, "start a remote box")
	command.Flags().BoolVar(&kubernetes, "kube", false, "start a kubernetes box")
	command.Flags().BoolVar(&docker, "docker", false, "start a docker box")
	// TODO podman, firecracker?
	command.MarkFlagsMutuallyExclusive("cloud", "kube", "docker")

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

func runCloudBoxCmd(name string, config BoxConfig) {
	log.Debug().Msgf("request cloud box: name=%s revision=%s", name, config.Revision)
}

func runKubeBoxCmd(name string, config BoxConfig) {
	log.Debug().Msgf("request kube box: name=%s revision=%s", name, config.Revision)
}

func runDockerBoxCmd(name string, config BoxConfig) {
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
