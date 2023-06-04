package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/util"
)

// go tool nm ./build/hckctl | grep commit
var (
	commit    string
	timestamp string
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print client version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version())
		},
	}
}

func version() string {
	if commit == "" || timestamp == "" {
		return "dev"
	}
	return versionJson()
}

func versionJson() string {
	type version struct{ Commit, Timestamp string }

	jsonString, _ := util.EncodeJson(version{
		Commit:    commit,
		Timestamp: timestamp,
	})

	return jsonString
}