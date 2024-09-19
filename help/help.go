package help

import (
	"flag"
	"fmt"
	"os"

	"github.com/NicksPatties/sweet/version"
)

const CommandName = "help"

func Run(args []string) {

	if len(args) == 0 {
		printHelpMessage()
		os.Exit(0)
	}

	// command parsing
	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)
	helpCmd.Usage = Usage

	helpCmd.Parse(args)

	if len(helpCmd.Args()) > 1 {
		fmt.Println("Too many arguments")
	}

	subcommand := args[0]

	// interpret arguments
	switch subcommand {
	case CommandName:
		printHelpHelpMessage()
	case version.CommandName:
		printVersionHelpMessage()
	default:
		fmt.Printf("Unrecognized sub command: %s", subcommand)
		PrintSweetUsage()
	}
}

// Prints help message for the main application
func printHelpMessage() {
	msg := "Sweet - The Software Engineer's Exercise for Typing\n" +
		"\n" +
		"SUB-COMMANDS\n" +
		"\n" +
		"\thelp\tPrints this helpful message\n" +
		"\tversion\tPrints the currently installed version of sweet\n" +
		"\n" +
		"For more information about specific sub-commands, use sweet help [sub-command]\n"

	fmt.Print(msg)
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

func PrintSweetUsage() {
	program := "sweet"
	p := fmt.Printf
	p("\n")
	p("Usage: %s", program)
}

func PrintHelpUsage() {
	fmt.Println("Help usage")
}

func PrintVersionUsage() {
	fmt.Println("version usage")
}

func PrintAddUsage() {
	fmt.Println("add usage")
}

func GetMessage() string {
	return "hey from help"
}
