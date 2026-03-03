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
	"os"

	"github.com/robotmaxtron/minimockbob"
)

func main() {
	os.Exit(minimockbob.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
