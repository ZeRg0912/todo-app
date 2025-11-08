// Package todo provides task management functionality including
// CRUD operations, filtering, and import/export capabilities.
package todo

// Task represents a single todo item in the system.
// ID is a unique auto-generated identifier.
// Description contains the task text content.
// Done indicates whether the task has been completed.
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}
