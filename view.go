package main

import (
	"github.com/rivo/tview"
)

// updateList replaces the contents of the given list with the
// entries from the specified directory.
//
// If path is empty, the current working directory is used.
// Any errors encountered while resolving the directory or reading
// its contents are displayed as list entries.
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

	for _, item := range items {
		list.AddItem(item, "", 0, nil)
	}
}

// setHeader updates the header with the given path.
func (a *App) setHeader(path string) {
	a.tHeader.SetText(a.i18n.T(MsgCurrentWd, path))
}
