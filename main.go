// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/Dhanush-varddhan-k/go-fiber-postgres/models"
// 	"github.com/Dhanush-varddhan-k/go-fiber-postgres/storage"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// 	"gorm.io/gorm"
// )

// type Book struct {
// 	Author    string `json:"author"`
// 	Title     string `json:"title"`
// 	Publisher string `json:"publisher"`
// }

// type Repository struct {
// 	DB *gorm.DB
// }

// func (r *Repository) CreateBook(context *fiber.Ctx) error {
// 	book := Book{}

// 	err := context.BodyParser(&book)

// 	if err != nil {
// 		context.Status(http.StatusUnprocessableEntity).JSON(
// 			&fiber.Map{"message": "request failed"})
// 		return err
// 	}

// 	err = r.DB.Create(&book).Error
// 	if err != nil {
// 		context.Status(http.StatusBadRequest).JSON(
// 			&fiber.Map{"message": "could not create book"})
// 		return err
// 	}

// 	context.Status(http.StatusOK).JSON(&fiber.Map{
// 		"message": "book has been added"})
// 	return nil
// }

// func (r *Repository) DeleteBook(context *fiber.Ctx) error {
// 	bookModel := models.Books{}
// 	id := context.Params("id")
// 	if id == "" {
// 		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
// 			"message": "id cannot be empty",
// 		})
// 		return nil
// 	}

// 	err := r.DB.Delete(bookModel, id)

// 	if err.Error != nil {
// 		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
// 			"message": "could not delete book",
// 		})
// 		return err.Error
// 	}
// 	context.Status(http.StatusOK).JSON(&fiber.Map{
// 		"message": "book delete successfully",
// 	})
// 	return nil
// }

// func (r *Repository) GetBooks(context *fiber.Ctx) error {
// 	bookModels := &[]models.Books{}

// 	err := r.DB.Find(bookModels).Error
// 	if err != nil {
// 		context.Status(http.StatusBadRequest).JSON(
// 			&fiber.Map{"message": "could not get books"})
// 		return err
// 	}

// 	context.Status(http.StatusOK).JSON(&fiber.Map{
// 		"message": "books fetched successfully",
// 		"data":    bookModels,
// 	})
// 	return nil
// }

// func (r *Repository) GetBookByID(context *fiber.Ctx) error {

// 	id := context.Params("id")
// 	bookModel := &models.Books{}
// 	if id == "" {
// 		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
// 			"message": "id cannot be empty",
// 		})
// 		return nil
// 	}

// 	fmt.Println("the ID is", id)

// 	err := r.DB.Where("id = ?", id).First(bookModel).Error
// 	if err != nil {
// 		context.Status(http.StatusBadRequest).JSON(
// 			&fiber.Map{"message": "could not get the book"})
// 		return err
// 	}
// 	context.Status(http.StatusOK).JSON(&fiber.Map{
// 		"message": "book id fetched successfully",
// 		"data":    bookModel,
// 	})
// 	return nil
// }

// func (r *Repository) SetupRoutes(app *fiber.App) {
// 	api := app.Group("/api")
// 	api.Post("/create_books", r.CreateBook)
// 	api.Delete("delete_book/:id", r.DeleteBook)
// 	api.Get("/get_books/:id", r.GetBookByID)
// 	api.Get("/books", r.GetBooks)
// }

// func main() {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	config := &storage.Config{
// 		Host:     os.Getenv("DB_HOST"),
// 		Port:     os.Getenv("DB_PORT"),
// 		Password: os.Getenv("DB_PASS"),
// 		User:     os.Getenv("DB_USER"),
// 		SSLMode:  os.Getenv("DB_SSLMODE"),
// 		DBName:   os.Getenv("DB_NAME"),
// 	}

// 	db, err := storage.NewConnection(config)

// 	if err != nil {
// 		log.Fatal("could not load the database")
// 	}
// 	err = models.MigrateBooks(db)
// 	if err != nil {
// 		log.Fatal("could not migrate db")
// 	}

//		r := Repository{
//			DB: db,
//		}
//		app := fiber.New()
//		r.SetupRoutes(app)
//		app.Listen(":8080")
//	}

// package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strconv"

// 	"github.com/Dhanush-varddhan-k/go-fiber-postgres/models"
// 	"github.com/Dhanush-varddhan-k/go-fiber-postgres/storage"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// 	"gorm.io/gorm"
// )

// type Repository struct {
// 	DB *gorm.DB
// }

