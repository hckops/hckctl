package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/template"
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
	command.Flags().StringVarP(&revision, "revision", "r", "main", "megalopolis git source version i.e. branch|tag|sha")
	command.Flags().StringVarP(&path, "path", "p", "", "load a template from a local path")
	command.MarkFlagsMutuallyExclusive("revision", "path")
	return command
}

func runTemplateLocalCmd(path string) {
	log.Info().Msgf("loading local template: path=%s", path)

	data, err := loadTemplate(path)
	if err != nil {
		printFatalError(err, "unable to load template")
	}

	err = schema.ValidateAllSchema(data)
	if err != nil {
		printFatalError(err, "invalid template")
	}

	fmt.Print(data)
}

func runTemplateRemoteCmd(name, revision string) {
	log.Info().Msgf("requesting remote template: name=%s, revision=%s", name, revision)

	// TODO handle all templates
	data, err := template.RequestTemplate(NewBoxParam(name, revision))
	if err != nil {
		printFatalError(err, "unable to fetch template")
	}

	err = schema.ValidateAllSchema(data)
	if err != nil {
		printFatalError(err, "invalid template")
	}

	fmt.Print(data)
}

func loadTemplate(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "unable to load template: %s", path)
	}
	return string(data), nil
}

// TODO shared with box cmd
func printFatalError(err error, message string) {
	fmt.Println(message)
	log.Fatal().Err(err).Msg(message)
}
