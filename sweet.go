package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/NicksPatties/sweet/commands"
	"github.com/NicksPatties/sweet/util"
)

func actualMain() {

	sweetName := commands.GetCommandSweet()

	// DEFAULT FLAGS
	sweetCmd := flag.NewFlagSet(sweetName, flag.ExitOnError)
	sweetCmd.Usage = commands.PrintSweetUsage
	sweetLang := sweetCmd.String("l", "", "The programming language to practice, based on extension name")
	sweetTopic := sweetCmd.String("t", "", "Do an exercise related to a given topic")

	// SUB-COMMANDS
	versionCmd := flag.NewFlagSet(commands.CommandVersion, flag.ExitOnError)
	versionCmd.Usage = commands.PrintVersionUsage

	helpCmd := flag.NewFlagSet(commands.CommandHelp, flag.ExitOnError)
	helpCmd.Usage = commands.PrintHelpUsage

	if len(os.Args) == 1 {
		// Default command
		RunExercise(*sweetLang, *sweetTopic)
	} else {
		args := os.Args[1:]
		switch args[0] {
		case commands.CommandVersion:
			versionCmd.Parse(args[1:])
			commands.PrintVersion()
		case commands.CommandHelp:
			helpCmd.Parse(args[1:])
			if len(helpCmd.Args()) > 1 {
				fmt.Println("Too many arguments")
			}
			commands.RunHelp(helpCmd.Arg(0))
		default:
			// Default command with flags
			sweetCmd.Parse(args)
			RunExercise(*sweetLang, *sweetTopic)
		}
	}
}

func main() {
	filePath := "./one.txt"
	checksum1, err := util.HashFile(filePath)
	if err != nil {
		fmt.Printf("Error hashing file: %v\n", err)
		return
	}

	otherFilePath := "./two.txt"
	checksum2, err := util.HashFile(otherFilePath)
	if err != nil {
		fmt.Printf("Error hashing file: %v\n", err)
		return
	}

	fmt.Printf("CRC-32 checksum for one.txt: %08x\n", checksum1)
	fmt.Printf("CRC-32 checksum for two.txt: %08x\n", checksum2)
}
