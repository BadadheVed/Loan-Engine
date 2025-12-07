package main

import (
	"log/slog"

	"github.com/BadadheVed/clickpe/database"
	"github.com/BadadheVed/clickpe/router"
)

func main() {

	database.DBConnect()
	slog.Info("Databae Connected")
	r := router.SetupRouter()

	slog.Info("Router Initialized")

	err := r.Run(":8080")
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		return
	}
	slog.Info("Server started on port 8080")
}
