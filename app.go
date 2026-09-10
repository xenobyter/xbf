package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App holds the application's UI components and shared state.
type App struct {
	tApp        *tview.Application
	tPages      *tview.Pages
	tHeader     *tview.TextView
	tFooter     *tview.TextView
	tLeft       *tview.List
	tRight      *tview.List
	tInputModal *tview.Flex
	tInputField *tview.InputField

	i18n *I18n

	leftWd     string
	leftItems  []FileInfo
	rightWd    string
	rightItems []FileInfo

	selection Selection
}

// newApp creates a new application instance with all UI components
// initialized in their default state.
func newApp() *App {
	a := &App{
		tApp:    tview.NewApplication(),
		tHeader: tview.NewTextView(),
		tFooter: tview.NewTextView(),
		tLeft:   tview.NewList().ShowSecondaryText(false).SetSelectedFocusOnly(true),
		tRight:  tview.NewList().ShowSecondaryText(false).SetSelectedFocusOnly(true),
		i18n:    newI18n(),
	}

	// Attempt to get the current working directory. If it fails, set an error message
	// in both leftWd and rightWd to inform the user.
	if wd, err := getWd(); err != nil {
		errText := a.i18n.T(MsgErrGetWd) + ": " + err.Error()
		a.leftWd = errText
		a.rightWd = errText
	} else {
		a.leftWd = wd
		a.rightWd = wd
	}

	a.setHeader(a.leftWd)
	a.updateList(a.tLeft, "")
	a.updateList(a.tRight, "")

	// Register inputhandlers
	a.tLeft.SetInputCapture(a.handleInput)
	a.tRight.SetInputCapture(a.handleInput)

	return a
}

// run builds the main layout, registers it as the root view,
// and starts the application's event loop.
//
// The layout consists of a header, two side-by-side content panes,
// and a footer. Initial focus is set to the left pane.
func (a *App) run() error {
	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(0, 0).
		SetBorders(true).
		AddItem(a.tHeader, 0, 0, 1, 2, 0, 0, false).
		AddItem(a.tLeft, 1, 0, 1, 1, 0, 0, true).
		AddItem(a.tRight, 1, 1, 1, 1, 0, 0, false).
		AddItem(a.tFooter, 2, 0, 1, 2, 0, 0, false)

	a.tInputField = tview.NewInputField().
		SetFieldWidth(40).
		SetAcceptanceFunc(nil)

	// Rahmen & Container für das Eingabefeld
	a.tInputModal = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.tInputField, 0, 1, true)

	a.tInputModal.
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft)

	// Zentrierung über ein 3x3 Grid
	inputGrid := tview.NewGrid().
		SetColumns(0, 50, 0).
		SetRows(0, 3, 0).
		AddItem(a.tInputModal, 1, 1, 1, 1, 0, 0, true)

	a.tPages = tview.NewPages().
		AddPage("main", grid, true, true).
		AddPage("inputModal", inputGrid, true, false)

	return a.tApp.SetRoot(a.tPages, true).SetFocus(a.tLeft).Run()
}

func (a *App) getActiveList() *tview.List {
	if a.tApp.GetFocus() == a.tLeft {
		return a.tLeft
	}
	return a.tRight
}

// getInactiveList returns the list belonging to the inactive pane.
func (a *App) getInactiveList() *tview.List {
	if a.tApp.GetFocus() == a.tLeft {
		return a.tRight
	}
	return a.tLeft
}

// getTargetWd returns the working directory of the inactive pane.
func (a *App) getTargetWd() string {
	return a.getWdFor(a.getInactiveList())
}

// getItemsFor returns the file list associated with the given UI list.
func (a *App) getItemsFor(list *tview.List) []FileInfo {
	if list == a.tLeft {
		return a.leftItems
	}
	return a.rightItems
}

// getWdFor returns the working directory associated with the given UI list.
func (a *App) getWdFor(list *tview.List) string {
	if list == a.tLeft {
		return a.leftWd
	}
	return a.rightWd
}

// setWdFor sets the working directory for the given UI list.
func (a *App) setWdFor(list *tview.List, wd string) {
	if list == a.tLeft {
		a.leftWd = wd
	} else {
		a.rightWd = wd
	}
}

// ShowInputDialog blendet das Modal ein, setzt den Titel und führt bei Enter `onSubmit` aus.
func (a *App) ShowInputDialog(title string, initialValue string, onSubmit func(text string)) {
	// 1. Titel & Werte setzen
	a.tInputModal.SetTitle(title).SetTitleAlign(tview.AlignCenter)
	a.tInputField.SetText(initialValue).SetFieldBackgroundColor(tcell.ColorBlack)
	activePane:= a.getActiveList()

	// 2. Key-Events steuern
	a.tInputField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			text := a.tInputField.GetText()
			a.hideInputDialog(activePane)
			if onSubmit != nil {
				onSubmit(text)
			}
		case tcell.KeyEscape:
			a.hideInputDialog(activePane)
		}
	})

	// 3. Modal einblenden und Fokus übergeben
	a.tPages.ShowPage("inputModal")
	a.tApp.SetFocus(a.tInputField)
}

func (a *App) hideInputDialog(activePane *tview.List) {
	a.tPages.HidePage("inputModal")
	// Fokus zurück auf das aktive Panel setzen
	a.tApp.SetFocus(activePane)
}
