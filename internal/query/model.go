package query

import "github.com/google/uuid"

type Request struct {
	Question string `json:"question"`
}

type Source struct {
	ChunkID    uuid.UUID `json:"chunk_id"`
	DocumentID uuid.UUID `json:"document_id"`
	Content    string    `json:"content"`
	Score      float64   `json:"score"`
}

type Response struct {
	Answer  string   `json:"answer"`
	Sources []Source `json:"sources"`
}
