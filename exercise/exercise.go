package exercise

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func getExercisesDirectory() (string, error) {
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(hd, ".sweet", "exercises"), nil
}

func getRandomExercise() (string, string, error) {
	dirPath, err := getExercisesDirectory()
	paths, err := getAllFilePathsInDirectory(dirPath)
	if err != nil {
		log.Fatalf("not sweet... an error ocurred: %s", err)
	}
	randI := rand.Intn(len(paths))

	return getExerciseFromFile(paths[randI])
}

func listExercises() (string, error) {
	ePath, err := getDefaultExercisesPath()
	if err != nil {
		return "", err
	}

	paths, err := getAllFilePathsInDirectory(ePath)
	if err != nil {
		return "", err
	}

	exercises := ""
	for _, path := range paths {
		str := strings.Replace(path, ePath, "", 1)
		exercises += fmt.Sprintln(str[1:])
	}
	return exercises, nil
}

// Gets an exercise for the matching lang file extension
func getExerciseForLang(lang string) (string, string, error) {
	dirPath, err := getExercisesDirectory()
	// get all the files in the exercises directory
	r, err := regexp.Compile("[[:alnum:]]+." + lang)
	exercisesList, err := listExercises()
	exercises := r.FindAllString(exercisesList, -1)
	if exercises == nil {
		return "", "", fmt.Errorf("Failed to find exercise of type %s", lang)
	}
	randI := rand.Intn(len(exercises))
	ex := exercises[randI]
	exPath := path.Join(dirPath, ex)
	if err != nil {
		return "", "", fmt.Errorf("Failed to find exercise of type %s", lang)
	}
	return getExerciseFromFile(exPath)
}

func getExerciseFromFile(fileName string) (string, string, error) {
	exercise, err := os.ReadFile(fileName)
	return fileName, string(exercise), err
}

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
	typedStyle     lipgloss.Style
	untypedStyle   lipgloss.Style
	cursorStyle    lipgloss.Style
	incorrectStyle lipgloss.Style
}

func defaultTheme() theme {
	return theme{
		typedStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")),
		untypedStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
		cursorStyle:    lipgloss.NewStyle().Background(lipgloss.Color("14")),
		incorrectStyle: lipgloss.NewStyle().Background(lipgloss.Color("1")).Foreground(lipgloss.Color("15")),
	}
}

// Returns the exercise string with the typed string overlaid on top of it. Renders
// correctly typed characters with white text, incorrectly typed characters with a
// red background, and characters that haven't been typed yet with gray text.
func (m sessionModel) exerciseView(args ...theme) string {
	var t theme
	if len(args) == 0 {
		t = defaultTheme()
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

func (m sessionModel) getExerciseRuneCount() int {
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

func Run() {

	// check if the $HOME/.sweet directory is there, create the directory, and then add the default exercises
	err := addDefaultExercises()
	if err != nil {
		log.Fatalf("Whoops! %s", err.Error())
	}

	// run the session
	m := RunSession()

	if m.quitEarly {
		fmt.Println("Goodbye!")
		os.Exit(0)
	}

	// show the results
	showResults(m)
}
