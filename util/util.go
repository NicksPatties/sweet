package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"strings"

	consts "github.com/NicksPatties/sweet/constants"
)

// Converts a string to an md5 hash. Used to
// convert the contents of an exercise into a string
// to verify if their contents are the same.
//
// see: https://stackoverflow.com/a/25286918
func MD5Hash(contents string) string {
	bytes := []byte(contents)
	hash := md5.Sum(bytes)
	return hex.EncodeToString(hash[:])
}

// Gets the language of the provided filename.
// Unlike `path.Ext`, the language doesn't include the
// leading dot.
func Lang(filename string) (lang string) {
	lang = ""
	split := strings.Split(filename, ".")
	if len(split) > 1 {
		lang = split[len(split)-1]
	}
	return
}

// Gets the path for sweet's configuration directory.
//
// See `os.UserConfigDir` for the default configuration
// location depending on the current operating system.
func SweetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %v", err)
	}
	return path.Join(configDir, "sweet"), nil
}

// Filters a list of file names by the given language extension.
func FilterFileNames(fileNames []string, language string) (found []string) {
	for _, f := range fileNames {
		ext := path.Ext(f)
		// Ignore files that don't have an extension.
		if len(ext) == 0 {
			continue
		}
		if ext[1:] == language {
			found = append(found, f)
		}
	}
	return found
}

func IsWhitespace(rn rune) bool {
	return rn == consts.Tab || rn == consts.Space
}

// Splits up a string of text by newlines.
// The newlines are preserved, since they'll be used
// in rendering, too.
func Lines(text string) []string {
	arr := strings.SplitAfter(text, "\n")
	if arr[len(arr)-1] == "" {
		arr = arr[:len(arr)-1]
	}
	return arr
}

// Returns an array of strings that map the typed characters
// to the exercise characters. If no characters have been typed
// on a current line, the typedLine will be nil.
func TypedLines(lines []string, typed string) []string {
	typedLines := []string{}
	i := 0
	for _, line := range lines {
		str := ""
		for range line {
			if i >= len(typed) {
				continue
			}
			str = str + string(typed[i])
			i = i + 1
		}
		if str != "" {
			typedLines = append(typedLines, str)
		}
	}
	return typedLines
}

func CurrentLine(lines []string, typed string) int {
	typedLen := len(typed)
	for i := range lines {
		for range lines[i] {
			if typedLen == 0 {
				return i
			}
			typedLen = typedLen - 1
		}
	}
	return 0
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

func RemoveLastNewline(str string) string {
	n := '\n'
	i := len(str) - 1
	for ; i >= 0 && rune(str[i]) != n; i = i - 1 {
	}

	if i < 0 {
		return str
	}

	return fmt.Sprintf("%s%s", str[:i], str[i+1:])
}
