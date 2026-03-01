# minimockbob is a sarcastic text generator written in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/robotmaxtron/minimockbob.svg)](https://pkg.go.dev/github.com/robotmaxtron/minimockbob) [![Go Report Card](https://goreportcard.com/badge/github.com/robotmaxtron/minimockbob)](https://goreportcard.com/report/github.com/robotmaxtron/minimockbob)

`minimockbob` transforms a string into one with alternating capitalization, it can be imported into your own package or 
compiled to a binary.

## Installation

Install the binary with Go:
```bash
go install github.com/robotmaxtron/minimockbob/cmd/minimockbob@latest
```

Or build from source:
```bash
cd cmd/minimockbob
go build
```

## Usage

The command line utility supports three usage modes:

1. **Quoted argument:**
   ```bash
   minimockbob "Hello, World!"
   # Output: hElLo, WoRlD!
   ```

2. **Multiple unquoted arguments:**
   ```bash
   minimockbob Hello World
   # Output: hElLo WoRlD
   ```

3. **Pipe input (no shell escaping required):**
   ```bash
   echo "Hello, World!" | minimockbob
   # Output: hElLo, WoRlD!
   ```

## Testing

Run the test suite:
```bash
go test ./...
```

Run CLI functional tests:
```bash
cd cmd/minimockbob
go test -v
```
