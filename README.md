# minimockbob is a sarcastic text generator written in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/robotmaxtron/minimockbob.svg)](https://pkg.go.dev/github.com/robotmaxtron/minimockbob) [![Go Report Card](https://goreportcard.com/badge/github.com/robotmaxtron/minimockbob)](https://goreportcard.com/report/github.com/robotmaxtron/minimockbob) [![Go](https://github.com/robotmaxtron/minimockbob/actions/workflows/go.yml/badge.svg)](https://github.com/robotmaxtron/minimockbob/actions/workflows/go.yml)

`minimockbob` transforms a string into one with alternating capitalization, it can be imported into your own package or 
compiled to a binary.

Usage examples (no shell escaping required when piping):

- As arguments (quotes optional):
  - `minimockbob "Hello, World!"`
  - `minimockbob Hello World`
- Via STDIN (recommended for complex strings/special characters):
  - `echo "Hello, World!" | minimockbob`
  - `printf 'special * ! ? and $(vars) stay literal' | minimockbob`

![compiled](https://raw.githubusercontent.com/robotmaxtron/minimockbob/main/cmd/demo-tape/demo.gif)

Install the binary with Go `go install github.com/robotmaxtron/minimockbob/cmd/minimockbob@v0.0.3`
