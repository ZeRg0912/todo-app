// Package storage provides persistence functionality for tasks
// in various formats including JSON and CSV.
package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"todo-app/internal/todo"
)

// LoadCSV reads tasks from a CSV file.
// The CSV file should have a header row with columns: ID, Description, Done.
// Returns an empty task slice if the file has only a header or is empty.
// Returns an error if file reading or CSV parsing fails.
func LoadCSV(path string) ([]todo.Task, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open file on path %s", path)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("can't read file on paht %s", path)
	}

	if len(records) <= 1 {
		return []todo.Task{}, nil
	}

	var tasks []todo.Task
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 3 {
			continue
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			continue
		}

		done, err := strconv.ParseBool(record[2])
		if err != nil {
			continue
		}
		task := todo.Task{
			ID:          id,
			Description: record[1],
			Done:        done,
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// SaveCSV writes tasks to a CSV file with a header row.
// The CSV format includes columns: ID, Description, Done.
// Returns an error if file creation or CSV writing fails.
func SaveCSV(path string, tasks []todo.Task) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("can't create file on path %s", path)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "Description", "Done"}
	err = writer.Write(header)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		record := []string{
			strconv.Itoa(task.ID),
			task.Description,
			strconv.FormatBool(task.Done),
		}
		err := writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}
