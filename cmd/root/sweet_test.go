package root

import (
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func mockCommand() (mockCmd *cobra.Command) {
	mockCmd = &cobra.Command{}
	// taken from sweet.go:51
	mockCmd.Flags().StringP("language", "l", "", "Language for the typing game")
	mockCmd.Flags().UintP("start", "s", 0, "Language for the typing game")
	mockCmd.Flags().UintP("end", "e", math.MaxUint, "Language for the typing game")
	return
}

func TestValidateFlags(t *testing.T) {
	want := ExerciseArgs{
		file: 
	}
	tmpDir := t.TempDir()
	var tmpFile *os.File
	var err error
	for retry := 0; err != nil; retry++ {
		tmpFile, err = os.CreateTemp(tmpDir, "")
		if retry > 2 {
			t.Fatal("Failed to create temp file")
		}
	}
	tmpFile.Write([]byte("hello there\n"))
	mockCmd := mockCommand()
	mockCmd.Run = func(cmd *cobra.Command, args []string) {
		got := getExerciseArgs(cmd, args)

		if got.language != want.language {
			
		}
	}
	if err := mockCmd.Execute(); err != nil {
		t.Fatal("mockCmd failed to execute. " + err.Error())
	}
}
