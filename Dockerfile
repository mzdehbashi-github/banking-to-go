# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .

# Download and install golang-migrate
RUN wget -O migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz && \
    tar -xzf migrate.tar.gz && \
    mv migrate /usr/local/bin/ && \
    rm migrate.tar.gz

# Build the main application
RUN go build -o main main.go

# Run stage
FROM alpine
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

# Copy commands.sh script
COPY commands.sh /app/commands.sh

# Set execute permissions for commands.sh
RUN chmod +x /app/commands.sh


COPY db/migrations ./db/migrations

# Expose port 8080
EXPOSE 8080

# Set commands.sh as the entrypoint
ENTRYPOINT ["/app/commands.sh"]
