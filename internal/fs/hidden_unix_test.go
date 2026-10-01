//go:build !windows

package fs

import "testing"

func TestIsHiddenUnix(t *testing.T) {
	if !isHidden(".git", "/path", nil) {
		t.Errorf("expected .git to be hidden on Unix")
	}
	if !isHidden(".env", "/path", nil) {
		t.Errorf("expected .env to be hidden on Unix")
	}
	if isHidden("main.go", "/path", nil) {
		t.Errorf("expected main.go not to be hidden on Unix")
	}
}

