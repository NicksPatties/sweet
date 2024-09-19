package util

import (
	"bytes"
	"io"
	"os"
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
