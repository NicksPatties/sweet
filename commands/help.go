package commands

import (
	"fmt"
)

// help command flowchart

// sweet help
// - prints sweet help message

// sweet help [subcommand]
// - does the subcommand exist?
//   - No
//     - `"Subcommand not found: [subcommand]"`
//     - Print sweet help message
//   - Yes
//     - Print subcommand help message

// sweet help [random flags]
// - Error from flags module
// - prints sweet help message
func RunHelp(subcommand string) {
	// interpret arguments
	switch subcommand {
	case "help":
		PrintHelpUsage()
	default:
		fmt.Printf("Unrecognized sub command: %s", subcommand)
		PrintSweetUsage()
	}
}

func PrintSweetUsage() {
	program := GetCommandSweet()
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
