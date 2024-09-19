package util

import (
	"fmt"
	"io"
)

// Creates a usage function given a format string and a sub command.
// The output should be assigned to the Usage field of a *flag.FlagSet variable.
//
// Example:
//
//	cmd.Usage = MakeUsage("%s -p [port]", "open")
//
// Results:
//
//	Usage:
func MakeUsage(w io.Writer, executableName string, subCommand string, usage string) func() {
	return func() {
		fmt.Fprintf(w, "Usage: %s %s "+usage+"\n", executableName, subCommand)
		fmt.Fprintf(w, "For more information, run: %s help %s", executableName, subCommand)
	}
}
