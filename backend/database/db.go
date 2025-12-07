package database

import (
	"log/slog"
	"os"

	"github.com/BadadheVed/clickpe/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() {
	slog.Info("Loading ENV")
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", err)
	}
	slog.Info("--------ENV loaded successfully-------")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		slog.Error("DATABASE_URL is not set in environment variables")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database:", err)
	}

	DB = db
	slog.Info("Database connected successfully")

	err = db.AutoMigrate(
		&models.User{},
		&models.LoanProduct{},
		&models.Match{},
	)

	if err != nil {
		slog.Error("Failed to auto-migrate tables:", err)
	}

	slog.Info("Tables migrated successfully")
}
