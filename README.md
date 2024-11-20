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

## Architecture diagram

  ![gorun](/readme_imag/flowchart.jpeg)

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
![gorun](/readme_imag/go_run.jpeg)
### To install rabbitMQ and Redis ( follow the documentation for installation and getting started )

###Run the Image Processor
```bash
cd image_processor
go run image_processor.go
```
![img](/readme_imag/imageprocessor.jpeg)
## API Endpoints

### **POST /products**
Description: Create a product with asynchronous image processing.
![img](/readme_imag/postmsn.jpeg)


```json
{
  {
  "user_id": 23,
  "product_name": "Test Image 23",
  "product_description": "Testing with a real image from the internet",
  "product_images": [
    "https://cdn.britannica.com/51/8151-003-FED5CB00/Bernese-mountain-dog.jpg"
  ],
  "product_price": 18.99
}

}
```

#### Response
```json
{
    "data": {
        "id": 44,
        "user_id": 23,
        "product_name": "Test Image 23",
        "product_description": "Testing with a real image from the internet",
        "product_images": [
            "https://cdn.britannica.com/51/8151-003-FED5CB00/Bernese-mountain-dog.jpg"
        ],
        "compressed_product_images": null
    },
    "message": "Product created successfully"
}
```

### **GET /products/**
Description: Retrieve product details by ID.

#### Response

```json
{
    "data": {
        "id": 31,
        "user_id": 13,
        "product_name": "Test Image 13",
        "product_description": "Testing with a real image from the internet viszal.",
        "product_images": null,
        "compressed_product_images": null
    },
    "message": "Product retrieved from database"
}
```
### if stored in cache
```json
{
    "data": "{\"id\":40,\"product_name\":\"Test Product\"}",
    "message": "Product retrieved from cache"
}
```


### GET /products
Description: List all products for a specific user with optional filters.

### Scalability
- RabbitMQ handles message queuing for high-throughput scenarios.
![img](/readme_imag/rabbitmqimage.jpeg)

- Redis and PostgreSQL optimize read/write performance.



### Caching
- Redis caches are used only for the GET /products/:id

![img](/readme_imag/rediscache.jpeg)
### if stored in cache
```json
{
    "data": "{\"id\":40,\"product_name\":\"Test Product\"}",
    "message": "Product retrieved from cache"
}
```




### Logging and Error Handling
- The Logrus library in your Go application provides structured and well-formatted logging output.
- example of how it looks
```
024/11/20 18:46:41 Image processor is running...
2024/11/20 18:46:56 Compressed images to be updated in DB: [https://zocket-dhan.s3.amazonaws.com/Bernese-mountain-dog.jpg]
2024/11/20 18:46:56 Marshalled JSON for DB update: ["https://zocket-dhan.s3.amazonaws.com/Bernese-mountain-dog.jpg"]
2024/11/20 18:46:56 Successfully updated database for product ID 44
```
#### Structured Logging
- Logs are captured using Logrus with fields for request methods, response times, and errors.
- this is how logrus output looks like.
```
{"duration_ms":41,"level":"info","method":"GET","msg":"Request processed","path":"/api/products","status":200,"time":"2024-11-20T20:49:58+05:30"}
```
#### Error Handling
- Includes retry mechanisms for RabbitMQ consumers and proper HTTP status codes for API responses.
![img](/readme_imag/errorhandling1.jpeg)
-

### dead lettet queue
![img](/readme_imag/deadletterque.png)
