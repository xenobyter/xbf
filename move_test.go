package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMoveFile(t *testing.T) {
	t.Run("move file successfully", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "source.txt")
		dst := filepath.Join(dir, "moved.txt")

		if err := os.WriteFile(src, []byte("move content"), 0o644); err != nil {
			t.Fatal(err)
		}

		if err := moveFile(src, dst); err != nil {
			t.Fatalf("moveFile failed: %v", err)
		}

		data, err := os.ReadFile(dst)
		if err != nil || string(data) != "move content" {
			t.Fatalf("expected 'move content', got %s (err: %v)", string(data), err)
		}

		if _, err := os.Stat(src); !os.IsNotExist(err) {
			t.Fatalf("expected source file to be removed, got err: %v", err)
		}
	})

	t.Run("move file to subfolder creating dirs", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "source.txt")
		dst := filepath.Join(dir, "sub", "deep", "moved.txt")

		if err := os.WriteFile(src, []byte("deep move"), 0o644); err != nil {
			t.Fatal(err)
		}

		if err := moveFile(src, dst); err != nil {
			t.Fatalf("moveFile failed: %v", err)
		}

		data, err := os.ReadFile(dst)
		if err != nil || string(data) != "deep move" {
			t.Fatalf("expected 'deep move', got %s (err: %v)", string(data), err)
		}

		if _, err := os.Stat(src); !os.IsNotExist(err) {
			t.Fatalf("expected source file to be removed")
		}
	})

	t.Run("overwrite existing destination", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "source.txt")
		dst := filepath.Join(dir, "existing.txt")

		if err := os.WriteFile(src, []byte("new"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dst, []byte("old"), 0o644); err != nil {
			t.Fatal(err)
		}

		if err := moveFile(src, dst); err != nil {
			t.Fatalf("moveFile failed: %v", err)
		}

		data, err := os.ReadFile(dst)
		if err != nil || string(data) != "new" {
			t.Fatalf("expected 'new', got %s (err: %v)", string(data), err)
		}

		if _, err := os.Stat(src); !os.IsNotExist(err) {
			t.Fatalf("expected source file to be removed")
		}
	})

	t.Run("source does not exist errors", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "nonexistent.txt")
		dst := filepath.Join(dir, "dest.txt")

		if err := moveFile(src, dst); err == nil {
			t.Fatal("expected error moving non-existent file, got nil")
		}
	})

	t.Run("source and destination same errors", func(t *testing.T) {
		dir := t.TempDir()
		file := filepath.Join(dir, "same.txt")
		if err := os.WriteFile(file, []byte("data"), 0o644); err != nil {
			t.Fatal(err)
		}

		if err := moveFile(file, file); err == nil {
			t.Fatal("expected error when src and dst are the same, got nil")
		}

		// ensure file was not deleted
		if _, err := os.Stat(file); err != nil {
			t.Fatalf("file should still exist: %v", err)
		}
	})

	t.Run("source is directory errors", func(t *testing.T) {
		dir := t.TempDir()
		subDir := filepath.Join(dir, "folder")
		if err := os.Mkdir(subDir, 0o755); err != nil {
			t.Fatal(err)
		}

		dst := filepath.Join(dir, "folder_dest")
		if err := moveFile(subDir, dst); err == nil {
			t.Fatal("expected error when moving directory with moveFile, got nil")
		}
	})

	t.Run("destination is directory errors", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "file.txt")
		if err := os.WriteFile(src, []byte("data"), 0o644); err != nil {
			t.Fatal(err)
		}

		subDir := filepath.Join(dir, "existing_folder")
		if err := os.Mkdir(subDir, 0o755); err != nil {
			t.Fatal(err)
		}

		if err := moveFile(src, subDir); err == nil {
			t.Fatal("expected error when destination is a directory, got nil")
		}
	})

	t.Run("move symlink", func(t *testing.T) {
		dir := t.TempDir()
		target := filepath.Join(dir, "target.txt")
		if err := os.WriteFile(target, []byte("symlink target"), 0o644); err != nil {
			t.Fatal(err)
		}

		linkSrc := filepath.Join(dir, "link.txt")
		if err := os.Symlink(target, linkSrc); err != nil {
			t.Skip("symlinks not supported")
		}

		linkDst := filepath.Join(dir, "link_moved.txt")
		if err := moveFile(linkSrc, linkDst); err != nil {
			t.Fatalf("moveFile failed for symlink: %v", err)
		}

		linkTarget, err := os.Readlink(linkDst)
		if err != nil {
			t.Fatalf("failed reading moved symlink: %v", err)
		}
		if linkTarget != target {
			t.Fatalf("expected target %s, got %s", target, linkTarget)
		}

		if _, err := os.Lstat(linkSrc); !os.IsNotExist(err) {
			t.Fatalf("expected source symlink to be removed")
		}
	})
}

