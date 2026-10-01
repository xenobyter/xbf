package app

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestNormalizeAction(t *testing.T) {
	tests := []struct {
		name     string
		event    *tcell.EventKey
		expected string
	}{
		{
			name:     "escape quits",
			event:    tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone),
			expected: "quit",
		},
		{
			name:     "q quits",
			event:    tcell.NewEventKey(tcell.KeyRune, 'q', tcell.ModNone),
			expected: "quit",
		},
		{
			name:     "tab switches pane",
			event:    tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone),
			expected: "tab",
		},
		{
			name:     "right arrow",
			event:    tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone),
			expected: "right",
		},
		{
			name:     "left arrow",
			event:    tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone),
			expected: "left",
		},
		{
			name:     "s selects",
			event:    tcell.NewEventKey(tcell.KeyRune, 's', tcell.ModNone),
			expected: "select",
		},
		{
			name:     "space selects",
			event:    tcell.NewEventKey(tcell.KeyRune, ' ', tcell.ModNone),
			expected: "select",
		},
		{
			name:     "c copies",
			event:    tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone),
			expected: "copy",
		},
		{
			name:     "m moves",
			event:    tcell.NewEventKey(tcell.KeyRune, 'm', tcell.ModNone),
			expected: "move",
		},
		{
			name:     "r renames",
			event:    tcell.NewEventKey(tcell.KeyRune, 'r', tcell.ModNone),
			expected: "rename",
		},
		{
			name:     "d deletes",
			event:    tcell.NewEventKey(tcell.KeyRune, 'd', tcell.ModNone),
			expected: "delete",
		},
		{
			name:     "n creates file",
			event:    tcell.NewEventKey(tcell.KeyRune, 'n', tcell.ModNone),
			expected: "newfile",
		},
		{
			name:     "N creates directory",
			event:    tcell.NewEventKey(tcell.KeyRune, 'N', tcell.ModNone),
			expected: "newdir",
		},
		{
			name:     "p previews",
			event:    tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			expected: "preview",
		},
		{
			name:     "g opens goto dialog",
			event:    tcell.NewEventKey(tcell.KeyRune, 'g', tcell.ModNone),
			expected: "goto",
		},
		{
			name:     "tilde goes home",
			event:    tcell.NewEventKey(tcell.KeyRune, '~', tcell.ModNone),
			expected: "home",
		},
		{
			name:     "slash searches",
			event:    tcell.NewEventKey(tcell.KeyRune, '/', tcell.ModNone),
			expected: "search",
		},
		{
			name:     "dot toggles hidden",
			event:    tcell.NewEventKey(tcell.KeyRune, '.', tcell.ModNone),
			expected: "toggle_hidden",
		},
		{
			name:     "KeyDelete deletes",
			event:    tcell.NewEventKey(tcell.KeyDelete, 0, tcell.ModNone),
			expected: "delete",
		},
		{
			name:     "D does not delete (only lowercase d or KeyDelete)",
			event:    tcell.NewEventKey(tcell.KeyRune, 'D', tcell.ModNone),
			expected: "",
		},
		{
			name:     "unknown key returns empty",
			event:    tcell.NewEventKey(tcell.KeyRune, 'x', tcell.ModNone),
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeAction(tc.event)
			if got != tc.expected {
				t.Errorf("normalizeAction() = %q, expected %q", got, tc.expected)
			}
		})
	}
}

func TestModalReturnsFocusToOriginalPane(t *testing.T) {
	a := &App{
		tApp:   tview.NewApplication(),
		tPages: tview.NewPages(),
		tLeft:  tview.NewList(),
		tRight: tview.NewList(),
	}
	a.tApp.SetFocus(a.tRight)

	a.showModal("testModal", tview.NewTextView(), nil)
	a.hideModal("testModal")

	if a.tApp.GetFocus() != a.tRight {
		t.Fatalf("expected focus to return to right pane, got %T", a.tApp.GetFocus())
	}
}

func TestMoveSearchSelection(t *testing.T) {
	results := tview.NewList()
	results.AddItem("one", "", 0, nil)
	results.AddItem("two", "", 0, nil)
	results.AddItem("three", "", 0, nil)

	moveSearchSelection(results, 1)
	if got := results.GetCurrentItem(); got != 1 {
		t.Fatalf("moving down selected item %d, want 1", got)
	}

	moveSearchSelection(results, 1)
	moveSearchSelection(results, 1)
	if got := results.GetCurrentItem(); got != 2 {
		t.Fatalf("moving down past the last result selected item %d, want 2", got)
	}

	moveSearchSelection(results, -1)
	if got := results.GetCurrentItem(); got != 1 {
		t.Fatalf("moving up selected item %d, want 1", got)
	}
}
