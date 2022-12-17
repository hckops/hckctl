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
}

func init() {
	// --server-url
	rootCmd.PersistentFlags().StringP(ServerUrl, "u", "https://api.hckops.com", "TODO ServerUrl")
	viper.BindPFlag(ServerUrl, rootCmd.PersistentFlags().Lookup(ServerUrl))

	// --token
	rootCmd.PersistentFlags().StringP(Token, "t", "", "TODO Token")
	viper.BindPFlag(Token, rootCmd.PersistentFlags().Lookup(Token))

	// --local
	rootCmd.PersistentFlags().BoolP(Local, "l", false, "TODO Local")
	viper.BindPFlag(Local, rootCmd.PersistentFlags().Lookup(Local))

	rootCmd.AddCommand(NewBoxCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
