FROM golang:1.23.1

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY . .

RUN go mod tidy
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /product-api

EXPOSE 8000

# Run
CMD ["/product-api"]