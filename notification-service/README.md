# Notification Service

A proof-of-concept service for sending email and push notifications using Asynq for message queueing and Gin for the HTTP API.

## Architecture

The system consists of four main components, each in its own directory:

1. **api-server**: API service that exposes an HTTP endpoint for accepting notification requests
2. **queue-service**: Redis-based message queue for delivering tasks between API and processors
3. **email-processor**: Service that processes only email notifications
4. **push-processor**: Service that processes only push notifications

## Project Structure

```
/
├── api-server/              # API server component
│   ├── internal/            # Internal implementation
│   │   ├── models.go        # Data structures
│   │   ├── controller.go    # HTTP controllers
│   │   └── service.go       # Business logic
│   ├── main.go              # Entry point for API server
│   └── Dockerfile           # Container configuration
│
├── queue-service/           # Queue service component
│   ├── internal/            # Internal implementation
│   │   └── queue.go         # Queue configuration
│   └── Dockerfile           # Redis container configuration
│
├── email-processor/         # Email processor component
│   ├── internal/            # Internal implementation
│   │   ├── models.go        # Data structures
│   │   └── processor.go     # Email processing logic
│   ├── main.go              # Entry point for email processor
│   └── Dockerfile           # Container configuration
│
├── push-processor/          # Push processor component
│   ├── internal/            # Internal implementation
│   │   ├── models.go        # Data structures
│   │   └── processor.go     # Push processing logic
│   ├── main.go              # Entry point for push processor
│   └── Dockerfile           # Container configuration
│
├── go.mod                   # Go module definition
├── docker-compose.yml       # Multi-container orchestration
└── README.md                # Project documentation
```

## Requirements

- Docker and Docker Compose
- Go 1.21 or later (for local development)

## Running the Service

Start the service with Docker Compose:

```bash
# Start with default number of processors (2 email, 2 push)
docker-compose up -d

# Start with custom number of processors
EMAIL_REPLICAS=3 PUSH_REPLICAS=5 docker-compose up -d
```

## API Usage

### Send a notification:

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hello World",
    "channel": "email",
    "timeZone": "UTC",
    "recipient": "user@example.com"
  }'
```

### For scheduled notifications, include a `sendAt` timestamp:

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hello World",
    "channel": "push",
    "timeZone": "UTC",
    "recipient": "user123",
    "sendAt": "2023-12-31T23:59:59Z"
  }'
```

### Monitor task status:

#### Get pending tasks (including scheduled tasks):
```bash
curl -X GET http://localhost:8080/tasks/pending
```

#### Get successfully completed tasks:
```bash
curl -X GET http://localhost:8080/tasks/completed
```

#### Get failed tasks (exceeded max retries):
```bash
curl -X GET http://localhost:8080/tasks/failed
```

## Features

- Validates notification requests
- Guarantees at-most-once delivery
- Retries failed notifications up to 3 times
- Supports scheduled notifications
- Simulates 50% failure rate for demonstration purposes
- Configurable number of processors via environment variables

## Component Details

- **api-server**: Implements a layered architecture (controller → service), validates requests, and enqueues tasks
- **queue-service**: Acts as a message broker using Redis
- **email-processor**: Processes email notifications with automatic retries
- **push-processor**: Processes push notifications with automatic retries