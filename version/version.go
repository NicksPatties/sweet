/*
version - Prints the installed version of sweet.

Usage:

	sweet version
*/
package version

import (
	"flag"
	"fmt"
	"os"

	"github.com/NicksPatties/sweet/util"
)

const CommandName = "version"

// Assigned via -ldflags.
// Example:
//
//	go build -ldflags "-X github.com/NicksPatties/sweet/version.version=`date -u +.%Y%m%d%H%M%S`" .
//
// See https://stackoverflow.com/a/11355611 for details.
var version string

func Run(args []string) int {
	cmd := flag.NewFlagSet(CommandName, flag.ExitOnError)
	cmd.Usage = util.MakeUsage(os.Args[0], CommandName, "")

	cmd.Parse(args)

	if len(cmd.Args()) > 0 {
		fmt.Println("Error: Too many arguments")
		cmd.Usage()
		return 1
	}

	printVersion()
	return 0
}

func printVersion() {
	if version == "" {
		fmt.Println("debug")
	} else {
		fmt.Println(version)
	}
}
