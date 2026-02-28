/*
command line binary

Usage:

	./minimockbob "<user input>"
	Or pipe text to it: echo \"<text>\" | minimockbob
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
