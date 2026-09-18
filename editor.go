package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

// HexEditor renders and edits raw bytes in a compact hex view.
type HexEditor struct {
	*tview.Box
	doc         *Document
	cursorByte  int
	cursorNibble int
	onSave      func()
	onClose     func()
	onDirty     func()
}

func newHexEditor(doc *Document, onSave func(), onClose func(), onDirty func()) *HexEditor {
	h := &HexEditor{
		Box:        tview.NewBox(),
		doc:        doc,
		onSave:     onSave,
		onClose:    onClose,
		onDirty:    onDirty,
		cursorByte: 0,
	}
	h.SetBorder(true)
	return h
}

func (h *HexEditor) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch {
		case event.Key() == tcell.KeyEscape:
			if h.onClose != nil {
				h.onClose()
			}
		case event.Key() == tcell.KeyCtrlS:
			if h.onSave != nil {
				h.onSave()
			}
		case event.Key() == tcell.KeyLeft:
			h.moveCursor(-1)
		case event.Key() == tcell.KeyRight:
			h.moveCursor(1)
		case event.Key() == tcell.KeyUp:
			h.moveCursor(-16)
		case event.Key() == tcell.KeyDown:
			h.moveCursor(16)
		case event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2:
			h.backspaceNibble()
		case event.Rune() == '0' || event.Rune() == '1' || event.Rune() == '2' || event.Rune() == '3' || event.Rune() == '4' || event.Rune() == '5' || event.Rune() == '6' || event.Rune() == '7' || event.Rune() == '8' || event.Rune() == '9' || event.Rune() == 'a' || event.Rune() == 'A' || event.Rune() == 'b' || event.Rune() == 'B' || event.Rune() == 'c' || event.Rune() == 'C' || event.Rune() == 'd' || event.Rune() == 'D' || event.Rune() == 'e' || event.Rune() == 'E' || event.Rune() == 'f' || event.Rune() == 'F':
			h.writeHexByte(event.Rune())
		}
	}
}

func (h *HexEditor) moveCursor(delta int) {
	if h.doc == nil || len(h.doc.Data) == 0 {
		return
	}
	newPos := h.cursorByte + delta
	if newPos < 0 {
		newPos = 0
	}
	if newPos >= len(h.doc.Data) {
		newPos = len(h.doc.Data) - 1
	}
	h.cursorByte = newPos
	h.cursorNibble = 0
}

func (h *HexEditor) backspaceNibble() {
	if h.doc == nil || len(h.doc.Data) == 0 {
		return
	}
	if h.cursorByte >= len(h.doc.Data) {
		h.cursorByte = len(h.doc.Data) - 1
	}
	if h.cursorNibble == 0 {
		h.doc.Data[h.cursorByte] &= 0x0F
	} else {
		h.doc.Data[h.cursorByte] &= 0xF0
	}
	h.markDirty()
	h.cursorNibble = 0
}

func (h *HexEditor) writeHexByte(r rune) {
	if h.doc == nil || len(h.doc.Data) == 0 {
		return
	}
	if h.cursorByte >= len(h.doc.Data) {
		return
	}
	val := hexValue(r)
	if h.cursorNibble == 0 {
		h.doc.Data[h.cursorByte] = (val << 4) | (h.doc.Data[h.cursorByte] & 0x0F)
		h.cursorNibble = 1
		h.markDirty()
		return
	}
	h.doc.Data[h.cursorByte] = (h.doc.Data[h.cursorByte] & 0xF0) | val
	h.cursorNibble = 0
	h.cursorByte++
	if h.cursorByte >= len(h.doc.Data) {
		h.cursorByte = len(h.doc.Data) - 1
	}
	h.markDirty()
}

func hexValue(r rune) byte {
	switch {
	case r >= '0' && r <= '9':
		return byte(r - '0')
	case r >= 'a' && r <= 'f':
		return byte(r - 'a' + 10)
	case r >= 'A' && r <= 'F':
		return byte(r - 'A' + 10)
	default:
		return 0
	}
}

func (h *HexEditor) markDirty() {
	if h.doc != nil {
		h.doc.Dirty = true
	}
	if h.onDirty != nil {
		h.onDirty()
	}
}

func (h *HexEditor) Draw(screen tcell.Screen) {
	h.Box.DrawForSubclass(screen, h)
	if h.doc == nil || len(h.doc.Data) == 0 {
		return
	}

	innerX, innerY, innerW, innerH := h.GetInnerRect()
	lineY := innerY
	for rowStart := 0; rowStart < len(h.doc.Data) && lineY < innerY+innerH; rowStart += 16 {
		rowEnd := rowStart + 16
		if rowEnd > len(h.doc.Data) {
			rowEnd = len(h.doc.Data)
		}

		line := fmt.Sprintf("%08x  ", rowStart)
		for i := 0; i < 16; i++ {
			if i < rowEnd-rowStart {
				idx := rowStart + i
				v := h.doc.Data[idx]
				text := fmt.Sprintf("%02x", v)
				if idx == h.cursorByte {
					line += "[" + text + "]"
				} else {
					line += text
				}
			} else {
				line += "  "
			}
			if i == 7 {
				line += " "
			}
			line += " "
		}

		line += " |"
		for i := rowStart; i < rowEnd; i++ {
			c := h.doc.Data[i]
			if c >= 32 && c <= 126 {
				line += string(rune(c))
			} else {
				line += "."
			}
		}
		line += "|"

		if innerW > 0 {
			if len(line) > innerW {
				line = line[:innerW]
			}
			drawText(screen, innerX, lineY, line, tcell.StyleDefault)
		}
		lineY++
	}
}

func drawText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, r := range text {
		screen.SetContent(x+i, y, r, nil, style)
	}
}
