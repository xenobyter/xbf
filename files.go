package main

import (
	"cmp"
	"os"
	"slices"
	"strings"
)

// FileInfo holds basic metadata about a file or directory entry.
type FileInfo struct {
	Name  string // Name of the file or directory
	IsDir bool   // Whether the entry is a directory
	Size  int64  // File size in bytes
}

// getWd returns the current working directory.
func getWd() (string, error) {
	return os.Getwd()
}

// sortFiles sorts files with directories first, then by name
// using a case-insensitive comparison.
func sortFiles(files []FileInfo) {
	slices.SortFunc(files, func(a, b FileInfo) int {
		if a.IsDir != b.IsDir {
			if a.IsDir {
				return -1
			}
			return 1
		}

		if n := cmp.Compare(
			strings.ToLower(a.Name),
			strings.ToLower(b.Name),
		); n != 0 {
			return n
		}

		return cmp.Compare(a.Name, b.Name)
	})
}

// readDir returns information about all entries in the specified directory.
func readDir(path string) ([]FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files := make([]FileInfo, 0, len(entries))

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		files = append(files, FileInfo{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Size:  info.Size(),
		})
	}

	sortFiles(files)

	return files, nil
}
