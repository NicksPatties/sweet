package main

import (
	"os"
	"testing"
)

const (
	gafpid = "getAllFilePathsInDirectory"
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
