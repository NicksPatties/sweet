#!/bin/bash

name="Select a random exercise"
dir="test-directory"
file="one.txt"
text="one"
cmd="./sweet"

# $1 The name of the test
# $2 The error message
fail () {
  echo "FAIL: $1: $2"
  teardown
  exit 1
}

setup () {
  mkdir $dir
  echo "$text" >> $dir/$file
}

run_test () {
  SWEET_EXERCISES_DIR=$dir $cmd
  local got=$?
  local want=0
  if [ $got -ne $want ]; then
    fail "$name" "Wanted $want, got $got." 
  fi
}

teardown () {
  rm -rf $dir
}

setup
run_test
teardown
