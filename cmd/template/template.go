package template

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/cmd/common"
)

type templateCmdOptions struct {
	global   *common.GlobalCmdOptions
	path     string
	revision string
}

func NewTemplateCmd(global *common.GlobalCmdOptions) *cobra.Command {

	opts := templateCmdOptions{
		global: global,
	}

	command := &cobra.Command{
		Use:   "template [NAME]",
		Short: "load and validate a template",
		RunE:  opts.run,
	}

	// TODO prefix file:// for local or https:// vs separate flag e.g. localPath and remotePath + revision
	command.Flags().StringVarP(&opts.path, "path", "p", "", "load a template from a local path")
	command.Flags().StringVarP(&opts.revision, "revision", "r", "main", "megalopolis git source version i.e. branch|tag|sha")
	command.MarkFlagsMutuallyExclusive("revision", "path")

	return command
}

func (opts *templateCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
