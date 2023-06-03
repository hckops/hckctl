package template

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"

	"github.com/hckops/hckctl/pkg/util"
)

type Format string

type templateShowCmdOptions struct {
	template *templateCmdOptions
	format   formatFlag
}

func NewTemplateShowCmd(templateOpts *templateCmdOptions) *cobra.Command {

	opts := &templateShowCmdOptions{
		template: templateOpts,
	}

	command := &cobra.Command{
		Use:   "show",
		Short: "validate and print template",
		RunE:  opts.run,
	}

	const (
		formatFlagName = "format"
	)

	formatValue := enumflag.New(&opts.format, formatFlagName, toFormatIds(), enumflag.EnumCaseInsensitive)
	formatUsage := fmt.Sprintf("output format, one of %s", strings.Join([]string{string(yamlFormat), string(jsonFormat)}, "|"))
	command.Flags().Var(formatValue, formatFlagName, formatUsage)

	return command
}

func (opts *templateShowCmdOptions) run(cmd *cobra.Command, args []string) error {
	format := opts.format.value()

	if opts.template.path != "" && len(args) > 0 {
		return errors.New(fmt.Sprintf("unexpected arguments: %v", args))
	} else if opts.template.path != "" {
		return showLocalTemplate(format, opts.template.path)
	} else if len(args) == 1 {
		return showRemoteTemplate(format, args[0], opts.template.revision)
	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func showLocalTemplate(format Format, path string) error {
	log.Debug().Msgf("show local template: format=%v path=%s", format, path)

	localTemplate, err := util.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "local template not found %s", localTemplate)
	}

	// TODO validation

	fmt.Print(localTemplate)
	return nil
}

func showRemoteTemplate(format Format, name, revision string) error {
	log.Debug().Msgf("show remote template: format=%v name=%s revision=%s", format, name, revision)
	return nil
}
