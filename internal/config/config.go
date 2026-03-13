package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL       string
	EmbeddingProvider string
	LLMProvider       string
	OllamaURL         string
	OllamaEmbedModel  string
	OllamaLLMModel    string
	OpenAIAPIKey      string
	ServerPort        string
}

func Load() (*Config, error) {
	databaseURL := getEnvOr("DATABASE_URL", "")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	embeddingProvider := getEnvOr("EMBEDDING_PROVIDER", "ollama")
	llmProvider := getEnvOr("LLM_PROVIDER", "ollama")

	openAIAPIKey := getEnvOr("OPENAI_API_KEY", "")
	if (embeddingProvider == "openai" || llmProvider == "openai") && openAIAPIKey == "" {
		return nil, errors.New("OPENAI_API_KEY is required when using OpenAI as provider")
	}

	return &Config{
		DatabaseURL:       databaseURL,
		EmbeddingProvider: embeddingProvider,
		LLMProvider:       llmProvider,
		OllamaURL:         getEnvOr("OLLAMA_URL", "http://localhost:11434"),
		OllamaEmbedModel:  getEnvOr("OLLAMA_EMBED_MODEL", "nomic-embed-text"),
		OllamaLLMModel:    getEnvOr("OLLAMA_LLM_MODEL", "llama3"),
		OpenAIAPIKey:      openAIAPIKey,
		ServerPort:        getEnvOr("SERVER_PORT", "8080"),
	}, nil
}

func getEnvOr(key, defaultValue string) string {
	var envValue = os.Getenv(key)
	if envValue == "" {
		return defaultValue
	}
	return envValue
}
