# minimockbob is a sarcastic text generator written in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/robotmaxtron/minimockbob.svg)](https://pkg.go.dev/github.com/robotmaxtron/minimockbob) [![Go Report Card](https://goreportcard.com/badge/github.com/robotmaxtron/minimockbob)](https://goreportcard.com/report/github.com/robotmaxtron/minimockbob)

`minimockbob` transforms a string into one with alternating capitalization. It can be imported as a Go package, 
compiled to a binary, or run as a container.

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

# Build with specific tags
KO_DOCKER_REPO=docker.io/yourusername ko build --tags latest,v0.0.1 ./cmd/minimockbob
```

### Build as a Wolfi APK package

Build with [melange](https://github.com/chainguard-dev/melange):
```bash
# Generate a signing key (first time only)
melange keygen

# Build the package
melange build melange.yaml \
  --signing-key melange.rsa \
  --runner docker \
  --arch aarch64

```

Install the package in a Wolfi container, ARM example below:
```bash
# Run Wolfi container with your packages directory mounted
docker run -v $(pwd)/packages:/packages --rm -it cgr.dev/chainguard/wolfi-base sh

# Inside the container, install the package
apk add --allow-untrusted /packages/aarch64/minimockbob-0.0.1-r0.apk

# Test the binary
minimockbob "Hello Wolfi"
hElLo WoLfI
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

Run the container built with ko:
```bash
# With arguments
docker run --rm ko.local/minimockbob:latest "Hello Container"

# With piped input
echo "Hello Container" | docker run --rm -i ko.local/minimockbob:latest

# Run from a remote registry
docker run --rm docker.io/yourusername/minimockbob:latest "Hello Container"
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

### Wolfi Package Testing

Test the Wolfi APK package build and functionality:
```bash
# Run all melange tests (requires melange and Docker)
go test -v -run TestMelange

# Run individual test suites
go test -v -run TestMelangePrerequisites  # Check prerequisites
go test -v -run TestMelangeYAMLValid      # Validate melange.yaml
go test -v -run TestMelangeBuild          # Build APK package
go test -v -run TestPackageArtifacts      # Verify build artifacts
go test -v -run TestPackageContents       # Check package contents
go test -v -run TestPackageInstallation   # Test installation
go test -v -run TestBinaryExecution       # Test binary execution

# Skip melange tests (they take time and require Docker)
go test -short ./...
```

The test suite validates:
- Prerequisites (melange and Docker availability)
- melange.yaml configuration syntax
- APK package build process
- Package artifacts (main package, index, subpackages)
- Package contents (binary, documentation, license)
- Package installation in Wolfi container
- Binary execution with arguments and piped input

**Note:** If using Docker alternatives like colima, rancher-desktop, or podman, the melange build tests may fail due to 
Docker mount path issues. In this case, run melange directly:
```bash
# Generate signing key if needed
melange keygen

# Build package
melange build melange.yaml --signing-key melange.rsa --runner docker --arch aarch64
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

## Container and Package Configuration

The project includes configuration for building optimized container images and packages:

- **`.ko.yaml`**: Configuration for [ko](https://ko.build) container builds
  - Uses distroless base image for minimal attack surface
  - CGO disabled for static binary
  - Trimpath for reproducible builds
  - Optimized ldflags for smaller binary size
  - Results in minimal container image (~10MB)

- **`melange.yaml`**: Configuration for [melange](https://github.com/chainguard-dev/melange) APK builds
  - Builds optimized APK packages for Wolfi OS
  - Supports x86_64 and aarch64 architectures
  - Includes stripped binaries with size optimization
  - Creates main package and documentation subpackage

### Dedication
For my friend, James.

### Inspiration
Shoutout to [mockbob](https://github.com/tlkamp/mockbob)