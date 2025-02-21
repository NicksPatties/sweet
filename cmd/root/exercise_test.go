package root

import (
	"testing"
	"time"

	consts "github.com/NicksPatties/sweet/constants"
	"github.com/NicksPatties/sweet/event"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var mockViewOptions = &viewOptions{
	styles:     defaultStyles(),
	windowSize: 0,
}

func Test_renderName(t *testing.T) {
	testModel := exerciseModel{
		name:        "exercise.go",
		text:        "",
		typedText:   "",
		startTime:   time.Time{},
		endTime:     time.Time{},
		quitEarly:   false,
		events:      []event.Event{},
		viewOptions: mockViewOptions,
	}
	want := "// exercise.go"
	got := testModel.renderName()
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func Test_renderText(t *testing.T) {
	// NOTE: if this function is indented with tabs, then this test fails
	testText := `func main() {
    fmt.Println("hello!")
}
`
	tt := []struct {
		name  string
		text  string
		typed string
		want  string
	}{
		{
			name:  "default test",
			text:  testText,
			typed: "",
			want:  testText,
		},
		{
			name:  "add newline character at end of line",
			text:  testText,
			typed: "func main() {",
			want: `func main() {↲
    fmt.Println("hello!")
}
`},
		{
			name:  "show arrow if there's an error on a newline",
			text:  testText,
			typed: "func main() {a",
			want: `func main() {↲
    fmt.Println("hello!")
}
`},
	}

	for _, test := range tt {
		testModel := exerciseModel{
			name:        "",
			text:        test.text,
			typedText:   test.typed,
			startTime:   time.Time{},
			endTime:     time.Time{},
			quitEarly:   false,
			events:      []event.Event{},
			viewOptions: mockViewOptions,
		}
		got := testModel.renderText()
		if got != test.want {
			t.Fatalf("%s failed\nexpected\n%s\ngot\n%s",
				test.name, test.want, got,
			)
		}
	}
}

func Test_renderText_cursorPosition(t *testing.T) {
	old := lg.ColorProfile()
	lg.SetColorProfile(termenv.TrueColor)
	defer lg.SetColorProfile(old)

	testViewOptions := &viewOptions{
		windowSize: 0,
		styles: styles{
			commentStyle: lg.NewStyle(),
			untypedStyle: lg.NewStyle(),
			cursorStyle:  lg.NewStyle().Foreground(lg.Color("1")),
			typedStyle:   lg.NewStyle(),
			mistakeStyle: lg.NewStyle(),
		},
	}
	escStart := "\033[31m"
	escEnd := "\033[0m"

	testCases := []struct {
		testName string
		text     string
		typed    string
		want     string
	}{
		{
			testName: "single line",
			text:     "asdf",
			typed:    "as",
			want:     "as" + escStart + "d" + escEnd + "f",
		},
		{
			testName: "multiple lines",
			text:     "def main:\n  print('hello')\n",
			typed:    "def main:\n  ",
			want:     "def main:\n  " + escStart + "p" + escEnd + "rint('hello')\n",
		},
	}

	for _, tc := range testCases {
		testModel := exerciseModel{
			name:        "",
			text:        tc.text,
			typedText:   tc.typed,
			startTime:   time.Time{},
			endTime:     time.Time{},
			quitEarly:   false,
			events:      []event.Event{},
			viewOptions: testViewOptions,
		}
		got := testModel.renderText()
		if got != tc.want {
			t.Fatalf("%s\ngot\n%s\nwant\n%s", tc.testName, got, tc.want)
		}
	}

}

func Test_addRuneToTypedText(t *testing.T) {
	tt := []struct {
		name      string
		text      string
		typed     string
		typedRune rune
		want      string
	}{
		{
			name:      "happy case",
			text:      "asdf",
			typed:     "",
			typedRune: 'a',
			want:      "a",
		},
		{
			name:      "ignore if typed text is the same length of text, but is incorrect",
			text:      "asdf",
			typed:     "asdq",
			typedRune: 'a',
			want:      "asdq",
		},
		{
			name: "adding a newline also adds whitespace up to rune",
			text: `def main:
  print("hey")
`,
			typed:     "def main:",
			typedRune: consts.Enter,
			want:      "def main:\n  ", // two whitespace indentation
		},
	}

	for _, test := range tt {
		testModel := exerciseModel{
			name:        "",
			text:        test.text,
			typedText:   test.typed,
			startTime:   time.Time{},
			endTime:     time.Time{},
			quitEarly:   false,
			events:      []event.Event{},
			viewOptions: mockViewOptions,
		}
		testModel = testModel.addRuneToTypedText(test.typedRune)
		if testModel.typedText != test.want {
			t.Fatalf("want %s, got %s", test.want, testModel.typedText)
		}
	}
}

func Test_deleteRuneFromTypedText(t *testing.T) {
	tt := []struct {
		name  string
		text  string
		typed string
		want  string
	}{
		{
			name:  "happy case",
			text:  "asdf",
			typed: "a",
			want:  "",
		},
		{
			name:  "no typed text yet",
			text:  "asdf",
			typed: "",
			want:  "",
		},
		{
			name: "remove all whitespace including the newline",
			text: `def main:
  print("hey")
`,
			typed: "def main:\n  ",
			want:  "def main:",
		},
	}

	for _, test := range tt {
		testModel := exerciseModel{
			name:        "",
			text:        test.text,
			typedText:   test.typed,
			startTime:   time.Time{},
			endTime:     time.Time{},
			quitEarly:   false,
			events:      []event.Event{},
			viewOptions: mockViewOptions,
		}
		testModel = testModel.deleteRuneFromTypedText()
		if testModel.typedText != test.want {
			t.Fatalf("want\n%s\ngot\n%s\n", test.want, testModel.typedText)
		}
	}
}

func Test_finished(t *testing.T) {
	var tt = []struct {
		name  string
		text  string
		typed string
		want  bool
	}{
		{
			name:  "finished",
			text:  "asdf",
			typed: "asdf",
			want:  true,
		},
		{
			name:  "not finished: didn't type enough characters",
			text:  "asdf",
			typed: "asd",
			want:  false,
		},
		{
			name:  "not finished: last character is wrong",
			text:  "asdf",
			typed: "asdq",
			want:  false,
		},
	}

	for _, test := range tt {
		testModel := exerciseModel{
			name:        "",
			text:        test.text,
			typedText:   test.typed,
			startTime:   time.Time{},
			endTime:     time.Time{},
			quitEarly:   false,
			events:      []event.Event{},
			viewOptions: mockViewOptions,
		}

		want := test.want
		got := testModel.finished()
		if got != want {
			t.Fatalf("want %t, got %t", want, got)
		}
	}
}
