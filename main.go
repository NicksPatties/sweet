package main

import (
	"fmt"
	"github.com/NicksPatties/sweet/cmd/root"
)

func main() {
	if err := root.Cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
