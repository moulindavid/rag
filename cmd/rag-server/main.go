package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/moulindavid/rag/internal/config"
	"github.com/moulindavid/rag/internal/database"
	"github.com/moulindavid/rag/internal/document"
	"github.com/moulindavid/rag/internal/embedding"
	"github.com/moulindavid/rag/internal/llm"
	"github.com/moulindavid/rag/internal/query"
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

	llmClient := llm.NewOllamaClient(cfg.OllamaURL, cfg.OllamaLLMModel)
	slog.Info("llm client created", "model", cfg.OllamaLLMModel)

	repo := document.NewRepository(pool)

	docService := document.NewService(repo, embedder)
	docHandler := document.NewHandler(docService)

	queryService := query.NewService(repo, embedder, llmClient)
	queryHandler := query.NewHandler(queryService)

	r := chi.NewRouter()
	r.Post("/documents", docHandler.Upload)
	r.Post("/query", queryHandler.Query)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	slog.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
