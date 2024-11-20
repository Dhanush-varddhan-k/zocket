# Product Management System with Asynchronous Image Processing

## Overview

This project is a backend system developed in Golang for managing products. It features RESTful APIs with asynchronous image processing, caching, structured logging, and robust error handling. The system is designed for high scalability and optimal performance, making use of PostgreSQL, Redis, RabbitMQ, and AWS S3.

---

## Features

1. **RESTful API Design**:
   - `POST /products`: Accepts product data, queues image processing, and stores product details.
   - `GET /products/:id`: Retrieves product details with caching for improved performance.
   - `GET /products`: Lists all products for a specific user with optional filters.

2. **Data Storage**:
   - PostgreSQL is used to store product and user details.
   - Products table includes fields for `compressed_product_images` for processed images.

3. **Asynchronous Image Processing**:
   - RabbitMQ handles message queuing for image processing tasks.
   - A microservice downloads, compresses, and uploads images to AWS S3, updating the database upon completion.

4. **Caching**:
   - Redis is used to cache `GET /products/:id` responses to reduce database load.
   - Cache invalidation ensures updated product data is reflected in real-time.

5. **Structured Logging**:
   - Logrus provides detailed structured logs for all components, including API requests and image processing events.

6. **Error Handling**:
   - Implements robust error handling for API operations and asynchronous processing with retry mechanisms and dead-letter queues.

---

## Architecture

- **Language**: Go
- **Database**: PostgreSQL
- **Message Queue**: RabbitMQ
- **Cache**: Redis
- **Cloud Storage**: AWS S3
- **Logging**: Logrus
- **Framework**: Fiber (HTTP framework for Go)

---

## Setup Instructions

### Prerequisites

Before running the project, ensure the following tools are installed:

1. **Install Go**:
   Download and install Go from [golang.org](https://golang.org/dl/).

2. **Install Docker**:
   Docker is required for PostgreSQL, Redis, and RabbitMQ containers. Install it from [docker.com](https://www.docker.com/).

3. **AWS S3 Setup**:
   Set up an S3 bucket for storing compressed images.

4. **Environment Variables**:
   Create a `.env` file in the root directory and copy-paste the following:

   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_db_user
   DB_PASS=your_db_password
   DB_NAME=your_db_name
   DB_SSLMODE=disable

   REDIS_HOST=localhost
   REDIS_PORT=6379

   RABBITMQ_URL=amqp://guest:guest@localhost:5672/
   S3_BUCKET_NAME=your_bucket_name
   AWS_REGION=ap-south-1
   AWS_ACCESS_KEY=your_access_key
   AWS_SECRET_KEY=your_secret_key
