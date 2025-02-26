package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"runtime/debug"
	"strings"

	consts "github.com/NicksPatties/sweet/constants"
)

// i.e. go build -ldflags "-X github.com/NicksPatties/sweet/util.version=v0.1.0" .
var version string

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

func GetVersion() string {
	if version == "" {
		info, ok := debug.ReadBuildInfo()
		gitHash := ""
		if ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					gitHash = "-" + setting.Value[:6]
				}
			}
		}
		return "dev" + gitHash
	}
	return version
}

// Splits up a string of text by newlines.
// The newlines are preserved, since they'll be used
// in rendering, too.
func Lines(text string) []string {
	arr := strings.SplitAfter(text, "\n")
	if last := len(arr) - 1; arr[last] == "" {
		arr = arr[:last]
	}
	return arr
}
