/*
help - Prints the help message for sweet and its subcommands.

Usage:

	sweet help [sub-command]
*/
package help

import (
	"flag"
	"fmt"

	"github.com/NicksPatties/sweet/about"
	"github.com/NicksPatties/sweet/util"
	"github.com/NicksPatties/sweet/version"
)

const CommandName = "help"

// Runs the help command and returns the status code.
// The status code should follow the conventions of os.Exit()
func Run(args []string, executableName string) int {

	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)
	helpCmd.Usage = util.MakeUsage(executableName, CommandName, "[sub-command]")

	if len(args) == 0 {
		printSweetHelpMessage(executableName)
		return 0
	}

	err := helpCmd.Parse(args)
	if err != nil {
		return 1
	}

	if len(helpCmd.Args()) > 1 {
		fmt.Println("Too many arguments")
		return 1
	}

	subcommand := args[0]

	// interpret arguments
	switch subcommand {
	case CommandName:
		printHelpHelpMessage(executableName)
		return 0
	case version.CommandName:
		printVersionHelpMessage(executableName)
		return 0
	case about.CommandName:
		printAboutHelpMessage(executableName)
		return 0
	default:
		fmt.Printf("Unrecognized sub-command: %s\n", subcommand)
		printSweetHelpMessage(executableName)
		return 1
	}
}

// Prints help message for the main application
func printSweetHelpMessage(executableName string) {
	fmt.Printf(`Usage: %s [sub-command] [flags] [exercise]

Sub-commands:
	help		Show this help message
	version		Show the currently installed version
	about		Show an about page

Flags:
	-l	The language of the exercise.
	-t	The topic of the exercise.
`, executableName)
}

func Usage(executableName string) {
	fmt.Printf("Usage: %s %s [sub-command]\n", executableName, CommandName)
	fmt.Printf("Run %s %s [sub-command] for more information", executableName, CommandName)
}

// Prints help message for the help subcommand
func printHelpHelpMessage(executableName string) {
	msg := "%s help - Displays help information\n"

	fmt.Printf(msg, executableName)
}

func printVersionHelpMessage(executableName string) {
	msg := "%s version - Displays version of the application"

	fmt.Printf(msg, executableName)
}

func printAboutHelpMessage(executableName string) {
	fmt.Printf("%s about - Shows details and the creators of sweet", executableName)
}
