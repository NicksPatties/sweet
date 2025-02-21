package root

import (
	lg "github.com/charmbracelet/lipgloss"
	"strconv"
)

type keymap struct {
	keys          [][]string
	modifiedKeys  [][]string
	margins       []int
	fingersMargin int
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
	margins:       []int{3, 4, 5, 0, 8},
	fingersMargin: 5,
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

// Rune to fingers
var rtfs = map[rune][]uint{
	'`': {lpinky},
	'1': {lpinky},
	'q': {lpinky},
	'a': {lpinky},
	'z': {lpinky},
	'~': {lpinky, rpinky},
	'!': {lpinky, rpinky},
	'Q': {lpinky, rpinky},
	'A': {lpinky, rpinky},
	'Z': {lpinky, rpinky},

	'2': {lring},
	'w': {lring},
	's': {lring},
	'x': {lring},
	'@': {lring, rpinky},
	'W': {lring, rpinky},
	'S': {lring, rpinky},
	'X': {lring, rpinky},

	'3': {lmiddle},
	'e': {lmiddle},
	'd': {lmiddle},
	'c': {lmiddle},
	'#': {lmiddle, rpinky},
	'E': {lmiddle, rpinky},
	'D': {lmiddle, rpinky},
	'C': {lmiddle, rpinky},

	'4': {lindex},
	'r': {lindex},
	'f': {lindex},
	'v': {lindex},
	'$': {lindex, rpinky},
	'R': {lindex, rpinky},
	'F': {lindex, rpinky},
	'V': {lindex, rpinky},
	'5': {lindex},
	't': {lindex},
	'g': {lindex},
	'b': {lindex},
	'%': {lindex, rpinky},
	'T': {lindex, rpinky},
	'G': {lindex, rpinky},
	'B': {lindex, rpinky},

	' ': {lthumb},

	'6': {rindex},
	'y': {rindex},
	'h': {rindex},
	'n': {rindex},
	'^': {rindex, lpinky},
	'Y': {rindex, lpinky},
	'H': {rindex, lpinky},
	'N': {rindex, lpinky},
	'7': {rindex},
	'u': {rindex},
	'j': {rindex},
	'm': {rindex},
	'&': {rindex, lpinky},
	'U': {rindex, lpinky},
	'J': {rindex, lpinky},
	'M': {rindex, lpinky},

	'8': {rmiddle},
	'i': {rmiddle},
	'k': {rmiddle},
	',': {rmiddle},
	'*': {rmiddle, lpinky},
	'I': {rmiddle, lpinky},
	'K': {rmiddle, lpinky},
	'<': {rmiddle, lpinky},

	'9': {rring},
	'o': {rring},
	'l': {rring},
	'.': {rring},
	'(': {rring, lpinky},
	'O': {rring, lpinky},
	'L': {rring, lpinky},
	'>': {rring, lpinky},

	'0': {rpinky},
	'p': {rpinky},
	';': {rpinky},
	'/': {rpinky},
	')': {rpinky, lpinky},
	'P': {rpinky, lpinky},
	':': {rpinky, lpinky},
	'?': {rpinky, lpinky},

	'-':  {rpinky},
	'[':  {rpinky},
	'\'': {rpinky},
	']':  {rpinky},
	'\n': {rpinky},
	'\\': {rpinky},
	'_':  {rpinky, lpinky},
	'{':  {rpinky, lpinky},
	'"':  {rpinky, lpinky},
	'+':  {rpinky, lpinky},
	'}':  {rpinky, lpinky},
	'|':  {rpinky, lpinky},
}

// Returns a view of the fingers for the keymap.
// margin is the spacing on the left to push
// fi is finger icon (which will actually be determined by)
// c is character to type
func renderFingers(margin int, fIcon rune, currChar rune) (view string) {

	// fingers view
	f := [][]rune{
		{rune(lpinky), fIcon},
		{rune(lring), fIcon, fIcon},
		{rune(lmiddle), fIcon, fIcon},
		{rune(lindex), fIcon},
		{rune(lthumb)},
		{rune(rthumb)},
		{rune(rindex), fIcon},
		{rune(rmiddle), fIcon, fIcon},
		{rune(rring), fIcon, fIcon},
		{rune(rpinky), fIcon},
	}

	activeFingers := rtfs[currChar]

	for row := len(f[lring]) - 1; row >= 0; row-- {
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
