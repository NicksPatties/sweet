package root

import (
	"os"
	"path"
	"testing"

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
		// Set an environment variable with an impossible path
		os.Setenv("SWEET_DB_LOCATION", "/absolutely/impossible/path/that/cannot/exist")
		defer os.Unsetenv("SWEET_DB_LOCATION")

		_, err := SweetDb()
		if err == nil {
			t.Error("Should return an error for invalid path")
		}
	})
}
