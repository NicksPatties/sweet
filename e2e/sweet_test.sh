#!/bin/bash

# sweet_test.sh
#
# End to end tests for the default sweet command.
# These tests verify the program exits with the intended error code
# depending on the arguments and environment. 
# Buils a copy of the executable, and removes it when tests are completed.

tempdir="temp-exercises"

echo "Starting sweet command tests."

fail () {
  echo "FAIL: $1: $2"
  teardown
  exit 1
}

teardown () {
  rm -rf $tempdir
  rm "./sweet"
}

# Buld sweet application
go build ..


# Select a random exercise from the exercises directory.
testname="sweet"
mkdir $tempdir
echo "one" >> $tempdir/one.txt
echo "two" >> $tempdir/two.txt
SWEET_EXERCISES_DIR="./$tempdir" ./sweet
got="$?"
if [ $got -ne 0 ]; then
  fail "$testname" "Exit code should be 0. Got $got."
fi
rm -rf $tempdir

testname="sweet, but no exercises"
mkdir $tempdir
SWEET_EXERCISES_DIR="./$tempdir" ./sweet
got="$?"
if [ $got -ne 1 ]; then
  fail "$testname" "Exit code should be 1. Got $got."
fi
rm -rf $tempdir

testname="sweet -s 1 -e 3, cannot assign start or end to a random exercise"
mkdir $tempdir
echo "one" >> $tempdir/one.txt
SWEET_EXERCISES_DIR="./$tempdir" ./sweet -s 1 -e 3
got="$?"
if [ $got -ne 1 ]; then
  fail "$testname" "Exit code should be 1. Got $got."
fi
rm -rf $tempdir

testname="sweet -s 1 -e 3, cannot assign start or end to a random exercise"
mkdir $tempdir
echo "one" >> $tempdir/one.txt
SWEET_EXERCISES_DIR="./$tempdir" ./sweet -s 1 -e 3
got="$?"
if [ $got -ne 1 ]; then
  fail "$testname" "Exit code should be 1. Got $got."
fi
rm -rf $tempdir

# Teardown
teardown

echo "Tests passed!"
