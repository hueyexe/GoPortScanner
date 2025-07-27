# Multi-stage build for GoPortScanner
# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o goportscanner cmd/goportscanner/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S goportscanner && \
    adduser -u 1001 -S goportscanner -G goportscanner

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/goportscanner .

# Change ownership to non-root user
RUN chown -R goportscanner:goportscanner /app

# Switch to non-root user
USER goportscanner

# Expose port (for documentation, not actually used by the scanner)
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["./goportscanner"]

# Default command (can be overridden)
CMD ["--help"]