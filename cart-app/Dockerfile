FROM golang:1.23.1

# Set destination for COPY
WORKDIR /app

# Copy cart-app
COPY . .

RUN go mod tidy
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /cart-app

EXPOSE 8081

# Run
CMD ["/cart-app"]