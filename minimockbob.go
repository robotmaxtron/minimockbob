// Package minimockbob is a sarcastic text generator that transforms a string into one with alternating capitalization.
//
// The package is primarily intended for use in testing and likely should not ever be used in production.
//
// # Overview
//
// The core functionality is provided by the Gen function, which transforms text by alternating
// between lowercase and uppercase for each letter, starting with lowercase. Non-letter characters
// (spaces, punctuation, numbers, etc.) are preserved as-is.
//
// # Package Usage
//
// Importing the package and calling the Gen function:
//
//	package main
//	import (
//		"fmt"
//		"github.com/robotmaxtron/minimockbob"
//	)
//	func main() {
//		result := minimockbob.Gen("Hello, World!")
//		fmt.Println(result)  // Output: hElLo, WoRlD!
//	}
//
// # Command Line Usage
//
// The package includes a command-line binary that can be used in three ways:
//
// Installation:
//
//	go install github.com/robotmaxtron/minimockbob/cmd/minimockbob@latest
//
// Or build from source:
//
//	cd cmd/minimockbob
//	go build
//
// Usage modes:
//
//	1. Quoted argument:
//	   minimockbob "Hello, World!"
//	   Output: hElLo, WoRlD!
//
//	2. Multiple unquoted arguments:
//	   minimockbob Hello World
//	   Output: hElLo WoRlD
//
//	3. Pipe input (no shell escaping required):
//	   echo "Hello, World!" | minimockbob
//	   Output: hElLo, WoRlD!
//
// # Container Usage
//
// The application can be built and run as a container using ko (https://ko.build):
//
//	# Build the container image locally
//	# ko creates images with a hash-based name and adds the 'latest' tag by default
//	KO_DOCKER_REPO=ko.local ko build ./cmd/minimockbob
//	# Output example: ko.local/minimockbob-8954540d57b52cc7913d1bf8fd346995:f491aaefe3ba5cfce00027e319ee1968d734a446b92345351abca33707545ff4
//
//	# Or push to a remote registry
//	KO_DOCKER_REPO=docker.io/yourusername ko build ./cmd/minimockbob
//
//	# Build with specific tags (hash-based name will still be created, with your custom tags added)
//	KO_DOCKER_REPO=ko.local ko build --tags latest,v0.0.1 ./cmd/minimockbob
//
// Running the container:
//
//	# Run with arguments (use the full image name from ko build output)
//	docker run --rm ko.local/minimockbob-8954540d57b52cc7913d1bf8fd346995:latest "Hello Container"
//
//	# Run with piped input
//	echo "Hello Container" | docker run --rm -i ko.local/minimockbob-8954540d57b52cc7913d1bf8fd346995:latest
//
//	# Run from a remote registry
//	docker run --rm docker.io/yourusername/minimockbob:latest "Hello Container"
//
// The ko build uses the .ko.yaml configuration for optimized builds with distroless base images.
//
// # Implementation Details
//
// The Gen function iterates through each character (rune) in the input string:
//   - If the character is a letter, it alternates between lowercase and uppercase
//   - Non-letter characters are preserved exactly as provided
//   - The alternation state is maintained only for letters
//   - Supports Unicode characters correctly
//
// Performance considerations:
//   - Uses strings.Builder for efficient string construction
//   - Pre-allocates capacity to avoid reallocations during iteration
//   - Handles Unicode characters correctly using unicode.IsLetter
//
// # Examples
//
// Basic usage:
//
//	Gen("Hello, World!")  // Returns: "hElLo, WoRlD!"
//	Gen("foo bar")        // Returns: "fOo BaR"
//	Gen("test123")        // Returns: "tEsT123"
//	Gen("ABC")            // Returns: "aBc"
//
// Special characters:
//
//	Gen("hello!")         // Returns: "hElLo!"
//	Gen("a+b=c")          // Returns: "a+B=c"
//	Gen("test@example")   // Returns: "tEsT@eXaMpLe"
//
// Unicode support:
//
//	Gen("café")           // Returns: "cAfÉ"
//	Gen("Привет")         // Returns: "пРиВеТ"
//
// # Testing
//
// Run the comprehensive test suite:
//
//	go test ./...
//
// Run with coverage:
//
//	go test -cover ./...
//
// Run benchmarks:
//
//	go test -bench=. -benchmem
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
