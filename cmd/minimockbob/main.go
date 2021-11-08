/*commandline binary

Usage:
	./minimockbob <user input>

*/
package main

import (
	"os"
	"strings"

	"github.com/robotmaxtron/minimockbob"
)

func main() {
	userInput := strings.Join(os.Args[1:], " ")
	var sb strings.Builder
	if userInput == "" {
		sb.WriteString("Usage:\n	./minimockbob foo bar\n")
		println(sb.String())
		os.Exit(1)
	}
	sb.Grow(len(userInput))
	output, err := minimockbob.Gen(userInput)
	if err != nil {
		panic(err)
	}

	sb.WriteString(output)
	println(sb.String())

}
