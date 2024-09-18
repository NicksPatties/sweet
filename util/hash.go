package util

import (
	"hash/crc32"
	"io"
	"os"
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
