package minimockbob

/*
Package minimockbob is a sarcastic text generator, transforming a string into one with alternating capitalization.

The package is primarily intended for use in testing, and likely should not ever be used in production.

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

		Will print: hElLo, WoRlD!



	2. Running the package as a command line utility. To build the binary, run `go build` in the cmd/minimockbobcli subdirectory.

		Usage: ./minimockbob Hello, World!

		Will print: hElLo, WoRlD!
*/
