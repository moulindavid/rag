package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/moulindavid/rag/internal/config"
	"github.com/moulindavid/rag/internal/database"
	"github.com/moulindavid/rag/internal/embedding"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	slog.Info("config loaded",
		"port", cfg.ServerPort,
		"embedding_provider", cfg.EmbeddingProvider,
		"ollama_url", cfg.OllamaURL,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	slog.Info("database connected")

	err = database.Migrate(ctx, pool)
	if err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("migrations applied")

	embedder := embedding.NewOllamaEmbedder(cfg.OllamaURL, cfg.OllamaEmbedModel)
	slog.Info("embedder created", "model", cfg.OllamaEmbedModel, "dimension", embedder.Dimension())

	vectors, err := embedder.Embed(ctx, []string{"Bip bop bip", "Ghosty ty ghost"})
	if err != nil {
		slog.Error("failed to embed texts", "error", err)
		os.Exit(1)
	}
	slog.Info("embedding test passed",
		"num_vectors", len(vectors),
		"dim_vector_0", len(vectors[0]),
		"first_3_values", fmt.Sprintf("%.4f, %.4f, %.4f", vectors[0][0], vectors[0][1], vectors[0][2]),
	)

	slog.Info("Bip bop no errorino")
}
