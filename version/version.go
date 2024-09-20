/*
version - Prints the installed version of sweet.

Usage:

	sweet version
*/
package version

import (
	"flag"
	"fmt"
)

const CommandName = "version"

func Run(args []string, executableName string, version string) int {
	cmd := flag.NewFlagSet(CommandName, flag.ExitOnError)

	cmd.Parse(args)

	if len(cmd.Args()) > 0 {
		fmt.Println("Error: Too many arguments")
		cmd.Usage()
		return 1
	}

	fmt.Print(version)
	return 0
}
