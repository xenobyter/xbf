package app

import (
	"testing"
)

func TestI18nTranslations(t *testing.T) {
	en := &I18n{lang: LangEN}
	if got := en.T(MsgErrReadDir); got != "Error reading directory" {
		t.Errorf("expected English text, got %q", got)
	}
	if got := en.T(MsgCurrentWd, "/tmp"); got != "Current Working Directory: /tmp" {
		t.Errorf("expected formatted English text, got %q", got)
	}
	if got := en.T(MsgErrRenameFile, "old.txt", "new.txt"); got != "rename from old.txt to new.txt failed" {
		t.Errorf("expected English rename error text, got %q", got)
	}
	if got := en.T(MsgRenameFile, "old.txt", "new.txt"); got != "Renamed old.txt to new.txt" {
		t.Errorf("expected formatted English rename text, got %q", got)
	}
	if got := en.T(MsgConfirmDeleteOne, "test.txt"); got != "Delete \"test.txt\"?" {
		t.Errorf("expected formatted English delete confirm text, got %q", got)
	}
	if got := en.T(MsgConfirmDeleteMulti, 3); got != "Delete 3 selected items?" {
		t.Errorf("expected formatted English delete multi confirm text, got %q", got)
	}
	if got := en.T(MsgDeleteFile, 2); got != "Deleted 2 items" {
		t.Errorf("expected formatted English delete success text, got %q", got)
	}
	if got := en.T(MsgEditorTextTitle, "/tmp/file.txt"); got != "Edit: /tmp/file.txt" {
		t.Errorf("expected formatted English editor title, got %q", got)
	}
	if got := en.T(MsgEditorHexTitle, "/tmp/file.bin"); got != "Hex editor: /tmp/file.bin" {
		t.Errorf("expected formatted English hex editor title, got %q", got)
	}
	if got := en.T(MsgErrEditorSave, "permission denied"); got != "Save failed: permission denied" {
		t.Errorf("expected formatted English save error, got %q", got)
	}
	if got := en.T(MsgEditorUnsaved); got != "Unsaved changes" {
		t.Errorf("expected English unsaved status, got %q", got)
	}
	if got := en.T(MsgEditorSaved, "/tmp/file.txt"); got != "Saved: /tmp/file.txt" {
		t.Errorf("expected formatted English saved status, got %q", got)
	}
	if got := en.T(MsgHiddenFilesShown); got != "Hidden files: shown" {
		t.Errorf("expected English hidden files shown text, got %q", got)
	}
	if got := en.T(MsgHiddenFilesHidden); got != "Hidden files: hidden" {
		t.Errorf("expected English hidden files hidden text, got %q", got)
	}
	if got := en.T(MsgErrShellSuspend); got != "Could not suspend UI for shell session" {
		t.Errorf("expected English shell suspend error, got %q", got)
	}
	if got := en.T(MsgErrShellRun, "exit status 1"); got != "Shell failed: exit status 1" {
		t.Errorf("expected formatted English shell error, got %q", got)
	}
	if got := en.T(MsgShellExited); got != "Shell session ended" {
		t.Errorf("expected English shell exit status, got %q", got)
	}

	de := &I18n{lang: LangDE}
	if got := de.T(MsgErrReadDir); got != "Fehler beim Lesen des Verzeichnisses" {
		t.Errorf("expected German text, got %q", got)
	}
	if got := de.T(MsgCurrentWd, "/tmp"); got != "Aktuelles Arbeitsverzeichnis: /tmp" {
		t.Errorf("expected formatted German text, got %q", got)
	}
	if got := de.T(MsgErrRenameFile, "alt.txt", "neu.txt"); got != "Umbenennen von alt.txt nach neu.txt fehlgeschlagen" {
		t.Errorf("expected German rename error text, got %q", got)
	}
	if got := de.T(MsgRenameFile, "alt.txt", "neu.txt"); got != "alt.txt in neu.txt umbenannt" {
		t.Errorf("expected formatted German rename text, got %q", got)
	}
	if got := de.T(MsgConfirmDeleteOne, "test.txt"); got != "\"test.txt\" wirklich löschen?" {
		t.Errorf("expected formatted German delete confirm text, got %q", got)
	}
	if got := de.T(MsgConfirmDeleteMulti, 3); got != "3 ausgewählte Elemente wirklich löschen?" {
		t.Errorf("expected formatted German delete multi confirm text, got %q", got)
	}
	if got := de.T(MsgDeleteFile, 2); got != "2 Elemente gelöscht" {
		t.Errorf("expected formatted German delete success text, got %q", got)
	}
	if got := de.T(MsgEditorTextTitle, "/tmp/datei.txt"); got != "Bearbeiten: /tmp/datei.txt" {
		t.Errorf("expected formatted German editor title, got %q", got)
	}
	if got := de.T(MsgEditorHexTitle, "/tmp/datei.bin"); got != "Hex-Editor: /tmp/datei.bin" {
		t.Errorf("expected formatted German hex editor title, got %q", got)
	}
	if got := de.T(MsgErrEditorSave, "Zugriff verweigert"); got != "Speichern fehlgeschlagen: Zugriff verweigert" {
		t.Errorf("expected formatted German save error, got %q", got)
	}
	if got := de.T(MsgEditorUnsaved); got != "Ungespeicherte Änderungen" {
		t.Errorf("expected German unsaved status, got %q", got)
	}
	if got := de.T(MsgEditorSaved, "/tmp/datei.txt"); got != "Gespeichert: /tmp/datei.txt" {
		t.Errorf("expected formatted German saved status, got %q", got)
	}
	if got := de.T(MsgHiddenFilesShown); got != "Versteckte Dateien: angezeigt" {
		t.Errorf("expected German hidden files shown text, got %q", got)
	}
	if got := de.T(MsgHiddenFilesHidden); got != "Versteckte Dateien: ausgeblendet" {
		t.Errorf("expected German hidden files hidden text, got %q", got)
	}
	if got := de.T(MsgErrShellSuspend); got != "UI konnte für die Shell-Sitzung nicht pausiert werden" {
		t.Errorf("expected German shell suspend error, got %q", got)
	}
	if got := de.T(MsgErrShellRun, "Status 1"); got != "Shell fehlgeschlagen: Status 1" {
		t.Errorf("expected formatted German shell error, got %q", got)
	}
	if got := de.T(MsgShellExited); got != "Shell-Sitzung beendet" {
		t.Errorf("expected German shell exit status, got %q", got)
	}
}

func TestI18nFallback(t *testing.T) {
	unknown := &I18n{lang: Lang("fr")}
	if got := unknown.T(MsgErrReadDir); got != "Error reading directory" {
		t.Errorf("expected fallback to English, got %q", got)
	}
}
