/*
Adds an exercise to Sweet.

This command can add exercises either from a file, or a remote source
like a GitHub repository.
*/
package add

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path"

	l "github.com/NicksPatties/sweet/log"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "add [flags] path",
	Short: "Add an exercise",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			l.PrintErr("No path specified. Exiting.")
			os.Exit(1)
		} else if len(args) > 1 {
			l.PrintWarn("Multiple paths specified. Ignoring paths %s", args[1:])
		}
		path := args[0]
		start, _ := cmd.Flags().GetInt("start")
		end, _ := cmd.Flags().GetInt("end")

		// ... the rest of the command
		addExercise(path, start, end)
	},
}

func init() {
	Command.Flags().IntP("start", "s", 1, "The start line number to extract the sample")
	Command.Flags().IntP("end", "e", math.MaxUint32, "The end line number to extract the sample")
}

// Adds an exercise to sweet's configured exercises directory.
func addExercise(pathName string, start int, end int) {
	// Open the exercise from the given path.
	inputFile, err := os.Open(pathName)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Exiting.")
		os.Exit(1)
	}
	defer inputFile.Close()

	// Create the new exercise file in the configuration directory.
	newExercisePath := path.Join("/home/nick/.config/sweet/exercises", path.Base(pathName))
	newExerciseFile, err := os.Create(newExercisePath)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Exiting.")
		os.Exit(1)
	}
	defer newExerciseFile.Close()

	// Scan the input file and add the lines to the new exercise file.
	scanner := bufio.NewScanner(inputFile)
	i := 1
	for scanner.Scan() {
		currLine := scanner.Text()
		if i >= start && i <= end {
			fmt.Printf("Adding line %d:\t%s\n", i, currLine)
			// Need to add the newline back in,
			// since the scanner splits the file
			newExerciseFile.WriteString(currLine + "\n")
		} else {
			fmt.Printf("Ignoring line %d:\t%s\n", i, currLine)
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Problem with scanning the file. %s", err)
	}
}
