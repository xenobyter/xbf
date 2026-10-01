//go:build !windows

package fs

import (
	"io/fs"
	"strings"
)

// isHidden determines whether a file or directory is hidden on Unix/Linux systems (starts with '.').
func isHidden(name, path string, info fs.FileInfo) bool {
	return strings.HasPrefix(name, ".")
}

