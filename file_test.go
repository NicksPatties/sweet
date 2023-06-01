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

	find := func(pathName string, expected []string) bool {
		for _, e := range expected {
			if pathName == e {
				return true
			}
		}
		return false
	}

	actual, err := getAllFilePathsInDirectory(root)

	if err != nil {
		t.Fatalf("%s: wanted []string, got error instead: %s", gafpid, err)
	}

	for _, a := range actual {
		if !find(a, expected) {
			t.Fatalf("%s: couldn't find %s in results", gafpid, a)
		}
	}
}
