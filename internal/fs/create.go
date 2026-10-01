package fs

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	errEmptyEntryName   = errors.New("entry name cannot be empty")
	errInvalidEntryName = errors.New("invalid entry name")
)

var (
	ErrEmptyEntryName   = errEmptyEntryName
	ErrInvalidEntryName = errInvalidEntryName
)

func ValidateNewEntryName(name string) error {
	return validateNewEntryName(name)
}

func EntryExists(path string) bool {
	return entryExists(path)
}

func CreateEmptyFile(path string) error {
	return createEmptyFile(path)
}

func CreateDirectory(path string) error {
	return createDirectory(path)
}

func validateNewEntryName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errEmptyEntryName
	}
	if filepath.Base(name) != name || name == "." || name == ".." {
		return errInvalidEntryName
	}
	return nil
}

func entryExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func createEmptyFile(path string) error {
	if strings.TrimSpace(path) == "" {
		return os.ErrInvalid
	}
	if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
		return err
	}
	return nil
}

func createDirectory(path string) error {
	if strings.TrimSpace(path) == "" {
		return os.ErrInvalid
	}
	if err := os.Mkdir(path, 0o755); err != nil {
		return err
	}
	return nil
}
