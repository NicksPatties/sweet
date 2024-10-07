package version

import "testing"

func TestGetVersion(t *testing.T) {
	type testCase struct {
		input string
		want  string
	}

	testCases := []testCase{
		{
			input: "",
			want:  "dev",
		},
		{
			input: "v0.0.1",
			want:  "v0.0.1",
		},
	}

	for _, tc := range testCases {
		got := getVersion(tc.input)
		if got != tc.want {
			t.Fatalf("getVersion: got %s, wanted %s", got, tc.want)
		}
	}
}
