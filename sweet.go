/*
sweet - The Software Engineer's Exercise in Typing.
Runs an interactive typing exercise.
Once complete, it displays statistics, including words per minute (WPM), accuracy, and number of mistakes.

Usage:

	sweet [sub-command] [flags] [exercise]

Sub-commands:

	help
		Prints a help message.
	version
		Prints the version of the application.
	about
		Prints general information about sweet.

For more information on sweet's subcommands, use sweet help [sub-command]

Flags:

	-l
		The language of the exercise. When given, sweet selects a random exercise that matches the given language. Examples: go, ts, rs
	-t
		The topic of the exercise. When given, sweet selects a random exercise that matches the give topic. Examples: sorting, algorithms, search
	-f
		The name of the exersize file
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/NicksPatties/sweet/exercise"
	"github.com/NicksPatties/sweet/help"
	"github.com/NicksPatties/sweet/version"
)

// Function types for each of the commands.
// Primarily used for dependency injection during tests.
type Commands struct {
	exercise func(string, string, string) int
	help     func([]string) int
	version  func([]string) int
}

// Runs the sweet top level command.
// Parses the arguments and runs the appropriate subcommands.
func Run(executableName string, args []string, commands Commands) int {
	sweetName := executableName

	sweetCmd := flag.NewFlagSet(sweetName, flag.ExitOnError)
	sweetLang := sweetCmd.String("l", "", "The programming language to practice, based on extension name")
	sweetTopic := sweetCmd.String("t", "", "Do an exercise related to a given topic")
	sweetFile := sweetCmd.String("f", "", "A specific file to practice")

	err := sweetCmd.Parse(args)
	if err != nil {
		return 1
	}

	code := 1
	if len(sweetCmd.Args()) == 0 {

		code = commands.exercise(*sweetLang, *sweetTopic, *sweetFile)
	} else {
		args := sweetCmd.Args()
		subCommand := args[0]
		switch subCommand {
		case version.CommandName:
			code = commands.version(args[1:])
		case help.CommandName:
			code = commands.help(args[1:])
		default:
			fmt.Printf("Unregognized command")
		}
	}

	return code
}

func main() {
	defaultCommands := Commands{
		exercise: exercise.Run,
		help:     help.Run,
		version:  version.Run,
	}
	code := Run(os.Args[0], os.Args[1:], defaultCommands)

	if code != 0 {
		os.Exit(code)
	}
}
