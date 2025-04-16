# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o websocket-server ./cmd/websocket-server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/websocket-server .

# Expose WebSocket port
EXPOSE 8081

# Run the server
CMD ["./websocket-server"]