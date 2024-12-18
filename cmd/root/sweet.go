package root

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path"

	"strings"
	"time"

	"github.com/NicksPatties/sweet/cmd/about"
	"github.com/NicksPatties/sweet/cmd/add"
	"github.com/NicksPatties/sweet/cmd/stats"
	"github.com/NicksPatties/sweet/cmd/version"
	. "github.com/NicksPatties/sweet/constants"
	. "github.com/NicksPatties/sweet/db"
	. "github.com/NicksPatties/sweet/util"
	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	lg "github.com/charmbracelet/lipgloss"
)

func getProductTagline() string {
	b := lg.NewStyle().Bold(true)
	return fmt.Sprintf(
		"%s: The %soft%sare %sngineer's %sxercise for %syping",
		b.Render("sweet"),
		b.Render("S"),
		b.Render("w"),
		b.Render("E"),
		b.Render("E"),
		b.Render("T"),
	)
}

func getExamples() (msg string) {
	msg += fmt.Sprintf("  Run a random exercise\n")
	msg += fmt.Sprintf("  $ sweet\n\n")
	msg += fmt.Sprintf("  Run an exercise from lines 2 to 10 of a file\n")
	msg += fmt.Sprintf("  $ sweet file -s 2 -e 10\n\n")
	msg += fmt.Sprintf("  Run an exercise with STDIN (use `-` as your file)\n")
	msg += fmt.Sprintf("  $ curl https://nickspatties.com/main.go | sweet -")
	return
}

var Cmd = &cobra.Command{
	Use:     "sweet [file]",
	Long:    fmt.Sprintf("%s.\nRuns an interactive touch typing game, and prints the results.", getProductTagline()),
	Args:    cobra.MaximumNArgs(1),
	Example: getExamples(),
	Run: func(cmd *cobra.Command, args []string) {
		ex, err := fromArgs(cmd, args)
		if err != nil {
			log.Fatal(err)
		}
		Run(ex)
	},
}

// Exercises that should be added to the
// exercises directory if it's empty.
var defaultExercises = []Exercise{
	{
		name: "sweet_cmd.go",
		text: `var Cmd = &cobra.Command{
	Use:   "sweet [file|-]",
	Short: "The Software Engineer Exercise for Typing.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ex, err := fromArgs(cmd, args)
		if err != nil {
			log.Fatal(err)
		}
		Run(ex)
	},
}
`,
	},
	{
		name: "resume-section.html",
		text: `<section id="themes">
  <h1>Themes</h1>
  <label class="has-checkbox-input">
    <input type="radio" name="resume-theme" value="default" checked />
    <span>Default</span>
  </label>
  <label class="has-checkbox-input">
    <input type="radio" name="resume-theme" value="monospace" />
    <span>Monospace</span>
  </label>
</section>
`,
	},
	{
		name: "portfolio-site-burger.css",
		text: `.hero .burger {
  position: absolute;
  height: 100%;
  top: 0;
  right: 0;
  opacity: 0.66;
  z-index: -1;
  transform: rotate(5deg);
}
`,
	},
}

func setRootCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("language", "l", "", "select a language by file extension")
	cmd.Flags().UintP("start", "s", 0, "start exercise at this line")
	cmd.Flags().UintP("end", "e", math.MaxUint, "end exercise at this line")
	cmd.Flags().SortFlags = false
}

func init() {

	// Add language flag to root command only.
	// The flags for other commands will be defined in their respective modules.
	setRootCmdFlags(Cmd)

	Cmd.CompletionOptions.DisableDefaultCmd = true

	commands := []*cobra.Command{
		about.Cmd,
		add.Cmd,
		version.Cmd,
		stats.Cmd,
	}

	for _, c := range commands {
		Cmd.AddCommand(c)
	}
}

// STRUCTS

// A single Exercise.
//
// This contains the data that is required to display and perform
// the typing Exercise.
type Exercise struct {
	// The name of the exercise. Usually the file name.
	name string
	// The contents of the exercise. The user types this.
	text string
	// A short description that shows when the user complete the exercise.
	completionDescription string
}

// Converts a bubbletea key message to a string.
// Used to properly record key events.
func teaKeyMsgToEventTyped(msg tea.KeyMsg) string {
	switch msg.Type {
	case tea.KeyEnter:
		return "enter"
	case tea.KeyBackspace:
		return "backspace"
	case tea.KeySpace:
		return "space"
	default:
		return string(msg.Runes[0])
	}
}

