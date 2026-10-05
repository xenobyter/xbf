package app

import (
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/xenobyter/xbf/internal/fs"
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
	case event.Rune() == 'd', event.Key() == tcell.KeyDelete:
		return "delete"
	case event.Rune() == 'n':
		return "newfile"
	case event.Rune() == 'N':
		return "newdir"
	case event.Rune() == 'p':
		return "preview"
	case event.Rune() == 'g':
		return "goto"
	case event.Rune() == 'e', event.Rune() == 'E':
		return "editor"
	case event.Rune() == '~':
		return "home"
	case event.Rune() == '/':
		return "search"
	case event.Rune() == '.':
		return "toggle_hidden"
	case event.Rune() == '!':
		return "shell"
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
		a.selection.ToggleItem(a.getItemsFor(activeList)[idx], activeList)
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
				if err := fs.RenameEntry(filepath.Join(item.Path, item.Name), filepath.Join(item.Path, newName)); err != nil {
					a.setFooter(a.i18n.T(MsgErrRenameFile, item.Name, newName), tcell.ColorRed)
					return
				}
				a.refreshPanes()
				a.setFooter(a.i18n.T(MsgRenameFile, item.Name, newName), tcell.ColorGreen)
			})

	},
	"delete": func(a *App, _ int) {
		a.deleteSelected()
	},
	"newfile": func(a *App, _ int) {
		activeList := a.getActiveList()
		wd := a.getWdFor(activeList)
		a.ShowInputDialog(a.i18n.T(MsgTitleInputNewFile), "new_file.txt", func(fileName string) {
			fileName = strings.TrimSpace(fileName)
			if err := fs.ValidateNewEntryName(fileName); err != nil {
				msgKey := MsgErrInvalidName
				if err == fs.ErrEmptyEntryName {
					msgKey = MsgErrEmptyFileName
				}
				a.setFooter(a.i18n.T(msgKey), tcell.ColorRed)
				return
			}
			newFilePath := filepath.Join(wd, fileName)
			if fs.EntryExists(newFilePath) {
				a.setFooter(a.i18n.T(MsgErrCreateFile, fileName), tcell.ColorRed)
				return
			}
			if err := fs.CreateEmptyFile(newFilePath); err != nil {
				a.setFooter(a.i18n.T(MsgErrCreateFile, fileName)+": "+err.Error(), tcell.ColorRed)
				return
			}
			a.refreshPanes()
			a.setFooter(a.i18n.T(MsgCreateFile, fileName), tcell.ColorGreen)
		})
	},
	"newdir": func(a *App, _ int) {
		activeList := a.getActiveList()
		wd := a.getWdFor(activeList)
		a.ShowInputDialog(a.i18n.T(MsgTitleInputNewDir), "new_directory", func(dirName string) {
			dirName = strings.TrimSpace(dirName)
			if err := fs.ValidateNewEntryName(dirName); err != nil {
				msgKey := MsgErrInvalidName
				if err == fs.ErrEmptyEntryName {
					msgKey = MsgErrEmptyDirName
				}
				a.setFooter(a.i18n.T(msgKey), tcell.ColorRed)
				return
			}
			newDirPath := filepath.Join(wd, dirName)
			if fs.EntryExists(newDirPath) {
				a.setFooter(a.i18n.T(MsgErrCreateDir, dirName), tcell.ColorRed)
				return
			}
			if err := fs.CreateDirectory(newDirPath); err != nil {
				a.setFooter(a.i18n.T(MsgErrCreateDir, dirName)+": "+err.Error(), tcell.ColorRed)
				return
			}
			a.refreshPanes()
			a.setFooter(a.i18n.T(MsgCreateDir, dirName), tcell.ColorGreen)
		})
	},
	"preview": func(a *App, _ int) {
		a.openPreview()
	},
	"goto": func(a *App, _ int) {
		activeList := a.getActiveList()
		a.ShowInputDialog(a.i18n.T(MsgTitleInputGoTo), a.getWdFor(activeList), func(path string) {
			a.goToPath(activeList, path)
		})
	},
	"home": func(a *App, _ int) {
		activeList := a.getActiveList()
		a.goHome(activeList)
	},
	"editor": func(a *App, _ int) {
		a.openEditor()
	},
	"search": func(a *App, _ int) {
		a.searchActiveList()
	},
	"toggle_hidden": func(a *App, _ int) {
		a.toggleHidden()
	},
	"shell": func(a *App, _ int) {
		a.openShell()
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
