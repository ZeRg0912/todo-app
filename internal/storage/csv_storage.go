// Package storage provides persistence functionality for tasks
// in various formats including JSON and CSV.
package storage

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"todo-app/internal/todo"

	"github.com/ZeRg0912/logger"
)

// LoadCSV reads tasks from a CSV file with logging support.
// The CSV file should have a header row with columns: ID, Description, Done.
// Returns an empty task slice if the file has only a header or is empty.
// Returns an error if file reading or CSV parsing fails.
func LoadCSV(path string) ([]todo.Task, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %w", path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var tasks []todo.Task
	lineNum := 0
	skippedCount := 0

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			// Логируем ошибку чтения, но продолжаем парсинг остальных строк
			skippedCount++
			logger.Warn("CSV read error at line %d: %v", lineNum+1, err)
			continue
		}

		lineNum++

		// Пропускаем заголовок
		if lineNum == 1 {
			continue
		}

		// Проверяем количество полей
		if len(record) < 3 {
			skippedCount++
			logger.Warn("Skipping record at line %d: expected 3 fields, got %d", lineNum, len(record))
			continue
		}

		// Парсим ID
		id, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			skippedCount++
			logger.Warn("Skipping record at line %d: invalid ID format '%s'", lineNum, record[0])
			continue
		}

		// Парсим статус Done
		done, err := strconv.ParseBool(strings.TrimSpace(record[2]))
		if err != nil {
			skippedCount++
			logger.Warn("Skipping record at line %d: invalid Done format '%s'", lineNum, record[2])
			continue
		}

		task := todo.Task{
			ID:          id,
			Description: strings.TrimSpace(record[1]),
			Done:        done,
		}
		tasks = append(tasks, task)
	}

	if skippedCount > 0 {
		logger.Info("Loaded %d tasks from CSV, skipped %d invalid records", len(tasks), skippedCount)
	} else {
		logger.Info("Successfully loaded %d tasks from CSV", len(tasks))
	}

	return tasks, nil
}

// SaveCSV writes tasks to a CSV file with a header row and logging.
// The CSV format includes columns: ID, Description, Done.
// Returns an error if file creation or CSV writing fails.
func SaveCSV(path string, tasks []todo.Task) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0644)
	if err != nil {
		return fmt.Errorf("can't create file on path %s: %w", path, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	header := []string{"ID", "Description", "Done"}
	err = writer.Write(header)
	if err != nil {
		return fmt.Errorf("can't write CSV header: %w", err)
	}

	successCount := 0
	for _, task := range tasks {
		record := []string{
			strconv.Itoa(task.ID),
			task.Description,
			strconv.FormatBool(task.Done),
		}
		err := writer.Write(record)
		if err != nil {
			logger.Warn("Failed to write task ID %d: %v", task.ID, err)
			continue
		}
		successCount++
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("CSV flush error: %w", err)
	}

	if err := file.Sync(); err != nil {
		return fmt.Errorf("cannot sync CSV file %s: %w", path, err)
	}

	logger.Info("Successfully exported %d/%d tasks to CSV file: %s", successCount, len(tasks), path)
	return nil
}
