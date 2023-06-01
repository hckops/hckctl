package command

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	commonCmd "github.com/hckops/hckctl/pkg/command/common"
	configCmd "github.com/hckops/hckctl/pkg/command/config"
)

// env HCK_CONFIG_LOG.FILEPATH=/tmp/example.log ./build/hckctl config --log-level debug

func NewRootCmd() *cobra.Command {

	opts := &commonCmd.CommonCmdOptions{}

	// TODO https://github.com/MakeNowJust/heredoc
	rootCmd := &cobra.Command{
		Use:   commonCmd.CliName,
		Short: "The Cloud Native HaCKing Tool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			// suppress messages on error
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			config, err := configCmd.SetupConfig()
			if err != nil {
				return errors.Wrap(err, "unable to init config")
			}
			opts.Config = config

			if err := commonCmd.SetupLogger(opts.Config.Log); err != nil {
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

	const (
		logLevelFlag      = "log-level"
		logLevelConfigKey = "log.level"
	)

	// --log-level
	rootCmd.PersistentFlags().String(logLevelFlag, "", "set the logging level, one of: debug|info|warning|error")
	viper.BindPFlag(logLevelConfigKey, rootCmd.PersistentFlags().Lookup(logLevelFlag))

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	//rootCmd.AddCommand(boxCmd.NewBoxCmd(opts))
	rootCmd.AddCommand(configCmd.NewConfigCmd(opts))
	//rootCmd.AddCommand(labCmd.NewLabCmd(opts))
	//rootCmd.AddCommand(templateCmd.NewTemplateCmd(opts))
	rootCmd.AddCommand(NewVersionCmd())
	return rootCmd
}
