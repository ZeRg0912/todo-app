package main

import (
	"fmt"
	"os"

	"todo-app/internal/storage"
	"todo-app/internal/todo"
	"todo-app/pkg/logging"
)

// main is the entry point of the To-Do Manager application.
// It initializes the logger, parses command line arguments,
// and routes to the appropriate command handler.
//
// The application supports the following commands:
//   - add: Add a new task
//   - list: List tasks with optional filtering
//   - complete: Mark a task as completed
//   - delete: Delete a task
//   - export: Export tasks to JSON or CSV
//   - load: Import tasks from JSON or CSV
//   - help: Show usage information
//
// Tasks are persisted in a JSON file and automatically saved after modifying commands.
func main() {
	// Initialize logger - LevelError to console, all levels to file
	err := logging.InitBoth(logging.LevelError, logging.LevelDebug, "logs/app.log", 10*1024*1024)
	if err != nil {
		// Before nitialize logger all info to console by fmt
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if r := recover(); r != nil {
			logging.Error("Application panic: %v", r)
			os.Exit(1)
		}
	}()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Parse args
	command := os.Args[1]
	args := os.Args[2:]

	logging.Info("Command executed: %s %v", command, args)
	logging.Debug("Full args: %#v", os.Args)

	// Load current tasks
	tasks, err := storage.LoadJSON("tasks.json")
	if err != nil {
		logging.Error("Failed to load tasks: %v", err)
		os.Exit(1)
	}

	var resultTasks []todo.Task

	// All available commands
	switch command {
	case "add":
		resultTasks = handleAdd(tasks, args)
	case "list":
		handleList(tasks, args)
	case "complete":
		resultTasks = handleComplete(tasks, args)
	case "delete":
		resultTasks = handleDelete(tasks, args)
	case "export":
		handleExport(tasks, args)
	case "load":
		resultTasks = handleLoad(args)
	case "help", "-h", "--help":
		printUsage()
	default:
		logging.Error("Unknown command: %s", command)
		printUsage()
		os.Exit(1)
	}

	// Save changes if command modified tasks
	if resultTasks != nil {
		err = storage.SaveJSON("tasks.json", resultTasks)
		if err != nil {
			logging.Error("Failed to save tasks: %v", err)
			os.Exit(1)
		}
		logging.Info("Tasks saved successfully, total tasks: %d", len(resultTasks))
	}
}
