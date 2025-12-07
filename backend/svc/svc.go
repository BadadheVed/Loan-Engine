package svc

import (
	"log/slog"

	"github.com/BadadheVed/clickpe/database"
	"github.com/BadadheVed/clickpe/models"
)

func SaveUsersBatch(batch []models.User) (int, error) {
	slog.Info("SaveUsersBatch: Starting insert", "batch_size", len(batch))
	result := database.DB.Create(&batch)

	if result.Error != nil {
		slog.Error("SaveUsersBatch: Insert failed", "error", result.Error, "batch_size", len(batch))
		return 0, result.Error
	}

	slog.Info("SaveUsersBatch: Insert successful", "rows_affected", result.RowsAffected, "batch_size", len(batch))
	return int(result.RowsAffected), nil
}
