#!/bin/bash

name="No exercises"
dir="test-directory"
cmd="./sweet"

# $1 The name of the test
# $2 The error message
fail() {
  echo "FAIL: $1: $2"
  teardown
  exit 1
}

setup() {
  mkdir $dir
}

run_test() {
  SWEET_EXERCISES_DIR=$dir $cmd
  local got=$?
  local want=1
  if [ $got -ne $want ]; then
    02_fail "$name" "Wanted $want, got $got."
  fi
}

teardown() {
  rm -rf $dir
}

setup
run_test
teardown
