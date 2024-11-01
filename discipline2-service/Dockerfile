# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o discipline2-server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/discipline2-server .

# Expose gRPC port
EXPOSE 50051

# Run the binary
CMD ["./discipline2-server"]