package main

import "testing"

func TestRun(t *testing.T) {
	// default executable name
	dex := "sweet"
	mockCommands := Commands{
		exercise: func(string, string, string) int {
			return 0
		},
		help: func([]string) int {
			return 0
		},
		version: func([]string) int {
			return 0
		},
	}
	type testCase struct {
		name           string
		executableName string
		args           []string
		want           int
	}
	testCases := []testCase{
		{
			// TODO: I may need to worry about mocking the
			// commands in case they capture the
			// user's input
			name:           "sweet command with no flags",
			executableName: dex,
			args:           []string{},
			want:           0,
		},
		{
			name:           "sweet command with valid flags",
			executableName: dex,
			args:           []string{"-l", "go"},
			want:           0,
		},
		// The "sweet command with invalid flags" case is already
		// handled by the flag module's behavior,
		// so it's skipped here.
		{
			name:           "help happy path",
			executableName: dex,
			args:           []string{"help"},
			want:           0,
		},
		{
			name:           "version happy path",
			executableName: dex,
			args:           []string{"version"},
			want:           0,
		},
		// TODO: Implement "about" path
		{
			name:           "incorrect sub-command",
			executableName: dex,
			args:           []string{"what"},
			want:           1,
		},
	}
	for _, tc := range testCases {
		got := Run(tc.executableName, tc.args, mockCommands)
		if got != tc.want {
			t.Errorf("%s: got: %d wanted: %d", tc.name, got, tc.want)
		}
	}
}
