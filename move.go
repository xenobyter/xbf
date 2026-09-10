package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// moveFiles moves the specified items to the target directory.
func moveFiles(items []FileInfo, target string) error {
	if len(items) == 0 {
		return nil
	}

	targetInfo, err := os.Stat(target)
	if err != nil {
		return err
	}
	if !targetInfo.IsDir() {
		return fmt.Errorf("target is not a directory: %s", target)
	}

	for _, item := range items {
		if item.Name == "" {
			continue
		}

		srcPath := filepath.Join(item.Path, item.Name)
		dstPath := filepath.Join(target, item.Name)

		if err := moveEntry(srcPath, dstPath); err != nil {
			return err
		}
	}

	return nil
}

// moveEntry moves a single filesystem entry (file, directory, or symlink) to dst.
func moveEntry(src, dst string) error {
	srcStat, err := os.Lstat(src)
	if err != nil {
		return err
	}

	srcClean := filepath.Clean(src)
	dstClean := filepath.Clean(dst)
	if srcClean == dstClean {
		return fmt.Errorf("source and destination are the same: %s", src)
	}

	if dstStat, err := os.Lstat(dst); err == nil {
		if os.SameFile(srcStat, dstStat) {
			return fmt.Errorf("source and destination are the same file: %s", src)
		}
	}

	if srcStat.IsDir() {
		rel, err := filepath.Rel(src, dst)
		if err == nil && !strings.HasPrefix(rel, "..") {
			return fmt.Errorf("cannot move directory into itself: %s", src)
		}

		if dstStat, err := os.Lstat(dst); err == nil && !dstStat.IsDir() {
			return fmt.Errorf("cannot overwrite non-directory %s with directory", dst)
		}
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	if err := copyEntry(src, dst); err != nil {
		return err
	}

	if srcStat.IsDir() {
		return os.RemoveAll(src)
	}

	return os.Remove(src)
}

// moveFile moves a regular file from src to dst. It attempts an atomic rename
// first, falling back to copyFile and removing the source if across filesystems.
func moveFile(src, dst string) error {
	srcStat, err := os.Lstat(src)
	if err != nil {
		return err
	}

	if srcStat.IsDir() {
		return fmt.Errorf("cannot move directory %s with moveFile", src)
	}

	srcClean := filepath.Clean(src)
	dstClean := filepath.Clean(dst)
	if srcClean == dstClean {
		return fmt.Errorf("source and destination are the same: %s", src)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	if dstStat, err := os.Lstat(dst); err == nil {
		if dstStat.IsDir() {
			return fmt.Errorf("cannot overwrite directory %s with file", dst)
		}
		if os.SameFile(srcStat, dstStat) {
			return fmt.Errorf("source and destination are the same file: %s", src)
		}
		_ = os.Remove(dst)
	}

	// Try atomic rename first (efficient when on the same filesystem)
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	// Fallback for cross-device moves
	if srcStat.Mode()&os.ModeSymlink != 0 {
		if err := copySymlink(src, dst); err != nil {
			return err
		}
	} else {
		if err := copyFile(src, dst); err != nil {
			return err
		}
	}

	return os.Remove(src)
}
