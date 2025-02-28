/*
version - Prints the installed version of sweet.

Usage:

	sweet version
*/
package version

import (
	"bytes"
	"fmt"

	"github.com/NicksPatties/sweet/constants"
	"github.com/NicksPatties/sweet/util"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "prints the version",
	RunE: func(cmd *cobra.Command, args []string) error {
		// just doing stupid version stuff
		return codeSnippet()
	},
}

func Red(s string) string {
	escStart := "\033[31m"
	escEnd := "\033[0m"
	return escStart + s + escEnd
}

func codeSnippet() error {
	code := `function calculateSum(a, b) {
`

	// mistakePoint := struct {
	// 	col int
	// }{
	// 	col: 3, // first "c" in "function"
	// }

	lexer := lexers.Match("source.js")
	for _, line := range util.Lines(code) {
		var buf bytes.Buffer
		err := quick.Highlight(
			&buf, line, lexer.Config().Name, "terminal", "vim",
		)
		if err != nil {
			return err
		}
		bts := buf.Bytes()
		currEscape := ""
		for i := 0; i < len(bts); i = i + 1 {
			b := bts[i]
			// If it's an escape, gather the escape character
			// and put it in the currEscape.
			if b == constants.ANSI_Escape {
				start := i
				for ; bts[i] != 'm' && i < len(bts); i = i + 1 {
				}
				currEscape = string(bts[start : i+1])
				fmt.Printf("currEscape: %#v %s\n", bts[start:i+1], currEscape)
			}
			if b >= 32 && b <= 128 {
				fmt.Printf("%#U\n", b)
			}
		}
	}

	return nil
}
