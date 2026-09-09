package main

import (
	"github.com/rivo/tview"
)

// selectionKey creates a unique key from name and pane
type selectionKey struct {
	name string
	pane *tview.List
}

// Selection represents items selected across the panes.
// It tracks which items are selected and from which pane they were selected.
type Selection struct {
	// items maps (name, pane) combination to the selected FileInfo
	items map[selectionKey]FileInfo
}

// ToggelItem adds a file to the selection if it's not already present,
// or removes it if it is (toggle behavior).
func (s *Selection) ToggelItem(item FileInfo, pane *tview.List) {
	if s.items == nil {
		s.items = make(map[selectionKey]FileInfo)
	}

	// Create a unique key from name and pane
	key := selectionKey{name: item.Name, pane: pane}

	if _, exists := s.items[key]; exists {
		// Item already selected, remove it (toggle)
		delete(s.items, key)
	} else {
		// Item not selected, add it
		s.items[key] = item
	}
}

// IsSelected checks if an item is currently selected in a specific pane.
func (s *Selection) IsSelected(name string, pane *tview.List) bool {
	if s.items == nil {
		return false
	}
	key := selectionKey{name: name, pane: pane}
	_, exists := s.items[key]
	return exists
}

// Clear removes all selected items for a specific pane.
func (s *Selection) Clear(pane *tview.List) {
	if s.items == nil {
		return
	}

	// Delete only items from the specified pane
	for key := range s.items {
		if key.pane == pane {
			delete(s.items, key)
		}
	}
}

// Count returns the number of selected items.
func (s *Selection) Count() int {
	if s.items == nil {
		return 0
	}
	return len(s.items)
}

// GetSelectedItems returns a slice of selected FileInfo items for a specific pane.
func (s *Selection) GetSelectedItems(pane *tview.List) []FileInfo {
	if s.items == nil {
		return nil
	}

	var selected []FileInfo
	for key, item := range s.items {
		if key.pane == pane {
			selected = append(selected, item)
		}
	}
	return selected
}
