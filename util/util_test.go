package util

import (
	"os"
	"path"
	"testing"

	lg "github.com/charmbracelet/lipgloss"
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
			name:  "empty string",
			input: "",
			want:  []string{},
		},
	}

	for _, tc := range testCases {
		got := Lines(tc.input)
		for i := 0; i < len(tc.want); i = i + 1 {
			if len(got) != len(tc.want) {
				t.Fatalf("Lengths don't match. Got %d, want %d\n", len(got), len(tc.want))
			}

			if got[i] != tc.want[i] {
				t.Errorf("%d got %s want %s", i, got[i], tc.want[i])
			}
		}
	}
}

func TestTypedLines(t *testing.T) {
	testCases := []struct {
		name  string
		lines []string
		typed string
		want  []string
	}{
		{
			name: "typed newline before end of line",
			lines: []string{
				"one\n",
				"two\n",
				"three\n",
			},
			typed: "on\n",
			want:  []string{"on\n"},
		},
		{
			name: "haven't started the exercise yet",
			lines: []string{
				"one\n",
				"two\n",
				"three\n",
			},
			typed: "",
			want:  []string{},
		},
		{
			name: "finished first line of the exercise",
			lines: []string{
				"one\n",
				"two\n",
				"three\n",
			},
			typed: "one\n",
			want:  []string{"one\n"},
		},
		{
			name: "finished two lines of the exercise",
			lines: []string{
				"one\n",
				"two\n",
				"three\n",
			},
			typed: "one\ntw",
			want:  []string{"one\n", "tw"},
		},
	}

	for _, tc := range testCases {
		got := TypedLines(tc.lines, tc.typed)
		for i, wantStr := range tc.want {
			if len(got) != len(tc.want) {
				t.Errorf("len got: %d len(want): %d", len(got), len(tc.want))
			}
			if got[i] != wantStr {
				t.Errorf("got: %s\n want: %s", got[i], wantStr)
			}
		}
	}
}

func TestCurrentLine(t *testing.T) {
	testCases := []struct {
		name  string
		text  string
		typed string
		want  int
	}{
		{
			name:  "no typed text",
			text:  "one\ntwo\nthree",
			typed: "",
			want:  0,
		},
		{
			name:  "one line of typed text",
			text:  "one\ntwo\nthree",
			typed: "one\n",
			want:  1,
		},
	}

	for _, tc := range testCases {
		lines := Lines(tc.text)
		got := CurrentLine(lines, tc.typed)
		if got != tc.want {
			t.Errorf("%s: got: %d, want: %d", tc.name, got, tc.want)
		}
	}
}

func TestRemoveLastNewline(t *testing.T) {
	style := lg.NewStyle().Foreground(lg.Color("8"))
	testCases := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "default",
			str:  "I am a string\n",
			want: "I am a string",
		},
		{
			name: "in the middle",
			str:  "I am a string\nin the middle",
			want: "I am a stringin the middle",
		},
		{
			name: "with styles",
			str:  "I am a string" + style.Render("\n") + "in the middle",
			want: "I am a string" + style.Render("") + "in the middle",
		},
		{
			name: "no newline",
			str:  "I am a string",
			want: "I am a string",
		},
	}

	for _, tc := range testCases {
		got := RemoveLastNewline(tc.str)
		if tc.want != got {
			t.Errorf("want\n%s\n%q\ngot\n%s\n%q", tc.want, tc.want, got, got)
		}
	}
}
