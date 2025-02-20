package root

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/NicksPatties/sweet/constants"
	db "github.com/NicksPatties/sweet/db"
	"github.com/NicksPatties/sweet/util"
	g "github.com/guptarohit/asciigraph"
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
		foundNewline = rn == constants.Enter
		arr = append(arr, rn)
	}
	return arr
}

// Remove backspaces from a list of events.
//
// Removing backspace events simplifies wpm calculations.
func removeBackspaces(events []util.Event) []util.Event {
	enb := []util.Event{}
	for _, e := range events {
		if e.Typed != "backspace" {
			enb = append(enb, e)
		}
	}
	return enb
}

// Calculates the words per minute (wpm) using the events in the list.
// Also allows the duration to be overridden, which is useful
// for calculating wpm per second, which is used in the `wpmGraph` function.
//
// If `d` is equal to 0, then divide by the duration
// between the first and last events.
//
// You should avoid using this function in favor of specific wpm
// functions, including `wpm`, `wpmRaw`, `wpmRawPerSecond`, and so on.
func wpmBase(e []util.Event, raw bool, d time.Duration) float64 {
	events := removeBackspaces(e)
	// cannot calculate wpm with less than 2 events
	if len(events) < 2 {
		return 0.0
	}
	start := events[0].Ts
	end := events[len(events)-1].Ts
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
func wpm(events []util.Event) float64 {
	return wpmBase(events, false, 0)
}

// Calculates the raw words per minute.
// Raw wpm does not subtract mistakes from the final
// wpm calculation
func wpmRaw(events []util.Event) float64 {
	return wpmBase(events, true, 0)
}

// Calculates the wpm of a series of events that
// lasted for `n` seconds. This is used to calculate the
// rolling average wpm during the course of the exercise.
func wpmForNSeconds(events []util.Event, n int) float64 {
	seconds := time.Duration(n) * time.Second
	return wpmBase(events, false, seconds)
}

// Calculates the raw wpm for a series of events.
// Assumes the events occurred in the same second.
func wpmRawPerSecond(events []util.Event) float64 {
	return wpmBase(events, true, time.Second)
}

func wpmGraph(events []util.Event) string {
	d := events[len(events)-1].Ts.Sub(events[0].Ts)
	seconds := int(d.Seconds()) + 1
	wpmData := make([]float64, seconds)
	wpmRawData := make([]float64, seconds)
	eventBuckets := make([][]util.Event, seconds)

	for _, event := range events {
		tsDiff := event.Ts.Sub(events[0].Ts)
		bucketId := int(tsDiff.Seconds())
		eventBuckets[bucketId] = append(eventBuckets[bucketId], event)
	}

	var currEvents []util.Event
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
func accuracy(events []util.Event) float64 {
	if len(events) == 0 {
		return 0.0
	}
	mistakes := float64(0)
	total := float64(0)
	for _, e := range events {
		if e.Typed == "backspace" {
			continue
		}
		if e.Typed != e.Expected {
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
func numUncorrectedErrors(events []util.Event) int {
	if len(events) == 0 {
		return 0
	}
	correct := map[int]bool{}
	for _, e := range events {
		if e.Typed == "backspace" {
			correct[e.I] = true
		} else {
			correct[e.I] = e.Typed == e.Expected
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
func numMistakes(events []util.Event) (mistakes int) {
	for _, e := range events {
		if e.Typed == "backspace" {
			continue
		}
		if e.Typed != e.Expected {
			mistakes++
		}
	}
	return
}

// Returns the duration between the first event and
// the last event of the array. If there are less than
// two events in the list, it returns zero duration.
func duration(events []util.Event) time.Duration {
	if len(events) < 2 {
		return 0.0
	}
	return events[len(events)-1].Ts.Sub(events[0].Ts)
}

// Finds the most missed key presses when completing
// an exercise. Missed keys are sorted alphabetically,
// and by the number of misses. Also, sets a limit
// of number of keys missed to avoid overflowing the line.
func mostMissedKeys(events []util.Event) string {
	misses := map[string]int{}
	for _, e := range events {
		if e.Typed != "backspace" && e.Typed != e.Expected {
			misses[e.Expected]++
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
func printExerciseResults(rep db.Rep) {
	fmt.Printf("results of %s:\n", rep.Name)
	fmt.Printf("wpm:                 %.f\n", rep.Wpm)
	fmt.Printf("uncorrected errors:  %d\n", rep.Errs)
	fmt.Printf("duration:            %s\n", rep.Dur)
	fmt.Printf("mistakes:            %d\n", rep.Miss)
	fmt.Printf("accuracy:            %.2f%%\n", rep.Acc)
	if rep.Miss > 0 {
		fmt.Printf("most missed keys:    %s\n", mostMissedKeys(rep.Events))
	}
	fmt.Printf("graph:\n%s", wpmGraph(rep.Events))
	fmt.Println()
}
