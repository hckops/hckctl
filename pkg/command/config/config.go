package config

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/util"
)

// TODO add command to "set" a field with dot notation

func NewConfigCmd(configRef *common.ConfigRef) *cobra.Command {

	command := &cobra.Command{
		Use:   "config",
		Short: "print current configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			if value, err := util.ToYaml(configRef.Config); err != nil {
				return errors.Wrap(err, "error encoding config")
			} else {
				fmt.Println(fmt.Sprintf("# %s", viper.ConfigFileUsed()))
				fmt.Print(value)
			}
			return nil
		},
	}

	resetCommand := &cobra.Command{
		Use:   "reset",
		Short: "restore default configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return initConfig(true)
		},
	}

	command.AddCommand(resetCommand)

	return command
}
