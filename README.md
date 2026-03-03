# minimockbob

[![Go Reference](https://pkg.go.dev/badge/github.com/robotmaxtron/minimockbob.svg)](https://pkg.go.dev/github.com/robotmaxtron/minimockbob) 
[![Go Report Card](https://goreportcard.com/badge/github.com/robotmaxtron/minimockbob)](https://goreportcard.com/report/github.com/robotmaxtron/minimockbob)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

`minimockbob` is a sarcastic text generator written in Go that transforms strings into alternating capitalization 
(e.g., "hello world" -> "hElLo WoRlD"). 

It can be imported as a Go package, compiled to a binary, or run as a container.

## Features

- **Alternating Capitalization**: The classic "mocking" text format.
- **Multiple Interfaces**: Use it as a CLI tool, a Go library, or a containerized service.
- **Modern Tooling**: Optimized for [ko](https://ko.build) (containers) and [melange](https://github.com/chainguard-dev/melange) (Wolfi APK packages).
- **Minimal Footprint**: Distroless-based container images (~10MB) and stripped binaries.
- **Flexible Input**: Supports CLI arguments, quoted strings, and piped input.
- **Go 1.26+**: Built using modern Go idioms and toolchain.

## Table of Contents

- [Quick Start](#quick-start)
- [Installation](#installation)
  - [Install as a Binary](#install-as-a-binary)
  - [Build as a Container](#build-as-a-container)
  - [Build as a Wolfi APK Package](#build-as-a-wolfi-apk-package)
- [Usage](#usage)
  - [Command Line](#command-line)
  - [Container Usage](#container-usage)
  - [As a Go Package](#as-a-go-package)
- [Testing](#testing)
  - [Go Test Suite](#go-test-suite)
  - [Wolfi Package Testing](#wolfi-package-testing)
- [Configuration](#configuration)
- [Documentation](#documentation)
- [Words](#words)
- [License](#license)

## Quick Start

```bash
go run github.com/robotmaxtron/minimockbob/cmd/minimockbob@latest "mocking intensifies"
# mOcKiNg InTeNsIfIeS
```

## Installation

### Install as a Binary

Install with Go:
```bash
go install github.com/robotmaxtron/minimockbob/cmd/minimockbob@latest
```

Or build from source:
```bash
git clone https://github.com/robotmaxtron/minimockbob.git
cd minimockbob/cmd/minimockbob
go build
```

### Build as a Container

Build with [ko](https://ko.build):
```bash
# Set the target repository (use ko.local for local builds)
KO_DOCKER_REPO=ko.local ko build ./cmd/minimockbob

# Or build with OCI image labels and tags
KO_DOCKER_REPO=ko.local ko build \
  --image-label org.opencontainers.image.source=https://github.com/robotmaxtron/minimockbob \
  --image-label org.opencontainers.image.description="A sarcastic text generator" \
  --image-label org.opencontainers.image.licenses=Apache-2.0 \
  --tags latest,v0.0.1 \
  ./cmd/minimockbob
```

### Build as a Wolfi APK Package

Build with [melange](https://github.com/chainguard-dev/melange):
```bash
# Generate a signing key (first time only)
melange keygen

# Build the package
melange build .melange.yaml \
  --signing-key melange.rsa \
  --runner docker \
  --arch aarch64
```

To install in a Wolfi container:
```bash
# Run Wolfi container with your packages directory mounted
docker run -v $(pwd)/packages:/packages --rm -it cgr.dev/chainguard/wolfi-base sh

# Inside the container, install the package
apk add --allow-untrusted /packages/aarch64/minimockbob-0.0.1-r0.apk

# Test the binary
minimockbob "Hello Wolfi"
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

3. **Pipe input:**
   ```bash
   echo "Hello, World!" | minimockbob
   # Output: hElLo, WoRlD!
   ```

### Container Usage

Run the container built with ko:
```bash
# With arguments (use the image name from ko build output)
docker run --rm ko.local/minimockbob-<hash>:latest "Hello Container"

# With piped input
echo "Hello Container" | docker run --rm -i ko.local/minimockbob-<hash>:latest
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

### Go Test Suite

Run the full test suite:
```bash
go test ./...
```

Run CLI functional tests:
```bash
cd cmd/minimockbob && go test -v
```

Run tests with coverage:
```bash
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out
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

# Skip melange tests (faster local development)
go test -short ./...
```

**Note:** If using Docker alternatives like colima, rancher-desktop, or podman, the melange build tests may fail due to Docker mount path issues. In this case, run melange directly as shown in the Installation section.

## Configuration

The project includes configuration for building optimized container images and packages:

- **`.ko.yaml`**: Configuration for [ko](https://ko.build) container builds.
- **`.melange.yaml`**: Configuration for [melange](https://github.com/chainguard-dev/melange) APK builds.

Both configurations focus on minimal attack surface (distroless), CGO disabled, and reproducible builds.

## Documentation

View the full package documentation:
```bash
# View package documentation locally
go doc -all github.com/robotmaxtron/minimockbob

# Or use godoc for a local server
# go install golang.org/x/tools/cmd/godoc@latest
godoc -http=:6060
```

Online documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/robotmaxtron/minimockbob).

## Words

This package is a joke. Shoutout to [mockbob](https://github.com/tlkamp/mockbob). For my friend, James.
> "Once men turned their thinking over to machines in the hope that this would set them free. But that only permitted 
> other men with machines to enslave them."
> 
> — **Reverend Mother Gaius Helen Mohiam**, *Dune* by Frank Herbert

## License

This project is licensed under the Apache License 2.0 – see the [LICENCE](LICENCE) file for details.