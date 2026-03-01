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
	var userInput string
	if len(os.Args) > 1 {
		userInput = strings.Join(os.Args[1:], " ")
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// Read from pipe
			b, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading from STDIN: %v\n", err)
				os.Exit(1)
			}
			userInput = strings.TrimSuffix(string(b), "\n")
		}
	}

	if userInput == "" {
		fmt.Println("Usage: minimockbob \"<text>\"")
		fmt.Println("Or pipe text to it: echo \"<text>\" | minimockbob")
		os.Exit(1)
	}
	output := minimockbob.Gen(userInput)
	fmt.Println(output)
}
