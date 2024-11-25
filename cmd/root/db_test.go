package root

import (
	"os"
	"path"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestGetStatsDb(t *testing.T) {
	t.Run("create database in default location", func(t *testing.T) {
		// Move current config to temporary location.
		// Move it back when test is done.
		tempDir := t.TempDir()
		userConfigDir, _ := os.UserConfigDir()
		prevSweetConfigLocation := path.Join(userConfigDir, "sweet")
		os.Rename(prevSweetConfigLocation, tempDir)
		defer os.Rename(tempDir, prevSweetConfigLocation)

		os.Unsetenv("SWEET_DB_LOCATION")

		db, err := SweetDb()
		if err != nil {
			t.Fatalf("Should create database without error: %v", err)
		}
		defer db.Close()

		expectedPath := path.Join(userConfigDir, "sweet", "sweet.db")

		_, err = os.Stat(expectedPath)
		if err != nil {
			t.Errorf("Database file should exist at expected location: %v", err)
		}
	})

	t.Run("create database in specified location", func(t *testing.T) {
		tempDir := t.TempDir()
		os.Setenv("SWEET_DB_LOCATION", tempDir)
		defer os.Unsetenv("SWEET_DB_LOCATION")

		db, err := SweetDb()
		if err != nil {
			t.Fatalf("Should create database at custom location: %v", err)
		}
		defer db.Close()

		expectedPath := path.Join(tempDir, "sweet.db")

		// Verify file exists
		_, err = os.Stat(expectedPath)
		if err != nil {
			t.Errorf("Database file should exist at custom location: %v", err)
		}
	})

	t.Run("validate db schema is correct", func(t *testing.T) {
		db, err := SweetDb()
		if err != nil {
			t.Fatalf("Should create database: %v", err)
		}
		defer db.Close()

		// Check if 'reps' table exists
		rows, err := db.Query("PRAGMA table_info(reps);")
		if err != nil {
			t.Fatalf("Should be able to query table info: %v", err)
		}
		defer rows.Close()

		// Expected columns
		expectedColumns := []string{
			"id", "hash", "start", "end", "name", "lang",
			"wpm", "raw", "dur", "acc", "miss", "errs", "events",
		}

		// Collect actual column names
		var columns []string
		for rows.Next() {
			var (
				cid       int
				name      string
				typ       string
				notNull   int
				dfltValue interface{}
				pk        int
			)
			err := rows.Scan(&cid, &name, &typ, &notNull, &dfltValue, &pk)
			if err != nil {
				t.Fatalf("Error scanning row: %v", err)
			}
			columns = append(columns, name)
		}

		// Verify all expected columns exist
		for _, expected := range expectedColumns {
			found := false
			for _, actual := range columns {
				if actual == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Column %s should exist in reps table", expected)
			}
		}
	})

	t.Run("should error if the location doesn't exist", func(t *testing.T) {
		os.Setenv("SWEET_DB_LOCATION", "/absolutely/impossible/path/that/cannot/exist")
		defer os.Unsetenv("SWEET_DB_LOCATION")

		_, err := SweetDb()
		if err == nil {
			t.Error("Should return an error for invalid path")
		}
	})
}

func TestEventsToColumn(t *testing.T) {
	want := "2024-10-07 16:29:26.916\t0\tc\tc\n" +
		"2024-10-07 16:29:27.004\t1\to\to\n" +
		"2024-10-07 16:29:27.095\t2\tn\tn\n" +
		"2024-10-07 16:29:27.279\t3\ts\ts\n" +
		"2024-10-07 16:29:27.416\t4\to\to\n" +
		"2024-10-07 16:29:27.667\t5\tl\tl\n" +
		"2024-10-07 16:29:27.784\t6\te\te\n" +
		"2024-10-07 16:29:31.538\t7\tenter\tenter"
	events := parseEvents(want)
	got := eventsToColumn(events)

	if got != want {
		t.Error("got != want")
	}
}

func TestInsertRep(t *testing.T) {
	t.Run("inserting a rep into a blank database", func(t *testing.T) {
		tempDir := t.TempDir()
		os.Setenv("SWEET_DB_LOCATION", tempDir)
		defer os.Unsetenv("SWEET_DB_LOCATION")

		// create a test database
		db, err := SweetDb()
		if err != nil {
			db.Close()
			t.Fatalf("failed to initialize sweet db: %v", err)
		}

		// initialize the rep
		start := time.Now()
		end := time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), start.Minute()+int(time.Minute), start.Second(), start.Nanosecond(), time.UTC)

		rep := Rep{
			hash:  "abcef123456",
			start: start,
			end:   end,
			name:  "exercise.go",
			lang:  "go",
			wpm:   60.333,
			raw:   65.5,
			dur:   end.Sub(start),
			acc:   98.98,
			miss:  2,
			errs:  1,
			events: parseEvents(
				"2024-10-07 13:46:47.679\t0\th\th\n" +
					"2024-10-07 13:46:56.521\t3\tenter\tenter",
			),
		}

		id, err := InsertRep(db, rep)
		if err != nil {
			t.Errorf("Inserting rep failed: %v", err)
		}

		if id == 0 {
			t.Errorf("Retrieving inserted id failed: %v", err)
		}

	})

	t.Run("inserting an invalid rep fails", func(t *testing.T) {
		tempDir := t.TempDir()
		os.Setenv("SWEET_DB_LOCATION", tempDir)
		defer os.Unsetenv("SWEET_DB_LOCATION")

		// create a test database
		db, err := SweetDb()
		if err != nil {
			db.Close()
			t.Fatalf("failed to initialize sweet db: %v", err)
		}

		// initialize the rep
		start := time.Now()
		end := time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), start.Minute()+int(time.Minute), start.Second(), start.Nanosecond(), time.UTC)

		rep := Rep{
			hash:  "abcef123456",
			start: start,
			end:   end,
			name:  "exercise.go",
			lang:  "go",
			wpm:   -60.333, // should fail
			raw:   65.5,
			dur:   end.Sub(start),
			acc:   98.98,
			miss:  2,
			errs:  1,
			events: parseEvents(
				"2024-10-07 13:46:47.679\t0\th\th\n" +
					"2024-10-07 13:46:56.521\t3\tenter\tenter",
			),
		}

		_, err = InsertRep(db, rep)
		if err == nil {
			t.Error("should have error")
		}
	})

}
