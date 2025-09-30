package db

import (
	"log"

	"gorm.io/driver/sqlite"

	"github.com/fazel/notebooq/internal/config"
	"gorm.io/gorm"
)

func Open(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Printf("[error] failed to initialize database, got error %v", err)
		return nil, err
	}

	return db, nil
}
