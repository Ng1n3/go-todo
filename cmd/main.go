package main

import (
	"github.com/Ng1n3/go-todo/cmd/menu"
)

func main() {
	controller := menu.NewMenuController()
	controller.Start()
}
