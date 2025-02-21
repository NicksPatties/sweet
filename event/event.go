package event

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// A recording of a keypress during the exercise.
//
// These are used to perform analysis on the user's performance,
// display stats, and keys that were causing the most trouble.
type Event struct {
	// The moment the event took place.
	Ts time.Time

	// The key that was Typed.
	Typed string

	// The rune that was Expected. Optional, since the user
	// may have pressed backspace.
	Expected string

	// The index of the exercise when the rune was typed.
	I int
}

const EventTsLayout = "2006-01-02 15:04:05.000"

// Converts an event to a string.
func (e Event) String() string {
	time := e.Ts.Format(EventTsLayout)
	return fmt.Sprintf("%s\t%d\t%s\t%s", time, e.I, e.Typed, e.Expected)
}

// Checks if an event has the same timestamp, index, typed
// and expected characters. Used primarily for testing.
func (a Event) Matches(b Event) bool {
	return a.Ts.Equal(b.Ts) &&
		a.I == b.I &&
		a.Typed == b.Typed &&
		a.Expected == b.Expected
}

// Converts an event string to an event struct.
func ParseEvent(line string) (e Event) {
	s := strings.Split(line, "\t")
	e.Ts, _ = time.Parse(EventTsLayout, s[0])
	e.I, _ = strconv.Atoi(s[1])
	e.Typed = s[2]
	if len(s) > 3 {
		e.Expected = s[3]
	}
	return
}

type Events []Event

// Same as above, but for a multi-line list of events.
func ParseEvents(list string) (events Events) {
	for _, line := range strings.Split(list, "\n") {
		if line != "" {
			events = append(events, ParseEvent(line))
		}
	}
	return
}

// Returns a string of an array of events.
func (events Events) String() (s string) {
	s += fmt.Sprintln("[")
	for _, e := range events {
		s += fmt.Sprintf("  %s\n", e)
	}
	s += fmt.Sprintln("]")
	return
}

// Creates a new event. Should be used when recording a keystroke
// to the model.
func NewEvent(typed string, expected string, i int) Event {
	return Event{
		Ts:       time.Now(),
		Typed:    typed,
		Expected: expected,
		I:        i,
	}
}

// Converts a bubbletea key message to a string.
// Used to properly record key events.
func TeaKeyMsgToEventTyped(msg tea.KeyMsg) string {
	switch msg.Type {
	case tea.KeyEnter:
		return "enter"
	case tea.KeyBackspace:
		return "backspace"
	case tea.KeySpace:
		return "space"
	default:
		return string(msg.Runes[0])
	}
}

func RuneToEventExpected(r rune) string {
	switch r {
	case '\n':
		return "enter"
	case ' ':
		return "space"
	default:
		return string(r)
	}
}
