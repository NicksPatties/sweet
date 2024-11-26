package stats

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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
		fmt.Println(args)
		start := cmd.Flag("start")
		fmt.Println(start)

		q, err := queryFromArgs(cmd, args)
		if err != nil {
			log.Fatal()
			os.Exit(-1)
		}
		printStats(q)
	},
}

type dateRange struct {
	start time.Time
	end   time.Time
}

// Converts the N[H,D,W,M,Y] format string to a dateRange.
// The end parameter represents the end time of the date range,
// usually time.Now().
//
// If the function fails to parse the arg variable, then it
// returns an error.
//
// H - hours, D - days, W - weeks, M - months, Y - years
func shorthandToDateRange(arg string, end time.Time) (dateRange, error) {
	failedToParse := errors.New("failed to parse argument " + arg)

	hours := 0
	days := 1
	weeks := 0
	months := 0
	years := 0

	// The number of a specific units
	nString := string(arg[:len(arg)-1])
	n, err := strconv.Atoi(nString)
	if err != nil || n <= 0 {
		return dateRange{}, failedToParse
	}

	// The unit of date range [H,D,W,M,Y]
	unit := rune(arg[len(arg)-1])

	switch unit {
	case 'H', 'h':
		hours = n
		break
	case 'D', 'd':
		days = n
		break
	case 'W', 'w':
		weeks = n
		break
	case 'M', 'm':
		months = n
		break
	case 'Y', 'y':
		years = n
		break
	default:
		return dateRange{}, failedToParse
	}

	correctHrs := time.Duration(int64(-1) * int64(hours) * int64(time.Hour))

	return dateRange{
		start: end.AddDate(-1*years, -1*months, -1*(days+7*weeks)).Add(correctHrs),
		end:   end,
	}, nil
}

func queryFromArgs(cmd *cobra.Command, args []string) (string, error) {
	query := `select * from reps order by start desc;`
	return query, nil
}

func argsToColumnFilter() []string {
	defaultCols := []string{"start", "name", "wpm", "errs", "miss", "acc"}
	return defaultCols
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

func setStatsCommandFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("start", "s", "", "the start date")
	cmd.Flags().String("since", "", "alias for \"start\" flag")
	cmd.Flags().SortFlags = false
}

func init() {
	setStatsCommandFlags(Cmd)
}
