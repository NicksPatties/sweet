package root

import (
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
			"A", "S", "D", "F", "G", "H", "J", "K", "L", ":", "\"", "\n",
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

func (k keymap) render() (km string) {
	rows := len(k.keys)
	for ri, row := range k.keys {
		km += spaces(k.margins[ri])
		for _, char := range row {
			km += char
		}
		if ri != rows-1 {
			km += "\n"
		}
	}
	return
}

func TestRender(t *testing.T) {
	name := "show keys, no modifiers, no "
	want := "   `1234567890-=\n" +
		"    qwertyuiop[]\\\n" +
		"     asdfghjkl;'↲\n" +
		"shift zxcvbnm,./\n" +
		"        space"
	got := qwerty.render()

	if got != want {
		t.Errorf("%s\ngot\n%s\nwant\n%s", name, got, want)
	}

}
