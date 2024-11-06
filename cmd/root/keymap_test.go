package root

import (
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

func TestFingerView(t *testing.T) {

	testCases := []struct {
		name   string
		char   rune
		margin int
		want   string
	}{
		{
			name:   "Stars for fingers",
			char:   '*',
			margin: 0,
			want: "" +
				" **     ** \n" +
				"****   ****\n" +
				"01234 56789",
		},
		{
			name:   "Stars for fingers, margins",
			char:   '*',
			margin: 2,
			want: "" +
				"   **     ** \n" +
				"  ****   ****\n" +
				"  01234 56789",
		},
	}
	for _, tc := range testCases {
		got := fingerView(tc.char, tc.margin)
		if got != tc.want {
			t.Errorf("%s:\ngot:\n%s\n\nwant:\n%s", tc.name, got, tc.want)
		}
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
