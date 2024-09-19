package help

import (
	"os"
	"testing"

	"github.com/NicksPatties/sweet/util"
)

func TestRun(t *testing.T) {

	type testCase struct {
		name         string
		args         []string
		wantFilename string
		codeWant     int
	}
	testCases := []testCase{
		{
			name:         "No sub-commands",
			args:         []string{},
			wantFilename: "sweet_help_want.txt",
			codeWant:     0,
		},
	}

	for _, tc := range testCases {
		f, err := os.ReadFile(tc.wantFilename)
		if err != nil {
			t.Errorf("Error opening file %s", tc.wantFilename)
		}
		want := string(f)
		var codeGot int

		got := util.GetStringFromStdout(func() {
			args := []string{}
			codeGot = Run(args)
		})
		if got != want {
			t.Errorf("%s: got\n%s\nwant\n%s", tc.name, got, want)
		}
		if codeGot != tc.codeWant {
			t.Errorf("%s: got error code %d, wanted error code %d", tc.name, codeGot, tc.codeWant)
		}

	}

}
