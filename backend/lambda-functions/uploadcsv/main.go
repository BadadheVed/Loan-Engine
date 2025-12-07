package main

import (
	"context"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"mime/multipart"
	"strconv"
	"strings"
	"sync"

	"github.com/BadadheVed/clickpe/lambda-functions/shared"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

func init() {
	// Initialize database connection on cold start
	if err := shared.InitDB(); err != nil {
		panic(err)
	}
}

func saveUsersBatch(batch []shared.User) (int, error) {
	slog.Info("SaveUsersBatch: Starting insert", "batch_size", len(batch))
	result := shared.DB.Create(&batch)

	if result.Error != nil {
		slog.Error("SaveUsersBatch: Insert failed", "error", result.Error, "batch_size", len(batch))
		return 0, result.Error
	}

	slog.Info("SaveUsersBatch: Insert successful", "rows_affected", result.RowsAffected, "batch_size", len(batch))
	return int(result.RowsAffected), nil
}

func userWorker(id int, jobs <-chan []shared.User, results chan<- shared.BatchResult, wg *sync.WaitGroup) {
	defer slog.Info("Worker finished", "worker_id", id)
	defer wg.Done()

	batchCount := 0
	for batch := range jobs {
		batchCount++
		slog.Info("Worker processing batch", "worker_id", id, "batch_num", batchCount, "batch_size", len(batch))
		inserted, err := saveUsersBatch(batch)

		if err != nil {
			slog.Error("Worker batch failed", "worker_id", id, "batch_num", batchCount, "error", err)
			results <- shared.BatchResult{Inserted: 0, Attempted: len(batch)}
		} else {
			slog.Info("Worker batch completed", "worker_id", id, "batch_num", batchCount, "inserted", inserted, "attempted", len(batch))
			results <- shared.BatchResult{Inserted: inserted, Attempted: len(batch)}
		}
	}
}

func processCSV(fileContent []byte) (map[string]interface{}, error) {
	reader := csv.NewReader(strings.NewReader(string(fileContent)))

	// Skip header
	_, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("invalid CSV: %w", err)
	}

	const workerCount = 5
	const channelBufferSize = 100
	jobs := make(chan []shared.User, channelBufferSize)
	results := make(chan shared.BatchResult, channelBufferSize)
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		slog.Info("Starting worker", "id", i)
		wg.Add(1)
		go userWorker(i, jobs, results, &wg)
	}

	var (
		batchSize      = 100
		batch          []shared.User
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

		user := shared.User{
			ID:            id,
			Name:          row[1],
			Email:         row[2],
			Age:           age,
			CreditScore:   creditScore,
			MonthlyIncome: income,
		}

		batch = append(batch, user)

		if len(batch) >= batchSize {
			batchesSent++
			slog.Info("Sending batch to jobs channel", "batch_num", batchesSent, "batch_size", len(batch), "total_rows_read", totalRowsRead)
			jobs <- batch
			batch = []shared.User{}
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

	return map[string]interface{}{
		"records_added":         addedCount,
		"records_failed":        failedCount,
		"records_skipped":       skippedCount,
		"duplicate_email_count": duplicateCount,
	}, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse multipart form data
	contentType := request.Headers["content-type"]
	if contentType == "" {
		contentType = request.Headers["Content-Type"]
	}

	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Invalid content type"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	boundary := params["boundary"]
	if boundary == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Missing boundary in content type"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Decode base64 body
	var body []byte
	if request.IsBase64Encoded {
		body, err = base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       `{"error": "Failed to decode base64 body"}`,
				Headers:    map[string]string{"Content-Type": "application/json"},
			}, nil
		}
	} else {
		body = []byte(request.Body)
	}

	// Parse multipart form
	mr := multipart.NewReader(strings.NewReader(string(body)), boundary)
	var fileContent []byte

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       `{"error": "Failed to parse multipart form"}`,
				Headers:    map[string]string{"Content-Type": "application/json"},
			}, nil
		}

		if part.FormName() == "file" {
			fileContent, err = io.ReadAll(part)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       `{"error": "Failed to read file content"}`,
					Headers:    map[string]string{"Content-Type": "application/json"},
				}, nil
			}
			break
		}
	}

	if len(fileContent) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "CSV file is required"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	slog.Info("Got the file", "size", len(fileContent))

	// Process CSV
	result, err := processCSV(fileContent)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"error": "Failed to process CSV: %s"}`, err.Error()),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	responseBody, err := json.Marshal(result)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to marshal response"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
