//go:build windows

package fs

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestIsHiddenWindows(t *testing.T) {
	if !isHidden(".git", "C:\\path", nil) {
		t.Errorf("expected .git to be hidden on Windows")
	}
	if isHidden("main.go", "C:\\path", nil) {
		t.Errorf("expected main.go not to be hidden on Windows")
	}

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "hidden_by_attr.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}

	ptr, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := syscall.SetFileAttributes(ptr, syscall.FILE_ATTRIBUTE_HIDDEN); err != nil {
		t.Skipf("cannot set hidden attribute: %v", err)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatal(err)
	}

	if !isHidden("hidden_by_attr.txt", tempDir, info) {
		t.Errorf("expected file with FILE_ATTRIBUTE_HIDDEN to be detected as hidden with info")
	}
	if !isHidden("hidden_by_attr.txt", tempDir, nil) {
		t.Errorf("expected file with FILE_ATTRIBUTE_HIDDEN to be detected as hidden without info")
	}
}