func TestMoveFiles(t *testing.T) {
	t.Run("empty items slice", func(t *testing.T) {
		tempDir := t.TempDir()
		if err := moveFiles(nil, tempDir); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
		if err := moveFiles([]FileInfo{}, tempDir); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("invalid target directory", func(t *testing.T) {
		tempDir := t.TempDir()
		nonExistentTarget := filepath.Join(tempDir, "does_not_exist")
		items := []FileInfo{{Name: "file.txt", Path: tempDir}}

		if err := moveFiles(items, nonExistentTarget); err == nil {
			t.Fatal("expected error when target does not exist, got nil")
		}

		filePath := filepath.Join(tempDir, "regular_file.txt")
		if err := os.WriteFile(filePath, []byte("hello"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := moveFiles(items, filePath); err == nil {
			t.Fatal("expected error when target is a file, got nil")
		}
	})

	t.Run("move single and multiple files", func(t *testing.T) {
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

		if err := moveFiles(items, dstDir); err != nil {
			t.Fatalf("moveFiles failed: %v", err)
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

		if _, err := os.Stat(file1); !os.IsNotExist(err) {
			t.Fatalf("expected source file1 to be removed, got err: %v", err)
		}
		if _, err := os.Stat(file2); !os.IsNotExist(err) {
			t.Fatalf("expected source file2 to be removed, got err: %v", err)
		}
	})

	t.Run("move directory recursively", func(t *testing.T) {
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

		items := []FileInfo{{Name: "subfolder", Path: srcDir, IsDir: true}}

		if err := moveFiles(items, dstDir); err != nil {
			t.Fatalf("moveFiles failed for directory: %v", err)
		}

		movedSub := filepath.Join(dstDir, "subfolder", "sub.txt")
		movedNested := filepath.Join(dstDir, "subfolder", "nested", "nested.txt")

		data1, err := os.ReadFile(movedSub)
		if err != nil || string(data1) != "sub content" {
			t.Fatalf("expected 'sub content', got %s (err: %v)", string(data1), err)
		}
		data2, err := os.ReadFile(movedNested)
		if err != nil || string(data2) != "nested content" {
			t.Fatalf("expected 'nested content', got %s (err: %v)", string(data2), err)
		}

		if _, err := os.Stat(subDir); !os.IsNotExist(err) {
			t.Fatalf("expected source directory to be removed, got err: %v", err)
		}
	})

	t.Run("move to same file errors", func(t *testing.T) {
		dir := t.TempDir()
		file := filepath.Join(dir, "self.txt")
		if err := os.WriteFile(file, []byte("same"), 0o644); err != nil {
			t.Fatal(err)
		}

		items := []FileInfo{{Name: "self.txt", Path: dir}}
		if err := moveFiles(items, dir); err == nil {
			t.Fatal("expected error when moving file to itself, got nil")
		}
	})

	t.Run("move directory into itself errors", func(t *testing.T) {
		dir := t.TempDir()
		parent := filepath.Join(dir, "parent")
		child := filepath.Join(parent, "child")
		if err := os.MkdirAll(child, 0o755); err != nil {
			t.Fatal(err)
		}

		items := []FileInfo{{Name: "parent", Path: dir, IsDir: true}}
		if err := moveFiles(items, parent); err == nil {
			t.Fatal("expected error when moving directory into itself, got nil")
		}
	})

	t.Run("move symlink", func(t *testing.T) {
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
		if err := moveFiles(items, dstDir); err != nil {
			t.Fatalf("moveFiles failed for symlink: %v", err)
		}

		dstLink := filepath.Join(dstDir, "link.txt")
		linkTarget, err := os.Readlink(dstLink)
		if err != nil {
			t.Fatalf("failed reading moved symlink: %v", err)
		}
		if linkTarget != targetFile {
			t.Fatalf("expected symlink target %s, got %s", targetFile, linkTarget)
		}

		if _, err := os.Lstat(symlinkPath); !os.IsNotExist(err) {
			t.Fatalf("expected source symlink to be removed")
		}
	})
}

