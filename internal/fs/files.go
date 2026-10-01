package fs

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"
)

// FileInfo holds basic metadata about a file or directory entry.
type FileInfo struct {
	Name     string
	Path     string
	IsDir    bool
	Size     int64
	IsHidden bool
}

// GetWd returns the current working directory.
func GetWd() (string, error) {
	return getWd()
}

// FilterFiles returns files whose names contain query (case-insensitive).
func FilterFiles(files []FileInfo, query string) []FileInfo {
	return filterFiles(files, query)
}

// FilterHidden removes hidden entries from the slice.
func FilterHidden(files []FileInfo) []FileInfo {
	return filterHidden(files)
}

// ReadDir returns information about all entries in the specified directory.
func ReadDir(path string) ([]FileInfo, error) {
	return readDir(path)
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

func filterFiles(files []FileInfo, query string) []FileInfo {
	query = strings.ToLower(query)
	filtered := make([]FileInfo, 0, len(files))
	for _, file := range files {
		if strings.Contains(strings.ToLower(file.Name), query) {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

// filterHidden removes hidden entries from the slice.
func filterHidden(files []FileInfo) []FileInfo {
	filtered := make([]FileInfo, 0, len(files))
	for _, file := range files {
		if !file.IsHidden && !strings.HasPrefix(file.Name, ".") {
			filtered = append(filtered, file)
		}
	}
	return filtered
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
			Name:     entry.Name(),
			Path:     path,
			IsDir:    entry.IsDir(),
			Size:     info.Size(),
			IsHidden: isHidden(entry.Name(), path, info),
		})
	}

	sortFiles(files)

	return files, nil
}

// FormatFileSize formats a file size in bytes into a human-readable string.
func FormatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}

	units := []string{"KB", "MB", "GB", "TB"}
	value := float64(size)
	for _, unit := range units {
		value /= 1024
		if value < 1024 || unit == units[len(units)-1] {
			return fmt.Sprintf("%.1f %s", value, unit)
		}
	}

	return fmt.Sprintf("%d B", size)
}