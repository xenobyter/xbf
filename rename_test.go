package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRenameEntry(t *testing.T) {
	t.Run("rename file successfully", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "old.txt")
		dst := filepath.Join(dir, "new.txt")

		if err := os.WriteFile(src, []byte("content"), 0o644); err != nil {
			t.Fatal(err)
		}

		if err := renameEntry(src, dst); err != nil {
			t.Fatalf("renameEntry failed: %v", err)
		}

		data, err := os.ReadFile(dst)
		if err != nil || string(data) != "content" {
			t.Fatalf("expected renamed file content, got %q (err: %v)", string(data), err)
		}

		if _, err := os.Stat(src); !os.IsNotExist(err) {
			t.Fatalf("expected source to not exist, got err: %v", err)
		}
	})

	t.Run("rename creates destination parent directories", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "old.txt")
		dst := filepath.Join(dir, "nested", "deep", "new.txt")

		if err := os.WriteFile(src, []byte("content"), 0o644); err != nil {
			t.Fatal(err)
		}

		if err := renameEntry(src, dst); err != nil {
			t.Fatalf("renameEntry failed: %v", err)
		}

		if _, err := os.Stat(dst); err != nil {
			t.Fatalf("expected destination file to exist, got err: %v", err)
		}
	})

	t.Run("rename directory successfully", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "folder_old")
		dst := filepath.Join(dir, "folder_new")
		nested := filepath.Join(src, "nested.txt")

		if err := os.MkdirAll(src, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(nested, []byte("nested"), 0o644); err != nil {
			t.Fatal(err)
		}

		if err := renameEntry(src, dst); err != nil {
			t.Fatalf("renameEntry failed: %v", err)
		}

		if _, err := os.Stat(filepath.Join(dst, "nested.txt")); err != nil {
			t.Fatalf("expected renamed directory content, got err: %v", err)
		}
		if _, err := os.Stat(src); !os.IsNotExist(err) {
			t.Fatalf("expected source directory to be removed, got err: %v", err)
		}
	})

	t.Run("rename symlink successfully", func(t *testing.T) {
		dir := t.TempDir()
		target := filepath.Join(dir, "target.txt")
		if err := os.WriteFile(target, []byte("target"), 0o644); err != nil {
			t.Fatal(err)
		}

		src := filepath.Join(dir, "link_old")
		dst := filepath.Join(dir, "link_new")
		if err := os.Symlink(target, src); err != nil {
			t.Skip("symlinks not supported")
		}

		if err := renameEntry(src, dst); err != nil {
			t.Fatalf("renameEntry failed for symlink: %v", err)
		}

		linkTarget, err := os.Readlink(dst)
		if err != nil {
			t.Fatalf("failed reading renamed symlink: %v", err)
		}
		if linkTarget != target {
			t.Fatalf("expected symlink target %s, got %s", target, linkTarget)
		}
	})

	t.Run("source does not exist", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "missing.txt")
		dst := filepath.Join(dir, "new.txt")

		if err := renameEntry(src, dst); err == nil {
			t.Fatal("expected error for missing source, got nil")
		}
	})

	t.Run("same source and destination errors", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "same.txt")
		if err := os.WriteFile(src, []byte("same"), 0o644); err != nil {
			t.Fatal(err)
		}

		if err := renameEntry(src, src); err == nil {
			t.Fatal("expected error when source and destination are identical, got nil")
		}
	})

	t.Run("cannot overwrite directory with file", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "file.txt")
		dst := filepath.Join(dir, "existing_dir")

		if err := os.WriteFile(src, []byte("data"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(dst, 0o755); err != nil {
			t.Fatal(err)
		}

		if err := renameEntry(src, dst); err == nil {
			t.Fatal("expected error when overwriting directory with file, got nil")
		}
	})

	t.Run("cannot overwrite file with directory", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "dir_src")
		dst := filepath.Join(dir, "file_dst")

		if err := os.MkdirAll(src, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dst, []byte("data"), 0o644); err != nil {
			t.Fatal(err)
		}

		if err := renameEntry(src, dst); err == nil {
			t.Fatal("expected error when overwriting file with directory, got nil")
		}
	})

	t.Run("cannot rename directory into itself", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "parent")
		dst := filepath.Join(src, "child")

		if err := os.MkdirAll(src, 0o755); err != nil {
			t.Fatal(err)
		}

		if err := renameEntry(src, dst); err == nil {
			t.Fatal("expected error when renaming directory into itself, got nil")
		}
	})
}
