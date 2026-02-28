/*
command line binary

Usage:

	./minimockbob <user input>
*/
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/robotmaxtron/minimockbob"
)

func main() {
	userInput := strings.Join(os.Args[1:], " ")
	if userInput == "" {
		fmt.Println("Usage: minimockbob <text>")
		os.Exit(1)
	}
	output := minimockbob.Gen(userInput)
	fmt.Println(output)
}
