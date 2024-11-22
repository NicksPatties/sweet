package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path"
)

// Converts a string to an md5 hash. Used to
// convert the contents of an exercise into a string
// to verify if their contents are the same.
func MD5Hash(contents string) string {
	bytes := []byte(contents)
	hash := md5.Sum(bytes)
	return hex.EncodeToString(hash[:])
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

func IsValidURL(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}

	// Check if scheme and host are present
	return u.Scheme != "" && u.Host != ""
}
