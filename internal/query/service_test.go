package query

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/moulindavid/rag/internal/document"
	"github.com/pgvector/pgvector-go"
)

type mockEmbedder struct {
	vecs [][]float32
	err  error
}

func (m *mockEmbedder) Embed(_ context.Context, _ []string) ([][]float32, error) {
	return m.vecs, m.err
}
func (m *mockEmbedder) Dimension() int { return 3 }

type mockSearcher struct {
	chunks []document.Chunk
	err    error
}

func (m *mockSearcher) SearchSimilar(_ context.Context, _ pgvector.Vector, _ int) ([]document.Chunk, error) {
	return m.chunks, m.err
}

type mockLLM struct {
	answer string
	err    error
}

func (m *mockLLM) Complete(_ context.Context, _, _ string) (string, error) {
	return m.answer, m.err
}

func TestQuery_Success(t *testing.T) {
	chunkID := uuid.New()
	docID := uuid.New()

	svc := NewService(
		&mockSearcher{chunks: []document.Chunk{{ID: chunkID, DocumentID: docID, Content: "relevant content", Score: 0.9}}},
		&mockEmbedder{vecs: [][]float32{{0.1, 0.2, 0.3}}},
		&mockLLM{answer: "the answer"},
	)

	resp, err := svc.Query(context.Background(), &Request{Question: "what?"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Answer != "the answer" {
		t.Errorf("expected answer %q, got %q", "the answer", resp.Answer)
	}
	if len(resp.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(resp.Sources))
	}
	if resp.Sources[0].ChunkID != chunkID {
		t.Errorf("unexpected chunk ID in sources")
	}
}

func TestQuery_EmptySources(t *testing.T) {
	svc := NewService(
		&mockSearcher{chunks: []document.Chunk{}},
		&mockEmbedder{vecs: [][]float32{{0.1}}},
		&mockLLM{answer: "I don't know"},
	)

	resp, err := svc.Query(context.Background(), &Request{Question: "anything?"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Sources) != 0 {
		t.Errorf("expected 0 sources, got %d", len(resp.Sources))
	}
}

func TestQuery_EmbedError(t *testing.T) {
	svc := NewService(
		&mockSearcher{},
		&mockEmbedder{err: errors.New("embed failed")},
		&mockLLM{},
	)
	_, err := svc.Query(context.Background(), &Request{Question: "what?"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestQuery_SearchError(t *testing.T) {
	svc := NewService(
		&mockSearcher{err: errors.New("db unavailable")},
		&mockEmbedder{vecs: [][]float32{{0.1}}},
		&mockLLM{},
	)
	_, err := svc.Query(context.Background(), &Request{Question: "what?"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestQuery_LLMError(t *testing.T) {
	svc := NewService(
		&mockSearcher{chunks: []document.Chunk{}},
		&mockEmbedder{vecs: [][]float32{{0.1}}},
		&mockLLM{err: errors.New("llm timeout")},
	)
	_, err := svc.Query(context.Background(), &Request{Question: "what?"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
