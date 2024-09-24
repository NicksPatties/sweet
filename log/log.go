package log

import (
	"fmt"
	"os"
	"path"
)

var l logger = NewLogger(path.Base(os.Args[0]))

type logger struct {
	exeName string
}

func NewLogger(exeName string) logger {
	return logger{
		exeName: exeName,
	}
}

func PrintErr(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, "%s ERROR: "+msg+"\n", l.exeName, args)
}

func PrintWarn(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, "%s WARN: "+msg+"\n", l.exeName, args)
}

func PrintInfo(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, "%s INFO: "+msg+"\n", l.exeName, args)
}
