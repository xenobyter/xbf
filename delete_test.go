package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeleteFiles(t *testing.T) {
	t.Run("empty items slice", func(t *testing.T) {
		if err := deleteFiles(nil); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
		if err := deleteFiles([]FileInfo{}); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("delete multiple files", func(t *testing.T) {
		dir := t.TempDir()
		file1 := filepath.Join(dir, "file1.txt")
		file2 := filepath.Join(dir, "file2.txt")

		if err := os.WriteFile(file1, []byte("content 1"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file2, []byte("content 2"), 0o644); err != nil {
			t.Fatal(err)
		}

		items := []FileInfo{
			{Name: "file1.txt", Path: dir},
			{Name: "", Path: dir}, // empty name should be skipped
			{Name: "file2.txt", Path: dir},
		}

		if err := deleteFiles(items); err != nil {
			t.Fatalf("deleteFiles failed: %v", err)
		}

		if _, err := os.Stat(file1); !os.IsNotExist(err) {
			t.Errorf("expected %s to be deleted, got err: %v", file1, err)
		}
		if _, err := os.Stat(file2); !os.IsNotExist(err) {
			t.Errorf("expected %s to be deleted, got err: %v", file2, err)
		}
	})

	t.Run("error on non-existent file", func(t *testing.T) {
		dir := t.TempDir()
		items := []FileInfo{
			{Name: "non_existent.txt", Path: dir},
		}

		if err := deleteFiles(items); err == nil {
			t.Fatal("expected error when deleting non-existent file, got nil")
		}
	})
}

func TestDeleteEntry(t *testing.T) {
	t.Run("delete regular file", func(t *testing.T) {
		dir := t.TempDir()
		file := filepath.Join(dir, "test.txt")
		if err := os.WriteFile(file, []byte("test"), 0o644); err != nil {
			t.Fatal(err)
		}

		if err := deleteEntry(file); err != nil {
			t.Fatalf("deleteEntry failed: %v", err)
		}

		if _, err := os.Stat(file); !os.IsNotExist(err) {
			t.Errorf("expected file to be deleted, got err: %v", err)
		}
	})

	t.Run("delete directory recursively", func(t *testing.T) {
		dir := t.TempDir()
		subDir := filepath.Join(dir, "nested_folder")
		fileInSubDir := filepath.Join(subDir, "inner.txt")
		subSubDir := filepath.Join(subDir, "deep")
		fileInDeep := filepath.Join(subSubDir, "deep.txt")

		if err := os.MkdirAll(subSubDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fileInSubDir, []byte("inner"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fileInDeep, []byte("deep"), 0o644); err != nil {
			t.Fatal(err)
		}

		if err := deleteEntry(subDir); err != nil {
			t.Fatalf("deleteEntry failed for directory: %v", err)
		}

		if _, err := os.Stat(subDir); !os.IsNotExist(err) {
			t.Errorf("expected directory to be deleted, got err: %v", err)
		}
	})

	t.Run("delete symlink preserves target", func(t *testing.T) {
		dir := t.TempDir()
		target := filepath.Join(dir, "target.txt")
		if err := os.WriteFile(target, []byte("important content"), 0o644); err != nil {
			t.Fatal(err)
		}

		link := filepath.Join(dir, "link_to_target")
		if err := os.Symlink(target, link); err != nil {
			t.Skip("symlinks not supported on this platform")
		}

		if err := deleteEntry(link); err != nil {
			t.Fatalf("deleteEntry failed for symlink: %v", err)
		}

		// Symlink itself should be gone
		if _, err := os.Lstat(link); !os.IsNotExist(err) {
			t.Errorf("expected symlink to be deleted, got err: %v", err)
		}

		// Target file MUST still exist!
		if data, err := os.ReadFile(target); err != nil || string(data) != "important content" {
			t.Errorf("expected target file to remain untouched, got err: %v, data: %q", err, string(data))
		}
	})

	t.Run("error on non-existent path", func(t *testing.T) {
		dir := t.TempDir()
		nonExistent := filepath.Join(dir, "does_not_exist")

		if err := deleteEntry(nonExistent); err == nil {
			t.Fatal("expected error on non-existent path, got nil")
		}
	})

	t.Run("prevent deleting root or current directory", func(t *testing.T) {
		if err := deleteEntry("/"); err == nil {
			t.Fatal("expected error when attempting to delete root, got nil")
		}
		if err := deleteEntry("."); err == nil {
			t.Fatal("expected error when attempting to delete current directory, got nil")
		}
	})
}

