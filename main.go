package main

import (
	"fmt"

	"github.com/Ng1n3/go-todo/cmd/menu"
)

func main() {
  controller := menu.NewMenuController()
  controller.Start()
  if confirm == "y" || confirm == "yes" {
    if err := os.Remove(fullPath)
  }
	fmt.Printf("\nWelcome to GO TODO app\n")
	menu.MainMenu()
}
