package root

import (
	"fmt"
	"os"
	"time"

	c "github.com/NicksPatties/sweet/constants"
	"github.com/NicksPatties/sweet/db"
	"github.com/NicksPatties/sweet/util"
	tea "github.com/charmbracelet/bubbletea"

	lg "github.com/charmbracelet/lipgloss"
)

type exercise struct {
	// The name of the exercise. Usually the file name.
	name string
	// The contents of the exercise. This is what the user types
	// during the typing game.
	text string
}

// The exercise model used by bubbletea.
//
// Implements tea.Model. Stores the state of the currently running exercise.
type exerciseModel struct {
	exercise exercise

	// The charcters that the user has typed during this exercise.
	typedText string

	// The point in time when the user started typing.
	startTime time.Time

	// The time the user completed the exercise.
	endTime time.Time

	quitEarly bool

	events []util.Event
}

func (m exerciseModel) exerciseNameView() string {
	commentStyle := lg.NewStyle().Foreground(lg.Color("7")).Italic(true)
	commentPrefix := "//"
	return commentStyle.Render(fmt.Sprintf("%s %s", commentPrefix, m.exercise.name))
}

func (m exerciseModel) exerciseTextView() (s string) {
	// typed style
	ts := lg.NewStyle()
	// untyped style
	us := lg.NewStyle().Foreground(lg.Color("7"))
	// cursor style
	cs := lg.NewStyle().Background(lg.Color("15")).Foreground(lg.Color("0"))
	// incorrest style
	is := lg.NewStyle().Background(lg.Color("1")).Foreground(lg.Color("15"))

	typed := m.typedText

	for i, exRune := range m.exercise.text {
		// Has this character been typed yet?
		if i > len(typed) {
			s += us.Render(string(exRune))
			continue
		}

		// Is this the cursor?
		if i == len(typed) {

			// Is the cursor on a newline?
			if exRune == c.Enter {
				s += fmt.Sprintf("%s\n", cs.Render(c.Arrow))
				continue
			}

			s += cs.Render(string(exRune))
			continue
		}

		// There's at least a typed character at this point...
		typedRune := rune(typed[i])

		// Is it incorrect?
		if typedRune != exRune {
			if exRune == c.Enter {
				s += fmt.Sprintf("%s\n", is.Render(c.Arrow))
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
	if len(m.typedText) == len(m.exercise.text) {
		return m
	}

	idx := len(m.typedText)

	// If the next character is an Enter,
	// then add the Enter and the following whitespace to the typedText.
	//
	// This provides the appearance of auto-indentation while typing.
	if rune(m.exercise.text[idx]) == c.Enter {
		whiteSpace := []rune{}
		for i := len(m.typedText) + 1; i < len(m.exercise.text) && util.IsWhitespace(rune(m.exercise.text[i])); i++ {
			whiteSpace = append(whiteSpace, rune(m.exercise.text[i]))
		}
		m.typedText += string(rn)
		m.typedText += string(whiteSpace)
	} else {
		m.typedText += string(rn)
	}
	return m
}

func (m exerciseModel) deleteRuneFromTypedText() exerciseModel {
	tex := m.typedText
	l := len(tex)

	if l <= 0 {
		m.typedText = tex
		return m
	}

	currRn := rune(tex[l-1])

	if !util.IsWhitespace(currRn) {
		m.typedText = tex[:l-1]
		return m
	}

	m.typedText = tex[:l-1]
	l = len(m.typedText)
	i := 1
	// move index backwards until a non-whitespace rune is found
	for ; util.IsWhitespace(rune(m.exercise.text[l-i])); i++ {
	}
	currRn = rune(m.exercise.text[l-i])
	if currRn == c.Enter {
		// remove all runes up to and including the newline rune
		m.typedText = tex[:l-i]
	}
	return m
}

func (m exerciseModel) finished() bool {
	// If the user hasn't reached the end of the exercise,
	// then they're not done yet.
	l := len(m.exercise.text)
	if len(m.typedText) < l {
		return false
	}

	// Handle the case where the user types the last character incorrectly
	exLast := rune(m.exercise.text[l-1])
	typedLast := rune(m.typedText[l-1])

	if exLast != typedLast {
		return false
	}
	return true
}

// Converts the exercise model to a Rep, in preparation for
// inserting it into the database.
func exerciseModelToRep(m exerciseModel) db.Rep {
	return db.Rep{
		Hash:   util.MD5Hash(m.exercise.text),
		Start:  m.events[0].Ts,
		End:    m.events[len(m.events)-1].Ts,
		Name:   m.exercise.name,
		Lang:   util.Lang(m.exercise.name),
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		var currTyped string
		currI := len(m.typedText)
		currExpected := util.RuneToEventExpected(rune(m.exercise.text[currI]))
		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitEarly = true
			return m, tea.Quit
		case tea.KeyBackspace:
			currTyped = util.TeaKeyMsgToEventTyped(msg)
			m = m.deleteRuneFromTypedText()
			// Create delete event and add it to events
			m.events = append(m.events, util.NewEvent("backspace", "", currI))
		case tea.KeyRunes, tea.KeySpace, tea.KeyEnter:
			currTyped = util.TeaKeyMsgToEventTyped(msg)

			if m.startTime.IsZero() {
				m.startTime = time.Now()
			}
			if msg.Type == tea.KeyEnter {
				m = m.addRuneToTypedText(c.Enter)
			} else {
				m = m.addRuneToTypedText(msg.Runes[0])
			}
			m.events = append(m.events, util.NewEvent(currTyped, currExpected, currI))
			if m.finished() {
				m.endTime = time.Now()
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

// Displays the text for the typing exercise.
// Hides the view once the exercise is complete or the user quits early.
func (m exerciseModel) View() (s string) {
	if !m.finished() {
		s += "\n"
		s += m.exerciseNameView()
		s += "\n\n"
		s += m.exerciseTextView()
		s += "\n"

		currKeyI := len(m.typedText)
		currKey := m.exercise.text[currKeyI]
		s += qwerty.render(string(currKey))
		s += "\n"
		s += fingerView(qwerty.fingersMargin, '*', rune(currKey))
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
func run(exercise exercise) {
	teaModel, err := tea.NewProgram(exerciseModel{
		exercise:  exercise,
		typedText: "",
		quitEarly: false,
		startTime: time.Time{},
		endTime:   time.Time{},
		events:    []util.Event{},
	}).Run()
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

	rep := exerciseModelToRep(exModel)

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
