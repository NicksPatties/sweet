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

	g "github.com/guptarohit/asciigraph"
	tw "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	c "github.com/NicksPatties/sweet/constants"
	db "github.com/NicksPatties/sweet/db"
)

type tmplData map[string]any

var Cmd = &cobra.Command{
	Use:   "stats",
	Short: "Print typing exercise statistics",
	Args:  cobra.MaximumNArgs(0),
	Example: "  get the stats from today\n" +
		"  sweet stats\n\n" +
		"  get stats from the past two weeks\n" +
		"  sweet stats --since=2w\n\n" +
		"  get stats for November 2024\n" +
		"  sweet stats --start=2024-11-01 --end=2024-11-30\n\n" +
		"  get stats for Go exercises only\n" +
		"  sweet stats --lang=go\n\n" +
		"  get stats for a specific exercise name\n" +
		"  sweet stats --name=hello.go\n\n" +
		"  get stats for exercises that contain the name \"hello\"\n" +
		"  sweet stats --name=hello*\n\n" +
		"  get stats for words per minute and mistakes only\n" +
		"  sweet stats --wpm --miss",
	RunE: func(cmd *cobra.Command, args []string) error {
		q, err := argsToQuery(cmd, time.Now())
		if err != nil {
			return err
		}
		reps, err := queryToReps(q)
		if err != nil {
			return err
		}
		render(cmd, reps)
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
	name := cmd.Flag(c.NAME).Value.String()
	lang := cmd.Flag(c.LANGUAGE).Value.String()

	if name != "" && lang != "" {
		return "", fmt.Errorf("both name and lang provided (please pick one of them!)")
	} else if lang != "" {
		filters = append(filters, fmt.Sprintf("%s='%s'", c.LANGUAGE, lang))
	} else if name != "" {
		nameFilter := fmt.Sprintf("%s like '%s'", c.NAME, name)
		nameFilter = strings.Replace(nameFilter, "*", "%", -1)
		filters = append(filters, nameFilter)
	}

	end := cmd.Flag(c.END).Value.String()
	since := cmd.Flag("since").Value.String()
	start := cmd.Flag(c.START).Value.String()

	if end != "" && since == "" && start == "" {
		return "", fmt.Errorf("must define %s if %s is provided", c.START, c.END)
	}

	if since != "" && start != "" {
		return "", fmt.Errorf("both since and start flags are provided. please use one or the other.")
	} else if since != "" && start == "" {
		start = since
	}

	startTime, err := parseDateFromArg(false, start, now)
	if err != nil {
		return "", fmt.Errorf("failed to parse start flag: %s", err)
	}
	filters = append(filters, fmt.Sprintf("%s >= %d", c.START, startTime.UnixMilli()))

	endTime, err := parseDateFromArg(true, end, now)
	if err != nil {
		return "", fmt.Errorf("failed to parse end flag: %s", err)
	}
	filters = append(filters, fmt.Sprintf("%s <= %d", c.END, endTime.UnixMilli()))

	if endTime.Before(startTime) {
		return "", fmt.Errorf("%s is before %s", c.END, c.START)
	}

	query := fmt.Sprintf("select * from reps where %s order by %s;", strings.Join(filters, " and "), c.START)

	return query, nil
}

func queryToReps(query string) (reps []db.Rep, err error) {
	statsDb, err := db.SweetDb()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s\n", err)
	}

	reps, err = db.GetReps(statsDb, query)

	if err != nil {
		return nil, fmt.Errorf("failed to get reps: %s\n", err)
	}

	return
}

func argsToColumnFilter(cmd *cobra.Command) []string {
	cols := []string{c.START}
	name := cmd.Flag(c.NAME).Value.String()
	showName := name == "" || strings.Contains(name, "*")
	if showName {
		cols = append(cols, c.NAME)
	}
	possibleCols := []string{
		c.WPM, c.RAW_WPM, c.ACCURACY,
		c.UNCORRECTED_ERRORS, c.MISTAKES, c.DURATION,
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
			c.WPM, c.RAW_WPM, c.ACCURACY, c.UNCORRECTED_ERRORS, c.MISTAKES,
		)
		return defaultCols
	} else {
		return cols
	}
}

