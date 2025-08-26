# PIXELZX POS EVM Chain Dockerfile
# Multi-stage build for production optimization

# Build stage
FROM golang:1.21-alpine AS builder

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

# Build the binary
RUN make build

# Production stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create app user for security
RUN addgroup -g 1001 -S pixelzx && \
    adduser -u 1001 -S pixelzx -G pixelzx

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/bin/pixelzx /usr/local/bin/pixelzx

# Create data directories
RUN mkdir -p /app/data /app/keystore /app/logs && \
    chown -R pixelzx:pixelzx /app

# Copy configuration files
COPY --chown=pixelzx:pixelzx configs/production.yaml /app/config.yaml

# Switch to app user
USER pixelzx

# Expose ports
EXPOSE 8545 8546 30303

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD pixelzx version || exit 1

# Set environment variables
ENV PIXELZX_HOME=/app
ENV PIXELZX_CONFIG=/app/config.yaml
ENV PIXELZX_DATA_DIR=/app/data

# Default command
CMD ["pixelzx", "start", "--config", "/app/config.yaml", "--datadir", "/app/data"]