package util

import (
	"os"
	"path"
	"testing"
)

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

func TestSweetConfigDir(t *testing.T) {
	t.Run("correct directory is returned", func(t *testing.T) {
		userConfig, _ := os.UserConfigDir()
		expected := path.Join(userConfig, "sweet")

		actual, err := SweetConfigDir()

		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		if expected != actual {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	})
}

func TestMD5Hash(t *testing.T) {
	testString := "what's up?"
	want := "0d0a3c37bc7f0d527b832c6460569d18"
	got := MD5Hash(testString)

	if want != got {
		t.Errorf("\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestLang(t *testing.T) {
	testCases := []struct {
		in   string
		want string
	}{
		{
			in:   "myfile",
			want: "",
		},
		{
			in:   "myfile.py",
			want: "py",
		},
		{
			in:   "myfile.today.go",
			want: "go",
		},
	}

	for _, tc := range testCases {
		got := Lang(tc.in)
		if got != tc.want {
			t.Errorf("got %s, want %s", got, tc.want)
		}
	}
}

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
		// manually verified that version build vars work
	}

	for _, tc := range testCases {
		got := GetVersion()
		if got != tc.want {
			t.Fatalf("getVersion: got %s, wanted %s", got, tc.want)
		}
	}
}

func TestLines(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "default case",
			input: "one\ntwo\nthree\nfour\nfive",
			want: []string{
				"one\n",
				"two\n",
				"three\n",
				"four\n",
				"five",
			},
		},
		{
			name:  "additional newline?",
			input: "one\ntwo\nthree\nfour\nfive\n",
			want: []string{
				"one\n",
				"two\n",
				"three\n",
				"four\n",
				"five\n",
			},
		},
		{
			name:  "empty string should return empty array",
			input: "",
			want:  []string{},
		},
	}

	for _, tc := range testCases {
		got := Lines(tc.input)
		for i := 0; i < len(tc.want); i = i + 1 {
			if len(got) != len(tc.want) {
				t.Errorf("Lengths don't match. Got %d, want %d\n", len(got), len(tc.want))
			}

			if got[i] != tc.want[i] {
				t.Errorf("%d got %s want %s", i, got[i], tc.want[i])
			}
		}
	}
}

func TestDedentLines(t *testing.T) {
	testCases := []struct {
		name  string
		lines []string
		want  []string
	}{
		{
			name: "do nothing if already dedented",
			lines: []string{
				"one\n",
				"two\n",
				"three\n",
			},
			want: []string{
				"one\n",
				"two\n",
				"three\n",
			},
		},
		{
			name: "dedent evenly dedented lines",
			lines: []string{
				"  one\n",
				"  two\n",
				"  three\n",
			},
			want: []string{
				"one\n",
				"two\n",
				"three\n",
			},
		},
		{
			name: "dedent a function",
			lines: []string{
				"  def main:\n",
				"    print('what')",
				"  end\n",
			},
			want: []string{
				"def main:\n",
				"  print('what')",
				"end\n",
			},
		},
		{
			name: "last line is already",
			lines: []string{
				"  def main:\n",
				"    print('what')",
				"  end\n",
				"end\n",
			},
			want: []string{
				"  def main:\n",
				"    print('what')",
				"  end\n",
				"end\n",
			},
		},
	}

	for _, tc := range testCases {
		gotLines := DedentLines(tc.lines)
		for i, got := range gotLines {
			want := tc.want[i]
			if got != want {
				t.Errorf("%s:\ngot\n%s\nwant\n%s", tc.name, got, want)
			}
		}
	}
}
