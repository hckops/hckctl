package old

import (
	"fmt"
	"github.com/hckops/hckctl/pkg/command"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hckops/hckctl/internal/config"
)

// TODO add version: git + timestamp
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
	rootCmd.PersistentFlags().String(LogLevelFlag, "", fmt.Sprintf("set the logging level for %s, one of: debug|info|warning|error", config.DefaultLogFile))
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup(LogLevelFlag))
}

func addCommands() {
	rootCmd.AddCommand(NewConfigCmd())
	rootCmd.AddCommand(NewBoxCmd())
	rootCmd.AddCommand(NewTemplateCmd())
}

func Execute() {
	if err := command.NewRootCmd().Execute(); err != nil {
		log.Fatal().Err(err)
	}
}
