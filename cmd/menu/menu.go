package menu

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ng1n3/go-todo/internal/store"
	"github.com/Ng1n3/go-todo/internal/types"
	"github.com/olekukonko/tablewriter"
)

var ts *store.TodoStorage

func MainMenu() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("1.) Create a new Todo file\n2.) Load from my todo files\n3.) List todo files\n4.) Delete todo files\n5.) Exit app \n")
		choice, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\nThere was an error reading a name for your file: %v\n", err)
			return
		}
		choice = strings.TrimSpace(choice)
		switch choice {
		case "1":
			CreateTodoFile()
		case "2":
			LoadTodo()
		case "3":
			ListTodoFiles()
		case "5":
			fmt.Println("Bye. Hope to see you soon!")
			return
		default:
			fmt.Printf("\nSorry this command is invalid\n")
		}
	}
}

func CreateTodoFile() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Give your new todo file a name:")
	rawName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nThere was an error reading a name for your file: %v\n", err)
		return
	}

	storageName, err := normalizeJSONFilename(rawName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if store.FileExists("save_todos.json", storageName) {
		fmt.Printf("Please find another name for your new todo, as %s has already being created by you.\n", storageName)
		return
	}

	storageDir := "storage"
	fullPath := filepath.Join(storageDir, storageName)

	ts = store.NewTodoStorage(fullPath)
	Menu()
}

func ListTodoFiles() {
	storageDir := "storage"

	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		fmt.Printf("\nno storage directory found. Create some todos first: %v\n", err)
		return
	}

	files, err := os.ReadDir(storageDir)
	if err != nil {
		fmt.Printf("error reading storage directory: %v\n", err)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"File Name", "Size (KB)", "Created At", "Updated At"})

	for _, file := range files {
		if !file.IsDir() {
			filepath := filepath.Join(storageDir, file.Name())
			info, err := os.Stat(filepath)
			if err != nil {
				fmt.Printf("\nerror getting file info for %s: %v\n", file.Name(), err)
				continue
			}
			table.Append([]string{
				file.Name(),
				fmt.Sprintf("%.2f", float64(info.Size())/1024.0),
				info.ModTime().Format("2006-01-02 15.04"),
				info.ModTime().Format("2006-01-02 15:04"),
			})
		}
	}
	table.Render()
}

func Menu() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("What would you like to do today\n1.) Create a todo\n2.) List my todos\n3.) Update a todo\n4.) Delete a todo\n5.) Back to Main menu\n")

		choice, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\nThere was an error reading the choices: %v\n", err)
			return
		}
		choice = strings.TrimSpace(choice)
		switch choice {
		case "1":
			CreateTodo()
		case "2":
			ListTodo()
		case "5":
			return
		default:
			fmt.Printf("\nSorry this command is invalid\n")
		}
	}
}

func normalizeJSONFilename(input string) (string, error) {
	name := strings.TrimSpace(input)
	name = filepath.Base(name) // avoid paths like foo/bar.json

	for strings.HasSuffix(strings.ToLower(name), ".json") {
		name = name[:len(name)-len(".json")]
	}
	if name == "" {
		return "", errors.New("please input a name for your todo file")
	}

	return name + ".json", nil
}

func toPriority(p string) types.Priority {
	p = strings.ToLower(p)
	switch p {
	case "low":
		return types.Low
	case "medium":
		return types.Medium
	case "high":
		return types.High
	default:
		return types.Low
	}
}

func CreateTodo() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("what tasks do you want to perform? \n")
	task, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nthere was an error reading your tasks: %v\n", err)
	}

	fmt.Printf("How important is this task, choose one between  HIGH, MEDIUM, LOW, choose one? \n")
	priority, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nthere was an error reading your priority: %v\n", err)
	}

	priority = strings.TrimSpace(priority)
	priority = strings.ToUpper(priority)

	fmt.Printf("When is this task due, use format (2022-06-12)\n")
	dueDate, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nthere was an error reading your due date: %v\n", err)
	}

	fmt.Printf("What labels would you give this todo \n")
	labelsInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nthere was an error reading your labels: %v\n", err)
	}

	labelsInput = strings.TrimSpace(labelsInput)
	labels := strings.Split(labelsInput, ",")

	todo, err := store.Create(task, dueDate, toPriority(priority), labels)
	if err != nil {
		fmt.Println(err)
	}

	ts.Save(todo)
	ts.Persist()
	ts.SaveSummary("save_todos.json")
	fmt.Printf("Todo was just created : %v\n", todo)
}

func ListTodo() {
	todos := ts.List()

	if len(todos) == 0 {
		fmt.Println("\nNo todos found.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "Task", "Due Date", "Priority", "Completed", "Labels", "Created At", "Updated At"})

	for _, todo := range todos {
		labels := strings.Join(todo.Labels, ", ")
		completed := "No"
		if todo.Completed {
			completed = "Yes"
		}

		table.Append([]string{
			todo.ID,
			todo.Task,
			todo.DueDate.Format("2006-01-02"),
			string(todo.Priority),
			completed,
			labels,
			todo.CreatedAt.Format("2006-01-02 15:04"),
			todo.UpdatedAt.Format("2006-01-02 15:04"),
		})
	}

	table.Render()
	// fmt.Printf("List of Todos: %v", todos)
}

func LoadTodo() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nEnter the name of the todo file you want to load \n")
	filename, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nthere was an error reading the name the filename: %v\n", err)
	}

	filename = strings.TrimSpace(filename)

	if _, err := os.Stat(filename); err != nil {
		fmt.Printf("\nthere was an error searching for the file you need: %v\n", err)
		return
	}

	ts = store.NewTodoStorage(filename)
	fmt.Printf("\nTodo file '%s' loaded successfully!\n", filename)

	Menu()
}
