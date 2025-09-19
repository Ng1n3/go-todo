package types

import "time"

type Priority string

const (
	High   Priority = "HIGH"
	Medium Priority = "MEDIUM"
	Low    Priority = "LOW"
)

type Todo struct {
	ID        string
	Task      string
	Labels    []string
	Completed bool
	DueDate   time.Time
	Priority  Priority
	CreatedAt time.Time
	UpdatedAt time.Time
}
