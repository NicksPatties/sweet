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

// Calculates the words per minute (wpm) using the events in the list.
// Also allows the duration to be overridden, which is useful
// for calculating wpm per second, which is used in the `wpmGraph` function.
func wpmBase(events []event, raw bool, d time.Duration) float64 {
	lenEventsNoBkspc := 0
	for _, e := range events {
		if e.typed != "backspace" {
			lenEventsNoBkspc++
		}
	}
	if lenEventsNoBkspc <= 1 {
		return 0.0
	}
	start := events[0].ts
	iOffset := events[0].i
	end := events[len(events)-1].ts
	if d == 0 {
		d = end.Sub(start)
	}
	mins := float64(d) / float64(time.Minute)
	wordSize := 5.0
	// TODO: This line smells really bad.
	// What do I need to do?
	// 1. Take a list of events, and return the number
	// of characters typed without backspaces and test it.

	chars := events[len(events)-1].i - iOffset + 1
	words := float64(chars) / wordSize
	var result float64
	if raw {
		result = (words) / mins
	} else {
		incorrect := float64(numIncorrect(events))
		if words-incorrect < 0 {
			return 0.0
		}
		result = (words - incorrect) / mins
	}
	return result
}

// Calculates the words per minute (wpm) based on the
// events that are passed into the
func wpm(events []event) float64 {
	return wpmBase(events, false, 0)
}

// Words per minute per second. Used to calculate the wpm of an
// array of events that occur within the same second of each other.
func wpmPs(events []event) float64 {
	return wpmBase(events, false, time.Second)
}

// Same as wpm, but doesn't subtract incorrect chars.
func wpmRaw(events []event) float64 {
	return wpmBase(events, true, 0)
}

func wpmRawPs(events []event) float64 {
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

	currSeconds := time.Second
	var currEvents []event
	for i, eventBucket := range eventBuckets {
		// Need to calculate wpm from 0 to i seconds
		// duration needs to be i + 1 seconds long
		currEvents = append(currEvents, eventBucket...)
		wpmData[i] = wpmBase(currEvents, false, currSeconds)
		currSeconds += time.Second
		wpmRawData[i] = wpmRawPs(eventBucket)
	}

	return g.PlotMany(
		[][]float64{wpmRawData, wpmData},
		g.SeriesColors(g.Gray, g.Default),
		g.SeriesLegends("wpm raw", "wpm"),
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
func accuracy(events []event) string {
	if len(events) == 0 {
		return "0.00"
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

	acc := (total - mistakes) / total * 100.0
	return fmt.Sprintf("%.2f", acc)
}

func numIncorrect(events []event) int {
	if len(events) == 0 {
		return 0
	}
	maxI, minI := events[0].i, events[0].i
	for _, e := range events {
		if e.i > maxI {
			maxI = e.i
		}
		if e.i < minI {
			minI = e.i
		}
	}
	size := maxI - minI + 1
	correct := make([]bool, size)
	for _, e := range events {
		if e.typed == "backspace" {
			correct[e.i-minI] = true
		} else {
			correct[e.i-minI] = e.typed == e.expected
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
// an exercise.
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

func duration(events []event) time.Duration {
	return events[len(events)-1].ts.Sub(events[0].ts)
}

func durationOld(startTime time.Time, endTime time.Time) string {
	nanos := (endTime.UnixMilli() - startTime.UnixMilli()) * int64(time.Millisecond)
	d := time.Duration(nanos)

	s := fmt.Sprintf("%.3fs", d.Seconds())
	return s
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

func showResults(m exerciseModel) {
	incorrect := numIncorrect(m.events)
	d := duration(m.events)
	chars := len(m.exercise.text)
	wordLen := 5
	fmt.Printf("Events: %s\n", eventsString(m.events))
	fmt.Printf("Results of %s:\n", m.exercise.name)
	fmt.Printf("Incorrect chars:  %d\n", numIncorrect(m.events))
	fmt.Printf("Duration:         %s\n", duration(m.events))
	fmt.Printf("WPM = ((%d/%d) - %d) / %f\n", chars, wordLen, incorrect, d.Minutes())
	fmt.Printf("    = %.f\n", wpm(m.events))
	fmt.Println()
	fmt.Printf("Mistakes made:    %d\n", numMistakes(m.events))
	fmt.Printf("Accuracy:         %s%%\n", accuracy(m.events))
	fmt.Printf("Most missed keys: %s\n", mostMissedKeys(m.events))
	fmt.Printf("Graph:\n%s", wpmGraph(m.events))
	fmt.Println()
}
