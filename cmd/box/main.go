package box

import (
	"log"

	"github.com/spf13/cobra"
)

func NewBoxCmd() *cobra.Command {
	var docker bool

	command := &cobra.Command{
		Use:   "box [NAME]",
		Short: "attach and tunnel a box",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				log.Println("BOX")
			} else {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	command.Flags().BoolVar(&docker, "docker", false, "start a docker container locally")
	//command.Flags().BoolVar(&docker, "podman", false, "start a podman container locally")
	//command.MarkFlagsMutuallyExclusive("docker", "podman")
	return command
}
