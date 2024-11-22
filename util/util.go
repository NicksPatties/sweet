package util

import (
	"fmt"
	"hash/crc32"
	"io"
	"net/url"
	"os"
	"path"
)

// Converts a file from the filePath into a hashed value
func HashFile(filePath string) (uint32, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Create a new CRC-32 hash object using the IEEE polynomial
	hasher := crc32.NewIEEE()

	// Copy the file contents to the hasher
	if _, err := io.Copy(hasher, file); err != nil {
		return 0, err
	}

	// Get the checksum
	checksum := hasher.Sum32()

	return checksum, nil
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
