package commands

import "os"

// Returns the name of the executable, or the default sweet command name
func GetCommandSweet() string {
	return os.Args[0]
}

// Command names
const (
	CommandSweet   = "sweet"
	CommandHelp    = "help"
	CommandList    = "list"
	CommandVersion = "version"
)
