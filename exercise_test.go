package main

import (
	"testing"
	"time"
)

func TestDeleteCharacter(t *testing.T) {
	type testCase struct {
		title         string
		exercise      string
		initialTyped  string
		expectedTyped string
	}

	testCases := []testCase{
		{
			title:         "Deleting single character",
			exercise:      "exercise",
			initialTyped:  "hello",
			expectedTyped: "hell",
		},
		{
			title:         "Deleting character from string of no length",
			exercise:      "exercise",
			initialTyped:  "",
			expectedTyped: "",
		},
		{
			title:         "Deleting symbol",
			exercise:      "exercise",
			initialTyped:  "{}",
			expectedTyped: "{",
		},
		{
			title:         "Deleting spaces in middle of line",
			exercise:      "cool spot",
			initialTyped:  "cool ",
			expectedTyped: "cool",
		},
		{
			title:         "Deleting whitespace after newline deletes up to and including that newline",
			exercise:      "indent with spaces\n  booyah",
			initialTyped:  "indent with spaces\n  ",
			expectedTyped: "indent with spaces",
		},
		{
			title:         "Deleting character directly after whitespace and newline",
			exercise:      "indent with spaces\n  booyah",
			initialTyped:  "indent with spaces\n  b",
			expectedTyped: "indent with spaces\n  ",
		},
		{
			title:         "Deleting character directly after whitespace and newline, but typed newline incorrectly",
			exercise:      "indent with spaces\n  booyah",
			initialTyped:  "indent with spacesa  b",
			expectedTyped: "indent with spacesa  ",
		},
		{
			title:         "Deleting character directly after whitespace (tabs) and newline",
			exercise:      "foo\n\tbar",
			initialTyped:  "foo\n\tb",
			expectedTyped: "foo\n\t",
		},
	}

	for _, tc := range testCases {
		m := Model{
			title:         tc.title,
			exercise:      tc.exercise,
			typedExercise: tc.initialTyped,
		}
		m = m.DeleteCharacter()
		if m.typedExercise != tc.expectedTyped {
			t.Fatalf("%s Model.DeleteCharacter() got %s, wanted %s", tc.title, m.typedExercise, tc.expectedTyped)
		}
	}
}

func TestAddRuneToExercise(t *testing.T) {
	type testCase struct {
		title         string
		exercise      string
		initialTyped  string
		rn            rune
		expectedTyped string
	}
	testCases := []testCase{
		{
			"happy case",
			`a
					b
				c`,
			"a",
			Enter,
			`a
					`, // a then newline then tab
		},
		{
			"Adding percent sign",
			"Test % yeah",
			"Test ",
			Percent,
			"Test %",
		},
		{
			"Adding more characters than the length of the exercise",
			"Test",
			"Tesq",
			rune(97), // 'a'
			"Tesq",
		},
		{
			"Newline at the end of input",
			`cd mydirectory
		`,
			`cd mydirectory`,
			Enter,
			`cd mydirectory
		`,
		},
		{
			"Indentation with spaces",
			"some text\n  indented",
			"some text",
			Enter,
			"some text\n  ",
		},
		{
			"Typing another rune instead of newline",
			"some text\n\tblah",
			"some text",
			Space,
			"some text \t",
		},
	}

	for _, tc := range testCases {
		m := Model{
			tc.title,
			tc.exercise,
			tc.initialTyped,
			false,
			time.Time{},
			time.Time{},
		}
		actual := m.AddRuneToExercise(tc.rn)
		if actual.typedExercise != tc.expectedTyped {
			t.Fatalf(`Model.AddRuneToExercise() %s  %s, want to match %s`, tc.title, actual.typedExercise, tc.expectedTyped)
		}
	}
}

func TestExerciseView(t *testing.T) {
	type testCase struct {
		m        Model
		expected string
	}
	testCases := []testCase{
		{
			Model{
				"Percent sign, stop (NOVERB) bug",
				`Hello % this is the percent sign`,
				"Hello %",
				false,
				time.Time{},
				time.Time{},
			},
			`Hello % this is the percent sign`,
		},
		{
			Model{
				"Render the newline character",
				"new\nline",
				"new",
				false,
				time.Time{},
				time.Time{},
			},
			"new" + Arrow + "\nline",
		},
		{
			Model{
				"Render incorrectly typed newline correctly",
				"What's going on?\nHahah",
				"What's going on?H",
				false,
				time.Time{},
				time.Time{},
			},
			"What's going on?" + Arrow + "\nHahah",
		},
	}
	for _, tc := range testCases {
		actual := tc.m.ExerciseView()
		if actual != tc.expected {
			t.Fatalf(`%s Model.ExerciseView() "%s", want to match "%s"`, tc.m.title, actual, tc.expected)
		}
	}
}

func TestGetExerciseRuneCount(t *testing.T) {
	type testCase struct {
		title         string
		exercise      string
		expectedCount int
	}

	testCases := []testCase{
		{
			title:         "Happy case, counts all characters in a string",
			exercise:      "this should be 17",
			expectedCount: 17,
		},
		{
			title:         "Counts newlines",
			exercise:      "a\nb",
			expectedCount: 3,
		},
		{
			title:         "Counts newlines, but not the spaces afterwards",
			exercise:      "a\n  b",
			expectedCount: 3,
		},
		{
			title:         "Counts newlines, but not the tabs afterwards",
			exercise:      "a\n\tb",
			expectedCount: 3,
		},
	}

	for _, tc := range testCases {
		m := Model{
			title:         tc.title,
			exercise:      tc.exercise,
			typedExercise: "",
		}

		actualCount := m.GetExerciseRuneCount()

		if actualCount != tc.expectedCount {
			t.Fatalf("%s Model.GetExerciseRuneCount() got %d expected %d", tc.title, actualCount, tc.expectedCount)
		}
	}
}
