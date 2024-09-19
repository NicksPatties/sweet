/*
help - Prints the help message for sweet and its subcommands.

Usage:

	sweet help [sub-command]
*/
package help

import (
	"flag"
	"fmt"
	"os"

	"github.com/NicksPatties/sweet/about"
	"github.com/NicksPatties/sweet/util"
	"github.com/NicksPatties/sweet/version"
)

const CommandName = "help"

// Runs the help command and returns the status code.
// The status code should follow the conventions of os.Exit()
func Run(args []string) int {

	if len(args) == 0 {
		printSweetHelpMessage()
		return 0
	}

	// command parsing
	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)
	helpCmd.Usage = util.MakeUsage(os.Args[0], CommandName, "[sub-command]")

	err := helpCmd.Parse(args)
	if err != nil {
		return 1
	}

	if len(helpCmd.Args()) > 1 {
		fmt.Println("Too many arguments")
	}

	subcommand := args[0]

	// interpret arguments
	switch subcommand {
	case CommandName:
		printHelpHelpMessage()
		return 0
	case version.CommandName:
		printVersionHelpMessage()
		return 0
	case about.CommandName:
		printAboutHelpMessage()
		return 0
	default:
		fmt.Printf("Unrecognized sub-command: %s\n", subcommand)
		printSweetHelpMessage()
		return 1
	}
}

// Prints help message for the main application
func printSweetHelpMessage() {
	executableName := "sweet"
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

func Usage() {
	fmt.Printf("Usage: %s [sub-command]\n", CommandName)
	fmt.Printf("Run %s %s [sub-command] for more information", "sweet", CommandName)
}

// Prints help message for the help subcommand
func printHelpHelpMessage() {
	msg := "help - Displays help information\n"

	fmt.Print(msg)
}

func printVersionHelpMessage() {
	msg := "version - Displays version of the application"

	fmt.Print(msg)
}

func printAboutHelpMessage() {
	fmt.Print("about - Shows details and the creators of sweet")
}
