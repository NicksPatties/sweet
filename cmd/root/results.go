package root

import (
	"fmt"
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

// Calculates the words per minute based on the calculations in this link:
// https://www.speedtypingonline.com/typing-equations
func wpm(start time.Time, end time.Time, typed string, exercise string, wordSize int) float64 {
	if start.After(end) {
		end = time.Now()
	}
	minLengthString := exercise
	if len(typed) < len(exercise) {
		minLengthString = typed
	}
	mins := end.Sub(start).Minutes()
	incorrect := float64(numIncorrectCharacters(typed, exercise))
	typedEntries := len(requiredRunes(minLengthString))
	words := float64(typedEntries / wordSize)
	return (words - incorrect) / mins
}

func cpm(start time.Time, end time.Time, typed string, exercise string) float64 {
	return wpm(start, end, typed, exercise, 1)
}

// Gives a percentage accuracy of the typed exercise
func accuracy(typed string, exercise string) float32 {
	if len(typed) == 0 || len(exercise) == 0 {
		return float32(0)
	}
	var accuracy float32
	minLengthString := exercise
	if len(typed) < len(exercise) {
		minLengthString = typed
	}

	m := float32(numIncorrectCharacters(typed, exercise))
	l := float32(len(requiredRunes(minLengthString)))
	accuracy = (l - m) / l

	return accuracy * 100
}

func numIncorrectCharacters(typed string, exercise string) (incorrect int) {
	r := min(len(typed), len(exercise))
	for i := 0; i < r; i++ {
		if typed[i] != exercise[i] {
			incorrect++
		}
	}
	return
}

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
	buckets := map[string]int{}
	for _, e := range events {
		if e.typed != "backspace" && e.typed != e.expected {
			buckets[e.expected]++
		}
	}

	keys := make([]string, 0, len(buckets))
	for key := range buckets {
		keys = append(keys, key)
	}
	// NOTE: Do I want to sort the keys with the same
	// character by time?
	sort.Strings(keys)
	sort.SliceStable(keys, func(i int, j int) bool {
		return buckets[keys[i]] > buckets[keys[j]]
	})

	// A miss looks like this: "a (2 times)"
	var misses []string
	limit := 3
	for i := 0; i < len(keys) && i < limit; i++ {
		key := keys[i]
		times := buckets[key]
		var t string
		if times == 1 {
			t = "time"
		} else {
			t = "times"
		}
		misses = append(misses, fmt.Sprintf("%s (%d %s)", key, buckets[key], t))
	}
	return strings.Join(misses, ", ")
}

func showResults(m exerciseModel) {
	fmt.Printf("Results of %s:\n", m.exercise.name)
	fmt.Printf("WPM: %.f\n", wpm(m.startTime, m.endTime, m.typedText, m.exercise.text, WORD_SIZE))
	fmt.Printf("Mistakes: %d\n", numMistakes(m.events))
	fmt.Printf("Accuracy: %.2f%%\n", accuracy(m.typedText, m.exercise.text))
	fmt.Printf("Duration: %s\n", duration(m.startTime, m.endTime))
	fmt.Printf("Most missed keys: %s", mostMissedKeys(m.events))
}
