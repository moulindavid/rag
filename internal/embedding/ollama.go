package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type OllamaEmbedder struct {
	url   string
	model string
}

type embedRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

var _ Embedder = (*OllamaEmbedder)(nil)

func NewOllamaEmbedder(url, model string) *OllamaEmbedder {
	return &OllamaEmbedder{
		url:   url,
		model: model,
	}
}

func (o *OllamaEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	reqBody, err := json.Marshal(embedRequest{Model: o.model, Input: texts})
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.url+"/api/embed", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Embeddings [][]float32 `json:"embeddings"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return result.Embeddings, nil
}

func (o *OllamaEmbedder) Dimension() int {
	return 768
}
