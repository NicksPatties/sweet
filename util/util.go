package util

import (
	"fmt"
)

// Helper function to create the usage string.
// Typically used for testing to verify command usage output.
func MakeUsageString(executableName string, subCommand string, usage string) string {
	msg := fmt.Sprintf("Usage: %s %s "+usage+"\n", executableName, subCommand) +
		fmt.Sprintf("For more information, run: %s help %s",
			executableName, subCommand)
	return msg
}

// Creates a usage function. This function should be assigned to
// the Usage property of a *flag.FlagSet variable.
//
// Example:
//
//	cmd := flag.NewFlagSet("open", flag.ExitOnError)
//	cmd.Usage = MakeUsage(os.Args[0], "open", "-p [port]")
func MakeUsage(executableName string, subCommand string, usage string) func() {
	return func() {
		fmt.Print(MakeUsageString(executableName, subCommand, usage))
	}
}
