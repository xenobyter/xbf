package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	previewSampleBytes       = 64 * 1024
	previewMaxBytes          = 2 * 1024 * 1024
	previewHighlightMaxBytes = 512 * 1024
)

type previewMode string

const (
	previewModeText previewMode = "text"
	previewModeHex  previewMode = "hex"
)

func (a *App) openPreview() {
	activeList := a.getActiveList()
	idx := activeList.GetCurrentItem()
	items := a.getItemsFor(activeList)
	if idx < 0 || idx >= len(items) {
		return
	}

	item := items[idx]
	if item.IsDir {
		a.setFooter(a.i18n.T(MsgErrPreviewDirectory), tcell.ColorYellow)
		return
	}

	path := filepath.Join(item.Path, item.Name)
	data, size, truncated, err := loadPreviewData(path)
	if err != nil {
		a.setFooter(a.i18n.T(MsgErrPreviewOpen)+": "+err.Error(), tcell.ColorRed)
		return
	}

	mode := detectPreviewMode(data)
	modeLabel := a.i18n.T(MsgPreviewModeText)
	body := ""
	if mode == previewModeHex {
		modeLabel = a.i18n.T(MsgPreviewModeHex)
		body = renderHex(data)
	} else {
		body = renderText(path, data)
	}

	truncatedSuffix := ""
	if truncated {
		truncatedSuffix = a.i18n.T(MsgPreviewTruncated, formatFileSize(previewMaxBytes))
	}

	a.tPreviewHeader.SetText(a.i18n.T(MsgPreviewHeader, path, modeLabel, formatFileSize(size), truncatedSuffix))
	a.tPreviewBody.SetText(body)
	a.tPreviewBody.ScrollToBeginning()

	a.previewReturn = activeList
	a.tPages.RemovePage("preview")
	a.tPages.AddPage("preview", a.tPreviewPage, true, true)
	a.tPages.SwitchToPage("preview")
	a.tApp.SetFocus(a.tPreviewBody)
}

func (a *App) closePreview() {
	if a.tPages == nil {
		return
	}

	a.tPages.RemovePage("preview")

	returnFocus := a.previewReturn
	a.previewReturn = nil
	if returnFocus == nil {
		returnFocus = a.tLeft
	}

	a.tApp.SetFocus(returnFocus)
	if list, ok := returnFocus.(*tview.List); ok {
		a.setHeader(a.getWdFor(list))
		a.setItemInfo(list, list.GetCurrentItem())
	}
}

func (a *App) handlePreviewInput(event *tcell.EventKey) *tcell.EventKey {
	switch {
	case event.Key() == tcell.KeyEscape:
		a.closePreview()
		return nil
	case event.Key() == tcell.KeyPgDn:
		a.scrollPreview(20)
		return nil
	case event.Key() == tcell.KeyPgUp:
		a.scrollPreview(-20)
		return nil
	case event.Key() == tcell.KeyDown:
		a.scrollPreview(1)
		return nil
	case event.Key() == tcell.KeyUp:
		a.scrollPreview(-1)
		return nil
	case event.Key() == tcell.KeyHome:
		a.tPreviewBody.ScrollToBeginning()
		return nil
	case event.Key() == tcell.KeyEnd:
		a.tPreviewBody.ScrollToEnd()
		return nil
	case event.Rune() == 'q' || event.Key() == tcell.KeyEscape || event.Rune() == 'p':
		a.closePreview()
		return nil
	}
	return event
}

func (a *App) scrollPreview(delta int) {
	row, col := a.tPreviewBody.GetScrollOffset()
	row += delta
	if row < 0 {
		row = 0
	}
	a.tPreviewBody.ScrollTo(row, col)
}

func loadPreviewData(path string) ([]byte, int64, bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, false, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, 0, false, err
	}

	buf, err := io.ReadAll(io.LimitReader(f, previewMaxBytes+1))
	if err != nil {
		return nil, 0, false, err
	}

	truncated := len(buf) > previewMaxBytes
	if truncated {
		buf = buf[:previewMaxBytes]
	}

	return buf, info.Size(), truncated, nil
}

func detectPreviewMode(data []byte) previewMode {
	if len(data) == 0 {
		return previewModeText
	}

	sample := data
	if len(sample) > previewSampleBytes {
		sample = sample[:previewSampleBytes]
	}

	controlCount := 0
	nullCount := 0
	for _, b := range sample {
		if b == 0 {
			nullCount++
		}
		if b < 32 && b != '\n' && b != '\r' && b != '\t' {
			controlCount++
		}
	}

	if nullCount > 0 {
		return previewModeHex
	}

	if !utf8.Valid(sample) {
		return previewModeHex
	}

	if float64(controlCount)/float64(len(sample)) > 0.10 {
		return previewModeHex
	}

	return previewModeText
}

func renderText(path string, data []byte) string {
	text := string(data)
	if len(data) > previewHighlightMaxBytes {
		return tview.Escape(text)
	}

	lexer := lexers.Match(path)
	if lexer == nil {
		lexer = lexers.Analyse(text)
	}
	if lexer == nil {
		return tview.Escape(text)
	}

	iter, err := lexer.Tokenise(nil, text)
	if err != nil {
		return tview.Escape(text)
	}

	var b strings.Builder
	for token := iter(); token != chroma.EOF; token = iter() {
		color := tokenColorTag(token.Type)
		if color != "" {
			b.WriteString(color)
		}
		b.WriteString(tview.Escape(token.Value))
		if color != "" {
			b.WriteString("[-]")
		}
	}

	return b.String()
}

