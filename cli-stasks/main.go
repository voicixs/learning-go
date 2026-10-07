package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
)

// Task represents a single to-do item.
type Task struct {
	ID        int
	Title     string
	Completed bool
}

// ToDoList manages tasks in memory and syncs them with a CSV file.
type ToDoList struct {
	filename string
	tasks    []Task
}

// Load reads all rows from the CSV file into memory.
func (l *ToDoList) Load() error {
	file, err := os.Open(l.filename)
	if errors.Is(err, os.ErrNotExist) {
		// If the file doesn't exist yet, start with an empty list
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read csv: %w", err)
	}

	const (
		colID        = 0
		colCompleted = 1
		colTitle     = 2
	)

	l.tasks = []Task{}
	for _, row := range records {
		if len(row) < 3 {
			continue
		}

		id, err := strconv.Atoi(row[colID])
		if err != nil {
			continue
		}

		completed, err := strconv.ParseBool(row[colCompleted])
		if err != nil {
			continue
		}

		l.tasks = append(l.tasks, Task{
			ID:        id,
			Completed: completed,
			Title:     row[colTitle],
		})
	}
	return nil
}

// Save writes all tasks from memory into the CSV file.
func (l *ToDoList) Save() error {
	file, err := os.Create(l.filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, t := range l.tasks {
		row := []string{
			strconv.Itoa(t.ID),
			strconv.FormatBool(t.Completed),
			t.Title,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}
	return nil
}

// getNextID finds the highest current ID and adds 1 so IDs never collide.
func (l *ToDoList) getNextID() int {
	maxID := 0
	for _, t := range l.tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

// Add appends a new task and saves to disk.
func (l *ToDoList) Add(title string) error {
	id := l.getNextID()
	l.tasks = append(l.tasks, Task{
		ID:        id,
		Title:     title,
		Completed: false,
	})
	fmt.Printf("Added task #%d: %q\n", id, title)
	return l.Save()
}

// List prints all tasks to the console formatted via text/tabwriter.
func (l *ToDoList) List() {
	if len(l.tasks) == 0 {
		fmt.Println("No tasks found. Use 'add <title>' to create one.")
		return
	}

	// Parameters for tabwriter.NewWriter:
	// - output: destination writer (os.Stdout)
	// - minwidth: minimum cell width (0)
	// - tabwidth: tab width (0)
	// - padding: number of padding spaces between columns (3)
	// - padchar: character used to pad cells (' ')
	// - flags: formatting flags (0 for standard)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush() // w.Flush() calculates alignment and renders output to stdout

	// Table Header
	fmt.Fprintln(w, "ID\tSTATUS\tTITLE")
	fmt.Fprintln(w, "--\t------\t-----")

	// Table Rows
	for _, t := range l.tasks {
		status := "[ ]"
		if t.Completed {
			status = "[x]"
		}
		fmt.Fprintf(w, "%d\t%s\t%s\n", t.ID, status, t.Title)
	}
}

// Complete marks a task as done and saves to disk.
func (l *ToDoList) Complete(id int) error {
	for i := range l.tasks {
		if l.tasks[i].ID == id {
			l.tasks[i].Completed = true
			fmt.Printf("Completed task #%d: %s\n", id, l.tasks[i].Title)
			return l.Save()
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// Delete removes a task by ID and saves to disk.
func (l *ToDoList) Delete(id int) error {
	for i, t := range l.tasks {
		if t.ID == id {
			l.tasks = append(l.tasks[:i], l.tasks[i+1:]...)
			fmt.Printf("Deleted task #%d: %s\n", id, t.Title)
			return l.Save()
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func printHelp() {
	fmt.Println("Usage: todo <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  add <title>      Add a new task")
	fmt.Println("  list             List all tasks")
	fmt.Println("  complete <id>    Mark a task as completed")
	fmt.Println("  delete <id>      Delete a task")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	todoList := &ToDoList{filename: "todos.csv"}

	if err := todoList.Load(); err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	const (
		cmdIdx = 1
		valIdx = 2
	)
	command := os.Args[cmdIdx]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: Task title is required.")
			return
		}
		if err := todoList.Add(os.Args[valIdx]); err != nil {
			fmt.Printf("Error saving task: %v\n", err)
		}

	case "list":
		todoList.List()

	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("Error: Task ID is required.")
			return
		}
		id, err := strconv.Atoi(os.Args[valIdx])
		if err != nil {
			fmt.Println("Error: ID must be an integer.")
			return
		}
		if err := todoList.Complete(id); err != nil {
			fmt.Println("Error:", err)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Error: Task ID is required.")
			return
		}
		id, err := strconv.Atoi(os.Args[valIdx])
		if err != nil {
			fmt.Println("Error: ID must be an integer.")
			return
		}
		if err := todoList.Delete(id); err != nil {
			fmt.Println("Error:", err)
		}

	default:
		fmt.Printf("Unknown command: %q\n\n", command)
		printHelp()
	}
}
// ```

// ---

// ### Test It in Your Terminal

// ```bash
// # 1. Add your first task
// go run main.go add "Buy groceries"

// # 2. Add a second task
// go run main.go add "Finish Go CLI project"

// # 3. List all tasks
// go run main.go list

// # 4. Mark task 1 as completed
// go run main.go complete 1

// # 5. Verify the updated list and inspect the CSV file
// go run main.go list
// cat todos.csv
// ```
