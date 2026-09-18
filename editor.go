package main

import (
	"os"
	"path/filepath"
	"strings"
)

// Document is the editor-ready model for a file and acts as a base for future
// text editing features.
type Document struct {
	Path      string
	Name      string
	Data      []byte
	Mode      previewMode
	Size      int64
	Truncated bool
	Dirty     bool
	CursorRow int
	CursorCol int
}

func newDocument(path string) (*Document, error) {
	data, size, truncated, err := loadPreviewData(path)
	if err != nil {
		return nil, err
	}

	mode := detectPreviewMode(data)
	return &Document{
		Path:      path,
		Name:      filepath.Base(path),
		Data:      data,
		Mode:      mode,
		Size:      size,
		Truncated: truncated,
		CursorRow: 0,
		CursorCol: 0,
	}, nil
}

func (d *Document) Text() string {
	return string(d.Data)
}

func (d *Document) SetText(text string) {
	d.Data = []byte(text)
	d.Dirty = true
}

func (d *Document) Save() error {
	if d == nil || d.Path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(d.Path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(d.Path, d.Data, 0o644); err != nil {
		return err
	}
	d.Dirty = false
	return nil
}

func (d *Document) Lines() []string {
	if len(d.Data) == 0 {
		return []string{""}
	}
	return strings.Split(string(d.Data), "\n")
}
