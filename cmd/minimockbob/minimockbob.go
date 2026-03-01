/*
Command line binary for minimockbob - a sarcastic text generator.

Usage:

	./minimockbob "Hello, World!"     # Quoted argument
	./minimockbob Hello World         # Multiple unquoted arguments
	echo "Hello, World!" | minimockbob  # Pipe input

All three methods will output: hElLo, WoRlD!

The binary transforms input text into alternating capitalization,
starting with lowercase for the first letter.
*/
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/robotmaxtron/minimockbob"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
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
	output := minimockbob.Gen(userInput)
	_, _ = fmt.Fprintln(stdout, output)
	return 0
}
