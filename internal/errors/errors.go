package errors

import "errors"

var (
	ErrTodoNotFound          = errors.New("todo not found")
	ErrFileNotFound          = errors.New("file not found")
	ErrInvalidInput          = errors.New("invalid input")
	ErrFileExists            = errors.New("file already exists")
	ErrInvalidDateFormat     = errors.New("invalid date format")
	ErrTaskTooShort          = errors.New("task must be at least 2 characters long")
	ErrInvalidCompletedValue = errors.New("invalid input")
)
