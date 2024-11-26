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
		fmt.Printf("failed to connect to database: %s\n", err)
		return
	}

	reps, err := GetReps(statsDb, query)

	if err != nil {
		fmt.Printf("failed to get reps: %s\n", err)
		return
	}

	for _, rep := range reps {
		fmt.Printf("%s | %s | %2fs\n", rep.Start.Format("1/02/2006 15:04:05"), rep.Name, rep.Wpm)
	}

}

func init() {
	// TODO: define query flags
}
