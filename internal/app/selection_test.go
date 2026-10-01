package app

import (
	"testing"

	"github.com/rivo/tview"
)

func TestSelection(t *testing.T) {
	var s Selection
	pane := tview.NewList()

	item1 := FileInfo{Name: "doc.txt", Path: "/home/user", IsDir: false, Size: 123}
	item2 := FileInfo{Name: "folder", Path: "/home/user", IsDir: true, Size: 4096}

	if s.Count() != 0 {
		t.Fatalf("expected count 0, got %d", s.Count())
	}

	// Toggle item1 on
	s.ToggleItem(item1, pane)
	if !s.IsSelected("doc.txt", pane) {
		t.Errorf("expected doc.txt to be selected")
	}
	if s.Count() != 1 {
		t.Errorf("expected count 1, got %d", s.Count())
	}

	// Toggle item2 on
	s.ToggleItem(item2, pane)
	if !s.IsSelected("folder", pane) {
		t.Errorf("expected folder to be selected")
	}
	if s.Count() != 2 {
		t.Errorf("expected count 2, got %d", s.Count())
	}

	// Verify GetSelectedItems preserves FileInfo fields
	items := s.GetSelectedItems(pane)
	if len(items) != 2 {
		t.Fatalf("expected 2 selected items, got %d", len(items))
	}

	foundDoc := false
	foundFolder := false
	for _, it := range items {
		if it.Name == "doc.txt" {
			foundDoc = true
			if it.Path != "/home/user" || it.IsDir != false || it.Size != 123 {
				t.Errorf("unexpected doc.txt item fields: %+v", it)
			}
		}
		if it.Name == "folder" {
			foundFolder = true
			if it.Path != "/home/user" || it.IsDir != true || it.Size != 4096 {
				t.Errorf("unexpected folder item fields: %+v", it)
			}
		}
	}
	if !foundDoc || !foundFolder {
		t.Errorf("expected to find both items, foundDoc=%v, foundFolder=%v", foundDoc, foundFolder)
	}

	// Toggle item1 off
	s.ToggleItem(item1, pane)
	if s.IsSelected("doc.txt", pane) {
		t.Errorf("expected doc.txt to no longer be selected")
	}
	if s.Count() != 1 {
		t.Errorf("expected count 1, got %d", s.Count())
	}

	// Clear pane
	s.Clear(pane)
	if s.Count() != 0 {
		t.Errorf("expected count 0 after Clear, got %d", s.Count())
	}
	if len(s.GetSelectedItems(pane)) != 0 {
		t.Errorf("expected empty slice from GetSelectedItems after Clear")
	}
}

func TestResolveItems(t *testing.T) {
	pane := tview.NewList()
	pane.AddItem("first", "", 0, nil)
	pane.AddItem("second", "", 0, nil)
	items := []FileInfo{
		{Name: "first"},
		{Name: "second"},
	}
	app := &App{tLeft: pane, leftItems: items}
	pane.SetCurrentItem(1)

	resolved := app.resolveItems(pane)
	if len(resolved) != 1 || resolved[0].Name != "second" {
		t.Fatalf("resolveItems() = %#v, want the current item", resolved)
	}

	app.selection.ToggleItem(items[0], pane)
	resolved = app.resolveItems(pane)
	if len(resolved) != 1 || resolved[0].Name != "first" {
		t.Fatalf("resolveItems() = %#v, want selected items to take precedence", resolved)
	}

	emptyPane := tview.NewList()
	emptyApp := &App{tLeft: emptyPane}
	if resolved := emptyApp.resolveItems(emptyPane); len(resolved) != 0 {
		t.Fatalf("resolveItems() for empty list = %#v, want no items", resolved)
	}
}

