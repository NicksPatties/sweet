/*
Sweet is a Software Engineering Exercise for Typing. In other words, it's a touch typing
exercise command line interface specifically designed for programmers.

Usage:

	sweet [subCommand] [flags]

Subcommands:

	help, h
		Prints a help message containing instructions on how to use the software

	js, ts, css, kt... and other file extensions
		Starts an exercise of that type of file

Flags:

	-f [path]
		Starts an exercise using the contents of the file.

	-u [url]
		Starts an exercise using the contents of a Github or Gitlab url.
		Line number fragment identifiers can be included to reduce the size of
		the exercise.
*/
package main

import (
	"fmt"
	"log"
	"os"
)

func printHelpMessage() {
	name := "sweet - a Software Engineering Exercise for Typing"
	type row struct {
		name        string
		description string
	}

	commands := []row{
		{"help, h", "Opens this help menu"},
		{"js, ts, go, ...", "Starts an exercise using the chosen language"},
	}

	flags := []row{
		{"-f [path]", "Starts an exercise with this file"},
		{"-u [url]", "Starts an exercise with the url"},
	}

	fmt.Printf("%s\n\n", name)
	fmt.Printf("  Commands\n\n")
	for _, cmd := range commands {
		fmt.Printf("    %-20s %s\n", cmd.name, cmd.description)
	}

	fmt.Printf("\n  Flags\n\n")
	for _, cmd := range flags {
		fmt.Printf("    %-20s %s\n", cmd.name, cmd.description)
	}
	fmt.Printf("\n")
}

func main() {
	var name string
	var exercise string
	var err error

	if len(os.Args) > 1 {
		switch arg := os.Args[1]; arg {
		case "help", "h":
			printHelpMessage()
			os.Exit(0)
		default:
			fmt.Printf("made a %s command\n", arg)
			name, exercise, err = GetExerciseFromDir(arg)
		}
	} else {
		name, exercise, err = GetRandomExercise()
	}

	if err != nil {
		log.Fatalf("Whoops! %s", err.Error())
	}

	// run the session
	m := RunSession(name, exercise)

	// show the results
	ShowResults(m)
}
