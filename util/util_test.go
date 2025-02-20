package util

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"
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
func getEventTs(s string) (t time.Time) {
	t, _ = time.Parse(EventTsLayout, s)
	return
}

func TestEventString(t *testing.T) {
	testCases := []struct {
		name string
		in   Event
		want string
	}{
		{
			name: "all fields",
			in: Event{
				Ts:       getEventTs("2024-10-07 13:46:47.679"),
				I:        0,
				Typed:    "a",
				Expected: "b",
			},
			want: "2024-10-07 13:46:47.679\t0\ta\tb",
		},
	}

	for _, tc := range testCases {
		got := fmt.Sprint(tc.in)
		if got != tc.want {
			t.Errorf("%s: got\n\t%s\nwant\n\t%s", tc.name, got, tc.want)
		}

	}

}

func TestParseEvent(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  Event
	}{
		{
			name:  "all fields",
			input: "2024-10-07 13:46:47.679\t0\ta\th",
			want: Event{
				Ts:       getEventTs("2024-10-07 13:46:47.679"),
				I:        0,
				Typed:    "a",
				Expected: "h",
			},
		},
		{
			name:  "backspace",
			input: "2024-10-07 13:46:47.679\t0\tbackspace",
			want: Event{
				Ts:       getEventTs("2024-10-07 13:46:47.679"),
				I:        0,
				Typed:    "backspace",
				Expected: "",
			},
		},
	}

	for _, tc := range testCases {
		got := ParseEvent(tc.input)
		if !got.Matches(tc.want) {
			t.Errorf("%s: got\n%s\n\nwant:\n%s", tc.name, got, tc.want)
		}
	}
}

func TestParseEvents(t *testing.T) {

	testCases := []struct {
		name string
		in   string
		want []Event
	}{
		{
			name: "two events",
			in: "2024-10-07 13:46:47.679\t0\ta\th\n" +
				"2024-10-07 13:46:48.298\t1\tbackspace",
			want: []Event{
				{
					Ts:       getEventTs("2024-10-07 13:46:47.679"),
					I:        0,
					Typed:    "a",
					Expected: "h",
				},
				{
					Ts:       getEventTs("2024-10-07 13:46:48.298"),
					I:        1,
					Typed:    "backspace",
					Expected: "",
				},
			},
		},
	}

	for _, tc := range testCases {
		gotEvents := ParseEvents(tc.in)
		for i, got := range gotEvents {
			if !got.Matches(tc.want[i]) {
				t.Errorf(
					"%s [%d]:\ngot\n  %s\nwant\n  %s",
					tc.name,
					i,
					got,
					tc.want[i],
				)
			}
		}
	}
}
