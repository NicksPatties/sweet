package root

import (
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path"
	"testing"

	"github.com/NicksPatties/sweet/util"
	"github.com/spf13/cobra"
)

func (got exerciseFile) matches(want exerciseFile) bool {
	return got.name == want.name && got.text == want.text
}

func (got exerciseFile) matchesOneOf(wants []exerciseFile) bool {
	for _, want := range wants {
		if got.matches(want) {
			return true
		}
	}
	return false
}

// Returns the name, text, and their arrays of bytes.
// Useful when printing out exercises when a test fails.
func (ex exerciseFile) details() (m string) {
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
		m += exerciseFile{
			name: name,
			text: string(textBytes),
		}.details()
	}
	return
}

func createExerciseFiles(t *testing.T, dir string, exercises []exerciseFile) {
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

func fileToExercise(t *testing.T, fileName string) (exercise exerciseFile) {
	exercise.name = path.Base(fileName)
	text, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("failed to read exercise file %s\n", fileName)
	}
	exercise.text = string(text)
	return
}

type fromArgsExerciseFileTestCase struct {
	args  []string
	check func(exerciseFile, error)
}

var mockExerciseFileCmd = func(tc fromArgsExerciseFileTestCase) *cobra.Command {
	cmd := &cobra.Command{
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			got, gotErr := exerciseFileFromArgs(cmd, args)
			tc.check(got, gotErr)
		},
	}
	setRootCmdFlags(cmd)
	cmd.SetArgs(tc.args)
	return cmd
}

