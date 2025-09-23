-----

# Go TODO CLI

A robust, file-based command-line interface (CLI) for managing all your to-do lists. Built with a clean, decoupled architecture in Go, this tool allows you to separate your tasks into different files (e.g., `work.json`, `personal.json`) and manage them efficiently from your terminal.

-----

## ✨ Features

  * **🗂️ Multi-File Management**: Create, load, list, and delete separate to-do list files for different projects or contexts.
  * **📝 Full CRUD Operations**: Complete Create, Read, Update, and Delete functionality for both to-do files and the tasks within them.
  * **🏷️ Rich Task Attributes**: Each task includes a description, due date, completion status, labels, and priority (`HIGH`, `MEDIUM`, `LOW`).
  * **💅 Clean Terminal UI**: All lists are displayed in clean, formatted tables for excellent readability.
  * **💾 Persistent JSON Storage**: Your lists are saved locally in a `storage/` directory, making them easy to inspect, backup, or version control.

-----

## 🏛️ Architecture

This application is built using a clean, layered architecture to ensure separation of concerns, making the code easier to maintain, test, and extend.

```
+--------------------------------+
|      cmd/menu (Controller)     |  <-- Handles user input and directs traffic
+--------------------------------+
               |
+--------------------------------+
|    internal/service (Service)  |  <-- Contains core business logic and validation
+--------------------------------+
               |
+--------------------------------+
|     internal/store (Store)     |  <-- Manages data persistence (reading/writing files)
+--------------------------------+
```

This structure is supported by several utility packages for handling `config`, `ui`, `types`, and `errors`.

-----

## 🚀 Getting Started

### Prerequisites

  * **Go**: Version `1.18` or higher.
  * **Git**: For cloning the repository.

### Installation

1.  **Clone the repository:**

    ```sh
    git clone https://github.com/your-username/go-todo.git
    cd go-todo
    ```

2.  **Build the binary:**
    The `Makefile` provides a convenient way to build the application.

    ```sh
    make build
    ```

    This will create an executable binary in the `bin/` directory.

### Usage

Run the application from the root of the project directory:

```sh
./bin/myapp-linux
```

You'll be greeted with the **Main Menu**, where you can manage your to-do files. After creating or loading a file, you'll enter the **Todo Menu** to manage the tasks within that specific list.

-----

## 📁 How Data is Stored

All your data is stored locally in the project directory:

  * **`storage/`**: This directory contains all the to-do list files you create (e.g., `storage/work.json`, `storage/shopping.json`). Each file holds a complete list of its own tasks.
  * **`save_todos.json`**: This file at the root level acts as a summary or index, containing a simple list of tasks from all files in the `storage` directory.

-----

## 🛣️ Future Work

This project has a solid foundation. Future enhancements could include:

  * **Comprehensive Unit Tests**: Adding a full suite of tests for the `service` and `store` layers to ensure maximum reliability.
  * **Advanced TUI**: Implementing a more interactive Text User Interface (TUI) with a library like `Bubble Tea` or `tview`.
  * **Sort & Filter**: Adding options to sort tasks by due date or priority and filter them by labels or completion status.

-----
