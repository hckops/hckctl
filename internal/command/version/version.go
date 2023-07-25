package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the client version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(jsonVersion())
		},
	}
}
