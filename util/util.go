package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// Converts a string to an md5 hash. Used to
// convert the contents of an exercise into a string
// to verify if their contents are the same.
//
// see: https://stackoverflow.com/a/25286918
func MD5Hash(contents string) string {
	bytes := []byte(contents)
	hash := md5.Sum(bytes)
	return hex.EncodeToString(hash[:])
}

// Gets the language of the provided filename.
// Unlike `path.Ext`, the language doesn't include the
// leading dot.
func Lang(filename string) (lang string) {
	lang = ""
	split := strings.Split(filename, ".")
	if len(split) > 1 {
		lang = split[len(split)-1]
	}
	return
}

// Gets the path for sweet's configuration directory.
//
// See `os.UserConfigDir` for the default configuration
// location depending on the current operating system.
func SweetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %v", err)
	}
	return path.Join(configDir, "sweet"), nil
}

// Filters a list of file names by the given language extension.
func FilterFileNames(fileNames []string, language string) (found []string) {
	for _, f := range fileNames {
		ext := path.Ext(f)
		// Ignore files that don't have an extension.
		if len(ext) == 0 {
			continue
		}
		if ext[1:] == language {
			found = append(found, f)
		}
	}
	return found
}

func IsValidURL(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}

	// Check if scheme and host are present
	return u.Scheme != "" && u.Host != ""
}

// A recording of a keypress during the exercise.
//
// These are used to perform analysis on the user's performance,
// display stats, and keys that were causing the most trouble.
type Event struct {
	// The moment the event took place.
	Ts time.Time

	// The key that was Typed.
	Typed string

	// The rune that was Expected. Optional, since the user
	// may have pressed backspace.
	Expected string

	// The index of the exercise when the rune was typed.
	I int
}

const EventTsLayout = "2006-01-02 15:04:05.000"

// Converts an event to a string.
func (e Event) String() string {
	time := e.Ts.Format(EventTsLayout)
	return fmt.Sprintf("%s\t%d\t%s\t%s", time, e.I, e.Typed, e.Expected)
}

// Checks if an event has the same timestamp, index, typed
// and expected characters. Used primarily for testing.
func (a Event) Matches(b Event) bool {
	return a.Ts.Equal(b.Ts) &&
		a.I == b.I &&
		a.Typed == b.Typed &&
		a.Expected == b.Expected
}

// Converts an event string to an event struct.
func ParseEvent(line string) (e Event) {
	s := strings.Split(line, "\t")
	e.Ts, _ = time.Parse(EventTsLayout, s[0])
	e.I, _ = strconv.Atoi(s[1])
	e.Typed = s[2]
	if len(s) > 3 {
		e.Expected = s[3]
	}
	return
}

// Same as above, but for a multi-line list of events.
func ParseEvents(list string) (events []Event) {
	for _, line := range strings.Split(list, "\n") {
		if line != "" {
			events = append(events, ParseEvent(line))
		}
	}
	return
}

// Returns a string of an array of events.
func EventsString(events []Event) (s string) {
	s += fmt.Sprintln("[")
	for _, e := range events {
		s += fmt.Sprintf("  %s\n", e)
	}
	s += fmt.Sprintln("]")
	return
}

// Creates a new event. Should be used when recording a keystroke
// to the model.
func NewEvent(typed string, expected string, i int) Event {
	return Event{
		Ts:       time.Now(),
		Typed:    typed,
		Expected: expected,
		I:        i,
	}
}
