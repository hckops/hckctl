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

	// define pointer/reference to pass around in all commands and initialize in each PersistentPreRunE
	configRef := &commonCmd.ConfigRef{}

	// TODO https://github.com/MakeNowJust/heredoc
	rootCmd := &cobra.Command{
		Use:   commonCmd.CliName,
		Short: "The Cloud Native HaCKing Tool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			// suppress messages on error
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			if config, err := configCmd.SetupConfig(); err != nil {
				return errors.Wrap(err, "unable to init config")
			} else {
				configRef.Config = config
			}

			if err := commonCmd.SetupLogger(configRef); err != nil {
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
	rootCmd.PersistentFlags().StringP(logLevelFlag, "l", commonCmd.NoneFlagShortHand, "set the logging level, one of: debug|info|warning|error")
	viper.BindPFlag(logLevelConfigKey, rootCmd.PersistentFlags().Lookup(logLevelFlag))

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	//rootCmd.AddCommand(boxCmd.NewBoxCmd(configRef))
	rootCmd.AddCommand(configCmd.NewConfigCmd(configRef))
	//rootCmd.AddCommand(labCmd.NewLabCmd(configRef))
	//rootCmd.AddCommand(templateCmd.NewTemplateCmd(configRef))
	rootCmd.AddCommand(NewVersionCmd())
	return rootCmd
}
