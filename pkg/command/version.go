package command

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// go tool nm ./build/hckctl | grep commit
var (
	commit    string
	timestamp string
)

// TODO move in cmd folder

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

func Version() string {
	if commit == "" || timestamp == "" {
		return "dev"
	}
	return fmt.Sprintf("commit=%s timestamp=%s", commit, timestamp)
}

func versionJson() string {
	type version struct{ Commit, Timestamp string }

	bytes, _ := json.Marshal(version{
		Commit:    commit,
		Timestamp: timestamp,
	})

	return string(bytes)
}
