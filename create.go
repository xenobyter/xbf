package main

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
