package hand_off

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"strings"
)

// HandoffHeader representa o cabeçalho obrigatório (§4.2.1)
type HandoffHeader struct {
	Timestamp string `json:"timestamp"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	TaskRef   string `json:"task_ref"`
	Intent    string `json:"intent"`
}

// HandoffSchema é a estrutura base genérica para todas as comunicações entre agentes
type HandoffSchema[T any] struct {
	Header  HandoffHeader `json:"header"`
	Payload T             `json:"payload"`
}

// CreateHandoff gera um arquivo JSON padronizado no diretório de handoffs.
// O tipo T deve corresponder a um dos payloads definidos no ARCHITECTURE.md.
func CreateHandoff[T any](dir string, header HandoffHeader, payload T) (string, error) {
	// Garante que o diretório existe
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("falha ao criar diretório de handoff: %w", err)
	}

	// Monta o objeto completo
	handoff := HandoffSchema[T]{
		Header:  header,
		Payload: payload,
	}

	// Se o timestamp estiver vazio, preenche automaticamente
	if handoff.Header.Timestamp == "" {
		handoff.Header.Timestamp = time.Now().Format("2006-01-02 15:04")
	}

	// Define o nome do arquivo: INTENT_SENDER_TASK_TIMESTAMP.json
	fileName := fmt.Sprintf("_%s_TO_%s_%s_%s.json",
		strings.Trim(header.Sender, "[]"),    // Remove os colchetes para o nome do arquivo
		strings.Trim(header.Recipient, "[]"),
		header.Intent,
		time.Now().Format("150405"),
	)
	
	filePath := filepath.Join(dir, fileName)

	// Serialização com indentação para legibilidade (essencial para auditoria humana/IA)
	data, err := json.MarshalIndent(handoff, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro na serialização do handoff: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("erro ao gravar arquivo de handoff: %w", err)
	}

	return filePath, nil
}