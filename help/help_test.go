package help

import "testing"

func TestGetMessage(t *testing.T) {
	want := "hey from help"
	got := GetMessage()
	if want != got {
		t.Errorf("got: %s, wanted: %s", got, want)
	}
}
