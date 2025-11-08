// Package todo provides task management functionality including
// CRUD operations, filtering, and import/export capabilities.
package todo

import (
	"fmt"
)

const (
	MinID                = 1
	MaxDescriptionLength = 1000
)

// Add creates a new task and appends it to the task list.
// Generates a unique ID by finding the maximum existing ID and incrementing it.
// Returns an error if description validation fails.
// Returns the updated task slice on success.
func Add(tasks []Task, desc string) ([]Task, error) {
	if err := ValidateDescription(desc); err != nil {
		return tasks, err
	}
	newTask := Task{
		ID:          generateID(tasks),
		Description: desc,
		Done:        false,
	}
	return append(tasks, newTask), nil
}

// List filters tasks based on the specified criteria.
// Supported filters: "all", "done", "pending".
// Returns a slice containing only tasks that match the filter.
func List(tasks []Task, filter string) []Task {
	switch filter {
	case "done":
		var result []Task
		for _, task := range tasks {
			if task.Done {
				result = append(result, task)
			}
		}
		return result
	case "pending":
		var result []Task
		for _, task := range tasks {
			if !task.Done {
				result = append(result, task)
			}
		}
		return result
	case "all":
		return tasks
	default:
		return tasks
	}
}

// Complete marks a task as done by its ID.
// Returns an error if ID is invalid or no task with the given ID is found.
// Returns the updated task slice on success.
func Complete(tasks []Task, id int) ([]Task, error) {
	if err := ValidateID(id); err != nil {
		return tasks, err
	}
	index := findTaskByID(tasks, id)
	if index == -1 {
		return tasks, fmt.Errorf("task with ID %d not found", id)
	}
	tasks[index].Done = true
	return tasks, nil
}

// Delete removes a task from the list by its ID.
// Returns an error if ID is invalid or no task with the given ID is found.
// Returns the updated task slice on success.
func Delete(tasks []Task, id int) ([]Task, error) {
	if err := ValidateID(id); err != nil {
		return tasks, err
	}
	index := findTaskByID(tasks, id)
	if index == -1 {
		return tasks, fmt.Errorf("task with ID %d not found", id)
	}

	return append(tasks[:index], tasks[index+1:]...), nil
}

// generateID creates a new unique ID for a task.
// It finds the maximum ID in the existing tasks and increments it by 1.
// Returns 1 if the task list is empty.
// Optimized: uses single pass through tasks with early exit optimization.
func generateID(tasks []Task) int {
	if len(tasks) == 0 {
		return MinID
	}

	maxID := MinID - 1
	for i := range tasks {
		if tasks[i].ID > maxID {
			maxID = tasks[i].ID
		}
	}
	return maxID + 1
}

// ValidateID validates that a task ID is within acceptable range.
// Returns an error if ID is less than MinID.
func ValidateID(id int) error {
	if id < MinID {
		return fmt.Errorf("task ID must be at least %d, got %d", MinID, id)
	}
	return nil
}

// ValidateDescription validates that a task description is within acceptable limits.
// Returns an error if description is empty or exceeds MaxDescriptionLength.
func ValidateDescription(desc string) error {
	if desc == "" {
		return fmt.Errorf("task description cannot be empty")
	}
	if len(desc) > MaxDescriptionLength {
		return fmt.Errorf("task description cannot exceed %d characters, got %d", MaxDescriptionLength, len(desc))
	}
	return nil
}

// findTaskByID searches for a task by its ID in the task slice.
// Returns the index of the task if found, or -1 if not found.
func findTaskByID(tasks []Task, id int) int {
	for i := range tasks {
		if tasks[i].ID == id {
			return i
		}
	}
	return -1
}