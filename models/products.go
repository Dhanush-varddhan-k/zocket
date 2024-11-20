// package models

// import "gorm.io/gorm"

// type Product struct {
// 	ID                      uint     `gorm:"primaryKey;autoIncrement" json:"id"`
// 	UserID                  uint     `json:"user_id"`
// 	ProductName             string   `json:"product_name"`
// 	ProductDescription      string   `json:"product_description"`
// 	ProductImages           []string `gorm:"type:text[]" json:"product_images"`
// 	ProductPrice            float64  `json:"product_price"`
// 	CompressedProductImages []string `gorm:"type:text[]" json:"compressed_product_images"`
// }

// func Migrate(db *gorm.DB) error {
// 	err := db.AutoMigrate(&Product{})
// 	return err
// }

// package models

// import "gorm.io/gorm"

// type Product struct {
// 	ID                      uint     `gorm:"primaryKey;autoIncrement" json:"id"`
// 	UserID                  uint     `json:"user_id"`
// 	ProductName             string   `json:"product_name"`
// 	ProductDescription      string   `json:"product_description"`
// 	ProductImages           []string `gorm:"-" json:"product_images"` // Use slice for JSON unmarshalling
// 	ProductPrice            float64  `json:"product_price"`
// 	CompressedProductImages []string `gorm:"-" json:"compressed_product_images"`
// 	ProductImagesJSON       string   `gorm:"type:json" json:"-"` // Store as JSON in DB
// 	CompressedImagesJSON    string   `gorm:"type:json" json:"-"` // Store as JSON in DB
// }

//	func Migrate(db *gorm.DB) error {
//		db.Migrator().DropTable(&Product{})
//		err := db.AutoMigrate(&Product{})
//		return err
//	}
package models

import "gorm.io/gorm"

// Product struct with proper JSON handling
type Product struct {
	ID                      uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID                  uint     `json:"user_id"`
	ProductName             string   `json:"product_name"`
	ProductDescription      string   `json:"product_description"`
	ProductImages           []string `gorm:"-" json:"product_images"`            // Ignore in DB, used for API
	CompressedProductImages []string `gorm:"-" json:"compressed_product_images"` // Ignore in DB
	ProductImagesJSON       string   `gorm:"type:jsonb" json:"-"`                // Stored in DB
	CompressedImagesJSON    string   `gorm:"type:jsonb" json:"-"`                // Stored in DB
}

// Migrate ensures the table is created properly
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Product{})
}
