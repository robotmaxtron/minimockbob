package minimockbob

import (
	"testing"
)

func TestNoCaps(t *testing.T) {
	want := "fOoBaR"
	if got, err := Gen("foobar"); got != want && err != nil {
		t.Errorf("Gen() = %q, want %q", got, want)
	}
}

func TestCaps(t *testing.T) {
	want := "fOoBaR"
	if got, err := Gen("Foobar"); got != want || err != nil {
		t.Errorf("Gen() = %q, want %q", got, want)
	}
}

func TestMultipleWords(t *testing.T) {
	want := "fOo BaR"
	if got, err := Gen("foo bar"); got != want || err != nil {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}

func TestPunctuation(t *testing.T) {
	want := "fOo, BaR?"
	if got, err := Gen("foo, bar?"); got != want || err != nil {
		t.Errorf("Gen() = %q, want %q", got, want)
	}
}

func TestTooShort(t *testing.T) {
	want := ""
	if got, err := Gen("f"); got != want {
		t.Errorf("Gen() = %q, want %q", got, want)
		t.Errorf("Error: %v", err)
	}
}
