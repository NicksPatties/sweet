package root

import (
	"strconv"
	"testing"
)

const (
	lpinky uint = iota
	lring
	lmiddle
	lindex
	lthumb
	rthumb
	rindex
	rmiddle
	rring
	rpinky
)

var fm map[string][]int = make(map[string][]int)

// func charToFingers() []int {

// }

func fingerView() (view string) {
	c := "*"
	// fingers
	f := [][]string{
		{"0", c},
		{"1", c, c},
		{"2", c, c},
		{"3", c},
		{"4"},
		{"5"},
		{"6", c},
		{"7", c, c},
		{"8", c, c},
		{"9", c},
	}

	for row := 2; row >= 0; row-- {
		// curr fingers
		for cf := 0; cf < len(f); cf++ {
			// if this is a finger spot, then I should print the character
			// in the finger view location
			if isFingerSpot := row < len(f[cf]); isFingerSpot {
				// label?
				if row == 0 {
					view += strconv.Itoa(cf)
				} else {
					view += c
				}
			} else {
				view += " "
			}
			// space in between the hands
			if cf == 4 {
				view += " "
			}
		}
		// not last row
		if row != 0 {
			view += "\n"
		}
	}
	return
}

func TestFingerView(t *testing.T) {
	name := "No highlighted characters, stars for fingers"
	want := "" +
		" **     ** \n" +
		"****   ****\n" +
		"01234 56789"
	got := fingerView()
	if got != want {
		t.Errorf("%s:\ngot:\n%s\n\nwant:\n%s", name, got, want)
	}
}

func TestFindKeyCombo(t *testing.T) {
	testCases := []struct {
		name string
		char string
		want []string
	}{
		{
			name: "on unmodified keys",
			char: "a",
			want: []string{"a"},
		},
		{
			name: "on modified keys",
			char: "W",
			want: []string{"shift", "w"},
		},
		{
			name: "newline",
			char: "\n",
			want: []string{"↲"},
		},
		{
			name: "space",
			char: " ",
			want: []string{"space"},
		},
		{
			name: "no character",
			char: "",
			want: []string{},
		},
	}

	for _, tc := range testCases {
		g := qwerty.findKeyCombo(tc.char)
		for i, got := range g {
			want := tc.want[i]
			if got != want {
				t.Errorf("%s:\n\tgot %s\n\twant %s\n", tc.name, g, tc.want)
			}
		}
	}

}
