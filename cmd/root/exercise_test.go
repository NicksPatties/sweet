package root

import (
	"testing"
	"time"

	consts "github.com/NicksPatties/sweet/constants"
	"github.com/NicksPatties/sweet/event"
	"github.com/NicksPatties/sweet/util"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var renderBytes = util.RenderBytes

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
			name: "do not render the last newline",
			text: `f()
  a()
end
`,
			typed: "",
			want: `f()
  a()
end`,
		},
		{
			name:  "add newline character at end of line",
			text:  testText,
			typed: "func main() {",
			want: `func main() {↲
    fmt.Println("hello!")
}`},
		{
			name:  "show arrow if there's a mistake on a newline",
			text:  testText,
			typed: "func main() {a",
			want: `func main() {↲
    fmt.Println("hello!")
}`},
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
			t.Fatalf("%s failed\nwant\n%s\n%q\ngot\n%s\n%q\n",
				test.name, test.want, test.want, got, got,
			)
		}
	}
}

func red(s string) string {
	escStart := "\033[31m"
	escEnd := "\033[0m"
	return escStart + s + escEnd
}

func Test_renderText_cursorPosition(t *testing.T) {
	oldProfile := lg.ColorProfile()
	lg.SetColorProfile(termenv.TrueColor)
	defer lg.SetColorProfile(oldProfile)

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

	testCases := []struct {
		testName string
		text     string
		typed    string
		want     string
	}{
		{
			testName: "single line: first character",
			text:     "asdf",
			typed:    "",
			want:     red("a") + "sdf",
		},
		{
			testName: "single line",
			text:     "asdf",
			typed:    "as",
			want:     "as" + red("d") + "f",
		},
		{
			testName: "multiple lines",
			text:     "def main:\n  print('hello')\nfunc yeah",
			typed:    "def main:\n  ",
			want:     "def main:\n  " + red("p") + "rint('hello')\nfunc yeah",
		},
		{
			testName: "multiple lines: blank line",
			text:     "#!/bin/bash\n\necho hello",
			typed:    "#!/bin/bash\n",
			want:     "#!/bin/bash\n" + red(consts.Arrow) + "\necho hello",
		},
		{
			testName: "multiple lines: last line",
			text:     "def main:\n  print('hello')\nfunc yeah",
			typed:    "def main:\n  print('hello')\n",
			want:     "def main:\n  print('hello')\n" + red("f") + "unc yeah",
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
			t.Errorf("%s\ngot\n%v\n%s\nwant\n%v\n%s", tc.testName, got, renderBytes(got), tc.want, renderBytes(tc.want))
		}
	}
}

func reds(s string) string {
	finished := ""
	for _, r := range s {
		finished += red(string(r))
	}
	return finished
}

func Test_renderText_typedAndUntyped(t *testing.T) {
	oldProfile := lg.ColorProfile()
	lg.SetColorProfile(termenv.TrueColor)
	defer lg.SetColorProfile(oldProfile)

	testViewOptions := &viewOptions{
		windowSize: 0,
		styles: styles{
			commentStyle: lg.NewStyle(),
			untypedStyle: lg.NewStyle(),
			cursorStyle:  lg.NewStyle(),
			typedStyle:   lg.NewStyle().Foreground(lg.Color("1")),
			mistakeStyle: lg.NewStyle(),
		},
	}

	testCases := []struct {
		testName string
		text     string
		typed    string
		want     string
	}{
		{
			testName: "partially typed line",
			text:     "asdf",
			typed:    "as",
			want:     reds("as") + "df",
		},
		{
			testName: "fully typed line",
			text:     "asdf",
			typed:    "asdf",
			want:     reds("asdf"),
		},
		// NOTE: test cases with newlines will fail
		// because of color reset escape sequences are
		// placed before newlines. You should
		// **visually test these cases!!**
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
			t.Fatalf("%s\ngot\n%v\n%s\nwant\n%v\n%s", tc.testName, got, renderBytes(got), tc.want, renderBytes(tc.want))
		}
	}
}

