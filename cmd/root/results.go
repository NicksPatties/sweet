package root

import (
	"fmt"
	g "github.com/guptarohit/asciigraph"
	"sort"
	"strings"
	"time"
)

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func requiredRunes(s string) []rune {
	arr := []rune{}
	foundNewline := false
	for _, rn := range s {
		if foundNewline && isWhitespace(rn) {
			continue
		}
		foundNewline = rn == Enter
		arr = append(arr, rn)
	}
	return arr
}

// Remove backspaces from a list of events.
//
// Removing backspace events simplifies wpm calculations.
func removeBackspaces(events []event) []event {
	enb := []event{}
	for _, e := range events {
		if e.typed != "backspace" {
			enb = append(enb, e)
		}
	}
	return enb
}

// Calculates the words per minute (wpm) using the events in the list.
// Also allows the duration to be overridden, which is useful
// for calculating wpm per second, which is used in the `wpmGraph` function.
//
// You should avoid using this function in favor of specific wpm
// functions, including `wpm`, `wpmRaw`, `wpmRawPerSecond`, and so on.
func wpmBase(e []event, raw bool, d time.Duration) float64 {
	events := removeBackspaces(e)
	// cannot calculate wpm with less than 2 events
	if len(events) < 2 {
		return 0.0
	}
	start := events[0].ts
	end := events[len(events)-1].ts
	if d == 0 {
		d = end.Sub(start)
	}
	mins := d.Minutes()
	wordSize := 5.0
	chars := len(events)
	words := float64(chars) / wordSize
	var result float64
	if raw {
		result = (words) / mins
	} else {
		incorrect := float64(numUncorrectedErrors(events))
		// avoid negative wpm
		if words-incorrect < 0 {
			return 0.0
		}
		result = (words - incorrect) / mins
	}
	return result
}

// Calculates the words per minute (wpm) based on the
// events that occurred during the exercise.
func wpm(events []event) float64 {
	return wpmBase(events, false, 0)
}

// Calculates the wpm of a series of events that
// lasted for `n` seconds. This is used to calculate the
// rolling average wpm during the course of the exercise.
func wpmForNSeconds(events []event, n int) float64 {
	seconds := time.Duration(n) * time.Second
	return wpmBase(events, false, seconds)
}

// Calculates the raw wpm for a series of events.
// Assumes the events occurred in the same second.
func wpmRawPerSecond(events []event) float64 {
	return wpmBase(events, true, time.Second)
}

func wpmGraph(events []event) string {
	d := events[len(events)-1].ts.Sub(events[0].ts)
	seconds := int(d.Seconds()) + 1
	wpmData := make([]float64, seconds)
	wpmRawData := make([]float64, seconds)
	eventBuckets := make([][]event, seconds)

	for _, event := range events {
		tsDiff := event.ts.Sub(events[0].ts)
		bucketId := int(tsDiff.Seconds())
		eventBuckets[bucketId] = append(eventBuckets[bucketId], event)
	}

	var currEvents []event
	for i, eventBucket := range eventBuckets {
		currEvents = append(currEvents, eventBucket...)
		currSeconds := i + 1
		wpmData[i] = wpmForNSeconds(currEvents, currSeconds)
		wpmRawData[i] = wpmRawPerSecond(eventBucket)
	}

	return g.PlotMany(
		[][]float64{wpmRawData, wpmData},
		g.SeriesColors(g.Gray, g.Default),
		g.SeriesLegends("raw wpm", "wpm"),
		g.Height(10),
		g.Width(0), // auto scaling
		g.LowerBound(0),
		g.Precision(0),
	)
}

// Gives a percentage accuracy of the typed exercise.
// Accuracy is the percentage of mistakes over the number
// of the total typed characters, excluding backspaces.
//
// Note, even if all characters at the end of an exercise
// are correct, you can have an accuracy of less than 100%
// if you made any mistakes.
func accuracy(events []event) float64 {
	if len(events) == 0 {
		return 0.0
	}
	mistakes := float64(0)
	total := float64(0)
	for _, e := range events {
		if e.typed == "backspace" {
			continue
		}
		if e.typed != e.expected {
			mistakes++
		}
		total++
	}

	return (total - mistakes) / total * 100.0
}

// Returns the number of uncorrected errors in
// a series of events.
//
// If a series of events only contains
// backspaces, then it's assumed no uncorrected
// errors have been made, because the user is in
// the process of correcting the error.
func numUncorrectedErrors(events []event) int {
	if len(events) == 0 {
		return 0
	}
	correct := map[int]bool{}
	for _, e := range events {
		if e.typed == "backspace" {
			correct[e.i] = true
		} else {
			correct[e.i] = e.typed == e.expected
		}
	}
	count := 0
	for _, c := range correct {
		if !c {
			count++
		}
	}
	return count
}

// Returns the number of mistakes made during
// an exercise. This includes both corrected and
// uncorrected errors.
//
// Backspaces do not count as mistakes.
func numMistakes(events []event) (mistakes int) {
	for _, e := range events {
		if e.typed == "backspace" {
			continue
		}
		if e.typed != e.expected {
			mistakes++
		}
	}
	return
}

// Returns the duration between the first event and
// the last event of the array. If there are less than
// two events in the list, it returns zero duration.
func duration(events []event) time.Duration {
	if len(events) < 2 {
		return 0.0
	}
	return events[len(events)-1].ts.Sub(events[0].ts)
}

// Finds the most missed key presses when completing
// an exercise. Missed keys are sorted alphabetically,
// and by the number of misses. Also, sets a limit
// of number of keys missed to avoid overflowing the line.
func mostMissedKeys(events []event) string {
	misses := map[string]int{}
	for _, e := range events {
		if e.typed != "backspace" && e.typed != e.expected {
			misses[e.expected]++
		}
	}

	keys := []string{}
	for key := range misses {
		keys = append(keys, key)
	}
	// NOTE: Do I want to sort the keys with the same
	// character by time?
	sort.Strings(keys)
	sort.SliceStable(keys, func(i int, j int) bool {
		return misses[keys[i]] > misses[keys[j]]
	})

	// A miss looks like this: "a (2 times)"
	var missesStrs []string
	limit := 3
	for i := 0; i < len(keys) && i < limit; i++ {
		key := keys[i]
		times := misses[key]
		var t string
		if times == 1 {
			t = "time"
		} else {
			t = "times"
		}
		missesStrs = append(missesStrs, fmt.Sprintf("%s (%d %s)", key, misses[key], t))
	}
	return strings.Join(missesStrs, ", ")
}

// Prints the results of a repetition.
func showResults(rep Rep) {
	fmt.Printf("results of %s:\n", rep.name)
	fmt.Printf("wpm:                 %.f\n", rep.wpm)
	fmt.Printf("uncorrected errors:  %d\n", rep.errs)
	fmt.Printf("duration:            %s\n", rep.dur)
	fmt.Printf("mistakes:            %d\n", rep.miss)
	fmt.Printf("accuracy:            %.2f%%\n", rep.acc)
	if rep.miss > 0 {
		fmt.Printf("most missed keys:    %s\n", mostMissedKeys(rep.events))
	}
	fmt.Printf("graph:\n%s", wpmGraph(rep.events))
	fmt.Println()
}
