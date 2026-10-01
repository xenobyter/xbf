package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rivo/tview"
)

func TestToggleHidden(t *testing.T) {
	tempDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tempDir, "visible.txt"), []byte("v"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, ".hidden.txt"), []byte("h"), 0o644); err != nil {
		t.Fatal(err)
	}

	a := &App{
		tApp:    tview.NewApplication(),
		tLeft:   tview.NewList(),
		tRight:  tview.NewList(),
		tFooter: tview.NewTextView(),
		i18n:    newI18n(),
		leftWd:  tempDir,
		rightWd: tempDir,
	}

	// Initially showHidden is false
	if a.showHidden {
		t.Fatalf("expected showHidden to be false initially")
	}

	a.updateList(a.tLeft, a.leftWd)
	if len(a.leftItems) != 1 || a.leftItems[0].Name != "visible.txt" {
		t.Fatalf("expected only visible.txt initially, got: %+v", a.leftItems)
	}

	// Toggle hidden files ON
	a.toggleHidden()
	if !a.showHidden {
		t.Fatalf("expected showHidden to be true after toggle")
	}
	if len(a.leftItems) != 2 {
		t.Fatalf("expected 2 items after toggling on, got: %+v", a.leftItems)
	}
	footerText := a.tFooter.GetText(true)
	if !strings.Contains(footerText, "Versteckte Dateien: angezeigt") && !strings.Contains(footerText, "Hidden files: shown") {
		t.Errorf("unexpected footer text after toggle on: %q", footerText)
	}

	// Toggle hidden files OFF
	a.toggleHidden()
	if a.showHidden {
		t.Fatalf("expected showHidden to be false after second toggle")
	}
	if len(a.leftItems) != 1 || a.leftItems[0].Name != "visible.txt" {
		t.Fatalf("expected only visible.txt after toggling off, got: %+v", a.leftItems)
	}
	footerText = a.tFooter.GetText(true)
	if !strings.Contains(footerText, "Versteckte Dateien: ausgeblendet") && !strings.Contains(footerText, "Hidden files: hidden") {
		t.Errorf("unexpected footer text after toggle off: %q", footerText)
	}
}

