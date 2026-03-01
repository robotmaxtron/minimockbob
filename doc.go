package minimockbob

/*
Package minimockbob is a sarcastic text generator that transforms a string into one with alternating capitalization.

The package is primarily intended for use in testing and likely should not ever be used in production.

The package can be used in two ways:

1. Importing the package and calling the function Gen with a string as an argument.

	package main
	import (
		"fmt"
		"github.com/robotmaxtron/minimockbob"
	)
	func main() {
		fmt.Println(minimockbob.Gen("Hello, World!"))
	}

	Output: hElLo, WoRlD!

2. Running the package as a command line utility.

Installation:

	go install github.com/robotmaxtron/minimockbob/cmd/minimockbob@latest

Or build from source:

	cd cmd/minimockbob
	go build

Usage examples:

	minimockbob "Hello, World!"          # Quoted argument
	minimockbob Hello World              # Multiple unquoted arguments
	echo "Hello, World!" | minimockbob   # Pipe input

	All three methods output: hElLo, WoRlD!

The Gen function transforms text by alternating between lowercase and uppercase for each letter,
starting with lowercase. Non-letter characters are preserved as-is.
*/
