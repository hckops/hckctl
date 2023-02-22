package cmd

import (
	"fmt"
	"log"

	b "github.com/hckops/hckctl/internal/box"
	t "github.com/hckops/hckctl/internal/template"
	l "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func NewBoxCmd() *cobra.Command {
	var revision string
	var cloud bool
	var kubernetes bool
	var docker bool

	command := &cobra.Command{
		Use:   "box [NAME]",
		Short: "attach and tunnel a box",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				name := args[0]

				if kubernetes {
					runKubeBoxCmd(name, revision)
				} else if docker {
					runDockerBoxCmd(name, revision)
				} else {
					runCloudBoxCmd(name, revision)
				}

			} else {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}

	command.Flags().StringVarP(&revision, "revision", "r", "main", "git source version i.e. branch|tag|sha")
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

func runCloudBoxCmd(name, revision string) {
	l.Info().Msg("CLOUD")
	log.Println("CLOUD")
}

func runKubeBoxCmd(name, revision string) {
	log.Println("KUBE")
}

func runDockerBoxCmd(name, revision string) {
	data, err := t.FetchTemplate(name, revision)
	if err != nil {
		log.Fatalln(err)
	}

	box, err := t.ParseValidBoxV1(data)
	if err != nil {
		log.Fatalln(err)
	}

	b.NewDockerBox(box).InitBox()
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
