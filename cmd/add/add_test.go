package add

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/cobra"
)

type addOptionsTestCase struct {
	name         string
	args         []string
	fileContents string
	wantErr      bool
	wantAdded    string
}

var mockAddCmd = func(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addExercise(cmd, args)
		},
	}
	setAddCmdFlags(cmd)
	cmd.SetArgs(args)
	return cmd
}

func TestAddCmd(t *testing.T) {
	testFileName := "add-me-please"
	testCases := []addOptionsTestCase{
		{
			name:         "no flags",
			args:         []string{testFileName},
			fileContents: "one\ntwo\nthree\nfour\nfive",
			wantAdded:    "one\ntwo\nthree\nfour\nfive",
		},
		{
			name:         "start flag provided",
			args:         []string{testFileName, "-s", "2"},
			fileContents: "one\ntwo\nthree\nfour\nfive",
			wantAdded:    "two\nthree\nfour\nfive",
		},
		{
			name:         "end flag provided",
			args:         []string{testFileName, "-e", "2"},
			fileContents: "one\ntwo\nthree\nfour\nfive",
			wantAdded:    "one\ntwo\n",
		},
		{
			name:         "both flags provided",
			args:         []string{testFileName, "-s", "2", "-e", "4"},
			fileContents: "one\ntwo\nthree\nfour\nfive",
			wantAdded:    "two\nthree\nfour\n",
		},
		{
			name:         "don't remove ending newline",
			args:         []string{testFileName},
			fileContents: "one\ntwo\nthree\nfour\nfive\n",
			wantErr:      false,
			wantAdded:    "one\ntwo\nthree\nfour\nfive\n",
		},
		{
			name:         "start cannot be greater than end",
			args:         []string{testFileName, "-s", "3", "-e", "2"},
			fileContents: "one\ntwo\nthree\nfour\nfive\n",
			wantErr:      true,
			wantAdded:    "",
		},
		{
			name:         "dedent input if it's already indented",
			args:         []string{testFileName},
			fileContents: "  one\n  two\n  three\n",
			wantErr:      false,
			wantAdded:    "one\ntwo\nthree\n",
		},
	}

	for _, tc := range testCases {
		tmpExercisesDir := t.TempDir()
		t.Setenv("SWEET_EXERCISES_DIR", tmpExercisesDir)

		// create the test file
		testFile, err := os.Create(testFileName)
		defer os.Remove(testFileName)
		if err != nil {
			t.Fatal("failed to create exercise file")
		}
		testFile.WriteString(tc.fileContents)
		err = testFile.Close()
		if err != nil {
			t.Fatalf("failed to close temporary file %s", testFile.Name())
		}

		cmd := mockAddCmd(tc.args)
		err = cmd.Execute()
		if tc.wantErr && err == nil {
			t.Fatalf("%s: wanted error, got nil", tc.name)
		}

		if !tc.wantErr {
			if err != nil {
				t.Fatalf("%s: wanted nil, got error: %s", tc.name, err)
			}

			// read contents from new created tmp file
			gotBytes, err := os.ReadFile(path.Join(tmpExercisesDir, testFileName))
			if err != nil {
				t.Fatalf("%s\nfailed to open expected file. %s", tc.name, err)
			}

			if got := string(gotBytes); got != tc.wantAdded {
				t.Errorf("%s:\nExpected to add:\n%s\n\n Added:\n%s", tc.name, tc.wantAdded, got)
			}
		}
	}
}

func TestAddCmd_addingSameFile(t *testing.T) {

	tmpExercisesDir := t.TempDir()
	t.Setenv("SWEET_EXERCISES_DIR", tmpExercisesDir)

	// create the test file
	testFileName := "add-me-please"
	testFile, err := os.Create(testFileName)
	defer os.Remove(testFile.Name())
	if err != nil {
		t.Fatal("failed to create exercise file")
	}
	testFile.WriteString("file\ncontents\n")

	cmd := mockAddCmd([]string{testFile.Name()})
	err = cmd.Execute()
	err = cmd.Execute()

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
