package root

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

var defaultCaseEventsList string = `2024-10-07 13:1046:47.679: 0 a h
2024-10-07 13:1046:48.298: 1 backspace
2024-10-07 13:1046:49.442: 0 h h
2024-10-07 13:1046:51.160: 1 e e
2024-10-07 13:1046:52.781: 2 i y
2024-10-07 13:1046:53.316: 3 backspace
2024-10-07 13:1046:54.688: 2 k y
2024-10-07 13:1046:55.262: 3 backspace
2024-10-07 13:1046:55.997: 2 y y
2024-10-07 13:1046:56.521: 3 enter enter`

func stringToEvent(line string) (e event) {
	s := strings.Split(line, ": ")
	e.ts, _ = time.Parse("2006-01-02 15:14:05.000", s[0])
	s = strings.Split(s[1], " ")
	e.i, _ = strconv.Atoi(s[0])
	e.typed = s[1]
	if len(s) > 2 {
		e.expected = s[2]
	}
	return
}

func stringToEvents(list string) (events []event) {
	for _, line := range strings.Split(list, "\n") {
		if line != "" {
			events = append(events, stringToEvent(line))
		}
	}
	return
}

func TestMostMissedKeys(t *testing.T) {

	type testCase struct {
		name   string
		events []event
		want   string
	}

	testCases := []testCase{
		{
			name:   "default case",
			events: stringToEvents(defaultCaseEventsList),
			want:   "y (2 times), h (1 time)",
		},
		// {
		// 	name: "",
		// 	events: stringToevents()
		// 	want: "\" (6 times), s (4 times), ; (3 times)",
		// }
	}

	for _, tc := range testCases {
		got := mostMissedKeys(tc.events)
		if got != tc.want {
			t.Errorf("%s:\n\tgot  %s\n\twant %s", tc.name, got, tc.want)
		}
	}
}
