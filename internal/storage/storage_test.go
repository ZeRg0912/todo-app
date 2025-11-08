package storage

import (
	"os"
	"testing"
	"todo-app/internal/todo"
)

func TestJSONSaveAndLoad(t *testing.T) {
	testFile := "test_tasks.json"
	defer os.Remove(testFile) // Cleanup after test

	tasks := []todo.Task{
		{ID: 1, Description: "Test task 1", Done: false},
		{ID: 2, Description: "Test task 2", Done: true},
	}

	// Test SaveJSON
	err := SaveJSON(testFile, tasks)
	if err != nil {
		t.Fatalf("SaveJSON failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("JSON file was not created")
	}

	// Test LoadJSON
	loaded, err := LoadJSON(testFile)
	if err != nil {
		t.Fatalf("LoadJSON failed: %v", err)
	}

	// Verify data integrity
	if len(loaded) != len(tasks) {
		t.Fatalf("Expected %d tasks, got %d", len(tasks), len(loaded))
	}

	for i, task := range loaded {
		if task.ID != tasks[i].ID {
			t.Errorf("Task %d: ID mismatch, expected %d, got %d", i, tasks[i].ID, task.ID)
		}
		if task.Description != tasks[i].Description {
			t.Errorf("Task %d: Description mismatch, expected '%s', got '%s'", i, tasks[i].Description, task.Description)
		}
		if task.Done != tasks[i].Done {
			t.Errorf("Task %d: Done mismatch, expected %t, got %t", i, tasks[i].Done, task.Done)
		}
	}
}

func TestJSONLoadNonExistentFile(t *testing.T) {
	// Test loading non-existent file
	loaded, err := LoadJSON("non_existent_file.json")
	if err != nil {
		t.Fatalf("LoadJSON should not fail for non-existent file: %v", err)
	}
	if len(loaded) != 0 {
		t.Errorf("Expected empty slice for non-existent file, got %d tasks", len(loaded))
	}
}

func TestJSONLoadEmptyFile(t *testing.T) {
	testFile := "empty_test.json"
	defer os.Remove(testFile)

	// Create empty file
	os.WriteFile(testFile, []byte{}, 0644)

	loaded, err := LoadJSON(testFile)
	if err != nil {
		t.Fatalf("LoadJSON failed for empty file: %v", err)
	}
	if len(loaded) != 0 {
		t.Errorf("Expected empty slice for empty file, got %d tasks", len(loaded))
	}
}

func TestJSONWithSpecialCharacters(t *testing.T) {
	testFile := "unicode_test.json"
	defer os.Remove(testFile)

	tasks := []todo.Task{
		{ID: 1, Description: "Задача с русскими буквами", Done: false},
		{ID: 2, Description: "Task with emoji 🚀 and symbols ©®", Done: true},
		{ID: 3, Description: "Task with \t tabs and \n newlines", Done: false},
	}

	err := SaveJSON(testFile, tasks)
	if err != nil {
		t.Fatalf("SaveJSON failed: %v", err)
	}

	loaded, err := LoadJSON(testFile)
	if err != nil {
		t.Fatalf("LoadJSON failed: %v", err)
	}

	if len(loaded) != len(tasks) {
		t.Fatalf("Expected %d tasks, got %d", len(tasks), len(loaded))
	}

	// Проверяем сохранение спецсимволов
	if loaded[0].Description != "Задача с русскими буквами" {
		t.Errorf("Russian text not preserved")
	}
	if loaded[1].Description != "Task with emoji 🚀 and symbols ©®" {
		t.Errorf("Emoji and symbols not preserved")
	}
}

func TestCSVSaveAndLoad(t *testing.T) {
	testFile := "test_tasks.csv"
	defer os.Remove(testFile) // Cleanup after test

	tasks := []todo.Task{
		{ID: 1, Description: "Test task 1", Done: false},
		{ID: 2, Description: "Test task 2", Done: true},
	}

	// Test SaveCSV
	err := SaveCSV(testFile, tasks)
	if err != nil {
		t.Fatalf("SaveCSV failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("CSV file was not created")
	}

	// Test LoadCSV
	loaded, err := LoadCSV(testFile)
	if err != nil {
		t.Fatalf("LoadCSV failed: %v", err)
	}

	// Verify data integrity
	if len(loaded) != len(tasks) {
		t.Fatalf("Expected %d tasks, got %d", len(tasks), len(loaded))
	}

	for i, task := range loaded {
		if task.ID != tasks[i].ID {
			t.Errorf("Task %d: ID mismatch, expected %d, got %d", i, tasks[i].ID, task.ID)
		}
		if task.Description != tasks[i].Description {
			t.Errorf("Task %d: Description mismatch, expected '%s', got '%s'", i, tasks[i].Description, task.Description)
		}
		if task.Done != tasks[i].Done {
			t.Errorf("Task %d: Done mismatch, expected %t, got %t", i, tasks[i].Done, task.Done)
		}
	}
}

func TestCSVLoadWithInvalidData(t *testing.T) {
	testFile := "invalid_test.csv"
	defer os.Remove(testFile)

	// Создаем CSV с разными типами невалидных данных
	invalidCSV := `ID,Description,Done
					1,Valid task,false
					invalid_id,Another task,true
					3,Task with invalid bool,invalid_bool
					5,Valid task 2,true
					`
	os.WriteFile(testFile, []byte(invalidCSV), 0644)

	// LoadCSV должен пропускать невалидные строки и загружать только валидные
	loaded, err := LoadCSV(testFile)
	if err != nil {
		t.Fatalf("LoadCSV should handle invalid data gracefully: %v", err)
	}

	// Должны загрузиться только валидные задачи (ID: 1 и 5)
	if len(loaded) != 2 {
		t.Errorf("Expected 2 valid tasks, got %d. Tasks: %+v", len(loaded), loaded)
		return
	}

	// Проверяем первую валидную задачу
	if loaded[0].ID != 1 {
		t.Errorf("Expected task with ID 1, got %d", loaded[0].ID)
	}
	if loaded[0].Description != "Valid task" {
		t.Errorf("Expected description 'Valid task', got '%s'", loaded[0].Description)
	}
	if loaded[0].Done {
		t.Error("Task 1 should not be done")
	}

	// Проверяем вторую валидную задачу
	if loaded[1].ID != 5 {
		t.Errorf("Expected task with ID 5, got %d", loaded[1].ID)
	}
	if loaded[1].Description != "Valid task 2" {
		t.Errorf("Expected description 'Valid task 2', got '%s'", loaded[1].Description)
	}
	if !loaded[1].Done {
		t.Error("Task 5 should be done")
	}
}

func TestCSVWithSpecialCharacters(t *testing.T) {
	testFile := "special_chars_test.csv"
	defer os.Remove(testFile)

	tasks := []todo.Task{
		{ID: 1, Description: "Task, with, commas", Done: false},
		{ID: 2, Description: "Task with \"quotes\"", Done: true},
		{ID: 3, Description: "Task with 'apostrophes'", Done: false},
		{ID: 4, Description: "Task with\nnewline", Done: true},
	}

	err := SaveCSV(testFile, tasks)
	if err != nil {
		t.Fatalf("SaveCSV failed: %v", err)
	}

	loaded, err := LoadCSV(testFile)
	if err != nil {
		t.Fatalf("LoadCSV failed: %v", err)
	}

	if len(loaded) != len(tasks) {
		t.Fatalf("Expected %d tasks, got %d", len(tasks), len(loaded))
	}

	// Проверяем сохранение специальных символов
	if loaded[0].Description != "Task, with, commas" {
		t.Errorf("Commas not preserved: expected 'Task, with, commas', got '%s'", loaded[0].Description)
	}
	if loaded[1].Description != "Task with \"quotes\"" {
		t.Errorf("Quotes not preserved: expected 'Task with \"quotes\"', got '%s'", loaded[1].Description)
	}
}
