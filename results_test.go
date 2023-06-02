package main

import (
	"fmt"
	"testing"
	"time"
)

func TestRequiredRunes(t *testing.T) {
	type testCase struct {
		name     string
		exercise string
		expected []rune
	}

	testCases := []testCase{
		{
			"String of no length",
			"",
			[]rune{},
		},
		{
			"String that has no whitespace after newline",
			"abc",
			[]rune{'a', 'b', 'c'},
		},
		{
			"String that does have whitespace after newline",
			"abc\n\tdef",
			[]rune{'a', 'b', 'c', '\n', 'd', 'e', 'f'},
		},
	}

	for _, tc := range testCases {
		actual := requiredRunes(tc.exercise)
		if len(tc.expected) != len(actual) {
			t.Errorf("requiredCharacters: %s: lengths of %v and %v don't match", tc.name, tc.expected, actual)
		}

		for i, aRune := range actual {
			exRune := tc.expected[i]
			if exRune != aRune {
				t.Errorf("requiredCharacters: %s: wanted %v, got %v", tc.name, exRune, aRune)
			}
		}
	}
}

func TestNumMistakes(t *testing.T) {
	type testCase struct {
		name     string
		typed    string
		exercise string
		expected int
	}

	testCases := []testCase{
		{
			"blank typed",
			"exercise",
			"",
			0,
		},
		{
			"blank exercise",
			"",
			"typed",
			0,
		},
		{
			"partly completed exercise, no mistakes",
			"exercise",
			"exer",
			0,
		},
		{
			"partly completed exercise, some mistakes",
			"exercise",
			"oxor",
			2,
		},
		{
			"over typed, no mistakes",
			"exercise",
			"exercise that is so cool",
			0,
		},
		{
			"over typed, some mistakes",
			"exercise",
			"oxercise wowowowo",
			1,
		},
	}

	for _, tc := range testCases {
		actual := NumMistakes(tc.typed, tc.exercise)
		if tc.expected != actual {
			t.Errorf("NumMistakes: %s: wanted %d, got %d", tc.name, tc.expected, actual)
		}
	}
}

func TestAccuracy(t *testing.T) {
	format := func(x float32) string {
		return fmt.Sprintf("%.2f", x)
	}
	type testCase struct {
		name     string
		typed    string
		exercise string
		expected string
	}
	testCases := []testCase{
		{
			"blank typed",
			"",
			"exercise",
			format(float32(0)),
		},
		{
			"blank exercise",
			"typed",
			"",
			format(float32(0)),
		},
		{
			"partly completed exercise, no mistakes",
			"exercise",
			"exer",
			format(float32(100)),
		},
		{
			"partly completed exercise, some mistakes",
			"exercise",
			"oxer",
			format(float32(75)),
		},
		{
			"over typed, no mistakes",
			"exercise",
			"exercise that is so cool",
			format(float32(100)),
		},
		{
			"over typed, some mistakes",
			"1234",
			"2234 wowowowo",
			format(float32(75)),
		},
	}

	for _, tc := range testCases {
		actual := Accuracy(tc.typed, tc.exercise)
		actualFormatted := format(actual)
		if tc.expected != actualFormatted {
			t.Errorf("TestAccuracy: %s: wanted %s, got %s", tc.name, tc.expected, actualFormatted)
		}
	}
}

func TestWPM(t *testing.T) {
	format := func(x float64) string {
		return fmt.Sprintf("%.f", x)
	}

	startTime := time.Now()
	endTime := startTime.Add(1 * time.Minute) // one minute duration
	exercise := "12345"
	typedExercise := "12345"
	expected := "1"
	actual := WPM(startTime, endTime, typedExercise, exercise, WORD_SIZE)
	formattedActual := format(actual)
	if expected != formattedActual {
		t.Errorf("WPM: wanted %s, got %s", expected, formattedActual)
	}

}

func TestCPM(t *testing.T) {
	format := func(x float64) string {
		return fmt.Sprintf("%.f", x)
	}

	startTime := time.Now()
	endTime := startTime.Add(1 * time.Minute) // one minute duration
	exercise := "12345"
	typedExercise := "12345"
	expected := "5"
	actual := CPM(startTime, endTime, typedExercise, exercise)
	formattedActual := format(actual)
	if expected != formattedActual {
		t.Errorf("WPM: wanted %s, got %s", expected, formattedActual)
	}
}
