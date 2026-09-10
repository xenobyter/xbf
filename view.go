package main

import (
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
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
	target := a.getInactiveList()
	a.tApp.SetFocus(target)
	a.setHeader(a.getWdFor(target))
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

// setFooter sets the status message and color in the footer bar.
func (a *App) setFooter(msg string, color tcell.Color) {
	a.tFooter.SetText(msg).SetTextColor(color)
}

// refreshPanes clears the active pane's selection and reloads both pane lists.
//
// It reserves the currently selected Items
func (a *App) refreshPanes() {
	activePane := a.getActiveList()
	activeCurrentItem := activePane.GetCurrentItem()
	inactivePane := a.getInactiveList()
	inactiveCurrentItem := inactivePane.GetCurrentItem()

	a.selection.Clear(activePane)
	a.updateList(activePane, a.getWdFor(activePane))
	activePane.SetCurrentItem(activeCurrentItem)
	a.updateList(inactivePane, a.getWdFor(inactivePane))
	inactivePane.SetCurrentItem(inactiveCurrentItem)
}

// copySelected copies all selected files from the active pane to the target pane.
//
// If no files are selected it copies only the currently highlighted file. Errors are displayed in the footer.
func (a *App) copySelected() {
	activeList := a.getActiveList()
	items := a.selection.GetSelectedItems(activeList)
	if len(items) == 0 {
		idx := activeList.GetCurrentItem()
		items = []FileInfo{a.getItemsFor(activeList)[idx]}
	}
	target := a.getTargetWd()

	if err := copyFiles(items, target); err != nil {
		a.setFooter(a.i18n.T(MsgErrCopyFile)+": "+err.Error(), tcell.ColorRed)
		return
	}
	a.setFooter(a.i18n.T(MsgCopyFile, len(items), target), tcell.ColorGreen)
	a.refreshPanes()
}

// moveSelected moves all selected files from the active pane to the target pane.
//
// If no files are selected it moves only the currently highlighted file. Errors are displayed in the footer.
func (a *App) moveSelected() {
	activeList := a.getActiveList()
	items := a.selection.GetSelectedItems(activeList)
	if len(items) == 0 {
		idx := activeList.GetCurrentItem()
		items = []FileInfo{a.getItemsFor(activeList)[idx]}
	}
	target := a.getTargetWd()

	if err := moveFiles(items, target); err != nil {
		a.setFooter(a.i18n.T(MsgErrMoveFile)+": "+err.Error(), tcell.ColorRed)
		return
	}
	a.setFooter(a.i18n.T(MsgMoveFile, len(items), target), tcell.ColorGreen)
	a.refreshPanes()
}

// renameSelected renames the currently selected item in the active pane.
func (a *App) renameSelected(idx int, newName string) {
	activeList := a.getActiveList()
	items := a.getItemsFor(activeList)
	if idx < 0 || idx >= len(items) {
		return
	}

	newName = strings.TrimSpace(newName)
	if newName == "" {
		a.setFooter(a.i18n.T(MsgErrInvalidName), tcell.ColorRed)
		return
	}

	item := items[idx]
	if item.Name == newName {
		return
	}

	src := filepath.Join(item.Path, item.Name)
	dst := filepath.Join(item.Path, newName)

	if err := renameEntry(src, dst); err != nil {
		a.setFooter(a.i18n.T(MsgErrRenameFile, item.Name, newName), tcell.ColorRed)
		return
	}

	a.setFooter(a.i18n.T(MsgRenameFile, item.Name, newName), tcell.ColorGreen)
	a.refreshPanes()
}