//go:build windows

package main

import (
	"io/fs"
	"path/filepath"
	"strings"
	"syscall"
)

// isHidden determines whether a file or directory is hidden on Windows systems.
// It returns true if the name starts with '.' or if the FILE_ATTRIBUTE_HIDDEN flag is set.
func isHidden(name, path string, info fs.FileInfo) bool {
	if strings.HasPrefix(name, ".") {
		return true
	}

	if info != nil {
		if sys, ok := info.Sys().(*syscall.Win32FileAttributeData); ok {
			return sys.FileAttributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0
		}
	}

	fullPath := filepath.Join(path, name)
	pointer, err := syscall.UTF16PtrFromString(fullPath)
	if err != nil {
		return false
	}

	attrs, err := syscall.GetFileAttributes(pointer)
	if err != nil || attrs == syscall.INVALID_FILE_ATTRIBUTES {
		return false
	}

	return attrs&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}

