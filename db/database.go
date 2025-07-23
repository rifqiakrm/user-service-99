package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"user-service/model"
)

// InitDB initializes the SQLite database using a pure Go driver.
func InitDB(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	return db, nil
}
