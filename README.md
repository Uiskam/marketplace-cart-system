# Distributed Marketplace System

## 1. Project Overview
This project implements a distributed marketplace system composed of three main microservices: `cart-app`, `product-app`, and `notification-service`. The system enables users to manage shopping carts, interact with product inventory, and receive notifications about their actions. Each service is independently deployable and communicates via well-defined interfaces, supporting scalability and maintainability.

## 2. Architecture Overview

### cart-app
- **Purpose:** Manages user shopping carts, supporting actions such as creating a cart, adding/removing items, and checking out.
- **Key Actions:**
  - `addToCart`: Add products to a user's cart.
  - `removeFromCart`: Remove products from a cart.
  - `getCart`: Retrieve the current state of a cart.
  - `checkoutCart`: Finalize the purchase and trigger product locking.
  - `createCart`: Initialize a new cart for a user.

### product-app
- **Purpose:** Handles product inventory, including querying product details, locking/unlocking products during checkout, and marking products as sold.
- **Key Actions:**
  - `getProduct`: Retrieve product information.
  - `lockProduct`: Reserve products for a cart during checkout.
  - `unlockProduct`: Release reserved products if checkout is canceled.
  - `sellProduct`: Mark products as sold after successful checkout.

### notification-service
- **Purpose:** Sends notifications to users via email or push messages. Operates independently but receives updates from `cart-app` (e.g., on checkout or cart changes).
- **Components:**
  - `api-server`: Receives notification requests.
  - `email-processor`: Handles email delivery.
  - `push-processor`: Handles push notifications.
  - `queue-service`: Manages notification queues.

## 3. Design Patterns Used
- **Command Query Responsibility Segregation (CQRS):**
  - Separates read (queries) and write (commands) operations for scalability and clarity.
- **Clean Architecture:**
  - Decouples business logic from infrastructure, using layers such as handlers, controllers, and command/query objects.
- **Microservices with Dockerized Deployments:**
  - Each service is containerized for independent deployment and scaling.

## 4. Detailed Breakdown

### actions Directories
- Each action (e.g., `addToCart`, `lockProduct`) is implemented as a directory containing:
  - `command.go` or `query.go`: Defines the command/query structure.
  - `controller.go`: Exposes HTTP endpoints and maps requests to handlers.
  - `handler.go`: Contains business logic for processing the command/query.

### app/cqrs
- Contains the CQRS infrastructure:
  - `bus.go`: Implements the command/query bus for dispatching requests.
  - `processor.go`: Handles the registration and execution of command/query handlers.

### external Interfaces
- Abstracts integration with external systems:
  - `postgresql.go`: Database access layer.
  - `product.go`: API client for product operations (in `cart-app`).
  - `notification.go`: API client for sending notifications (in `cart-app`).
  - `base.sql`: Database schema (in `product-app`).

### repository
- Encapsulates persistence logic:
  - `cart/` and `product/` directories contain `read.go` and `write.go` for data access.
  - `model/` subdirectories define data models (e.g., `cart.go`, `product.go`).

### Docker and Entrypoints
- Each service has its own `Dockerfile` and `docker-compose.yaml` for containerized deployment.
- `main.go` in each service is the entrypoint, initializing the server and dependencies.

## 5. Inter-module Coupling
- **cart-app ↔ product-app:**
  - `cart-app` depends on `product-app` to lock/unlock products during checkout, ensuring inventory consistency.
- **Shared CQRS Concepts:**
  - Both `cart-app` and `product-app` use similar CQRS processors for handling commands and queries, promoting consistency and code reuse.
- **notification-service:**
  - Operates independently but receives updates (e.g., checkout events) from `cart-app` to notify users.

---

## Getting Started
0. **Create docker network:**
   ```zsh
   docker network create shared-network
   ```
1. **Build and Run All Services:**
   ```zsh
   ./run.sh
   ```
2. **API Documentation:**
   - Postman collections are provided in each service directory for testing endpoints.

## License
This project is for educational purposes.
