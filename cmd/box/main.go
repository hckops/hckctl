package box

import (
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
					RunKubeBoxCmd(name, revision)
				} else if docker {
					RunDockerBoxCmd(name, revision)
				} else {
					RunCloudBoxCmd(name, revision)
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
			RunBoxListCmd()
		},
	}

	command.AddCommand(listCmd)
	return command
}
