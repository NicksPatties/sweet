package help

import (
	"testing"

	"github.com/NicksPatties/sweet/util"
)

func TestRun(t *testing.T) {

	type testCase struct {
		name           string
		args           []string
		executableName string
		wantFilename   string
		codeWant       int
	}
	testCases := []testCase{
		{
			name:           "No sub-commands",
			args:           []string{},
			executableName: "sweet",
			wantFilename:   "sweet_help_want.txt",
			codeWant:       0,
		},
	}

	for _, tc := range testCases {
		want := util.GetWantFile(tc.wantFilename, t)
		var codeGot int

		got := util.GetStringFromStdout(func() {
			codeGot = Run(tc.args, tc.executableName)
		})
		if got != want {
			t.Errorf("%s: got\n%s\nwant\n%s", tc.name, got, want)
		}
		if codeGot != tc.codeWant {
			t.Errorf("%s: got error code %d, wanted error code %d", tc.name, codeGot, tc.codeWant)
		}
	}

}
