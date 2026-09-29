package main

import (
	"fmt"
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

	a.setListItems(list, items)
}

func (a *App) setListItems(list *tview.List, items []FileInfo) {
	if list == a.tLeft {
		a.leftItems = items
	} else {
		a.rightItems = items
	}
	list.Clear()
	for _, item := range items {
		icon := IconFile
		if item.IsDir {
			icon = IconFolder
		}
		text := icon + item.Name
		if a.selection.IsSelected(item.Name, list) {
			text = "[red]" + text
		}
		list.AddItem(text, "", 0, nil)
	}
}

func (a *App) searchActiveList() {
	list := a.getActiveList()
	items := append([]FileInfo(nil), a.getItemsFor(list)...)
	matches := items
	inputField := tview.NewInputField().
		SetFieldWidth(48).
		SetFieldBackgroundColor(tcell.ColorBlack)
	results := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFocusOnly(false)
	status := tview.NewTextView().SetTextAlign(tview.AlignCenter)

	updateResults := func(query string) {
		matches = filterFiles(items, query)
		results.Clear()
		for _, item := range matches {
			icon := IconFile
			if item.IsDir {
				icon = IconFolder
			}
			results.AddItem(icon+item.Name, "", 0, nil)
		}
		if len(matches) == 0 {
			status.SetText(a.i18n.T(MsgSearchNoMatches))
			return
		}
		results.SetCurrentItem(0)
		status.SetText(a.i18n.T(MsgSearchMatches, len(matches)))
	}
	updateResults("")
	inputField.SetChangedFunc(updateResults)
	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			moveSearchSelection(results, -1)
			return nil
		case tcell.KeyDown:
			moveSearchSelection(results, 1)
			return nil
		default:
			return event
		}
	})
	inputField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			index := results.GetCurrentItem()
			if index < 0 || index >= len(matches) {
				return
			}
			selected := matches[index]
			a.hideModal("searchModal")
			for itemIndex, item := range a.getItemsFor(list) {
				if item.Name == selected.Name {
					list.SetCurrentItem(itemIndex)
					a.setItemInfo(list, itemIndex)
					return
				}
			}
		case tcell.KeyEscape:
			a.hideModal("searchModal")
		}
	})

	searchModal := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputField, 1, 0, true).
		AddItem(results, 0, 1, false).
		AddItem(status, 1, 0, false)
	searchModal.SetBorder(true).SetTitle(a.i18n.T(MsgTitleSearch)).SetTitleAlign(tview.AlignCenter)
	searchModal.Box.SetBackgroundColor(tcell.ColorBlack)
	searchModal.Box.SetBorderColor(tcell.ColorWhite)
	searchGrid := tview.NewGrid().
		SetColumns(0, 60, 0).
		SetRows(0, 14, 0).
		AddItem(searchModal, 1, 1, 1, 1, 0, 0, true)
	a.showModal("searchModal", searchGrid, inputField)
}

func moveSearchSelection(results *tview.List, delta int) {
	count := results.GetItemCount()
	if count == 0 {
		return
	}
	index := results.GetCurrentItem() + delta
	if index < 0 {
		index = 0
	} else if index >= count {
		index = count - 1
	}
	results.SetCurrentItem(index)
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
	a.setItemInfo(target, target.GetCurrentItem())
}

// changeDir updates the working directory for the list, refreshes the view,
// and updates the header.
func (a *App) changeDir(list *tview.List, wd string) {
	a.setWdFor(list, wd)
	a.updateList(list, wd)
	a.setHeader(wd)
	a.setItemInfo(list, list.GetCurrentItem())
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

// setItemInfo displays metadata for the item currently under the cursor.
func (a *App) setItemInfo(list *tview.List, idx int) {
	items := a.getItemsFor(list)
	if idx < 0 || idx >= len(items) {
		a.setFooter("", tcell.ColorWhite)
		return
	}

	item := items[idx]
	if item.IsDir {
		a.setFooter(a.i18n.T(MsgDirectoryInfo, item.Name), tcell.ColorWhite)
		return
	}

	a.setFooter(a.i18n.T(MsgFileInfo, item.Name, formatFileSize(item.Size)), tcell.ColorWhite)
}

func formatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}

	units := []string{"KB", "MB", "GB", "TB"}
	value := float64(size)
	for _, unit := range units {
		value /= 1024
		if value < 1024 || unit == units[len(units)-1] {
			return fmt.Sprintf("%.1f %s", value, unit)
		}
	}

	return fmt.Sprintf("%d B", size)
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

func (a *App) resolveItems(list *tview.List) []FileInfo {
	items := a.selection.GetSelectedItems(list)
	if len(items) > 0 {
		return items
	}

	paneItems := a.getItemsFor(list)
	idx := list.GetCurrentItem()
	if idx < 0 || idx >= len(paneItems) {
		return nil
	}
	return []FileInfo{paneItems[idx]}
}

// copySelected copies all selected files from the active pane to the target pane.
//
// If no files are selected it copies only the currently highlighted file. Errors are displayed in the footer.
func (a *App) copySelected() {
	activeList := a.getActiveList()
	items := a.resolveItems(activeList)
	if len(items) == 0 {
		return
	}
	target := a.getTargetWd()

	if err := copyFiles(items, target); err != nil {
		a.setFooter(a.i18n.T(MsgErrCopyFile)+": "+err.Error(), tcell.ColorRed)
		return
	}
	a.refreshPanes()
	a.setFooter(a.i18n.T(MsgCopyFile, len(items), target), tcell.ColorGreen)
}

// moveSelected moves all selected files from the active pane to the target pane.
//
// If no files are selected it moves only the currently highlighted file. Errors are displayed in the footer.
func (a *App) moveSelected() {
	activeList := a.getActiveList()
	items := a.resolveItems(activeList)
	if len(items) == 0 {
		return
	}
	target := a.getTargetWd()

	if err := moveFiles(items, target); err != nil {
		a.setFooter(a.i18n.T(MsgErrMoveFile)+": "+err.Error(), tcell.ColorRed)
		return
	}
	a.refreshPanes()
	a.setFooter(a.i18n.T(MsgMoveFile, len(items), target), tcell.ColorGreen)
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

	a.refreshPanes()
	a.setFooter(a.i18n.T(MsgRenameFile, item.Name, newName), tcell.ColorGreen)
}

// deleteSelected prompts the user for confirmation and deletes selected (or highlighted) files.
func (a *App) deleteSelected() {
	activeList := a.getActiveList()
	items := a.resolveItems(activeList)
	if len(items) == 0 {
		return
	}

	var message string
	if len(items) == 1 {
		message = a.i18n.T(MsgConfirmDeleteOne, items[0].Name)
	} else {
		message = a.i18n.T(MsgConfirmDeleteMulti, len(items))
	}

	btnCancel := a.i18n.T(MsgBtnCancel)
	btnDelete := a.i18n.T(MsgBtnDelete)

	// Buttons: [Cancel, Delete]. Index 0 (Cancel) is focused by default for safety.
	a.ShowConfirmDialog(message, []string{btnCancel, btnDelete}, 0, func(btnIndex int) {
		if btnIndex == 1 {
			if err := deleteFiles(items); err != nil {
				a.setFooter(a.i18n.T(MsgErrDeleteFile)+": "+err.Error(), tcell.ColorRed)
				return
			}
			a.refreshPanes()
			a.setFooter(a.i18n.T(MsgDeleteFile, len(items)), tcell.ColorGreen)
		}
	})
}
