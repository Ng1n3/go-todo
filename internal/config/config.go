// Package config provides configuration management for storage directories, summary file paths, and default file permission
package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	StorageDir  string
	SummaryFile string
	FileMode    os.FileMode
}

func Default() *Config {
	return &Config{
		StorageDir:  "storage",
		SummaryFile: "save_todos.json",
		FileMode:    0644,
	}
}

// EnsureStorageDir create storage Directory if it doesn't exist
func (c *Config) EnsureStorageDir() error {
	if _, err := os.Stat(c.StorageDir); os.IsNotExist(err) {
		return os.MkdirAll(c.StorageDir, 0755)
	}
	return nil
}

// GetFullPath return full path for filename
func (c *Config) GetFullPath(filename string) string {
	return filepath.Join(c.StorageDir, filename)
}
