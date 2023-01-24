package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hckops/hckctl/internal/template"
	"github.com/spf13/cobra"
)

func NewTemplateCmd() *cobra.Command {
	var revision string
	var path string

	command := &cobra.Command{
		Use:   "template [NAME]",
		Short: "load and validate a template",
		Run: func(cmd *cobra.Command, args []string) {
			if path != "" {
				runTemplateLocalCmd(path)
			} else if len(args) == 1 {
				runTemplateRemoteCmd(args[0], revision)
			} else {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	command.Flags().StringVarP(&revision, "revision", "r", "main", "git source version i.e. branch|tag|sha")
	command.Flags().StringVarP(&path, "path", "p", "", "load a template from a local path")
	command.MarkFlagsMutuallyExclusive("revision", "path")
	return command
}

func runTemplateLocalCmd(path string) {
	data, err := loadTemplate(path)
	if err != nil {
		log.Fatalln(err)
	}

	err = template.ValidateAllSchema(data)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(data)
}

func runTemplateRemoteCmd(name, revision string) {
	data, err := template.FetchTemplate(name, revision)
	if err != nil {
		log.Fatalln(err)
	}

	err = template.ValidateAllSchema(data)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(data)
}

func loadTemplate(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("unable to load the template")
	}
	return string(data), nil
}
