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

	// TODO >>> add config to opts?! not sure
	var opts = &commonCmd.GlobalCmdOptions{}
	var config *commonCmd.ConfigV1

	// TODO https://github.com/MakeNowJust/heredoc
	rootCmd := &cobra.Command{
		Use:   "hckctl",
		Short: "The Cloud Native HaCKing Tool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			opts.LogLevel = "FIXME"

			// suppress messages on error
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			fmt.Println("PersistentPreRunE")
			if err := configCmd.InitConfig(); err != nil {
				return errors.Wrap(err, "unable to init config")
			}
			var err error
			if config, err = configCmd.LoadConfig(); err != nil {
				return errors.Wrap(err, "unable to load config")
			}
			opts.InternalConfig = config
			if err := commonCmd.InitFileLogger(opts, config.Log); err != nil {
				return errors.Wrap(err, "unable to init log")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println("Run")
			cmd.HelpFunc()(cmd, args)
		},
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}

	// TODO --log-level

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	rootCmd.AddCommand(boxCmd.NewBoxCmd(opts))
	rootCmd.AddCommand(configCmd.NewConfigCmd(opts, config))
	rootCmd.AddCommand(labCmd.NewLabCmd(opts))
	rootCmd.AddCommand(templateCmd.NewTemplateCmd(opts))
	rootCmd.AddCommand(NewVersionCmd())
	return rootCmd
}
