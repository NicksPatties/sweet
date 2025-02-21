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
	lg "github.com/charmbracelet/lipgloss"
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

	// The point in time when the user started typing.
	startTime time.Time

	// The time the user completed the exercise.
	endTime time.Time

	// True if the user quit before the exercise was complete
	quitEarly bool

	// The user's keystrokes during the exercise
	events []event.Event
}

var styles = struct {
	commentStyle lg.Style
	untypedStyle lg.Style
	cursorStyle  lg.Style
	typedStyle   lg.Style
	mistakeStyle lg.Style
}{
	commentStyle: lg.NewStyle().Foreground(lg.Color("7")).Italic(true),
	untypedStyle: lg.NewStyle().Foreground(lg.Color("7")),
	cursorStyle:  lg.NewStyle().Background(lg.Color("15")).Foreground(lg.Color("0")),
	typedStyle:   lg.NewStyle(),
	mistakeStyle: lg.NewStyle().Background(lg.Color("1")).Foreground(lg.Color("15")),
}

func (m exerciseModel) renderName() string {
	commentStyle := styles.commentStyle
	commentPrefix := "//"
	return commentStyle.Render(fmt.Sprintf("%s %s", commentPrefix, m.name))
}

func (m exerciseModel) renderText() (s string) {
	// typed style
	ts := styles.typedStyle
	// untyped style
	us := styles.untypedStyle
	// cursor style
	cs := styles.cursorStyle
	// incorrest style
	is := styles.mistakeStyle

	typed := m.typedText

	for i, exRune := range m.text {
		// Has this character been typed yet?
		if i > len(typed) {
			s += us.Render(string(exRune))
			continue
		}

		// Is this the cursor?
		if i == len(typed) {

			// Is the cursor on a newline?
			if exRune == consts.Enter {
				s += fmt.Sprintf("%s\n", cs.Render(consts.Arrow))
				continue
			}

			s += cs.Render(string(exRune))
			continue
		}

		// There's at least a typed character at this point...
		typedRune := rune(typed[i])

		// Is it incorrect?
		if typedRune != exRune {
			if exRune == consts.Enter {
				s += fmt.Sprintf("%s\n", is.Render(consts.Arrow))
			} else {
				s += is.Render(string(exRune))
			}

			continue
		}

		s += ts.Render(string(exRune))
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
		s += "\n"

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
func run(name string, text string) {
	newModel := exerciseModel{
		name:      name,
		text:      text,
		typedText: "",
		quitEarly: false,
		startTime: time.Time{},
		endTime:   time.Time{},
		events:    []event.Event{},
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
