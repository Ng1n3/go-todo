package utils

import (
	"strings"
	"time"

	"github.com/Ng1n3/go-todo/internal/errors"
)

func ValidateTask(task string) (string, error) {
	task = strings.TrimSpace(task)
	if len(task) < 2 {
		return "", errors.ErrTaskTooShort
	}

	return task, nil
}

func ValidateDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Time{}, errors.ErrInvalidDateFormat
	}

	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, errors.ErrInvalidDateFormat
	}

	return parsedDate, nil
}

func ValidateLabels(labelsInput string) []string {
	if labelsInput == "" {
		return []string{}
	}

	labelsInput = strings.TrimSpace(labelsInput)
	labels := strings.Split(labelsInput, ",")

	var cleanLabel []string
	for _, label := range labels {
		label = strings.TrimSpace(label)
		if label != "" {
			cleanLabel = append(cleanLabel, label)
		}
	}

	return cleanLabel
}

func ValidateCompleted(completedInput string) (bool, error) {
	completedInput = strings.TrimSpace(strings.ToLower(completedInput))

	switch completedInput {
	case "yes", "y", "true", "1":
		return true, nil
	case "no", "n", "false", "0":
		return false, nil
	case "":
		return false, nil
	default:
		return false, errors.ErrInvalidCompletedValue
	}
}