// // Serialize slices to JSON strings for storage
// func toJSONString(input interface{}) string {
// 	jsonData, err := json.Marshal(input)
// 	if err != nil {
// 		log.Printf("Error serializing to JSON: %v", err)
// 		return "[]"
// 	}
// 	return string(jsonData)
// }

// func (r *Repository) CreateProduct(context *fiber.Ctx) error {
// 	product := models.Product{}

// 	// Parse the request body into the product struct
// 	err := context.BodyParser(&product)
// 	if err != nil {
// 		context.Status(http.StatusUnprocessableEntity).JSON(
// 			&fiber.Map{"message": "request failed"})
// 		return err
// 	}

// 	// Serialize product_images and compressed_product_images to JSON strings
// 	product.ProductImages = toJSONString(product.ProductImages)
// 	product.CompressedProductImages = toJSONString(product.CompressedProductImages)

// 	// Save the product to the database
// 	err = r.DB.Create(&product).Error
// 	if err != nil {
// 		context.Status(http.StatusBadRequest).JSON(
// 			&fiber.Map{"message": "could not create product"})
// 		return err
// 	}

// 	context.Status(http.StatusOK).JSON(&fiber.Map{
// 		"message": "product has been added",
// 	})
// 	return nil
// }

// func (r *Repository) GetProductByID(context *fiber.Ctx) error {
// 	id := context.Params("id")
// 	product := &models.Product{}
// 	if id == "" {
// 		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
// 			"message": "id cannot be empty",
// 		})
// 		return nil
// 	}

// 	// Query the product by ID
// 	err := r.DB.Where("id = ?", id).First(product).Error
// 	if err != nil {
// 		context.Status(http.StatusBadRequest).JSON(
// 			&fiber.Map{"message": "could not get the product"})
// 		return err
// 	}

// 	context.Status(http.StatusOK).JSON(&fiber.Map{
// 		"message": "product fetched successfully",
// 		"data":    product,
// 	})
// 	return nil
// }

// func (r *Repository) GetProducts(context *fiber.Ctx) error {
// 	userID := context.Query("user_id")
// 	minPrice := context.Query("min_price")
// 	maxPrice := context.Query("max_price")
// 	productName := context.Query("product_name")

// 	products := &[]models.Product{}
// 	query := r.DB

// 	// Add filters to the query
// 	if userID != "" {
// 		query = query.Where("user_id = ?", userID)
// 	}
// 	if minPrice != "" && maxPrice != "" {
// 		min, _ := strconv.ParseFloat(minPrice, 64)
// 		max, _ := strconv.ParseFloat(maxPrice, 64)
// 		query = query.Where("product_price BETWEEN ? AND ?", min, max)
// 	}
// 	if productName != "" {
// 		query = query.Where("product_name ILIKE ?", "%"+productName+"%")
// 	}

// 	// Execute the query
// 	err := query.Find(products).Error
// 	if err != nil {
// 		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
// 			"message": "could not get products",
// 		})
// 		return err
// 	}

// 	context.Status(http.StatusOK).JSON(&fiber.Map{
// 		"message": "products fetched successfully",
// 		"data":    products,
// 	})
// 	return nil
// }

// func (r *Repository) SetupRoutes(app *fiber.App) {
// 	api := app.Group("/api")
// 	api.Post("/products", r.CreateProduct)
// 	api.Get("/products/:id", r.GetProductByID)
// 	api.Get("/products", r.GetProducts)
// }

// func main() {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	config := &storage.Config{
// 		Host:     os.Getenv("DB_HOST"),
// 		Port:     os.Getenv("DB_PORT"),
// 		Password: os.Getenv("DB_PASS"),
// 		User:     os.Getenv("DB_USER"),
// 		SSLMode:  os.Getenv("DB_SSLMODE"),
// 		DBName:   os.Getenv("DB_NAME"),
// 	}

// 	db, err := storage.NewConnection(config)
// 	if err != nil {
// 		log.Fatal("could not load the database")
// 	}

// 	err = models.Migrate(db)
// 	if err != nil {
// 		log.Fatal("could not migrate db")
// 	}

// 	r := Repository{
// 		DB: db,
// 	}

// 	app := fiber.New()
// 	r.SetupRoutes(app)
// 	app.Listen(":8080")
// }

// package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strconv"

// 	"github.com/Dhanush-varddhan-k/go-fiber-postgres/models"
// 	"github.com/Dhanush-varddhan-k/go-fiber-postgres/storage"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// 	"gorm.io/gorm"
// )

