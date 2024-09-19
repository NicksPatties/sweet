package util

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// Grabs a string from a function that prints to os.Stdout.
// Especially useful for testing the output of print functions.
func GetStringFromStdout(function func()) string {
	prevStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	function()

	w.Close()

	os.Stdout = prevStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	return buf.String()
}

// Return the contents of a file as a string.
// Useful for comparing the output of commands that print
// a lot of text, like "about" or "help."
func GetWantFile(wantFile string, t *testing.T) string {
	f, err := os.ReadFile(wantFile)
	if err != nil {
		t.Errorf("Error opening file %s", wantFile)
	}
	return string(f)
}
