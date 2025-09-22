package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/Ng1n3/go-todo/internal/types"
	"github.com/olekukonko/tablewriter"
)

type Display struct {}

func NewDisplay() *Display {
  return &Display{}
}

func (d *Display) ShowTodos(todos []types.Todo) {
  if len(todos) == 0 {
    fmt.Println("No todos found.")
    return
  }

  table := tablewriter.NewWriter(os.Stdout)
  table.Header([]string{"ID", "Task", "Due Date", "Priority", "Completed", "Labels", "Created", "Updated"})

  for _, todo := range todos {
    labels := strings.Join(todo.Labels, ", ")
    completed := "No"
    if todo.Completed {
      completed = "Yes"
    }
    
    table.Append([]string {
      todo.ID,
      todo.Task,
      todo.DueDate.Format("2006-01-02"),
      string(todo.Priority),
      completed,
      labels,
      todo.CreatedAt.Format("2006-01-02"),
      todo.UpdatedAt.Format("2006-01-02"),
    })
  }

  table.Render()
}

func (d *Display) ShowFiles(files []os.FileInfo, storageDir string) {
  if len(files) == 0 {
    fmt.Println("No todos files found.")
    return
  }
  
  table := tablewriter.NewWriter(os.Stdout)
  table.Header([]string{"File Name", "Size (KB)", "Modified"})

  for _, file := range files {
    if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
      table.Append([]string {
        file.Name(),
        fmt.Sprintf("%0.2f", float64(file.Size())/1024.0),
        file.ModTime().Format("2006-01-02 15:04"),
      })
    }
  }
  table.Render()
}

func (d *Display) ShowError(err error) {
  fmt.Printf("Error: %v\n", err)
}

func (d *Display) ShowSuccess(message string) {
  fmt.Printf("✓ %s\n", message)
}

func (d *Display) ShowInfo (message string) {
  fmt.Printf("i %s\n",message)
}
