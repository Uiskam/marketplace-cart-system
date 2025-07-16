FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the push processor
RUN CGO_ENABLED=0 GOOS=linux go build -o /push-processor ./push-processor

# Create a minimal production image
FROM alpine:latest

WORKDIR /

# Copy binary from builder
COPY --from=builder /push-processor /push-processor

# Run the push processor
CMD ["/push-processor"]