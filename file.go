package main

import (
	"math/rand"
	"os"
	"path"
	"strings"
)

const (
	CONFIG_DIR_NAME = ".sweet"
)

// Copies a file defined in the srcFilePath to the destDirPath, and returns the full
// path of the added file for future use.
func addFileToDirectory(srcFilePath string, destDirPath string) (string, error) {
	// open the source file
	sData, err := os.ReadFile(srcFilePath)
	if err != nil {
		return "", nil
	}
	// create the new file
	dFileName := path.Join(destDirPath, path.Base(srcFilePath))
	df, err := os.Create(dFileName)
	if err != nil {
		return "", err
	}
	defer df.Close()

	// Write the contents of the file into destination file
	df.Write(sData)
	return dFileName, nil
}

func getDefaultConfigPath() (string, error) {
	// using the home directory for now, but later this directory
	// should be assigned in a config file somewhere
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// check if the .sweet config directory exists
	return path.Join(homeDir, CONFIG_DIR_NAME), nil
}

func createDefaultConfigDirectory() (string, error) {
	// using the home directory for now, but later this directory
	// should be assigned in a config file somewhere
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// check if the .sweet config directory exists
	sweetPath := path.Join(homeDir, CONFIG_DIR_NAME)
	_, err = os.Stat(sweetPath)
	if os.IsNotExist(err) {
		// it doesn't exist, so create it
		err = os.Mkdir(sweetPath, 0777)
		if err != nil {
			return "", err
		}

	} else if err != nil {
		return "", err
	}
	return sweetPath, nil
}

func addDefaultExercises() error {
	sweetPath, err := createDefaultConfigDirectory()
	if err != nil {
		return err
	}

	// by now, the .sweet directory should exist, so let's make subdirectories
	sweetExercisesPath := path.Join(sweetPath, "exercises")
	sweetConfigPath := path.Join(sweetPath, "config")
	os.MkdirAll(sweetExercisesPath, 0777)
	os.MkdirAll(sweetConfigPath, 0777)

	exercises := getSampleExercises()

	for _, ex := range exercises {
		sweetExerciseFilePath := path.Join(sweetExercisesPath, ex.name)
		err = os.WriteFile(sweetExerciseFilePath, []byte(ex.text), 0666)
		if err != nil {
			break
		}
	}
	return err
}

// Returns all the pathnames of files in the given root directory.
func getAllFilePathsInDirectory(dirPath string) ([]string, error) {
	filePaths := []string{}
	if !path.IsAbs(dirPath) {
		wd, err := os.Getwd()
		if err != nil {
			return filePaths, err
		}
		// note, this breaks when trying to make a change in a hidden directory (i.e. /home/user/.sweet/exercises)
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
