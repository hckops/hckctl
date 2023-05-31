package command

import (
	"fmt"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"

	boxCmd "github.com/hckops/hckctl/pkg/command/box"
	commonCmd "github.com/hckops/hckctl/pkg/command/common"
	configCmd "github.com/hckops/hckctl/pkg/command/config"
	labCmd "github.com/hckops/hckctl/pkg/command/lab"
	templateCmd "github.com/hckops/hckctl/pkg/command/template"
)

func NewRootCmd() *cobra.Command {

	opts := &commonCmd.GlobalCmdOptions{}

	// TODO https://github.com/MakeNowJust/heredoc
	description := fmt.Sprintf("The Cloud Native HaCKing Tool - %s", version())

	rootCmd := &cobra.Command{
		Use:   "hckctl",
		Short: description,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			if err := configCmd.InitConfig(opts); err != nil {
				return errors.Wrap(err, "unable to init config")
			}
			if err := commonCmd.InitFileLogger(opts); err != nil {
				return errors.Wrap(err, "unable to init log")
			}

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
