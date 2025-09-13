package store

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Ng1n3/go-todo/internal/types"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWZYZ0123456789"

type TodoStorage struct {
	store map[string]types.Todo
	file  string
}

func NewTodoStorage(file string) *TodoStorage {
	ts := &TodoStorage{store: make(map[string]types.Todo), file: file}
	ts.Load()
	return ts
}

func (ts *TodoStorage) Load() {
	data, err := os.ReadFile(ts.file)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &ts.store)
	if err != nil {
		fmt.Printf("error loading todos: %v", err)
	}
}

func (ts *TodoStorage) Persist() {
	data, err := json.MarshalIndent(ts.store, "", " ")
	if err != nil {
		fmt.Printf("error saving file: %v", err)
		return
	}
	err = os.WriteFile(ts.file, data, 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v", err)
	}
}

func (ts *TodoStorage) Save(todo types.Todo) {
	ts.store[todo.ID] = todo
}

func (ts *TodoStorage) SaveSummary(file string) {
	tasks := make([]string, 0, len(ts.store)) // title arrays
	for _, todo := range ts.store {
		tasks = append(tasks, strings.TrimSpace(todo.Task))
	}

	summary := map[string][]string{
		ts.file: tasks,
	}

	data, err := json.MarshalIndent(summary, "", " ")
	if err != nil {
		fmt.Printf("\nan error occured while marshalling: %v\n", err)
	}

	err = os.WriteFile(file, data, 0644)
	if err != nil {
		fmt.Printf("\nthere was an error writting to your json file: %v\n", err)
	}

}

func (ts *TodoStorage) Get(id string) (types.Todo, bool) {
	todo, ok := ts.store[id]
	return todo, ok
}

func (ts *TodoStorage) Delete(id string) bool {
	_, ok := ts.store[id]
	if !ok {
		return false
	}

	delete(ts.store, id)
	return true
}

func (ts *TodoStorage) List() []types.Todo {
	todos := make([]types.Todo, 0, len(ts.store))
	for _, todo := range ts.store {
		todos = append(todos, todo)
	}
	return todos
}

func FileExists(summaryFile, filename string) bool {
	if _, err := os.Stat(summaryFile); os.IsNotExist(err) {
		emptySummary := make(map[string][]string)

		data, _ := json.MarshalIndent(emptySummary, "", " ")

		if err := os.WriteFile(summaryFile, data, 0644); err != nil {
			fmt.Printf("\nfailed to create summary file: %v\n", err)
			return false
		}

		fmt.Printf("\nCreated a new summary file: %s\n", summaryFile)
	}

	data, err := os.ReadFile(summaryFile)
	if err != nil {
		fmt.Printf("\nthere was an error reading your summary file: %v\n", err)
		return false
	}

	var summary map[string][]string
	if err := json.Unmarshal(data, &summary); err != nil {
		fmt.Printf("\nthere was an error unmarshaling: %v\n", err)
		return false
	}
	_, exists := summary[filename]
	return exists

}

func Create(task, dueDate string, priority types.Priority) (types.Todo, error) {
	if len(task) < 2 {
		return types.Todo{}, fmt.Errorf("length of task must be above 2 characters")
	} else if priority == "" {
		priority = types.Low
	}

	dueDate = strings.TrimSpace(dueDate)

	parsedDate, err := time.Parse("2006-01-02", dueDate)
	if err != nil {
		return types.Todo{}, fmt.Errorf("invalid error format, it should be like (2022-12-06: %v)", err)
	}
	return types.Todo{
		ID:        generateID(6),
		Task:      task,
		DueDate:   parsedDate,
		Priority:  priority,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func generateID(n int) string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}
