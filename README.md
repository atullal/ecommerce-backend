# E-commerce Backend Microservices

## Overview

This project is a microservices-based backend system for an e-commerce application. It is designed to handle user authentication, product management, order processing, and inventory management.

## Services

The system is composed of the following services:

1. **User Service**: Manages user registration and authentication.
2. **Product Service**: Handles product information, including CRUD operations and inventory management.
3. **Order Service**: Responsible for order creation, retrieval, and management.

## Technologies

- **Programming Language**: Go (Golang)
- **Database**: PostgreSQL
- **Microservices Framework**: gRPC for inter-service communication and Gin for RESTful endpoints.
- **Other Tools**: Docker for containerization, `protoc` for generating gRPC code from Protobuf files.

## Getting Started

### Prerequisites

- Go (version 1.x)
- Docker and Docker Compose
- PostgreSQL (if running locally without Docker)
- Protobuf compiler (`protoc`)

### Installation and Setup

1. **Clone the Repository**:
   ```sh
   git clone [repository URL]
   cd ecommerce-backend
   ```
2. **Set Up Environment Variables**:

   Create .env files for each service with necessary configurations like database connection strings.
3. **Build the Services**:

   Use Docker Compose to build and run the services:
   ```sh
   Copy code
   docker-compose up --build
   ```
### API Documentation
   Swagger is used for API documentation. Access the Swagger UI at [service URL]/swagger/index.html for RESTful services.
### Usage
   User Service: Register new users, authenticate existing users.

   Product Service: Add new products, retrieve product information, update product details, and manage inventory.

   Order Service: Place orders, retrieve order details, and manage orders.
   Development and Contribution

   To contribute to this project, create a new branch and submit pull requests for review.

   Ensure to write tests and maintain coding standards.
