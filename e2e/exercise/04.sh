#!/bin/bash

name="Run the exercise with standard input"

# $1 The name of the test
# $2 The error message
fail () {
  echo "FAIL: $1: $2"
  teardown
  exit 1
}

setup () {
  true
}

# Note, this exercise should run the game with "stdin" text.
# This may pass erroneously if it runs the wrong exercise.
run_test () {
  echo stdin | ./sweet -
  local got=$?
  local want=0
  if [ $got -ne $want ]; then
    fail "$name" "Wanted $want, got $got." 
  fi
}

teardown () {
  true
}

setup
run_test
teardown
