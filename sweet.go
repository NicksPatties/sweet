package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/NicksPatties/sweet/exercise"
	_ "github.com/NicksPatties/sweet/help"
	_ "github.com/NicksPatties/sweet/util"
)

func main() {

	sweetName := "sweet"

	// DEFAULT FLAGS
	sweetCmd := flag.NewFlagSet(sweetName, flag.ExitOnError)
	sweetLang := sweetCmd.String("l", "", "The programming language to practice, based on extension name")
	sweetTopic := sweetCmd.String("t", "", "Do an exercise related to a given topic")

	// SUB-COMMANDS
	versionCmd := flag.NewFlagSet("version", flag.ExitOnError)

	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)

	if len(os.Args) == 1 {
		// Default command
		exercise.Run(*sweetLang, *sweetTopic)
	} else {
		args := os.Args[1:]
		switch args[0] {
		case "version":
			versionCmd.Parse(args[1:])

		case "help":
			helpCmd.Parse(args[1:])
			if len(helpCmd.Args()) > 1 {
				fmt.Println("Too many arguments")
			}

		default:
			// Default command with flags
			sweetCmd.Parse(args)
			exercise.Run(*sweetLang, *sweetTopic)
		}
	}
}
