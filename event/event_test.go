package event

import (
	"fmt"
	"testing"
	"time"
)

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
