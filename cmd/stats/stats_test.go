package stats

import (
	"fmt"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestParseDateFromArg(t *testing.T) {
	// 2024-12-06 17:36:20.000000 -0700
	now := time.Date(2024, 12, 6, 17, 36, 20, 0, time.Now().Location())
	nowAtMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	nowBeforeMidnight := nowAtMidnight.AddDate(0, 0, 1).Add(-1 * time.Nanosecond)

	testCases := []struct {
		name    string
		isEnd   bool
		arg     string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "blank string",
			isEnd:   false,
			arg:     "",
			want:    nowAtMidnight,
			wantErr: false,
		},
		{
			name:    "blank string, end date",
			isEnd:   true,
			arg:     "",
			want:    nowBeforeMidnight,
			wantErr: false,
		},
		{
			name:    "--start=2H",
			isEnd:   false,
			arg:     "2H",
			want:    time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-2, now.Minute(), now.Second(), now.Nanosecond(), now.Location()),
			wantErr: false,
		},
		{
			name:    "--end=2H",
			isEnd:   true,
			arg:     "2H",
			want:    time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-2, now.Minute(), now.Second(), now.Nanosecond(), now.Location()),
			wantErr: false,
		},
		{
			name:    "--start=1D",
			isEnd:   false,
			arg:     "1D",
			want:    nowAtMidnight.AddDate(0, 0, -1),
			wantErr: false,
		},
		{
			name:    "--end=1D",
			isEnd:   true,
			arg:     "1D",
			want:    nowBeforeMidnight.AddDate(0, 0, -1),
			wantErr: false,
		},
		{
			name:    "--start=2W",
			isEnd:   false,
			arg:     "2W",
			want:    nowAtMidnight.AddDate(0, 0, -14),
			wantErr: false,
		},
		{
			name:    "--end=2W",
			isEnd:   true,
			arg:     "2W",
			want:    nowBeforeMidnight.AddDate(0, 0, -14),
			wantErr: false,
		},
		{
			name:    "--start=1M",
			isEnd:   false,
			arg:     "1M",
			want:    nowAtMidnight.AddDate(0, -1, 0),
			wantErr: false,
		},
		{
			name:    "--end=1M",
			isEnd:   true,
			arg:     "1M",
			want:    nowBeforeMidnight.AddDate(0, -1, 0),
			wantErr: false,
		},
		{
			name:    "--start=1Y",
			isEnd:   false,
			arg:     "1Y",
			want:    nowAtMidnight.AddDate(-1, 0, 0),
			wantErr: false,
		},
		{
			name:    "--end=1Y",
			isEnd:   true,
			arg:     "1Y",
			want:    nowBeforeMidnight.AddDate(-1, 0, 0),
			wantErr: false,
		},
		{
			name:    "--start=1X (invalid shorthand)",
			isEnd:   false,
			arg:     "1X",
			want:    nowAtMidnight,
			wantErr: true,
		},
		{
			name:    "--end=1X (invalid shorthand)",
			isEnd:   true,
			arg:     "1X",
			want:    nowBeforeMidnight,
			wantErr: true,
		},
		{
			name:    "--start=barf (invalid input)",
			isEnd:   false,
			arg:     "barf",
			want:    nowAtMidnight,
			wantErr: true,
		},
		{
			name:    "--end=barf (invalid input)",
			isEnd:   true,
			arg:     "barf",
			want:    nowBeforeMidnight,
			wantErr: true,
		},
		{
			name:    "--start=2011-10-01",
			isEnd:   false,
			arg:     "2011-10-01",
			want:    time.Date(2011, time.October, 1, 0, 0, 0, 0, now.Location()),
			wantErr: false,
		},
		{
			name:    "--end=2011-10-01",
			isEnd:   true,
			arg:     "2011-10-01",
			want:    time.Date(2011, time.October, 1, 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1).Add(-1 * time.Nanosecond),
			wantErr: false,
		},
		{
			name:    "--start=2222-10-10 (invalid future date)",
			isEnd:   false,
			arg:     "2222-10-10",
			want:    nowAtMidnight,
			wantErr: true,
		},
		{
			name:    "--end=2222-10-10 (invalid future date)",
			isEnd:   true,
			arg:     "2222-10-10",
			want:    nowBeforeMidnight,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		got, gotErr := parseDateFromArg(tc.isEnd, tc.arg, now)

		if tc.wantErr && gotErr == nil {
			t.Errorf("%s: wanted error, but got nil", tc.name)
		}

		if !tc.want.Equal(got) {
			t.Errorf("%s:\n got  %v\n want %v", tc.name, got, tc.want)
		}
	}
}

