/*
version - Prints the installed version of sweet.

Usage:

	sweet version
*/
package version

import (
	"fmt"

	"github.com/NicksPatties/sweet/util"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "prints the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(util.GetVersion())
	},
}
