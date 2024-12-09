package stats

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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
		q, err := queryFromArgs(cmd, time.Now())
		if err != nil {
			log.Fatal()
			os.Exit(-1)
		}
		printStats(q)
	},
}

// Returns the correct date depending on the start argument.
// The time returned is at midnight of the date specified by `start`,
// unless an hour query is used, in which case it will return the time
// specified by
func parseDateFromArg(isStart bool, arg string, now time.Time) (time.Time, error) {
	var todayAtMidnight time.Time
	todayAtMidnight = time.Date(now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location())

	if !isStart {
		todayAtMidnight = todayAtMidnight.AddDate(0, 0, 1).Add(-1 * time.Nanosecond)
	}

	if arg == "" {
		return todayAtMidnight, nil
	}

	unit := rune(arg[len(arg)-1])
	amount, err := strconv.Atoi(arg[:len(arg)-1])

	if err != nil {
		// maybe it's in date format?
		argTime, err := time.Parse(time.DateOnly, arg)

		if err != nil {
			return todayAtMidnight, fmt.Errorf("error parsing date: %s", err)
		}

		argTime = time.Date(
			argTime.Year(), argTime.Month(), argTime.Day(),
			0, 0, 0, 0, now.Location())

		if !isStart {
			argTime = argTime.AddDate(0, 0, 1).Add(-1 * time.Nanosecond)
		}

		if argTime.After(now) {
			return todayAtMidnight, fmt.Errorf("invalid date: %s hasn't happend yet!", arg)
		}

		return argTime, nil
	}

	switch unit {
	case 'H', 'h':
		return time.Date(
			now.Year(), now.Month(), now.Day(),
			now.Hour()-amount, now.Minute(), now.Second(), now.Nanosecond(),
			now.Location()), nil
	case 'D', 'd':
		return todayAtMidnight.AddDate(0, 0, -1*amount), nil
	case 'W', 'w':
		return todayAtMidnight.AddDate(0, 0, amount*-7), nil
	case 'M', 'm':
		return todayAtMidnight.AddDate(0, -1*amount, 0), nil
	case 'Y', 'y':
		return todayAtMidnight.AddDate(-1*amount, 0, 0), nil
	default:
		return todayAtMidnight, fmt.Errorf("invalid date format %s\nsee \"sweet stats --help\" for more details", arg)
	}
}

// Writes a query that retrieves the entries specified by the start and end
// date functions. Also handles the `since` variable, which is an alias for
// start.
func queryFromArgs(cmd *cobra.Command, now time.Time) (string, error) {
	// since := cmd.Flag("since").Value.String()
	start := cmd.Flag("start").Value.String()
	end := cmd.Flag("end").Value.String()

	startTime, _ := parseDateFromArg(true, start, now)
	endTime, _ := parseDateFromArg(false, end, now)

	query := fmt.Sprintf("select * from reps where start >= %d and end <= %d order by start desc;", startTime.UnixMilli(), endTime.UnixMilli())
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
	cmd.Flags().String("start", "", "find stats starting from this date")
	cmd.Flags().String("since", "", "alias for \"start\" flag")
	cmd.Flags().String("end", "", "find stats ending at this date")
	cmd.Flags().SortFlags = false
}

func init() {
	setStatsCommandFlags(Cmd)
}
