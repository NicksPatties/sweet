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

const CommandName = "version"

var Command = &cobra.Command{
	Use:   "version",
	Short: "prints the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("version")
	},
}
