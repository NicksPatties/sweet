/*
sweet - The Software Engineer's Exercise in Typing.
Runs an interactive typing exercise.
Once complete, it displays statistics, including words per minute (WPM), accuracy, and number of mistakes.
*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path"

	"github.com/NicksPatties/sweet/about"
	"github.com/NicksPatties/sweet/add"

	"github.com/NicksPatties/sweet/exercise"
	"github.com/NicksPatties/sweet/stats"
	"github.com/spf13/cobra"
)

type Exercise = exercise.Exercise

// Validates and returns the exercise arguments.
func getExerciseFromArgs(cmd *cobra.Command, args []string) (exercise exercise.Exercise, err error) {
	start, _ := cmd.Flags().GetUint("start")
	end, _ := cmd.Flags().GetUint("end")

	if start > end {
		err = errors.New(fmt.Sprintf("start flag %d cannot be greater than end flag %d", start, end))
		return
	}

	var file *os.File
	defer file.Close()
	if len(args) > 0 { // get the file from the argument
		if args[0] == "-" {
			file = os.Stdin
		} else {
			file, err = os.Open(args[0])
			if err != nil {
				return
			}

		}
	} else { // get a random exercise

		if start != 0 || end != math.MaxUint {
			err = errors.New("start and end should not be assigned for random exercise")
			return
		}

		var exercisesDir string
		if envDir := os.Getenv("SWEET_EXERCISES_DIR"); envDir != "" {
			exercisesDir = envDir

		} else {
			var configDir string
			configDir, err = os.UserConfigDir()
			if err != nil {
				return
			}
			exercisesDir = path.Join(configDir, "sweet", "exercises")
		}

		if err = os.MkdirAll(exercisesDir, 0775); err != nil {
			return
		}

		var entries []os.DirEntry
		entries, err = os.ReadDir(exercisesDir)
		if err != nil {
			return
		}

		var files []os.DirEntry
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			files = append(files, entry)
		}

		numFiles := len(files)
		if numFiles == 0 {
			err = errors.New("no files in the exercises directory")
			return
		}
		randI := rand.Intn(numFiles)
		filePath := path.Join(exercisesDir, files[randI].Name())
		file, err = os.Open(filePath)
		if err != nil {
			return
		}

	}

	var text string
	scanner := bufio.NewScanner(file)
	for line := uint(1); line <= end && scanner.Scan(); line++ {
		if line >= start {
			text += scanner.Text() + "\n"
		}
	}

	if text == "" {
		err = errors.New("no input text selected")
		return
	}

	exercise.Text = text
	exercise.Name = path.Base(file.Name())
	return
}

var rootCmd = &cobra.Command{
	Use:  "sweet",
	Long: "The Software Engineer Exercise for Typing.",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ex, err := getExerciseFromArgs(cmd, args)
		if err != nil {
			log.Fatal(err)
		}

		exercise.Run(ex)
	},
}

func init() {
	// Add language flag to root command only.
	// The flags for other commands will be defined in their respective modules.
	rootCmd.Flags().StringP("language", "l", "", "Language for the typing game")
	rootCmd.Flags().UintP("start", "s", 0, "Language for the typing game")
	rootCmd.Flags().UintP("end", "e", math.MaxUint, "Language for the typing game")

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	commands := []*cobra.Command{
		about.Command,
		stats.Command,
		add.Command,
	}

	for _, c := range commands {
		rootCmd.AddCommand(c)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
