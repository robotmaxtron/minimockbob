// Package minimockbob transforms a string into a new one with alternating capitalization.
package minimockbob

import (
	"strings"
	"unicode"
)

// Gen transforms a string into one with alternating capitalization.
// The result will always begin with a lowercase letter if the first character is a letter.
//
// The function iterates through each character (rune) in the input string:
//   - If the character is a letter, it alternates between lowercase and uppercase
//   - Non-letter characters (spaces, punctuation, numbers, etc.) are preserved as-is
//   - The alternation state is maintained only for letters
//
// Example usage:
//
//	Gen("Hello, World!")  // Returns: "hElLo, WoRlD!"
//	Gen("foo bar")        // Returns: "fOo BaR"
//	Gen("test123")        // Returns: "tEsT123"
//
// Performance considerations:
//   - Uses strings.Builder for efficient string construction
//   - Pre-allocates capacity to avoid reallocations
//   - Handles Unicode characters correctly using unicode.IsLetter
func Gen(input string) string {
	var b strings.Builder
	b.Grow(len(input)) // Pre-allocate capacity for efficiency
	upper := false     // Track alternation state (false = lowercase, true = uppercase)

	for _, c := range input {
		if unicode.IsLetter(c) {
			// For letters, alternate between lower and upper case
			if upper {
				b.WriteRune(unicode.ToUpper(c))
			} else {
				b.WriteRune(unicode.ToLower(c))
			}
			upper = !upper // Toggle state for next letter
		} else {
			// Preserve non-letter characters as-is
			b.WriteRune(c)
		}
	}
	return b.String()
}
