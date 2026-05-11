package embeddings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	ollamaURL string
	model     string
	client    = &http.Client{Timeout: 600 * time.Second}
)

// Init configura o cliente de embeddings.
func Init(url, embedModel string) {
	ollamaURL = url
	model = embedModel
}

type embedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type embedResponse struct {
	Embedding []float64 `json:"embedding"`
}

// Generate gera um embedding para o texto fornecido.
func Generate(text string) ([]float64, error) {
	body, err := json.Marshal(embedRequest{Model: model, Prompt: text})
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(
		ollamaURL+"/api/embeddings",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("ollama embeddings: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama embeddings: status %d", resp.StatusCode)
	}

	var result embedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decodificar embedding: %w", err)
	}
	if len(result.Embedding) == 0 {
		return nil, fmt.Errorf("ollama retornou embedding vazio")
	}
	return result.Embedding, nil
}
