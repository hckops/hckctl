package config

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hckops/hckctl/pkg/util"
)

// TODO add command to "set" a field with dot notation
// TODO add confirmation prompt before reset
type configCmdOptions struct {
	configRef *ConfigRef
}

func NewConfigCmd(configRef *ConfigRef) *cobra.Command {

	opts := &configCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "config",
		Short: "Print the current configurations",
		Example: heredoc.Doc(`

			# prints current configs, the default path is XDG_CONFIG_HOME/hck/config.yml
			hckctl config

			# override configs path: nested folders are automatically created if they don't exist
			HCK_CONFIG_DIR=/tmp/<PATH>/<SUB_PATH> hckctl config

			# config value override precedence (add "env" prefix to use dot notation): flag > env > config
			env HCK_CONFIG_LOG.LEVEL=error hckctl config --log-level debug
		`),
		RunE: opts.run,
	}

	resetCommand := &cobra.Command{
		Use:   "reset",
		Short: "Restore default configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return InitConfig(true)
			} else {
				cmd.HelpFunc()(cmd, args)
			}
			return nil
		},
	}

	command.AddCommand(resetCommand)

	return command
}

func (opts *configCmdOptions) run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		if value, err := util.EncodeYaml(opts.configRef.Config); err != nil {
			return errors.Wrap(err, "error encoding config")
		} else {
			fmt.Println(fmt.Sprintf("# %s", viper.ConfigFileUsed()))
			fmt.Print(value)
		}
	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}
