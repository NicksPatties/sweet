package main

import (
	"os"
	"path"
	"strings"
)

// Returns all the pathnames of files in the given root directory.
func getAllFilePathsInDirectory(root string) ([]string, error) {
	filePaths := []string{}
	if !path.IsAbs(root) {
		wd, err := os.Getwd()
		if err != nil {
			return filePaths, err
		}
		root = strings.Replace(root, ".", wd, 1)
	}

	// read the contents of the current directory
	currDirContents, err := os.ReadDir(root)
	if err != nil {
		return filePaths, err
	}

	for _, c := range currDirContents {
		fileName := path.Join(root, c.Name())
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

func getRandomFilePathFromDirectory(root string) string {
	return "string"
}
