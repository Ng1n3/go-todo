package store

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
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

func create(title, body string) (types.Todo, error) {
	if len(title) < 3 {
		return types.Todo{}, fmt.Errorf("length of title must be above 4")
	} else if len(body) < 2 {
		return types.Todo{}, fmt.Errorf("length of body must be above 4")
	}
	return types.Todo{
		ID:        generateID(6),
		Title:     title,
		Body:      body,
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
