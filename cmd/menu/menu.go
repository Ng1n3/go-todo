package menu

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
		case "4":
			DeleteTodoFile()
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

func DeleteTodoFile() {
	storageDir := "storage/"

	reader := bufio.NewReader(os.Stdin)
	ListTodoFiles()

	fmt.Printf("\nEnter the name of the name of the file you want to delete\n")
	filename, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nthere was an error reading your filename: %v\n", err)
		return
	}

	filename = storageDir + strings.TrimSpace(filename)
	err = os.Remove(filename)
	if err != nil {
		fmt.Printf("\nThere was an error deleting the file: %v\n", err)
		return
	}

	fmt.Printf("\nFile '%v'  was successfully deleted!\n", filename)
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
		case "3":
			UpdateTodo()
		case "4":
			DeleteTodo()
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

func UpdateTodo() {

	todos := ts.List()
	reader := bufio.NewReader(os.Stdin)

	// show the list of todos with their IDs
	ListTodo()

	// As the user for the ID of the todo to update
	fmt.Printf("\n Enter the Todo ID you would like to update from.\n")
	userId, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nthere was an error reading userId for updating: %v\n", err)
		return
	}

	userId = strings.TrimSpace(userId)

	// Find the todo with the given id
	var selectedTodo *types.Todo

	for _, todo := range todos {
		if todo.ID == userId {
			selectedTodo = &todo
		}
	}

	if selectedTodo == nil {
		fmt.Printf("Todo with id '%s' not found", userId)
		return
	}

	fmt.Printf("\nWhat field would you like to update \n1.) Title \n2.) Due Date \n3.) Priority \n4.) Labels \n5.) Completed Status \n")
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nthere was an error reading options for updating: %v\n", err)
		return
	}
	choice = strings.TrimSpace(choice)
	switch choice {
	case "1":
		fmt.Printf("\nEnter a new Task title:  \n")
		newTitle, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\nthere was an error reading your title for the update: %v\n", err)
			return
		}
		selectedTodo.Task = newTitle
	case "2":
		fmt.Printf("\nEnter a new due Date: \n")
		dueDateStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\nthere was an error reading your due date for the update: %v\n", err)
			return
		}
		dueDateStr = strings.TrimSpace(dueDateStr)

		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil {
			fmt.Printf("\ninvalid date format. Please use YYYY-MM-DD\n")
			return
		}

		selectedTodo.DueDate = dueDate

	case "3":
		fmt.Printf("\nEnter the new priority: \n")
		priority, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\nthere was an error reading your priority for the update: %v\n", err)
			return
		}
		priority = strings.TrimSpace(priority)
		priority = strings.ToUpper(priority)

		selectedTodo.Priority = toPriority(priority)
	case "4":
		fmt.Printf("\nEnter new lables: \n")
		labelsInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\nthere was an error reading your priority for the update: %v\n", err)
			return
		}

		labelsInput = strings.TrimSpace(labelsInput)
		labels := strings.Split(labelsInput, ",")
		selectedTodo.Labels = labels
	case "5":
		fmt.Printf("\nEnter a completed status: \n")
		completedStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\nthere was an error reading your priority for the update: %v\n", err)
			return
		}
		completedStr = strings.TrimSpace(completedStr)

		completed, err := strconv.ParseBool(completedStr)
		if err != nil {
			fmt.Printf("\nthere was an error converting string to boolean: %v\n", err)
			return
		}
		selectedTodo.Completed = completed
	default:
		fmt.Printf("\nSorry this command is invalid\n")
		return
	}

	selectedTodo.UpdatedAt = time.Now()

	ts.Save(*selectedTodo)
	ts.Persist()
	ts.SaveSummary("save_todos.json")

	fmt.Println("Todo updated successfully")
}

func DeleteTodo() {
	todos := ts.List()

	ListTodo()

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\n Enter the id of the todo you want to delete: \n")
	todoID, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nthere was an error reading the todo id: %v\n", err)
		return
	}

	todoID = strings.TrimSpace(todoID)

	var selectedTodo *types.Todo

	for _, todo := range todos {
		if todo.ID == todoID {
			selectedTodo = &todo
		}
	}

	if selectedTodo == nil {
		fmt.Printf("\n selected todo with id '%v' not found! \n", todoID)
		return
	}

	fmt.Printf("\nAre you sure you want to delete this todo with task: %v\n1.) Yes \n2.) No\n", selectedTodo.Task)
	deletedChoice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nthere was an error reading your delete Choice: %v\n", err)
	}

	deletedChoice = strings.TrimSpace(deletedChoice)
	switch deletedChoice {
	case "1":
		if ts.Delete(selectedTodo.ID) {
			fmt.Printf("\nTodo successfully Deleted!\n")
			ts.Persist()
			ts.SaveSummary("save_todos.json")
		} else {
			fmt.Printf("\nFailed to delete todo!\n")
		}
	case "2":
		fmt.Printf("\nDeletion cancelled.\n")
	default:
		fmt.Printf("\nInvalid choice.\n")
	}

}

func LoadTodo() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nEnter the name of the todo file you want to load \n")
	filename, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nthere was an error reading the name the filename: %v\n", err)
	}

	filename = "storage/" + strings.TrimSpace(filename)

	if _, err := os.Stat(filename); err != nil {
		fmt.Printf("\nthere was an error searching for the file you need: %v\n", err)
		return
	}

	ts = store.NewTodoStorage(filename)
	fmt.Printf("\nTodo file '%s' loaded successfully!\n", filename)

	Menu()
}
