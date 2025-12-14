# Multi-stage build for Go application
FROM golang:1.25-alpine AS builder

# Install build dependencies (добавляем gcc и musl-dev для CGO)
RUN apk add --no-cache git ca-certificates tzdata gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o /app/bin/spell_bot ./cmd/bot

# Final stage
FROM alpine:latest

# Install runtime dependencies (libgcc для sqlite3)
RUN apk --no-cache add ca-certificates tzdata libgcc

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/bin/spell_bot .

# Copy configuration files if any
COPY --from=builder /app/.env.example .env.example

# Change ownership to appuser
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Set environment variables
ENV DEBUG_MODE=false

# Run the application
CMD ["./spell_bot"]