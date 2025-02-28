/*
version - Prints the installed version of sweet.

Usage:

	sweet version
*/
package version

import (
	"os"

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

	lexer := lexers.Match("source.js")

	return quick.Highlight(os.Stdout, code, lexer.Config().Name, "terminal", "default")
}