// type Repository struct {
// 	DB *gorm.DB
// }

// // Converts a slice of strings to a JSON string for storage
// func toJSONString(input []string) string {
// 	jsonData, err := json.Marshal(input)
// 	if err != nil {
// 		log.Printf("Error serializing to JSON: %v", err)
// 		return "[]"
// 	}
// 	return string(jsonData)
// }

// // Converts a JSON string back to a slice of strings for API response
// func fromJSONString(input string) []string {
// 	var output []string
// 	err := json.Unmarshal([]byte(input), &output)
// 	if err != nil {
// 		log.Printf("Error deserializing JSON: %v", err)
// 		return []string{}
// 	}
// 	return output
// }

// // CreateProduct handles the creation of a new product
// func (r *Repository) CreateProduct(context *fiber.Ctx) error {
// 	product := models.Product{}

// 	// Parse request body into product struct
// 	err := context.BodyParser(&product)
// 	if err != nil {
// 		context.Status(http.StatusUnprocessableEntity).JSON(
// 			&fiber.Map{"message": "request failed"})
// 		return err
// 	}

// 	// Serialize arrays into JSON strings for storage in the database
// 	product.ProductImagesJSON = toJSONString(product.ProductImages)
// 	product.CompressedImagesJSON = toJSONString(product.CompressedProductImages)

// 	// Save the product to the database
// 	err = r.DB.Create(&product).Error
// 	if err != nil {
// 		context.Status(http.StatusBadRequest).JSON(
// 			&fiber.Map{"message": "could not create product"})
// 		return err
// 	}

// 	context.Status(http.StatusOK).JSON(&fiber.Map{
// 		"message": "product has been added",
// 	})
// 	return nil
// }

// // GetProductByID retrieves a single product by ID
// func (r *Repository) GetProductByID(context *fiber.Ctx) error {
// 	id := context.Params("id")
// 	product := &models.Product{}
// 	if id == "" {
// 		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
// 			"message": "id cannot be empty",
// 		})
// 		return nil
// 	}

// 	// Fetch product by ID
// 	err := r.DB.Where("id = ?", id).First(product).Error
// 	if err != nil {
// 		context.Status(http.StatusBadRequest).JSON(
// 			&fiber.Map{"message": "could not get the product"})
// 		return err
// 	}

// 	// Deserialize JSON strings back to arrays for response
// 	product.ProductImages = fromJSONString(product.ProductImagesJSON)
// 	product.CompressedProductImages = fromJSONString(product.CompressedImagesJSON)

// 	context.Status(http.StatusOK).JSON(&fiber.Map{
// 		"message": "product fetched successfully",
// 		"data":    product,
// 	})
// 	return nil
// }

// // GetProducts retrieves all products with optional filters
// func (r *Repository) GetProducts(context *fiber.Ctx) error {
// 	userID := context.Query("user_id")
// 	minPrice := context.Query("min_price")
// 	maxPrice := context.Query("max_price")
// 	productName := context.Query("product_name")

// 	products := &[]models.Product{}
// 	query := r.DB

// 	// Apply filters based on query parameters
// 	if userID != "" {
// 		query = query.Where("user_id = ?", userID)
// 	}
// 	if minPrice != "" && maxPrice != "" {
// 		min, _ := strconv.ParseFloat(minPrice, 64)
// 		max, _ := strconv.ParseFloat(maxPrice, 64)
// 		query = query.Where("product_price BETWEEN ? AND ?", min, max)
// 	}
// 	if productName != "" {
// 		query = query.Where("product_name ILIKE ?", "%"+productName+"%")
// 	}

// 	// Fetch products
// 	err := query.Find(products).Error
// 	if err != nil {
// 		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
// 			"message": "could not get products",
// 		})
// 		return err
// 	}

// 	// Deserialize JSON strings back to arrays for response
// 	for i := range *products {
// 		(*products)[i].ProductImages = fromJSONString((*products)[i].ProductImagesJSON)
// 		(*products)[i].CompressedProductImages = fromJSONString((*products)[i].CompressedImagesJSON)
// 	}

// 	context.Status(http.StatusOK).JSON(&fiber.Map{
// 		"message": "products fetched successfully",
// 		"data":    products,
// 	})
// 	return nil
// }

