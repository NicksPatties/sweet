package util

import (
	"fmt"
)

// Creates a usage function given a format string and a sub command.
// The output should be assigned to the Usage field of a *flag.FlagSet variable.
//
// Example:
//
//	cmd.Usage = MakeUsage("%s -p [port]", "open")
func MakeUsage(executableName string, subCommand string, usage string) func() {
	return func() {
		fmt.Printf("Usage: %s %s "+usage+"\n", executableName, subCommand)
		fmt.Printf("For more information, run: %s help %s", executableName, subCommand)
	}
}
