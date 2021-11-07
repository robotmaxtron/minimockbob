package main

import (
	"os"
	"strings"

	"github.com/robotmaxtron/minimockbob"
)

func main() {
	userInput := strings.Join(os.Args[1:], " ")
	var sb strings.Builder
	sb.Grow(len(userInput))
	sb.WriteString(minimockbob.Gen(userInput))
	println(sb.String())

}