// // SetupRoutes sets up the API routes
// func (r *Repository) SetupRoutes(app *fiber.App) {
// 	api := app.Group("/api")
// 	api.Post("/products", r.CreateProduct)
// 	api.Get("/products/:id", r.GetProductByID)
// 	api.Get("/products", r.GetProducts)
// }

// func main() {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	config := &storage.Config{
// 		Host:     os.Getenv("DB_HOST"),
// 		Port:     os.Getenv("DB_PORT"),
// 		Password: os.Getenv("DB_PASS"),
// 		User:     os.Getenv("DB_USER"),
// 		SSLMode:  os.Getenv("DB_SSLMODE"),
// 		DBName:   os.Getenv("DB_NAME"),
// 	}

// 	db, err := storage.NewConnection(config)
// 	if err != nil {
// 		log.Fatal("could not load the database")
// 	}

// 	err = models.Migrate(db)
// 	if err != nil {
// 		log.Fatal("could not migrate db")
// 	}

// 	r := Repository{
// 		DB: db,
// 	}

// 	app := fiber.New()
// 	r.SetupRoutes(app)
// 	app.Listen(":8080")
// }

// package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/Dhanush-varddhan-k/go-fiber-postgres/models"
// 	"github.com/Dhanush-varddhan-k/go-fiber-postgres/storage"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// 	"github.com/streadway/amqp"
// 	"gorm.io/gorm"
// )

// type Repository struct {
// 	DB *gorm.DB
// }

// // Converts a slice of strings to a JSON string
// func toJSONString(input []string) string {
// 	jsonData, err := json.Marshal(input)
// 	if err != nil {
// 		log.Printf("Error serializing to JSON: %v", err)
// 		return "[]"
// 	}
// 	return string(jsonData)
// }

// func setupRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
// 	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	_, err = ch.QueueDeclare(
// 		"image-processing-queue", // Queue name
// 		true,                     // Durable
// 		false,                    // Delete when unused
// 		false,                    // Exclusive
// 		false,                    // No-wait
// 		nil,                      // Arguments
// 	)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return conn, ch, nil
// }

// func (r *Repository) CreateProduct(context *fiber.Ctx) error {
// 	product := models.Product{}
// 	err := context.BodyParser(&product)
// 	if err != nil {
// 		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Invalid data"})
// 		return err
// 	}

// 	// Serialize ProductImages into JSON
// 	product.ProductImagesJSON = toJSONString(product.ProductImages)

// 	// Save product to the database
// 	err = r.DB.Create(&product).Error
// 	if err != nil {
// 		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not save product"})
// 		return err
// 	}

// 	// Publish product images to RabbitMQ for processing
// 	r.publishToQueue(product.ProductImages)

// 	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Product created successfully"})
// 	return nil
// }

// func (r *Repository) publishToQueue(images []string) {
// 	conn, ch, err := setupRabbitMQ()
// 	if err != nil {
// 		log.Printf("RabbitMQ connection error: %v", err)
// 		return
// 	}
// 	defer conn.Close()
// 	defer ch.Close()

// 	// Serialize the images slice into a JSON message
// 	message, err := json.Marshal(images)
// 	if err != nil {
// 		log.Printf("Failed to serialize images: %v", err)
// 		return
// 	}

// 	// Publish the message to the RabbitMQ queue
// 	err = ch.Publish(
// 		"",                       // Exchange
// 		"image-processing-queue", // Queue name
// 		false,                    // Mandatory
// 		false,                    // Immediate
// 		amqp.Publishing{
// 			ContentType: "application/json",
// 			Body:        message,
// 		},
// 	)
// 	if err != nil {
// 		log.Printf("Failed to publish message: %v", err)
// 	}
// }

// func main() {
// 	// Load environment variables
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Database configuration
// 	config := &storage.Config{
// 		Host:     os.Getenv("DB_HOST"),
// 		Port:     os.Getenv("DB_PORT"),
// 		Password: os.Getenv("DB_PASS"),
// 		User:     os.Getenv("DB_USER"),
// 		SSLMode:  os.Getenv("DB_SSLMODE"),
// 		DBName:   os.Getenv("DB_NAME"),
// 	}

// 	// Establish database connection
// 	db, err := storage.NewConnection(config)
// 	if err != nil {
// 		log.Fatal("Could not connect to the database")
// 	}

// 	// Run database migrations
// 	err = models.Migrate(db)
// 	if err != nil {
// 		log.Fatal("Could not migrate database")
// 	}

// 	// Set up the Fiber web server
// 	r := Repository{DB: db}
// 	app := fiber.New()

