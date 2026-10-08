package app

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func (a *App) openShell() {
	workingDir := a.getWdFor(a.getActiveList())
	var shellErr error
	suspended := a.tApp.Suspend(func() {
		shellErr = runShell(workingDir)
	})
	if !suspended {
		a.setFooter(a.i18n.T(MsgErrShellSuspend), tcell.ColorRed)
		return
	}

	a.refreshPanes()
	if shellErr != nil {
		a.setFooter(a.i18n.T(MsgErrShellRun, shellErr.Error()), tcell.ColorRed)
	} else {
		a.setFooter(a.i18n.T(MsgShellExited), tcell.ColorGreen)
	}
}

func runShell(wd string) error {
	cmd := newShell()
	cmd.Dir = wd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func newShell() *exec.Cmd {
	if runtime.GOOS == "windows" {
		shell := strings.TrimSpace(os.Getenv("COMSPEC"))
		if shell == "" {
			shell = "cmd"
		}
		return exec.Command(shell)
	}

	shell := strings.TrimSpace(os.Getenv("SHELL"))
	if shell == "" {
		shell = "/bin/sh"
	}
	return exec.Command(shell)
}
