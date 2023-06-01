package main

import (
	"math/rand"
	"os"
	"path"
	"strings"
)

// Returns all the pathnames of files in the given root directory.
func getAllFilePathsInDirectory(dirPath string) ([]string, error) {
	filePaths := []string{}
	if !path.IsAbs(dirPath) {
		wd, err := os.Getwd()
		if err != nil {
			return filePaths, err
		}
		dirPath = strings.Replace(dirPath, ".", wd, 1)
	}

	// read the contents of the current directory
	currDirContents, err := os.ReadDir(dirPath)
	if err != nil {
		return filePaths, err
	}

	for _, c := range currDirContents {
		fileName := path.Join(dirPath, c.Name())
		if c.IsDir() {
			subPaths, err := getAllFilePathsInDirectory(fileName)
			if err != nil {
				return filePaths, err
			}
			for _, sp := range subPaths {
				filePaths = append(filePaths, sp)
			}
			continue
		}
		// this is a file, so I should add the path to the array
		filePaths = append(filePaths, fileName)
	}

	return filePaths, nil
}

// Gets a random file from within a directory as specified by dirPath
func getRandomFilePathFromDirectory(dirPath string) (string, error) {
	contents, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}
	filePathNames := []string{}
	if len(contents) == 0 {
		return "", nil
	}
	for _, c := range contents {
		if c.IsDir() {
			continue
		}
		filePathNames = append(filePathNames, path.Join(dirPath, c.Name()))
	}
	randI := rand.Intn(len(filePathNames))

	return filePathNames[randI], nil
}
