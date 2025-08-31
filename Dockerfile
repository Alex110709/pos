# PIXELZX POS EVM Chain Dockerfile
# Multi-stage build for production optimization with multi-architecture support

# Build arguments for multi-platform support
ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

# Build stage with cross-compilation support
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Set cross-compilation environment variables
ENV CGO_ENABLED=0
ARG TARGETOS
ARG TARGETARCH
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

# Build the binary with cross-compilation support
RUN echo "Building for platform: ${TARGETPLATFORM} (OS: ${TARGETOS}, ARCH: ${TARGETARCH})" && \
    go build -ldflags "-s -w -X main.version=$(git describe --tags --always --dirty)" \
    -o bin/pixelzx ./cmd/pixelzx

# Production stage with target platform
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create app user for security with explicit UID/GID
RUN addgroup -g 1001 -S pixelzx && \
    adduser -u 1001 -S pixelzx -G pixelzx

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/bin/pixelzx /usr/local/bin/pixelzx

# Make binary executable
RUN chmod +x /usr/local/bin/pixelzx

# Create data directories with proper permissions
RUN mkdir -p /app/data /app/keystore /app/logs && \
    chown -R 1001:1001 /app && \
    chmod -R 755 /app

# Copy configuration files with proper ownership
COPY --chown=1001:1001 configs/production.yaml /app/config.yaml

# Ensure pixelzx user can write to necessary directories
RUN chown -R pixelzx:pixelzx /app/data /app/keystore /app/logs

# Switch to app user
USER pixelzx

# Expose ports
EXPOSE 8545 8546 30303

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD /usr/local/bin/pixelzx version || exit 1

# Set environment variables
ENV PIXELZX_HOME=/app
ENV PIXELZX_CONFIG=/app/config.yaml
ENV PIXELZX_DATA_DIR=/app/data
ENV PIXELZX_PLATFORM=${TARGETPLATFORM}

# Display platform information
ARG TARGETPLATFORM
RUN echo "Container platform: ${TARGETPLATFORM}" && \
    echo "Binary architecture: $(file /usr/local/bin/pixelzx)"

# Default command
CMD ["/usr/local/bin/pixelzx", "start", "--config", "/app/config.yaml", "--datadir", "/app/data"]