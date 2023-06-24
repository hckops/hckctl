package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/util"
)

// go tool nm ./build/hckctl | grep commit
var (
	release   string
	commit    string
	timestamp string
)

const devVersion = "dev"

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the client version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(readVersion())
		},
	}
}

func readVersion() string {
	if release == "" || commit == "" || timestamp == "" {
		return devVersion
	}
	return versionJson()
}

func versionJson() string {
	type model struct{ Version, Commit, Timestamp string }

	jsonString, _ := util.EncodeJson(model{
		Version:   release,
		Commit:    commit,
		Timestamp: timestamp,
	})

	return jsonString
}

// ClientVersion returns the ".Artifacts.Name" available in the PRO version only
// https://goreleaser.com/customization/templates/#artifacts
func ClientVersion() string {
	var version string
	if release == "" {
		version = devVersion
	} else {
		version = release
	}
	return fmt.Sprintf("%s-%s-%s-%s", common.CliName, version, runtime.GOOS, runtime.GOARCH)
}
