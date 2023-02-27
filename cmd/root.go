package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "hckctl",
	Short: "The Cloud Native HaCKing Tool",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	cobra.OnInitialize(initConfig)
	addGlobalFlags()
	addCommands()
}

func initConfig() {
	InitCliConfig()
	InitFileLogger()
}

func addGlobalFlags() {
	const (
		LogLevelFlag = "log-level"
	)

	// --log-level
	rootCmd.PersistentFlags().String(LogLevelFlag, "", fmt.Sprintf("set the logging level for %s, one of: debug|info|warning|error", DefaultLogFile))
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup(LogLevelFlag))
}

func addCommands() {
	rootCmd.AddCommand(NewConfigCmd())
	rootCmd.AddCommand(NewBoxCmd())
	rootCmd.AddCommand(NewTemplateCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err)
	}
}
