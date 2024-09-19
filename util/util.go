package util

import (
	"fmt"
)

// Creates a usage function. This function should be assigned to
// the Usage property of a *flag.FlagSet variable.
//
// Example:
//
//	cmd := flag.NewFlagSet("open", flag.ExitOnError)
//	cmd.Usage = MakeUsage(os.Args[0], "open", "-p [port]")
func MakeUsage(executableName string, subCommand string, usage string) func() {
	return func() {
		fmt.Printf("Usage: %s %s "+usage+"\n", executableName, subCommand)
		fmt.Printf("For more information, run: %s help %s", executableName, subCommand)
	}
}
