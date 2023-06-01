package main

import "fmt"

func (m Model) WPM() int {
	return 50
}

func (m Model) CPM() int {
	return 100
}

func (m Model) Accuracy() float32 {
	return .90
}

func (m Model) NumMistakes() int {
	return 1
}

func ShowResults(m Model) {
	fmt.Printf("Results of %s:\n", m.title)
	fmt.Printf("WPM: %d\nCPM: %d\nAccuracy: %.2f\nMistakes: %d\n", m.WPM(), m.CPM(), m.Accuracy(), m.NumMistakes())
}
