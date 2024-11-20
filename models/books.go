package models

import "gorm.io/gorm"

type Books struct {
	ID        uint    `gorm:"primary key;autoIncrement" json:"id"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
}

// i d would be createed gorm but others we should create so no *
// in mongobb easy to createdb. In postgres we need to have db so we do automigration.
func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Books{})
	return err
}
