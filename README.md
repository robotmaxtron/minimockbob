# minimockbob is a sarcastic text generator written in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/robotmaxtron/minimockbob.svg)](https://pkg.go.dev/github.com/robotmaxtron/minimockbob) [![Go Report Card](https://goreportcard.com/badge/github.com/robotmaxtron/minimockbob)](https://goreportcard.com/report/github.com/robotmaxtron/minimockbob)

`minimockbob` transforms a string into one with alternating capitalization. It can be imported as a Go package, compiled to a binary, or run as a container.

## Installation

### Install as a binary

Install with Go:
```bash
go install github.com/robotmaxtron/minimockbob/cmd/minimockbob@latest
```

Or build from source:
```bash
cd cmd/minimockbob
go build
```

### Build as a container

Build with [ko](https://ko.build):
```bash
# Set the target repository (use ko.local for local builds)
KO_DOCKER_REPO=ko.local ko build ./cmd/minimockbob

# Or push to a remote registry
KO_DOCKER_REPO=docker.io/yourusername ko build ./cmd/minimockbob
```

Or build with Docker:
```bash
docker build -t minimockbob .
```

## Usage

### Command Line

The binary supports three usage modes:

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

### Container Usage

Run with Docker:
```bash
# With arguments
docker run --rm ko.local/minimockbob:latest "Hello Container"

# With piped input
echo "Hello Container" | docker run --rm -i ko.local/minimockbob:latest
```

### As a Go Package

Import and use in your code:
```go
package main

import (
    "fmt"
    "github.com/robotmaxtron/minimockbob"
)

func main() {
    result := minimockbob.Gen("Hello, World!")
    fmt.Println(result)  // Output: hElLo, WoRlD!
}
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

Run tests with coverage:
```bash
go test -cover ./...
```

Generate detailed coverage report:
```bash
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out
```

Generate HTML coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

Run performance benchmarks:
```bash
go test -bench=. -benchmem
```

## Documentation

View the full package documentation:
```bash
# View package documentation locally
go doc -all github.com/robotmaxtron/minimockbob

# Or use godoc to start a local documentation server
godoc -http=:6060
# Then visit http://localhost:6060/pkg/github.com/robotmaxtron/minimockbob/
```

Online documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/robotmaxtron/minimockbob).

## Container Configuration

The project includes configuration for building optimized container images:

- **`.ko.yaml`**: Configuration for [ko](https://ko.build) builds
  - Uses distroless base image for minimal attack surface
  - CGO disabled for static binary
  - Trimpath for reproducible builds

- **`Dockerfile`**: Multi-stage Docker build
  - Alpine-based build stage
  - Distroless runtime for security
  - Results in ~9MB container image

### Dedication
For my friend, James.

### Inspiration
Shoutout to [mockbob](https://github.com/tlkamp/mockbob)