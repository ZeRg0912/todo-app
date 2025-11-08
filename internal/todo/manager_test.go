package todo

import (
	"testing"
)

func TestAdd(t *testing.T) {
	tasks := []Task{}

	// Test adding first task
	var err error
	tasks, err = Add(tasks, "First task")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	if tasks[0].ID != 1 {
		t.Errorf("Expected ID 1, got %d", tasks[0].ID)
	}
	if tasks[0].Description != "First task" {
		t.Errorf("Expected description 'First task', got '%s'", tasks[0].Description)
	}
	if tasks[0].Done {
		t.Error("New task should not be done")
	}

	// Test adding second task
	tasks, err = Add(tasks, "Second task")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
	if tasks[1].ID != 2 {
		t.Errorf("Expected ID 2, got %d", tasks[1].ID)
	}
}

func TestList(t *testing.T) {
	tasks := []Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: true},
		{ID: 3, Description: "Task 3", Done: false},
	}

	// Test "all" filter
	allTasks := List(tasks, "all")
	if len(allTasks) != 3 {
		t.Errorf("Expected 3 tasks for 'all' filter, got %d", len(allTasks))
	}

	// Test "done" filter
	doneTasks := List(tasks, "done")
	if len(doneTasks) != 1 {
		t.Errorf("Expected 1 task for 'done' filter, got %d", len(doneTasks))
	}
	if !doneTasks[0].Done {
		t.Error("Done filter should return only done tasks")
	}

	// Test "pending" filter
	pendingTasks := List(tasks, "pending")
	if len(pendingTasks) != 2 {
		t.Errorf("Expected 2 tasks for 'pending' filter, got %d", len(pendingTasks))
	}
	if pendingTasks[0].Done || pendingTasks[1].Done {
		t.Error("Pending filter should return only not done tasks")
	}

	// Test unknown filter (should return all)
	unknownTasks := List(tasks, "unknown")
	if len(unknownTasks) != 3 {
		t.Errorf("Unknown filter should return all tasks, got %d", len(unknownTasks))
	}
}

func TestComplete(t *testing.T) {
	tasks := []Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: false},
	}

	// Test completing existing task
	result, err := Complete(tasks, 1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result[0].Done {
		t.Error("Task should be marked as done")
	}
	if result[1].Done {
		t.Error("Other task should not be affected")
	}

	// Test completing non-existing task
	_, err = Complete(tasks, 999)
	if err == nil {
		t.Error("Expected error for non-existing task")
	}
}

func TestDelete(t *testing.T) {
	tasks := []Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: false},
		{ID: 3, Description: "Task 3", Done: false},
	}

	// Test deleting middle task
	result, err := Delete(tasks, 2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 tasks after deletion, got %d", len(result))
	}
	if result[0].ID != 1 || result[1].ID != 3 {
		t.Error("Wrong tasks remaining after deletion")
	}

	// Test deleting non-existing task
	_, err = Delete(tasks, 999)
	if err == nil {
		t.Error("Expected error for non-existing task")
	}
}

func TestGenerateID(t *testing.T) {
	// Test empty tasks
	emptyTasks := []Task{}
	if id := generateID(emptyTasks); id != 1 {
		t.Errorf("Expected ID 1 for empty tasks, got %d", id)
	}

	// Test with existing tasks
	tasks := []Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 5, Description: "Task 5", Done: false}, // Gap in IDs
		{ID: 3, Description: "Task 3", Done: false},
	}
	if id := generateID(tasks); id != 6 {
		t.Errorf("Expected ID 6 (max+1), got %d", id)
	}
}

func TestCompleteEdgeCases(t *testing.T) {
	tasks := []Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: true}, // Уже выполнена
	}

	// Тест: выполнение уже выполненной задачи
	result, err := Complete(tasks, 2)
	if err != nil {
		t.Errorf("Should not error when completing already done task: %v", err)
	}
	if !result[1].Done {
		t.Error("Task should remain done")
	}

	// Тест: несуществующий ID
	_, err = Complete(tasks, 999)
	if err == nil {
		t.Error("Expected error for non-existing task ID")
	}
}

