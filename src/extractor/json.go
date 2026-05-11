package extractor

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ExtractJSON extrai texto de um arquivo JSON.
// Ele tenta decodificar o conteúdo e concatenar valores de strings encontrados.
func ExtractJSON(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("ler arquivo json: %w", err)
	}

	// Usamos uma interface genérica para suportar qualquer estrutura JSON
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return "", fmt.Errorf("decodificar json: %w", err)
	}

	var sb strings.Builder
	flattenJSON(&sb, v)

	result := strings.TrimSpace(sb.String())
	if result == "" {
		return "", fmt.Errorf("json sem conteúdo textual extraível: %s", path)
	}

	return result, nil
}

// flattenJSON percorre recursivamente a estrutura do JSON para extrair strings.
func flattenJSON(sb *strings.Builder, v interface{}) {
	switch val := v.(type) {
	case string:
		// Limpa espaços e adiciona ao builder se não for vazio
		s := strings.TrimSpace(val)
		if s != "" {
			sb.WriteString(s)
			sb.WriteString("\n\n")
		}
	case []interface{}:
		// Se for uma lista, processa cada item
		for _, item := range val {
			flattenJSON(sb, item)
		}
	case map[string]interface{}:
		// Se for um objeto, processa cada valor
		for _, value := range val {
			flattenJSON(sb, value)
		}
	}
}