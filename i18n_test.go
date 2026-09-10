package main

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
}

func TestI18nFallback(t *testing.T) {
	unknown := &I18n{lang: Lang("fr")}
	if got := unknown.T(MsgErrReadDir); got != "Error reading directory" {
		t.Errorf("expected fallback to English, got %q", got)
	}
}
