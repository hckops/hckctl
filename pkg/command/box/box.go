package box

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag/v2"
	"strings"

	"github.com/hckops/hckctl/pkg/command/common"
)

type boxCmdOptions struct {
	configRef *common.ConfigRef
	path      string
	revision  string
}

func NewBoxCmd(configRef *common.ConfigRef) *cobra.Command {

	opts := &boxCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "box [name]",
		Short: "attach and tunnel a box",
		RunE:  opts.run,
	}

	const (
		pathFlag     = "path"
		revisionFlag = "revision"
		providerFlag = "provider"
	)

	// --path
	command.Flags().StringVarP(&opts.path, pathFlag, "p", "", "load a local template")

	// --revision
	command.Flags().StringVarP(&opts.revision, revisionFlag, "r", common.DefaultMegalopolisBranch, "megalopolis version i.e. branch|tag|sha")
	viper.BindPFlag(fmt.Sprintf("box.%s", revisionFlag), command.Flags().Lookup(revisionFlag))

	command.MarkFlagsMutuallyExclusive(pathFlag, revisionFlag)

	// --provider
	// possible bug: &provider reference (previously in opts) is always 0
	// enumflag is used only for validation, retrieve validated and merged value between config and flag from configRef
	var provider common.ProviderFlag
	providerValue := enumflag.New(&provider, providerFlag, common.ProviderIds, enumflag.EnumCaseInsensitive)
	command.Flags().Var(providerValue, providerFlag, fmt.Sprintf("set the box provider, one of %s",
		strings.Join([]string{string(common.Docker), string(common.Kubernetes), string(common.Argo), string(common.Cloud)}, "|")))
	viper.BindPFlag(fmt.Sprintf("box.%s", providerFlag), command.Flags().Lookup(providerFlag))

	command.AddCommand(NewBoxCopyCmd(opts))
	command.AddCommand(NewBoxCreateCmd(opts))
	command.AddCommand(NewBoxDeleteCmd(opts))
	command.AddCommand(NewBoxExecCmd(opts))
	command.AddCommand(NewBoxListCmd(opts))
	command.AddCommand(NewBoxOpenCmd(opts))
	command.AddCommand(NewBoxTunnelCmd(opts))

	return command
}

func (opts *boxCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println(fmt.Sprintf("not implemented: path=%s revision=%s providerConfig=%v",
		opts.path, opts.revision, opts.configRef.Config.Box.Provider))
	return nil
}
