package root

import (
	"fmt"
	"os"
	"time"

	consts "github.com/NicksPatties/sweet/constants"
	"github.com/NicksPatties/sweet/db"
	"github.com/NicksPatties/sweet/event"
	"github.com/NicksPatties/sweet/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// The exercise model used by bubbletea.
//
// Implements tea.Model. Stores the state of the currently running exercise.
type exerciseModel struct {
	// The name of the exercise (typically a file name)
	name string

	// The text to type in the exercise (typically the file's contents)
	text string

	// The charcters that the user has typed during this exercise.
	typedText string

	startTime time.Time
	endTime   time.Time
	quitEarly bool

	// The user's keystrokes during the exercise
	events []event.Event

	viewOptions *viewOptions
}

func (m exerciseModel) renderName() string {
	commentStyle := m.viewOptions.styles.commentStyle
	commentPrefix := "//"
	return commentStyle.Render(fmt.Sprintf("%s %s", commentPrefix, m.name))
}

func (m exerciseModel) renderText() (s string) {
	lines := util.Lines(m.text)
	typedLines := typedLines(lines, m.typedText)

	windowSize := int(m.viewOptions.windowSize)
	currLine := currentLineI(lines, m.typedText)
	linesBefore := windowSize / 3
	linesAfter := windowSize * 2 / 3
	var windowStart, windowEnd int
	switch {
	case currLine < linesBefore:
		windowStart = 0
		windowEnd = windowSize
	case currLine >= linesBefore && currLine < len(lines)-linesAfter:
		windowStart = currLine - linesBefore
		windowEnd = windowStart + windowSize
	case currLine >= len(lines)-linesAfter:
		windowEnd = len(lines)
		windowStart = windowEnd - windowSize
	}

	// should show the whole exercise
	if windowSize == 0 {
		windowStart = 0
		windowEnd = len(lines)
	}

	vignetteLastLine := true
	if windowEnd == len(lines) {
		vignetteLastLine = false
	}

	for i := windowStart; i < windowEnd; i = i + 1 {
		text := lines[i]
		var typed *string = nil
		if i < len(typedLines) {
			typed = &typedLines[i]
		}
		isCurrLine := i == currLine
		shouldVignette := false
		if vignetteLastLine && i == windowEnd-1 && i != windowStart {
			shouldVignette = true
		}
		line := renderLine(text, typed, m.viewOptions.styles, shouldVignette, isCurrLine)
		if lastLine := i == windowEnd-1; lastLine {
			line = removeLastNewline(line)
		}
		s += line
	}
	return
}

// Returns an array of strings that map the typed characters
// to the exercise characters. If no characters have been typed
// on a current line, the typedLine will be nil.
func typedLines(lines []string, typed string) []string {
	typedLines := []string{}
	i := 0
	for _, line := range lines {
		str := ""
		for range line {
			if i >= len(typed) {
				continue
			}
			str = str + string(typed[i])
			i = i + 1
		}
		if str != "" {
			typedLines = append(typedLines, str)
		}
	}
	return typedLines
}

func currentLineI(lines []string, typed string) int {
	typedLen := len(typed)
	for i := range lines {
		for range lines[i] {
			if typedLen == 0 {
				return i
			}
			typedLen = typedLen - 1
		}
	}
	return 0
}

func removeLastNewline(str string) string {
	n := '\n'
	i := len(str) - 1
	for ; i >= 0 && rune(str[i]) != n; i = i - 1 {
	}

	if i < 0 {
		return str
	}

	return str[:i] + str[i+1:]
}

func renderLine(text string, typedP *string, style styles, vignette bool, currLine bool) (s string) {
	typedStyle := style.typedStyle
	untypedStyle := style.untypedStyle
	cursorStyle := style.cursorStyle
	mistakeStyle := style.mistakeStyle

	if vignette {
		typedStyle = style.vignetteStyle
		untypedStyle = style.vignetteStyle
		cursorStyle = style.vignetteStyle
	}

	if typedP == nil {
		for i, c := range text {
			currChar := string(c)
			if c != '\n' {
				currChar = untypedStyle.Render(string(c))
			}
			if i == 0 && currLine {
				currChar = renderVisibleRune(cursorStyle, c)
			}
			s += currChar
		}
		return
	}

	typed := *typedP

	for i, exRune := range text {
		typedYet := i > len(typed)
		isCursor := i == len(typed) && currLine
		isMistake := false
		if i < len(typed) {
			typedRune := rune(typed[i])
			isMistake = typedRune != exRune
		}
		switch {
		case typedYet:
			s += untypedStyle.Render(string(exRune))
		case isCursor:
			s += renderVisibleRune(cursorStyle, exRune)
		case isMistake:
			s += renderVisibleRune(mistakeStyle, exRune)
		default:
			s += typedStyle.Render(string(exRune))
		}
	}
	return
}

// If the rune is a newline, and it needs to be visible
// (i.e. it's a cursor character, or a mistake), then use this function
func renderVisibleRune(style lipgloss.Style, exRune rune) (s string) {
	s = style.Render(string(exRune))
	if exRune == '\n' {
		s = fmt.Sprintf("%s\n", style.Render(consts.Arrow))
	}
	return
}

func (m exerciseModel) addRuneToTypedText(rn rune) exerciseModel {
	if len(m.typedText) == len(m.text) {
		return m
	}

	idx := len(m.typedText)

	// If the next character is an Enter,
	// then add the Enter and the following whitespace to the typedText.
	//
	// This provides the appearance of auto-indentation while typing.
	if rune(m.text[idx]) == consts.Enter {
		whiteSpace := []rune{}
		for i := len(m.typedText) + 1; i < len(m.text) && util.IsWhitespace(rune(m.text[i])); i++ {
			whiteSpace = append(whiteSpace, rune(m.text[i]))
		}
		m.typedText += string(rn)
		m.typedText += string(whiteSpace)
		return m
	}
	m.typedText += string(rn)
	return m
}

func (m exerciseModel) deleteRuneFromTypedText() exerciseModel {
	typed := m.typedText
	l := len(typed)

	if l <= 0 {
		m.typedText = typed
		return m
	}

	currRn := rune(typed[l-1])

	if !util.IsWhitespace(currRn) {
		m.typedText = typed[:l-1]
		return m
	}

	m.typedText = typed[:l-1]
	l = len(m.typedText)
	i := 1
	// move index backwards until a non-whitespace rune is found
	for ; util.IsWhitespace(rune(m.text[l-i])); i++ {
	}
	currRn = rune(m.text[l-i])
	if currRn == consts.Enter {
		// remove all runes up to and including the newline rune
		m.typedText = typed[:l-i]
	}
	return m
}

func (m exerciseModel) finished() bool {
	// If the user hasn't reached the end of the exercise,
	// then they're not done yet.
	l := len(m.text)
	if len(m.typedText) < l {
		return false
	}

	// Handle the case where the user types the last character incorrectly
	exLast := rune(m.text[l-1])
	typedLast := rune(m.typedText[l-1])

	if exLast != typedLast {
		return false
	}
	return true
}

// Converts the exercise model to a Rep, in preparation for
// inserting it into the database.
func (m exerciseModel) Rep() db.Rep {
	return db.Rep{
		Hash:   util.MD5Hash(m.text),
		Start:  m.events[0].Ts,
		End:    m.events[len(m.events)-1].Ts,
		Name:   m.name,
		Lang:   util.Lang(m.name),
		Wpm:    wpm(m.events),
		Raw:    wpmRaw(m.events),
		Dur:    duration(m.events),
		Acc:    accuracy(m.events),
		Miss:   numMistakes(m.events),
		Errs:   numUncorrectedErrors(m.events),
		Events: m.events,
	}
}

func (m exerciseModel) Init() tea.Cmd {
	return nil
}

func (m exerciseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	var currTyped string
	currI := len(m.typedText)
	currExpected := event.RuneToEventExpected(rune(m.text[currI]))
	switch keyMsg.Type {
	case tea.KeyCtrlC:
		m.quitEarly = true
		return m, tea.Quit
	case tea.KeyBackspace:
		currTyped = event.TeaKeyMsgToEventTyped(keyMsg)
		m = m.deleteRuneFromTypedText()
		// Create delete event and add it to events
		m.events = append(m.events, event.NewEvent("backspace", "", currI))
	case tea.KeyRunes, tea.KeySpace, tea.KeyEnter:
		currTyped = event.TeaKeyMsgToEventTyped(keyMsg)
		if m.startTime.IsZero() {
			m.startTime = time.Now()
		}
		if keyMsg.Type == tea.KeyEnter {
			m = m.addRuneToTypedText(consts.Enter)
		} else {
			m = m.addRuneToTypedText(keyMsg.Runes[0])
		}
		m.events = append(m.events, event.NewEvent(currTyped, currExpected, currI))
		if m.finished() {
			m.endTime = time.Now()
			return m, tea.Quit
		}
	}
	return m, nil
}

// Displays the text for the typing exercise.
// Hides the view once the exercise is complete or the user quits early.
func (m exerciseModel) View() (s string) {
	if !m.finished() {
		s += "\n"
		s += m.renderName()
		s += "\n\n"
		s += m.renderText()
		s += "\n\n"

		currKeyI := len(m.typedText)
		currKey := m.text[currKeyI]
		s += qwerty.render(string(currKey))
		s += "\n"
		s += renderFingers(qwerty.fingersMargin, '*', rune(currKey))
	}
	return
}

// Runs the exercise, which does the following
//
// 1. Runs the bubbletea interactive typing application.
//
// 2. Depending on the outcome of the exercise, either completes
// or returns an error.
//
// 3. If the exercise is completed, gather the results, print them, and
// save them to the database
func run(name string, text string, options *viewOptions) {
	newModel := exerciseModel{
		name:        name,
		text:        text,
		typedText:   "",
		quitEarly:   false,
		startTime:   time.Time{},
		endTime:     time.Time{},
		events:      []event.Event{},
		viewOptions: options,
	}
	teaModel, err := tea.NewProgram(newModel).Run()
	if err != nil {
		fmt.Printf("Error running typing exercise: %v\n", err)
		os.Exit(1)
	}

	exModel, ok := teaModel.(exerciseModel)
	if !ok {
		fmt.Printf("Error casting bubbletea model.\n")
	}
	if exModel.quitEarly {
		os.Exit(0)
	}

	rep := exModel.Rep()

	printExerciseResults(rep)

	// open connection to db once exercise is complete
	statsDb, err := db.SweetDb()
	if err != nil {
		fmt.Println(err)
	}
	// insert the row into the database
	var repId int64
	repId, err = db.InsertRep(statsDb, rep)
	if err != nil {
		fmt.Printf("Error saving rep to the database: %v\n", err)
	} else {
		fmt.Printf("Rep %d saved to the database! Keep it up!\n", repId)
	}
}
