package db

import (
	"github.com/fazel/notebooq/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(cfg *config.Config) (*gorm.DB, error) {

	return gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
}
