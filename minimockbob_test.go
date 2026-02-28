package minimockbob

import (
	"testing"
)

func TestGen(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"foobar", "fOoBaR"},
		{"Foobar", "fOoBaR"},
		{"foo bar", "fOo BaR"},
		{"foo, bar?", "fOo, BaR?"},
		{"f", "f"},
		{"", ""},
		{"AB", "aB"},
		{"123", "123"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Gen(tt.input); got != tt.want {
				t.Errorf("Gen(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
