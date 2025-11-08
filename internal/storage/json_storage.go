// Package storage provides persistence functionality for tasks
// in various formats including JSON and CSV.
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"todo-app/internal/todo"

	"github.com/ZeRg0912/logger"
)

// LoadJSON reads tasks from a JSON file with logging.
// Returns an empty task slice if the file doesn't exist or is empty.
// Returns an error if file reading or JSON parsing fails.
func LoadJSON(path string) ([]todo.Task, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		logger.Info("JSON file %s does not exist, returning empty task list", path)
		return []todo.Task{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("unexpected error accessing path %s: %w", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %s: %w", path, err)
	}

	if len(data) == 0 {
		logger.Info("JSON file %s is empty, returning empty task list", path)
		return []todo.Task{}, nil
	}

	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
		logger.Debug("Removed UTF-8 BOM from JSON file")
	}

	var tasks []todo.Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("cannot parse JSON from %s: %w", path, err)
	}

	logger.Info("Successfully loaded %d tasks from JSON file: %s", len(tasks), path)
	return tasks, nil
}

// SaveJSON writes tasks to a JSON file with indentation and logging.
// Uses atomic write (temp file + rename) to protect data from corruption.
// Uses file locking to prevent concurrent access conflicts.
// Returns an error if JSON marshaling or file writing fails.
func SaveJSON(path string, tasks []todo.Task) error {
	lock, err := AcquireLock(path)
	if err != nil {
		return fmt.Errorf("cannot acquire lock for %s: %w", path, err)
	}
	defer lock.Release()

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal tasks to JSON: %w", err)
	}

	dir := filepath.Dir(path)
	if dir == "." {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("cannot get absolute path for %s: %w", path, err)
		}
		dir = filepath.Dir(absPath)
	}
	tmpFile, err := os.CreateTemp(dir, filepath.Base(path)+".tmp.*")
	if err != nil {
		return fmt.Errorf("cannot create temporary file for %s: %w", path, err)
	}
	tmpPath := tmpFile.Name()

	defer func() {
		tmpFile.Close()
		if _, err := os.Stat(tmpPath); err == nil {
			os.Remove(tmpPath)
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		return fmt.Errorf("cannot write to temporary file %s: %w", tmpPath, err)
	}

	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("cannot sync temporary file %s: %w", tmpPath, err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("cannot close temporary file %s: %w", tmpPath, err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("cannot rename temporary file to %s: %w", path, err)
	}

	logger.Info("Successfully saved %d tasks to JSON file: %s", len(tasks), path)
	return nil
}
