
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv" // Import this package for .env loading
	"github.com/streadway/amqp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Payload struct {
	ProductID uint     `json:"product_id"`
	Images    []string `json:"images"`
}

func main() {
	// Load .env file
	err := godotenv.Load("../.env") // Use relative path to load the .env file
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// RabbitMQ setup
	conn, ch, err := setupRabbitMQ()
	if err != nil {
		log.Fatalf("RabbitMQ connection failed: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	// Consume messages from RabbitMQ queue
	msgs, err := ch.Consume("image-processing-queue", "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Set up AWS S3 client
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}
	s3Client := s3.New(sess)

	// Set up PostgreSQL connection
	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Process messages from RabbitMQ
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var payload Payload
			err := json.Unmarshal(d.Body, &payload)
			if err != nil {
				log.Printf("Failed to unmarshal RabbitMQ message: %v", err)
				continue
			}

			// Process and upload images to S3
			compressedImages := processAndUploadImages(s3Client, payload.Images)

			// Update database with compressed image URLs
			updateDatabaseWithCompressedImages(db, payload.ProductID, compressedImages)
		}
	}()
	log.Println("Image processor is running...")
	<-forever
}

func setupDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback for individual credentials from .env
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_SSLMODE"),
		)
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func processAndUploadImages(s3Client *s3.S3, images []string) []string {
	var compressedImages []string
	for _, imageURL := range images {
		filename := filepath.Base(imageURL)
		localPath := fmt.Sprintf("./compressed/%s", filename)

		// Download the image
		err := downloadImage(imageURL, localPath)
		if err != nil {
			log.Printf("Failed to download image: %v", err)
			continue
		}

		// Upload the compressed image to S3
		s3URL := uploadToS3(s3Client, localPath, filename)
		if s3URL != "" {
			compressedImages = append(compressedImages, s3URL)
		}
	}
	return compressedImages
}

func uploadToS3(s3Client *s3.S3, localPath, filename string) string {
	file, err := os.Open(localPath)
	if err != nil {
		log.Printf("Failed to open file for upload: %v", err)
		return ""
	}
	defer file.Close()

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("zocket-dhan"),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		log.Printf("S3 upload failed for %s: %v", filename, err)
		return ""
	}
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", "zocket-dhan", filename)
}

func downloadImage(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func updateDatabaseWithCompressedImages(db *gorm.DB, productID uint, compressedImages []string) {
	// Check if the compressedImages array contains valid URLs
	log.Printf("Compressed images to be updated in DB: %v", compressedImages)

	// Convert the compressed images array to JSON
	compressedImagesJSON, err := json.Marshal(compressedImages)
	if err != nil {
		log.Printf("Failed to marshal compressed images: %v", err)
		return
	}

	// Debug: Log the marshaled JSON
	log.Printf("Marshalled JSON for DB update: %s", string(compressedImagesJSON))

	// Perform the database update
	err = db.Model(&Product{}).Where("id = ?", productID).Update("compressed_images_json", string(compressedImagesJSON)).Error
	if err != nil {
		log.Printf("Failed to update database for product ID %d: %v", productID, err)
		return
	}

	log.Printf("Successfully updated database for product ID %d", productID)
}

func setupRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	return conn, ch, nil
}

type Product struct {
	ID                   uint   `gorm:"primaryKey"`
	CompressedImagesJSON string `gorm:"type:jsonb"`
}
