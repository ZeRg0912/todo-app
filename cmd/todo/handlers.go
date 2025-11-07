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
func handleAdd(tasks []todo.Task, args []string) []todo.Task {
	logger.Debug("handleAdd called with %d args", len(args))

	addCmd := flag.NewFlagSet("add", flag.ContinueOnError)
	desc := addCmd.String("desc", "", "Task description")
	setupCommandConfig(addCmd)

	if len(args) == 0 {
		logger.Error("add command requires --desc flag")
		printCommandUsage("add", addCmd, "add a new task")
		os.Exit(1)
	}

	err := addCmd.Parse(args)
	if err != nil {
		logger.Error("Invalid arguments: %v", err)
		printCommandUsage("add", addCmd, "add a new task")
		os.Exit(1)
	}

	if *desc == "" {
		logger.Error("task description cannot be empty")
		printCommandUsage("add", addCmd, "add a new task")
		os.Exit(1)
	}

	newTasks := todo.Add(tasks, *desc)
	logger.ConsoleSuccess("Task added: %s", *desc)
	return newTasks
}

// handleList processes the list command to display tasks.
// Supports --filter flag with values: all, done, pending.
// Tasks are displayed with status emojis and IDs.
func handleList(tasks []todo.Task, args []string) {
	logger.Debug("handleList called with %d args", len(args))

	listCmd := flag.NewFlagSet("list", flag.ContinueOnError)
	filter := listCmd.String("filter", "all", "Task filter: all, done, pending")
	setupCommandConfig(listCmd)

	err := listCmd.Parse(args)
	if err != nil {
		logger.Error("Invalid arguments: %v", err)
		printCommandUsage("list", listCmd, "list tasks")
		os.Exit(1)
	}

	validFilters := map[string]bool{"all": true, "done": true, "pending": true}
	if !validFilters[*filter] {
		logger.Error("invalid filter value '%s'", *filter)
		printCommandUsage("list", listCmd, "list tasks")
		os.Exit(1)
	}

	filteredTasks := todo.List(tasks, *filter)
	if len(filteredTasks) == 0 {
		logger.Info("No tasks found with filter '%s'", *filter)
		logger.ConsoleHelp("No tasks found")
		return
	}

	logger.Info("Displaying %d tasks with filter '%s'", len(filteredTasks), *filter)
	logger.ConsoleHelpf("Task list (%s):", *filter)
	for _, task := range filteredTasks {
		status := "❌"
		if task.Done {
			status = "✅"
		}
		logger.ConsoleHelpf("%s [ID:%d] %s", status, task.ID, task.Description)
	}
}

// handleComplete processes the complete command to mark a task as done.
// It expects a --id flag with the task ID to complete.
// Returns the updated task slice.
func handleComplete(tasks []todo.Task, args []string) []todo.Task {
	logger.Debug("handleComplete called with %d args", len(args))

	completeCmd := flag.NewFlagSet("complete", flag.ContinueOnError)
	id := completeCmd.Int("id", 0, "Task ID to mark as completed")
	setupCommandConfig(completeCmd)

	if len(args) == 0 {
		logger.Error("complete command requires --id flag")
		printCommandUsage("complete", completeCmd, "mark task as completed")
		os.Exit(1)
	}

	err := completeCmd.Parse(args)
	if err != nil {
		logger.Error("Invalid arguments: %v", err)
		printCommandUsage("complete", completeCmd, "mark task as completed")
		os.Exit(1)
	}

	if *id == 0 {
		logger.Error("task ID is required and must be greater than 0")
		printCommandUsage("complete", completeCmd, "mark task as completed")
		os.Exit(1)
	}

	resultTasks, err := todo.Complete(tasks, *id)
	if err != nil {
		logger.Error("cannot complete task %d: %v", *id, err)
		os.Exit(1)
	}

	logger.ConsoleSuccess("Task %d marked as completed", *id)
	return resultTasks
}

