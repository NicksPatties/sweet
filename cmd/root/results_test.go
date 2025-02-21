package root

import (
	"fmt"
	"testing"
	"time"

	"github.com/NicksPatties/sweet/event"
)

// mistakes:          3
// total key presses: 7
// typed:             hey↲
// expected:          hey↲
var defaultCaseEventsList string = "2024-10-07 13:46:47.679\t0\ta\th\n" + // miss
	"2024-10-07 13:46:48.298\t1\tbackspace\n" +
	"2024-10-07 13:46:49.442\t0\th\th\n" +
	"2024-10-07 13:46:51.160\t1\te\te\n" +
	"2024-10-07 13:46:52.781\t2\ti\ty\n" + // miss
	"2024-10-07 13:46:53.316\t3\tbackspace\n" +
	"2024-10-07 13:46:54.688\t2\tk\ty\n" + // miss
	"2024-10-07 13:46:55.262\t3\tbackspace\n" +
	"2024-10-07 13:46:55.997\t2\ty\ty\n" +
	"2024-10-07 13:46:56.521\t3\tenter\tenter"

func TestAccuracy(t *testing.T) {
	testCases := []struct {
		name   string
		events []event.Event
		want   string
	}{
		{
			name:   "default case",
			events: event.ParseEvents(defaultCaseEventsList),
			// (7 - 3) / 7 = 4 / 7 ~ 57%
			want: "57.14",
		},
		{
			name: "100 percent",
			events: event.ParseEvents(
				"2024-10-07 13:46:47.679\t0\th\th\n" +
					"2024-10-07 13:46:56.521\t3\tenter\tenter",
			),
			want: "100.00",
		},
		{
			name:   "no events",
			events: []event.Event{},
			want:   "0.00",
		},
	}

	for _, tc := range testCases {
		got := accuracy(tc.events)
		fGot := fmt.Sprintf("%.2f", got)
		if fGot != tc.want {
			t.Errorf("%s want %s, got %s\n", tc.name, tc.want, fGot)
		}
	}

}

func TestUncorrectedErrors(t *testing.T) {
	type testCase struct {
		name   string
		events []event.Event
		want   int
	}

	testCases := []testCase{
		{
			name:   "no uncorrected errors",
			events: event.ParseEvents(defaultCaseEventsList),
			want:   0,
		},
		{
			name: "no uncorrected errors, i offset",
			// typed:    console↲
			// expected: console↲
			events: event.ParseEvents(
				"2024-10-07 16:29:26.916\t10\tc\tc\n" +
					"2024-10-07 16:29:27.004\t11\to\to\n" +
					"2024-10-07 16:29:27.095\t12\tn\tn\n" +
					"2024-10-07 16:29:27.279\t13\ts\ts\n" +
					"2024-10-07 16:29:27.416\t14\to\to\n" +
					"2024-10-07 16:29:27.667\t15\tl\tl\n" +
					"2024-10-07 16:29:27.784\t16\te\te\n" +
					"2024-10-07 16:29:31.538\t17\tenter\tenter",
			),
			want: 0,
		},
		{
			name: "all backspaces",
			events: event.ParseEvents(
				"2024-10-07 13:46:47.679\t4\tbackspace\n" +
					"2024-10-07 13:46:48.298\t3\tbackspace\n" +
					"2024-10-07 13:46:49.442\t2\tbackspace\n" +
					"2024-10-07 13:46:51.160\t1\tbackspace\n",
			),
			want: 0,
		},
		{
			name: "some uncorrected errors",
			// typed:    hei↲
			// expected: hey↲
			events: event.ParseEvents(
				"2024-10-07 13:46:49.442\t0\th\th\n" +
					"2024-10-07 13:46:51.160\t1\te\te\n" +
					"2024-10-07 13:46:52.781\t2\ti\ty\n" +
					"2024-10-07 13:46:56.521\t3\tenter\tenter",
			),
			want: 1,
		},
		{
			name: "all uncorrected errors",
			// expected: hey↲
			// typed:    omg!
			events: event.ParseEvents(
				"2024-10-07 13:46:49.442\t0\to\th\n" +
					"2024-10-07 13:46:51.160\t1\tm\te\n" +
					"2024-10-07 13:46:52.781\t2\tg\ty\n" +
					"2024-10-07 13:46:56.521\t3\t!\tenter",
			),
			want: 4,
		},
		{
			name:   "empty event list",
			events: []event.Event{},
			want:   0,
		},
		{
			name: "No uncorrected errors, but with newline and tab",
			// expected and typed:
			// Todo
			//    - Drink milk
			events: event.ParseEvents(
				"2024-10-25 13:30:22.846\t0\tT\tT\n" +
					"2024-10-25 13:30:22.958\t1\to\to\n" +
					"2024-10-25 13:30:23.029\t2\td\td\n" +
					"2024-10-25 13:30:23.121\t3\to\to\n" +
					"2024-10-25 13:30:23.540\t4\tenter\tenter\n" +
					"2024-10-25 13:30:24.394\t6\t-\t-\n" +
					"2024-10-25 13:30:24.591\t7\tspace\tspace\n" +
					"2024-10-25 13:30:24.826\t8\tD\tD\n" +
					"2024-10-25 13:30:25.030\t9\tr\tr\n" +
					"2024-10-25 13:30:25.076\t10\ti\ti\n" +
					"2024-10-25 13:30:25.171\t11\tn\tn\n" +
					"2024-10-25 13:30:25.255\t12\tk\tk\n" +
					"2024-10-25 13:30:25.408\t13\tspace\tspace\n" +
					"2024-10-25 13:30:25.611\t14\tm\tm\n" +
					"2024-10-25 13:30:25.665\t15\ti\ti\n" +
					"2024-10-25 13:30:25.737\t16\tl\tl\n" +
					"2024-10-25 13:30:25.942\t17\tk\tk\n" +
					"2024-10-25 13:30:26.199\t18\tenter\tenter",
			),
			want: 0,
		},
	}

	for _, tc := range testCases {
		got := numUncorrectedErrors(tc.events)
		if got != tc.want {
			t.Errorf("%s: got %d, want %d", tc.name, got, tc.want)
		}
	}
}

