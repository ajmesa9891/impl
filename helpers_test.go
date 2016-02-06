package impl

import (
	"testing"
)

func TestFormatInterface(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		isError bool
	}{
		{"io.Blinger", "io.Blinger", false},
		{"io.Reader", "io.Reader", false},
		{"imports.Madeuper", "imports.Madeuper", false},

		{"io..Reader", "", true},
	}
	for _, c := range cases {
		got, err := formatInterface(c.in)
		if c.isError && err == nil {
			t.Errorf("formatInterface(%q) ==  \"nil\", want \"error\"", c.in)
		} else if !c.isError && err != nil {
			t.Errorf("formatInterface(%q) == \"error\", want \"nil\": %q", c.in, err)
		} else if got != c.want {
			t.Errorf("formatInterface(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
