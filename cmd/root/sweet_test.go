package root

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/spf13/cobra"
)

func (got Exercise) matches(want Exercise) bool {
	return got.name == want.name && got.text == want.text
}

func (got Exercise) matchesOneOf(wants []Exercise) bool {
	for _, want := range wants {
		if got.matches(want) {
			fmt.Printf("I match!\n")
			return true
		} else {
			fmt.Printf("no match\n")
		}
	}
	return false
}

func (ex Exercise) details() (m string) {
	m = fmt.Sprintf("\tname %s\n", ex.name)
	m += fmt.Sprintf("\tname bytes %v\n", []byte(ex.name))
	m += fmt.Sprintf("\ttext %s\n", ex.text)
	m += fmt.Sprintf("\ttext bytes %v\n", []byte(ex.text))
	return
}

func printExerciseFiles(dir string) (m string) {
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
		textBytes, _ := os.ReadFile(path.Join(dir, file.Name()))
		m += Exercise{
			name: name,
			text: string(textBytes),
		}.details()
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
			Args: cobra.MaximumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				got, gotErr := fromArgs(cmd, args)
				tc.check(got, gotErr)
			},
		}
		setCmdFlags(cmd)
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
				m := fmt.Sprintf("\n%s got\n", name)
				m += got.details()
				m += printExerciseFiles(tmpExercisesDir)
				t.Fatal(m)
			}
		},
		wantErr: nil,
		args:    []string{},
	}

	cmd := mockCmd(tc)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("%s: mock command failed to run: %s", tc.name, err)
	}
}
