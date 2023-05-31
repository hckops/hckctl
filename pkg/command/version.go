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

	command := &cobra.Command{
		Use:   "version",
		Short: "print client version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(versionJson())
		},
	}

	return command
}

func version() string {
	if commit == "" || timestamp == "" {
		return "dev"
	}
	return fmt.Sprintf("commit=%s timestamp=%s", commit, timestamp)
}

func versionJson() string {
	type version struct{ Commit, Timestamp string }

	jsonString, _ := util.ToJson(version{
		Commit:    commit,
		Timestamp: timestamp,
	})

	return jsonString
}
