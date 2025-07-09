# Universal Dockerfile for all burrowctl server examples
# Usage: docker build --build-arg EXAMPLE_DIR=basic -t burrowctl-server-basic .
# Usage: docker build --build-arg EXAMPLE_DIR=advanced -t burrowctl-server-advanced .
# Usage: docker build --build-arg EXAMPLE_DIR=advanced/cache-server -t burrowctl-server-cache .

# Multi-stage build for optimal image size
FROM golang:1.22-alpine AS builder

# Install git for go mod download
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files from root
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build argument to specify which example to build
ARG EXAMPLE_DIR=basic
ENV EXAMPLE_DIR=${EXAMPLE_DIR}

# Build the specific example
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./examples/server/${EXAMPLE_DIR}/main.go

# Production stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S burrowctl && \
    adduser -S -D -H -u 1001 -h /app -s /sbin/nologin -G burrowctl burrowctl

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Change ownership to non-root user
RUN chown burrowctl:burrowctl /app/server

# Switch to non-root user
USER burrowctl

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Default command
CMD ["./server"]