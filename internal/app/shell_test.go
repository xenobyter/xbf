package app

import (
	"runtime"
	"testing"
)

func TestNewShellUsesConfiguredShell(t *testing.T) {
	envName := "SHELL"
	shellPath := "/custom/user-shell"
	if runtime.GOOS == "windows" {
		envName = "COMSPEC"
		shellPath = `C:\custom\user-shell.exe`
	}
	t.Setenv(envName, shellPath)

	command := newShell()
	if len(command.Args) != 1 || command.Args[0] != shellPath {
		t.Fatalf("newShell() args = %v, want only configured shell %q", command.Args, shellPath)
	}
}
