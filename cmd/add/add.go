/*
Adds an exercise to Sweet.

This command can add exercises either from a file, or a remote source
like a GitHub repository.
*/
package add

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"path"
	"strings"

	"github.com/NicksPatties/sweet/util"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "add [flags] path",
	Short: "Add an exercise",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return addExercise(cmd, args)
	},
}

func addExercise(cmd *cobra.Command, args []string) (err error) {
	if len(args) != 1 {
		return errors.New("Incorrect number of args")
	}
	pathName := args[0]
	start, _ := cmd.Flags().GetUint("start")
	end, _ := cmd.Flags().GetUint("end")

	if start > end {
		fmt.Printf("YOU ARE BEING BAD!!!!")
		return errors.New("start flag cannot be greater than end flag")
	}

	inputFile, err := os.Open(pathName)
	if err != nil {
		return
	}
	defer inputFile.Close()

	// Create the new exercise file in the configuration directory.
	sweetConfigDir, err := util.SweetConfigDir()
	if err != nil {
		return
	}
	exercisesDir := path.Join(sweetConfigDir, "exercises")
	if envDir := os.Getenv("SWEET_EXERCISES_DIR"); envDir != "" {
		exercisesDir = envDir
	}
	newExercisePath := path.Join(exercisesDir, path.Base(pathName))
	newExerciseFile, err := os.Create(newExercisePath)
	if err != nil {
		return
	}
	defer newExerciseFile.Close()

	// Scan the input file and add the lines to the new exercise file.
	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanBytes)
	text := ""
	for scanner.Scan() {
		text += scanner.Text()
	}
	lines := util.Lines(text)
	if end > uint(len(lines)) {
		end = uint(len(lines))
	}
	newExerciseFileString := strings.Join(lines[start-1:end], "")
	newExerciseFile.WriteString(newExerciseFileString)

	if err = scanner.Err(); err != nil {
		return
	}
	return nil
}

func init() {
	setAddCmdFlags(Cmd)
}

func setAddCmdFlags(cmd *cobra.Command) {
	cmd.Flags().UintP("start", "s", 1, "The start line number to extract the sample")
	cmd.Flags().UintP("end", "e", math.MaxUint32, "The end line number to extract the sample")
}