// 	// Define routes
// 	app.Post("/api/products", r.CreateProduct)

// 	// Start the server
// 	app.Listen(":8080")
// }

//working

// package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/Dhanush-varddhan-k/go-fiber-postgres/models"
// 	"github.com/Dhanush-varddhan-k/go-fiber-postgres/storage"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// 	"github.com/streadway/amqp"
// 	"gorm.io/gorm"
// )

// type Repository struct {
// 	DB *gorm.DB
// }

// // Converts a slice of strings to a JSON string
// func toJSONString(input []string) string {
// 	jsonData, err := json.Marshal(input)
// 	if err != nil {
// 		log.Printf("Error serializing to JSON: %v", err)
// 		return "[]"
// 	}
// 	return string(jsonData)
// }

// func setupRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
// 	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	_, err = ch.QueueDeclare(
// 		"image-processing-queue", // Queue name
// 		true,                     // Durable
// 		false,                    // Delete when unused
// 		false,                    // Exclusive
// 		false,                    // No-wait
// 		nil,                      // Arguments
// 	)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return conn, ch, nil
// }

// func (r *Repository) CreateProduct(context *fiber.Ctx) error {
// 	product := models.Product{}
// 	err := context.BodyParser(&product)
// 	if err != nil {
// 		return context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Invalid data"})
// 	}

// 	// Serialize ProductImages into JSON for database
// 	product.ProductImagesJSON = toJSONString(product.ProductImages)
// 	product.CompressedImagesJSON = "[]" // Placeholder for now

// 	// Save product to the database
// 	err = r.DB.Create(&product).Error
// 	if err != nil {
// 		log.Printf("Database error: %v", err)
// 		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not save product"})
// 	}

// 	// Publish product images to RabbitMQ for processing
// 	r.publishToQueue(product.ProductImages)

// 	return context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Product created successfully", "data": product})
// }

// func (r *Repository) publishToQueue(images []string) {
// 	conn, ch, err := setupRabbitMQ()
// 	if err != nil {
// 		log.Printf("RabbitMQ connection error: %v", err)
// 		return
// 	}
// 	defer conn.Close()
// 	defer ch.Close()

// 	message, err := json.Marshal(images)
// 	if err != nil {
// 		log.Printf("Failed to serialize images: %v", err)
// 		return
// 	}

// 	err = ch.Publish(
// 		"",                       // Exchange
// 		"image-processing-queue", // Queue name
// 		false,                    // Mandatory
// 		false,                    // Immediate
// 		amqp.Publishing{
// 			ContentType: "application/json",
// 			Body:        message,
// 		},
// 	)
// 	if err != nil {
// 		log.Printf("Failed to publish message: %v", err)
// 	}
// }

// func main() {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	config := &storage.Config{
// 		Host:     os.Getenv("DB_HOST"),
// 		Port:     os.Getenv("DB_PORT"),
// 		Password: os.Getenv("DB_PASS"),
// 		User:     os.Getenv("DB_USER"),
// 		SSLMode:  os.Getenv("DB_SSLMODE"),
// 		DBName:   os.Getenv("DB_NAME"),
// 	}

// 	db, err := storage.NewConnection(config)
// 	if err != nil {
// 		log.Fatal("Could not connect to the database")
// 	}

// 	err = models.Migrate(db)
// 	if err != nil {
// 		log.Fatal("Could not migrate database")
// 	}

// 	r := Repository{DB: db}
// 	app := fiber.New()

// 	app.Post("/api/products", r.CreateProduct)

//		log.Fatal(app.Listen(":8080"))
//	}

// package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/Dhanush-varddhan-k/go-fiber-postgres/models"
// 	"github.com/Dhanush-varddhan-k/go-fiber-postgres/storage"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// 	"github.com/streadway/amqp"
// 	"gorm.io/gorm"
// )

// type Repository struct {
// 	DB *gorm.DB
// }

// // Converts a slice of strings to a JSON string
// func toJSONString(input []string) string {
// 	jsonData, err := json.Marshal(input)
// 	if err != nil {
// 		log.Printf("Error serializing to JSON: %v", err)
// 		return "[]"
// 	}
// 	return string(jsonData)
// }

// func setupRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
// 	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	_, err = ch.QueueDeclare(
// 		"image-processing-queue", // Queue name
// 		true,                     // Durable
// 		false,                    // Delete when unused
// 		false,                    // Exclusive
// 		false,                    // No-wait
// 		nil,                      // Arguments
// 	)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return conn, ch, nil
// }

