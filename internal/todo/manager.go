// Package todo provides task management functionality including
// CRUD operations, filtering, and import/export capabilities.
package todo

import (
	"fmt"
)

// Add creates a new task and appends it to the task list.
// Generates a unique ID by finding the maximum existing ID and incrementing it.
// Returns the updated task slice.
func Add(tasks []Task, desc string) []Task {
	newTask := Task{
		ID:          generateID(tasks),
		Description: desc,
		Done:        false,
	}
	return append(tasks, newTask)
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
// Returns an error if no task with the given ID is found.
// Returns the updated task slice on success.
func Complete(tasks []Task, id int) ([]Task, error) {
	index := findTaskByID(tasks, id)
	if index == -1 {
		return tasks, fmt.Errorf("task with ID %d not found", id)
	}
	tasks[index].Done = true
	return tasks, nil
}

// Delete removes a task from the list by its ID.
// Returns an error if no task with the given ID is found.
// Returns the updated task slice on success.
func Delete(tasks []Task, id int) ([]Task, error) {
	index := findTaskByID(tasks, id)
	if index == -1 {
		return tasks, fmt.Errorf("task with ID %d not found", id)
	}

	return append(tasks[:index], tasks[index+1:]...), nil
}

// generateID creates a new unique ID for a task.
// It finds the maximum ID in the existing tasks and increments it by 1.
// Returns 1 if the task list is empty.
func generateID(tasks []Task) int {
	if len(tasks) == 0 {
		return 1
	}

	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
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
