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
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	initFlags()
	initCommands()

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	flags := NewFlags()

	InitLogger(flags)
}

func initCommands() {
	rootCmd.AddCommand(box.NewBoxCmd())
	rootCmd.AddCommand(template.NewTemplateCmd())
}

func initFlags() {
	// --log-level
	rootCmd.PersistentFlags().String(LogLevelFlag, "", "Set the logging level, one of: debug|info|warning|error")
	viper.BindPFlag(LogLevelFlag, rootCmd.PersistentFlags().Lookup(LogLevelFlag))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

// TODO logging
// https://github.com/derailed/k9s/blob/master/cmd/root.go
// https://github.com/derailed/k9s/tree/0249f7cf2c2b403348e98f03a26355aadfbdfdda/internal/config
// https://github.com/rs/zerolog

// HOME https://github.com/adrg/xdg

// TODO config
// https://cobra.dev
// https://medium.com/@bnprashanth256/reading-configuration-files-and-environment-variables-in-go-golang-c2607f912b63
// https://github.com/kubernetes/minikube/blob/master/cmd/minikube/cmd/root.go
