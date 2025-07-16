FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the email processor
RUN CGO_ENABLED=0 GOOS=linux go build -o /email-processor ./email-processor

# Create a minimal production image
FROM alpine:latest

WORKDIR /

# Copy binary from builder
COPY --from=builder /email-processor /email-processor

# Run the email processor
CMD ["/email-processor"]