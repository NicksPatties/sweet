package main

import (
	"os"
	"path"
	"testing"
)

const (
	gafpid = "getAllFilePathsInDirectory"
	grfpfd = "getRandomFilePathFromDirectory"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func find(pathName string, expected []string) bool {
	for _, e := range expected {
		if pathName == e {
			return true
		}
	}
	return false
}

func TestGetAllFilePathsInDirectory_EmptyDirectory(t *testing.T) {
	// create a temp directory
	root, err := os.MkdirTemp(".", "tmp")
	check(err)
	defer os.Remove(root)

	expected := []string{}

	// read the file
	actual, err := getAllFilePathsInDirectory(root)
	check(err)

	if err != nil {
		t.Fatalf("%s: wanted %s, got error instead: %s", gafpid, expected, err)
	}

	if len(actual) != 0 && err != nil {
		t.Fatalf("%s: wanted empty slice, got %s", gafpid, actual)
	}
}

func TestGetAllFilePathsInDirectory_MultipleEmptyDirectories(t *testing.T) {
	root, err := os.MkdirTemp(".", "root")
	check(err)
	defer os.Remove(root)

	child, err := os.MkdirTemp(root, "child")
	check(err)
	defer os.Remove(child)

	expected := []string{}

	// read the file
	actual, err := getAllFilePathsInDirectory(root)
	check(err)

	if err != nil {
		t.Fatalf("%s: wanted %s, got error instead: %s", gafpid, expected, err)
	}

	if len(actual) != 0 && err != nil {
		t.Fatalf("%s: wanted empty slice, got %s", gafpid, actual)
	}
}

func TestGetAllFilePathsInDirectory_FilesInsideRootDirectory(t *testing.T) {
	wd, err := os.Getwd()
	check(err)

	root, err := os.MkdirTemp(wd, "root")
	check(err)
	defer os.Remove(root)

	file, err := os.CreateTemp(root, "tmp")
	check(err)

	filePath := file.Name()
	defer os.Remove(filePath)

	expected := []string{
		filePath,
	}

	actual, err := getAllFilePathsInDirectory(root)

	if err != nil {
		t.Fatalf("%s: wanted %s, got error instead: %s", gafpid, expected[0], actual)
	}

	if expected[0] != actual[0] {
		t.Fatalf("%s: wanted %s, got %s", gafpid, expected[0], actual[0])
	}
}

func TestGetAllFilePathsInDirectory_FilesInsideRootAndChildDirectories(t *testing.T) {
	wd, err := os.Getwd()
	check(err)

	// root directory
	root, err := os.MkdirTemp(wd, "root")
	check(err)
	defer os.Remove(root)

	// tmp file in root directory
	file, err := os.CreateTemp(root, "tmp")
	check(err)
	filePath := file.Name()
	defer os.Remove(filePath)

	// child directory
	child, err := os.MkdirTemp(root, "child")
	check(err)
	defer os.Remove(child)

	// tmp file in child directory
	childFile, err := os.CreateTemp(child, "childtmp")
	check(err)
	childFilePath := childFile.Name()
	defer os.Remove(childFilePath)

	expected := []string{
		filePath,
		childFilePath,
	}

	actual, err := getAllFilePathsInDirectory(root)

	if err != nil {
		t.Fatalf("%s: wanted []string, got error instead: %s", gafpid, err)
	}

	if len(expected) != len(actual) {
		t.Fatalf("%s: wanted %d total files, but only found %d", gafpid, len(expected), len(actual))
	}

	for _, ex := range expected {
		if !find(ex, actual) {
			t.Fatalf("%s: couldn't find %s in results", gafpid, ex)
		}
	}
}

func TestGetAllFilePathsInDirectory_FilesInsideMultipleChildDirectories(t *testing.T) {
	wd, err := os.Getwd()
	check(err)

	// root directory
	root, err := os.MkdirTemp(wd, "root")
	check(err)
	defer os.Remove(root)

	// child1 directory
	child1, err := os.MkdirTemp(root, "child1")
	check(err)
	defer os.Remove(child1)

	// tmp file in child directory
	child1File, err := os.CreateTemp(child1, "childtmp")
	check(err)
	child1FilePath := child1File.Name()
	defer os.Remove(child1FilePath)

	// child2 directory
	child2, err := os.MkdirTemp(root, "child2")
	check(err)
	defer os.Remove(child2)

	// tmp file in child directory
	child2File, err := os.CreateTemp(child2, "childtmp")
	check(err)
	child2FilePath := child2File.Name()
	defer os.Remove(child2FilePath)

	expected := []string{
		child1FilePath,
		child2FilePath,
	}

	actual, err := getAllFilePathsInDirectory(root)

	if err != nil {
		t.Fatalf("%s: wanted []string, got error instead: %s", gafpid, err)
	}

	if len(expected) != len(actual) {
		t.Fatalf("%s: wanted %d total files, but only found %d", gafpid, len(expected), len(actual))
	}

	for _, ex := range expected {
		if !find(ex, actual) {
			t.Fatalf("%s: couldn't find %s in results", gafpid, ex)
		}
	}
}

func TestGetRandomFilePathFromDirectory_EmptyDirectoryReturnsEmptyString(t *testing.T) {

	expected := ""

	wd, err := os.Getwd()
	check(err)

	root, err := os.MkdirTemp(wd, "root")
	check(err)
	defer os.Remove(root)
	actual, err := getRandomFilePathFromDirectory(root)

	if err != nil {
		t.Fatalf("%s: expected %s, but found error %s", grfpfd, expected, err)
	}

	// need to search inside of the files that were created in the root directory
	if expected != actual {
		t.Fatalf("%s: expected %s, got %s", grfpfd, expected, actual)
	}
}

func TestGetRandomFilePathFromDirectory_ReturnsAFilePathInTheDirectory(t *testing.T) {
	wd, err := os.Getwd()
	check(err)

	root, err := os.MkdirTemp(wd, "root")
	check(err)

	// create files
	fileNames := []string{
		"a",
		"b",
		"c",
	}
	filePathNames := []string{}
	files := []*os.File{}

	for _, fn := range fileNames {
		f, err := os.CreateTemp(root, fn)
		check(err)
		filePathNames = append(filePathNames, f.Name())
		files = append(files, f)
	}

	removeFiles := func(files []*os.File) error {
		for _, f := range files {
			err := os.Remove(f.Name())
			if err != nil {
				return err
			}
		}
		return nil
	}

	// run the test
	actual, err := getRandomFilePathFromDirectory(root)

	if err != nil {
		t.Errorf("%s: supposed to return pathname, but got error instead: %s", grfpfd, err)
	}

	if !find(actual, filePathNames) {
		t.Errorf("%s: couldn't find %s in temporary directory", grfpfd, actual)
	}

	t.Cleanup(func() {
		removeFiles(files)
		os.Remove(root)
	})
}

func TestAddFileToDirectory(t *testing.T) {
	wd, err := os.Getwd()
	check(err)

	// mock directory one
	m1, err := os.MkdirTemp(wd, "m1")
	check(err)

	// tmp file inside mock directory one
	tmp, err := os.CreateTemp(m1, "tmp")
	tmpFileName := path.Base(tmp.Name())
	check(err)

	// mock directory two
	m2, err := os.MkdirTemp(wd, "m2")
	check(err)

	expectedContents := "Hello"
	expectedPath := path.Join(m2, tmpFileName)

	tmp.WriteString(expectedContents)

	srcFilePath := path.Join(m1, tmpFileName)
	destDirPath := m2

	actualPath, err := addFileToDirectory(srcFilePath, destDirPath)
	check(err)

	actualContents, err := os.ReadFile(path.Join(m2, tmpFileName))
	check(err)

	if err != nil {
		t.Fatalf("addFileToDirectory: %s\n", err)
	}
	if expectedContents != string(actualContents) {
		t.Fatalf("addFileToDirectory: wanted %s, got %s", expectedContents, string(actualContents))
	}

	if expectedPath != actualPath {
		t.Fatalf("addFileToDirectory: wanted %s path, got %s path", expectedPath, actualPath)
	}

	// remove tmp files and folders
	os.Remove(tmp.Name())
	os.Remove(path.Join(m2, tmpFileName))
	os.Remove(m1)
	os.Remove(m2)
}