func Test_exerciseFileFromArgs(t *testing.T) {
	// Test environment setup
	tmpExercisesDir := t.TempDir()
	t.Setenv("SWEET_EXERCISES_DIR", tmpExercisesDir)
	testExercises := []exerciseFile{
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

	testCases := []fromArgsExerciseFileTestCase{
		{
			args: []string{},
			check: func(got exerciseFile, gotErr error) {
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
			check: func(got exerciseFile, gotErr error) {
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
			check: func(got exerciseFile, gotErr error) {
				name := "specific exercise, start flag"
				want := exerciseFile{
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
			check: func(got exerciseFile, gotErr error) {
				name := "specific exercise, end flag"
				want := exerciseFile{
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
			check: func(got exerciseFile, gotErr error) {
				name := "specific exercise, start and end flag"
				want := exerciseFile{
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
			check: func(got exerciseFile, gotErr error) {
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
			check: func(got exerciseFile, gotErr error) {
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
			check: func(got exerciseFile, gotErr error) {
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
			check: func(got exerciseFile, gotErr error) {
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
		cmd := mockExerciseFileCmd(tc)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("mock command no. %d failed to run: %s", i, err)
		}
	}
}

func Test_exerciseFileFromArgs_withStdin(t *testing.T) {
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
	want := exerciseFile{
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

	tc := fromArgsExerciseFileTestCase{
		args: []string{"-"},
		check: func(got exerciseFile, gotErr error) {
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

	cmd := mockExerciseFileCmd(tc)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("mock command failed to run: %s", err)
	}
}

func Test_exerciseFileFromArgs_withEmptyExerciseFiles(t *testing.T) {
	type testCase struct {
		testExercises []exerciseFile
		args          []string
		check         func(got exerciseFile, gotErr error)
	}

	testCases := []testCase{
		{
			testExercises: []exerciseFile{
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
			check: func(got exerciseFile, gotErr error) {
				name := "blank exercise file when randomly selecting should select a non-empty exercise"
				want := exerciseFile{
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
			testExercises: []exerciseFile{
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
			check: func(got exerciseFile, gotErr error) {
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
			testExercises: []exerciseFile{},
			args:          []string{},
			check: func(got exerciseFile, gotErr error) {
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

		fromArgsTC := fromArgsExerciseFileTestCase{
			args:  tc.args,
			check: tc.check,
		}
		cmd := mockExerciseFileCmd(fromArgsTC)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("mock command failed to run: %s", err)
		}
	}

}

type fromArgsViewOptionsTestCase struct {
	args         []string
	exerciseText string
	check        func(*viewOptions, error)
}

var mockViewOptionsCmd = func(tc fromArgsViewOptionsTestCase) *cobra.Command {
	cmd := &cobra.Command{
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, _ []string) {
			got, gotErr := viewOptionsFromArgs(cmd, tc.exerciseText)
			tc.check(got, gotErr)
		},
	}
	setRootCmdFlags(cmd)
	cmd.SetArgs(tc.args)
	return cmd
}

func Test_viewOptionsFromArgs(t *testing.T) {
	testCases := []fromArgsViewOptionsTestCase{
		{
			args:         []string{""},
			exerciseText: "an exercise",
			check: func(got *viewOptions, gotErr error) {
				name := "default window size should be 0"
				want := uint(0)
				if gotErr != nil {
					t.Fatalf("%s wanted no error, got %s\n", name, gotErr)
				}
				if got.windowSize != want {
					t.Fatalf(
						"got window size %d, wanted window size %d",
						got.windowSize,
						want,
					)
				}
			},
		},
		{
			args:         []string{"--window-size", "3"},
			exerciseText: "one\ntwo\nthree\nfour\nfive",
			check: func(got *viewOptions, gotErr error) {
				name := "passing window size of 3"
				want := uint(3)
				if gotErr != nil {
					t.Fatalf("%s wanted no error, got %s\n", name, gotErr)
				}
				if got.windowSize != want {
					t.Fatalf(
						"got window size %d, wanted window size %d",
						got.windowSize,
						want,
					)
				}
			},
		},
		{
			args:         []string{"--window-size", "5"},
			exerciseText: "one\ntwo\nthree\nfour\nfive",
			check: func(got *viewOptions, gotErr error) {
				name := "window size greater or equal to num lines returns 0"
				want := uint(0)
				if gotErr != nil {
					t.Fatalf("%s wanted no error, got %s\n", name, gotErr)
				}
				if got.windowSize != want {
					t.Fatalf(
						"got window size %d, wanted window size %d",
						got.windowSize,
						want,
					)
				}
			},
		},
	}

	for _, tc := range testCases {
		cmd := mockViewOptionsCmd(tc)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("mock command failed to run: %s", err)
		}
	}
}

func Test_scanFileText(t *testing.T) {
	// see sweet.go#setRootCmdFlags
	defaultStart := uint(1)
	var defaultEnd uint = math.MaxUint

	testCases := []struct {
		name     string
		start    uint
		end      uint
		contents string
		want     string
	}{
		{
			name:     "default case",
			start:    defaultStart,
			end:      defaultEnd,
			contents: "This is a test file.",
			want:     "This is a test file.",
		},
		{
			name:     "multiple lines, no flags",
			start:    defaultStart,
			end:      defaultEnd,
			contents: "This is a test file\nthat has two lines!",
			want:     "This is a test file\nthat has two lines!",
		},
		{
			name:     "multiple lines, start only",
			start:    2,
			end:      defaultEnd,
			contents: "one\ntwo\nthree",
			want:     "two\nthree",
		},
		{
			name:     "multiple lines, end only",
			start:    defaultStart,
			end:      2,
			contents: "one\ntwo\nthree",
			want:     "one\ntwo\n",
		},
		{
			name:     "both variables",
			start:    2,
			end:      3,
			contents: "one\ntwo\nthree",
			want:     "two\nthree",
		},
		{
			name:     "both variables, but one line lol",
			start:    2,
			end:      2,
			contents: "one\ntwo\nthree",
			want:     "two\n",
		},
	}

	for _, tc := range testCases {
		// Create a temporary file for testing
		tempFile, err := os.CreateTemp("", "")
		defer tempFile.Close()
		_, err = tempFile.WriteString(tc.contents)
		_, err = tempFile.Seek(0, 0)
		if err != nil {
			t.Fatalf("something went wrong with creating the tmp file")
		}

		got := scanFileText(tempFile, tc.start, tc.end)
		if got != tc.want {
			t.Fatalf("%s\nwant: \n%q\n(%v)\ngot:\n%q\n(%v)", tc.name, tc.contents, util.RenderBytes(tc.contents), got, util.RenderBytes(got))
		}
	}
}
