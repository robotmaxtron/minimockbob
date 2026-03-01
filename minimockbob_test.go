package minimockbob

import (
	"testing"
)

func TestGen(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Basic cases
		{"Hello, World!", "hElLo, WoRlD!"},
		{"foobar", "fOoBaR"},
		{"Foobar", "fOoBaR"},
		{"foo bar", "fOo BaR"},
		{"foo, bar?", "fOo, BaR?"},
		{"f", "f"},
		{"", ""},
		{"AB", "aB"},
		{"123", "123"},

		// Edge cases
		{"a", "a"},
		{"A", "a"},
		{"aB", "aB"},
		{"Ab", "aB"},

		// Multiple spaces
		{"a  b  c", "a  B  c"},
		{"   leading", "   lEaDiNg"},
		{"trailing   ", "tRaIlInG   "},

		// Special characters
		{"!@#$%^&*()", "!@#$%^&*()"},
		{"test!@#", "tEsT!@#"},
		{"a-b_c", "a-B_c"},
		{"x[y]z", "x[Y]z"},
		{"m&n*o", "m&N*o"},
		{"a:b;c", "a:B;c"},
		{"p+q=r", "p+Q=r"},
		{"s/t\\u", "s/T\\u"},
		{"v<w>x", "v<W>x"},
		{"y{z}a", "y{Z}a"},
		{"test'quote", "tEsT'qUoTe"},
		{"test\"quote", "tEsT\"qUoTe"},

		// Numbers mixed with letters
		{"abc123def456ghi", "aBc123DeF456gHi"},
		{"123abc", "123aBc"},
		{"abc123", "aBc123"},
		{"1a2b3c", "1a2B3c"},

		// Punctuation
		{"hello!", "hElLo!"},
		{"hello?", "hElLo?"},
		{"hello.", "hElLo."},
		{"hello,world", "hElLo,WoRlD"},
		{"a.b.c.d", "a.B.c.D"},

		// Unicode and international characters
		{"café", "cAfÉ"},
		{"naïve", "nAïVe"},
		{"Ñoño", "ñOñO"},

		// Newlines and tabs
		{"hello\nworld", "hElLo\nWoRlD"},
		{"hello\tworld", "hElLo\tWoRlD"},
		{"a\nb\nc", "a\nB\nc"},

		// Long strings
		{"abcdefghijklmnopqrstuvwxyz", "aBcDeFgHiJkLmNoPqRsTuVwXyZ"},
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZ", "aBcDeFgHiJkLmNoPqRsTuVwXyZ"},

		// Mixed case scenarios
		{"HeLLo WoRLd", "hElLo WoRlD"},
		{"ALLCAPS", "aLlCaPs"},
		{"alllower", "aLlLoWeR"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Gen(tt.input); got != tt.want {
				t.Errorf("Gen(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestGenConsistency ensures Gen produces consistent results
func TestGenConsistency(t *testing.T) {
	input := "Hello, World!"
	want := "hElLo, WoRlD!"

	// Run multiple times to ensure consistency
	for i := 0; i < 100; i++ {
		got := Gen(input)
		if got != want {
			t.Errorf("Gen(%q) iteration %d = %q, want %q", input, i, got, want)
		}
	}
}

// TestGenAlternation ensures a proper alternation pattern
func TestGenAlternation(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"even letters", "abcd", "aBcD"},
		{"odd letters", "abcde", "aBcDe"},
		{"with non-letters", "a1b2c3d", "a1B2c3D"},
		{"letters at start and end", "xabc123defx", "xAbC123dEfX"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Gen(tt.input)
			if got != tt.want {
				t.Errorf("Gen(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestGenStartsWithLower ensures the first letter is always lowercase
func TestGenStartsWithLower(t *testing.T) {
	tests := []string{
		"Hello",
		"HELLO",
		"hELLO",
		"A",
		"Z",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			got := Gen(input)
			if got == "" {
				t.Fatalf("Gen(%q) returned empty string", input)
			}
			firstChar := rune(got[0])
			if firstChar >= 'A' && firstChar <= 'Z' {
				t.Errorf("Gen(%q) = %q, first letter should be lowercase", input, got)
			}
		})
	}
}
