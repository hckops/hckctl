package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

type globalFlags struct {
	serverUrl *string
	token     *string
	local     *bool
}

var rootCmd = &cobra.Command{
	Use:   "hckctl",
	Short: "Cloud Native HaCKing Tool",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

func init() {
	var globalFlags = &globalFlags{
		serverUrl: rootCmd.PersistentFlags().String("server-url", "", "TODO ServerUrl"),
		token:     rootCmd.PersistentFlags().StringP("token", "t", "", "TODO Token"),
		local:     rootCmd.PersistentFlags().Bool("local", false, "TODO Local"),
	}

	rootCmd.AddCommand(NewBoxCmd(globalFlags))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
