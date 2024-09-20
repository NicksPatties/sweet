package exercise

import (
	"fmt"
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
	mistakes := float64(numMistakes(typed, exercise))
	typedEntries := len(requiredRunes(minLengthString))
	words := float64(typedEntries / wordSize)
	return (words - mistakes) / mins
}

func cpm(start time.Time, end time.Time, typed string, exercise string) float64 {
	return wpm(start, end, typed, exercise, 1)
}

// Gives a percentage accuracy of the typed exercise
func accuracy(typed string, exercise string) float32 {
	if len(typed) == 0 || len(exercise) == 0 {
		return float32(0)
	}
	accuracy := float32(0)
	minLengthString := exercise
	if len(typed) < len(exercise) {
		minLengthString = typed
	}

	m := float32(numMistakes(typed, exercise))
	l := float32(len(requiredRunes(minLengthString)))
	accuracy = (l - m) / l

	return accuracy * 100
}

// Counts the number of mistakes made in an exercise. Only counts up to the number
// of characters typed into the exercise. If the number of characters typed exceeds
// the length of the exercise, then this function only counts up to the length of
// the exercise, and the remaining characters are discarded
func numMistakes(typed string, exercise string) int {
	mistakes := 0
	r := min(len(typed), len(exercise))

	for i := 0; i < r; i++ {
		if typed[i] != exercise[i] {
			mistakes++
		}
	}

	return mistakes
}

func duration(startTime time.Time, endTime time.Time) string {
	nanos := (endTime.UnixMilli() - startTime.UnixMilli()) * int64(time.Millisecond)
	d := time.Duration(nanos)

	s := fmt.Sprintf("%.3fs", d.Seconds())
	return s
}

func showResults(m exerciseModel) {
	fmt.Printf("Results of %s:\n", m.title)
	fmt.Printf("WPM: %.f\n", wpm(m.startTime, m.endTime, m.typedExercise, m.exercise, WORD_SIZE))
	fmt.Printf("Mistakes: %d\n", numMistakes(m.typedExercise, m.exercise))
	fmt.Printf("Accuracy: %.2f%%\n", accuracy(m.typedExercise, m.exercise))
	fmt.Printf("Duration: %s\n", duration(m.startTime, m.endTime))
}
