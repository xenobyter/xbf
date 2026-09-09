package main

import (
	"path/filepath"

	"github.com/rivo/tview"
)

const (
	IconFolder = "\U0001F4C1" // 📁
	IconFile   = "\U0001F4C4" // 📄
)

// updateList replaces the contents of a list with the entries
// from the specified directory.
//
// If path is empty, the current working directory is used.
// Errors are displayed as list items.
func (a *App) updateList(list *tview.List, path string) {
	list.Clear()

	showError := func(msgKey MsgKey, err error) {
		list.AddItem(a.i18n.T(msgKey)+": "+err.Error(), "", 0, nil)
	}

	if path == "" {
		p, err := getWd()
		if err != nil {
			showError(MsgErrGetWd, err)
			return
		}
		path = p
	}

	items, err := readDir(path)
	if err != nil {
		showError(MsgErrReadDir, err)
		return
	}

	switch list {
	case a.tLeft:
		a.leftItems = items
	case a.tRight:
		a.rightItems = items
	}

	for _, item := range items {
		icon := IconFile
		if item.IsDir {
			icon = IconFolder
		}
		list.AddItem(icon+item.Name, "", 0, nil)
	}
}

// updateListItem updates a single item in the list with the current selection state.
// This avoids reloading the entire directory and preserves the cursor position.
func (a *App) updateListItem(list *tview.List, idx int) {
	items := a.getItemsFor(list)
	if idx < 0 || idx >= len(items) {
		return
	}

	item := items[idx]
	icon := IconFile
	if item.IsDir {
		icon = IconFolder
	}
	style := ""
	if a.selection.IsSelected(item.Name, list) {
		style = "[red]"
	}
	list.SetItemText(idx, style+icon+item.Name, "")
}

// setHeader updates the header text.
func (a *App) setHeader(path string) {
	a.tHeader.SetText(a.i18n.T(MsgCurrentWd, path))
}

// switchPane moves focus between the left and right directory panes
// and updates the header to reflect the active pane.
func (a *App) switchPane() {
	curr := a.tApp.GetFocus()

	var target *tview.List
	var wd string

	if curr == a.tLeft {
		target = a.tRight
		wd = a.rightWd
	} else {
		target = a.tLeft
		wd = a.leftWd
	}

	a.tApp.SetFocus(target)
	a.setHeader(wd)
}

// changeDir updates the working directory for the list, refreshes the view,
// and updates the header.
func (a *App) changeDir(list *tview.List, wd string) {
	a.setWdFor(list, wd)
	a.updateList(list, wd)
	a.setHeader(wd)
}

// enterDirectory navigates into the selected directory entry.
func (a *App) enterDirectory(list *tview.List, idx int) {
	if idx < 0 || idx >= list.GetItemCount() {
		return
	}

	items := a.getItemsFor(list)
	item := items[idx]
	if !item.IsDir {
		return
	}

	a.changeDir(list, filepath.Join(a.getWdFor(list), item.Name))
}

// leaveDirectory navigates to the parent Directory.
func (a *App) leaveDirectory(list *tview.List) {
	wd := a.getWdFor(list)
	parent := filepath.Dir(wd)
	if parent == wd {
		return // already at root
	}

	a.changeDir(list, parent)
}
