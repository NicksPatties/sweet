package root

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"time"

	. "github.com/NicksPatties/sweet/util"
	_ "github.com/mattn/go-sqlite3"
)

// A single repetition of an exercise.
//
// This is the structure of a row in the sweet db module.
// See `db/db.go` for more details.
type Rep struct {
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
	Events []Event // one string is one event
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
		sweetDir, err := SweetConfigDir()
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
	db, err := sql.Open("sqlite3", path.Join(dbPath, "sweet.db"))
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
	_, err = db.Exec(`
CREATE TABLE if not exists reps(
  id integer primary key autoincrement not null,
  -- md5 hash of the exercise file's contents
  hash string NOT NULL,
  -- start time in unix milliseconds
  start integer not null,
  -- end time in unix milliseconds
  end integer not null,
  -- name of the exercise file, includes extension if present
  name text not null,
  -- language: extension of the exercise file, or "" if there is none.
  lang text,
  -- words per minute
  wpm real not null check(wpm >= 0.0),
  -- raw words per minute
  raw real not null check(raw >= 0.0),
  -- duration: duration of rep in **nanoseconds**
  dur integer not null check(dur >= 0),
  -- accuracy: float between [0, 100]
  acc real not null check(acc >= 0.0),
  -- mistakes: must be gte 0
  miss integer not null check(miss >= 0),
  -- uncorrected errors: must be gte 0
  errs integer not null check(errs >= 0),
  -- array of events, events are separated by '\n'
  events text not null
);`)

	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return db, nil
}

func eventsStringToColumn(events []Event) (s string) {
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
	query := `insert into reps (
	    hash, start, end, name, lang, wpm,
	    raw, dur, acc, miss, errs, events
	   ) values (
	   	?, ?, ?, ?, ?, ?,
	   	?, ?, ?, ?, ?, ?
	   );`

	result, err := db.Exec(query,
		hash, start, end, name, lang, wpm,
		raw, dur, acc, miss, errs, events,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
