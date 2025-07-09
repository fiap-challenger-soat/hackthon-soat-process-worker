package postgres

import (
	"fmt"
	"time"

	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresClient() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Vars.DBHost, config.Vars.DBPort, config.Vars.DBUser, config.Vars.DBPassword, config.Vars.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(config.Vars.DbMaxIdleConns)
	sqlDB.SetMaxOpenConns(config.Vars.DbMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
