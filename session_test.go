package main

import (
	"fmt"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func oneOfTheseCmdsIsNil(a tea.Cmd, b tea.Cmd) bool {
	return (a == nil && b != nil) || (a != nil && b == nil)
}

func modelsAreTheSame(a Model, b Model) bool {
	return a.typedExercise == b.typedExercise
}

func (m Model) String() string {
	c := "\n"
	c += fmt.Sprintf("title: %s\n", m.title)
	c += fmt.Sprintf("exercise: %s\n", m.exercise)
	c += fmt.Sprintf("typedExercise: %s\n", m.typedExercise)
	c += fmt.Sprintf("quitEarly: %t\n", m.quitEarly)
	c += fmt.Sprintf("startTime: %s\n", m.startTime.String())
	c += fmt.Sprintf("endTime: %s\n", m.endTime)
	return c
}

func TestUpdate(t *testing.T) {
	type input struct {
		model Model
		msg   tea.Msg
	}
	type output struct {
		model Model
		cmd   tea.Cmd
	}
	type check func(output, output) bool
	type errMsg func(output, output) string
	type testCase struct {
		name     string
		input    input
		expected output
		check    check
		errMsg   errMsg
	}

	fakeStartTime := time.Now()
	testCases := []testCase{
		{
			"Handle Enter key message",
			input{
				Model{
					"",
					"a test\nwith a newline",
					"a test",
					false,
					fakeStartTime,
					time.Time{},
				},
				tea.KeyMsg{
					Type:  tea.KeyEnter,
					Runes: []rune{},
					Alt:   false,
				},
			},
			output{
				Model{
					"",
					"a test\nwith a newline",
					"a test\n",
					false,
					fakeStartTime,
					time.Time{},
				},
				nil,
			},
			func(expected output, actual output) bool {
				return expected.model.String() == actual.model.String()
			},
			func(expected output, actual output) string {
				return fmt.Sprintf("%s \ndoes not match %s", expected.model.String(), actual.model.String())
			},
		},
		{
			"Start the timer when the user begins typing",
			input{
				Model{
					"timer test",
					"a test",
					"",
					false,
					time.Time{},
					time.Time{},
				},
				tea.KeyMsg{
					Type:  tea.KeyRunes,
					Runes: []rune{'a'},
					Alt:   false,
				},
			},
			output{
				Model{
					"timer test",
					"a test",
					"a",
					false,
					fakeStartTime,
					time.Time{},
				},
				nil,
			},
			func(expected output, actual output) bool {
				return !expected.model.startTime.IsZero()
			},
			func(expected output, actual output) string {
				return "startTime is supposed to be non-zero"
			},
		},
		{
			"Save the end time when the user finished the exercise",
			input{
				Model{
					"",
					"exercise",
					"exercis",
					false,
					fakeStartTime,
					time.Time{},
				},
				tea.KeyMsg{
					Type:  tea.KeyRunes,
					Runes: []rune{'e'},
					Alt:   false,
				},
			},
			output{
				Model{
					"",
					"exercise",
					"exercise",
					false,
					fakeStartTime,
					time.Now(),
				},
				tea.Quit,
			},
			func(expected output, actual output) bool {
				actualEnd := actual.model.endTime
				expectedStart := expected.model.startTime
				return !actual.model.endTime.IsZero() && actualEnd.After(expectedStart)
			},
			func(expected output, actual output) string {
				return fmt.Sprintf("either actual end time is zero, or not after expected start time")
			},
		},
		{
			"End the exercise early when pressing Ctrl+c",
			input{
				Model{
					"",
					"a test",
					"",
					false,
					time.Time{},
					time.Time{},
				},
				tea.KeyMsg{
					Type:  tea.KeyCtrlC,
					Runes: []rune{},
					Alt:   false,
				},
			},
			output{
				Model{
					"",
					"a test\nwith a newline",
					"",
					true,
					fakeStartTime,
					time.Time{},
				},
				tea.Quit,
			},
			func(expected output, actual output) bool {
				// want to compare each of the commands to each other to verify
				// if they're the same, but not sure how to do that...
				return expected.model.quitEarly == actual.model.quitEarly &&
					actual.cmd != nil
			},
			func(expected output, actual output) string {
				return fmt.Sprintf(
					"wanted quitEarly to be %t, but got %t, or the actual tea.Cmd is nil",
					expected.model.quitEarly,
					actual.model.quitEarly,
				)
			},
		},
	}

	for _, tc := range testCases {
		m, actualCmd := tc.input.model.Update(tc.input.msg)
		actualModel := m.(Model)
		if oneOfTheseCmdsIsNil(tc.expected.cmd, actualCmd) {
			t.Errorf("Update: %s: tea commands don't match, one is nil!", tc.name)
		}
		actual := output{actualModel, actualCmd}
		if !tc.check(tc.expected, output{actualModel, actualCmd}) {
			errMsg := tc.errMsg(tc.expected, actual)
			t.Errorf("Update: %s: %s", tc.name, errMsg)
		}
	}
}
