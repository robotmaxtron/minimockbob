// Package minimockbob transforms a string into a new one with alternating capitalization.
package minimockbob

import (
	"strings"
	"unicode"
)

// Gen transforms a string into one with alternating capitalization.
// The result will always begin with a lowercase letter if the first character is a letter.
func Gen(input string) string {
	var b strings.Builder
	b.Grow(len(input))
	upper := false
	for _, c := range input {
		if unicode.IsLetter(c) {
			if upper {
				b.WriteRune(unicode.ToUpper(c))
			} else {
				b.WriteRune(unicode.ToLower(c))
			}
			upper = !upper
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}
