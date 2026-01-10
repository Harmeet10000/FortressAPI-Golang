# Development Dockerfile for Go 1.25
# Optimized for hot reload with air and debugging

FROM golang:1.25-alpine

WORKDIR /app

# Install required tools
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata

# Install air for hot reload
RUN go install github.com/cosmtrek/air@latest

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Expose port
EXPOSE 8080 8081

# Run with air for hot reload
CMD ["air", "-c", ".air.toml"]
