package minimockbob

import (
	"testing"
)

func TestNoCaps(t *testing.T) {
	want := "fOoBaR"
	if got := minimockbob("foobar"); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}

func TestCaps(t *testing.T) {
	want := "fOoBaR"
	if got := minimockbob("Foobar"); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}

func TestMultipleWords(t *testing.T) {
	want := "fOo BaR"
	if got := minimockbob("foo bar"); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}

func TestPunctuation(t *testing.T) {
	want := "fOo, BaR?"
	if got := minimockbob("foo, bar?"); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
