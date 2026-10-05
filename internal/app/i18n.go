package app

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
	MsgFileInfo
	MsgDirectoryInfo
	MsgErrGetWd
	MsgErrReadDir
	MsgErrCopyFile
	MsgCopyFile
	MsgErrMoveFile
	MsgMoveFile
	MsgTitleInputRename
	MsgErrRenameFile
	MsgRenameFile
	MsgErrInvalidName
	MsgErrSameSrcDst
	MsgErrSameFile
	MsgErrOverwriteNonDirWithDir
	MsgErrOverwriteDirWithFile
	MsgRenameDirIntoItself
	MsgRenameOpFailed
	MsgErrDeleteFile
	MsgDeleteFile
	MsgConfirmDeleteOne
	MsgConfirmDeleteMulti
	MsgBtnCancel
	MsgBtnDelete
	MsgErrPreviewOpen
	MsgErrPreviewDirectory
	MsgErrEditDirectory
	MsgPreviewHeader
	MsgPreviewModeText
	MsgPreviewModeHex
	MsgPreviewTruncated
	MsgTitleInputNewFile
	MsgTitleInputGoTo
	MsgErrEmptyFileName
	MsgErrCreateFile
	MsgCreateFile
	MsgTitleInputNewDir
	MsgErrEmptyDirName
	MsgErrCreateDir
	MsgCreateDir
	MsgTitleSearch
	MsgErrGoToPath
	MsgSearchNoMatches
	MsgSearchMatches
	MsgEditorTextTitle
	MsgEditorHexTitle
	MsgErrEditorSave
	MsgEditorUnsaved
	MsgEditorSaved
	MsgHiddenFilesShown
	MsgHiddenFilesHidden
	MsgTitleInputCommand
	MsgErrEmptyCommand
	MsgErrCommandRun
	MsgCommandRan
	MsgErrInteractiveCommandBlocked
	MsgErrCommandSuspend
	MsgCommandOutputTitle
	MsgCommandOutputHint
	MsgCommandOutputEmpty
	MsgCommandOutputTruncated
)

