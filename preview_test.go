package main

import (
	"strings"
	"testing"
)

func TestDetectPreviewMode(t *testing.T) {
	t.Run("text", func(t *testing.T) {
		data := []byte("hello\nworld\n")
		if got := detectPreviewMode(data); got != previewModeText {
			t.Fatalf("expected text mode, got %q", got)
		}
	})

	t.Run("binary by null byte", func(t *testing.T) {
		data := []byte{0x41, 0x00, 0x42}
		if got := detectPreviewMode(data); got != previewModeHex {
			t.Fatalf("expected hex mode, got %q", got)
		}
	})

	t.Run("binary by invalid utf8", func(t *testing.T) {
		data := []byte{0xff, 0xfe, 0xfd}
		if got := detectPreviewMode(data); got != previewModeHex {
			t.Fatalf("expected hex mode, got %q", got)
		}
	})
}

func TestRenderHex(t *testing.T) {
	out := renderHex([]byte("ABC"))
	if !strings.Contains(out, "00000000") {
		t.Fatalf("expected offset in hex output, got %q", out)
	}
	if !strings.Contains(out, "41 42 43") {
		t.Fatalf("expected hex bytes in output, got %q", out)
	}
	if !strings.Contains(out, "|ABC|") {
		t.Fatalf("expected ascii column in output, got %q", out)
	}
}