func TestArgsToQuery(t *testing.T) {
	// 2024-12-06 17:36:20.000000 -0700
	now := time.Date(2024, 12, 6, 17, 36, 20, 0, time.Now().Location())
	nowAtMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	nowBeforeMidnight := nowAtMidnight.AddDate(0, 0, 1).Add(-1 * time.Nanosecond)

	type testCase struct {
		name    string
		in      []string
		want    string
		wantErr bool
	}

	var mockCmd = func(tc testCase) *cobra.Command {
		cmd := &cobra.Command{
			Run: func(cmd *cobra.Command, args []string) {
				got, err := argsToQuery(cmd, now)
				if err == nil && tc.wantErr {
					t.Errorf("%s wanted error, got nil", tc.name)
				}

				if got != tc.want {
					t.Errorf("%s\n"+
						"  got:  %s\n"+
						"  want: %s\n",
						tc.name, got, tc.want)
				}
			},
		}
		setStatsCommandFlags(cmd)
		cmd.SetArgs(tc.in)
		return cmd
	}

	testCases := []testCase{
		{
			name: "default case (get stats from today only)",
			in:   []string{},
			want: fmt.Sprintf(
				"select * from reps where start >= %d and end <= %d order by start;",
				nowAtMidnight.UnixMilli(),
				nowBeforeMidnight.UnixMilli(),
			),
			wantErr: false,
		},
		{
			name: "since is an alias for start",
			in:   []string{"--since=2D"},
			want: fmt.Sprintf(
				"select * from reps where start >= %d and end <= %d order by start;",
				nowAtMidnight.AddDate(0, 0, -2).UnixMilli(),
				nowBeforeMidnight.UnixMilli(),
			),
			wantErr: false,
		},
		{
			name:    "both since and start are given. error",
			in:      []string{"--since=2D", "--start=1D"},
			want:    "",
			wantErr: true,
		},
		{
			name: "start provided",
			in:   []string{"--start=1D"},
			want: fmt.Sprintf(
				"select * from reps where start >= %d and end <= %d order by start;",
				nowAtMidnight.AddDate(0, 0, -1).UnixMilli(),
				nowBeforeMidnight.UnixMilli(),
			),
			wantErr: false,
		},
		{
			name:    "end provided, but no start",
			in:      []string{"--end=1D"},
			want:    "",
			wantErr: true,
		},
		{
			name: "start and end provided",
			in:   []string{"--start=2024-10-01", "--end=2024-11-01"},
			want: fmt.Sprintf(
				"select * from reps where start >= %d and end <= %d order by start;",
				time.Date(2024, time.October, 1, 0, 0, 0, 0, now.Location()).UnixMilli(),
				time.Date(2024, time.November, 1, 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1).Add(-1*time.Nanosecond).UnixMilli()),
			wantErr: false,
		},
		{
			name:    "start and end provided, but end is before start",
			in:      []string{"--start=1D", "--end=3D"},
			want:    "",
			wantErr: true,
		},
		{
			name: "language provided",
			in:   []string{"--lang=py"},
			want: fmt.Sprintf(
				"select * from reps where lang='py' and start >= %d and end <= %d order by start;",
				nowAtMidnight.UnixMilli(),
				nowBeforeMidnight.UnixMilli(),
			),
			wantErr: false,
		},
		{
			name: "name provided",
			in:   []string{"--name=filename.go"},
			want: fmt.Sprintf(
				"select * from reps where name like 'filename.go' and start >= %d and end <= %d order by start;",
				nowAtMidnight.UnixMilli(),
				nowBeforeMidnight.UnixMilli(),
			),
			wantErr: false,
		},
		{
			name: "name with wildcard",
			in:   []string{"--name=file*"},
			want: fmt.Sprintf(
				"select * from reps where name like 'file%%' and start >= %d and end <= %d order by start;",
				nowAtMidnight.UnixMilli(),
				nowBeforeMidnight.UnixMilli(),
			),
			wantErr: false,
		},
		{
			name:    "both name and language provided, error",
			in:      []string{"--name=filename.go", "--lang=py"},
			want:    "",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		cmd := mockCmd(tc)

		if err := cmd.Execute(); err != nil {
			t.Fatalf("mock command failed to run: %s", err)
		}
	}
}

func TestArgsToColumnFilter(t *testing.T) {

	type testCase struct {
		name string
		in   []string
		want []string
	}

	var mockCmd = func(tc testCase) *cobra.Command {
		cmd := &cobra.Command{
			Run: func(cmd *cobra.Command, args []string) {
				got := argsToColumnFilter(cmd)
				if len(got) != len(tc.want) {
					t.Fatalf("%s\n"+
						"  got:  %s\n"+
						"  want: %s\n",
						tc.name, got, tc.want)

				}
				for i, want := range tc.want {
					if got[i] != want {
						t.Errorf("%s\n"+
							"  got:  %s\n"+
							"  want: %s\n",
							tc.name, got, tc.want)
						break
					}
				}

			},
		}
		setStatsCommandFlags(cmd)
		cmd.SetArgs(tc.in)
		return cmd
	}

	testCases := []testCase{
		{
			name: "default columns",
			in:   []string{},
			want: []string{"start", "name", "wpm", "raw", "acc", "errs", "miss"},
		},
		{
			name: "only one column",
			in:   []string{"--wpm"},
			want: []string{"start", "name", "wpm"},
		},
		{
			name: "some columns",
			in:   []string{"--raw", "--dur", "--miss"},
			want: []string{"start", "name", "raw", "miss", "dur"},
		},
		{
			name: "name provided, hides name column",
			in:   []string{"--name=hello.go"},
			want: []string{"start", "wpm", "raw", "acc", "errs", "miss"},
		},
		{
			name: "name wildcard provided, shows name column",
			in:   []string{"--name=hello*"},
			want: []string{"start", "name", "wpm", "raw", "acc", "errs", "miss"},
		},
		{
			name: "name and other columns provided, hides name column and only shows provided columns",
			in:   []string{"--raw", "--name=hello.go"},
			want: []string{"start", "raw"},
		},
	}
	for _, tc := range testCases {
		cmd := mockCmd(tc)

		if err := cmd.Execute(); err != nil {
			t.Fatalf("mock command failed to run: %s", err)
		}
	}
}