var translations = map[Lang]map[MsgKey]string{
	LangEN: {
		MsgCurrentWd:           "Current Working Directory: %s",
		MsgFileInfo:            "%s | %s",
		MsgDirectoryInfo:       "%s | Directory",
		MsgErrGetWd:            "Error retrieving working directory",
		MsgErrReadDir:          "Error reading directory",
		MsgErrCopyFile:         "Error copying",
		MsgCopyFile:            "Copied %d items to %s",
		MsgErrMoveFile:         "Error moving",
		MsgMoveFile:            "Moved %d items to %s",
		MsgTitleInputRename:    "Rename Item",
		MsgErrRenameFile:       "rename from %s to %s failed",
		MsgRenameFile:          "Renamed %s to %s",
		MsgErrDeleteFile:       "Error deleting",
		MsgDeleteFile:          "Deleted %d items",
		MsgConfirmDeleteOne:    "Delete \"%s\"?",
		MsgConfirmDeleteMulti:  "Delete %d selected items?",
		MsgBtnCancel:           "Cancel",
		MsgBtnDelete:           "Delete",
		MsgErrPreviewOpen:      "Error opening preview",
		MsgErrPreviewDirectory: "Cannot preview directories",
		MsgErrEditDirectory:    "Cannot edit directories",
		MsgPreviewHeader:       "Preview: %s | %s | %s%s",
		MsgPreviewModeText:     "TEXT",
		MsgPreviewModeHex:      "HEX",
		MsgPreviewTruncated:    " | truncated to %s",
		MsgTitleInputNewFile:   "New File",
		MsgTitleInputGoTo:      "Go to Path",
		MsgErrEmptyFileName:    "File name cannot be empty",
		MsgErrCreateFile:       "Error creating file",
		MsgCreateFile:          "Created file: %s",
		MsgTitleInputNewDir:    "New Directory",
		MsgErrEmptyDirName:     "Directory name cannot be empty",
		MsgErrCreateDir:        "Error creating directory",
		MsgCreateDir:           "Created directory: %s",
		MsgTitleSearch:         "Search current directory",
		MsgErrGoToPath:         "Error opening path",
		MsgSearchNoMatches:     "No matching files",
		MsgSearchMatches:       "Results: %d | Up/Down to select, Enter to confirm",
		MsgEditorTextTitle:     "Edit: %s",
		MsgEditorHexTitle:      "Hex editor: %s",
		MsgErrEditorSave:       "Save failed: %s",
		MsgEditorUnsaved:       "Unsaved changes",
		MsgEditorSaved:         "Saved: %s",
		MsgHiddenFilesShown:    "Hidden files: shown",
		MsgHiddenFilesHidden:   "Hidden files: hidden",
		MsgTitleInputCommand:   "Run Command",
		MsgErrEmptyCommand:     "Command cannot be empty",
		MsgErrCommandRun:       "Command failed: %s",
		MsgCommandRan:          "Command finished: %s",
		MsgErrInteractiveCommandBlocked: "Interactive commands are blocked in this mode",
		MsgErrCommandSuspend:   "Could not suspend UI for command execution",
		MsgCommandOutputTitle:  "Command Output",
		MsgCommandOutputHint:   "Esc/q close | Up/Down scroll | PgUp/PgDn page",
		MsgCommandOutputEmpty:  "(no output)",
		MsgCommandOutputTruncated: "Output truncated",
	},
	LangDE: {
		MsgCurrentWd:           "Aktuelles Arbeitsverzeichnis: %s",
		MsgFileInfo:            "%s | %s",
		MsgDirectoryInfo:       "%s | Verzeichnis",
		MsgErrGetWd:            "Fehler beim Ermitteln des Arbeitsverzeichnisses",
		MsgErrReadDir:          "Fehler beim Lesen des Verzeichnisses",
		MsgErrCopyFile:         "Fehler beim Kopieren",
		MsgCopyFile:            "%d Elemente nach %s kopiert",
		MsgErrMoveFile:         "Fehler beim Verschieben",
		MsgMoveFile:            "%d Elemente nach %s verschoben",
		MsgTitleInputRename:    "Element umbenennen",
		MsgErrRenameFile:       "Umbenennen von %s nach %s fehlgeschlagen",
		MsgRenameFile:          "%s in %s umbenannt",
		MsgErrDeleteFile:       "Fehler beim Löschen",
		MsgDeleteFile:          "%d Elemente gelöscht",
		MsgConfirmDeleteOne:    "\"%s\" wirklich löschen?",
		MsgConfirmDeleteMulti:  "%d ausgewählte Elemente wirklich löschen?",
		MsgBtnCancel:           "Abbrechen",
		MsgBtnDelete:           "Löschen",
		MsgErrPreviewOpen:      "Fehler beim Öffnen der Vorschau",
		MsgErrPreviewDirectory: "Verzeichnisse können nicht in der Vorschau angezeigt werden",
		MsgErrEditDirectory:    "Verzeichnisse können nicht bearbeitet werden",
		MsgPreviewHeader:       "Vorschau: %s | %s | %s%s",
		MsgPreviewModeText:     "TEXT",
		MsgPreviewModeHex:      "HEX",
		MsgPreviewTruncated:    " | gekürzt auf %s",
		MsgTitleInputNewFile:   "Neue Datei",
		MsgTitleInputGoTo:      "Zu Pfad springen",
		MsgErrEmptyFileName:    "Dateiname darf nicht leer sein",
		MsgErrCreateFile:       "Fehler beim Erstellen der Datei",
		MsgCreateFile:          "Datei erstellt: %s",
		MsgTitleInputNewDir:    "Neues Verzeichnis",
		MsgErrEmptyDirName:     "Verzeichnisname darf nicht leer sein",
		MsgErrCreateDir:        "Fehler beim Erstellen des Verzeichnisses",
		MsgCreateDir:           "Verzeichnis erstellt: %s",
		MsgTitleSearch:         "Aktuelles Verzeichnis durchsuchen",
		MsgErrGoToPath:         "Fehler beim Öffnen des Pfads",
		MsgSearchNoMatches:     "Keine passenden Einträge",
		MsgSearchMatches:       "%d Treffer | Hoch/Runter wählen, Enter bestätigen",
		MsgEditorTextTitle:     "Bearbeiten: %s",
		MsgEditorHexTitle:      "Hex-Editor: %s",
		MsgErrEditorSave:       "Speichern fehlgeschlagen: %s",
		MsgEditorUnsaved:       "Ungespeicherte Änderungen",
		MsgEditorSaved:         "Gespeichert: %s",
		MsgHiddenFilesShown:    "Versteckte Dateien: angezeigt",
		MsgHiddenFilesHidden:   "Versteckte Dateien: ausgeblendet",
		MsgTitleInputCommand:   "Befehl ausführen",
		MsgErrEmptyCommand:     "Befehl darf nicht leer sein",
		MsgErrCommandRun:       "Befehl fehlgeschlagen: %s",
		MsgCommandRan:          "Befehl abgeschlossen: %s",
		MsgErrInteractiveCommandBlocked: "Interaktive Befehle sind in diesem Modus gesperrt",
		MsgErrCommandSuspend:   "UI konnte für die Befehlsausführung nicht pausiert werden",
		MsgCommandOutputTitle:  "Befehlsausgabe",
		MsgCommandOutputHint:   "Esc/q schließen | Hoch/Runter scrollen | PgUp/PgDn seitenweise",
		MsgCommandOutputEmpty:  "(keine Ausgabe)",
		MsgCommandOutputTruncated: "Ausgabe gekürzt",
	},
}

func defaultMsg(key MsgKey) string {
	return translations[LangEN][key]
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