func runeToEventExpected(r rune) string {
	switch r {
	case '\n':
		return "enter"
	case ' ':
		return "space"
	default:
		return string(r)
	}
}

// The exercise model used by bubbletea.
//
// Implements tea.Model. Stores the state of the currently running exercise.
type exerciseModel struct {
	exercise Exercise

	// The charcters that the user has typed during this exercise.
	typedText string

	// The point in time when the user started typing.
	startTime time.Time

	// The time the user completed the exercise.
	endTime time.Time

	quitEarly bool

	events []Event
}

// INITIALIZATION
func NewExerciseModel(ex Exercise) exerciseModel {
	return exerciseModel{
		exercise:  ex,
		typedText: "",
		quitEarly: false,
		startTime: time.Time{},
		endTime:   time.Time{},
		events:    []Event{},
	}
}

func (m exerciseModel) Init() tea.Cmd {
	return nil
}

// UPDATE

func isWhitespace(rn rune) bool {
	return rn == Tab || rn == Space
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
	if rune(m.exercise.text[idx]) == Enter {
		whiteSpace := []rune{}
		for i := len(m.typedText) + 1; i < len(m.exercise.text) && isWhitespace(rune(m.exercise.text[i])); i++ {
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

	if !isWhitespace(currRn) {
		m.typedText = tex[:l-1]
		return m
	}

	m.typedText = tex[:l-1]
	l = len(m.typedText)
	i := 1
	// move index backwards until a non-whitespace rune is found
	for ; isWhitespace(rune(m.exercise.text[l-i])); i++ {
	}
	currRn = rune(m.exercise.text[l-i])
	if currRn == Enter {
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

func (m exerciseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		var currTyped string
		currI := len(m.typedText)
		currExpected := runeToEventExpected(rune(m.exercise.text[currI]))
		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitEarly = true
			return m, tea.Quit
		case tea.KeyBackspace:
			currTyped = teaKeyMsgToEventTyped(msg)
			m = m.deleteRuneFromTypedText()
			// Create delete event and add it to events
			m.events = append(m.events, NewEvent("backspace", "", currI))
		case tea.KeyRunes, tea.KeySpace, tea.KeyEnter:
			currTyped = teaKeyMsgToEventTyped(msg)

			if m.startTime.IsZero() {
				m.startTime = time.Now()
			}
			if msg.Type == tea.KeyEnter {
				m = m.addRuneToTypedText(Enter)
			} else {
				m = m.addRuneToTypedText(msg.Runes[0])
			}
			m.events = append(m.events, NewEvent(currTyped, currExpected, currI))
			if m.finished() {
				m.endTime = time.Now()
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

// VIEWS

func (m exerciseModel) exerciseNameView() string {
	commentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Italic(true)
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

	return
}

// Displays the text for the typing exercise.
// Hides the view once the exercise is complete or the user quits early.
//
// TODO: Consider leaving the exercise on screen if I'd like to include
// a description of the exercise once it's done.
func (m exerciseModel) View() (s string) {
	if !m.finished() {
		currKeyI := len(m.typedText)
		currKey := m.exercise.text[currKeyI]
		s += "\n"
		s += m.exerciseNameView()
		s += "\n\n"
		s += m.exerciseTextView()
		s += "\n"
		s += qwerty.render(string(currKey))
		s += "\n"
		s += fingerView(qwerty.fingersMargin, '*', rune(currKey))
	}
	return
}

// Scans a file and returns its text as a string.
// If start or end is defined, only returns the lines between start and end.
// If the file is empty, it returns an empty string.
func scanFileText(file *os.File, start uint, end uint) (text string) {
	scanner := bufio.NewScanner(file)
	for line := uint(1); line <= end && scanner.Scan(); line++ {
		if line >= start {
			text += scanner.Text() + "\n"
		}
	}
	return
}

// Add some default exercises to the dir directory.
// Assumes the contents of the directory are empty.
// Returns the dirEntries of the newly added files.
func addDefaultExercises(dir string) (files []os.DirEntry) {
	for _, ex := range defaultExercises {
		os.WriteFile(path.Join(dir, ex.name), []byte(ex.text), 0600)
	}
	files, _ = os.ReadDir(dir)

	return
}

// Validates and returns the exercise from command line arguments.
// If the flags are incorrect, an error is returned.
func fromArgs(cmd *cobra.Command, args []string) (exercise Exercise, err error) {
	start, _ := cmd.Flags().GetUint("start")
	end, _ := cmd.Flags().GetUint("end")

	if start > end {
		err = errors.New(fmt.Sprintf("start flag %d cannot be greater than end flag %d", start, end))
		return
	}

	var file *os.File
	var text string
	defer file.Close()
	if len(args) > 0 { // get the file from the argument
		if args[0] == "-" {
			file = os.Stdin
		} else {
			file, err = os.Open(args[0])
			if err != nil {
				return
			}

		}
		text = scanFileText(file, start, end)
		if text == "" {
			msg := fmt.Sprintf("no text found in file %s. are you sure it's not empty?", file.Name())
			err = errors.New(msg)
			return
		}
	} else { // get a random exercise
		if start != 0 || end != math.MaxUint {
			err = errors.New("start and end should not be assigned for random exercise")
			return
		}

		var exercisesDir string
		if envDir := os.Getenv("SWEET_EXERCISES_DIR"); envDir != "" {
			exercisesDir = envDir

		} else {
			var sweetConfigDir string
			sweetConfigDir, err = SweetConfigDir()
			if err != nil {
				return
			}
			exercisesDir = path.Join(sweetConfigDir, "exercises")
		}

		if err = os.MkdirAll(exercisesDir, 0775); err != nil {
			return
		}

		var entries []os.DirEntry
		entries, err = os.ReadDir(exercisesDir)
		if err != nil {
			return
		}

		language, _ := cmd.Flags().GetString("language")
		var files []os.DirEntry
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			// Gets the file extension
			ext := strings.Split(entry.Name(), ".")[1]
			if language != "" && language != ext {
				continue
			}

			files = append(files, entry)
		}

		numFiles := len(files)
		if numFiles == 0 {
			if language != "" {
				err = errors.New("failed to find exercise matching language " + language)
				return
			}
			fmt.Printf("adding default exercises to the %s directory...\n", exercisesDir)
			files = addDefaultExercises(exercisesDir)
			numFiles = len(files)
		}
		// finding a valid exercise file
		for text == "" {
			randI := rand.Intn(numFiles)
			filePath := path.Join(exercisesDir, files[randI].Name())
			file, err = os.Open(filePath)
			if err != nil {
				return
			}
			text = scanFileText(file, start, end)
			// If there's an empty file in the directory,
			// then warn the user of that weird behavior.
			if text == "" {
				fmt.Printf("warn: found an empty file in the exercises directory: %s\n", exercisesDir)
				numFiles--
				if numFiles == 0 {
					msg := fmt.Sprintf("all files found in the following exercises directory are empty: %s\n", exercisesDir)
					err = errors.New(msg)
					return
				}
				fmt.Println("trying another exercise file...")
				files = append(files[:randI], files[randI+1:]...)
			}
		}
	}

	if text == "" {
		err = errors.New("no input text selected")
		return
	}

	exercise.text = text
	exercise.name = path.Base(file.Name())
	return
}

// Converts the exercise model to a Rep, in preparation for
// inserting it into the database.
func exerciseModelToRep(m exerciseModel) Rep {
	return Rep{
		Hash:   MD5Hash(m.exercise.text),
		Start:  m.events[0].Ts,
		End:    m.events[len(m.events)-1].Ts,
		Name:   m.exercise.name,
		Lang:   Lang(m.exercise.name),
		Wpm:    wpm(m.events),
		Raw:    wpmRaw(m.events),
		Dur:    duration(m.events),
		Acc:    accuracy(m.events),
		Miss:   numMistakes(m.events),
		Errs:   numUncorrectedErrors(m.events),
		Events: m.events,
	}
}

func Run(exercise Exercise) {
	exModel := NewExerciseModel(exercise)
	teaModel, err := tea.NewProgram(exModel).Run()

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
	statsDb, err := SweetDb()
	if err != nil {
		fmt.Println(err)
	}
	// insert the row into the database
	var repId int64
	repId, err = InsertRep(statsDb, rep)
	if err != nil {
		fmt.Printf("Error saving rep to the database: %v\n", err)
	} else {
		fmt.Printf("Rep %d saved to the database! Keep it up!\n", repId)
	}
}
