// Package service provides the business logic layer for the go-todo application.
// 
// It acts as an intermediary between the CLI/UI layer and the storage layer.
// The TodoService type exposes high-level operations for managing todos such as
// creating, updating, deleting, and listing tasks. It validates input using the
// utils package, interacts with the store package for persistence, and applies
// application-level rules like default priorities and summary file generation.
//
// In short, service orchestrates todo management while keeping validation and
// persistence concerns separated into their respective packages.
package service

import (
	"fmt"
	"time"

	"github.com/Ng1n3/go-todo/internal/config"
	"github.com/Ng1n3/go-todo/internal/store"
	"github.com/Ng1n3/go-todo/internal/types"
	"github.com/Ng1n3/go-todo/internal/utils"
)

type TodoService struct {
  storage *store.TodoStorage
  config *config.Config
}

func NewTodoService(filename string, cfg *config.Config) (*TodoService, error) {
  if cfg == nil {
    cfg = config.Default()
  }

  storage, err := store.NewTodoStorage(filename, cfg)
  if err != nil {
    return nil, fmt.Errorf("failed to create create todo storage: %w", err)
  }

  return &TodoService{
    storage: storage,
    config: cfg,
  }, nil
}

func (ts *TodoService) CreateTodo (task, dueDate, completed  string, priority types.Priority, labels string) (*types.Todo, error) {

  validTask, err := utils.ValidateTask(task)
  if err != nil {
    return nil, err
  }

  validDate, err := utils.ValidateDate(dueDate)
  if err != nil {
    return nil, err
  }

  validCompleted, err := utils.ValidateCompleted(completed)
  if err != nil {
    return nil, err
  }

  if priority == "" {
    priority = types.Low
  }

  if err := priority.Validate(); err != nil {
    return nil, err
  }

  validLabels := utils.ValidateLabels(labels)

  todo := &types.Todo{
    ID: utils.GenerateID(6),
    Task: validTask,
    Labels: validLabels,
    Completed: validCompleted,
    DueDate: validDate,
    Priority: priority,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
  }

  if err := ts.storage.Save(todo); err != nil {
    return nil, fmt.Errorf("failed to save todo: %w",err)
  }

  return todo, nil
  
}

func (ts *TodoService) UpdateTodo (id string, updates map[string]any) error {
  todo, err := ts.storage.Get(id)
  if err != nil {
    return err
  }

  for field, value := range updates {
    switch field {
      case "task":
        if task, ok := value.(string); ok {
          validTask, err := utils.ValidateTask(task)
          if err != nil {
            return err
          }
          todo.Task = validTask
        }
      case "due_date":
        if dateStr, ok := value.(string); ok {
          validDate, err := utils.ValidateDate(dateStr)
          if err != nil {
            return err
          }
          todo.DueDate = validDate
        }
      case "priority":
        if priority, ok := value.(types.Priority); ok {
          if err := priority.Validate(); err != nil {
            return err
          }
          todo.Priority = priority
        }

      case "labels":
        if labels, ok := value.(string); ok {
          validatedLabels := utils.ValidateLabels(labels)
          todo.Labels = validatedLabels
        }
      case "completed":
        if completed, ok := value.(string); ok {
          validatedCompleted, err := utils.ValidateCompleted(completed)
          if err != nil {
            return err
          }
          todo.Completed = validatedCompleted
        }
    }
  }
  return ts.storage.Save(&todo)
}

func (ts *TodoService) DeleteTodo(id string) error {
  return ts.storage.Delete(id)
}

func (ts *TodoService) ListTodos() []types.Todo {
  return ts.storage.List()
}


func (ts *TodoService) Save() error {
  if err := ts.storage.Persist(); err != nil {
    return fmt.Errorf("failed to persist todos: %w",err)
  }

  if err := ts.storage.SaveSummary(ts.config.SummaryFile); err != nil {
    return fmt.Errorf("failed to save summary: %w",err)
  }
  return nil
}
