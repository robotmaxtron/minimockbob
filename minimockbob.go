// Package minimockbob transforms a string into a new one with alternating capitalization.
package minimockbob

import (
	"fmt"
	"io"
	"os"
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

// Run is the main logic of the CLI, separated for testability.
// It accepts command-line arguments and I/O streams as parameters.
//
// The function supports three input modes:
//  1. Command-line arguments: args are joined with spaces
//  2. Piped input: reads from stdin when no args provided and stdin is a pipe
//  3. Empty input: shows usage information
//
// Returns exit code: 0 for success, 1 for error or usage display
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	var userInput string
	if len(args) > 0 {
		userInput = strings.Join(args, " ")
	} else {
		// Check if we're reading from a pipe
		if f, ok := stdin.(*os.File); ok {
			stat, _ := f.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				// Read from pipe
				b, err := io.ReadAll(stdin)
				if err != nil {
					_, _ = fmt.Fprintf(stderr, "Error reading from STDIN: %v\n", err)
					return 1
				}
				userInput = strings.TrimSuffix(string(b), "\n")
			}
		} else {
			// For testing with non-file readers
			b, err := io.ReadAll(stdin)
			if err != nil {
				_, _ = fmt.Fprintf(stderr, "Error reading from STDIN: %v\n", err)
				return 1
			}
			userInput = strings.TrimSuffix(string(b), "\n")
		}
	}

	if userInput == "" {
		_, _ = fmt.Fprintln(stdout, "Usage: minimockbob \"<text>\"")
		_, _ = fmt.Fprintln(stdout, "Or pipe text to it: echo \"<text>\" | minimockbob")
		return 1
	}
	output := Gen(userInput)
	_, _ = fmt.Fprintln(stdout, output)
	return 0
}
