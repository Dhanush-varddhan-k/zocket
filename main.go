
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"
	"strconv"

	"github.com/Dhanush-varddhan-k/go-fiber-postgres/models"
	"github.com/Dhanush-varddhan-k/go-fiber-postgres/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type Repository struct {
	DB    *gorm.DB
	Cache *redis.Client
}

var log = logrus.New()

// Middleware for logging requests
func loggingMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	duration := time.Since(start)
	log.WithFields(logrus.Fields{
		"method":      c.Method(),
		"path":        c.Path(),
		"status":      c.Response().StatusCode(),
		"duration_ms": duration.Milliseconds(),
	}).Info("Request processed")
	return err
}

// Converts a slice of strings to a JSON string
func toJSONString(input []string) string {
	jsonData, err := json.Marshal(input)
	if err != nil {
		log.Errorf("Error serializing to JSON: %v", err)
		return "[]"
	}
	return string(jsonData)
}

// RabbitMQ setup
func setupRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	_, err = ch.QueueDeclare(
		"image-processing-queue", // Queue name
		true,                     // Durable
		false,                    // Delete when unused
		false,                    // Exclusive
		false,                    // No-wait
		nil,                      // Arguments
	)
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

// Redis setup
func setupRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})
	return client
}

// Create Product Endpoint
// func (r *Repository) CreateProduct(ctx *fiber.Ctx) error {
// 	product := models.Product{}
// 	err := ctx.BodyParser(&product)
// 	if err != nil {
// 		return ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Invalid data"})
// 	}

// 	// Serialize ProductImages into JSON for database
// 	product.ProductImagesJSON = toJSONString(product.ProductImages)
// 	product.CompressedImagesJSON = "[]" // Placeholder for now

// 	// Save product to the database
// 	err = r.DB.Create(&product).Error
// 	if err != nil {
// 		log.Errorf("Database error: %v", err)
// 		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not save product"})
// 	}

// 	// Publish product images to RabbitMQ for processing
// 	payload := map[string]interface{}{
// 		"product_id": product.ID,
// 		"images":     product.ProductImages,
// 	}
// 	message, _ := json.Marshal(payload)
// 	r.publishToQueue(message)

// 	// Invalidate cache for all products
// 	r.Cache.Del(ctx.Context(), "all_products")

// 	log.WithFields(logrus.Fields{"product_id": product.ID}).Info("Product created successfully")
// 	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "Product created successfully", "data": product})
// }


func (r *Repository) CreateProduct(ctx *fiber.Ctx) error {
	product := models.Product{}
	err := ctx.BodyParser(&product)
	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Invalid data"})
	}

	// Serialize ProductImages into JSON for database
	product.ProductImagesJSON = toJSONString(product.ProductImages)
	product.CompressedImagesJSON = "[]" // Placeholder for now

	// Save product to the database
	err = r.DB.Create(&product).Error
	if err != nil {
		log.Errorf("Database error: %v", err)
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not save product"})
	}

	// Construct cache key
	cacheKey := "product:" + strconv.Itoa(int(product.ID))

	// Cache the product in Redis
	productJSON, _ := json.Marshal(product) // Serialize product
	err = r.Cache.Set(ctx.Context(), cacheKey, productJSON, 0).Err() // Set in Redis
	if err != nil {
		log.Errorf("Failed to cache product in Redis: %v", err)
	} else {
		log.WithFields(logrus.Fields{
			"cache_key": cacheKey,
			"product_id": product.ID,
		}).Info("Caching product in Redis")
	}

	// Publish product images to RabbitMQ for processing
	payload := map[string]interface{}{
		"product_id": product.ID,
		"images":     product.ProductImages,
	}
	message, _ := json.Marshal(payload)
	r.publishToQueue(message)

	// Invalidate cache for all products
	r.Cache.Del(ctx.Context(), "all_products")

	log.WithFields(logrus.Fields{"product_id": product.ID}).Info("Product created successfully")
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "Product created successfully", "data": product})
}


// Update Product Endpoint
func (r *Repository) UpdateProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	product := models.Product{}

	err := ctx.BodyParser(&product)
	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Invalid data"})
	}

	// Update the database
	err = r.DB.Model(&models.Product{}).Where("id = ?", id).Updates(&product).Error
	if err != nil {
		log.Errorf("Failed to update product with ID %s: %v", id, err)
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Failed to update product"})
	}

	// Invalidate cache for the specific product and all products
	r.Cache.Del(ctx.Context(), "product:"+id)
	r.Cache.Del(ctx.Context(), "all_products")

	log.WithFields(logrus.Fields{"product_id": id}).Info("Product updated successfully")
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "Product updated successfully"})
}