// func (r *Repository) CreateProduct(context *fiber.Ctx) error {
// 	product := models.Product{}
// 	err := context.BodyParser(&product)
// 	if err != nil {
// 		return context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Invalid data"})
// 	}

// 	// Serialize ProductImages into JSON for database
// 	product.ProductImagesJSON = toJSONString(product.ProductImages)
// 	product.CompressedImagesJSON = "[]" // Placeholder for now

// 	// Save product to the database
// 	err = r.DB.Create(&product).Error
// 	if err != nil {
// 		log.Printf("Database error: %v", err)
// 		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not save product"})
// 	}

// 	// Publish product images to RabbitMQ for processing
// 	payload := map[string]interface{}{
// 		"product_id": product.ID,
// 		"images":     product.ProductImages,
// 	}
// 	message, _ := json.Marshal(payload)
// 	r.publishToQueue(message)

// 	return context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Product created successfully", "data": product})
// }

// func (r *Repository) publishToQueue(message []byte) {
// 	conn, ch, err := setupRabbitMQ()
// 	if err != nil {
// 		log.Printf("RabbitMQ connection error: %v", err)
// 		return
// 	}
// 	defer conn.Close()
// 	defer ch.Close()

// 	err = ch.Publish(
// 		"",                       // Exchange
// 		"image-processing-queue", // Queue name
// 		false,                    // Mandatory
// 		false,                    // Immediate
// 		amqp.Publishing{
// 			ContentType: "application/json",
// 			Body:        message,
// 		},
// 	)
// 	if err != nil {
// 		log.Printf("Failed to publish message: %v", err)
// 	}
// }

// func main() {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	config := &storage.Config{
// 		Host:     os.Getenv("DB_HOST"),
// 		Port:     os.Getenv("DB_PORT"),
// 		Password: os.Getenv("DB_PASS"),
// 		User:     os.Getenv("DB_USER"),
// 		SSLMode:  os.Getenv("DB_SSLMODE"),
// 		DBName:   os.Getenv("DB_NAME"),
// 	}

// 	db, err := storage.NewConnection(config)
// 	if err != nil {
// 		log.Fatal("Could not connect to the database")
// 	}

// 	err = models.Migrate(db)
// 	if err != nil {
// 		log.Fatal("Could not migrate database")
// 	}

// 	r := Repository{DB: db}
// 	app := fiber.New()

// 	app.Post("/api/products", r.CreateProduct)

// 	log.Fatal(app.Listen(":8080"))
// }

//running perfectly
// package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/Dhanush-varddhan-k/go-fiber-postgres/models"
// 	"github.com/Dhanush-varddhan-k/go-fiber-postgres/storage"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// 	"github.com/streadway/amqp"
// 	"gorm.io/gorm"
// )

// type Repository struct {
// 	DB *gorm.DB
// }

// // Converts a slice of strings to a JSON string
// func toJSONString(input []string) string {
// 	jsonData, err := json.Marshal(input)
// 	if err != nil {
// 		log.Printf("Error serializing to JSON: %v", err)
// 		return "[]"
// 	}
// 	return string(jsonData)
// }

// func setupRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
// 	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	_, err = ch.QueueDeclare(
// 		"image-processing-queue", // Queue name
// 		true,                     // Durable
// 		false,                    // Delete when unused
// 		false,                    // Exclusive
// 		false,                    // No-wait
// 		nil,                      // Arguments
// 	)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return conn, ch, nil
// }

// func (r *Repository) CreateProduct(context *fiber.Ctx) error {
// 	product := models.Product{}
// 	err := context.BodyParser(&product)
// 	if err != nil {
// 		return context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Invalid data"})
// 	}

// 	// Serialize ProductImages into JSON for database
// 	product.ProductImagesJSON = toJSONString(product.ProductImages)
// 	product.CompressedImagesJSON = "[]" // Placeholder for now

// 	// Save product to the database
// 	err = r.DB.Create(&product).Error
// 	if err != nil {
// 		log.Printf("Database error: %v", err)
// 		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not save product"})
// 	}

// 	// Publish product images to RabbitMQ for processing
// 	payload := map[string]interface{}{
// 		"product_id": product.ID,
// 		"images":     product.ProductImages,
// 	}
// 	message, _ := json.Marshal(payload)
// 	r.publishToQueue(message)

// 	return context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Product created successfully", "data": product})
// }

