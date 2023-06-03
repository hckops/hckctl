package config

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
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
		Example: heredoc.Doc(`

			# prints current config, default path is XDG_CONFIG_HOME/hck/config.yml
			hckctl config

			# override config path: nested folders are automatically created if they don't exist
			HCK_CONFIG_DIR=/tmp/<PATH>/<SUB_PATH> hckctl config

			# config value override precedence (add "env" prefix to use dot notation): flag > env > config
			env HCK_CONFIG_LOG.LEVEL=error hckctl config --log-level debug
		`),
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
