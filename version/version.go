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
const version = "v0.0.2"

func Run(args []string) int {
	cmd := flag.NewFlagSet(CommandName, flag.ExitOnError)
	cmd.Usage = util.MakeUsage(os.Args[0], CommandName, "")

	cmd.Parse(args)

	if len(cmd.Args()) > 0 {
		fmt.Println("Error: Too many arguments")
		cmd.Usage()
		return 1
	}

	fmt.Println(version)
	return 0
}

func Usage() {
	fmt.Printf("Usage: %s\n", CommandName)
	fmt.Printf("Run %s %s for more information", "sweet", CommandName)
}
