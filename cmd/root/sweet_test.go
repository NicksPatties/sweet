package root

import (
	"errors"
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
			return true
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

	// Test environment setup
	tmpExercisesDir := t.TempDir()
	t.Setenv("SWEET_EXERCISES_DIR", tmpExercisesDir)
	testExercises := []Exercise{
		{
			name: "test.txt",
			text: "the test\n",
		},
		{
			name: "hello.js",
			text: "console.log('Hello!');\n",
		},
		{
			name: "hello.go",
			text: "fmt.Println(\"Hello!\")\n",
		},
		{
			name: "threelines.txt",
			text: "this file\nhas three lines\nof text.\n",
		},
	}
	createExerciseFiles(t, tmpExercisesDir, testExercises)

	type testCase struct {
		args  []string
		check func(Exercise, error)
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

	testCases := []testCase{
		{
			args: []string{},
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
		},
		{
			args: []string{path.Join(tmpExercisesDir, testExercises[0].name)},
			check: func(got Exercise, gotErr error) {
				name := "specific exercise, no args"
				want := testExercises[0]
				if gotErr != nil {
					t.Fatalf("%s wanted no error, got %s\n", name, gotErr)
				}
				if !got.matches(want) {
					m := fmt.Sprintf("\n%s got\n", name)
					m += got.details()
					m += want.details()
					t.Fatal(m)
				}
			},
		},
		{
			args: []string{
				path.Join(tmpExercisesDir, testExercises[3].name),
				"-s",
				"2",
			},
			check: func(got Exercise, gotErr error) {
				name := "specific exercise, start flag"
				want := Exercise{
					name: testExercises[3].name,
					text: "has three lines\nof text.\n",
				}
				if gotErr != nil {
					t.Fatalf("%s wanted no error, got %s\n", name, gotErr)
				}
				if !got.matches(want) {
					m := fmt.Sprintf("\n%s got\n", name)
					m += got.details()
					m += want.details()
					t.Fatal(m)
				}
			},
		},
		{
			args: []string{
				path.Join(tmpExercisesDir, testExercises[3].name),
				"-e",
				"2",
			},
			check: func(got Exercise, gotErr error) {
				name := "specific exercise, start flag"
				want := Exercise{
					name: testExercises[3].name,
					text: "this file\nhas three lines\n",
				}
				if gotErr != nil {
					t.Fatalf("%s wanted no error, got %s\n", name, gotErr)
				}
				if !got.matches(want) {
					m := fmt.Sprintf("\n%s got\n", name)
					m += got.details()
					m += fmt.Sprintf("\nwanted\n")
					m += want.details()
					t.Fatal(m)
				}
			},
		},
		{
			args: []string{
				path.Join(tmpExercisesDir, testExercises[3].name),
				"-s",
				"2",
				"-e",
				"2",
			},
			check: func(got Exercise, gotErr error) {
				name := "specific exercise, start and end flag"
				want := Exercise{
					name: testExercises[3].name,
					text: "has three lines\n",
				}
				if gotErr != nil {
					t.Fatalf("%s wanted no error, got %s\n", name, gotErr)
				}
				if !got.matches(want) {
					m := fmt.Sprintf("\n%s got\n", name)
					m += got.details()
					m += fmt.Sprintf("\nwanted\n")
					m += want.details()
					t.Fatal(m)
				}
			},
		},
		{
			args: []string{
				path.Join(tmpExercisesDir, testExercises[3].name),
				"-s",
				"2",
				"-e",
				"1",
			},
			check: func(got Exercise, gotErr error) {
				wantErr := errors.New("start flag 2 cannot be greater than end flag 1")
				name := "specific exercise, start and end flag, incorrect output"
				if gotErr == nil {
					t.Fatalf("%s wanted error, got nil\n", name)
				}
				if gotErr.Error() != wantErr.Error() {
					t.Fatalf("%s wanted error msg \"%s\", got \"%s\"", name, wantErr.Error(), gotErr.Error())
				}
			},
		},
		{
			args: []string{},
			check: func(got Exercise, gotErr error) {
				t.Errorf("TODO: standard input")
			},
		},
		{
			args: []string{},
			check: func(got Exercise, gotErr error) {
				t.Errorf("TODO: random exercise, with start and end flags")
			},
		},
		{
			args: []string{},
			check: func(got Exercise, gotErr error) {
				t.Errorf("TODO: random exercise, specific language")
			},
		},
		{
			args: []string{},
			check: func(got Exercise, gotErr error) {
				t.Errorf("TODO: No files in exercise directory")
			},
		},
		{
			args: []string{},
			check: func(got Exercise, gotErr error) {
				t.Errorf("TODO: Empty files in exercise directory")
			},
		},
	}

	for i, tc := range testCases {
		cmd := mockCmd(tc)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("mock command no %d failed to run: %s", i, err)
		}
	}
}
