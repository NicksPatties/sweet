package db

import (
	"database/sql"
	"fmt"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

// Gets a pointer to the stats database. If the
// database doesn't exist already, then create it, and place
// it in the location specified by the SWEET_DB_LOCATION environment
// variable.
//
// The default location of the database file `sweet.db` is
// `~/.config/sweet/`
func GetStatsDb() (*sql.DB, error) {
	// get the path for the database
	var dbPath string
	if envDir := os.Getenv("SWEET_DB_LOCATION"); envDir != "" {
		dbPath = envDir
	} else {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("failed to find user config directory: %v", err)
		}
		dbPath = path.Join(configDir, "sweet")
	}

	// create the sweet config directory
	if err := os.MkdirAll(dbPath, 0775); err != nil {
		return nil, fmt.Errorf("failed to find or create sweet config directory: %v", err)
	}
	// The directory should be available by the time
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", path.Join(dbPath, "sweet.db"))
	defer db.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
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
		    wpm real not null,
		    -- raw words per minute
		    raw real not null,
		    -- duration: duration of rep in **nanoseconds**
		    dur integer not null,
		    -- accuracy: float between [0, 100]
		    acc real not null,
		    -- mistakes: must be gte 0
		    miss integer not null,
		    -- uncorrected errors: must be gte 0
		    errs integer not null,
		    -- array of events, events are separated by '\n'
		    events text not null
		);
    `)

	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return db, nil
}
