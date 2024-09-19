package help

import (
	"os"
	"testing"

	"github.com/NicksPatties/sweet/util"
)

func TestRun(t *testing.T) {

	testName := "No sub-commands"
	filename := "sweet_help_want.txt"
	f, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("Error opening file %s", filename)
	}

	want := string(f)

	codeWant := 0
	codeGot := -1
	got := util.GetStringFromStdout(func() {
		args := []string{}
		codeGot = Run(args)
	})

	if got != want {
		t.Errorf("%s: got\n%s\nwant\n%s", testName, got, want)
	}
	if codeWant != codeGot {
		t.Errorf("%s: got error code %d, wanted error code %d", testName, codeGot, codeWant)
	}

}
