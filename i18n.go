package main

import (
	"fmt"
	"os"
	"strings"
)

// Lang represents a supported language code.
type Lang string

const (
	LangEN Lang = "en"
	LangDE Lang = "de"
)

// MsgKey identifies a translatable string.
type MsgKey int

const (
	MsgCurrentWd MsgKey = iota
	MsgErrGetWd
	MsgErrReadDir
)

var translations = map[Lang]map[MsgKey]string{
	LangEN: {
		MsgCurrentWd:  "Current Working Directory: %s",
		MsgErrGetWd:   "Error retrieving working directory",
		MsgErrReadDir: "Error reading directory",
	},
	LangDE: {
		MsgCurrentWd:  "Aktuelles Arbeitsverzeichnis: %s",
		MsgErrGetWd:   "Fehler beim Ermitteln des Arbeitsverzeichnisses",
		MsgErrReadDir: "Fehler beim Lesen des Verzeichnisses",
	},
}

// I18n handles translations for the application.
type I18n struct {
	lang Lang
}

// newI18n creates an I18n instance with language detected from environment variables.
func newI18n() *I18n {
	env := os.Getenv("LC_ALL")
	if env == "" {
		env = os.Getenv("LANG")
	}

	lang := LangEN
	if strings.HasPrefix(strings.ToLower(env), "de") {
		lang = LangDE
	}

	return &I18n{lang: lang}
}

// T returns the translated string for the given key, formatting with optional args.
func (i *I18n) T(key MsgKey, args ...any) string {
	dict, ok := translations[i.lang]
	if !ok {
		dict = translations[LangEN]
	}

	text, ok := dict[key]
	if !ok {
		text = translations[LangEN][key]
	}

	if len(args) > 0 {
		return fmt.Sprintf(text, args...)
	}
	return text
}

