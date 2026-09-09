package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSortFiles(t *testing.T) {
	input := []FileInfo{
		{Name: "zebra.txt", IsDir: false},
		{Name: "Beta_dir", IsDir: true},
		{Name: "apple.txt", IsDir: false},
		{Name: "alpha_dir", IsDir: true},
		{Name: "B.txt", IsDir: false},
		{Name: "a.txt", IsDir: false},
	}

	expected := []FileInfo{
		{Name: "alpha_dir", IsDir: true},
		{Name: "Beta_dir", IsDir: true},
		{Name: "a.txt", IsDir: false},
		{Name: "apple.txt", IsDir: false},
		{Name: "B.txt", IsDir: false},
		{Name: "zebra.txt", IsDir: false},
	}

	sortFiles(input)

	for i := range expected {
		if !reflect.DeepEqual(input[i], expected[i]) {
			t.Errorf("at index %d: expected %+v, got %+v", i, expected[i], input[i])
		}
	}
}

func TestCopyFiles(t *testing.T) {
	t.Run("empty items slice", func(t *testing.T) {
		tempDir := t.TempDir()
		if err := copyFiles(nil, tempDir); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
		if err := copyFiles([]FileInfo{}, tempDir); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("invalid target directory", func(t *testing.T) {
		tempDir := t.TempDir()
		nonExistentTarget := filepath.Join(tempDir, "does_not_exist")
		items := []FileInfo{{Name: "file.txt", Path: tempDir}}

		if err := copyFiles(items, nonExistentTarget); err == nil {
			t.Fatal("expected error when target does not exist, got nil")
		}

		filePath := filepath.Join(tempDir, "regular_file.txt")
		if err := os.WriteFile(filePath, []byte("hello"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := copyFiles(items, filePath); err == nil {
			t.Fatal("expected error when target is a file, got nil")
		}
	})

	t.Run("copy single and multiple files", func(t *testing.T) {
		srcDir := t.TempDir()
		dstDir := t.TempDir()

		file1 := filepath.Join(srcDir, "file1.txt")
		file2 := filepath.Join(srcDir, "file2.txt")
		if err := os.WriteFile(file1, []byte("content 1"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file2, []byte("content 2"), 0o644); err != nil {
			t.Fatal(err)
		}

		items := []FileInfo{
			{Name: "file1.txt", Path: srcDir},
			{Name: "file2.txt", Path: srcDir},
		}

		if err := copyFiles(items, dstDir); err != nil {
			t.Fatalf("copyFiles failed: %v", err)
		}

		dst1 := filepath.Join(dstDir, "file1.txt")
		dst2 := filepath.Join(dstDir, "file2.txt")

		data1, err := os.ReadFile(dst1)
		if err != nil || string(data1) != "content 1" {
			t.Fatalf("expected 'content 1', got %s (err: %v)", string(data1), err)
		}
		data2, err := os.ReadFile(dst2)
		if err != nil || string(data2) != "content 2" {
			t.Fatalf("expected 'content 2', got %s (err: %v)", string(data2), err)
		}
	})

	t.Run("copy directory recursively", func(t *testing.T) {
		srcDir := t.TempDir()
		dstDir := t.TempDir()

		subDir := filepath.Join(srcDir, "subfolder")
		nestedDir := filepath.Join(subDir, "nested")
		if err := os.MkdirAll(nestedDir, 0o755); err != nil {
			t.Fatal(err)
		}

		fileInSub := filepath.Join(subDir, "sub.txt")
		fileInNested := filepath.Join(nestedDir, "nested.txt")
		if err := os.WriteFile(fileInSub, []byte("sub content"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fileInNested, []byte("nested content"), 0o644); err != nil {
			t.Fatal(err)
		}

		items := []FileInfo{
			{Name: "subfolder", Path: srcDir, IsDir: true},
		}

		if err := copyFiles(items, dstDir); err != nil {
			t.Fatalf("copyFiles failed for directory: %v", err)
		}

		copiedSub := filepath.Join(dstDir, "subfolder", "sub.txt")
		copiedNested := filepath.Join(dstDir, "subfolder", "nested", "nested.txt")

		data1, err := os.ReadFile(copiedSub)
		if err != nil || string(data1) != "sub content" {
			t.Fatalf("expected 'sub content', got %s (err: %v)", string(data1), err)
		}
		data2, err := os.ReadFile(copiedNested)
		if err != nil || string(data2) != "nested content" {
			t.Fatalf("expected 'nested content', got %s (err: %v)", string(data2), err)
		}
	})

	t.Run("overwrite existing destination file", func(t *testing.T) {
		srcDir := t.TempDir()
		dstDir := t.TempDir()

		srcFile := filepath.Join(srcDir, "test.txt")
		dstFile := filepath.Join(dstDir, "test.txt")

		if err := os.WriteFile(srcFile, []byte("new version"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dstFile, []byte("old version"), 0o644); err != nil {
			t.Fatal(err)
		}

		items := []FileInfo{{Name: "test.txt", Path: srcDir}}
		if err := copyFiles(items, dstDir); err != nil {
			t.Fatalf("copyFiles failed: %v", err)
		}

		data, err := os.ReadFile(dstFile)
		if err != nil || string(data) != "new version" {
			t.Fatalf("expected 'new version', got %s (err: %v)", string(data), err)
		}
	})

	t.Run("copy to same file errors", func(t *testing.T) {
		dir := t.TempDir()
		file := filepath.Join(dir, "self.txt")
		if err := os.WriteFile(file, []byte("same"), 0o644); err != nil {
			t.Fatal(err)
		}

		items := []FileInfo{{Name: "self.txt", Path: dir}}
		if err := copyFiles(items, dir); err == nil {
			t.Fatal("expected error when copying file to itself, got nil")
		}
	})

	t.Run("copy directory into itself errors", func(t *testing.T) {
		dir := t.TempDir()
		parent := filepath.Join(dir, "parent")
		child := filepath.Join(parent, "child")
		if err := os.MkdirAll(child, 0o755); err != nil {
			t.Fatal(err)
		}

		items := []FileInfo{{Name: "parent", Path: dir, IsDir: true}}
		if err := copyFiles(items, parent); err == nil {
			t.Fatal("expected error when copying directory into itself, got nil")
		}
	})

	t.Run("copy symlink", func(t *testing.T) {
		srcDir := t.TempDir()
		dstDir := t.TempDir()

		targetFile := filepath.Join(srcDir, "target.txt")
		if err := os.WriteFile(targetFile, []byte("symlink target"), 0o644); err != nil {
			t.Fatal(err)
		}

		symlinkPath := filepath.Join(srcDir, "link.txt")
		if err := os.Symlink(targetFile, symlinkPath); err != nil {
			t.Skip("symlinks not supported on this environment")
		}

		items := []FileInfo{{Name: "link.txt", Path: srcDir}}
		if err := copyFiles(items, dstDir); err != nil {
			t.Fatalf("copyFiles failed for symlink: %v", err)
		}

		dstLink := filepath.Join(dstDir, "link.txt")
		linkTarget, err := os.Readlink(dstLink)
		if err != nil {
			t.Fatalf("failed reading copied symlink: %v", err)
		}
		if linkTarget != targetFile {
			t.Fatalf("expected symlink target %s, got %s", targetFile, linkTarget)
		}
	})
}

