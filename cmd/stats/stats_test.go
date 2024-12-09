package stats

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestParseDateFromArg(t *testing.T) {
	// 2024-12-06 17:36:20.000000 -0700
	now := time.Date(2024, 12, 6, 17, 36, 20, 0, time.Now().Location())

	testCases := []struct {
		name    string
		isStart bool
		arg     string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "blank string",
			isStart: true,
			arg:     "",
			want:    time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
			wantErr: false,
		},
		{
			name:    "blank string, end date",
			isStart: false,
			arg:     "",
			want:    time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1).Add(-1 * time.Nanosecond),
			wantErr: false,
		},
		{
			name:    "--start=2H",
			isStart: true,
			arg:     "2H",
			want:    time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-2, now.Minute(), now.Second(), now.Nanosecond(), now.Location()),
			wantErr: false,
		},
		{
			name:    "--end=2H",
			isStart: false,
			arg:     "2H",
			want:    time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-2, now.Minute(), now.Second(), now.Nanosecond(), now.Location()),
			wantErr: false,
		},
		{
			name:    "--start=1D",
			isStart: true,
			arg:     "1D",
			want:    time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location()),
			wantErr: false,
		},
		{
			name:    "--end=1D",
			isStart: false,
			arg:     "1D",
			want:    time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1).Add(-1 * time.Nanosecond),
			wantErr: false,
		},
		{
			name:    "--start=2W",
			isStart: true,
			arg:     "2W",
			want:    time.Date(now.Year(), now.Month(), now.Day()-14, 0, 0, 0, 0, now.Location()),
			wantErr: false,
		},
		{
			name:    "--end=2W",
			isStart: false,
			arg:     "2W",
			want:    time.Date(now.Year(), now.Month(), now.Day()-14, 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1).Add(-1 * time.Nanosecond),
			wantErr: false,
		},
		{
			name:    "--start=1M",
			isStart: true,
			arg:     "1M",
			want:    time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, -1, 0),
			wantErr: false,
		},
		{
			name:    "--end=1M",
			isStart: false,
			arg:     "1M",
			want:    time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, -1, 0).AddDate(0, 0, 1).Add(-1 * time.Nanosecond),
			wantErr: false,
		},
		{
			name:    "--start=1Y",
			isStart: true,
			arg:     "1Y",
			want:    time.Date(now.Year()-1, now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
			wantErr: false,
		},
		{
			name:    "--end=1Y",
			isStart: false,
			arg:     "1Y",
			want:    time.Date(now.Year()-1, now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1).Add(-1 * time.Nanosecond),
			wantErr: false,
		},
		{
			name:    "--start=1X (invalid shorthand)",
			isStart: true,
			arg:     "1X",
			want:    time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
			wantErr: true,
		},
		{
			name:    "--end=1X (invalid shorthand)",
			isStart: false,
			arg:     "1X",
			want:    time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1).Add(-1 * time.Nanosecond),
			wantErr: true,
		},
		{
			name:    "--start=barf (invalid input)",
			isStart: true,
			arg:     "barf",
			want:    time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
			wantErr: true,
		},
		{
			name:    "--end=barf (invalid input)",
			isStart: false,
			arg:     "barf",
			want:    time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1).Add(-1 * time.Nanosecond),
			wantErr: true,
		},
		{
			name:    "--start=2011-10-01",
			isStart: true,
			arg:     "2011-10-01",
			want:    time.Date(2011, time.October, 1, 0, 0, 0, 0, now.Location()),
			wantErr: false,
		},
		{
			name:    "--end=2011-10-01",
			isStart: false,
			arg:     "2011-10-01",
			want:    time.Date(2011, time.October, 1, 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1).Add(-1 * time.Nanosecond),
			wantErr: false,
		},
		{
			name:    "--start=2222-10-10 (invalid future date)",
			isStart: true,
			arg:     "2222-10-10",
			want:    time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
			wantErr: true,
		},
		{
			name:    "--end=2222-10-10 (invalid future date)",
			isStart: false,
			arg:     "2222-10-10",
			want:    time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1).Add(-1 * time.Nanosecond),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		got, gotErr := parseDateFromArg(tc.isStart, tc.arg, now)

		if tc.wantErr && gotErr == nil {
			t.Errorf("%s: wanted error, but got nil", tc.name)
		}

		if !tc.want.Equal(got) {
			t.Errorf("%s:\n got  %v\n want %v", tc.name, got, tc.want)
		}
	}
}

func TestQueryFromArgs(t *testing.T) {
	type testCase struct {
		name    string
		in      []string
		want    string
		wantErr bool
	}

	var mockCmd = func(tc testCase) *cobra.Command {
		cmd := &cobra.Command{
			Args: cobra.MaximumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				got, err := queryFromArgs(cmd, args)
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
			name:    "default case",
			in:      []string{},
			want:    "select * from reps order by start desc;",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		cmd := mockCmd(tc)

		if err := cmd.Execute(); err != nil {
			t.Fatalf("mock command failed to run: %s", err)
		}
	}
}
