package minimockbob

import (
	"testing"
)

func TestNoCaps(t *testing.T) {
	want := "fOoBaR"
	if got := sarcasmGen("foobar"); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}

func TestCaps(t *testing.T) {
	want := "fOoBaR"
	if got := sarcasmGen("Foobar"); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}

func TestMultipleWords(t *testing.T) {
	want := "fOo BaR"
	if got := sarcasmGen("foo bar"); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}

func TestPunctuation(t *testing.T) {
	want := "fOo, BaR?"
	if got := sarcasmGen("foo, bar?"); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
