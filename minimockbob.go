package minimockbob

import (
	"errors"
	"strings"
	"unicode"
)

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
