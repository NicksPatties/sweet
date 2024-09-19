package util

import (
	"testing"
)

func TestMakeUsage(t *testing.T) {
	type makeUsageTestCase struct {
		testName       string
		executableName string
		subCommand     string
		usage          string
		want           string
	}

	testCases := []makeUsageTestCase{
		{
			testName:       "no usage",
			executableName: "sweet",
			subCommand:     "version",
			usage:          "",
			want: "Usage: sweet version \n" +
				"For more information, run: sweet help version",
		},
		{
			testName:       "all variables",
			executableName: "./sweet",
			subCommand:     "help",
			usage:          "[sub-command]",
			want: "Usage: ./sweet help [sub-command]\n" +
				"For more information, run: ./sweet help help",
		},
	}

	for _, tc := range testCases {
		got := GetStringFromStdout(
			// Don't forget that MakeUsage returns a function,
			// so I can pass it here
			MakeUsage(tc.executableName, tc.subCommand, tc.usage),
		)

		if got != tc.want {
			t.Errorf("%s:\ngot\n%s\nwant\n%s", tc.testName, got, tc.want)
		}
	}

}
