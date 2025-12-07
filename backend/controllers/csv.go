package controllers

import (
	"bufio"
	"encoding/csv"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/BadadheVed/clickpe/job"
	"github.com/BadadheVed/clickpe/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadCSVUsers(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV file is required"})
		return
	}
	slog.Info("got the file with the size", file.Size)
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	reader := csv.NewReader(bufio.NewReader(f))
	// Skip header
	_, err = reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CSV"})
		return
	}

	const workerCount = 5
	const channelBufferSize = 100 // Large enough to hold all potential batches
	jobs := make(chan []models.User, channelBufferSize)
	results := make(chan job.BatchResult, channelBufferSize)
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		slog.Info("Starting worker", "id", i)
		wg.Add(1)
		go job.UserWorker(i, jobs, results, &wg)
	}

	var (
		batchSize      = 100
		batch          []models.User
		addedCount     int
		failedCount    int
		skippedCount   int
		duplicateCount int
		totalRowsRead  int
		batchesSent    int
	)

	slog.Info("Starting CSV processing")

	for {
		row, err := reader.Read()
		if err == io.EOF {
			slog.Info("Reached EOF", "total_rows_read", totalRowsRead, "batches_sent", batchesSent)
			break
		}

		if err != nil {
			slog.Warn("Error reading CSV row", "error", err, "row_num", totalRowsRead)
			failedCount++
			continue
		}

		totalRowsRead++

		for i := range row {
			row[i] = strings.TrimSpace(row[i])
		}

		if row[2] == "" {
			skippedCount++
			continue
		}

		id, err := uuid.Parse(row[0])
		if err != nil {
			skippedCount++
			continue
		}

		age, _ := strconv.Atoi(row[6])
		creditScore, _ := strconv.Atoi(row[4])
		income, _ := strconv.ParseFloat(row[3], 64)
		employmentStatus := row[5]
		user := models.User{
			ID:               id,
			Name:             row[1],
			Email:            row[2],
			Age:              age,
			CreditScore:      creditScore,
			MonthlyIncome:    income,
			EmploymentStatus: employmentStatus,
		}

		batch = append(batch, user)

		if len(batch) >= batchSize {
			batchesSent++
			slog.Info("Sending batch to jobs channel", "batch_num", batchesSent, "batch_size", len(batch), "total_rows_read", totalRowsRead)
			jobs <- batch
			batch = []models.User{}
		}
	}

	if len(batch) > 0 {
		batchesSent++
		slog.Info("Sending final batch to jobs channel", "batch_num", batchesSent, "batch_size", len(batch), "total_rows_read", totalRowsRead)
		jobs <- batch
	}
	slog.Info("All batches sent to jobs channel")

	close(jobs)
	slog.Info("Jobs channel closed, waiting for workers to finish")
	wg.Wait()
	slog.Info("All workers finished, closing results channel")
	close(results)

	resultCount := 0
	for r := range results {
		resultCount++
		slog.Info("Processing result", "result_num", resultCount, "inserted", r.Inserted, "attempted", r.Attempted)
		addedCount += r.Inserted
		duplicatesInBatch := r.Attempted - r.Inserted
		if duplicatesInBatch > 0 {
			duplicateCount += duplicatesInBatch
		}
	}
	slog.Info("All results processed", "total_results", resultCount, "total_added", addedCount)

	c.JSON(http.StatusCreated, gin.H{
		"records_added":         addedCount,
		"records_failed":        failedCount,
		"records_skipped":       skippedCount,
		"duplicate_email_count": duplicateCount,
	})
}
