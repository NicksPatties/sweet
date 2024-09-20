package exercise

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	lg "github.com/charmbracelet/lipgloss"
)

func isWhitespace(rn rune) bool {
	return rn == Tab || rn == Space
}

type exercise struct {
	name string
	text string
}

type exerciseModel struct {
	title         string
	exercise      string
	typedExercise string
	quitEarly     bool
	startTime     time.Time
	endTime       time.Time
}

func NewExerciseModel(t string, ex string) exerciseModel {
	return exerciseModel{
		title:         t,
		exercise:      ex,
		typedExercise: "",
		quitEarly:     false,
		startTime:     time.Time{},
		endTime:       time.Time{},
	}
}

func (m exerciseModel) addRuneToExercise(rn rune) exerciseModel {
	if len(m.typedExercise) == len(m.exercise) {
		return m
	}

	idx := len(m.typedExercise)
	if rune(m.exercise[idx]) == Enter {
		whiteSpace := []rune{}
		for i := len(m.typedExercise) + 1; i < len(m.exercise) && isWhitespace(rune(m.exercise[i])); i++ {
			whiteSpace = append(whiteSpace, rune(m.exercise[i]))
		}
		m.typedExercise += string(rn)
		m.typedExercise += string(whiteSpace)
		return m
	}
	m.typedExercise += string(rn)
	return m
}

func (m exerciseModel) deleteCharacter() exerciseModel {
	tex := m.typedExercise
	l := len(tex)

	if l <= 0 {
		m.typedExercise = tex
		return m
	}

	currRn := rune(tex[l-1])

	if !isWhitespace(currRn) {
		m.typedExercise = tex[:l-1]
		return m
	}

	m.typedExercise = tex[:l-1]
	l = len(m.typedExercise)
	i := 1
	// move index backwards until a non-whitespace rune is found
	for ; isWhitespace(rune(m.exercise[l-i])); i++ {
	}
	currRn = rune(m.exercise[l-i])
	if currRn == Enter {
		// remove all runes up to and including the newline rune
		m.typedExercise = tex[:l-i]
	}
	return m
}

type theme struct {
	typedStyle     lg.Style
	untypedStyle   lg.Style
	cursorStyle    lg.Style
	incorrectStyle lg.Style
}

func (m exerciseModel) finished() bool {
	// If the user hasn't reached the end of the exercise,
	// then they're not done yet.
	l := len(m.exercise)
	if len(m.typedExercise) < l {
		return false
	}

	// Handle the case where the user types the last character incorrectly
	exLast := rune(m.exercise[l-1])
	typedLast := rune(m.typedExercise[l-1])

	if exLast != typedLast {
		return false
	}
	return true
}

func (m exerciseModel) Init() tea.Cmd {
	return nil
}

func (m exerciseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitEarly = true
			return m, tea.Quit
		case tea.KeyBackspace:
			m = m.deleteCharacter()
		case tea.KeyRunes, tea.KeySpace, tea.KeyEnter:
			if m.startTime.IsZero() {
				m.startTime = time.Now()
			}
			if msg.Type == tea.KeyEnter {
				m = m.addRuneToExercise(Enter)
			} else {
				m = m.addRuneToExercise(msg.Runes[0])
			}
			if m.finished() {
				m.endTime = time.Now()
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m exerciseModel) currentCharacterView() string {
	typedEnd := min(len(m.typedExercise), len(m.exercise)-1)
	currChar := rune(m.exercise[typedEnd])
	charString := string(currChar)
	if currChar == Enter {
		charString = Arrow
	}
	return fmt.Sprintf("Curr character: %#U %d %s", currChar, currChar, charString)
}

func (m exerciseModel) nameView() string {
	commentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Italic(true)
	commentPrefix := "//"
	return commentStyle.Render(fmt.Sprintf("%s %s", commentPrefix, m.title))
}

func (m exerciseModel) View() string {
	s := ""
	if !m.finished() {
		s += "\n"
		s += m.nameView()
		s += "\n\n"
		s += m.exerciseView()
		s += "\n"
	}
	return s
}

// Returns the exercise string with the typed string overlaid on top of it. Renders
// correctly typed characters with white text, incorrectly typed characters with a
// red background, and characters that haven't been typed yet with gray text.
func (m exerciseModel) exerciseView() string {
	t := theme{
		typedStyle:     lg.NewStyle().Foreground(lg.Color("#FFFFFF")),
		untypedStyle:   lg.NewStyle().Foreground(lg.Color("7")),
		cursorStyle:    lg.NewStyle().Background(lg.Color("255")).Foreground(lg.Color("0")),
		incorrectStyle: lg.NewStyle().Background(lg.Color("1")).Foreground(lg.Color("15")),
	}
	ts, us, cs, is := t.typedStyle, t.untypedStyle, t.cursorStyle, t.incorrectStyle
	s := ""

	typed := m.typedExercise

	for i, exRune := range m.exercise {
		// Has this character been typed yet?
		if i > len(typed) {
			s += us.Render(string(exRune))
			continue
		}

		// Is this the cursor?
		if i == len(typed) {

			// Is the cursor on a newline?
			if exRune == Enter {
				s += fmt.Sprintf("%s\n", cs.Render(Arrow))
				continue
			}

			s += cs.Render(string(exRune))
			continue
		}

		// There's at least a typed character at this point...
		typedRune := rune(typed[i])

		// Is it incorrect?
		if typedRune != exRune {
			if exRune == Enter {
				s += fmt.Sprintf("%s\n", is.Render(Arrow))
			} else {
				s += is.Render(string(exRune))
			}

			continue
		}

		s += ts.Render(string(exRune))
	}

	return s
}

func Run() {

	// Get an exercise.
	theExercise := exercise{
		name: "simple",
		text: "the text of the exercise",
	}

	// Create the new Exercise Model
	model := NewExerciseModel(theExercise.name, theExercise.text)
	teaModel, err := tea.NewProgram(model).Run()

	model, _ = teaModel.(exerciseModel)

	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	if model.quitEarly {
		fmt.Println("Goodbye!")
		os.Exit(0)
	}

	// show the results
	showResults(model)
}
