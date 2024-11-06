package root

import (
	lg "github.com/charmbracelet/lipgloss"
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
	spaces := func(n int) (s string) {
		for i := 0; i < n; i++ {
			s += " "
		}
		return
	}

	// Highlighted key style
	var hk = lg.NewStyle().Reverse(true).Bold(true)
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
				km += hk.Render(key)
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
