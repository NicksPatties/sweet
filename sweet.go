/*
sweet - The Software Engineer's Exercise in Typing.
Runs an interactive typing exercise.
Once complete, it displays statistics, including words per minute (WPM), accuracy, and number of mistakes.
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/NicksPatties/sweet/about"
	"github.com/NicksPatties/sweet/exercise"
	"github.com/NicksPatties/sweet/version"
)

// Assigned via -ldflags.
// Example:
//
//	go build -ldflags "-X github.com/NicksPatties/sweet/version.version=`date -u +.%Y%m%d%H%M%S`" .
//
// See https://stackoverflow.com/a/11355611 for details.
var sweetVersion string

const issueLink string = "issue-link"
const supportLink string = "support-link"

// Function types for each of the commands.
// Primarily used for dependency injection during tests.
//
// By convention, these functions accept the following parameters:
//
//	args []string - The arguments required to run the command.
//	executableName string - The name of the calling executable. Usually "sweet."
//
// Afterwards, pass whatever parameters you'd like into the function.
type Commands struct {
	exercise func(string, string, string) int
	version  func([]string, string, string) int
	about    func([]string, string, string, string, string) int
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
			code = commands.version(args[1:], sweetName, getVersion())
		case about.CommandName:
			code = commands.about(
				args[1:],
				getVersion(),
				issueLink,
				supportLink,
				os.Args[0],
			)
		default:
			fmt.Printf("Unregognized command")
		}
	}

	return code
}

func main() {
	defaultCommands := Commands{
		exercise: exercise.Run,
		version:  version.Run,
		about:    about.Run,
	}
	code := Run(os.Args[0], os.Args[1:], defaultCommands)

	if code != 0 {
		os.Exit(code)
	}
}

func getVersion() string {
	if sweetVersion == "" {
		return "debug"
	} else {
		return sweetVersion
	}
}
