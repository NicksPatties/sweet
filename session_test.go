package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestUpdate(t *testing.T) {
	title := "Handle Enter key message"
	exercise := `a test
with a newline`
	initialTyped := "a test"
	msg := tea.KeyMsg{
		Type:  tea.KeyEnter,
		Runes: []rune{},
		Alt:   false,
	}
	expectedTyped := `a test
`
	m := Model{
		title:         title,
		exercise:      exercise,
		typedExercise: initialTyped,
	}

	nuM, cmd := m.Update(msg)
	m = nuM.(Model)
	if cmd != nil {
		t.Fatalf(`Update() %s, returning end command, but don't want return a command`, title)
	}
	if m.typedExercise != expectedTyped {
		t.Fatalf(`Update() %s , %s want to match %s`, title, m.typedExercise, expectedTyped)
	}
}
