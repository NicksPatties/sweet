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

// Returns the name, text, and their arrays of bytes.
// Useful when printing out exercises when a test fails.
func (ex Exercise) details() (m string) {
	m = fmt.Sprintf("\tname %s\n", ex.name)
	m += fmt.Sprintf("\tname bytes %v\n", []byte(ex.name))
	m += fmt.Sprintf("\ttext %s\n", ex.text)
	m += fmt.Sprintf("\ttext bytes %v\n", []byte(ex.text))
	return
}

// Returs the details of each exercise in a tmp directory.
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

type fromArgsTestCase struct {
	args  []string
	check func(Exercise, error)
}

var mockCmd = func(tc fromArgsTestCase) *cobra.Command {
	cmd := &cobra.Command{
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			got, gotErr := fromArgs(cmd, args)
			tc.check(got, gotErr)
		},
	}
	setRootCmdFlags(cmd)
	cmd.SetArgs(tc.args)
	return cmd
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

	testCases := []fromArgsTestCase{
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
				name := "specific exercise, end flag"
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
			args: []string{
				"-s",
				"1",
				"-e",
				"2",
			},
			check: func(got Exercise, gotErr error) {
				wantErr := errors.New("start and end should not be assigned for random exercise")
				name := "random exercise, start and end flag"
				if gotErr == nil {
					t.Fatalf("%s wanted error, got nil\n", name)
				}
				if gotErr.Error() != wantErr.Error() {
					t.Fatalf("%s wanted error msg \"%s\", got \"%s\"", name, wantErr.Error(), gotErr.Error())
				}
			},
		},
		{
			// WARNING: This test might be flaky if the random
			// selected exercise happens to be the go exercise.
			args: []string{
				"-l",
				"go",
			},
			check: func(got Exercise, gotErr error) {
				name := "random exercise, with language"
				want := testExercises[2]
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
				"-l",
				"ts",
			},
			check: func(got Exercise, gotErr error) {
				name := "random exercise, with language, but not available"
				wantErr := errors.New("failed to find exercise matching language ts")
				if gotErr == nil {
					t.Fatalf("%s wanted error, got nil\n", name)
				}
				if gotErr.Error() != wantErr.Error() {
					t.Fatalf("%s wanted error msg \"%s\", got \"%s\"", name, wantErr.Error(), gotErr.Error())
				}

			},
		},
	}

	for i, tc := range testCases {
		cmd := mockCmd(tc)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("mock command no. %d failed to run: %s", i, err)
		}
	}
}

func TestFromArgsWithStdin(t *testing.T) {
	// Test environment setup
	tmpExercisesDir := t.TempDir()
	t.Setenv("SWEET_EXERCISES_DIR", tmpExercisesDir)

	// Create a tmp file. This will replace os.Stdin
	tmp, err := os.CreateTemp(".", "stdin")
	if err != nil {
		t.Error("Failed to create tmp file")
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	wantText := "Hello from stdin!\n"
	want := Exercise{
		name: path.Join(".", tmp.Name()),
		text: wantText,
	}

	// Write to tmp file
	_, err = tmp.Write([]byte(wantText))
	if err != nil {
		t.Error("Failed to write to tmp file")
	}
	// Go back to the beginning of the file
	tmp.Seek(0, 0)

	// Replace Stdin with the tmp file
	oldStdin := os.Stdin
	os.Stdin = tmp
	defer func() {
		os.Stdin = oldStdin
	}()

	tc := fromArgsTestCase{
		args: []string{"-"},
		check: func(got Exercise, gotErr error) {
			name := "from stdin"
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
	}

	cmd := mockCmd(tc)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("mock command failed to run: %s", err)
	}
}

func TestFromArgsWithEmptyExerciseFiles(t *testing.T) {
	type testCase struct {
		testExercises []Exercise
		args          []string
		check         func(got Exercise, gotErr error)
	}

	testCases := []testCase{
		{
			testExercises: []Exercise{
				{
					name: "empty.txt",
					text: "",
				},
				{
					name: "not-empty.txt",
					text: "Hello!\n",
				},
			},
			args: []string{},
			check: func(got Exercise, gotErr error) {
				name := "blank exercise file when randomly selecting should select a non-empty exercise"
				want := Exercise{
					name: "not-empty.txt",
					text: "Hello!\n",
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
			testExercises: []Exercise{
				{
					name: "empty.txt",
					text: "",
				},
				{
					name: "also-empty.txt",
					text: "",
				},
			},
			args: []string{},
			check: func(got Exercise, gotErr error) {
				name := "multiple blank exercise files when randomly selecting, should exit if there are no remaining files"
				dir := os.Getenv("SWEET_EXERCISES_DIR")
				wantErrMsg := fmt.Sprintf("all files found in the following exercises directory are empty: %s\n", dir)
				if gotErr == nil {
					t.Fatalf("%s wanted error, got nil\n", name)
				}
				if gotErr.Error() != wantErrMsg {

					t.Fatalf("got error msg:\n\t%s\nwanted error msg\n\t%s", gotErr.Error(), wantErrMsg)
				}

			},
		},
		{
			testExercises: []Exercise{},
			args:          []string{},
			check: func(got Exercise, gotErr error) {
				name := "no exercise files, should select a default exercise"
				dir := os.Getenv("SWEET_EXERCISES_DIR")
				if gotErr != nil {
					t.Fatalf("%s wanted no error, got %s\n", name, gotErr.Error())
				}
				if !got.matchesOneOf(defaultExercises) {
					m := fmt.Sprintf("\n%s got\n", name)
					m += got.details()
					m += printExerciseFiles(dir)
					t.Fatal(m)
				}
			},
		},
	}

	for _, tc := range testCases {
		tmpExercisesDir := t.TempDir()
		t.Setenv("SWEET_EXERCISES_DIR", tmpExercisesDir)

		testExercises := tc.testExercises
		createExerciseFiles(t, tmpExercisesDir, testExercises)

		fromArgsTC := fromArgsTestCase{
			args:  tc.args,
			check: tc.check,
		}
		cmd := mockCmd(fromArgsTC)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("mock command failed to run: %s", err)
		}
	}

}
