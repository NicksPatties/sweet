package main

import (
	"fmt"
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

func (m Model) WPM() int {
	return 50
}

func (m Model) CPM() int {
	return 100
}

// Gives a percentage accuracy of the typed exercise
func Accuracy(typed string, exercise string) float32 {
	if len(typed) == 0 || len(exercise) == 0 {
		return float32(0)
	}
	accuracy := float32(0)
	minLengthString := exercise
	if len(typed) < len(exercise) {
		minLengthString = typed
	}

	m := float32(NumMistakes(typed, exercise))
	l := float32(len(requiredRunes(minLengthString)))
	accuracy = (l - m) / l

	return accuracy
}

// Counts the number of mistakes made in an exercise. Only counts up to the number
// of characters typed into the exercise. If the number of characters typed exceeds
// the length of the exercise, then this function only counts up to the length of
// the exercise, and the remaining characters are discarded
func NumMistakes(typed string, exercise string) int {
	mistakes := 0
	r := min(len(typed), len(exercise))

	for i := 0; i < r; i++ {
		if typed[i] != exercise[i] {
			mistakes++
		}
	}

	return mistakes
}

func ShowResults(m Model) {
	fmt.Printf("Results of %s:\n", m.title)
	fmt.Printf("Mistakes: %d\n", NumMistakes(m.typedExercise, m.exercise))
	fmt.Printf("Accuracy: %.2f\n", Accuracy(m.typedExercise, m.exercise))
}
