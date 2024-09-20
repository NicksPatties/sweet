package exercise

import (
	"fmt"
	"os"

	lg "github.com/charmbracelet/lipgloss"
)

func isWhitespace(rn rune) bool {
	return rn == Tab || rn == Space
}

func (m sessionModel) addRuneToExercise(rn rune) sessionModel {
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

func (m sessionModel) deleteCharacter() sessionModel {
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

// Returns the exercise string with the typed string overlaid on top of it. Renders
// correctly typed characters with white text, incorrectly typed characters with a
// red background, and characters that haven't been typed yet with gray text.
func (m sessionModel) exerciseView() string {
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

	// check if the $HOME/.sweet directory is there, create the directory, and then add the default exercises

	// run the session
	m := RunSession()

	if m.quitEarly {
		fmt.Println("Goodbye!")
		os.Exit(0)
	}

	// show the results
	showResults(m)
}