// Converts a date argument to human readable format.
// Assumes the dates that are passed into the function
// happened in the past, so words like "ago" are
// expected in the output.
//
// This also assumes the arg _will_ be in either the
// N[HDWMY] format or YYYY-MM-DD format
func dateArgToHumanReadable(arg string) string {
	unit := rune(arg[len(arg)-1])
	amount, err := strconv.Atoi(arg[:len(arg)-1])
	if err != nil {
		return arg
	}

	// pluralized strings
	type pString struct {
		singular string
		plural   string
	}

	var amountStr string
	switch amount {
	case 1:
		amountStr = "one"
	case 2:
		amountStr = "two"
	case 3:
		amountStr = "three"
	case 4:
		amountStr = "four"
	case 5:
		amountStr = "five"
	case 6:
		amountStr = "six"
	case 7:
		amountStr = "seven"
	case 8:
		amountStr = "eight"
	case 9:
		amountStr = "nine"
	default:
		amountStr = strconv.Itoa(amount)
	}

	var unitStr string
	unitStrs := map[rune]pString{
		'h': {
			singular: "hour",
			plural:   "hours",
		},
		'd': {
			singular: "day",
			plural:   "days",
		},
		'w': {
			singular: "week",
			plural:   "weeks",
		},
		'm': {
			singular: "month",
			plural:   "months",
		},
		'y': {
			singular: "year",
			plural:   "years",
		},
	}

	switch unit {
	case 'H', 'h':
		if amount == 1 {
			unitStr = unitStrs['h'].singular
		} else {
			unitStr = unitStrs['h'].plural
		}
	case 'D', 'd':
		if amount == 1 {
			unitStr = unitStrs['d'].singular
		} else {
			unitStr = unitStrs['d'].plural
		}
	case 'W', 'w':
		if amount == 1 {
			unitStr = unitStrs['w'].singular
		} else {
			unitStr = unitStrs['w'].plural
		}
	case 'M', 'm':
		if amount == 1 {
			unitStr = unitStrs['w'].singular
		} else {
			unitStr = unitStrs['w'].plural
		}
	}

	return fmt.Sprintf("%s %s ago", amountStr, unitStr)
}

// Creates the header of the stats command. This is a summary
// of the flags that were used to execute the command in human
// readable format.
func renderHeader(name string, lang string, start string, end string) {
	nameSection := ""
	if name != "" {
		nameSection = name
	} else if lang != "" {
		// TODO make a map of languages to human readable names
		nameSection = lang
	}

	dateSection := "from today"
	if start != "" {
		startStr := dateArgToHumanReadable(start)
		if end != "" {
			endStr := dateArgToHumanReadable(end)
			dateSection = fmt.Sprintf("from %s to %s", startStr, endStr)
		} else {
			dateSection = fmt.Sprintf("from %s", startStr)
		}
	}

	tmplString := `stats {{if .Name}}for {{.Name}} {{end}}{{.Date}}:`

	tmpl := template.Must(template.New("header").Parse(tmplString))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tmplData{
		"Name": nameSection,
		"Date": dateSection,
	}); err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
}

// Calculates the average, min, max, first, last, and delta of a stat.
// Returns that information in an array of strings.
// Formats the data based on the column name selected.
func getColumnStats(reps []db.Rep, colName string) []string {
	getCol := func(r db.Rep, colName string) float64 {
		switch colName {
		case c.WPM:
			return r.Wpm
		case c.RAW_WPM:
			return r.Raw
		case c.DURATION:
			return float64(r.Dur)
		case c.ACCURACY:
			return r.Acc
		case c.MISTAKES:
			return float64(r.Miss)
		case c.UNCORRECTED_ERRORS:
			return float64(r.Errs)
		}
		return -math.MaxFloat64
	}

	columnString := func(col string, value float64) string {
		switch col {
		case c.WPM:
			return fmt.Sprintf("%.f", value)
		case c.RAW_WPM:
			return fmt.Sprintf("%.f", value)
		case c.DURATION:
			d := time.Duration(value)
			return d.Round(time.Millisecond).String()
		case c.ACCURACY:
			return fmt.Sprintf("%.2f%%", value)
		case c.MISTAKES:
			return fmt.Sprintf("%.f", value)
		case c.UNCORRECTED_ERRORS:
			return fmt.Sprintf("%.f", value)
		default:
			return ""
		}
	}

	avg := 0.0
	min := math.MaxFloat64
	max := -math.MaxFloat64
	first := getCol(reps[0], colName)
	last := getCol(reps[len(reps)-1], colName)
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
		row = append(row, columnString(colName, d))
	}
	return row
}

