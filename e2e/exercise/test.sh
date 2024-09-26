#!/bin/bash

# exercise_test.sh
#
# End to end tests for the default sweet command.
# These tests verify the program exits with the intended error code
# depending on the arguments and environment. 
# Buils a copy of the executable, and removes it when tests are completed.

# Buld sweet application
sweetpath="$HOME/Documents/dev/sweet"
go build "$sweetpath"
sweet="./sweet"

./01.sh # With exercises
./02.sh # No exercises
./03.sh # No start and end flags for random exercises
./04.sh # Standard input

echo "exercise tests: pass!"
rm $sweet