// Get Product Endpoint
// func (r *Repository) GetProduct(ctx *fiber.Ctx) error {
// 	id := ctx.Params("id")
// 	cacheKey := "product:" + id

// 	// Check if the product is in cache
// 	cachedProduct, err := r.Cache.Get(context.Background(), cacheKey).Result()
// 	if err == nil {
// 		log.WithFields(logrus.Fields{"product_id": id}).Info("Cache hit for product")
// 		return ctx.Status(http.StatusOK).JSON(&fiber.Map{"data": cachedProduct, "message": "Product retrieved from cache"})
// 	}

// 	// If not in cache, retrieve from database
// 	product := models.Product{}
// 	err = r.DB.First(&product, id).Error
// 	if err != nil {
// 		log.Errorf("Product with ID %s not found: %v", id, err)
// 		return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "Product not found"})
// 	}

// 	// Serialize product to JSON and store in cache
// 	productJSON, _ := json.Marshal(product)
// 	r.Cache.Set(context.Background(), cacheKey, productJSON, 0) // Cache indefinitely

// 	log.WithFields(logrus.Fields{"product_id": id}).Info("Product retrieved from database")
// 	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"data": product, "message": "Product retrieved from database"})
// }

func (r *Repository) GetProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	cacheKey := "product:" + id

	// Check if the product is in cache
	cachedProduct, err := r.Cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		log.WithFields(logrus.Fields{
			"cache_key": cacheKey,
			"product_id": id,
		}).Info("Cache hit for product")
		return ctx.Status(http.StatusOK).JSON(&fiber.Map{"data": cachedProduct, "message": "Product retrieved from cache"})
	}

	// If not in cache, retrieve from database
	product := models.Product{}
	err = r.DB.First(&product, id).Error
	if err != nil {
		log.Errorf("Product with ID %s not found: %v", id, err)
		return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "Product not found"})
	}

	// Serialize product to JSON and store in cache
	productJSON, _ := json.Marshal(product)
	err = r.Cache.Set(context.Background(), cacheKey, productJSON, 0).Err() // Cache indefinitely
	if err != nil {
		log.Errorf("Failed to cache product in Redis: %v", err)
	} else {
		log.WithFields(logrus.Fields{
			"cache_key": cacheKey,
			"product_id": id,
		}).Info("Caching retrieved product in Redis")
	}

	log.WithFields(logrus.Fields{"product_id": id}).Info("Product retrieved from database")
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"data": product, "message": "Product retrieved from database"})
}



// Get All Products Endpoint
func (r *Repository) GetAllProducts(ctx *fiber.Ctx) error {
	cacheKey := "all_products"

	// Check if all products are in cache
	cachedProducts, err := r.Cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		log.Info("Cache hit for all products")
		return ctx.Status(http.StatusOK).JSON(&fiber.Map{"data": cachedProducts, "message": "Products retrieved from cache"})
	}

	// If not in cache, retrieve from database
	var products []models.Product
	err = r.DB.Find(&products).Error
	if err != nil {
		log.Errorf("Failed to retrieve products from database: %v", err)
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Failed to retrieve products"})
	}

	// Serialize products to JSON and store in cache
	productsJSON, _ := json.Marshal(products)
	r.Cache.Set(context.Background(), cacheKey, productsJSON, 0) // Cache indefinitely

	log.Info("Products retrieved from database")
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"data": products, "message": "Products retrieved from database"})
}

// Publish to RabbitMQ
func (r *Repository) publishToQueue(message []byte) {
	conn, ch, err := setupRabbitMQ()
	if err != nil {
		log.Errorf("RabbitMQ connection error: %v", err)
		return
	}
	defer conn.Close()
	defer ch.Close()

	err = ch.Publish(
		"",                       // Exchange
		"image-processing-queue", // Queue name
		false,                    // Mandatory
		false,                    // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		log.Errorf("Failed to publish message to RabbitMQ: %v", err)
	}
}

// Main function
func main() {
	// Initialize structured logging
	log.Formatter = &logrus.JSONFormatter{}
	log.Out = os.Stdout
	log.Level = logrus.InfoLevel

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not connect to the database")
	}

	err = models.Migrate(db)
	if err != nil {
		log.Fatal("Could not migrate database")
	}

	// Initialize Redis client
	cache := setupRedis()

	r := Repository{DB: db, Cache: cache}
	app := fiber.New()

	// Middleware
	app.Use(loggingMiddleware)

	// Endpoints
	app.Post("/api/products", r.CreateProduct)
	app.Put("/api/products/:id", r.UpdateProduct)
	app.Get("/api/products/:id", r.GetProduct)
	app.Get("/api/products", r.GetAllProducts)

	log.Fatal(app.Listen(":8080"))
}
