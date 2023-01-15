package box

import (
	"github.com/spf13/cobra"
)

func NewBoxCmd() *cobra.Command {
	var revision string
	var docker bool

	command := &cobra.Command{
		Use:   "box [NAME]",
		Short: "attach and tunnel a box",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				name := args[0]

				if docker {
					RunBoxDockerCmd(name, revision)
				} else {
					RunBoxCloudCmd(name)
				}

			} else {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}

	command.Flags().StringVarP(&revision, "revision", "r", "main", "git source version i.e. branch|tag|sha")
	command.Flags().BoolVar(&docker, "docker", false, "start a docker container locally")
	//command.Flags().BoolVar(&docker, "podman", false, "start a podman container locally")
	//command.MarkFlagsMutuallyExclusive("docker", "podman")

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
