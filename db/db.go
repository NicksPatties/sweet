package db

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/NicksPatties/sweet/constants"
	"github.com/NicksPatties/sweet/event"
	"github.com/NicksPatties/sweet/util"

	_ "modernc.org/sqlite"
)

// Column names

// A single repetition of an exercise.
//
// This is the structure of a row in the sweet db module.
// See `db/db.go` for more details.
type Rep struct {
	Id     int
	Hash   string
	Start  time.Time // use
	End    time.Time
	Name   string
	Lang   string
	Wpm    float64
	Raw    float64
	Dur    time.Duration
	Acc    float64
	Miss   int
	Errs   int
	Events event.Events // one string is one event
}

func (r Rep) String() (s string) {
	s += fmt.Sprintf("rep:\n")
	s += fmt.Sprintf("  hash:  %s\n", r.Hash)
	s += fmt.Sprintf("  start: %s\n", r.Start)
	s += fmt.Sprintf("  end:   %s\n", r.End)
	s += fmt.Sprintf("  name:  %s\n", r.Name)
	s += fmt.Sprintf("  lang:  %s\n", r.Lang)
	s += fmt.Sprintf("  wpm:   %.f\n", r.Wpm)
	s += fmt.Sprintf("  dur:   %s\n", r.Dur)
	s += fmt.Sprintf("  acc:   %2.f%%\n", r.Acc)
	s += fmt.Sprintf("  miss:  %d\n", r.Miss)
	s += fmt.Sprintf("  errs:  %d\n", r.Errs)
	s += fmt.Sprintf("  events: %d events\n", len(r.Events))
	return
}

// Returns a string with a decorated version of the asked
// column corresponding to the matching property of a Rep struct.
//
// If the column string doesn't match any Rep column, then
// an empty string is returned.
func (r Rep) ColumnString(col string) string {
	switch col {
	case constants.ID:
		return strconv.Itoa(r.Id)
	case constants.HASH:
		return r.Hash
	case constants.START:
		return r.Start.Format(event.EventTsLayout)
	case constants.END:
		return r.End.Format(event.EventTsLayout)
	case constants.NAME:
		return r.Name
	case constants.LANGUAGE:
		return r.Lang
	case constants.WPM:
		return fmt.Sprintf("%.f", r.Wpm)
	case constants.RAW_WPM:
		return fmt.Sprintf("%.f", r.Raw)
	case constants.DURATION:
		return r.Dur.Round(time.Millisecond).String()
	case constants.ACCURACY:
		return fmt.Sprintf("%.2f%%", r.Acc)
	case constants.MISTAKES:
		return strconv.Itoa(r.Miss)
	case constants.UNCORRECTED_ERRORS:
		return strconv.Itoa(r.Errs)
	case constants.EVENTS:
		return fmt.Sprintf("%s", r.Events)
	default:
		return ""
	}
}

// Gets a pointer to the stats database. If the database file
// doesn't exist already, it will be created at sweet's default
// configuration location (`~/.config/sweet`), or at the path
// specified by `SWEET_DB_LOCATION`, if it's defined.
//
// If an error is returned from this function, then the pointer
// will be `nil`.
func SweetDb() (*sql.DB, error) {
	// get the path for the database
	var dbPath string
	if envDir := os.Getenv("SWEET_DB_LOCATION"); envDir != "" {
		dbPath = envDir
	} else {
		sweetDir, err := util.SweetConfigDir()
		if err != nil {
			return nil, fmt.Errorf("failed to find user config directory: %v", err)
		}
		dbPath = sweetDir
	}

	// create the sweet config directory
	if err := os.MkdirAll(dbPath, 0775); err != nil {
		return nil, fmt.Errorf("failed to find or create sweet config directory: %v", err)
	}

	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite", path.Join(dbPath, "sweet.db"))
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// create the reps table if it doesn't exist
	createTableStr := fmt.Sprintf(`
CREATE TABLE if not exists reps(
  %s integer primary key autoincrement not null,
  -- md5 hash of the exercise file's contents
  %s string NOT NULL,
  -- start time in unix milliseconds
  %s integer not null,
  -- end time in unix milliseconds
  %s integer not null,
  -- name of the exercise file, includes extension if present
  %s text not null,
  -- language: extension of the exercise file, or "" if there is none.
  %s text,
  -- words per minute
  %s real not null check(wpm >= 0.0),
  -- raw words per minute
  %s real not null check(raw >= 0.0),
  -- duration: duration of rep in **nanoseconds**
  %s integer not null check(dur >= 0),
  -- accuracy: float between [0, 100]
  %s real not null check(acc >= 0.0),
  -- mistakes: must be gte 0
  %s integer not null check(miss >= 0),
  -- uncorrected errors: must be gte 0
  %s integer not null check(errs >= 0),
  -- array of events, events are separated by '\n'
  %s text not null
);`,
		constants.ID, constants.HASH, constants.START, constants.END, constants.NAME, constants.LANGUAGE, constants.WPM, constants.RAW_WPM,
		constants.DURATION, constants.ACCURACY, constants.MISTAKES, constants.UNCORRECTED_ERRORS, constants.EVENTS,
	)

	_, err = db.Exec(createTableStr)

	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return db, nil
}

