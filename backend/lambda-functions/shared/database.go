package shared

import (
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection using DATABASE_URL env variable
func InitDB() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		slog.Error("DATABASE_URL environment variable is not set")
		return &DBError{Message: "DATABASE_URL is required"}
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return err
	}

	// Set connection pool settings
	sqlDB, err := DB.DB()
	if err != nil {
		slog.Error("Failed to get database instance", "error", err)
		return err
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)

	// Auto-migrate all models
	err = DB.AutoMigrate(
		&User{},
		&LoanProduct{},
		&Match{},
	)
	if err != nil {
		slog.Error("Failed to auto-migrate tables", "error", err)
		return err
	}

	slog.Info("Database connected successfully")
	return nil
}

// DBError represents a database error
type DBError struct {
	Message string
}

func (e *DBError) Error() string {
	return e.Message
}
