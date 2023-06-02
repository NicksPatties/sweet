package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	title         string
	exercise      string
	typedExercise string
	quitEarly     bool
	startTime     time.Time
	endTime       time.Time
}

func initialModel(t string, ex string) Model {
	return Model{
		title:         t,
		exercise:      ex,
		typedExercise: "",
		quitEarly:     false,
		startTime:     time.Time{},
		endTime:       time.Time{},
	}
}

func (m Model) finished() bool {
	l := len(m.exercise)
	if len(m.typedExercise) < l {
		return false
	}
	// they're the same length, so check the last characters
	exLast := rune(m.exercise[l-1])
	typedLast := rune(m.typedExercise[l-1])

	if exLast != typedLast {
		return false
	}
	return true
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitEarly = true
			return m, tea.Quit
		case tea.KeyBackspace:
			m = m.DeleteCharacter()
		case tea.KeyRunes, tea.KeySpace, tea.KeyEnter:
			if m.startTime.IsZero() {
				m.startTime = time.Now()
			}
			if msg.Type == tea.KeyEnter {
				m = m.AddRuneToExercise(Enter)
			} else {
				m = m.AddRuneToExercise(msg.Runes[0])
			}
			if m.finished() {
				m.endTime = time.Now()
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m Model) currentCharacterView() string {
	min := func(a int, b int) int {
		if a <= b {
			return a
		}
		return b
	}
	typedEnd := min(len(m.typedExercise), len(m.exercise)-1)
	currChar := rune(m.exercise[typedEnd])
	charString := string(currChar)
	if currChar == Enter {
		charString = Arrow
	}
	return fmt.Sprintf("Curr character: %#U %d %s", currChar, currChar, charString)
}

func (m Model) nameView() string {
	commentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Italic(true)
	commentPrefix := "//"
	return commentStyle.Render(fmt.Sprintf("%s %s", commentPrefix, m.title))
}

func (m Model) View() string {
	s := ""
	if !m.finished() {
		s += "\n"
		s += m.nameView()
		s += "\n\n"
		s += m.ExerciseView()
		s += "\n\n"
		s += m.currentCharacterView()
		s += "\n"
	}
	return s
}

func RunSession(t string, ex string) (m Model) {
	title := t
	exercise := ex
	model, err := tea.NewProgram(initialModel(title, exercise)).Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	return model.(Model)
}
