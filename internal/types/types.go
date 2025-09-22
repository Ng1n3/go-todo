package types

import (
	"fmt"
	"time"
)

type Priority string

const (
	High   Priority = "HIGH"
	Medium Priority = "MEDIUM"
	Low    Priority = "LOW"
)

type Todo struct {
	ID        string    `json:"id"`
	Task      string    `json:"task"`
	Labels    []string  `json:"labels"`
	Completed bool      `json:"completed"`
	DueDate   time.Time `json:"due_date"`
	Priority  Priority  `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p Priority) Validate() error {
	switch p {
	case High, Medium, Low:
		return nil
	default:
		return fmt.Errorf("\ninvalid Priority: %s. Must be HIGHT, MEDIUM, OR LOW", p)
	}
}

func (t *Todo) Validate() error {
	if len(t.Task) < 2 {
		return fmt.Errorf("\ntask must be at least 2 characters long\n")
	}

	if t.ID == "" {
		return fmt.Errorf("\ntodo ID cannot be empty\n")
	}

	if err := t.Priority.Validate(); err != nil {
		return err
	}

	return nil
}
