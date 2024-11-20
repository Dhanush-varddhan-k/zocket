# Product Management System with Asynchronous Image Processing

## Overview
This project is a backend system developed in Golang for managing products. It features RESTful APIs with asynchronous image processing, caching, structured logging, and robust error handling. The system is designed for high scalability and optimal performance, making use of PostgreSQL, Redis, RabbitMQ, and AWS S3.

---

## Features

### RESTful API Design
- **POST /products**: Accepts product data, queues image processing, and stores product details.
- **GET /products/:id**: Retrieves product details with caching for improved performance.
- **GET /products**: Lists all products for a specific user with optional filters.

### Data Storage
- PostgreSQL is used to store product and user details.
- `Products` table includes fields for `compressed_product_images` for processed images.

### Asynchronous Image Processing
- RabbitMQ handles message queuing for image processing tasks.
- A microservice downloads, compresses, and uploads images to AWS S3, updating the database upon completion.

### Caching
- Redis is used to cache `GET /products/:id` responses to reduce database load.
- Cache invalidation ensures updated product data is reflected in real-time.

### Structured Logging
- Logrus provides detailed structured logs for all components, including API requests and image processing events.

### Error Handling
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
1. **Install Go**: Ensure Go is installed on your system.
2. **InstallPostgreSQL, RabbitMQ, and Redis**
3. **AWS S3**: Set up an S3 bucket for image storage with an IAM user.
4. **Environment Variables**: Create a `.env` file in the root directory with the following:

```plaintext
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
```
# Steps to Run the Project

## Clone the Repository
```bash
git clone https://github.com/your_username/go-fiber-postgres.git
cd go-fiber-postgres
```
##install dependencies
```bash
go mod tidy
```

###Run the Main Application
```bash
go run main.go
```
### To install rabbitMQ and Redis ( follow the documentation for installation and getting started )

###Run the Image Processor
```bash
cd image_processor
go run image_processor.go
```
## API Endpoints

### **POST /products**
Description: Create a product with asynchronous image processing.

#### Request Body
```json
{
  "user_id": 1,
  "product_name": "Sample Product",
  "product_description": "Description here",
  "product_images": ["https://example.com/image1.jpg"],
  "product_price": 19.99
}
```

#### Response
```json
{
  "message": "Product created successfully",
  "data": {
    "id": 42,
    "user_id": 1,
    "product_name": "Sample Product"
  }
}
```

### **GET /products/**
Description: Retrieve product details by ID.

#### Response

```json
{
  "data": {
    "id": 42,
    "user_id": 1,
    "product_name": "Sample Product",
    "compressed_product_images": ["https://bucket.s3.amazonaws.com/compressed_image.jpg"]
  },
  "message": "Product retrieved successfully"
}
```

### GET /products
Description: List all products for a specific user with optional filters.

#### Query Parameters:
- user_id (required)
- price_min (optional)
- price_max (optional)
- product_name (optional)


### Caching
- Redis caches are used only for the GET /products/:id endpoint.

### Scalability
- RabbitMQ handles message queuing for high-throughput scenarios.
- Redis and PostgreSQL optimize read/write performance.

### Logging and Error Handling
#### Structured Logging
- Logs are captured using Logrus with fields for request methods, response times, and errors.
#### Error Handling
- Includes retry mechanisms for RabbitMQ consumers and proper HTTP status codes for API responses.
