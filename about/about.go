/*
about - Prints some details about sweet.
This shows the name of the application, some author contact details,
and repository information.
*/
package about

import (
	"flag"
	"fmt"
	"os"

	"github.com/NicksPatties/sweet/util"
)

const CommandName = "about"

func Run(args []string) int {

	cmd := flag.NewFlagSet(CommandName, flag.ExitOnError)
	cmd.Usage = util.MakeUsage(os.Args[0], CommandName, "")

	cmd.Parse(args)

	fmt.Printf("This is sweet!\n")

	return 0
}