func Test_renderText_windowSize(t *testing.T) {
	oldProfile := lg.ColorProfile()
	lg.SetColorProfile(termenv.TrueColor)
	defer lg.SetColorProfile(oldProfile)
	blankStyles := styles{
		commentStyle:         lg.NewStyle(),
		untypedStyle:         lg.NewStyle(),
		cursorStyle:          lg.NewStyle().Foreground(lg.Color("1")), // red
		typedStyle:           lg.NewStyle(),
		mistakeStyle:         lg.NewStyle(),
		vignetteStyle:        lg.NewStyle(),
		vignetteMistakeStyle: lg.NewStyle(),
	}
	mockText := "one\ntwo\nthree\nfour\nfive"
	testCases := []struct {
		name       string
		windowSize uint
		text       string
		typed      string
		want       string
	}{
		{
			name:       "zero windowSize should show the entire exercise",
			windowSize: 0,
			text:       mockText,
			typed:      "",
			want:       red(string(mockText[0])) + mockText[1:],
		},
		{
			name:       "should only show one line",
			windowSize: 1,
			text:       mockText,
			typed:      "",
			want:       red("o") + "ne",
		},
		{
			name:       "two lines: start of exercise",
			windowSize: 2,
			text:       mockText,
			typed:      "",
			want:       red("o") + "ne\ntwo",
		},
		{
			name:       "two lines: partway through",
			windowSize: 2,
			text:       mockText,
			typed:      "one\n",
			want:       red("t") + "wo\nthree",
		},
		{
			name:       "two lines: last line",
			windowSize: 2,
			text:       mockText,
			typed:      "one\ntwo\nthree\nfour\n",
			want:       "four\n" + red("f") + "ive",
		},
		{
			name:       "three lines: exercise start",
			windowSize: 3,
			text:       mockText,
			typed:      "",
			want:       red("o") + "ne\ntwo\nthree",
		},
		{
			name:       "three lines: partway",
			windowSize: 3,
			text:       mockText,
			typed:      "one\ntwo\n",
			want:       "two\n" + red("t") + "hree\nfour",
		},
		{
			name:       "three lines: end",
			windowSize: 3,
			text:       mockText,
			typed:      "one\ntwo\nthree\nfour\n",
			want:       "three\nfour\n" + red("f") + "ive",
		},
		{
			name:       "four lines: end",
			windowSize: 4,
			text:       mockText,
			typed:      "one\ntwo\nthree\nfour\n",
			want:       "two\nthree\nfour\n" + red("f") + "ive",
		},
		{
			name:       "four lines: prevent content shift with final newline",
			windowSize: 4,
			text:       mockText + "\n",
			typed:      "one\ntwo\nthree\nfour\n",
			want:       "two\nthree\nfour\n" + red("f") + "ive",
		},
	}

	for _, tc := range testCases {
		mockViewOptions := &viewOptions{
			styles:     blankStyles,
			windowSize: tc.windowSize,
		}
		mockModel := exerciseModel{
			name:        "",
			text:        tc.text,
			typedText:   tc.typed,
			startTime:   time.Time{},
			endTime:     time.Time{},
			quitEarly:   false,
			events:      []event.Event{},
			viewOptions: mockViewOptions,
		}
		got := mockModel.renderText()

		if got != tc.want {
			t.Fatalf("%s\nwindowSize: %d\ngot:\n%s\nwant:\n%s",
				tc.name, tc.windowSize, got, tc.want)
		}
	}
}

func Test_renderText_mistakes(t *testing.T) {
	// TODO test case where I make a mistake on a newline
	// If I make a mistake on a newline, then I have to let
	// the next line know to render a cursor on the next line,
	// which is totally cancerous.
	//
	// This may be why I counted indeces instead of splitting
	// my typed lines by newlines, since they may not be there!
}

func Test_renderLine(t *testing.T) {
	oldProfile := lg.ColorProfile()
	lg.SetColorProfile(termenv.TrueColor)
	defer lg.SetColorProfile(oldProfile)

	testCaseStyles := styles{
		commentStyle:         lg.NewStyle().Reverse(true),
		untypedStyle:         lg.NewStyle().Reverse(true),
		cursorStyle:          lg.NewStyle().Reverse(true),
		typedStyle:           lg.NewStyle().Reverse(true),
		mistakeStyle:         lg.NewStyle().Reverse(true),
		vignetteStyle:        lg.NewStyle(),
		vignetteMistakeStyle: lg.NewStyle().Foreground(lg.Color("1")),
	}

	testCases := []struct {
		name      string
		text      string
		typed     string
		style     styles
		vignette  bool
		hasCursor bool
		want      string
	}{
		{
			name:     "vignette correctly",
			text:     "my text",
			typed:    "",
			style:    testCaseStyles,
			vignette: true,
			want:     "my text",
		},
		{
			name:     "vignette with mistakes",
			text:     "my text",
			typed:    "my next",
			style:    testCaseStyles,
			vignette: true,
			want:     "my " + red("t") + "ext",
		},
	}
	for _, tc := range testCases {
		typed := &tc.typed
		if tc.typed == "" {
			typed = nil
		}
		got := renderLine(tc.text, typed, tc.style, tc.vignette, false)
		want := tc.want
		if got != want {
			t.Fatalf("%s\ngot:  %s\nwant: %s", tc.name, got, want)
		}
	}
}

func Test_lines(t *testing.T) {
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
		got := lines(tc.input)
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

func Test_typedLines(t *testing.T) {
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
		got := typedLines(tc.lines, tc.typed)
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

func Test_currentLine(t *testing.T) {
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
		lines := lines(tc.text)
		got := currentLine(lines, tc.typed)
		if got != tc.want {
			t.Errorf("%s: got: %d, want: %d", tc.name, got, tc.want)
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
