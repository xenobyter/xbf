package main

import (
	"path/filepath"

	"github.com/gdamore/tcell/v2"
)

// actionFunc represents a handler for a keyboard action.
type actionFunc func(a *App, idx int)

// normalizeAction maps a key event to a symbolic action name.
func normalizeAction(event *tcell.EventKey) string {
	switch {
	case event.Key() == tcell.KeyEscape, event.Rune() == 'q':
		return "quit"
	case event.Key() == tcell.KeyTab:
		return "tab"
	case event.Key() == tcell.KeyRight:
		return "right"
	case event.Key() == tcell.KeyLeft:
		return "left"
	case event.Rune() == 's', event.Rune() == ' ':
		return "select"
	case event.Rune() == 'c':
		return "copy"
	case event.Rune() == 'm':
		return "move"
	case event.Rune() == 'r':
		return "rename"
	}
	return ""
}

// dirListKeyActions maps action names to handlers for the directory view.
var dirListKeyActions = map[string]actionFunc{
	"quit": func(a *App, _ int) {
		a.tApp.Stop()
	},
	"tab": func(a *App, _ int) {
		a.switchPane()
	},
	"right": func(a *App, idx int) {
		activeList := a.getActiveList()
		a.enterDirectory(activeList, idx)
		a.selection.Clear(activeList)
	},
	"left": func(a *App, _ int) {
		activeList := a.getActiveList()
		a.leaveDirectory(activeList)
		a.selection.Clear(activeList)
	},
	"select": func(a *App, idx int) {
		activeList := a.getActiveList()
		a.selection.ToggelItem(a.getItemsFor(activeList)[idx], activeList)
		a.updateListItem(activeList, idx)
	},
	"copy": func(a *App, _ int) {
		a.copySelected()
	},
	"move": func(a *App, _ int) {
		a.moveSelected()
	},
	"rename": func(a *App, idx int) {
		activeList := a.getActiveList()
		item := a.getItemsFor(activeList)[idx]
		a.ShowInputDialog(a.i18n.T(MsgTitleInputRename),
			item.Name, // Initialer Wert
			func(newName string) { // Callback bei Enter
				if err := renameEntry(filepath.Join(item.Path, item.Name), filepath.Join(item.Path, newName)); err != nil {
					a.setFooter(a.i18n.T(MsgErrRenameFile, item.Name, newName), tcell.ColorRed)
					return
				}
				a.setFooter(a.i18n.T(MsgRenameFile, item.Name, newName), tcell.ColorGreen)
				a.refreshPanes()
			})

	},
}

// handleInput dispatches keyboard events to registered action handlers.
//
// Unhandled events are returned so tview can process them normally.
func (a *App) handleInput(event *tcell.EventKey) *tcell.EventKey {
	action := normalizeAction(event)
	idx := a.getActiveList().GetCurrentItem()

	if fn, ok := dirListKeyActions[action]; ok {
		fn(a, idx)
		return nil // event consumed
	}

	return event // fall through for unhandled keys
}
