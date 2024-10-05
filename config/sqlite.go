package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	// Db *gorm.DB
}

func InitDbSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("/app/db/library.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
