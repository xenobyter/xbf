package main

import (
	"github.com/rivo/tview"
)

// App holds the application's UI components and shared state.
type App struct {
	tApp    *tview.Application
	tPages  *tview.Pages
	tHeader *tview.TextView
	tFooter *tview.TextView
	tLeft   *tview.List
	tRight  *tview.List

	i18n *I18n

	leftWd  string // Current working directory for the left pane
	rightWd string // Current working directory for the right pane
}

// newApp creates a new application instance with all UI components
// initialized in their default state.
func newApp() *App {
	a := &App{
		tApp:    tview.NewApplication(),
		tHeader: tview.NewTextView(),
		tFooter: tview.NewTextView(),
		tLeft:   tview.NewList().ShowSecondaryText(false),
		tRight:  tview.NewList().ShowSecondaryText(false),
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

	a.tPages = tview.NewPages().
		AddPage("main", grid, true, true)

	return a.tApp.SetRoot(a.tPages, true).SetFocus(a.tLeft).Run()
}
