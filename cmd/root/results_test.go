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

var sortedSpecialCharsEventList string = `2024-10-07 16:1029:26.916: 0 c c 
2024-10-07 16:1029:27.004: 1 o o 
2024-10-07 16:1029:27.095: 2 n n 
2024-10-07 16:1029:27.279: 3 s s 
2024-10-07 16:1029:27.416: 4 o o 
2024-10-07 16:1029:27.667: 5 l l 
2024-10-07 16:1029:27.784: 6 e e 
2024-10-07 16:1029:31.538: 7 d . 
2024-10-07 16:1029:32.243: 8 backspace  
2024-10-07 16:1029:33.216: 7 e . 
2024-10-07 16:1029:33.432: 8 backspace  
2024-10-07 16:1029:33.811: 7 r . 
2024-10-07 16:1029:34.175: 8 backspace  
2024-10-07 16:1029:34.768: 7 . . 
2024-10-07 16:1029:35.313: 8 l l 
2024-10-07 16:1029:35.502: 9 o o 
2024-10-07 16:1029:35.676: 10 g g 
2024-10-07 16:1029:37.565: 11 8 ( 
2024-10-07 16:1029:38.374: 12 backspace  
2024-10-07 16:1029:38.750: 11 9 ( 
2024-10-07 16:1029:39.810: 12 backspace  
2024-10-07 16:1029:41.048: 11 0 ( 
2024-10-07 16:1029:41.380: 12 backspace  
2024-10-07 16:1029:42.058: 11 ( ( 
2024-10-07 16:1029:45.428: 12 2 " 
2024-10-07 16:1029:45.991: 13 backspace  
2024-10-07 16:1029:46.178: 12 3 " 
2024-10-07 16:1029:46.502: 13 backspace  
2024-10-07 16:1029:48.972: 12 " " 
2024-10-07 16:1029:49.427: 13 E E 
2024-10-07 16:1029:50.641: 18 y y 
2024-10-07 16:1029:55.056: 19 4 " 
2024-10-07 16:1029:55.797: 20 backspace  
2024-10-07 16:1029:56.540: 21 " " 
2024-10-07 16:1029:57.101: 22 ) ) 
2024-10-07 16:1029:58.765: 23 enter enter`

var whitespaceEventsList string = `2024-10-07 16:109:16.628: 0 h h 
2024-10-07 16:109:17.177: 1 e e 
2024-10-07 16:109:17.274: 2 y y 
2024-10-07 16:109:18.290: 3 d space 
2024-10-07 16:109:19.222: 4 backspace  
2024-10-07 16:109:20.319: 3 e space 
2024-10-07 16:109:20.837: 4 backspace  
2024-10-07 16:109:21.151: 3 space space 
2024-10-07 16:109:21.487: 4 t t 
2024-10-07 16:109:21.562: 5 h h 
2024-10-07 16:109:21.691: 6 e e 
2024-10-07 16:109:21.832: 7 r r 
2024-10-07 16:109:21.913: 8 e e 
2024-10-07 16:109:22.902: 9 i enter 
2024-10-07 16:109:23.937: 10 backspace  
2024-10-07 16:109:24.429: 9 s enter 
2024-10-07 16:109:25.283: 10 backspace  
2024-10-07 16:109:26.132: 9 enter enter 
2024-10-07 16:109:26.745: 10 t t 
2024-10-07 16:109:28.081: 11 w w 
2024-10-07 16:109:28.148: 12 o o 
2024-10-07 16:109:28.843: 13 space space 
2024-10-07 16:109:29.498: 14 l l 
2024-10-07 16:109:29.537: 15 i i 
2024-10-07 16:109:29.611: 16 n n 
2024-10-07 16:109:29.789: 17 e e 
2024-10-07 16:109:29.831: 18 s s 
2024-10-07 16:109:30.694: 19 enter enter`

var limitedMissesEventsList string = `2024-10-07 16:1046:36.929: 0 q a 
2024-10-07 16:1046:37.331: 1 backspace  
2024-10-07 16:1046:38.067: 0 a a 
2024-10-07 16:1046:39.145: 1 w s 
2024-10-07 16:1046:39.408: 2 backspace  
2024-10-07 16:1046:39.831: 1 s s 
2024-10-07 16:1046:40.827: 2 e d 
2024-10-07 16:1046:41.089: 3 backspace  
2024-10-07 16:1046:41.658: 2 d d 
2024-10-07 16:1046:43.471: 3 r f 
2024-10-07 16:1046:43.942: 4 backspace  
2024-10-07 16:1046:44.862: 3 f f 
2024-10-07 16:1046:46.290: 4 enter enter`

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
		{
			// " = 32, ( = 40, . = 46
			name:   "sorted special characters",
			events: stringToEvents(sortedSpecialCharsEventList),
			want:   "\" (3 times), ( (3 times), . (3 times)",
		},
		{
			name:   "whitespace mistakes",
			events: stringToEvents(whitespaceEventsList),
			want:   "enter (2 times), space (2 times)",
		},
		{
			name:   "only show limited number of misses",
			events: stringToEvents(limitedMissesEventsList),
			want:   "a (1 time), d (1 time), f (1 time)",
		},
	}

	for _, tc := range testCases {
		got := mostMissedKeys(tc.events)
		if got != tc.want {
			t.Errorf("%s:\n\tgot  %s\n\twant %s", tc.name, got, tc.want)
		}
	}
}
