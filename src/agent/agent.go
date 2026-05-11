package agent

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ollamaURL    string
	model        string
	systemPrompt string // Carregado dinamicamente de AGEND.md
	client       = &http.Client{Timeout: 60 * time.Minute}
)

// Init configura o cliente do agente LLM e carrega as definições do AGEND.md.
func Init(url, llmModel string) {
	ollamaURL = url
	model = llmModel

	// Carrega o System Prompt do arquivo local para permitir ajustes sem recompilação
	p, err := CarregarVariavel("SYSTEM_PROMPT")
	if err != nil {
		fmt.Printf("⚠️  Aviso: SYSTEM_PROMPT não encontrado em AGENTS. Usando configuração padrão.\n")
		systemPrompt = "Você é um assistente útil e fiel aos dados fornecidos."
	} else {
		systemPrompt = p
	}
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type streamChunk struct {
	Message message `json:"message"`
	Done    bool    `json:"done"`
}

// GerarArquivo transforma chunks de contexto em uma nota estruturada para Obsidian.
func GerarArquivo(tema string, chunks []string) (string, error) {
	if len(chunks) == 0 {
		return "", fmt.Errorf("nenhum chunk de contexto fornecido")
	}

	contexto := strings.Join(chunks, "\n\n---\n\n")
	prompt := fmt.Sprintf(
		"Crie uma nota Obsidian completa sobre o tema: **%s**\n\nContexto extraído do documento:\n\n%s",
		tema, contexto,
	)

	// Salva o prompt exato para auditoria (evita alucinações como 'Dr. Amorim')
	
	body, err := json.Marshal(chatRequest{
		Model: model,
		Messages: []message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: prompt},
		},
		Stream: true,
	})
	if err != nil {
		return "", err
	}

	resp, err := client.Post(
		ollamaURL+"/api/chat",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("ollama connection: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama api error: status %d", resp.StatusCode)
	}

	var sb strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 512*1024)

	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	tokenCount := 0
	inicio := time.Now()

	fmt.Printf("  ⠋ [%s] Aguardando resposta do modelo...", tema)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var chunk streamChunk
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			continue
		}

		sb.WriteString(chunk.Message.Content)
		tokenCount++

		if tokenCount%15 == 0 {
			spin := spinner[(tokenCount/15)%len(spinner)]
			fmt.Printf("\r  %s [%s] Processando fluxo... %d tokens", spin, tema, tokenCount)
		}

		if chunk.Done {
			break
		}
	}

	elapsed := time.Since(inicio).Round(time.Second)
	fmt.Printf("\r  ✅ [%s] Nota concluída — %d tokens em %s          \n", tema, tokenCount, elapsed)

	return strings.TrimSpace(sb.String()), nil
}

// salvarPromptLocal gera um log do prompt em /prompts com timestamp para depuração.
func salvarPromptLocal(tema string, conteudo string) error {
	dir := "prompts"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	
	// Sanitização básica do nome do arquivo
	temaLimpo := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, tema)

	fileName := fmt.Sprintf("%s_%s.md", timestamp, temaLimpo)
	return os.WriteFile(filepath.Join(dir, fileName), []byte(conteudo), 0644)
}

// CarregarVariavel extrai configurações de AGENTS/*.md, suportando blocos de texto.
func CarregarVariavel(chave string) (string, error) {
	file, err := os.ReadFile("AGENTS/" + chave + ".md")
	if err != nil {
		return "", err
	}

	linhas := strings.Split(string(file), "\n")
	var capturando bool
	var resultado []string

	for _, linha := range linhas {
		linhaTrim := strings.TrimSpace(linha)
		
		// Inicia captura ao encontrar a Chave:
		if strings.HasPrefix(strings.ToUpper(linhaTrim), strings.ToUpper(chave)+":") {
			capturando = true
			conteudoMesmaLinha := strings.TrimSpace(strings.TrimPrefix(linhaTrim, chave+":"))
			if conteudoMesmaLinha != "" {
				resultado = append(resultado, conteudoMesmaLinha)
			}
			continue
		}

		// Para a captura ao encontrar outra definição de variável (Padrão Chave:)
		if capturando && strings.Contains(linhaTrim, ":") && !strings.HasPrefix(linha, " ") && !strings.HasPrefix(linha, "\t") {
			parts := strings.Split(linhaTrim, ":")
			if len(parts[0]) < 25 && !strings.Contains(parts[0], " ") {
				break
			}
		}

		if capturando {
			resultado = append(resultado, linha)
		}
	}

	if len(resultado) == 0 {
		return "", fmt.Errorf("variável %s não encontrada em AGEND.md", chave)
	}

	return strings.TrimSpace(strings.Join(resultado, "\n")), nil
}

