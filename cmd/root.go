package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/cmd/template"
)

var rootCmd = &cobra.Command{
	Use:   "hckctl",
	Short: "The Cloud Native HaCKing Tool",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(template.NewTemplateCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
