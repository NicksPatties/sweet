/*
Sweet is a Software Engineering Exercise for Typing. In other words, it's a touch typing
exercise command line interface specifically designed for programmers.

Upon first execution, Sweet will create a configuration folder at $HOME/.sweet. Some sample
exercises will also be created in the $HOME/.sweet/exercises directory, and the typing game
will begin, selecting a random exercise from the aforementioned directory.

Type the letters that are highlighted by the cursor, following along with the code that
appears in the exercise. Once the final character of the exercise is inputted correctly,
your WPM (words per minute), mistakes, and accuracy are printed in the console, and the
program ends.

If you'd like some more exercises, you can use the "add" command! Provide a path to "add",
and the file will be added to the exercises directory and available to use. Additionally,
sweet will immediately run an exercise with the provided file.

You can focus on only testing specific languages by using the "lang" command. Provide a
file extension corresponding to the kind of file you'd like to practice. For instance, want
to practice writing go code? Try using `sweet lang go`, and a random go exercise will
be provided for you.

Not sure what exercises are available to use? Use the "list" command to see which exercises
are available to practice. Then, run "sweet <exercise-name>" to begin your exercise!

Subcommands

	help                            Opens this help menu
	add [path]                      Adds a file to the exercise list
	lang [go|js|ts|java...]	        Finds a random exercise with the specified extension
	list                            Lists the available exercises to run
	[exercise name]                 Runs this exercise
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
func listExercises() (string, error) {
	ePath, err := getDefaultExercisesPath()
	if err != nil {
		return "", err
	}

	paths, err := getAllFilePathsInDirectory(ePath)
	if err != nil {
		return "", err
	}

	exercises := ""
	for _, path := range paths {
		str := strings.Replace(path, ePath, "", 1)
		exercises += fmt.Sprintln(str[1:])
	}
	return exercises, nil
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
			name, exercise, err = getExerciseFromFile(exPath)
		case "list":
			exs, err := listExercises()
			if err != nil {
				fmt.Printf("Something went wrong with listing the exercises")
				os.Exit(1)
			}
			fmt.Print(exs)
			os.Exit(0)
		case "lang":
			if len(args) != 2 {
				fmt.Println("Print usage message for lang command")
				os.Exit(1)
			}
			name, exercise, err = getExerciseForLang(args[1])
		default:
			exName := arg
			exPath, err := getDefaultExercisesPath()
			if err != nil {
				fmt.Printf("Something went wrong with getting this exercise\n")
				os.Exit(1)
			}
			name, exercise, err = getExerciseFromFile(path.Join(exPath, exName))
			if err != nil {
				fmt.Printf("Something went wrong with getting this exercise\n")
				os.Exit(1)
			}
		}
	} else {
		name, exercise, err = getRandomExercise()
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
	showResults(m)
}