// ────────────────────────────────────────────────────────────────────
// Resolução de Modelos por Agente (baseado em ENV_SETUP.md)
// ────────────────────────────────────────────────────────────────────

// ResolverModeloAgente lê ENV_SETUP.md e retorna o modelo para o agente especificado.
// Mapeia agentes a seus respectivos modelos conforme seção 3.1.
func ResolverModeloAgente(nomeAgente string) (string, error) {
	// Mapeamento de agentes a modelos/tiers
	// Baseado em ENV_SETUP.md seção 3.1 AI Agents Configuration
	agentModelMap := map[string]string{
		"CEO":           "deepseek-r1:7b",       // TIER_3_EXPERT_MODEL
		"BA":            "qwen2.5-coder",        // TIER_2_DEVELOPMENT_MODEL
		"CTO":           "deepseek-r1:7b",       // TIER_3_EXPERT_MODEL
		"DEV_FRONTEND":  "qwen2.5-coder",        // TIER_2_DEVELOPMENT_MODEL
		"DEV_BACKEND":   "qwen2.5-coder",        // TIER_2_DEVELOPMENT_MODEL
		"DBA":           "qwen2.5-coder",        // TIER_2_DEVELOPMENT_MODEL
		"DS_ML":         "qwen2.5-coder",        // TIER_2_DEVELOPMENT_MODEL
		"SECURITY":      "deepseek-r1:7b",       // TIER_3_EXPERT_MODEL
		"QA":            "gemma4:latest",        // TIER_1_EFFICIENCY_MODEL
		"DATA_ENGINEER": "qwen2.5-coder",        // TIER_2_DEVELOPMENT_MODEL
		"PM":            "deepseek-r1:7b",       // TIER_3_EXPERT_MODEL
		"UX_RESEARCHER": "qwen2.5-coder",        // TIER_2_DEVELOPMENT_MODEL
		"WRITER":        "gemma4:latest",        // TIER_1_EFFICIENCY_MODEL
		"DOCUMENTATION":"gemma4:latest",         // TIER_1_EFFICIENCY_MODEL
		"ARTIST":        "gemma4:latest",        // TIER_1_EFFICIENCY_MODEL
		"DEVOPS":        "qwen2.5-coder",        // TIER_2_DEVELOPMENT_MODEL
		"CODE_REVIEWER": "deepseek-r1:7b",       // TIER_3_EXPERT_MODEL
		"CMO":           "gemma4:latest",        // TIER_1_EFFICIENCY_MODEL
	}

	// Normalizar nome do agente (uppercase)
	agentUpper := strings.ToUpper(nomeAgente)

	if m, exists := agentModelMap[agentUpper]; exists {
		return m, nil
	}

	// Se não encontrar no mapa, retorna erro
	return "", fmt.Errorf("agente %s não mapeado em ENV_SETUP.md", nomeAgente)
}

// CarregarConfigAgente carrega modelo e system prompt para um agente específico.
// Tenta ler de AGENTS/{nomeAgente}.md primeiro, depois aplica o modelo do mapeamento.
func CarregarConfigAgente(nomeAgente string) (modeloResolvido, systemPromptResolvido string, err error) {
	// Resolver modelo do agente
	modeloResolvido, err = ResolverModeloAgente(nomeAgente)
	if err != nil {
		return "", "", err
	}

	// Tentar carregar system prompt do arquivo de configuração do agente
	systemPromptResolvido, err = CarregarVariavel(nomeAgente + ":SYSTEM_PROMPT")
	if err != nil {
		// Se não encontrar, usar um padrão genérico
		systemPromptResolvido = fmt.Sprintf(
			"Você é o agente %s. Responda com precisão e clareza, fornecendo respostas estruturadas.",
			nomeAgente,
		)
	}

	return modeloResolvido, systemPromptResolvido, nil
}

// InitComAgente configura o cliente LLM específico para um agente.
// Carrega modelo e system prompt baseado no arquivo ENV_SETUP.md.
func InitComAgente(url, nomeAgente string) error {
	ollamaURL = url

	modelResolvido, promptResolvido, err := CarregarConfigAgente(nomeAgente)
	if err != nil {
		return fmt.Errorf("erro ao carregar config do agente %s: %w", nomeAgente, err)
	}

	model = modelResolvido
	systemPrompt = promptResolvido

	fmt.Printf("✅ Agente %s carregado: modelo=%s\n", nomeAgente, model)
	return nil
}