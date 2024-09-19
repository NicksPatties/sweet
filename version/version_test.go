package version

import (
	"testing"

	"github.com/NicksPatties/sweet/util"
)

func TestRun(t *testing.T) {
	// default executable
	defaulExe := "sweet"
	defaultVersion := "version"

	type testCase struct {
		name           string
		args           []string
		executableName string
		version        string
		want           string
		codeWant       int
	}
	testCases := []testCase{
		{
			name:           "No sub-commands",
			args:           []string{},
			executableName: defaulExe,
			version:        defaultVersion,
			want:           defaultVersion,
			codeWant:       0,
		},
		{
			name:           "With sub commands",
			args:           []string{"but", "with", "more", "commands"},
			executableName: defaulExe,
			version:        defaultVersion,
			want:           "Error: Too many arguments\n" + util.MakeUsageString(defaulExe, defaultVersion, ""),
			codeWant:       1,
		},
	}

	for _, tc := range testCases {
		var codeGot int
		got := util.GetStringFromStdout(func() {
			codeGot = Run(tc.args, tc.executableName, tc.version)
		})
		if got != tc.want {
			t.Errorf("%s: got\n%s,\nwant\n%s", tc.name, got, tc.want)
		}
		if codeGot != tc.codeWant {
			t.Errorf("%s: got error code %d, wanted error code %d", tc.name, codeGot, tc.codeWant)
		}

	}

}
