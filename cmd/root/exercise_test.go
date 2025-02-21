package root

import (
	"testing"
	"time"

	"github.com/NicksPatties/sweet/event"
)

func TestRenderName(t *testing.T) {
	testModel := exerciseModel{
		name:      "exercise.go",
		text:      "",
		typedText: "",
		startTime: time.Time{},
		endTime:   time.Time{},
		quitEarly: false,
		events:    []event.Event{},
	}
	want := "// exercise.go"
	got := testModel.renderName()
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestRenderText(t *testing.T) {
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
	}

	for _, test := range tt {
		// TODO: if this function is indented with tabs, then this test fails
		testModel := exerciseModel{
			name:      "",
			text:      test.text,
			typedText: test.typed,
			startTime: time.Time{},
			endTime:   time.Time{},
			quitEarly: false,
			events:    []event.Event{},
		}
		got := testModel.renderText()
		if got != test.want {
			t.Fatalf("expected\n%s\ngot\n%s", test.want, got)
		}

	}
}
