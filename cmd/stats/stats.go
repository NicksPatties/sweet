package stats

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	. "github.com/NicksPatties/sweet/db"
	"github.com/NicksPatties/sweet/util"
	g "github.com/guptarohit/asciigraph"
	tw "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "stats",
	Short: "Print statistics about typing exercises",
	RunE: func(cmd *cobra.Command, args []string) error {
		q, err := argsToQuery(cmd, time.Now())
		if err != nil {
			return err
		}
		reps, err := queryToReps(q)
		if err != nil {
			return err
		}
		printStats(cmd, reps)
		return nil
	},
}

// Returns the date depending on the provided argument.
// If this is an argument for the `start` or `since` flag, then return the date
// corresponding to midnight of the provided argument, otherwise return
// the date at one nanosecond before midnight of the next day.
// This function supports both `YYYY-MM-DD` arguments as well as `N[HDWMY]` arguments,
// where `N` is the number of hours, days, weeks, months, or years before `now`.
// The `now` argument is used primarily for testing, and is typically assigned to
// `time.Now()` in normal usage.
func parseDateFromArg(isEnd bool, arg string, now time.Time) (time.Time, error) {
	todayAtMidnight := time.Date(now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location())

	if isEnd {
		todayAtMidnight = todayAtMidnight.AddDate(0, 0, 1).Add(-1 * time.Nanosecond)
	}

	if arg == "" {
		return todayAtMidnight, nil
	}

	unit := rune(arg[len(arg)-1])
	amount, err := strconv.Atoi(arg[:len(arg)-1])

	if err != nil {
		// maybe it's in YYYY-MM-DD format?
		argTime, err := time.Parse(time.DateOnly, arg)

		if err != nil {
			return todayAtMidnight, fmt.Errorf("error parsing date: %s", err)
		}

		argTime = time.Date(
			argTime.Year(), argTime.Month(), argTime.Day(),
			0, 0, 0, 0, now.Location())

		if isEnd {
			argTime = argTime.AddDate(0, 0, 1).Add(-1 * time.Nanosecond)
		}

		if argTime.After(now) {
			return todayAtMidnight, fmt.Errorf("invalid date: %s hasn't happened yet!", arg)
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

// Converts the flags assigned to the stats command into an SQLite query,
// retrieving all entries from the database that match the query.
func argsToQuery(cmd *cobra.Command, now time.Time) (string, error) {
	filters := []string{}
	name := cmd.Flag(NAME).Value.String()
	lang := cmd.Flag(LANGUAGE).Value.String()

	if name != "" && lang != "" {
		return "", fmt.Errorf("both name and lang provided (please pick one of them!)")
	} else if lang != "" {
		filters = append(filters, fmt.Sprintf("%s='%s'", LANGUAGE, lang))
	} else if name != "" {
		nameFilter := fmt.Sprintf("%s like '%s'", NAME, name)
		nameFilter = strings.Replace(nameFilter, "*", "%", -1)
		filters = append(filters, nameFilter)
	}

	end := cmd.Flag(END).Value.String()
	since := cmd.Flag("since").Value.String()
	start := cmd.Flag(START).Value.String()

	if end != "" && since == "" && start == "" {
		return "", fmt.Errorf("must define %s if %s is provided", START, END)
	}

	if since != "" && start != "" {
		return "", fmt.Errorf("both since and start flags are provided. please only one of them.")
	} else if since != "" && start == "" {
		start = since
	}

	startTime, err := parseDateFromArg(false, start, now)
	if err != nil {
		return "", fmt.Errorf("failed to parse start flag: %s", err)
	}
	filters = append(filters, fmt.Sprintf("%s >= %d", START, startTime.UnixMilli()))

	endTime, err := parseDateFromArg(true, end, now)
	if err != nil {
		return "", err
	}
	filters = append(filters, fmt.Sprintf("%s <= %d", END, endTime.UnixMilli()))

	if endTime.Before(startTime) {
		return "", fmt.Errorf("%s is before %s", END, START)
	}

	query := fmt.Sprintf("select * from reps where %s order by %s desc;", strings.Join(filters, " and "), START)

	return query, nil
}

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
	cols := []string{START}
	name := cmd.Flag(NAME).Value.String()
	showName := name == "" || strings.Contains(name, "*")
	if showName {
		cols = append(cols, NAME)
	}
	possibleCols := []string{
		WPM, RAW_WPM, ACCURACY,
		UNCORRECTED_ERRORS, MISTAKES, DURATION,
	}
	selectedColCount := 0

	for _, col := range possibleCols {
		showCol := cmd.Flag(col).Value.String() == "true"
		if showCol {
			cols = append(cols, col)
			selectedColCount++
		}
	}

	if selectedColCount == 0 {
		defaultCols := append(cols,
			WPM, RAW_WPM, ACCURACY, UNCORRECTED_ERRORS, MISTAKES,
		)
		return defaultCols
	} else {
		return cols
	}
}

func getStatsHeader(name string, lang string, start string, end string) (header string) {
	nameSection := ""
	if name != "" {
		nameSection = name
	} else if lang != "" {
		// TODO make a map of languages to human readable names
		nameSection = lang
	}

	dateSection := "from today"
	if start != "" {
		if end != "" {
			dateSection = fmt.Sprintf("from %s to %s", start, end)
		} else {
			dateSection = fmt.Sprintf("since %s", start)
		}
	}

	tmplData := struct {
		// name of the file or language of the file
		Name string
		// "since X days", or "from YYYY-MM-DD to YYYY-MM-DD", etc
		Date string
	}{
		Name: nameSection,
		Date: dateSection,
	}

	tmplString := "" +
		`stats {{if .Name}}for {{.Name}} {{end}}{{.Date}}:`

	tmpl := template.Must(template.New("header").Parse(tmplString))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, tmplData)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// Calculates the average, min, max, first, last, and delta of a stat.
// Returns that information in an array of strings.
// Formats the data based on the column name selected.
func getColumnStats(reps []Rep, colName string) []string {
	// TODO if there is only one rep, return just those stats.

	getCol := func(r Rep, colName string) float64 {
		switch colName {
		case WPM:
			return r.Wpm
		case RAW_WPM:
			return r.Raw
		case DURATION:
			return float64(r.Dur)
		case ACCURACY:
			return r.Acc
		case MISTAKES:
			return float64(r.Miss)
		case UNCORRECTED_ERRORS:
			return float64(r.Errs)
		}
		return -math.MaxFloat64
	}

	avg := 0.0
	min := math.MaxFloat64
	max := -math.MaxFloat64
	first := getCol(reps[len(reps)-1], colName)
	last := getCol(reps[0], colName)
	delta := last - first

	sum := 0.0
	for _, rep := range reps {
		curr := getCol(rep, colName)
		sum += curr
		if curr < min {
			min = curr
		}
		if curr > max {
			max = curr
		}
	}
	avg = sum / float64(len(reps))

	colData := []float64{avg, min, max, first, last, delta}

	row := []string{}

	for _, d := range colData {
		row = append(row, util.ColumnString(colName, d))
	}
	return row
}

func printStats(cmd *cobra.Command, reps []Rep) {
	name := cmd.Flag(NAME).Value.String()
	lang := cmd.Flag(LANGUAGE).Value.String()
	start := cmd.Flag(START).Value.String()
	if since := cmd.Flag("since").Value.String(); since != "" {
		start = since
	}
	end := cmd.Flag(END).Value.String()

	// print the header
	fmt.Println(getStatsHeader(name, lang, start, end))

	cols := argsToColumnFilter(cmd)

	if len(reps) == 0 {
		fmt.Println("no stats")
	} else if len(reps) > 1 {

		// print the stats table
		table := tw.NewWriter(os.Stdout)
		table.SetHeader([]string{"", "avg", "min", "max", "first", "last", "delta"})
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tw.ALIGN_LEFT)
		table.SetAlignment(tw.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetTablePadding("  ") // pad with tabs
		table.SetNoWhiteSpace(true)
		for _, col := range cols {
			if col == START || col == NAME {
				continue
			}
			row := []string{col}
			row = append(row, getColumnStats(reps, col)...)
			table.Append(row)
		}
		table.Render()

		wpmGraphData := []float64{}
		rawWpmGraphData := []float64{}
		mistakesGraphData := []float64{}
		accuracyGraphData := []float64{}
		errorsGraphData := []float64{}
		for i := range reps {
			// reverse the order of the reps
			// so the date increases as X increases
			currRep := reps[len(reps)-1-i]
			wpmGraphData = append(wpmGraphData, currRep.Wpm)
			mistakesGraphData = append(mistakesGraphData, float64(currRep.Miss))
			rawWpmGraphData = append(rawWpmGraphData, currRep.Raw)
			accuracyGraphData = append(accuracyGraphData, currRep.Acc)
			errorsGraphData = append(errorsGraphData, float64(currRep.Errs))
		}

		// the data for each column's graph
		graphDatum := [][]float64{accuracyGraphData, errorsGraphData, mistakesGraphData, rawWpmGraphData, wpmGraphData}
		graphColors := []g.AnsiColor{g.Green, g.Red, g.Yellow, g.Gray, g.Default}
		graphLegends := []string{ACCURACY, UNCORRECTED_ERRORS, MISTAKES, RAW_WPM, WPM}

		graph := g.PlotMany(
			graphDatum,
			g.SeriesColors(graphColors...),
			g.SeriesLegends(graphLegends...),
			g.Height(10),
			g.Width(0), // auto scaling
			g.LowerBound(0),
			g.Precision(0),
		)

		fmt.Println(graph)

	}

	// print the reps table
	table := tw.NewWriter(os.Stdout)
	table.SetHeader(cols)
	for _, rep := range reps {
		repCols := []string{}
		for _, c := range cols {
			repCols = append(repCols, rep.ColumnString(c))
		}
		table.Append(repCols)
	}
	table.Render()
}

func setStatsCommandFlags(cmd *cobra.Command) {
	// date selection flags
	cmd.Flags().String(START, "", "find stats starting from this date")
	cmd.Flags().String("since", "", "alias for \"start\" flag")
	cmd.Flags().String(END, "", "find stats ending at this date")

	// column filtering flags
	cmd.Flags().String(NAME, "", "filter by exercise name")
	cmd.Flags().String(LANGUAGE, "", "filter by language")
	cmd.Flags().Bool(WPM, false, "show words per minute (wpm)")
	cmd.Flags().Bool(RAW_WPM, false, "show raw words per minute")
	cmd.Flags().Bool(ACCURACY, false, "show accuracy (acc)")
	cmd.Flags().Bool(MISTAKES, false, "show mistakes")
	cmd.Flags().Bool(UNCORRECTED_ERRORS, false, "show uncorrected errors")
	cmd.Flags().Bool(DURATION, false, "show duration")

	cmd.Flags().SortFlags = false
}

func init() {
	setStatsCommandFlags(Cmd)
}
