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
	templateCmd "github.com/hckops/hckctl/pkg/command/template"
	versionCmd "github.com/hckops/hckctl/pkg/command/version"
	"github.com/hckops/hckctl/pkg/logger"
)

func NewRootCmd() *cobra.Command {

	// define pointer/reference to pass around in all commands and initialize in each PersistentPreRunE
	configRef := &configCmd.ConfigRef{}
	var logCallback func() error

	rootCmd := &cobra.Command{
		Use:   commonCmd.CliName,
		Short: "The Cloud Native HaCKing Tool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			// suppress messages on error
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			if config, err := setupConfig(); err != nil {
				return errors.Wrap(err, "unable to init config")
			} else {
				configRef.Config = config
			}

			if callback, err := setupLogger(configRef); err != nil {
				return errors.Wrap(err, "unable to init logger")
			} else {
				logCallback = callback
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			// close log file properly
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
	logLevelUsage := fmt.Sprintf("set the logging level, one of %s", strings.Join(logger.LevelValues(), "|"))
	rootCmd.PersistentFlags().StringP(logLevelFlag, "l", logger.InfoLogLevel.String(), logLevelUsage)
	viper.BindPFlag(logLevelConfigKey, rootCmd.PersistentFlags().Lookup(logLevelFlag))

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	rootCmd.AddCommand(boxCmd.NewBoxCmd(configRef))
	rootCmd.AddCommand(configCmd.NewConfigCmd(configRef))
	// TODO rootCmd.AddCommand(labCmd.NewLabCmd(configRef))
	rootCmd.AddCommand(templateCmd.NewTemplateCmd(configRef))
	rootCmd.AddCommand(versionCmd.NewVersionCmd())
	return rootCmd
}

// loads configs or initialize the default
func setupConfig() (*configCmd.ConfigV1, error) {
	err := configCmd.InitConfig(false)
	if err != nil {
		return nil, err
	}
	return configCmd.LoadConfig()
}

func setupLogger(configRef *configCmd.ConfigRef) (func() error, error) {
	logConfig := configRef.Config.Log
	logger.SetTimestamp()
	logger.SetLevel(logger.ParseLevel(logConfig.Level))
	logger.SetContext(commonCmd.CliName)
	return logger.SetFileOutput(logConfig.FilePath)
}
