package main

import (
	"github.com/gdamore/tcell/v2"
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

	a.tPages = tview.NewPages().
		AddPage("main", grid, true, true)

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

// showModal adds a primitive as a modal page over the current view and sets focus.
func (a *App) showModal(name string, item tview.Primitive, focus tview.Primitive) {
	a.tPages.RemovePage(name)
	a.tPages.AddPage(name, item, true, true)
	if focus == nil {
		focus = item
	}
	a.tApp.SetFocus(focus)
}

// hideModal removes a modal page and returns focus to the active pane.
func (a *App) hideModal(name string) {
	a.tPages.RemovePage(name)
	a.tApp.SetFocus(a.getActiveList())
}

// ShowInputDialog displays an input modal, prompts for text, and calls onSubmit upon Enter.
func (a *App) ShowInputDialog(title string, initialValue string, onSubmit func(text string)) {
	inputField := tview.NewInputField().
		SetFieldWidth(40).
		SetText(initialValue).
		SetFieldBackgroundColor(tcell.ColorBlack)

	inputField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			text := inputField.GetText()
			a.hideModal("inputModal")
			if onSubmit != nil {
				onSubmit(text)
			}
		case tcell.KeyEscape:
			a.hideModal("inputModal")
		}
	})

	inputModal := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputField, 0, 1, true)
	inputModal.
		SetBorder(true).
		SetTitle(title).
		SetTitleAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorBlack)
	inputModal.Box.SetBackgroundColor(tcell.ColorBlack)
	inputModal.Box.SetBorderColor(tcell.ColorWhite)

	inputGrid := tview.NewGrid().
		SetColumns(0, 50, 0).
		SetRows(0, 3, 0).
		AddItem(inputModal, 1, 1, 1, 1, 0, 0, true)

	a.showModal("inputModal", inputGrid, inputField)
}

// ShowConfirmDialog displays a modal confirmation dialog with the given message and buttons.
// defaultFocus sets which button is focused initially (e.g. 0 for Cancel).
// onDone is called with the chosen button index.
func (a *App) ShowConfirmDialog(message string, buttons []string, defaultFocus int, onDone func(buttonIndex int)) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons(buttons).
		SetFocus(defaultFocus).
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(tcell.ColorBlack).
		SetButtonBackgroundColor(tcell.ColorBlack).
		SetButtonTextColor(tcell.ColorWhite).
		SetButtonStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)).
		SetButtonActivatedStyle(tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack))

	modal.Box.SetBackgroundColor(tcell.ColorBlack)
	modal.Box.SetBorderColor(tcell.ColorWhite)
	modal.Box.SetBorderStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack))

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		a.hideModal("confirmModal")
		if onDone != nil {
			onDone(buttonIndex)
		}
	})

	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyEscape || event.Rune() == 'n' || event.Rune() == 'N':
			a.hideModal("confirmModal")
			return nil
		case event.Rune() == 'y' || event.Rune() == 'Y' || event.Rune() == 'j' || event.Rune() == 'J':
			a.hideModal("confirmModal")
			if onDone != nil {
				onDone(1)
			}
			return nil
		}
		return event
	})

	a.showModal("confirmModal", modal, modal)
}
