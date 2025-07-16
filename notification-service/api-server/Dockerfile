FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the API server
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-server ./api-server

# Create a minimal production image
FROM alpine:latest

WORKDIR /

# Copy binary from builder
COPY --from=builder /api-server /api-server

# Expose API port
EXPOSE 8080

# Run the API server
CMD ["/api-server"]