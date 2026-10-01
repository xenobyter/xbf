package fs

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")

	if err := createEmptyFile(path); err != nil {
		t.Fatalf("createEmptyFile failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("created file does not exist: %v", err)
	}
	if info.IsDir() {
		t.Fatal("created path is a directory, want a file")
	}
	if info.Size() != 0 {
		t.Fatalf("created file size = %d, want 0", info.Size())
	}
}

func TestCreateEmptyFileRejectsInvalidPaths(t *testing.T) {
	for _, path := range []string{"", " \t\n"} {
		if err := createEmptyFile(path); !errors.Is(err, os.ErrInvalid) {
			t.Errorf("createEmptyFile(%q) error = %v, want os.ErrInvalid", path, err)
		}
	}
}

func TestCreateDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "new-directory")

	if err := createDirectory(path); err != nil {
		t.Fatalf("createDirectory failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("created directory does not exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("created path is not a directory")
	}
}

func TestCreateDirectoryRejectsInvalidPaths(t *testing.T) {
	for _, path := range []string{"", " \t\n"} {
		if err := createDirectory(path); !errors.Is(err, os.ErrInvalid) {
			t.Errorf("createDirectory(%q) error = %v, want os.ErrInvalid", path, err)
		}
	}
}

func TestValidateNewEntryName(t *testing.T) {
	tests := []struct {
		name string
		want error
	}{
		{name: "document.txt"},
		{name: " document.txt "},
		{name: "", want: errEmptyEntryName},
		{name: " \t\n", want: errEmptyEntryName},
		{name: ".", want: errInvalidEntryName},
		{name: "..", want: errInvalidEntryName},
		{name: "folder/document.txt", want: errInvalidEntryName},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := validateNewEntryName(test.name); got != test.want {
				t.Errorf("validateNewEntryName(%q) = %v, want %v", test.name, got, test.want)
			}
		})
	}
}

func TestEntryExists(t *testing.T) {
	dir := t.TempDir()
	existingPath := filepath.Join(dir, "existing")
	if err := os.WriteFile(existingPath, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	if !entryExists(existingPath) {
		t.Errorf("entryExists(%q) = false, want true", existingPath)
	}
	if entryExists(filepath.Join(dir, "missing")) {
		t.Error("entryExists returned true for a missing path")
	}
}
