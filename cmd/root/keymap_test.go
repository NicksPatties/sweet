package root

import "testing"

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
