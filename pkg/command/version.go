package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/util"
)

// go tool nm ./build/hckctl | grep commit
var (
	version   string
	commit    string
	timestamp string
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print client version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(readVersion())
		},
	}
}

func readVersion() string {
	if version == "" || commit == "" || timestamp == "" {
		return "dev"
	}
	return versionJson()
}

func versionJson() string {
	type model struct{ Version, Commit, Timestamp string }

	jsonString, _ := util.EncodeJson(model{
		Version:   version,
		Commit:    commit,
		Timestamp: timestamp,
	})

	return jsonString
}
