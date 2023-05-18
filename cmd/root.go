package cmd

import (
	"github.com/spf13/cobra"

	boxCmd "github.com/hckops/hckctl/cmd/box"
	"github.com/hckops/hckctl/cmd/common"
	configCmd "github.com/hckops/hckctl/cmd/config"
	templateCmd "github.com/hckops/hckctl/cmd/template"
)

func NewRoodCmd() *cobra.Command {

	opts := &common.GlobalCmdOptions{}

	rootCmd := &cobra.Command{
		Use:   "hckctl",
		Short: "The Cloud Native HaCKing Tool",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}

	// TODO --log-level

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	rootCmd.AddCommand(boxCmd.NewBoxCmd(opts))
	rootCmd.AddCommand(configCmd.NewConfigCmd())
	rootCmd.AddCommand(templateCmd.NewTemplateCmd(opts))
	return rootCmd
}
