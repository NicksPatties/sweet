package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"

	"github.com/charmbracelet/lipgloss"
)

func getExercisesDirectory() (string, error) {
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(hd, ".sweet", "exercises"), nil
}

func GetRandomExercise() (string, string, error) {
	dirPath, err := getExercisesDirectory()
	paths, err := getAllFilePathsInDirectory(dirPath)
	if err != nil {
		log.Fatalf("not sweet... an error ocurred: %s", err)
	}
	randI := rand.Intn(len(paths))

	return GetExerciseFromFile(paths[randI])
}

// Gets an exercise for the matching lang file extension
func GetExerciseForLang(lang string) (string, string, error) {
	dirPath, err := getExercisesDirectory()
	exPath, err := getRandomFilePathFromDirectory(dirPath + "/" + lang)
	if err != nil {
		return "", "", errors.New("Failed to find exercise of type " + lang)
	}
	return GetExerciseFromFile(exPath)
}

func GetExerciseFromFile(fileName string) (string, string, error) {
	exercise, err := os.ReadFile(fileName)
	return fileName, string(exercise), err
}

func isWhitespace(rn rune) bool {
	return rn == Tab || rn == Space
}

func (m Model) AddRuneToExercise(rn rune) Model {
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

func (m Model) DeleteCharacter() Model {
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

type Theme struct {
	typedStyle     lipgloss.Style
	untypedStyle   lipgloss.Style
	cursorStyle    lipgloss.Style
	incorrectStyle lipgloss.Style
}

func DefaultTheme() Theme {
	return Theme{
		typedStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")),
		untypedStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
		cursorStyle:    lipgloss.NewStyle().Background(lipgloss.Color("14")),
		incorrectStyle: lipgloss.NewStyle().Background(lipgloss.Color("1")).Foreground(lipgloss.Color("15")),
	}
}

// Returns the exercise string with the typed string overlaid on top of it. Renders
// correctly typed characters with white text, incorrectly typed characters with a
// red background, and characters that haven't been typed yet with gray text.
func (m Model) ExerciseView(args ...Theme) string {
	var t Theme
	if len(args) == 0 {
		t = DefaultTheme()
	} else {
		t = args[0]
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

func (m Model) GetExerciseRuneCount() int {
	ex := m.exercise
	c := 0
	hitNewline := false
	for _, rn := range ex {
		if isWhitespace(rn) && hitNewline {
			continue
		}
		hitNewline = rn == Enter
		c++
	}
	return c
}
