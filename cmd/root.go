package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hckops/hckctl/cmd/box"
	"github.com/hckops/hckctl/cmd/template"
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
	addCommands()

	cobra.OnInitialize(initConfig)
}

func addGlobalFlags() {
	// --log-level
	rootCmd.PersistentFlags().String(LogLevelFlag, "", "Set the logging level, one of: debug|info|warning|error")
	viper.BindPFlag(LogLevelFlag, rootCmd.PersistentFlags().Lookup(LogLevelFlag))
}

func addCommands() {
	rootCmd.AddCommand(box.NewBoxCmd())
	rootCmd.AddCommand(template.NewTemplateCmd())
}

func initConfig() {
	flags := NewFlags()

	InitFileLogger(flags)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
