# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy source code and proto files
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o agronomia

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/agronomia .

# Expose port
EXPOSE 8081

# Run the binary
CMD ["./agronomia"]