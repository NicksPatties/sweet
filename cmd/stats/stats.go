package stats

import (
	"fmt"
	"strings"

	. "github.com/NicksPatties/sweet/db"
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
		printStats()
	},
}

func argsToColumnFilter() []string {
	defaultCols := []string{"start", "name", "wpm", "errs", "miss", "acc"}
	return defaultCols
}

// Prints the columns
func printStats() {
	// connect to db
	statsDb, err := SweetDb()
	if err != nil {
		fmt.Printf("failed to connect to database: %s\n", err)
		return
	}

	reps, err := GetReps(statsDb)

	if err != nil {
		fmt.Printf("failed to get reps: %s\n", err)
		return
	}

	cols := argsToColumnFilter()
	// print the header
	fmt.Printf("%s\n", strings.Join(cols, "\t"))

	for _, rep := range reps {
		repCols := []string{}
		for _, c := range cols {
			repCols = append(repCols, rep.ColumnString(c))
		}
		fmt.Printf("%s\n", strings.Join(repCols, "\t"))
	}

}

func init() {
	// TODO: define query flags
}