// handleDelete processes the delete command to remove a task.
// It expects a --id flag with the task ID to delete.
// Returns the updated task slice.
func handleDelete(tasks []todo.Task, args []string) []todo.Task {
	logger.Debug("handleDelete called with %d args", len(args))

	deleteCmd := flag.NewFlagSet("delete", flag.ContinueOnError)
	id := deleteCmd.Int("id", 0, "Task ID to delete")
	setupCommandConfig(deleteCmd)

	if len(args) == 0 {
		logger.Error("delete command requires --id flag")
		printCommandUsage("delete", deleteCmd, "delete a task")
		os.Exit(1)
	}

	err := deleteCmd.Parse(args)
	if err != nil {
		logger.Error("Invalid arguments: %v", err)
		printCommandUsage("delete", deleteCmd, "delete a task")
		os.Exit(1)
	}

	if *id == 0 {
		logger.Error("task ID is required and must be greater than 0")
		printCommandUsage("delete", deleteCmd, "delete a task")
		os.Exit(1)
	}

	resultTasks, err := todo.Delete(tasks, *id)
	if err != nil {
		logger.Error("cannot delete task %d: %v", *id, err)
		os.Exit(1)
	}

	logger.ConsoleSuccess("Task %d deleted", *id)
	return resultTasks
}

// handleExport processes the export command to save tasks to a file.
// Supports --format flag (json or csv) and --out flag for output file.
// Automatically adds file extension if not specified.
func handleExport(tasks []todo.Task, args []string) {
	logger.Debug("handleExport called with %d args", len(args))

	exportCmd := flag.NewFlagSet("export", flag.ContinueOnError)
	format := exportCmd.String("format", "json", "Export format: json or csv")
	outFile := exportCmd.String("out", "tasks_export", "Output file")
	setupCommandConfig(exportCmd)

	if len(args) == 0 {
		logger.Error("export command requires --format and --out flags")
		printCommandUsage("export", exportCmd, "export tasks to file")
		os.Exit(1)
	}

	err := exportCmd.Parse(args)
	if err != nil {
		logger.Error("Invalid arguments: %v", err)
		printCommandUsage("export", exportCmd, "export tasks to file")
		os.Exit(1)
	}

	validFormats := map[string]bool{"json": true, "csv": true}
	if !validFormats[*format] {
		logger.Error("invalid format '%s'", *format)
		printCommandUsage("export", exportCmd, "export tasks to file")
		os.Exit(1)
	}

	// Add file extension if not specified
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
		logger.Error("export error: %v", err)
		os.Exit(1)
	}

	logger.Info("Tasks exported to %s", *outFile)
	logger.ConsoleHelpf("Tasks exported to %s", *outFile)
}

// handleLoad processes the load command to import tasks from a file.
// It expects a --file flag with the path to import from.
// Supports JSON and CSV formats based on file extension.
// Returns the imported tasks slice.
func handleLoad(args []string) []todo.Task {
	logger.Debug("handleLoad called with %d args", len(args))

	loadCmd := flag.NewFlagSet("load", flag.ContinueOnError)
	file := loadCmd.String("file", "", "File to import from")
	setupCommandConfig(loadCmd)

	if len(args) == 0 {
		logger.Error("load command requires --file flag")
		printCommandUsage("load", loadCmd, "import tasks from file")
		os.Exit(1)
	}

	err := loadCmd.Parse(args)
	if err != nil {
		logger.Error("Invalid arguments: %v", err)
		printCommandUsage("load", loadCmd, "import tasks from file")
		os.Exit(1)
	}

	if *file == "" {
		logger.Error("import file is required")
		printCommandUsage("load", loadCmd, "import tasks from file")
		os.Exit(1)
	}

	if _, err := os.Stat(*file); os.IsNotExist(err) {
		logger.Error("file does not exist: %s", *file)
		os.Exit(1)
	}

	// Determine format by file extension
	ext := strings.ToLower(filepath.Ext(*file))
	var importedTasks []todo.Task

	switch ext {
	case ".json":
		importedTasks, err = storage.LoadJSON(*file)
	case ".csv":
		importedTasks, err = storage.LoadCSV(*file)
	default:
		logger.Error("unsupported file format: %s", ext)
		os.Exit(1)
	}

	if err != nil {
		logger.Error("import error: %v", err)
		os.Exit(1)
	}

	logger.Info("Successfully imported %d tasks from %s", len(importedTasks), *file)
	logger.ConsoleHelpf("Successfully imported %d tasks from %s", len(importedTasks), *file)
	return importedTasks
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
