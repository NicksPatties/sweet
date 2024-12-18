package constants

// Typing exercise game symbols
const (
	Enter   = '\n'
	Tab     = '\t'
	Space   = ' '
	Percent = '%'
	Arrow   = `↲`
)

// Used for words per minute (WPM) calculations.
const WORD_SIZE = 5

// Reps database table column names.
// They're also used as flag names the commands.
const (
	ID                 string = "id"
	HASH               string = "hash"
	START              string = "start"
	END                string = "end"
	NAME               string = "name"
	LANGUAGE           string = "lang"
	WPM                string = "wpm"
	RAW_WPM            string = "raw"
	DURATION           string = "dur"
	ACCURACY           string = "acc"
	MISTAKES           string = "miss"
	UNCORRECTED_ERRORS string = "errs"
	EVENTS             string = "events"
)
