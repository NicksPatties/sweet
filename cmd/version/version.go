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
	lg "github.com/charmbracelet/lipgloss"
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

  // This is a comment
  const result = a + b;
  if (result > 10) {
    console.log("Result is greater than 10");
    return true;
  } else {
    return false;
  }
}

function add(a, b) {
  return a + b
}
`

	missRow := 2
	missCol := 3
	render := lg.NewStyle().Background(lg.Color("1")).Foreground(lg.Color("15")).Render
	lexer := lexers.Match("source.js")
	for row, line := range util.Lines(code) {
		var buf bytes.Buffer
		err := quick.Highlight(
			&buf, line, lexer.Config().Name, "terminal", "vim",
		)
		if err != nil {
			return err
		}
		// render the individual line
		bts := buf.Bytes()
		currEscape := ""
		col := 0
		for i := 0; i < len(bts); i = i + 1 {
			b := bts[i]
			// If it's an escape, gather the escape character
			// and put it in the currEscape.
			if b == constants.ANSI_Escape {
				start := i
				for ; bts[i] != 'm' && i < len(bts); i = i + 1 {
				}
				currEscape = string(bts[start : i+1])
				fmt.Printf("%s", currEscape)
			}
			if b >= 32 && b <= 128 ||
				b == constants.Enter ||
				b == constants.Tab {
				char := string(b)
				if col == missCol && row == missRow {
					char = render(char)
				}
				fmt.Printf("%s", char)
				col = col + 1
			}
		}
	}

	return nil
}
