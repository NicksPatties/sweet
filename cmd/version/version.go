/*
version - Prints the installed version of sweet.

Usage:

	sweet version
*/
package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Included via a compiler flag.
// i.e. go build -ldflags "-X github.com/NicksPatties/sweet/cmd/version.version=v0.1.0" .
var version string

var Command = &cobra.Command{
	Use:   "version",
	Short: "prints the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(getVersion(version))
	},
}

func getVersion(v string) string {
	if v == "" {
		return "dev"
	}
	return v
}
