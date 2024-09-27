package root

import (
	"fmt"
	"math"
	"os"
	"path"
	"testing"

	"github.com/spf13/cobra"
)

func createExerciseFiles(t *testing.T, dir string, exercises []Exercise) (files []*os.File) {
	for _, ex := range exercises {
		tmpFile, err := os.CreateTemp(dir, ex.name)
		if err != nil {
			t.Fatal("failed to create exercise file")
		}
		tmpFile.WriteString(ex.text)
		err = tmpFile.Close()
		if err != nil {
			t.Fatalf("failed to close temporary file %s", tmpFile.Name())
		}
		files = append(files, tmpFile)
	}
	return
}

func fileToExercise(t *testing.T, fileName string) (exercise Exercise) {
	exercise.name = path.Base(fileName)
	text, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("failed to read exercise file %s\n", fileName)
	}
	exercise.text = string(text)
	return
}

func TestFromArgs(t *testing.T) {

	type testCase struct {
		name    string
		want    Exercise
		wantErr error
		args    []string
	}

	var mockCmd = func(tc testCase) *cobra.Command {
		cmd := &cobra.Command{
			Use:  "sweet",
			Long: "The Software Engineer Exercise for Typing.",
			Args: cobra.MaximumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				// Actually run the test here
				got, gotErr := FromArgs(cmd, args)

				if gotErr == nil {
					if tc.wantErr != nil {
						t.Fatalf("%s: got no error, want error %s", tc.name, tc.wantErr)
					} else if got.name != tc.want.name || got.text != tc.want.text {
						m := fmt.Sprintf("%s: expected exercise don't match\n", tc.name)
						m += fmt.Sprintf("got name  %s\n", got.name)
						m += fmt.Sprintf("want name %s\n", tc.want.name)
						m += fmt.Sprintf("name bytes  %v\n", []byte(got.name))
						m += fmt.Sprintf("want bytes  %v\n\n", []byte(tc.want.name))
						m += fmt.Sprintf("got  text %s\n", got.text)
						m += fmt.Sprintf("want text %s\n", tc.want.text)
						m += fmt.Sprintf("text bytes  %v\n", []byte(got.text))
						m += fmt.Sprintf("want bytes  %v\n\n", []byte(tc.want.text))
						t.Fatal(m)
					}
				} else {
					if tc.wantErr == nil {
						t.Fatalf("%s: got error %s, wanted no error ", tc.name, gotErr.Error())
					}
				}
			},
		}

		cmd.Flags().StringP("language", "l", "", "Language for the typing game")
		cmd.Flags().UintP("start", "s", 0, "Language for the typing game")
		cmd.Flags().UintP("end", "e", math.MaxUint, "Language for the typing game")
		cmd.SetArgs(tc.args)
		return cmd
	}

	tmpExercisesDir := t.TempDir()
	prevExercisesDir := os.Getenv("SWEET_EXERCISES_DIR")
	defer os.Setenv("SWEET_EXERCISES_DIR", prevExercisesDir)
	os.Setenv("SWEET_EXERCISES_DIR", tmpExercisesDir)

	testExercises := []Exercise{
		{
			name: "tmpExercise",
			text: "the test\n",
		},
	}
	testFiles := createExerciseFiles(t, tmpExercisesDir, testExercises)

	tc := testCase{
		name:    "random exercise, no args",
		want:    fileToExercise(t, testFiles[0].Name()),
		wantErr: nil,
		args:    []string{},
	}

	cmd := mockCmd(tc)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("%s: mock command failed to run: %s", tc.name, err)
	}
}
