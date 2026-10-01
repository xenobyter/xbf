package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rivo/tview"
)

func TestResolveDirectoryPath(t *testing.T) {
	baseDir := t.TempDir()
	childDir := filepath.Join(baseDir, "child")
	if err := os.Mkdir(childDir, 0o755); err != nil {
		t.Fatal(err)
	}

	homeDir := filepath.Join(baseDir, "home")
	if err := os.Mkdir(homeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", homeDir)

	a := &App{}

	resolved, err := a.resolveDirectoryPath(baseDir, "child")
	if err != nil {
		t.Fatalf("resolveDirectoryPath relative failed: %v", err)
	}
	if resolved != childDir {
		t.Fatalf("resolveDirectoryPath relative = %q, want %q", resolved, childDir)
	}

	resolved, err = a.resolveDirectoryPath(baseDir, "~")
	if err != nil {
		t.Fatalf("resolveDirectoryPath home failed: %v", err)
	}
	if resolved != homeDir {
		t.Fatalf("resolveDirectoryPath home = %q, want %q", resolved, homeDir)
	}
}

func TestGoHomeUsesUserHomeDir(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	a := &App{
		tApp:   tview.NewApplication(),
		tHeader: tview.NewTextView(),
		tLeft:  tview.NewList(),
		tRight: tview.NewList(),
		tFooter: tview.NewTextView(),
		i18n:   newI18n(),
		leftWd: "/tmp/left",
	}

	a.goHome(a.tLeft)
	if got := a.getWdFor(a.tLeft); got != homeDir {
		t.Fatalf("goHome() set wd = %q, want %q", got, homeDir)
	}
}