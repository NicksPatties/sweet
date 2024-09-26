#!/bin/bash

name="No start and end on random exercise"
dir="test-directory"
file="one.txt"
text=$(printf "one\ntwo")
cmd="./sweet -s 1 -e 2"

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
  local want=1
  if [ $got -ne $want ]; then
    fail "$name" "Wanted $want, got $got." 
  fi
}

teardown () {
  rm -rf $dir
}

# Runs the first sweet test
01 () {
  setup
  run_test
  teardown
}
