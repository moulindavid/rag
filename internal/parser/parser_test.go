package parser

import (
	"strings"
	"testing"
)

func TestParse_Txt(t *testing.T) {
	text, err := Parse("doc.txt", strings.NewReader("hello world"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", text)
	}
}

func TestParse_TxtTrimsWhitespace(t *testing.T) {
	text, err := Parse("doc.txt", strings.NewReader("  hello  "))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "hello" {
		t.Errorf("expected %q, got %q", "hello", text)
	}
}

func TestParse_TxtNormalizesCRLF(t *testing.T) {
	text, err := Parse("doc.txt", strings.NewReader("line1\r\nline2"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(text, "\r\n") {
		t.Error("expected CRLF to be normalized to LF")
	}
	if text != "line1\nline2" {
		t.Errorf("expected %q, got %q", "line1\nline2", text)
	}
}

func TestParse_TxtEmpty(t *testing.T) {
	text, err := Parse("doc.txt", strings.NewReader("   "))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "" {
		t.Errorf("expected empty string after trim, got %q", text)
	}
}

func TestParse_UnsupportedExtension(t *testing.T) {
	_, err := Parse("data.csv", strings.NewReader("a,b,c"))
	if err == nil {
		t.Fatal("expected error for unsupported extension, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported file extension") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestParse_NoExtension(t *testing.T) {
	_, err := Parse("README", strings.NewReader("some text"))
	if err == nil {
		t.Fatal("expected error for missing extension, got nil")
	}
}
