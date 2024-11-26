package stats

import (
	"testing"
	"time"
)

func TestShorthandToDateRange(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		// input N[H,D,W,M,Y]
		in      string
		want    dateRange
		wantErr bool
	}{
		{
			// incorrect number
			in:      "-1D",
			want:    dateRange{},
			wantErr: true,
		},
		{
			// incorrect date range code
			in:      "-1F",
			want:    dateRange{},
			wantErr: true,
		},
		{
			// one days (lowercase letters work)
			in: "1d",
			want: dateRange{
				start: now.AddDate(0, 0, -1),
				end:   now,
			},
			wantErr: false,
		},
		{
			// two days (n works)
			in: "2D",
			want: dateRange{
				start: now.AddDate(0, 0, -2),
				end:   now,
			},
			wantErr: false,
		},
		{
			// one week (W works)
			in: "1H",
			want: dateRange{
				start: now.Add(-1 * time.Hour),
				end:   now,
			},
			wantErr: false,
		},
		{
			// one week (W works)
			in: "1W",
			want: dateRange{
				start: now.AddDate(0, 0, -7),
				end:   now,
			},
			wantErr: false,
		},
		{
			// one month (Y works)
			in: "1M",
			want: dateRange{
				end:   now,
				start: now.AddDate(0, -1, 0),
			},
			wantErr: false,
		},
		{
			// one year (Y works)
			in: "1Y",
			want: dateRange{
				end:   now,
				start: now.AddDate(-1, 0, 0),
			},
			wantErr: false,
		},
		{
			// 100 days (n > 9 works)
			in: "100D",
			want: dateRange{
				end:   now,
				start: now.AddDate(0, -100, 0),
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		out, err := shorthandToDateRange(tc.in, now)

		if err == nil && tc.wantErr {
			t.Errorf("%s wanted error, but it's nil", tc.in)
		}

		if !(out.start.Equal(tc.want.start) || out.end.Equal(tc.want.end)) {
			t.Errorf(
				"%s\n"+
					"  got  start %s\n"+
					"  got  end   %s\n"+
					"  want start %s\n"+
					"  want end   %s\n",
				tc.in,
				out.start, out.end,
				tc.want.start, tc.want.end,
			)
		}
	}
}
