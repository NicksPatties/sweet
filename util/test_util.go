package util

import (
	"bytes"
	"fmt"
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

func Red(s string) string {
	escStart := "\033[31m"
	escEnd := "\033[0m"
	return escStart + s + escEnd
}

func Reds(s string) string {
	finished := ""
	for _, r := range s {
		finished += Red(string(r))
	}
	return finished
}

// Renders either the base 16 byte code of a byte,
// or it's visual representation. Useful for debugging
// rendering errors.
func RenderBytes(str string) (s string) {
	bytes := []byte(str)
	for i, b := range bytes {
		c := fmt.Sprintf("\\x%x", str[i])
		if b >= 32 && b <= 128 {
			c = fmt.Sprintf("%s", string(str[i]))
		}
		s += c
	}
	return
}
