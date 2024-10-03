package root

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path"
	"testing"

	"github.com/spf13/cobra"
)

func (got Exercise) matches(want Exercise) bool {
	return got.name != want.name || got.text != want.text
}

func (got Exercise) matchesOneOf(wants []Exercise) bool {
	for _, want := range wants {
		if got.matches(want) {
			return true
		}
	}
	return false
}

func printExerciseFiles(t *testing.T, dir string) (m string) {
	entries, _ := os.ReadDir(dir)
	var files []fs.DirEntry
	for _, ent := range entries {
		if !ent.IsDir() {
			files = append(files, ent)
		}
	}

	m = fmt.Sprintf("wanted one of\n")
	for _, file := range files {
		name := file.Name()
		text, _ := os.ReadFile(path.Join(dir, file.Name()))
		m += fmt.Sprintf("\tname %s\n", name)
		m += fmt.Sprintf("\ttext %s\n\n", text)
	}
	return
}

func createExerciseFiles(t *testing.T, dir string, exercises []Exercise) {
	for _, ex := range exercises {
		tmpFile, err := os.Create(path.Join(dir, ex.name))
		if err != nil {
			t.Fatal("failed to create exercise file")
		}
		tmpFile.WriteString(ex.text)
		err = tmpFile.Close()
		if err != nil {
			t.Fatalf("failed to close temporary file %s", tmpFile.Name())
		}
	}
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
		check   func(Exercise, error)
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
				tc.check(got, gotErr)

				// if gotErr == nil {
				// 	if tc.wantErr != nil {
				// 		t.Fatalf("%s: got no error, want error %s", tc.name, tc.wantErr)
				// 	} else if got.name != tc.want.name || got.text != tc.want.text {
				// 		m := fmt.Sprintf("%s: expected exercise don't match\n", tc.name)
				// 		m += fmt.Sprintf("got name  %s\n", got.name)
				// 		m += fmt.Sprintf("want name %s\n", tc.want.name)
				// 		m += fmt.Sprintf("name bytes  %v\n", []byte(got.name))
				// 		m += fmt.Sprintf("want bytes  %v\n\n", []byte(tc.want.name))
				// 		m += fmt.Sprintf("got  text %s\n", got.text)
				// 		m += fmt.Sprintf("want text %s\n", tc.want.text)
				// 		m += fmt.Sprintf("text bytes  %v\n", []byte(got.text))
				// 		m += fmt.Sprintf("want bytes  %v\n\n", []byte(tc.want.text))
				// 		t.Fatal(m)
				// 	}
				// } else {
				// 	if tc.wantErr == nil {
				// 		t.Fatalf("%s: got error %s, wanted no error ", tc.name, gotErr.Error())
				// 	}
				// }
			},
		}

		cmd.Flags().StringP("language", "l", "", "Language for the typing game")
		cmd.Flags().UintP("start", "s", 0, "Language for the typing game")
		cmd.Flags().UintP("end", "e", math.MaxUint, "Language for the typing game")
		cmd.SetArgs(tc.args)
		return cmd
	}

	tmpExercisesDir := t.TempDir()
	t.Setenv("SWEET_EXERCISES_DIR", tmpExercisesDir)
	testExercises := []Exercise{
		{
			name: "tmpExercise",
			text: "the test\n",
		},
	}
	createExerciseFiles(t, tmpExercisesDir, testExercises)

	tc := testCase{
		name: "random exercise, no args",
		check: func(got Exercise, gotErr error) {
			name := "random exercise, no args"
			if gotErr != nil {
				t.Fatalf("%s wanted no error, got %s\n", name, gotErr)
			}
			if !got.matchesOneOf(testExercises) {
				m += fmt.Sprintf("%s got\n", name)

				m += fmt.Sprintf("\tname       %s\n", got.name)
				m += fmt.Sprintf("\tname bytes %v\n", []byte(got.name))
				m += fmt.Sprintf("\ttext       %s\n", got.text)
				m += fmt.Sprintf("\ttext bytes %v\n", []byte(got.text))
				m += printExerciseFiles(t, tmpExercisesDir)
				t.Fatal(m)
			}
		},
		args: []string{},
	}

	cmd := mockCmd(tc)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("%s: mock command failed to run: %s", tc.name, err)
	}
}
