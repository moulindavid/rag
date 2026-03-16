package document

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type Document struct {
	ID        uuid.UUID
	Filename  string
	CreatedAt time.Time
}

type Chunk struct {
	ID         uuid.UUID
	DocumentID uuid.UUID
	Content    string
	Embedding  pgvector.Vector
	CreatedAt  time.Time
	Score      float64
}
