package document

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) InsertDocument(ctx context.Context, doc *Document) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO documents (id, filename) VALUES ($1, $2)`,
		doc.ID, doc.Filename,
	)
	if err != nil {
		return fmt.Errorf("insert document: %w", err)
	}
	return nil
}

func (r *Repository) InsertChunks(ctx context.Context, chunks []Chunk) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	_, err = conn.Conn().CopyFrom(
		ctx,
		pgx.Identifier{"chunks"},
		[]string{"id", "document_id", "content", "embedding"},
		pgx.CopyFromSlice(len(chunks), func(i int) ([]any, error) {
			return []any{
				chunks[i].ID,
				chunks[i].DocumentID,
				chunks[i].Content,
				chunks[i].Embedding,
			}, nil
		}),
	)
	if err != nil {
		return fmt.Errorf("copy chunks: %w", err)
	}
	return nil
}

func (r *Repository) SearchSimilar(ctx context.Context, embedding pgvector.Vector, limit int) ([]Chunk, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, document_id, content, 1 - (embedding <=> $1) AS score
		 FROM chunks
		 ORDER BY embedding <=> $1
		 LIMIT $2`,
		embedding, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("search similar: %w", err)
	}
	defer rows.Close()

	var chunks []Chunk
	for rows.Next() {
		var c Chunk
		if err := rows.Scan(&c.ID, &c.DocumentID, &c.Content); err != nil {
			return nil, fmt.Errorf("scan chunk: %w", err)
		}
		chunks = append(chunks, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return chunks, nil
}
