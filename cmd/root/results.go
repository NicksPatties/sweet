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

func wpm(events []event) float64 {
	if len(events) <= 1 {
		return 0.0
	}
	start := events[0].ts
	iOffset := events[0].i
	end := events[len(events)-1].ts
	duration := end.Sub(start)
	mins := float64(duration) / float64(time.Minute)
	wordSize := 5.0
	chars := events[len(events)-1].i - iOffset + 1
	words := float64(chars) / wordSize
	incorrect := float64(numIncorrect(events))
	return (words - incorrect) / mins
}

// Same as wpm, but doesn't subtract incorrect chars.
func wpmRaw(events []event) float64 {
	start := events[0].ts
	end := events[len(events)-1].ts
	duration := end.Sub(start)
	mins := float64(duration) / float64(time.Minute)
	wordSize := 5.0
	chars := events[len(events)-1].i + 1
	words := float64(chars) / wordSize
	return words / mins
}

func wpmGraph() string {
	events := parseEvents(`2024-10-07 13:46:47.679: 0 a h
2024-10-07 13:46:48.298: 1 backspace
2024-10-07 13:46:49.442: 0 h h
2024-10-07 13:46:51.160: 1 e e
2024-10-07 13:46:52.781: 2 i y
2024-10-07 13:46:53.316: 3 backspace
2024-10-07 13:46:54.688: 2 k y
2024-10-07 13:46:55.262: 3 backspace
2024-10-07 13:46:55.997: 2 y y
2024-10-07 13:46:56.521: 3 enter enter`)
	d := events[len(events)-1].ts.Sub(events[0].ts)
	seconds := int(d.Seconds()) + 1
	wpmData := make([]float64, seconds)
	eventBuckets := make([][]event, seconds)

	for _, event := range events {
		tsDiff := event.ts.Sub(events[0].ts)
		bucketId := int(tsDiff.Seconds())
		eventBuckets[bucketId] = append(eventBuckets[bucketId], event)
	}

	for i, eventBucket := range eventBuckets {
		wpmData[i] = wpm(eventBucket)
	}

	// get plot data for events
	// create an array of floats of size of seconds length
	// for each second
	//   get the events that are within this second
	//   calculate the wpm and wpmRaw for those events
	//   add the wpm and wpmRaw for that second into the plot data
	// create the plot and end the data

	fmt.Printf("wpmData: %v\n", wpmData)

	return g.Plot(wpmData)
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

func duration(startTime time.Time, endTime time.Time) string {
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
	fmt.Printf("Results of %s:\n", m.exercise.name)
	fmt.Printf("WPM: %.f\n", wpm(m.events))
	fmt.Printf("Mistakes: %d\n", numMistakes(m.events))
	fmt.Printf("Accuracy: %s%%\n", accuracy(m.events))
	fmt.Printf("Duration: %s\n", duration(m.startTime, m.endTime))
	fmt.Printf("Most missed keys: %s\n", mostMissedKeys(m.events))
	fmt.Printf("Graph:\n%s", wpmGraph())
}
