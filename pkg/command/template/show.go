package template

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
)

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

	// --format (enum)
	formatValue := enumflag.New(&opts.format, formatFlagName, formatIds, enumflag.EnumCaseInsensitive)
	formatUsage := fmt.Sprintf("output format, one of %s", strings.Join(formatValues(), "|"))
	command.Flags().Var(formatValue, formatFlagName, formatUsage)

	return command
}

func (opts *templateShowCmdOptions) run(cmd *cobra.Command, args []string) error {
	format := opts.format.value()

	if opts.template.path != "" && len(args) > 0 {
		return errors.New(fmt.Sprintf("unexpected arguments: %v", args))
	} else if opts.template.path != "" {
		return printLocalTemplate(format, opts.template.path)
	} else if len(args) == 1 {
		return printRemoteTemplate(format, args[0], opts.template.revision)
	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}
