package minimockbob

import (
	"strings"
	"unicode"
)

func sarcasmGen(input string) string {
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
	return b.String()
}
