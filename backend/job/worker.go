package job

import (
	"log/slog"
	"sync"

	"github.com/BadadheVed/clickpe/models"
	"github.com/BadadheVed/clickpe/svc"
)

type BatchResult struct {
	Inserted  int
	Attempted int
}

func UserWorker(id int, jobs <-chan []models.User, results chan<- BatchResult, wg *sync.WaitGroup) {
	defer slog.Info("Worker finished", "worker_id", id)
	defer wg.Done()

	batchCount := 0
	for batch := range jobs {
		batchCount++
		slog.Info("Worker processing batch", "worker_id", id, "batch_num", batchCount, "batch_size", len(batch))
		inserted, err := svc.SaveUsersBatch(batch)

		if err != nil {
			slog.Error("Worker batch failed", "worker_id", id, "batch_num", batchCount, "error", err)
			results <- BatchResult{Inserted: 0, Attempted: len(batch)}
		} else {
			slog.Info("Worker batch completed", "worker_id", id, "batch_num", batchCount, "inserted", inserted, "attempted", len(batch))
			results <- BatchResult{Inserted: inserted, Attempted: len(batch)}
		}
	}
}
