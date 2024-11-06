package root

import (
	lg "github.com/charmbracelet/lipgloss"
	"testing"
)

type keymap struct {
	keys         [][]string
	modifiedKeys [][]string
	margins      []int
}

var qwerty = keymap{
	keys: [][]string{
		{
			"`", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-", "=",
		},
		{
			"q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "[", "]", "\\",
		},
		{
			"a", "s", "d", "f", "g", "h", "j", "k", "l", ";", "'", "↲",
		},
		// Spaces in this array are purely cosmetic. They're used to
		// add padding between "shift" and "z" in the keymap.
		{
			"shift", " ", "z", "x", "c", "v", "b", "n", "m", ",", ".", "/",
		},
		{
			"space",
		},
	},
	modifiedKeys: [][]string{
		{
			"~", "!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "+",
		},
		{
			"Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P", "{", "}", "|",
		},
		{
			"A", "S", "D", "F", "G", "H", "J", "K", "L", ":", "\"", "↲",
		},
		{
			"shift", " ", "Z", "X", "C", "V", "B", "N", "M", "<", ">", "?",
		},
		{
			"space",
		},
	},
	margins: []int{3, 4, 5, 0, 8},
}

func spaces(n int) (s string) {
	for i := 0; i < n; i++ {
		s += " "
	}
	return
}

func (k keymap) findKeyCombo(char string) (combo []string) {

	if char == "\n" {
		return []string{"↲"}
	}

	if char == " " {
		return []string{"space"}
	}

	for _, row := range k.keys {
		for _, key := range row {
			if key == char {
				return []string{key}
			}
		}
	}

	for i, row := range k.modifiedKeys {
		for j, key := range row {
			if key == char {
				return []string{"shift", k.keys[i][j]}
			}
		}
	}

	return
}

// Renders the keymap. Uses a key that needs to be
// rendered as input.
func (k keymap) render(char string) (km string) {
	// red style
	r := lg.NewStyle().Foreground(lg.Color("#FF0000")).Bold(true)
	combo := k.findKeyCombo(char)
	var currKey string
	if len(combo) == 0 {
		currKey = ""
	} else {
		currKey = combo[len(combo)-1]
	}
	isShift := len(combo) > 1
	rows := len(k.keys)
	for ri, row := range k.keys {
		km += spaces(k.margins[ri])
		for _, key := range row {
			if key == currKey || key == "shift" && isShift {
				km += r.Render(key)
			} else {
				km += key
			}
		}
		if ri != rows-1 {
			km += "\n"
		}
	}
	return
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

// This function doesn't actually test what I want it to.
// I would like to see if the correct colors are set to red and bold,
// but my tests pass in both cases, even though the output is different.
func TestRender(t *testing.T) {
	r := lg.NewStyle().Foreground(lg.Color("#FF0000")).Bold(true)

	testCases := []struct {
		name string
		char string
		want string
	}{
		{
			name: "show keys, 'a' is highlighted",
			char: "a",
			want: "   `1234567890-=\n" +
				"    qwertyuiop[]\\\n" +
				"     " + r.Render("a") + "sdfghjkl;'↲\n" +
				"shift zxcvbnm,./\n" +
				"        space",
		},
		{
			name: "show keys, no keys highlighted",
			char: "",
			want: "   `1234567890-=\n" +
				"    qwertyuiop[]\\\n" +
				"     asdfghjkl;'↲\n" +
				"shift zxcvbnm,./\n" +
				"        space",
		},
	}
	for _, tc := range testCases {
		got := qwerty.render(tc.char)
		if got != tc.want {
			t.Errorf("%s\ngot\n%s\nwant\n%s", tc.name, got, tc.want)
		}
	}
}
