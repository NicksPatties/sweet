/*
version - Prints the installed version of sweet.

Usage:

	sweet version
*/
package version

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	sitter "github.com/tree-sitter/go-tree-sitter"
	javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "prints the version",
	Run: func(cmd *cobra.Command, args []string) {
		// just doing stupid version stuff
		codeSnippet()
	},
}

const (
	Reset      = "\033[0m"
	Bold       = "\033[1m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Magenta    = "\033[35m"
	Cyan       = "\033[36m"
	White      = "\033[37m"
	BrightBlue = "\033[94m"
)

// NodeType to color mapping
var colorMap = map[string]string{
	"identifier":           Cyan,
	"property_identifier":  Cyan,
	"string":               Green,
	"number":               Magenta,
	"true":                 Yellow,
	"false":                Yellow,
	"null":                 Yellow,
	"comment":              White,
	"function":             Blue,
	"function_declaration": Blue,
	"arrow_function":       Blue,
	"method_definition":    Blue,
	"keyword":              Yellow,
	"return":               Yellow,
	"if":                   Yellow,
	"else":                 Yellow,
	"for":                  Yellow,
	"while":                Yellow,
	"var":                  Yellow,
	"let":                  Yellow,
	"const":                Yellow,
}

func codeSnippet() {
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
	// This would probably be loaded from a file instead
	// of from a variable
	queryString := `
(function_expression
  name: (identifier) @function)
(function_declaration
  name: (identifier) @function)
(method_definition
  name: (property_identifier) @function.method)
`
	parser := sitter.NewParser()
	defer parser.Close()
	js := sitter.NewLanguage(javascript.Language())
	parser.SetLanguage(js)

	tree := parser.Parse([]byte(code), nil)
	defer tree.Close()
	n := tree.RootNode()

	q, err := sitter.NewQuery(js, queryString)
	if err != nil {
		fmt.Printf("error setting up query: %s", err)
		os.Exit(1)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	captures := qc.Captures(q, n, []byte(code))

	// iterate over all the captures
	// this will probably be the building the information to color the string
	for cap, i := captures.Next(); cap != nil; cap, i = captures.Next() {
		nodes := cap.NodesForCaptureIndex(i)

		for j, node := range nodes {
			start, end := node.Range().StartPoint, node.Range().EndPoint
			node.Range()
			fmt.Printf("  %d: %s %s start (%d, %d) end (%d, %d)\n", j, node.Parent().Kind(), node.Kind(), start.Row, start.Column, end.Row, end.Column)
		}
	}
}
