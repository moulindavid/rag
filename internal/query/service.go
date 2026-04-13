package query

import (
	"context"
	"fmt"
	"strings"

	"github.com/moulindavid/rag/internal/document"
	"github.com/moulindavid/rag/internal/embedding"
	"github.com/moulindavid/rag/internal/llm"
	"github.com/pgvector/pgvector-go"
)

type chunkSearcher interface {
	SearchSimilar(ctx context.Context, embedding pgvector.Vector, limit int) ([]document.Chunk, error)
}

type Service struct {
	repo     chunkSearcher
	embedder embedding.Embedder
	llm      llm.Client
}

func NewService(repo chunkSearcher, embedder embedding.Embedder, llmClient llm.Client) *Service {
	return &Service{repo: repo, embedder: embedder, llm: llmClient}
}

func (s *Service) Query(ctx context.Context, req *Request) (*Response, error) {
	vecs, err := s.embedder.Embed(ctx, []string{req.Question})
	if err != nil {
		return nil, fmt.Errorf("embed question: %w", err)
	}
	vec := pgvector.NewVector(vecs[0])
	chunks, err := s.repo.SearchSimilar(ctx, vec, 5)
	if err != nil {
		return nil, fmt.Errorf("search similar: %w", err)
	}

	var sb strings.Builder
	for i, c := range chunks {
		fmt.Fprintf(&sb, "[%d] %s\n\n", i+1, c.Content)
	}

	system := "You are a helpful assistant. Answer the user's question using only the context below. If the answer is not in the context, say so.\n\nContext:\n" + sb.String()

	answer, err := s.llm.Complete(ctx, system, req.Question)
	if err != nil {
		return nil, fmt.Errorf("llm complete: %w", err)
	}

	sources := make([]Source, len(chunks))
	for i, c := range chunks {
		sources[i] = Source{
			ChunkID:    c.ID,
			DocumentID: c.DocumentID,
			Content:    c.Content,
			Score:      c.Score,
		}
	}

	return &Response{Answer: answer, Sources: sources}, nil
}
