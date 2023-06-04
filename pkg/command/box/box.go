package box

import (
	"fmt"
	"github.com/hckops/hckctl/pkg/command/config"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hckops/hckctl/pkg/command/common"
)

type boxCmdOptions struct {
	configRef *config.ConfigRef
	path      string
	revision  string
}

func NewBoxCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "box [name]",
		Short: "attach and tunnel containers",
		RunE:  opts.run,
	}

	const (
		pathFlag     = "path"
		revisionFlag = "revision"
		providerFlag = "provider"
	)

	// --path
	command.Flags().StringVarP(&opts.path, pathFlag, "p", "", "local path")

	// --revision
	command.Flags().StringVarP(&opts.revision, revisionFlag, "r", common.TemplateRevision, common.TemplateRevisionUsage)
	viper.BindPFlag(fmt.Sprintf("template.%s", revisionFlag), command.Flags().Lookup(revisionFlag))

	command.MarkFlagsMutuallyExclusive(pathFlag, revisionFlag)

	// --provider
	command.Flags().StringP(providerFlag, common.NoneFlagShortHand, string(config.Docker), fmt.Sprintf("change box provider, one of %s",
		strings.Join([]string{string(config.Docker), string(config.Kubernetes), string(config.Argo), string(config.Cloud)}, "|")))
	viper.BindPFlag(fmt.Sprintf("box.%s", providerFlag), command.Flags().Lookup(providerFlag))

	command.AddCommand(NewBoxCopyCmd(opts))
	command.AddCommand(NewBoxCreateCmd(opts))
	command.AddCommand(NewBoxDeleteCmd(opts))
	command.AddCommand(NewBoxExecCmd(opts))
	command.AddCommand(NewBoxListCmd(opts))
	command.AddCommand(NewBoxOpenCmd(opts)) // default
	command.AddCommand(NewBoxTunnelCmd(opts))

	return command
}

func (opts *boxCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println(fmt.Sprintf("not implemented: path=%s revision=%s provider=%v",
		opts.path, opts.revision, opts.configRef.Config.Box.Provider))

	// TODO validation

	return nil
}
