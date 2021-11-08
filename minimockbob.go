// package minimockbob transforms a string of length of at least 2 letters into a new one with alternating capitalization.
package minimockbob

import (
	"errors"
	"strings"
	"unicode"
)

// Gen transforms a string into one with alternating capitalization.
// If the string length is less than 2, an empty string and an error are returned.
// result will always begin with a lowercase letter.
func Gen(input string) (string, error) {
	if len(input) < 2 {
		return "", errors.New("input is too short to process")
	}
	var b strings.Builder
	b.Grow(len(input))
	upper := true
	for _, c := range input {
		if unicode.IsLetter(c) {
			upper = !upper
			if upper {
				b.WriteRune(unicode.ToUpper(c))
			} else {
				b.WriteRune(unicode.ToLower(c))
			}
		} else {
			b.WriteRune(c)
		}
	}
	return b.String(), nil
}
