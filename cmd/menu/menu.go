package menu

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ng1n3/go-todo/internal/config"
	"github.com/Ng1n3/go-todo/internal/errors"
	"github.com/Ng1n3/go-todo/internal/service"
	"github.com/Ng1n3/go-todo/internal/ui"
)

type MenuController struct {
	input       *ui.InputReader
	display     *ui.Display
	config      *config.Config
	todoService *service.TodoService
}

func NewMenuController() *MenuController {
	return &MenuController{
		input:   ui.NewInputReader(),
		display: ui.NewDisplay(),
		config:  config.Default(),
	}
}

func (mc *MenuController) Start() {
	fmt.Println("\nWelcome to Go TODO app")

	if err := mc.config.EnsureStorageDir(); err != nil {
		mc.display.ShowError(fmt.Errorf("failed to create storage directory: %w", err))
		return
	}

	for {
		choice, err := mc.input.ReadChoice("\n1.) Create a new Todo file\n2.) Load from my todo files\n3.) List todo files\n4.) Delete todo files\n5.) Exit app\nChoice: ",
			[]string{"1", "2", "3", "4", "5"})

		if err != nil {
			mc.display.ShowError(err)
			continue
		}

		switch choice {
		case "1":
			mc.createTodoFile()
		case "2":
			mc.loadTodoFile()
		case "3":
			mc.listTodoFiles()
		case "4":
			mc.deleteTodoFile()
		case "5":
			fmt.Println("Bye. Hope to see you soon!")
			return
		}
	}

}

func (mc *MenuController) createTodoFile() {
	filename, err := mc.input.ReadString("Give your new todo file a name: ")
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	normalizedName, err := mc.normalizeFileName(filename)
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	fullPath := mc.config.GetFullPath(normalizedName)

	// check if file already exists
	if _, err := os.Stat(fullPath); err == nil {
		mc.display.ShowError(fmt.Errorf("file %s already exists", normalizedName))
		return
	} else if !os.IsNotExist(err) {
		mc.display.ShowError(fmt.Errorf("error checking file %s: %w", normalizedName, err))
		return
	}

	todoService, err := service.NewTodoService(fullPath, mc.config)
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	mc.todoService = todoService
	mc.display.ShowSuccess(fmt.Sprintf("Created todo file :%s", normalizedName))
	mc.todoMenu()
}

func (mc *MenuController) loadTodoFile() {
	filename, err := mc.input.ReadString("Enter the name of the todo file you wan to load: ")
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	if !strings.HasSuffix(filename, ".json") {
		filename = ".json"
	}

	fullPath := mc.config.GetFullPath(filename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		mc.display.ShowError(errors.ErrFileNotFound)
		return
	}

	todoService, err := service.NewTodoService(fullPath, mc.config)
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	mc.todoService = todoService
	mc.display.ShowSuccess(fmt.Sprintf("Loaded todo file: %s", filename))
	mc.todoMenu()
}

func (mc *MenuController) listTodoFiles() {
	files, err := os.ReadDir(mc.config.StorageDir)
	if err != nil {
		mc.display.ShowError(fmt.Errorf("failed to read storage directory: %w", err))
		return
	}

	var fileInfos []os.FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, info)
	}
	mc.display.ShowFiles(fileInfos, mc.config.StorageDir)
}

func (mc *MenuController) deleteTodoFile() {
	mc.listTodoFiles()

	filename, err := mc.input.ReadString("Enter the name of the file you want to delete: ")
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	if !strings.HasSuffix(filename, ".json") {
		filename += ".json"
	}

	fullPath := mc.config.GetFullPath(filename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		mc.display.ShowError(errors.ErrFileNotFound)
		return
	}

	confirm, err := mc.input.ReadChoice(fmt.Sprintf("Are you sure you want to delete '%s'? (y/n): ", filename),
		[]string{"y", "n", "yes", "no"})
	if err != nil {
		mc.display.ShowError(err)
		return
	}

	if confirm == "y" || confirm == "yes" {
		if err := os.Remove(fullPath); err != nil {
			mc.display.ShowError(fmt.Errorf("failed to delete file %s: %w", filename, err))
			return
		}
		mc.display.ShowSuccess("Todo file deleted Successfully")
	} else {
		mc.display.ShowInfo("Deletion cancelled")
	}
}

func (mc *MenuController) normalizeFileName(input string) (string, error) {
	name := strings.TrimSpace(input)
	if name == "" {
		return "", errors.ErrInvalidInput
	}

	name = strings.TrimSuffix(name, ".json")

	name = filepath.Base(name)

	if name == "" || name == "." {
		return "", errors.ErrInvalidInput
	}

	return name + ".json", nil
}
