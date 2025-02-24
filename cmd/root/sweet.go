// This package includes code for initializing and running the
// `sweet` command. It defines the flags and subcommands for the CLI,
// handles processing arguments, and reading files.
package root

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"strings"

	"github.com/NicksPatties/sweet/cmd/about"
	"github.com/NicksPatties/sweet/cmd/add"
	"github.com/NicksPatties/sweet/cmd/stats"
	"github.com/NicksPatties/sweet/cmd/version"
	"github.com/NicksPatties/sweet/util"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type exerciseFile struct {
	name string
	text string
}

// Controls the appearance of the exercise performed
// with this command.
type viewOptions struct {
	// Colors and styles for the exercise text.
	styles styles

	// The number of lines that are visible
	// in the exericse. Default value is 0, meaning
	// show the entire exercise.
	windowSize uint
}

type styles struct {
	commentStyle         lg.Style
	untypedStyle         lg.Style
	cursorStyle          lg.Style
	typedStyle           lg.Style
	mistakeStyle         lg.Style
	vignetteStyle        lg.Style
	vignetteMistakeStyle lg.Style
}

func defaultStyles() styles {
	return styles{
		commentStyle:         lg.NewStyle().Foreground(lg.Color("7")).Italic(true),
		untypedStyle:         lg.NewStyle().Foreground(lg.Color("7")),
		cursorStyle:          lg.NewStyle().Background(lg.Color("15")).Foreground(lg.Color("0")),
		typedStyle:           lg.NewStyle().Foreground(lg.Color("15")),
		mistakeStyle:         lg.NewStyle().Background(lg.Color("1")).Foreground(lg.Color("15")),
		vignetteStyle:        lg.NewStyle().Foreground(lg.Color("8")),
		vignetteMistakeStyle: lg.NewStyle().Background(lg.Color("1")).Foreground(lg.Color("8")),
	}
}

var Cmd = &cobra.Command{
	Use:     "sweet [file]",
	Long:    fmt.Sprintf("%s.\nRuns an interactive touch typing game, and prints the results.", tagline()),
	Args:    cobra.MaximumNArgs(1),
	Example: examples(),
	RunE: func(cmd *cobra.Command, args []string) error {
		exercise, err := exerciseFileFromArgs(cmd, args)
		if err != nil {
			return err
		}
		viewOptions, err := viewOptionsFromArgs(cmd, exercise.text)
		if err != nil {
			return err
		}
		run(exercise.name, exercise.text, viewOptions)
		return nil
	},
}

func tagline() string {
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

func examples() (msg string) {
	msg += fmt.Sprintf("  Run a random exercise\n")
	msg += fmt.Sprintf("  $ sweet\n\n")
	msg += fmt.Sprintf("  Run an exercise from lines 2 to 10 of a file\n")
	msg += fmt.Sprintf("  $ sweet file -s 2 -e 10\n\n")
	msg += fmt.Sprintf("  Run an exercise with STDIN (use `-` as your file)\n")
	msg += fmt.Sprintf("  $ curl https://nickspatties.com/main.go | sweet -")
	return
}

// Validates and returns the exercise from command line arguments.
// If the flags are incorrect, an error is returned.
func exerciseFileFromArgs(cmd *cobra.Command, args []string) (exercise exerciseFile, err error) {
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
			sweetConfigDir, err = util.SweetConfigDir()
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

// Scans a file and returns its text as a string.
// If start or end is defined, only returns the lines between start and end.
// If the file is empty, it returns an empty string.
func scanFileText(file *os.File, start uint, end uint) (text string) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		text += scanner.Text()
	}
	lines := util.Lines(text)
	if end >= uint(len(lines)) {
		end = uint(len(lines))
	}
	return strings.Join(lines[start-1:end], "")
}

var defaultExercises = []exerciseFile{
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

func addDefaultExercises(dir string) (files []os.DirEntry) {
	for _, ex := range defaultExercises {
		os.WriteFile(path.Join(dir, ex.name), []byte(ex.text), 0600)
	}
	files, _ = os.ReadDir(dir)
	return
}

func viewOptionsFromArgs(cmd *cobra.Command, exerciseText string) (*viewOptions, error) {
	windowSize, err := cmd.Flags().GetUint("window-size")
	if err != nil {
		return nil, err
	}
	numLines := uint(len(util.Lines(exerciseText)))
	if windowSize >= numLines {
		windowSize = 0
	}
	return &viewOptions{
		styles:     defaultStyles(),
		windowSize: windowSize,
	}, nil
}

func init() {
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

func setRootCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("language", "l", "", "select a language by file extension")
	cmd.Flags().UintP("start", "s", 1, "start exercise at this line")
	cmd.Flags().UintP("end", "e", math.MaxUint, "end exercise at this line")
	cmd.Flags().UintP("window-size", "w", 0, "set the number of visible lines for the exercise")
	cmd.Flags().SortFlags = false
}