func TestDuration(t *testing.T) {
	testCases := []struct {
		name   string
		events []event.Event
		want   time.Duration
	}{
		{
			name: "default case",
			events: event.ParseEvents("2024-10-25 13:30:25.000\t17\tk\tk\n" +
				"2024-10-25 13:30:26.000\t18\tenter\tenter"),
			want: time.Second,
		},
		{
			name:   "no events",
			events: []event.Event{},
			want:   0.0,
		},
		{
			name:   "one event",
			events: event.ParseEvents("2024-10-25 13:30:26.000\t18\tenter\tenter"),
			want:   0.0,
		},
	}

	for _, tc := range testCases {
		got := duration(tc.events)
		if got != tc.want {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}

// Used to compare wpm values.
func aboutTheSame(a float64, b float64) bool {
	af := fmt.Sprintf("%.2f", a)
	bf := fmt.Sprintf("%.2f", b)
	return af == bf
}

func TestWpm(t *testing.T) {
	testCases := []struct {
		name   string
		events []event.Event
		want   float64
	}{
		{
			name: "no mistakes",
			// typed:    console↲
			// expected: console↲
			events: event.ParseEvents(
				"2024-10-07 16:29:26.916\t0\tc\tc\n" +
					"2024-10-07 16:29:27.004\t1\to\to\n" +
					"2024-10-07 16:29:27.095\t2\tn\tn\n" +
					"2024-10-07 16:29:27.279\t3\ts\ts\n" +
					"2024-10-07 16:29:27.416\t4\to\to\n" +
					"2024-10-07 16:29:27.667\t5\tl\tl\n" +
					"2024-10-07 16:29:27.784\t6\te\te\n" +
					"2024-10-07 16:29:31.538\t7\tenter\tenter",
			),
			// ((8/5) - 0) / ((31.538 - 26.916) / 60) = 20.77
			want: 20.77,
		},
		{
			name: "no mistakes, but i is offset",
			// typed:    console↲
			// expected: console↲
			events: event.ParseEvents(
				"2024-10-07 16:29:26.916\t10\tc\tc\n" +
					"2024-10-07 16:29:27.004\t11\to\to\n" +
					"2024-10-07 16:29:27.095\t12\tn\tn\n" +
					"2024-10-07 16:29:27.279\t13\ts\ts\n" +
					"2024-10-07 16:29:27.416\t14\to\to\n" +
					"2024-10-07 16:29:27.667\t15\tl\tl\n" +
					"2024-10-07 16:29:27.784\t16\te\te\n" +
					"2024-10-07 16:29:31.538\t17\tenter\tenter",
			),
			// ((8/5) - 0) / ((31.538 - 26.916) / 60) = 20.77
			want: 20.77,
		},
		{
			name: "with mistakes",
			// uncorrected: 1
			// typed:       consold↲
			// expected:    console↲
			events: event.ParseEvents(
				"2024-10-07 16:29:26.916\t0\tc\tc\n" +
					"2024-10-07 16:29:27.004\t1\to\to\n" +
					"2024-10-07 16:29:27.095\t2\tn\tn\n" +
					"2024-10-07 16:29:27.279\t3\ts\ts\n" +
					"2024-10-07 16:29:27.416\t4\to\to\n" +
					"2024-10-07 16:29:27.667\t5\tl\tl\n" +
					"2024-10-07 16:29:27.784\t6\td\te\n" + // miss
					"2024-10-07 16:29:31.538\t7\tenter\tenter",
			),
			// ((8/5) - 1) / ((31.538 - 26.916) / 60) = 7.79
			want: 7.79,
		},
		{
			name: "longer than one minute",
			// typed:       consold↲
			// expected:    console↲
			events: event.ParseEvents(
				"2024-10-07 16:29:26.916\t0\tc\tc\n" +
					"2024-10-07 16:29:27.004\t1\to\to\n" +
					"2024-10-07 16:29:27.095\t2\tn\tn\n" +
					"2024-10-07 16:29:27.279\t3\ts\ts\n" +
					"2024-10-07 16:29:27.416\t4\to\to\n" +
					"2024-10-07 16:29:27.667\t5\tl\tl\n" +
					"2024-10-07 16:29:27.784\t6\te\te\n" +
					"2024-10-07 16:31:26.916\t7\tenter\tenter",
			),
			// ((8/5) - 0) / 2 = 0.8
			want: 0.8,
		},
		{
			name:   "no events",
			events: []event.Event{},
			want:   0.0,
		},
		{
			name:   "one event",
			events: event.ParseEvents("2024-10-07 16:29:26.916\t0\tc\tc"),
			want:   0.0,
		},
		{
			name: "going backwards in index, only one typed character",
			// I type one character in this time period, so this should
			// return 0.0 becuase I don't have two characters to compare.
			events: event.ParseEvents(
				"2024-10-07 16:29:27.279\t3\tbackspace\n" +
					"2024-10-07 16:29:27.416\t2\tbackspace\n" +
					"2024-10-07 16:29:27.667\t1\tbackspace\n" +
					"2024-10-07 16:29:27.784\t0\tc\tc",
			),
			want: 0.0,
		},
		{
			name: "going backwards in index, typed multiple characters",
			events: event.ParseEvents(
				"2024-10-07 16:29:27.279\t3\tbackspace\n" +
					"2024-10-07 16:29:27.416\t2\tbackspace\n" +
					"2024-10-07 16:29:27.667\t1\tbackspace\n" +
					"2024-10-07 16:29:27.784\t0\tc\tc\n" +
					"2024-10-07 16:29:28.000\t1\to\to",
			),
			// Note the duration is between the 4th and 5th events.
			// (2/5 - 0) / ((28.00 - 27.784) / 60) = 111.11
			want: 111.11,
		},
	}

	for _, tc := range testCases {
		got := wpm(tc.events)
		if !aboutTheSame(got, tc.want) {
			t.Errorf("%s: got about %.2f, wanted about %.2f\n", tc.name, got, tc.want)
		}
	}
}

func TestWpmRaw(t *testing.T) {
	events := event.ParseEvents(
		"2024-10-07 16:29:26.916\t0\tc\tc\n" +
			"2024-10-07 16:29:27.004\t1\to\to\n" +
			"2024-10-07 16:29:27.095\t2\tn\tn\n" +
			"2024-10-07 16:29:27.279\t3\ts\ts\n" +
			"2024-10-07 16:29:27.416\t4\to\to\n" +
			"2024-10-07 16:29:27.667\t5\tl\tl\n" +
			"2024-10-07 16:29:27.784\t6\td\te\n" + // miss
			"2024-10-07 16:29:31.538\t7\tenter\tenter",
	)
	// ((8/5) - 0) / ((31.538 - 26.916) / 60)
	// 1.6 / ((31.538 - 26.916) / 60) = 20.770229338
	want := 20.77

	got := wpmRaw(events)

	if !aboutTheSame(got, want) {
		t.Errorf("got about %.2f, wanted about %.2f\n", got, want)
	}

}

func TestMostMissedKeys(t *testing.T) {

	type testCase struct {
		name   string
		events []event.Event
		want   string
	}

	testCases := []testCase{
		{
			name:   "default case",
			events: event.ParseEvents(defaultCaseEventsList),
			want:   "y (2 times), h (1 time)",
		},
		{
			// " = 32, ( = 40, . = 46
			name: "sorted special characters",
			// mistakes: 9
			// typed:    console.log("Ey")↲
			// expected: console.log("Ey")↲
			events: event.ParseEvents("2024-10-07 16:29:26.916\t0\tc\tc\n" +
				"2024-10-07 16:29:27.004\t1\to\to\n" +
				"2024-10-07 16:29:27.095\t2\tn\tn\n" +
				"2024-10-07 16:29:27.279\t3\ts\ts\n" +
				"2024-10-07 16:29:27.416\t4\to\to\n" +
				"2024-10-07 16:29:27.667\t5\tl\tl\n" +
				"2024-10-07 16:29:27.784\t6\te\te\n" +
				"2024-10-07 16:29:31.538\t7\td\t.\n" + // miss
				"2024-10-07 16:29:32.243\t8\tbackspace\n" +
				"2024-10-07 16:29:33.216\t7\te\t.\n" + // miss
				"2024-10-07 16:29:33.432\t8\tbackspace\n" +
				"2024-10-07 16:29:33.811\t7\tr\t.\n" + // miss
				"2024-10-07 16:29:34.175\t8\tbackspace\n" +
				"2024-10-07 16:29:34.768\t7\t.\t.\n" +
				"2024-10-07 16:29:35.313\t8\tl\tl\n" +
				"2024-10-07 16:29:35.502\t9\to\to\n" +
				"2024-10-07 16:29:35.676\t10\tg\tg\n" +
				"2024-10-07 16:29:37.565\t11\t8\t(\n" + // miss
				"2024-10-07 16:29:38.374\t12\tbackspace\n" +
				"2024-10-07 16:29:38.750\t11\t9\t(\n" + // miss
				"2024-10-07 16:29:39.810\t12\tbackspace\n" +
				"2024-10-07 16:29:41.048\t11\t0\t(\n" + // miss
				"2024-10-07 16:29:41.380\t12\tbackspace\n" +
				"2024-10-07 16:29:42.058\t11\t(\t(\n" +
				"2024-10-07 16:29:45.428\t12\t2\t\"\n" + // miss
				"2024-10-07 16:29:45.991\t13\tbackspace\n" +
				"2024-10-07 16:29:46.178\t12\t3\t\"\n" + // miss
				"2024-10-07 16:29:46.502\t13\tbackspace\n" +
				"2024-10-07 16:29:48.972\t12\t\"\t\"\n" +
				"2024-10-07 16:29:49.427\t13\tE\tE\n" +
				"2024-10-07 16:29:50.641\t14\ty\ty\n" +
				"2024-10-07 16:29:55.056\t15\t4\t\"\n" + // miss
				"2024-10-07 16:29:55.797\t16\tbackspace\n" +
				"2024-10-07 16:29:56.540\t15\t\"\t\"\n" +
				"2024-10-07 16:29:57.101\t16\t)\t)\n" +
				"2024-10-07 16:29:58.765\t17\tenter\tenter"),
			want: "\" (3 times), ( (3 times), . (3 times)",
		},
		{
			name: "whitespace mistakes",
			// expected:
			// hey there↲
			// two lines↲
			events: event.ParseEvents("2024-10-07 16:09:16.628\t0\th\th\n" +
				"2024-10-07 16:09:17.177\t1\te\te\n" +
				"2024-10-07 16:09:17.274\t2\ty\ty\n" +
				"2024-10-07 16:09:18.290\t3\td\tspace\n" + // miss
				"2024-10-07 16:09:19.222\t4\tbackspace\n" +
				"2024-10-07 16:09:20.319\t3\te\tspace\n" + // miss
				"2024-10-07 16:09:20.837\t4\tbackspace\n" +
				"2024-10-07 16:09:21.151\t3\tspace\tspace\n" +
				"2024-10-07 16:09:21.487\t4\tt\tt\n" +
				"2024-10-07 16:09:21.562\t5\th\th\n" +
				"2024-10-07 16:09:21.691\t6\te\te\n" +
				"2024-10-07 16:09:21.832\t7\tr\tr\n" +
				"2024-10-07 16:09:21.913\t8\te\te\n" +
				"2024-10-07 16:09:22.902\t9\ti\tenter\n" + // miss
				"2024-10-07 16:09:23.937\t10\tbackspace\n" +
				"2024-10-07 16:09:24.429\t9\ts\tenter\n" + // miss
				"2024-10-07 16:09:25.283\t10\tbackspace\n" +
				"2024-10-07 16:09:26.132\t9\tenter\tenter\n" +
				"2024-10-07 16:09:26.745\t10\tt\tt\n" +
				"2024-10-07 16:09:28.081\t11\tw\tw\n" +
				"2024-10-07 16:09:28.148\t12\to\to\n" +
				"2024-10-07 16:09:28.843\t13\tspace\tspace\n" +
				"2024-10-07 16:09:29.498\t14\tl\tl\n" +
				"2024-10-07 16:09:29.537\t15\ti\ti\n" +
				"2024-10-07 16:09:29.611\t16\tn\tn\n" +
				"2024-10-07 16:09:29.789\t17\te\te\n" +
				"2024-10-07 16:09:29.831\t18\ts\ts\n" +
				"2024-10-07 16:09:30.694\t19\tenter\tenter"),
			want: "enter (2 times), space (2 times)",
		},
		{
			name: "only show limited number of misses",
			// mistakes: 4
			// typed:    asdf↲
			// expected: asdf↲
			events: event.ParseEvents("2024-10-07 16:46:36.929\t0\tq\ta\n" + // miss
				"2024-10-07 16:46:37.331\t1\tbackspace\n" +
				"2024-10-07 16:46:38.067\t0\ta\ta\n" +
				"2024-10-07 16:46:39.145\t1\tw\ts\n" + // miss
				"2024-10-07 16:46:39.408\t2\tbackspace\n" +
				"2024-10-07 16:46:39.831\t1\ts\ts\n" +
				"2024-10-07 16:46:40.827\t2\te\td\n" + // miss
				"2024-10-07 16:46:41.089\t3\tbackspace\n" +
				"2024-10-07 16:46:41.658\t2\td\td\n" +
				"2024-10-07 16:46:43.471\t3\tr\tf\n" + // miss
				"2024-10-07 16:46:43.942\t4\tbackspace\n" +
				"2024-10-07 16:46:44.862\t3\tf\tf\n" +
				"2024-10-07 16:46:46.290\t4\tenter\tenter"),
			want: "a (1 time), d (1 time), f (1 time)",
		},
		{
			name:   "no events",
			events: []event.Event{},
			want:   "",
		},
		{
			name: "no mistakes",
			events: event.ParseEvents(
				"2024-10-07 16:29:26.916\t0\tc\tc\n" +
					"2024-10-07 16:29:27.004\t1\to\to\n" +
					"2024-10-07 16:29:27.095\t2\tn\tn\n" +
					"2024-10-07 16:29:27.279\t3\ts\ts\n" +
					"2024-10-07 16:29:27.416\t4\to\to\n" +
					"2024-10-07 16:29:27.667\t5\tl\tl\n" +
					"2024-10-07 16:29:27.784\t6\te\te\n" +
					"2024-10-07 16:29:31.538\t7\tenter\tenter",
			),
			want: "",
		},
	}

	for _, tc := range testCases {
		got := mostMissedKeys(tc.events)
		if got != tc.want {
			t.Errorf("%s:\n\tgot  %s\n\twant %s", tc.name, got, tc.want)
		}
	}
}
