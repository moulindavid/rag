package document

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/moulindavid/rag/internal/chunker"
	"github.com/moulindavid/rag/internal/embedding"
	"github.com/moulindavid/rag/internal/parser"
	"github.com/pgvector/pgvector-go"
)

type Service struct {
	repo     *Repository
	embedder embedding.Embedder
}

func NewService(repo *Repository, embedder embedding.Embedder) *Service {
	return &Service{repo: repo, embedder: embedder}
}

func (s *Service) Ingest(ctx context.Context, filename string, content io.Reader) (*Document, error) {

	text, err := parser.Parse(filename, content)
	if err != nil {
		return nil, fmt.Errorf("parse document: %w", err)
	}
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("document is empty")
	}

	chunks := chunker.Chunk(text)
	if len(chunks) == 0 {
		return nil, fmt.Errorf("no chunks produced")
	}

	embeddings, err := s.embedder.Embed(ctx, chunks)
	if err != nil {
		return nil, fmt.Errorf("embed chunks: %w", err)
	}

	doc := Document{
		ID:       uuid.New(),
		Filename: filename,
	}

	chunkRecords := make([]Chunk, len(chunks))
	for i, text := range chunks {
		chunkRecords[i] = Chunk{
			ID:         uuid.New(),
			DocumentID: doc.ID,
			Content:    text,
			Embedding:  pgvector.NewVector(embeddings[i]),
		}
	}

	if err := s.repo.InsertDocument(ctx, &doc); err != nil {
		return nil, fmt.Errorf("ingest: %w", err)
	}
	if err := s.repo.InsertChunks(ctx, chunkRecords); err != nil {
		return nil, fmt.Errorf("ingest: %w", err)
	}

	return &doc, nil
}
