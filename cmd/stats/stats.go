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

var Cmd = &cobra.Command{
	Use:   "stats",
	Short: "Print statistics about typing exercises",
	Run: func(cmd *cobra.Command, args []string) {
		q, err := argsToQuery(cmd, time.Now())
		if err != nil {
			log.Fatal()
			os.Exit(-1)
		}
		reps, err := queryToReps(q)
		if err != nil {
			log.Fatal()
			os.Exit(-1)
		}
		cols := argsToColumnFilter(cmd)
		printStats(reps, cols)
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
func argsToQuery(cmd *cobra.Command, now time.Time) (string, error) {
	filters := []string{}
	lang := cmd.Flag("lang").Value.String()

	if lang != "" {
		filters = append(filters, fmt.Sprintf("lang='%s'", lang))
	}

	end := cmd.Flag("end").Value.String()
	since := cmd.Flag("since").Value.String()
	start := cmd.Flag("start").Value.String()

	if end != "" && since == "" && start == "" {
		return "", fmt.Errorf("must define start if end is provided")
	}

	if since != "" && start != "" {
		fmt.Printf("both `--since` and `--start` variables provided (you only need one of them!)")
	} else if since != "" && start == "" {
		start = since
	}

	startTime, _ := parseDateFromArg(true, start, now)
	filters = append(filters, fmt.Sprintf("start >= %d", startTime.UnixMilli()))
	endTime, _ := parseDateFromArg(false, end, now)
	filters = append(filters, fmt.Sprintf("end <= %d", endTime.UnixMilli()))
	// TODO finish creating filters for language and name properties

	if endTime.Before(startTime) {
		return "", fmt.Errorf("end is before start")
	}

	filterText := strings.Join(filters, " and ")

	query := fmt.Sprintf("select * from reps where %s order by start desc;", filterText)
	return query, nil
}

// Converts a query to an array of Reps, which are the individual
// rows that appear in stats database
func queryToReps(query string) (reps []Rep, err error) {
	statsDb, err := SweetDb()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s\n", err)
	}

	reps, err = GetReps(statsDb, query)

	if err != nil {
		return nil, fmt.Errorf("failed to get reps: %s\n", err)
	}

	return
}

func argsToColumnFilter(cmd *cobra.Command) []string {
	cols := []string{"start"}
	if cmd.Flag("name").Value.String() == "" {
		cols = append(cols, "name")
	}
	defaultCols := append(cols, "wpm", "raw", "acc", "errs", "miss")
	possibleCols := []string{"wpm", "raw", "acc", "errs", "miss", "dur"}
	selectedColCount := 0

	for _, col := range possibleCols {
		isThere := cmd.Flag(col).Value.String() == "true"
		if isThere {
			cols = append(cols, col)
			selectedColCount++
		}
	}

	if selectedColCount == 0 {
		return defaultCols
	} else {
		return cols
	}
}

// Prints the columns
func printStats(reps []Rep, cols []string) {
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
	// date selection flags
	cmd.Flags().String("start", "", "find stats starting from this date")
	cmd.Flags().String("since", "", "alias for \"start\" flag")
	cmd.Flags().String("end", "", "find stats ending at this date")

	// column filtering flags
	cmd.Flags().String("name", "", "filter by exercise name")
	cmd.Flags().String("lang", "", "filter by language")
	cmd.Flags().Bool("wpm", false, "show words per minute (wpm)")
	cmd.Flags().Bool("raw", false, "show raw words per minute")
	cmd.Flags().Bool("acc", false, "show accuracy (acc)")
	cmd.Flags().Bool("miss", false, "show mistakes")
	cmd.Flags().Bool("errs", false, "show uncorrected errors")
	cmd.Flags().Bool("dur", false, "show duration")

	cmd.Flags().SortFlags = false
}

func init() {
	setStatsCommandFlags(Cmd)
}
