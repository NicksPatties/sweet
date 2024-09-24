package util

import "testing"

func TestFilterFileNames(t *testing.T) {
	type testCase struct {
		name      string
		fileNames []string
		language  string
		want      []string
	}

	testCases := []testCase{
		{
			name:      "Happy case",
			fileNames: []string{"one.js", "two.js", "three.go"},
			language:  "js",
			want:      []string{"one.js", "two.js"},
		},
		{
			name:      "No match",
			fileNames: []string{"one.js", "two.go"},
			language:  "c",
			want:      []string{},
		},
		{
			name:      "Has files with no extension",
			fileNames: []string{"one.js", "two"},
			language:  "js",
			want:      []string{"one.js"},
		},
	}

	for _, tc := range testCases {
		got := FilterFileNames(tc.fileNames, tc.language)
		if len(got) != len(tc.want) {
			t.Errorf("%s: wanted length of %d, got length of %d\n", tc.name, len(tc.want), len(got))
		}

		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("%s: wanted result[%d] = %s, got result[%d] = %s\n", tc.name, i, tc.want[i], i, got[i])
			}
		}
	}
}