// func (r *Repository) publishToQueue(message []byte) {
// 	conn, ch, err := setupRabbitMQ()
// 	if err != nil {
// 		log.Printf("RabbitMQ connection error: %v", err)
// 		return
// 	}
// 	defer conn.Close()
// 	defer ch.Close()

// 	err = ch.Publish(
// 		"",                       // Exchange
// 		"image-processing-queue", // Queue name
// 		false,                    // Mandatory
// 		false,                    // Immediate
// 		amqp.Publishing{
// 			ContentType: "application/json",
// 			Body:        message,
// 		},
// 	)
// 	if err != nil {
// 		log.Printf("Failed to publish message: %v", err)
// 	}
// }

// // Function to update the database with compressed images
// func updateDatabaseWithCompressedImages(db *gorm.DB, productID uint, compressedImages []string) {
// 	// Check if compressedImages is empty
// 	if len(compressedImages) == 0 {
// 		log.Printf("No compressed images to update for product ID %d", productID)
// 		return
// 	}

// 	// Debug: Log the images to be updated
// 	log.Printf("Compressed images to be updated in DB for product ID %d: %v", productID, compressedImages)

// 	// Convert compressed images to JSON
// 	compressedImagesJSON, err := json.Marshal(compressedImages)
// 	if err != nil {
// 		log.Printf("Failed to marshal compressed images: %v", err)
// 		return
// 	}

// 	// Debug: Log the JSON payload
// 	log.Printf("Marshalled JSON for DB update: %s", string(compressedImagesJSON))

// 	// Perform the update in the database
// 	err = db.Model(&models.Product{}).Where("id = ?", productID).Update("compressed_images_json", string(compressedImagesJSON)).Error
// 	if err != nil {
// 		log.Printf("Failed to update database for product ID %d: %v", productID, err)
// 		return
// 	}

// 	log.Printf("Successfully updated database for product ID %d", productID)
// }

// func main() {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	config := &storage.Config{
// 		Host:     os.Getenv("DB_HOST"),
// 		Port:     os.Getenv("DB_PORT"),
// 		Password: os.Getenv("DB_PASS"),
// 		User:     os.Getenv("DB_USER"),
// 		SSLMode:  os.Getenv("DB_SSLMODE"),
// 		DBName:   os.Getenv("DB_NAME"),
// 	}

// 	db, err := storage.NewConnection(config)
// 	if err != nil {
// 		log.Fatal("Could not connect to the database")
// 	}

// 	err = models.Migrate(db)
// 	if err != nil {
// 		log.Fatal("Could not migrate database")
// 	}

// 	r := Repository{DB: db}
// 	app := fiber.New()

// 	// Endpoint to create a product
// 	app.Post("/api/products", r.CreateProduct)

// 	log.Fatal(app.Listen(":8080"))
// }

//working caching using redis

// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/Dhanush-varddhan-k/go-fiber-postgres/models"
// 	"github.com/Dhanush-varddhan-k/go-fiber-postgres/storage"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// 	"github.com/redis/go-redis/v9"
// 	"github.com/streadway/amqp"
// 	"gorm.io/gorm"
// )

// type Repository struct {
// 	DB    *gorm.DB
// 	Cache *redis.Client
// }

// // Converts a slice of strings to a JSON string
// func toJSONString(input []string) string {
// 	jsonData, err := json.Marshal(input)
// 	if err != nil {
// 		log.Printf("Error serializing to JSON: %v", err)
// 		return "[]"
// 	}
// 	return string(jsonData)
// }

// // RabbitMQ setup
// func setupRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
// 	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	_, err = ch.QueueDeclare(
// 		"image-processing-queue", // Queue name
// 		true,                     // Durable
// 		false,                    // Delete when unused
// 		false,                    // Exclusive
// 		false,                    // No-wait
// 		nil,                      // Arguments
// 	)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return conn, ch, nil
// }

// // Redis setup
// func setupRedis() *redis.Client {
// 	client := redis.NewClient(&redis.Options{
// 		Addr: "localhost:6379", // Redis server address
// 	})
// 	return client
// }

// // Create Product Endpoint
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
// 		log.Printf("Database error: %v", err)
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

// 	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "Product created successfully", "data": product})
// }

// // Update Product Endpoint
// func (r *Repository) UpdateProduct(ctx *fiber.Ctx) error {
// 	id := ctx.Params("id")
// 	product := models.Product{}

