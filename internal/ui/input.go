// Package ui provides the user input interface for the go-todo application.
//
// It defines utilities for reading and validating user input from the command line.
// The InputReader type wraps a bufio.Reader and exposes helper methods to read
// strings, choices, priorities, labels, and boolean values. These methods are
// designed to validate input and guide the user when invalid values are entered.
package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Ng1n3/go-todo/internal/types"
	"github.com/Ng1n3/go-todo/internal/utils"
)

type InputReader struct {
  reader *bufio.Reader
}

func NewInputReader () *InputReader {
  return &InputReader{
    reader: bufio.NewReader(os.Stdin),
  }
}

func (ir *InputReader) ReadString (prompt string) (string, error) {
  fmt.Print(prompt)
  input, err := ir.reader.ReadString('\n')
  if err != nil {
    return "", fmt.Errorf("failed to read input: %w", err)
  }
  return strings.TrimSpace(input), nil
}


func (ir *InputReader) ReadChoice(prompt string, validChoices []string) (string, error) {
  for {
    choice, err := ir.ReadString(prompt)
    if err != nil {
      return "", err
    }

    for _, valid := range validChoices {
      if choice == valid {
        return choice, nil
      }
    }

    fmt.Printf("Invalid choice. Valid options: %s\n", strings.Join(validChoices, ", "))
  }
}


func (ir *InputReader) ReadPriority(prompt string) (types.Priority, error) {
  priorityStr, err := ir.ReadString(prompt)
  if err != nil {
    return "", err
  }

  priorityStr = strings.ToUpper(strings.TrimSpace(priorityStr))
  
  switch priorityStr {
  case "LOW", "L":
    return types.Low, nil

  case "HIGH", "H":
    return types.High, nil

  case "MEDIUM", "M":
    return types.Medium, nil

  case "":
    return types.Low, nil

  default:
    return types.Low, fmt.Errorf("invalid priority: %s. Use HIGH, MEDIUM, ", priorityStr)
  }
}

func (ir *InputReader) ReadLabels (prompt string) []string {
  labelsInput, err := ir.ReadString(prompt)
  if err != nil {
    return []string{}
  }
  return utils.ValidateLabels(labelsInput)
}

func (ir *InputReader) ReadBool (prompt string) (bool, error) {
  input, err := ir.ReadString(prompt)
  if err != nil {
    return false, err
  }

  input = strings.ToLower(input)
  switch input {
  case "true", "t", "yes", "y", "1":
    return true, nil
  case "false", "f", "no", "n", "2", "0":
    return false, nil
  default:
    parsed, err := strconv.ParseBool(input)
    if err != nil {
      return false, fmt.Errorf("invalid boolean value: %s", input)
    }
    return parsed, nil
  }
}
