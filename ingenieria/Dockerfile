# Build Stage
FROM rust:1.72.1-slim-bullseye as builder

WORKDIR /usr/src/app

# Install protobuf compiler
RUN apt-get update && \
    apt-get install -y protobuf-compiler && \
    rm -rf /var/lib/apt/lists/*

# Copy all files first
COPY . .

# Build the application with release profile
RUN cargo build --release

# Runtime Stage - Using debian slim for smaller image size
FROM debian:bullseye-slim

# Install necessary runtime dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends libssl-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /usr/src/app/target/release/ingenieria /app/ingenieria

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./ingenieria"]