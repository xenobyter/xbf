package app

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const commandOutputMaxBytes = 256 * 1024

var interactiveCommands = map[string]struct{}{
	"vi":         {},
	"vim":        {},
	"nvim":       {},
	"nano":       {},
	"emacs":      {},
	"less":       {},
	"more":       {},
	"man":        {},
	"top":        {},
	"htop":       {},
	"btop":       {},
	"watch":      {},
	"ssh":        {},
	"sftp":       {},
	"tmux":       {},
	"screen":     {},
	"fzf":        {},
	"lazygit":    {},
	"tig":        {},
	"gitk":       {},
	"python":     {},
	"python3":    {},
	"ipython":    {},
	"node":       {},
	"mysql":      {},
	"psql":       {},
	"sqlite3":    {},
	"kubectl":    {},
	"k9s":        {},
	"journalctl": {},
}

type commandResult struct {
	stdout    string
	stderr    string
	err       error
	truncated bool
}

func (a *App) openShell() {
	activeList := a.getActiveList()
	workingDir := a.getWdFor(activeList)

	a.ShowInputDialog(a.i18n.T(MsgTitleInputCommand), "", func(command string) {
		command = strings.TrimSpace(command)
		if command == "" {
			a.setFooter(a.i18n.T(MsgErrEmptyCommand), tcell.ColorRed)
			return
		}
		if isInteractiveCommand(command) {
			a.setFooter(a.i18n.T(MsgErrInteractiveCommandBlocked), tcell.ColorRed)
			return
		}

		var result commandResult
		suspended := a.tApp.Suspend(func() {
			result = runCommand(workingDir, command)
		})
		if !suspended {
			a.setFooter(a.i18n.T(MsgErrCommandSuspend), tcell.ColorRed)
			return
		}

		a.refreshPanes()
		a.showCommandOutput(command, workingDir, result)

		if result.err != nil {
			a.setFooter(a.i18n.T(MsgErrCommandRun, result.err.Error()), tcell.ColorRed)
		} else {
			a.setFooter(a.i18n.T(MsgCommandRan, command), tcell.ColorGreen)
		}
	})

}

func isInteractiveCommand(command string) bool {
	tokens := tokenizeCommand(command)
	for _, token := range tokens {
		lower := strings.ToLower(token)
		if _, blocked := interactiveCommands[lower]; blocked {
			return true
		}
	}
	return false
}

func tokenizeCommand(command string) []string {
	replacer := strings.NewReplacer(
		"&&", " ",
		"||", " ",
		";", " ",
		"|", " ",
	)
	clean := replacer.Replace(command)
	parts := strings.Fields(clean)
	tokens := make([]string, 0, len(parts))

	for _, part := range parts {
		if part == "sudo" || strings.Contains(part, "=") {
			continue
		}
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		part = strings.Trim(part, "\"'`")
		part = strings.TrimPrefix(part, "./")
		part = strings.TrimSuffix(part, ",")
		if part == "" {
			continue
		}
		tokens = append(tokens, part)
	}

	return tokens
}

func runCommand(wd, command string) commandResult {
	cmd := newCommand(command)
	cmd.Dir = wd

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	out, outTruncated := limitString(stdout.String(), commandOutputMaxBytes)
	errOut, errTruncated := limitString(stderr.String(), commandOutputMaxBytes)

	return commandResult{
		stdout:    out,
		stderr:    errOut,
		err:       err,
		truncated: outTruncated || errTruncated,
	}
}

func newCommand(command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		shell := strings.TrimSpace(os.Getenv("COMSPEC"))
		if shell == "" {
			shell = "cmd"
		}
		return exec.Command(shell, "/C", command)
	}

	shell := strings.TrimSpace(os.Getenv("SHELL"))
	if shell == "" {
		shell = "/bin/sh"
	}
	return exec.Command(shell, "-c", command)

}

func limitString(text string, max int) (string, bool) {
	if len(text) <= max {
		return text, false
	}
	return text[:max], true
}

func (a *App) showCommandOutput(command, wd string, result commandResult) {
	body := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(false)
	body.SetBorder(false)
	body.SetText(buildCommandOutput(command, wd, result, a.i18n.T(MsgCommandOutputEmpty), a.i18n.T(MsgCommandOutputTruncated)))

	hint := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(a.i18n.T(MsgCommandOutputHint))

	body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, col := body.GetScrollOffset()
		switch event.Key() {
		case tcell.KeyUp:
			if row > 0 {
				body.ScrollTo(row-1, col)
			}
			return nil
		case tcell.KeyDown:
			body.ScrollTo(row+1, col)
			return nil
		case tcell.KeyPgUp:
			step := 10
			next := row - step
			if next < 0 {
				next = 0
			}
			body.ScrollTo(next, col)
			return nil
		case tcell.KeyPgDn:
			body.ScrollTo(row+10, col)
			return nil
		case tcell.KeyHome:
			body.ScrollToBeginning()
			return nil
		case tcell.KeyEnd:
			body.ScrollToEnd()
			return nil
		case tcell.KeyEscape:
			a.hideModal("commandOutputModal")
			return nil
		}

		if event.Rune() == 'q' {
			a.hideModal("commandOutputModal")
			return nil
		}

		return event
	})

	content := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(body, 0, 1, true).
		AddItem(hint, 1, 0, false)
	content.SetBorder(true).
		SetTitle(a.i18n.T(MsgCommandOutputTitle)).
		SetTitleAlign(tview.AlignCenter)
	content.Box.SetBackgroundColor(tcell.ColorBlack)
	content.Box.SetBorderColor(tcell.ColorWhite)

	grid := tview.NewGrid().
		SetColumns(0, 100, 0).
		SetRows(0, 28, 0).
		AddItem(content, 1, 1, 1, 1, 0, 0, true)

	a.showModal("commandOutputModal", grid, body)
}

func buildCommandOutput(command, wd string, result commandResult, emptyText, truncatedText string) string {
	var b strings.Builder
	b.WriteString("$ ")
	b.WriteString(command)
	b.WriteString("\n")
	b.WriteString("cwd: ")
	b.WriteString(wd)
	b.WriteString("\n")
	if result.err != nil {
		b.WriteString("status: ")
		b.WriteString(result.err.Error())
	} else {
		b.WriteString("status: ok")
	}
	b.WriteString("\n\n")

	b.WriteString("stdout:\n")
	if strings.TrimSpace(result.stdout) == "" {
		b.WriteString(emptyText)
	} else {
		b.WriteString(result.stdout)
	}
	b.WriteString("\n\n")

	b.WriteString("stderr:\n")
	if strings.TrimSpace(result.stderr) == "" {
		b.WriteString(emptyText)
	} else {
		b.WriteString(result.stderr)
	}

	if result.truncated {
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("%s (%d bytes)", truncatedText, commandOutputMaxBytes))
	}

	return b.String()
}



