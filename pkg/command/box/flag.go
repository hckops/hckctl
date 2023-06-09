package box

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"

	"github.com/spf13/viper"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/command/common"
)

func addBoxProviderFlag(command *cobra.Command) {
	const (
		providerFlagName = "provider"
	)
	command.Flags().StringP(providerFlagName, common.NoneFlagShortHand, string(box.Docker),
		fmt.Sprintf("switch box provider, one of %s", strings.Join(box.BoxProviderValues(), "|")))
	viper.BindPFlag(fmt.Sprintf("box.%s", providerFlagName), command.Flags().Lookup(providerFlagName))
}
