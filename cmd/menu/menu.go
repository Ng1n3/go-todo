package menu

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ng1n3/go-todo/internal/store"
)

var ts *store.TodoStorage

func MainMenu() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("1.) Create a new Todo file\n2.) Load from my todo files\n3.) List todo files\n4.) Exit app\n")
		choice, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\nThere was an error reading a name for your file: %v\n", err)
			return
		}
		choice = strings.TrimSpace(choice)
		switch choice {
		case "1":
			CreateTodoFile()
		case "4":
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
	ts = store.NewTodoStorage(storageName)

	Menu()
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

func CreateTodo() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Give your Todo a title \n")
	title, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nthere was an error reading the string: %v\n", err)
	}

	fmt.Printf("What is your todo about? \n")
	body, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nthere was an error reading the string: %v\n", err)
	}
	todo, err := store.Create(title, body)
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
	fmt.Printf("List of Todos: %v", todos)
}
