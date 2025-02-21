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
	// TODO: if this function is indented with tabs, then this test fails
	testText := `func main() {
    fmt.Println("hello!")
}
`
	testModel := exerciseModel{
		name:      "",
		text:      testText,
		typedText: "",
		startTime: time.Time{},
		endTime:   time.Time{},
		quitEarly: false,
		events:    []event.Event{},
	}
	want := testText
	got := testModel.renderText()
	if got != want {
		t.Fatalf("expected\n%s\ngot\n%s", want, got)
	}
}
