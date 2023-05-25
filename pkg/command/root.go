package command

import (
	"fmt"

	"github.com/spf13/cobra"

	boxCmd "github.com/hckops/hckctl/pkg/command/box"
	"github.com/hckops/hckctl/pkg/command/common"
	configCmd "github.com/hckops/hckctl/pkg/command/config"
	labCmd "github.com/hckops/hckctl/pkg/command/lab"
	templateCmd "github.com/hckops/hckctl/pkg/command/template"
)

func NewRoodCmd() *cobra.Command {

	opts := &common.GlobalCmdOptions{}

	description := fmt.Sprintf("The Cloud Native HaCKing Tool - %s", Version())

	rootCmd := &cobra.Command{
		Use:   "hckctl",
		Short: description,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			// TODO init config
			// TODO init logger
			common.InitFileLogger(opts)

			return nil
		},
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
	rootCmd.AddCommand(NewVersionCmd())
	return rootCmd
}