func tokenColorTag(tt chroma.TokenType) string {
	typeName := tt.String()
	switch {
	case strings.HasPrefix(typeName, "Keyword"):
		return "[yellow]"
	case strings.HasPrefix(typeName, "NameFunction"):
		return "[deepskyblue]"
	case strings.HasPrefix(typeName, "NameClass"):
		return "[deepskyblue]"
	case strings.HasPrefix(typeName, "LiteralString"):
		return "[green]"
	case strings.HasPrefix(typeName, "LiteralNumber"):
		return "[turquoise]"
	case strings.HasPrefix(typeName, "Comment"):
		return "[gray]"
	case strings.HasPrefix(typeName, "Operator"):
		return "[gold]"
	case strings.HasPrefix(typeName, "NameBuiltin"):
		return "[orchid]"
	default:
		return ""
	}
}

func renderHex(data []byte) string {
	const width = 16
	var b strings.Builder
	for off := 0; off < len(data); off += width {
		end := off + width
		if end > len(data) {
			end = len(data)
		}
		chunk := data[off:end]

		b.WriteString(fmt.Sprintf("%08x  ", off))
		for i := 0; i < width; i++ {
			if i < len(chunk) {
				b.WriteString(fmt.Sprintf("%02x ", chunk[i]))
			} else {
				b.WriteString("   ")
			}
			if i == 7 {
				b.WriteByte(' ')
			}
		}

		b.WriteString(" |")
		for _, c := range chunk {
			if c >= 32 && c <= 126 {
				b.WriteByte(c)
			} else {
				b.WriteByte('.')
			}
		}
		b.WriteString("|\n")
	}
	return b.String()
}

func (a *App) openTextEditor() {
	activeList := a.getActiveList()
	idx := activeList.GetCurrentItem()
	items := a.getItemsFor(activeList)
	if idx < 0 || idx >= len(items) {
		return
	}

	item := items[idx]
	if item.IsDir {
		a.setFooter(a.i18n.T(MsgErrEditDirectory), tcell.ColorYellow)
		return
	}

	path := filepath.Join(item.Path, item.Name)
	doc, err := newDocument(path)
	if err != nil {
		a.setFooter(a.i18n.T(MsgErrPreviewOpen)+": "+err.Error(), tcell.ColorRed)
		return
	}
	if doc.Mode == previewModeHex {
		a.openHexEditor(doc)
		return
	}

	a.document = doc
	a.tEditorHeader.SetText("Edit: " + path)
	textArea := tview.NewTextArea()
	textArea.SetBorder(true)
	textArea.SetInputCapture(a.handleEditorInput)
	textArea.SetText(doc.Text(), false)
	textArea.SetChangedFunc(func() {
		if a.document != nil {
			a.document.SetText(textArea.GetText())
			a.handleDocumentDirty()
		}
	})
	a.editorReturn = activeList
	a.rebuildEditorPage(textArea)
	a.tPages.SwitchToPage("editor")
	a.tApp.SetFocus(textArea)
}

func (a *App) openHexEditor(doc *Document) {
	a.document = doc
	a.tEditorHeader.SetText("Hex Edit: " + doc.Path)
	hexEditor := newHexEditor(doc, a.saveEditor, a.closeEditor, a.handleDocumentDirty)
	hexEditor.SetBorder(true)
	a.editorReturn = a.getActiveList()
	a.rebuildEditorPage(hexEditor)
	a.tPages.SwitchToPage("editor")
	a.tApp.SetFocus(hexEditor)
}

func (a *App) closeEditor() {
	if a.tPages == nil {
		return
	}
	a.tPages.RemovePage("editor")
	returnFocus := a.editorReturn
	a.editorReturn = nil
	if returnFocus == nil {
		returnFocus = a.tLeft
	}
	a.tApp.SetFocus(returnFocus)
	if list, ok := returnFocus.(*tview.List); ok {
		a.setHeader(a.getWdFor(list))
		a.setItemInfo(list, list.GetCurrentItem())
	}
}

func (a *App) saveEditor() {
	if a.document == nil {
		return
	}
	if err := a.document.Save(); err != nil {
		a.setFooter("Save failed: "+err.Error(), tcell.ColorRed)
		return
	}
	a.document.Dirty = false
	a.handleDocumentSaved()
}

func (a *App) handleEditorInput(event *tcell.EventKey) *tcell.EventKey {
	switch {
	case event.Key() == tcell.KeyEscape:
		a.closeEditor()
		return nil
	case event.Key() == tcell.KeyCtrlS:
		a.saveEditor()
		return nil
	case event.Key() == tcell.KeyCtrlW:
		a.closeEditor()
		return nil
	case event.Rune() == 'q':
		a.tApp.Stop()
		return nil
	}
	return event
}

func (a *App) openEditor() {
	activeList := a.getActiveList()
	idx := activeList.GetCurrentItem()
	items := a.getItemsFor(activeList)
	if idx < 0 || idx >= len(items) {
		return
	}

	item := items[idx]
	if item.IsDir {
		a.setFooter(a.i18n.T(MsgErrEditDirectory), tcell.ColorYellow)
		return
	}

	path := filepath.Join(item.Path, item.Name)
	doc, err := newDocument(path)
	if err != nil {
		a.setFooter(a.i18n.T(MsgErrPreviewOpen)+": "+err.Error(), tcell.ColorRed)
		return
	}
	if doc.Mode == previewModeHex {
		a.openHexEditor(doc)
		return
	}
	a.openTextEditor()
}

func (a *App) handleDocumentDirty() {
	if a.document != nil && a.document.Dirty {
		a.setFooter("Unsaved changes", tcell.ColorYellow)
	}
}

func (a *App) handleDocumentSaved() {
	if a.document != nil && !a.document.Dirty {
		a.setFooter("Saved: "+a.document.Path, tcell.ColorGreen)
	}
}
