# Multi-stage Dockerfile for minimockbob
# This Dockerfile is provided for reference, but ko can build without it

# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o minimockbob \
    ./cmd/minimockbob

# Final stage - use distroless for minimal attack surface
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copy the binary from builder
COPY --from=builder /build/minimockbob /minimockbob

# Use non-root user
USER nonroot:nonroot

# Set the entrypoint
ENTRYPOINT ["/minimockbob"]
