package stats

import (
	"fmt"
	"strings"

	. "github.com/NicksPatties/sweet/cmd/root"
	"github.com/spf13/cobra"
)

// var Cmd = &cobra.Command{
// 	Use:     "sweet [file]",
// 	Long:    fmt.Sprintf("%s.\nRuns an interactive touch typing game, and prints the results.", getProductTagline()),
// 	Args:    cobra.MaximumNArgs(1),
// 	Example: getExamples(),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		ex, err := fromArgs(cmd, args)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		Run(ex)
// 	},
// }

var Cmd = &cobra.Command{
	Use:   "stats",
	Short: "Print statistics about typing exercises",
	Run: func(cmd *cobra.Command, args []string) {
		query := argsToQuery()
		printStats(query)
	},
}

func argsToQuery() (q string) {
	defaultCols := []string{"start", "name", "wpm", "errs", "dur", "miss", "acc"}
	selectString := strings.Join(defaultCols, ", ")
	fmt.Printf("selectString: %v\n", selectString)

	q = fmt.Sprintf("select %s from reps;", selectString)
	return
}

// Prints the columns
func printStats(query string) {
	// connect to db
	statsDb, err := SweetDb()
	if err != nil {
		fmt.Println(err)
	}

	// query the data

	// print the data
}

func init() {
	// TODO: define query flags
}
