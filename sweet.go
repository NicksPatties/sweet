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
*/
package main

import (
	"flag"
	"os"

	"github.com/NicksPatties/sweet/exercise"
	"github.com/NicksPatties/sweet/help"
	"github.com/NicksPatties/sweet/version"
)

func main() {

	sweetName := "sweet"

	// DEFAULT FLAGS
	sweetCmd := flag.NewFlagSet(sweetName, flag.ExitOnError)
	sweetLang := sweetCmd.String("l", "", "The programming language to practice, based on extension name")
	sweetTopic := sweetCmd.String("t", "", "Do an exercise related to a given topic")

	// SUB-COMMANDS

	code := 1
	if len(os.Args) == 1 {
		// Default command
		exercise.Run(*sweetLang, *sweetTopic)
	} else {
		args := os.Args[1:]
		subCommand := args[0]
		switch subCommand {
		case version.CommandName:
			version.Run(args[1:])
		case help.CommandName:
			code = help.Run(args[1:])
		default:
			// Default command with flags
			sweetCmd.Parse(args)
			exercise.Run(*sweetLang, *sweetTopic)
		}
	}
	os.Exit(code)
}
