package command

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	boxCmd "github.com/hckops/hckctl/pkg/command/box"
	commonCmd "github.com/hckops/hckctl/pkg/command/common"
	configCmd "github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/command/config/setup"
	labCmd "github.com/hckops/hckctl/pkg/command/lab"
	templateCmd "github.com/hckops/hckctl/pkg/command/template"
	"github.com/hckops/hckctl/pkg/logger"
)

func NewRootCmd() *cobra.Command {

	// define pointer/reference to pass around in all commands and initialize in each PersistentPreRunE
	configRef := &commonCmd.ConfigRef{}
	var logCallback func() error

	rootCmd := &cobra.Command{
		Use:   commonCmd.CliName,
		Short: "The Cloud Native HaCKing Tool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			// suppress messages on error
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			if config, err := setup.SetupConfig(); err != nil {
				return errors.Wrap(err, "unable to init config")
			} else {
				configRef.Config = config
			}

			if callback, err := setup.SetupLogger(configRef); err != nil {
				return errors.Wrap(err, "unable to init log")
			} else {
				logCallback = callback
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			// close properly log file
			return logCallback()
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
	logLevelUsage := fmt.Sprintf("set the logging level, one of: %s", strings.Join(logger.LevelValues(), "|"))
	rootCmd.PersistentFlags().StringP(logLevelFlag, "l", commonCmd.NoneFlagShortHand, logLevelUsage)
	viper.BindPFlag(logLevelConfigKey, rootCmd.PersistentFlags().Lookup(logLevelFlag))

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	rootCmd.AddCommand(boxCmd.NewBoxCmd(configRef))
	rootCmd.AddCommand(configCmd.NewConfigCmd(configRef))
	rootCmd.AddCommand(labCmd.NewLabCmd(configRef))
	rootCmd.AddCommand(templateCmd.NewTemplateCmd())
	rootCmd.AddCommand(NewVersionCmd())
	return rootCmd
}