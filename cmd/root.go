package cmd

import (
	"fmt"
	"github.com/spf13/cobra"

	boxCmd "github.com/hckops/hckctl/cmd/box"
	commonCmd "github.com/hckops/hckctl/cmd/common"
	configCmd "github.com/hckops/hckctl/cmd/config"
	labCmd "github.com/hckops/hckctl/cmd/lab"
	templateCmd "github.com/hckops/hckctl/cmd/template"
)

func NewRoodCmd() *cobra.Command {

	opts := &commonCmd.GlobalCmdOptions{}

	description := fmt.Sprintf("The Cloud Native HaCKing Tool - %s", commonCmd.Version())

	rootCmd := &cobra.Command{
		Use:   "hckctl",
		Short: description,
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
	rootCmd.AddCommand(configCmd.NewConfigCmd(opts))
	rootCmd.AddCommand(labCmd.NewLabCmd(opts))
	rootCmd.AddCommand(templateCmd.NewTemplateCmd(opts))
	rootCmd.AddCommand(commonCmd.NewVersionCmd())
	return rootCmd
}
