package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"todo-app/internal/storage"
	"todo-app/internal/todo"

	"github.com/ZeRg0912/logger"
)

// handleAdd processes the add command to create a new task.
// It expects a --desc flag with the task description.
// Returns the updated task slice.
func handleAdd(tasks []todo.Task, args []string) ([]todo.Task, error) {
	logger.Debug("handleAdd called with %d args", len(args))

	addCmd := flag.NewFlagSet("add", flag.ContinueOnError)
	desc := addCmd.String("desc", "", "Task description")
	setupCommandConfig(addCmd)

	err := addCmd.Parse(args)
	if err != nil {
		printCommandUsage("add", addCmd, "add a new task")
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if *desc == "" {
		printCommandUsage("add", addCmd, "add a new task")
		return nil, fmt.Errorf("task description cannot be empty: use --desc flag")
	}

	newTasks := todo.Add(tasks, *desc)
	logger.ConsoleSuccess("Task added: %s", *desc)
	return newTasks, nil
}

// handleList processes the list command to display tasks.
// Supports --filter flag with values: all, done, pending.
// Tasks are displayed with status emojis and IDs.
func handleList(tasks []todo.Task, args []string) error {
	logger.Debug("handleList called with %d args", len(args))

	listCmd := flag.NewFlagSet("list", flag.ContinueOnError)
	filter := listCmd.String("filter", "all", "Task filter: all, done, pending")
	setupCommandConfig(listCmd)

	err := listCmd.Parse(args)
	if err != nil {
		printCommandUsage("list", listCmd, "list tasks")
		return fmt.Errorf("invalid arguments: %w", err)
	}

	validFilters := map[string]bool{"all": true, "done": true, "pending": true}
	if !validFilters[*filter] {
		printCommandUsage("list", listCmd, "list tasks")
		return fmt.Errorf("invalid filter value '%s'", *filter)
	}

	filteredTasks := todo.List(tasks, *filter)
	if len(filteredTasks) == 0 {
		logger.Info("No tasks found with filter '%s'", *filter)
		logger.ConsoleHelp("No tasks found")
		return nil
	}

	logger.Info("Displaying %d tasks with filter '%s'", len(filteredTasks), *filter)
	logger.ConsoleHelpf("Task list (%s):", *filter)
	for _, task := range filteredTasks {
		status := "[ ]"
		if task.Done {
			status = "[X]"
		}
		logger.ConsoleHelpf("%s [ID:%d] %s", status, task.ID, task.Description)
	}
	return nil
}

// handleComplete processes the complete command to mark a task as done.
// It expects a --id flag with the task ID to complete.
// Returns the updated task slice.
func handleComplete(tasks []todo.Task, args []string) ([]todo.Task, error) {
	logger.Debug("handleComplete called with %d args", len(args))

	completeCmd := flag.NewFlagSet("complete", flag.ContinueOnError)
	id := completeCmd.Int("id", 0, "Task ID to mark as completed")
	setupCommandConfig(completeCmd)

	err := completeCmd.Parse(args)
	if err != nil {
		printCommandUsage("complete", completeCmd, "mark task as completed")
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if *id == 0 {
		printCommandUsage("complete", completeCmd, "mark task as completed")
		return nil, fmt.Errorf("task ID is required and must be greater than 0")
	}

	resultTasks, err := todo.Complete(tasks, *id)
	if err != nil {
		return nil, fmt.Errorf("cannot complete task %d: %w", *id, err)
	}

	logger.ConsoleSuccess("Task %d marked as completed", *id)
	return resultTasks, nil
}

// handleDelete processes the delete command to remove a task.
// It expects a --id flag with the task ID to delete.
// Returns the updated task slice.
func handleDelete(tasks []todo.Task, args []string) ([]todo.Task, error) {
	logger.Debug("handleDelete called with %d args", len(args))

	deleteCmd := flag.NewFlagSet("delete", flag.ContinueOnError)
	id := deleteCmd.Int("id", 0, "Task ID to delete")
	setupCommandConfig(deleteCmd)

	err := deleteCmd.Parse(args)
	if err != nil {
		printCommandUsage("delete", deleteCmd, "delete a task")
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if *id == 0 {
		printCommandUsage("delete", deleteCmd, "delete a task")
		return nil, fmt.Errorf("task ID is required and must be greater than 0")
	}

	resultTasks, err := todo.Delete(tasks, *id)
	if err != nil {
		return nil, fmt.Errorf("cannot delete task %d: %w", *id, err)
	}

	logger.ConsoleSuccess("Task %d deleted", *id)
	return resultTasks, nil
}

// handleExport processes the export command to save tasks to a file.
// Supports --format flag (json or csv) and --out flag for output file.
// Automatically adds file extension if not specified.
func handleExport(tasks []todo.Task, args []string) error {
	logger.Debug("handleExport called with %d args", len(args))

	exportCmd := flag.NewFlagSet("export", flag.ContinueOnError)
	format := exportCmd.String("format", "json", "Export format: json or csv")
	outFile := exportCmd.String("out", "tasks_export", "Output file")
	setupCommandConfig(exportCmd)

	err := exportCmd.Parse(args)
	if err != nil {
		printCommandUsage("export", exportCmd, "export tasks to file")
		return fmt.Errorf("invalid arguments: %w", err)
	}

	validFormats := map[string]bool{"json": true, "csv": true}
	if !validFormats[*format] {
		printCommandUsage("export", exportCmd, "export tasks to file")
		return fmt.Errorf("invalid format '%s'", *format)
	}

	if !strings.HasSuffix(*outFile, "."+*format) {
		*outFile = *outFile + "." + *format
	}

	switch *format {
	case "json":
		err = storage.SaveJSON(*outFile, tasks)
	case "csv":
		err = storage.SaveCSV(*outFile, tasks)
	}

	if err != nil {
		return fmt.Errorf("export error: %w", err)
	}

	logger.Info("Tasks exported to %s", *outFile)
	logger.ConsoleHelpf("Tasks exported to %s", *outFile)
	return nil
}

// handleLoad processes the load command to import tasks from a file.
// It expects a --file flag with the path to import from.
// Supports JSON and CSV formats based on file extension.
// Returns the imported tasks slice and error if any.
func handleLoad(args []string) ([]todo.Task, error) {
	logger.Debug("handleLoad called with %d args", len(args))

	loadCmd := flag.NewFlagSet("load", flag.ContinueOnError)
	file := loadCmd.String("file", "", "File to import from")
	setupCommandConfig(loadCmd)

	if len(args) == 0 {
		return nil, fmt.Errorf("load command requires --file flag: specify file to import")
	}

	err := loadCmd.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if *file == "" {
		return nil, fmt.Errorf("import file is required")
	}

	if _, err := os.Stat(*file); os.IsNotExist(err) {
		if _, err := os.Stat(*file + ".csv"); err == nil {
			*file = *file + ".csv"
		} else if _, err := os.Stat(*file + ".json"); err == nil {
			*file = *file + ".json"
		} else {
			return nil, fmt.Errorf("file does not exist: %s", *file)
		}
	}

	// Determine format by file extension
	ext := strings.ToLower(filepath.Ext(*file))
	var importedTasks []todo.Task

	logger.Info("Starting import from file: %s (format: %s)", *file, ext)

	switch ext {
	case ".json":
		importedTasks, err = storage.LoadJSON(*file)
	case ".csv":
		importedTasks, err = storage.LoadCSV(*file)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("import error: %w", err)
	}

	logger.Info("Successfully imported %d tasks from %s", len(importedTasks), *file)
	logger.ConsoleHelpf("Successfully imported %d tasks from %s", len(importedTasks), *file)
	return importedTasks, nil
}

// printCommandUsage displays formatted help for a specific command.
// It shows command syntax, available flags, and usage examples.
func printCommandUsage(cmd string, flags *flag.FlagSet, description string) {
	var flagLines []string
	flags.VisitAll(func(f *flag.Flag) {
		flagLines = append(flagLines, fmt.Sprintf("  --%-12s %s", f.Name, f.Usage))
	})

	exampleFlag := "--id=1"
	if cmd == "add" {
		exampleFlag = "--desc=\"Your task description\""
	} else if cmd == "list" {
		exampleFlag = "--filter=pending"
	} else if cmd == "export" {
		exampleFlag = "--format=csv|json --out=backup"
	} else if cmd == "load" {
		exampleFlag = "--file=tasks.csv | tasks.json"
	}

	message := fmt.Sprintf(
		"Usage: <app> %s [flags]\nDescription: %s\nFlags:\n%s\nExample: todo %s %s",
		cmd,
		description,
		strings.Join(flagLines, "\n"),
		cmd,
		exampleFlag,
	)

	logger.ConsoleHelp(message)
}

// printUsage displays the main help message with all available commands.
// It provides an overview of the application and usage examples.
func printUsage() {
	fmt.Println("To-Do Manager - command line task management")
	fmt.Println("Usage: <app_name> <command> [arguments]")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("-  add --desc=\"description\"          - add a new task")
	fmt.Println("-  list [--filter=all|done|pending]    - list tasks")
	fmt.Println("-  complete --id=ID                    - mark task as completed")
	fmt.Println("-  delete --id=ID                      - delete a task")
	fmt.Println("-  export --format=json|csv --out=file - export tasks")
	fmt.Println("-  load --file=file                    - import tasks from file")
	fmt.Println("-  help                                - show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  <app_name> add --desc=\"Buy milk\"")
	fmt.Println("  <app_name> list --filter=pending")
	fmt.Println("  <app_name> complete --id=3")
	fmt.Println("  <app_name> delete --id=3")
	fmt.Println("  <app_name> export --format=csv --out=backup")
	fmt.Println("  <app_name> load --file=tasks.csv")
	fmt.Println("  <app_name> help")
}

// setupCommandConfig configures command flags to suppress default output.
// It disables automatic help printing and error output from the flag package.
func setupCommandConfig(cmd *flag.FlagSet) {
	cmd.SetOutput(io.Discard)
	cmd.Usage = func() {}
}
