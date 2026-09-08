package main

import "github.com/gdamore/tcell/v2"

// actionFunc abstracts a handler for a specific key action.
// It receives the *App instance and the currently selected index in the list.
// Returning a function pointer enables a clean map of key strings to behaviours.
type actionFunc func(a *App, idx int)

// normalizeAction converts a raw tcell.EventKey into the symbolic action key used
// by `keyActions`. This isolates platform‑specific key codes from the business
// logic, making the mapping table easier to understand and extend.
func normalizeAction(event *tcell.EventKey) string {
	switch {
	case event.Key() == tcell.KeyEscape, event.Rune() == 'q':
		return "quit"
	case event.Key() == tcell.KeyTab:
		return "tab"
	}
	return ""
}

// handleInput is the central entry point for all keyboard events.
// It normalises the raw tcell.EventKey into a symbolic action string and
// dispatches the request to the corresponding handler from the `keyActions`
// map. If the key does not map to a known action, the original event is
// returned unchanged so tview can handle default behaviour.
func (a *App) handleInput(event *tcell.EventKey) *tcell.EventKey {
	action := normalizeAction(event)
	idx := 0

	if fn, ok := dirListKeyActions[action]; ok {
		fn(a, idx)
		return nil // event consumed
	}

	return event // fall‑through for unhandled keys
}

// dirListKeyActions maps symbolic action names to concrete functions for the main UI
var dirListKeyActions = map[string]actionFunc{
	"quit": func(a *App, _ int) {
		// Gracefully stop the tview application.
		a.tApp.Stop()
	},
	"tab": func(a *App, _ int) {
		// Switch focus between the left and right panes.
		if a.tApp.GetFocus() == a.tLeft {
			a.tApp.SetFocus(a.tRight)
		} else {
			a.tApp.SetFocus(a.tLeft)
		}
	},
	// "right": func(a *App, idx int) {
	// 	// Open a directory when the selected entry is a folder.
	// 	if idx >= 0 && idx < len(a.dirEntries) {
	// 		item := a.dirEntries[idx]
	// 		if item.isDir && a.navigateToDir(item.name) {
	// 			a.dirPathView.SetText(a.currentDir)
	// 			a.updateDirEntries(true)
	// 		}
	// 	}
	// },
}
