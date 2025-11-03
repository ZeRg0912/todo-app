// Package storage provides persistence functionality for tasks
// in various formats including JSON and CSV.
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"todo-app/internal/todo"
)

// LoadJSON reads tasks from a JSON file.
// Returns an empty task slice if the file doesn't exist or is empty.
// Returns an error if file reading or JSON parsing fails.
func LoadJSON(path string) ([]todo.Task, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return []todo.Task{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("ERROR: unexcepted error on path %s ", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ERROR: can't read file on path %s", path)
	}

	if len(data) == 0 {
		return []todo.Task{}, nil
	}

	var tasks []todo.Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("ERROR: can't parsing JSON on path %s", path)
	}

	return tasks, nil
}

// SaveJSON writes tasks to a JSON file with indentation for readability.
// Returns an error if JSON marshaling or file writing fails.
func SaveJSON(path string, tasks []todo.Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return fmt.Errorf("ERROR: can't parse JSON")
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("ERROR: can't save to path %s", path)
	}

	return nil
}