// 	err := ctx.BodyParser(&product)
// 	if err != nil {
// 		return ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Invalid data"})
// 	}

// 	// Update the database
// 	err = r.DB.Model(&models.Product{}).Where("id = ?", id).Updates(&product).Error
// 	if err != nil {
// 		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Failed to update product"})
// 	}

// 	// Invalidate cache for the specific product and all products
// 	r.Cache.Del(ctx.Context(), "product:"+id)
// 	r.Cache.Del(ctx.Context(), "all_products")

// 	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "Product updated successfully"})
// }

// // Get Product Endpoint
// func (r *Repository) GetProduct(ctx *fiber.Ctx) error {
// 	id := ctx.Params("id")
// 	cacheKey := "product:" + id

// 	// Check if the product is in cache
// 	cachedProduct, err := r.Cache.Get(context.Background(), cacheKey).Result()
// 	if err == nil {
// 		log.Printf("Cache hit for product ID %s", id)
// 		return ctx.Status(http.StatusOK).JSON(&fiber.Map{"data": cachedProduct, "message": "Product retrieved from cache"})
// 	}

// 	// If not in cache, retrieve from database
// 	product := models.Product{}
// 	err = r.DB.First(&product, id).Error
// 	if err != nil {
// 		return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "Product not found"})
// 	}

// 	// Serialize product to JSON and store in cache
// 	productJSON, _ := json.Marshal(product)
// 	r.Cache.Set(context.Background(), cacheKey, productJSON, 0) // Cache indefinitely

// 	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"data": product, "message": "Product retrieved from database"})
// }

// // Get All Products Endpoint
// func (r *Repository) GetAllProducts(ctx *fiber.Ctx) error {
// 	cacheKey := "all_products"

// 	// Check if all products are in cache
// 	cachedProducts, err := r.Cache.Get(context.Background(), cacheKey).Result()
// 	if err == nil {
// 		log.Println("Cache hit for all products")
// 		return ctx.Status(http.StatusOK).JSON(&fiber.Map{"data": cachedProducts, "message": "Products retrieved from cache"})
// 	}

// 	// If not in cache, retrieve from database
// 	var products []models.Product
// 	err = r.DB.Find(&products).Error
// 	if err != nil {
// 		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Failed to retrieve products"})
// 	}

// 	// Serialize products to JSON and store in cache
// 	productsJSON, _ := json.Marshal(products)
// 	r.Cache.Set(context.Background(), cacheKey, productsJSON, 0) // Cache indefinitely

// 	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"data": products, "message": "Products retrieved from database"})
// }

// // Publish to RabbitMQ
// func (r *Repository) publishToQueue(message []byte) {
// 	conn, ch, err := setupRabbitMQ()
// 	if err != nil {
// 		log.Printf("RabbitMQ connection error: %v", err)
// 		return
// 	}
// 	defer conn.Close()
// 	defer ch.Close()

// 	err = ch.Publish(
// 		"",                       // Exchange
// 		"image-processing-queue", // Queue name
// 		false,                    // Mandatory
// 		false,                    // Immediate
// 		amqp.Publishing{
// 			ContentType: "application/json",
// 			Body:        message,
// 		},
// 	)
// 	if err != nil {
// 		log.Printf("Failed to publish message: %v", err)
// 	}
// }

// // Main function
// func main() {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	config := &storage.Config{
// 		Host:     os.Getenv("DB_HOST"),
// 		Port:     os.Getenv("DB_PORT"),
// 		Password: os.Getenv("DB_PASS"),
// 		User:     os.Getenv("DB_USER"),
// 		SSLMode:  os.Getenv("DB_SSLMODE"),
// 		DBName:   os.Getenv("DB_NAME"),
// 	}

// 	db, err := storage.NewConnection(config)
// 	if err != nil {
// 		log.Fatal("Could not connect to the database")
// 	}

// 	err = models.Migrate(db)
// 	if err != nil {
// 		log.Fatal("Could not migrate database")
// 	}

// 	// Initialize Redis client
// 	cache := setupRedis()

// 	r := Repository{DB: db, Cache: cache}
// 	app := fiber.New()

// 	// Endpoints
// 	app.Post("/api/products", r.CreateProduct)
// 	app.Put("/api/products/:id", r.UpdateProduct)
// 	app.Get("/api/products/:id", r.GetProduct)
// 	app.Get("/api/products", r.GetAllProducts)

// 	log.Fatal(app.Listen(":8080"))
// }

// working till logger
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
