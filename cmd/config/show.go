package config

import (
	"fmt"
	"github.com/spf13/cobra"
)

// TODO format e.g. convert to json
type configShowCmdOptions struct {
	config *configCmdOptions
}

func NewConfigShowCmd(config *configCmdOptions) *cobra.Command {

	opts := &configShowCmdOptions{
		config: config,
	}

	command := &cobra.Command{
		Use:   "show",
		Short: "validate and print current configurations",
		RunE:  opts.run,
	}

	return command
}

func (opts *configShowCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
