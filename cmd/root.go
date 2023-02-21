package cmd

import (
	"log"

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
	// removes timestamps
	log.SetFlags(0)

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	addGlobalFlags()

	config := InitCliConfig()
	InitFileLogger(&config.Log)
	addCommands(config)
}

func addGlobalFlags() {
	const (
		LogLevelFlag = "log-level"
	)

	// --log-level
	rootCmd.PersistentFlags().String(LogLevelFlag, "", "Set the logging level, one of: debug|info|warning|error")
	viper.BindPFlag(LogLevelFlag, rootCmd.PersistentFlags().Lookup(LogLevelFlag))
}

func addCommands(config *CliConfig) {
	rootCmd.AddCommand(NewConfigCmd())
	rootCmd.AddCommand(NewBoxCmd(&config.Box))
	rootCmd.AddCommand(NewTemplateCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
