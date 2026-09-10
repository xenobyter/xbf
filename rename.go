package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// renameEntry renames a single filesystem entry (file, directory, or symlink).
func renameEntry(src, dst string) error {
	srcStat, err := os.Lstat(src)
	if err != nil {
		return err
	}

	srcClean := filepath.Clean(src)
	dstClean := filepath.Clean(dst)
	if srcClean == dstClean {
		return errors.New("source and destination are the same")
	}

	if dstStat, err := os.Lstat(dst); err == nil {
		if os.SameFile(srcStat, dstStat) {
			return errors.New("source and destination are the same file")
		}

		if srcStat.IsDir() && !dstStat.IsDir() {
			return errors.New("cannot overwrite non-directory with directory")
		}
		if !srcStat.IsDir() && dstStat.IsDir() {
			return errors.New("cannot overwrite directory with file")
		}
	}

	if srcStat.IsDir() {
		rel, err := filepath.Rel(src, dst)
		if err == nil && !strings.HasPrefix(rel, "..") {
			return errors.New("cannot rename directory into itself")
		}
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	if err := os.Rename(src, dst); err != nil {
		return errors.New("rename from " + src + " to " + dst + " failed: " + err.Error())
	}

	return nil
}
