package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
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
	loadEnvFile(".env")

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

// loadEnvFile reads a .env file and sets environment variables.
func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

func getEnvOr(key, defaultValue string) string {
	var envValue = os.Getenv(key)
	if envValue == "" {
		return defaultValue
	}
	return envValue
}
