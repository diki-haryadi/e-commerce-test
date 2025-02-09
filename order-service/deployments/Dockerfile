# Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git build-base

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY .. .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/envs/.env ./envs/

# Expose default port
EXPOSE 8080

# Run with serve argument
ENTRYPOINT ["./main"]
CMD ["serve"]