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
			help.Run(args[1:])
		default:
			// Default command with flags
			sweetCmd.Parse(args)
			exercise.Run(*sweetLang, *sweetTopic)
		}
	}
}
