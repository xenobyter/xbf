package main

import (
	"errors"
	"os"
	"path/filepath"
)

// deleteFiles deletes the specified items from the filesystem.
func deleteFiles(items []FileInfo) error {
	if len(items) == 0 {
		return nil
	}

	for _, item := range items {
		if item.Name == "" {
			continue
		}

		itemPath := filepath.Join(item.Path, item.Name)
		if err := deleteEntry(itemPath); err != nil {
			return err
		}
	}

	return nil
}

// deleteEntry deletes a single filesystem entry (file, directory, or symlink).
// Directories are deleted recursively.
func deleteEntry(path string) error {
	clean := filepath.Clean(path)
	if clean == "/" || clean == "." {
		return errors.New("cannot delete root or current directory: " + path)
	}

	if _, err := os.Lstat(clean); err != nil {
		return err
	}

	return os.RemoveAll(clean)
}

