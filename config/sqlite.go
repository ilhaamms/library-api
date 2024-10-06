package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
}

func InitDbSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("/app/db/library.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Exec("PRAGMA foreign_keys = ON;")

	return db, nil
}
