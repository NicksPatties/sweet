/*
version - Prints the installed version of sweet.

Usage:

	sweet version
*/
package version

import (
	"fmt"
	"os"
	"strings"

	// lg "github.com/charmbracelet/lipgloss"
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
	"identifier":          Cyan,
	"property_identifier": Cyan,
	"string":              Green,
	"number":              Magenta,
	"true":                Yellow,
	"false":               Yellow,
	"null":                Yellow,
	"comment":             White,
	// "function":             Blue,
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

var jsKeywords = map[string]bool{
	"break":      true,
	"case":       true,
	"catch":      true,
	"class":      true,
	"const":      true,
	"continue":   true,
	"debugger":   true,
	"default":    true,
	"delete":     true,
	"do":         true,
	"else":       true,
	"export":     true,
	"extends":    true,
	"finally":    true,
	"for":        true,
	"function":   true,
	"if":         true,
	"import":     true,
	"in":         true,
	"instanceof": true,
	"new":        true,
	"return":     true,
	"super":      true,
	"switch":     true,
	"this":       true,
	"throw":      true,
	"try":        true,
	"typeof":     true,
	"var":        true,
	"void":       true,
	"while":      true,
	"with":       true,
	"yield":      true,
	"async":      true,
	"await":      true,
	"let":        true,
	"static":     true,
	"implements": true,
	"interface":  true,
	"package":    true,
	"private":    true,
	"protected":  true,
	"public":     true,
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

	for cap, i := captures.Next(); cap != nil; cap, i = captures.Next() {
		nodes := cap.NodesForCaptureIndex(i)
		fmt.Println(nodes)

		for j, node := range nodes {
			start, end := node.Range().StartPoint, node.Range().EndPoint
			node.Range()
			fmt.Printf("  %d: %s start (%d, %d) end (%d, %d)\n", j, node.Kind(), start.Row, start.Column, end.Row, end.Column)
		}
	}

	highlighted := processSource(n, []byte(code))

	// write a query
	// execute the query
	// get the ranges to color for the query
	// get the source code lines
	// for each line
	//   color the text within range
	fmt.Println(highlighted)
}

func processSource(rootNode *sitter.Node, source []byte) string {
	var result strings.Builder

	// Create a map of byte positions to node information
	positionMap := make(map[uint]*nodeInfo)
	populatePositionMap(rootNode, source, positionMap)

	// Process each byte in the source
	for i := uint(0); i < uint(len(source)); i++ {
		if nodeInfo, exists := positionMap[i]; exists {
			// This is the start of a node we want to highlight
			nodeEnd := nodeInfo.endByte
			nodeText := string(source[i:nodeEnd])

			// Add the highlighted text
			result.WriteString(nodeInfo.color + nodeText + Reset)

			// Skip to the end of this node
			i = nodeEnd - 1
		} else {
			// This is whitespace or a character not part of a syntax node we're highlighting
			result.WriteByte(source[i])
		}
	}

	return result.String()
}

type nodeInfo struct {
	endByte uint
	color   string
}

func populatePositionMap(node *sitter.Node, source []byte, positionMap map[uint]*nodeInfo) {
	if node == nil {
		return
	}

	nodeType := node.Kind()

	// Check if this is a node we want to highlight
	if color, exists := colorMap[nodeType]; exists {
		startByte := node.StartByte()
		endByte := node.EndByte()

		// Special handling for keywords
		if nodeType == "identifier" {
			nodeText := string(source[startByte:endByte])
			if jsKeywords[nodeText] {
				color = Yellow + Bold
			}
		}

		// Only add to the map if we haven't processed this position yet
		if _, exists := positionMap[startByte]; !exists {
			positionMap[startByte] = &nodeInfo{
				endByte: endByte,
				color:   color,
			}
		}
	}

	// Process child nodes
	for i := uint(0); i < node.ChildCount(); i++ {
		populatePositionMap(node.Child(i), source, positionMap)
	}
}
