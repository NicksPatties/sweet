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
	"path"
	"strings"
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
// it available to use for the exercises. This function returns the path of the exercises
// once it's created, and a possible error if something goes wrong.
func addExercise(srcPath string) (string, error) {
	sweetPath, err := getDefaultConfigPath()
	if err != nil {
		return "", err
	}
	destPath := path.Join(sweetPath, EXERCISES_DIR_NAME)
	addedPath, err := addFileToDirectory(srcPath, destPath)

	return addedPath, err
}

// lists all of the available exercises in the exercises directory
func listExercises() error {
	ePath, err := getDefaultExercisesPath()
	if err != nil {
		return err
	}

	paths, err := getAllFilePathsInDirectory(ePath)
	if err != nil {
		return err
	}

	for _, path := range paths {
		str := strings.Replace(path, ePath, "", 1)
		fmt.Println(str[1:])
	}
	return nil
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

	args := os.Args[1:]

	// making a command
	if len(args) > 0 {
		switch arg := args[0]; arg {
		case "help":
			printHelpMessage()
			os.Exit(0)
		case "add":
			if len(args) != 2 {
				fmt.Println("Print usage message for add command")
				os.Exit(1)
			}
			srcPath := args[1]
			exPath, err := addExercise(srcPath)
			if err != nil {
				fmt.Printf("Something went wrong adding the exercise... %s", err)
				os.Exit(1)
			}
			name, exercise, err = GetExerciseFromFile(exPath)
		case "list":
			err := listExercises()
			if err != nil {
				fmt.Printf("Something went wrong with listing the exercises")
				os.Exit(1)
			}
			os.Exit(0)
		default:
			exName := arg
			exPath, err := getDefaultExercisesPath()
			if err != nil {
				fmt.Printf("Something went wrong with getting this exercise")
				os.Exit(1)
			}
			name, exercise, err = GetExerciseFromFile(path.Join(exPath, exName))
			if err != nil {
				fmt.Printf("Something went wrong with getting this exercise")
				os.Exit(1)
			}
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