func renderStatsTable(cols []string, reps []db.Rep) {
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
	table.SetTablePadding("  ")
	table.SetNoWhiteSpace(true)
	for _, col := range cols {
		if col == c.START || col == c.NAME {
			continue
		}
		row := []string{col}
		row = append(row, getColumnStats(reps, col)...)
		table.Append(row)
	}
	table.Render()
}

func renderGraph(cols []string, reps []db.Rep) {
	wpmGraphData := []float64{}
	rawWpmGraphData := []float64{}
	mistakesGraphData := []float64{}
	accuracyGraphData := []float64{}
	errorsGraphData := []float64{}
	for _, currRep := range reps {
		wpmGraphData = append(wpmGraphData, currRep.Wpm)
		mistakesGraphData = append(mistakesGraphData, float64(currRep.Miss))
		rawWpmGraphData = append(rawWpmGraphData, currRep.Raw)
		accuracyGraphData = append(accuracyGraphData, currRep.Acc)
		errorsGraphData = append(errorsGraphData, float64(currRep.Errs))
	}

	type plot struct {
		data   []float64
		color  g.AnsiColor
		legend string
	}
	plots := map[string]*plot{
		c.ACCURACY: {
			data:   accuracyGraphData,
			color:  g.Green,
			legend: "accuracy",
		},
		c.UNCORRECTED_ERRORS: {
			data:   errorsGraphData,
			color:  g.Red,
			legend: "errors",
		},
		c.MISTAKES: {
			data:   mistakesGraphData,
			color:  g.Yellow,
			legend: "mistakes",
		},
		c.RAW_WPM: {
			data:   rawWpmGraphData,
			color:  g.Gray,
			legend: "raw wpm",
		},
		c.WPM: {
			data:   wpmGraphData,
			color:  g.Default,
			legend: "wpm",
		},
	}

	colsRenderOrder := []string{c.ACCURACY, c.UNCORRECTED_ERRORS, c.MISTAKES, c.RAW_WPM, c.WPM}

	var graphDatum [][]float64
	var graphColors []g.AnsiColor
	var graphLegends []string
	for _, renderCol := range colsRenderOrder {
		argColFound := false
		for _, col := range cols {
			if col == renderCol {
				argColFound = true
				break
			}
		}

		if argColFound && plots[renderCol] != nil {
			graphDatum = append(graphDatum, plots[renderCol].data)
			graphColors = append(graphColors, plots[renderCol].color)
			graphLegends = append(graphLegends, plots[renderCol].legend)
		}
	}

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

func renderReps(cols []string, reps []db.Rep) {
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

func render(cmd *cobra.Command, reps []db.Rep) {
	name := cmd.Flag(c.NAME).Value.String()
	lang := cmd.Flag(c.LANGUAGE).Value.String()
	start := cmd.Flag(c.START).Value.String()
	if since := cmd.Flag("since").Value.String(); since != "" {
		start = since
	}
	end := cmd.Flag(c.END).Value.String()

	renderHeader(name, lang, start, end)

	cols := argsToColumnFilter(cmd)

	if len(reps) == 0 {
		fmt.Println("no stats")
	} else {
		renderStatsTable(cols, reps)
		renderGraph(cols, reps)
		renderReps(cols, reps)
	}
}

func setStatsCommandFlags(cmd *cobra.Command) {
	// date selection flags
	cmd.Flags().StringP(c.START, "s", "", "find stats starting from this date")
	cmd.Flags().String("since", "", "alias for \"start\" flag")
	cmd.Flags().StringP(c.END, "n", "", "find stats ending at this date")

	// column filtering flags
	cmd.Flags().String(c.NAME, "", "filter by exercise name")
	cmd.Flags().StringP(c.LANGUAGE, "l", "", "filter by language")
	cmd.Flags().BoolP(c.WPM, "w", false, "show words per minute (wpm)")
	cmd.Flags().BoolP(c.RAW_WPM, "r", false, "show raw words per minute")
	cmd.Flags().BoolP(c.ACCURACY, "a", false, "show accuracy (acc)")
	cmd.Flags().BoolP(c.MISTAKES, "m", false, "show mistakes")
	cmd.Flags().BoolP(c.UNCORRECTED_ERRORS, "e", false, "show uncorrected errors")
	cmd.Flags().BoolP(c.DURATION, "d", false, "show duration")

	cmd.Flags().SortFlags = false
}

func init() {
	setStatsCommandFlags(Cmd)
}
