package menu

import (
	"fmt"
	"strings"

	"github.com/Ng1n3/go-todo/internal/types"
)

func (mc *MenuController) todoMenu() {
	for {
		choice, err := mc.input.ReadChoice("\n1.) Create Todo\n2.)List Todos \n3.)Update Todo\n4.)Delete Todo\n5.)main menu \nChoice: ", []string{"1", "2", "3", "4", "5"})
		if err != nil {
			mc.display.ShowError(err)
			continue
		}

		switch choice {
		case "1":
			mc.createTodo()
		case "2":
			mc.listTodo()
		case "3":
			mc.updateTodo()
		case "4":
			mc.deleteTodo()
		case "5":
			mc.display.ShowInfo("Returning to Main menu ...")
			return
		default:
			mc.display.ShowError(fmt.Errorf("invalid input"))
		}
	}
}

func (mc *MenuController) createTodo() {
	task, err := mc.input.ReadString("Enter task: ")
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	dueDate, err := mc.input.ReadString("Enter due date (YYYY-MD-DD format): ")
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	priority, err := mc.input.ReadChoice("Enter priority (low/medium/high, default low): ", []string{"low", "medium", "high", ""})
	priority = strings.TrimSpace(priority)
	if err != nil {
		mc.display.ShowError(err)
		return
	}
	if priority == "" {
		priority = "low"
	}

	labels, err := mc.input.ReadString("Enter labels (comma-separated, optional): ")
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	completed, err := mc.input.ReadChoice("Is the task completed? (y/n, default n): ", []string{"y", "n", "yes", "no", ""})
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	if completed == "" || completed == "n" || completed == "no" {
		completed = "false"
	} else {
		completed = "true"
	}

	todo, err := mc.todoService.CreateTodo(task, dueDate, completed, types.Priority(strings.ToLower(priority)), labels)
	if err != nil {
		mc.display.ShowError(fmt.Errorf("failed to save todo: %w", err))
		return
	}

	if err := mc.todoService.Save(); err != nil {
		mc.display.ShowError(fmt.Errorf("failed to save todo: %w", err))
		return
	}

	mc.display.ShowSuccess(fmt.Sprintf("Todo created successfully: %s", todo.Task))

}

func (mc *MenuController) listTodo() {
	todos := mc.todoService.ListTodos()
	mc.display.ShowTodos(todos)
}

func (mc *MenuController) updateTodo() {
	todos := mc.todoService.ListTodos()

	mc.display.ShowTodos(todos)

	todoID, err := mc.input.ReadString("Enter the id of the todo to update: ")
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	todo, err := mc.todoService.GetTodo(todoID)
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	mc.display.ShowTodo(todo)

	field, err := mc.input.ReadChoice("\nwhich field woud you like to update?\n1.) Task\n2.)Due Date \n3.)Priority\n4.)Labels\n5.)Completed Status \n6.) back \nChoice: ", []string{"1", "2", "3", "4", "5", "6"})
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	updates := make(map[string]any)

	switch field {
	case "1":
		newTask, err := mc.input.ReadString("📝 Enter new task title: ")
		if err != nil {
			mc.display.ShowError(err)
			return
		}
		updates["task"] = newTask
	case "2":
		newDate, err := mc.input.ReadString("📅 Enter new due date (YYYY-MM-DD): ")
		if err != nil {
			mc.display.ShowError(err)
			return
		}
		updates["due_date"] = newDate
	case "3":
		newPriority, err := mc.input.ReadPriority("⭐ Enter new priority (HIGH/MEDIUM/LOW): ")
		if err != nil {
			mc.display.ShowError(err)
			return
		}
		updates["priority"] = newPriority
	case "4":
		newLabels := mc.input.ReadLabels("🏷️  Enter new labels (comma-separated): ")
		updates["labels"] = newLabels
	case "5":
		completed, err := mc.input.ReadBool("✅ Is the task completed? (true/false): ")
		if err != nil {
			mc.display.ShowError(err)
			return
		}
		updates["completed"] = completed
	case "6":
		mc.display.ShowInfo("Returning to menu...")
		return
	}

	if err := mc.todoService.UpdateTodo(todoID, updates); err != nil {
		mc.display.ShowError(err)
		return
	}

	if err := mc.todoService.Save(); err != nil {
		mc.display.ShowError(fmt.Errorf("failed to save updates: %w", err))
		return
	}

	mc.display.ShowSuccess("Todo updated successfully!")
}

func (mc *MenuController) deleteTodo() {
	todos := mc.todoService.ListTodos()

	mc.display.ShowTodos(todos)

	todoID, err := mc.input.ReadString("Enter the id of the todo to delete: ")
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	todo, err := mc.todoService.GetTodo(todoID)
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	mc.display.ShowTodo(todo)
	err = mc.todoService.DeleteTodo(todoID)
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	mc.display.ShowSuccess("Todo deleted successfully!")

}
