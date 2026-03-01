package minimockbob

/*
Package minimockbob is a sarcastic text generator that transforms a string into one with alternating capitalization.

The package is primarily intended for use in testing and likely should not ever be used in production.

# Overview

The core functionality is provided by the Gen function, which transforms text by alternating
between lowercase and uppercase for each letter, starting with lowercase. Non-letter characters
(spaces, punctuation, numbers, etc.) are preserved as-is.

# Package Usage

Importing the package and calling the Gen function:

	package main
	import (
		"fmt"
		"github.com/robotmaxtron/minimockbob"
	)
	func main() {
		result := minimockbob.Gen("Hello, World!")
		fmt.Println(result)  // Output: hElLo, WoRlD!
	}

# Command Line Usage

The package includes a command-line binary that can be used in three ways:

Installation:

	go install github.com/robotmaxtron/minimockbob/cmd/minimockbob@latest

Or build from source:

	cd cmd/minimockbob
	go build

Usage modes:

	1. Quoted argument:
	   minimockbob "Hello, World!"
	   Output: hElLo, WoRlD!

	2. Multiple unquoted arguments:
	   minimockbob Hello World
	   Output: hElLo WoRlD

	3. Pipe input (no shell escaping required):
	   echo "Hello, World!" | minimockbob
	   Output: hElLo, WoRlD!

# Container Usage

The application can be built and run as a container using ko (https://ko.build):

	# Build the container image
	KO_DOCKER_REPO=ko.local ko build ./cmd/minimockbob

	# Run with arguments
	docker run --rm ko.local/minimockbob:latest "Hello Container"

	# Run with piped input
	echo "Hello Container" | docker run --rm -i ko.local/minimockbob:latest

Alternatively, build with Docker:

	docker build -t minimockbob .
	docker run --rm minimockbob "Hello Docker"

# Implementation Details

The Gen function iterates through each character (rune) in the input string:
  - If the character is a letter, it alternates between lowercase and uppercase
  - Non-letter characters are preserved exactly as provided
  - The alternation state is maintained only for letters
  - Supports Unicode characters correctly

Performance considerations:
  - Uses strings.Builder for efficient string construction
  - Pre-allocates capacity to avoid reallocations during iteration
  - Handles Unicode characters correctly using unicode.IsLetter

# Examples

Basic usage:

	Gen("Hello, World!")  // Returns: "hElLo, WoRlD!"
	Gen("foo bar")        // Returns: "fOo BaR"
	Gen("test123")        // Returns: "tEsT123"
	Gen("ABC")            // Returns: "aBc"

Special characters:

	Gen("hello!")         // Returns: "hElLo!"
	Gen("a+b=c")          // Returns: "a+B=c"
	Gen("test@example")   // Returns: "tEsT@eXaMpLe"

Unicode support:

	Gen("café")           // Returns: "cAfÉ"
	Gen("Привет")         // Returns: "пРиВеТ"

# Testing

Run the comprehensive test suite:

	go test ./...

Run with coverage:

	go test -cover ./...

Run benchmarks:

	go test -bench=. -benchmem
*/
