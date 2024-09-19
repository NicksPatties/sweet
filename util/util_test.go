package util

import (
	"bytes"
	"testing"
)

type makeUsageTestCase struct {
	testName       string
	executableName string
	subCommand     string
	usage          string
	want           string
}

func TestMakeUsage(t *testing.T) {

	testCases := []makeUsageTestCase{
		{
			testName:       "no usage",
			executableName: "sweet",
			subCommand:     "version",
			usage:          "",
			want: "Usage: sweet version \n" +
				"For more information, run: sweet help version",
		},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		// Remember to call the function returned by MakeUsage!
		MakeUsage(&buf, tc.executableName, tc.subCommand, tc.usage)()
		got := buf.String()

		if got != tc.want {
			t.Errorf("%s:\ngot\n%s\nwant\n%s", tc.testName, got, tc.want)
		}
	}

}
