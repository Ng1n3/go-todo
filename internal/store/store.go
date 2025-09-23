// Package store provides functionality for managing and persisting todos.
//
// It defines the TodoStorage type, which handles CRUD operations (Create,
// Read, Update, Delete) on todos, as well as saving them to disk and
// generating summary files. The package abstracts away the low-level
// details of reading from and writing to JSON files, allowing higher-level
// components (like the CLI menu) to interact with todos through a simple API.
//
// Key features:
//   - Load and persist todos to a JSON file
//   - Save individual todos after validation
//   - Delete todos by ID with error handling
//   - List all stored todos
//   - Save a summary file containing just the todo tasks
//   - Retrieve todos by ID
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Ng1n3/go-todo/internal/config"
	"github.com/Ng1n3/go-todo/internal/errors"
	"github.com/Ng1n3/go-todo/internal/types"
)

type TodoStorage struct {
	store  map[string]types.Todo
	file   string
	config *config.Config
}

func NewTodoStorage(file string, cfg *config.Config) (*TodoStorage, error) {
	if cfg == nil {
		cfg = config.Default()
	}

	ts := &TodoStorage{store: make(map[string]types.Todo), file: file, config: cfg}
	if err := ts.Load(); err != nil {
		return nil, fmt.Errorf("failed to load todos: %w", err)
	}
	return ts, nil
}

func (ts *TodoStorage) Load() error {
	data, err := os.ReadFile(ts.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, &ts.store); err != nil {
		return fmt.Errorf("failed to unmarshal todos : %w", err)
	}

	return nil

}

func (ts *TodoStorage) Persist() error {
	data, err := json.MarshalIndent(ts.store, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal todos: %w", err)
	}

	if err := os.WriteFile(ts.file, data, ts.config.FileMode); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func (ts *TodoStorage) Save(todo *types.Todo) error {
	if err := todo.Validate(); err != nil {
		return fmt.Errorf("invalid todo: %w", err)
	}

	todo.UpdatedAt = time.Now()
	ts.store[todo.ID] = *todo
	return nil
}

func (ts *TodoStorage) SaveSummary(summaryFile string) error {
	tasks := make([]string, 0, len(ts.store)) // title arrays
	for _, todo := range ts.store {
		tasks = append(tasks, strings.TrimSpace(todo.Task))
	}

	summary := map[string][]string{
		ts.file: tasks,
	}

	data, err := json.MarshalIndent(summary, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marashal summary: %w", err)
	}

	if err := os.WriteFile(summaryFile, data, ts.config.FileMode); err != nil {
		return fmt.Errorf("failed to write summary file: %w", err)
	}

	return nil

}

func (ts *TodoStorage) Get(id string) (types.Todo, error) {
	todo, exists := ts.store[id]
	if !exists {
		return types.Todo{}, errors.ErrTodoNotFound
	}

	return todo, nil
}

func (ts *TodoStorage) Delete(id string) error {
	if _, exists := ts.store[id]; !exists {
		return errors.ErrTodoNotFound
	}

	delete(ts.store, id)
	return nil
}

func (ts *TodoStorage) List() []types.Todo {
	todos := make([]types.Todo, 0, len(ts.store))
	for _, todo := range ts.store {
		todos = append(todos, todo)
	}
	return todos
}

func (ts *TodoStorage) Count() int {
	return len(ts.store)
}
