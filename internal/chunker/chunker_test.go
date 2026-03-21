package chunker

import (
	"strings"
	"testing"
)

func TestChunk_EmptyText(t *testing.T) {
	result := Chunk("")
	if len(result) != 0 {
		t.Errorf("expected 0 chunks, got %d", len(result))
	}
}

func TestChunk_SingleParagraph(t *testing.T) {
	result := Chunk("Hello world")
	if len(result) != 1 {
		t.Errorf("expected 1 chunk, got %d", len(result))
	}
	if result[0] != "Hello world" {
		t.Errorf("unexpected chunk content: %q", result[0])
	}
}

func TestChunk_SmallParagraphsMerged(t *testing.T) {
	text := "First paragraph.\n\nSecond paragraph.\n\nThird paragraph."
	result := Chunk(text)
	if len(result) != 1 {
		t.Errorf("expected 1 chunk (all fit under 2000 chars), got %d", len(result))
	}
}

func TestChunk_LargeTextSplitsIntoMultipleChunks(t *testing.T) {
	// Build a text with paragraphs of ~300 chars each
	// 7 paragraphs x 300 chars = 2100 chars total -> should produce at least 2 chunks
	para := strings.Repeat("a", 300)
	text := strings.Join([]string{para, para, para, para, para, para, para}, "\n\n")

	result := Chunk(text)
	if len(result) < 2 {
		t.Errorf("expected at least 2 chunks, got %d", len(result))
	}
}

func TestChunk_OverlapAppearsInNextChunk(t *testing.T) {
	// First chunk will be ~2000 chars, second chunk should start with the last 200 chars of it
	para := strings.Repeat("a", 600)
	text := strings.Join([]string{para, para, para, para}, "\n\n")

	result := Chunk(text)
	if len(result) < 2 {
		t.Fatalf("expected at least 2 chunks, got %d", len(result))
	}

	firstChunk := result[0]
	secondChunk := result[1]
	expectedOverlap := firstChunk[len(firstChunk)-overlapSize:]

	if !strings.Contains(secondChunk, expectedOverlap) {
		t.Errorf("expected second chunk to start with overlap from first chunk")
	}
}

func TestChunk_SkipsEmptyParagraphs(t *testing.T) {
	text := "First.\n\n\n\nSecond."
	result := Chunk(text)
	// \n\n\n\n splits into ["First.", "", "Second."] — empty one should be skipped
	if len(result) != 1 {
		t.Errorf("expected 1 chunk, got %d", len(result))
	}
}
