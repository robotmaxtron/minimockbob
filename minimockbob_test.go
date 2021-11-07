package minimockbob

import (
	"testing"
)

func TestNoCaps(t *testing.T) {
	want := "fOoBaR"
	if got := Gen("foobar"); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}

func TestCaps(t *testing.T) {
	want := "fOoBaR"
	if got := Gen("Foobar"); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}

func TestMultipleWords(t *testing.T) {
	want := "fOo BaR"
	if got := Gen("foo bar"); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}

func TestPunctuation(t *testing.T) {
	want := "fOo, BaR?"
	if got := Gen("foo, bar?"); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
