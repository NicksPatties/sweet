package root

import (
	lg "github.com/charmbracelet/lipgloss"
	"strconv"
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

// Returns a view of the fingers for the keymap.
// margin is the spacing on the left to push
// fi is finger icon (which will actually be determined by)
// c is character to type
func fingerView(margin int, fIcon rune, currChar rune) (view string) {

	// rune to fingers
	// doing this _for each character change_ a good idea?
	rtfs := make(map[rune][]uint)
	rtfs['`'] = []uint{lpinky}
	rtfs['1'] = []uint{lpinky}
	rtfs['q'] = []uint{lpinky}
	rtfs['a'] = []uint{lpinky}
	rtfs['z'] = []uint{lpinky}
	rtfs['~'] = []uint{lpinky, rpinky}
	rtfs['!'] = []uint{lpinky, rpinky}
	rtfs['Q'] = []uint{lpinky, rpinky}
	rtfs['A'] = []uint{lpinky, rpinky}
	rtfs['Z'] = []uint{lpinky, rpinky}

	rtfs['2'] = []uint{lring}
	rtfs['w'] = []uint{lring}
	rtfs['s'] = []uint{lring}
	rtfs['x'] = []uint{lring}
	rtfs['@'] = []uint{lring, rpinky}
	rtfs['W'] = []uint{lring, rpinky}
	rtfs['S'] = []uint{lring, rpinky}
	rtfs['X'] = []uint{lring, rpinky}

	rtfs['3'] = []uint{lmiddle}
	rtfs['e'] = []uint{lmiddle}
	rtfs['d'] = []uint{lmiddle}
	rtfs['c'] = []uint{lmiddle}
	rtfs['#'] = []uint{lmiddle, rpinky}
	rtfs['E'] = []uint{lmiddle, rpinky}
	rtfs['D'] = []uint{lmiddle, rpinky}
	rtfs['C'] = []uint{lmiddle, rpinky}

	rtfs['4'] = []uint{lindex}
	rtfs['r'] = []uint{lindex}
	rtfs['f'] = []uint{lindex}
	rtfs['v'] = []uint{lindex}
	rtfs['$'] = []uint{lindex, rpinky}
	rtfs['R'] = []uint{lindex, rpinky}
	rtfs['F'] = []uint{lindex, rpinky}
	rtfs['V'] = []uint{lindex, rpinky}
	rtfs['5'] = []uint{lindex}
	rtfs['t'] = []uint{lindex}
	rtfs['g'] = []uint{lindex}
	rtfs['b'] = []uint{lindex}
	rtfs['%'] = []uint{lindex, rpinky}
	rtfs['T'] = []uint{lindex, rpinky}
	rtfs['G'] = []uint{lindex, rpinky}
	rtfs['B'] = []uint{lindex, rpinky}

	rtfs[' '] = []uint{lthumb}

	rtfs['6'] = []uint{rindex}
	rtfs['y'] = []uint{rindex}
	rtfs['h'] = []uint{rindex}
	rtfs['n'] = []uint{rindex}
	rtfs['^'] = []uint{rindex, lpinky}
	rtfs['Y'] = []uint{rindex, lpinky}
	rtfs['H'] = []uint{rindex, lpinky}
	rtfs['N'] = []uint{rindex, lpinky}
	rtfs['7'] = []uint{rindex}
	rtfs['u'] = []uint{rindex}
	rtfs['j'] = []uint{rindex}
	rtfs['m'] = []uint{rindex}
	rtfs['&'] = []uint{rindex, lpinky}
	rtfs['U'] = []uint{rindex, lpinky}
	rtfs['J'] = []uint{rindex, lpinky}
	rtfs['M'] = []uint{rindex, lpinky}

	rtfs['8'] = []uint{rmiddle}
	rtfs['i'] = []uint{rmiddle}
	rtfs['k'] = []uint{rmiddle}
	rtfs[','] = []uint{rmiddle}
	rtfs['*'] = []uint{rmiddle, lpinky}
	rtfs['I'] = []uint{rmiddle, lpinky}
	rtfs['K'] = []uint{rmiddle, lpinky}
	rtfs['<'] = []uint{rmiddle, lpinky}

	rtfs['9'] = []uint{rring}
	rtfs['o'] = []uint{rring}
	rtfs['l'] = []uint{rring}
	rtfs['.'] = []uint{rring}
	rtfs['('] = []uint{rring, lpinky}
	rtfs['O'] = []uint{rring, lpinky}
	rtfs['L'] = []uint{rring, lpinky}
	rtfs['>'] = []uint{rring, lpinky}

	rtfs['0'] = []uint{rpinky}
	rtfs['p'] = []uint{rpinky}
	rtfs[';'] = []uint{rpinky}
	rtfs['/'] = []uint{rpinky}
	rtfs[')'] = []uint{rpinky, lpinky}
	rtfs['P'] = []uint{rpinky, lpinky}
	rtfs[':'] = []uint{rpinky, lpinky}
	rtfs['?'] = []uint{rpinky, lpinky}

	rtfs['-'] = []uint{rpinky}
	rtfs['['] = []uint{rpinky}
	rtfs['\''] = []uint{rpinky}
	rtfs['='] = []uint{rpinky}
	rtfs[']'] = []uint{rpinky}
	rtfs['\n'] = []uint{rpinky}
	rtfs['\\'] = []uint{rpinky}
	rtfs['_'] = []uint{rpinky, lpinky}
	rtfs['{'] = []uint{rpinky, lpinky}
	rtfs['"'] = []uint{rpinky, lpinky}
	rtfs['+'] = []uint{rpinky, lpinky}
	rtfs['}'] = []uint{rpinky, lpinky}
	rtfs['|'] = []uint{rpinky, lpinky}

	// fingers
	f := [][]rune{
		{'0', fIcon},
		{'1', fIcon, fIcon},
		{'2', fIcon, fIcon},
		{'3', fIcon},
		{'4'},
		{'5'},
		{'6', fIcon},
		{'7', fIcon, fIcon},
		{'8', fIcon, fIcon},
		{'9', fIcon},
	}

	activeFingers := rtfs[currChar]

	for row := 2; row >= 0; row-- {
		for space := 0; space < margin; space++ {
			view += " "
		}
		// curr fingers
		for cf := 0; cf < len(f); cf++ {
			// if this is a finger spot, then I should print the character
			// in the finger view location
			if isFingerSpot := row < len(f[cf]); isFingerSpot {
				style := lg.NewStyle()

				isActive := false
				for _, finger := range activeFingers {
					if finger == uint(cf) {
						isActive = true
					}
				}
				if isActive {
					style = style.Reverse(true)
				}
				if row == 0 {
					view += style.Render(strconv.Itoa(cf))
				} else {
					view += style.Render(string(fIcon))
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
