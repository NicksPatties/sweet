#!/bin/bash

# name="Select a random exercise"
# cmd="./sweet"

# $1 The name of the test
# $2 The error message
fail () {
  echo "FAIL: $1: $2"
  teardown
  exit 1
}

setup () {
  echo "setup"
}

run_test () {
  echo "run test"
}

teardown () {
  echo "teardown"
}

setup
run_test
teardown
