FROM golang:1.21-alpine

WORKDIR /app

# Copy go mod files
COPY go.* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o agronomia-monitor

# Run the binary
CMD ["./agronomia-monitor"]