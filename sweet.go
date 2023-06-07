/*
Sweet is a Software Engineering Exercise for Typing. In other words, it's a touch typing
exercise command line interface specifically designed for programmers.

	Subcommands

	  help, h              Opens this help menu
	  add, a [path]        Adds a file to the exercise list
	  list, l              Lists the available exercises to run
	  [exercise name]      Runs this exercise

	Flags

	  [-go|-js|-ts|...]    Runs an exercise of the given file extension
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

	subcommands := []row{
		{"help, h", "Opens this help menu"},
		{"add, a [path]", "Adds a file to the exercise list"},
		{"list, l", "Lists the available exercises to run"},
		{"[exercise name]", "Runs this exercise"},
	}

	flags := []row{
		{"[-go|-js|-ts|...]", "Runs an exercise of the given file extension"},
	}

	fmt.Printf("%s\n\n", name)
	fmt.Printf("  Subcommands\n\n")
	for _, scmd := range subcommands {
		fmt.Printf("    %-20s %s\n", scmd.name, scmd.description)
	}

	fmt.Printf("\n  Flags\n\n")
	for _, f := range flags {
		fmt.Printf("    %-20s %s\n", f.name, f.description)
	}
	fmt.Printf("\n")
}

// Adds an exercise from the specified path into sweet's exercise directory, making
// it available to use for the exercises. If an exercise is already present, then
// it will just return the given exercise.
func addExercise(srcPath string) (string, string, error) {
	destPath, err := getDefaultConfigPath()
	if err != nil {
		return "", "", err
	}

	exPath, err := addFileToDirectory(srcPath, destPath)
	if err != nil {
		return "", "", err
	}

	return GetExerciseFromFile(exPath)
}

func main() {

	var name string
	var exercise string
	var err error

	// check if the $HOME/.sweet directory is there, create the directory, and then add the default exercises
	err = addDefaultExercises()
	if err != nil {
		log.Fatalf("Whoops! %s", err.Error())
	}

	// making a command
	if len(os.Args) > 1 {
		switch arg := os.Args[1]; arg {
		case "help":
			printHelpMessage()
			os.Exit(0)
		case "add":
			name, exercise, err = addExercise(arg)
		case "list":
			fmt.Print("This is the list command\n")
			os.Exit(0)
		default:
			fmt.Printf("Search for the %s exercise in the .sweet directory\n", arg)
			os.Exit(0)
		}
	} else {
		name, exercise, err = GetRandomExercise()
	}

	if err != nil {
		log.Fatalf("Whoops! %s", err.Error())
	}

	// run the session
	m := RunSession(name, exercise)

	if m.quitEarly {
		fmt.Println("Goodbye!")
		os.Exit(0)
	}

	// show the results
	ShowResults(m)
}
