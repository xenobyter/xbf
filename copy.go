package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// copyFiles copies the specified items to the target directory.
func copyFiles(items []FileInfo, target string) error {
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

		if err := copyEntry(srcPath, dstPath); err != nil {
			return err
		}
	}

	return nil
}

// copyEntry copies a single filesystem entry (file, directory, or symlink) to dst.
func copyEntry(src, dst string) error {
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
		return copyDir(src, dst)
	}
	if srcStat.Mode()&os.ModeSymlink != 0 {
		return copySymlink(src, dst)
	}
	return copyFile(src, dst)
}

// copyFile copies a regular file from src to dst, preserving file mode permissions.
func copyFile(src, dst string) (err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	if fi, err := os.Lstat(dst); err == nil {
		if fi.IsDir() {
			return fmt.Errorf("cannot overwrite directory %s with file", dst)
		}
		_ = os.Remove(dst)
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode().Perm())
	if err != nil {
		return err
	}
	defer func() {
		if cErr := dstFile.Close(); err == nil {
			err = cErr
		}
	}()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	if err = os.Chmod(dst, srcInfo.Mode().Perm()); err != nil {
		return err
	}

	return dstFile.Sync()
}

// copyDir recursively copies a directory tree from src to dst.
func copyDir(src, dst string) error {
	if _, err := os.Stat(src); err != nil {
		return err
	}

	if dstStat, err := os.Lstat(dst); err == nil && !dstStat.IsDir() {
		return fmt.Errorf("cannot overwrite non-directory %s with directory", dst)
	}

	rel, err := filepath.Rel(src, dst)
	if err == nil && !strings.HasPrefix(rel, "..") {
		return fmt.Errorf("cannot copy directory into itself: %s", src)
	}

	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			return os.MkdirAll(targetPath, info.Mode().Perm())
		}

		if d.Type()&os.ModeSymlink != 0 {
			return copySymlink(path, targetPath)
		}

		return copyFile(path, targetPath)
	})
}

// copySymlink replicates a symbolic link at dst pointing to src's link target.
func copySymlink(src, dst string) error {
	target, err := os.Readlink(src)
	if err != nil {
		return err
	}
	_ = os.RemoveAll(dst)
	return os.Symlink(target, dst)
}