func TestDeleteEdgeCases(t *testing.T) {
	tasks := []Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: true},
	}

	// Тест: удаление несуществующей задачи
	_, err := Delete(tasks, 999)
	if err == nil {
		t.Error("Expected error for non-existing task ID")
	}

	// Тест: удаление из пустого списка
	_, err = Delete([]Task{}, 1)
	if err == nil {
		t.Error("Expected error when deleting from empty list")
	}
}

func TestValidateID(t *testing.T) {
	// Тест: валидный ID
	if err := ValidateID(1); err != nil {
		t.Errorf("Expected no error for valid ID 1, got %v", err)
	}

	// Тест: валидный ID больше MinID
	if err := ValidateID(100); err != nil {
		t.Errorf("Expected no error for valid ID 100, got %v", err)
	}

	// Тест: невалидный ID (меньше MinID)
	if err := ValidateID(0); err == nil {
		t.Error("Expected error for ID 0")
	}

	// Тест: невалидный ID (отрицательный)
	if err := ValidateID(-1); err == nil {
		t.Error("Expected error for negative ID")
	}
}

func TestValidateDescription(t *testing.T) {
	// Тест: валидное описание
	if err := ValidateDescription("Valid task description"); err != nil {
		t.Errorf("Expected no error for valid description, got %v", err)
	}

	// Тест: пустое описание
	if err := ValidateDescription(""); err == nil {
		t.Error("Expected error for empty description")
	}

	// Тест: описание на границе максимальной длины
	maxDesc := string(make([]byte, MaxDescriptionLength))
	if err := ValidateDescription(maxDesc); err != nil {
		t.Errorf("Expected no error for description at max length, got %v", err)
	}

	// Тест: описание превышает максимальную длину
	tooLongDesc := string(make([]byte, MaxDescriptionLength+1))
	if err := ValidateDescription(tooLongDesc); err == nil {
		t.Error("Expected error for description exceeding max length")
	}
}

func TestAddValidation(t *testing.T) {
	tasks := []Task{}

	// Тест: добавление с пустым описанием
	_, err := Add(tasks, "")
	if err == nil {
		t.Error("Expected error for empty description")
	}

	// Тест: добавление с описанием превышающим максимальную длину
	tooLongDesc := string(make([]byte, MaxDescriptionLength+1))
	_, err = Add(tasks, tooLongDesc)
	if err == nil {
		t.Error("Expected error for description exceeding max length")
	}
}

func TestCompleteValidation(t *testing.T) {
	tasks := []Task{
		{ID: 1, Description: "Task 1", Done: false},
	}

	// Тест: валидный ID
	_, err := Complete(tasks, 1)
	if err != nil {
		t.Errorf("Unexpected error for valid ID: %v", err)
	}

	// Тест: невалидный ID (0)
	_, err = Complete(tasks, 0)
	if err == nil {
		t.Error("Expected error for ID 0")
	}

	// Тест: невалидный ID (отрицательный)
	_, err = Complete(tasks, -1)
	if err == nil {
		t.Error("Expected error for negative ID")
	}
}

func TestDeleteValidation(t *testing.T) {
	tasks := []Task{
		{ID: 1, Description: "Task 1", Done: false},
	}

	// Тест: валидный ID
	_, err := Delete(tasks, 1)
	if err != nil {
		t.Errorf("Unexpected error for valid ID: %v", err)
	}

	// Тест: невалидный ID (0)
	_, err = Delete(tasks, 0)
	if err == nil {
		t.Error("Expected error for ID 0")
	}

	// Тест: невалидный ID (отрицательный)
	_, err = Delete(tasks, -1)
	if err == nil {
		t.Error("Expected error for negative ID")
	}
}