func eventsStringToColumn(events []event.Event) (s string) {
	for i, event := range events {
		s += event.String()
		if i != len(events)-1 {
			s += "\n"
		}
	}
	return
}

// Inserts a repetition into the database.
// On successful insert, returns the id of the inserted row and nil.
// If an error is returned, the returned id is 0.
func InsertRep(db *sql.DB, rep Rep) (int64, error) {
	hash := rep.Hash
	start := rep.Start.UnixMilli()
	end := rep.End.UnixMilli()
	name := rep.Name
	lang := rep.Lang
	wpm := rep.Wpm
	raw := rep.Raw
	dur := rep.Dur
	acc := rep.Acc
	miss := rep.Miss
	errs := rep.Errs
	events := eventsStringToColumn(rep.Events)
	query := fmt.Sprintf(`insert into reps (
	    %s, %s, %s, %s, %s, %s,
	    %s, %s, %s, %s, %s, %s
	   ) values (
	   	?, ?, ?, ?, ?, ?,
	   	?, ?, ?, ?, ?, ?
	   );`,
		constants.HASH, constants.START, constants.END, constants.NAME, constants.LANGUAGE, constants.WPM, constants.RAW_WPM,
		constants.DURATION, constants.ACCURACY, constants.MISTAKES, constants.UNCORRECTED_ERRORS, constants.EVENTS,
	)

	result, err := db.Exec(query,
		hash, start, end, name, lang, wpm,
		raw, dur, acc, miss, errs, events,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// This should accept an array of columns to show, a start and an end range,
// and return an array of anything
func GetReps(db *sql.DB, query string) ([]Rep, error) {
	if query == "" {
		query = fmt.Sprintf(`select * from reps order by %s;`, constants.START)
	}

	var reps []Rep

	// Query to retrieve data
	rows, err := db.Query(query)
	if err != nil {
		return reps, err
	}
	defer rows.Close()

	// Iterate through results
	for rows.Next() {
		// looking for start, name, wpm, errs, dur, miss, acc
		var (
			id     int64
			hash   string
			start  int64
			end    int64
			name   string
			lang   string
			wpm    float64
			raw    float64
			dur    int64
			acc    float64
			miss   int
			errs   int
			events string
		)

		// this should match the columns from the query input
		queriedCols := []any{
			&id,
			&hash,
			&start,
			&end,
			&name,
			&lang,
			&wpm,
			&raw,
			&dur,
			&acc,
			&miss,
			&errs,
			&events,
		}

		// The scan is dependent on the query that is performed
		err := rows.Scan(
			queriedCols...,
		)

		if err != nil {
			return reps, err
		}

		r := Rep{
			Id:     int(id),
			Hash:   hash,
			Start:  time.UnixMilli(start),
			End:    time.UnixMilli(end),
			Name:   name,
			Lang:   lang,
			Wpm:    wpm,
			Raw:    raw,
			Dur:    time.Duration(dur),
			Acc:    acc,
			Miss:   miss,
			Errs:   errs,
			Events: event.ParseEvents(events),
		}

		reps = append(reps, r)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return reps, err
	}

	return reps, nil
}
