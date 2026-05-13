package agent

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kballard/go-shellquote"
)

var (
	mu           sync.RWMutex
	ollamaURL    string
	model        string
	systemPrompt string // Carregado dinamicamente de AGENTS/
	agentsDir    = "AGENTS"
	client       = &http.Client{Timeout: 60 * time.Minute}
)

func SetConfigDir(dir string) {
	mu.Lock()
	defer mu.Unlock()
	agentsDir = dir
}

// Init configura o cliente do agente LLM e carrega as definições de AGENTS/.
func Init(url, llmModel string) {
	mu.Lock()
	defer mu.Unlock()

	ollamaURL = url
	model = llmModel
	systemPrompt = "Você é um assistente útil e fiel aos dados fornecidos."
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
		"Tema: %s\n\nContexto extraído do documento:\n\n%s",
		tema, contexto,
	)

	mu.RLock()
	currentModel := model
	currentURL := ollamaURL
	currentPrompt := systemPrompt
	mu.RUnlock()

	body, err := json.Marshal(chatRequest{
		Model: currentModel,
		Messages: []message{
			{Role: "system", Content: currentPrompt},
			{Role: "user", Content: prompt},
		},
		Stream: true,
	})
	if err != nil {
		return "", err
	}

	resp, err := client.Post(
		currentURL+"/api/chat",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("ollama connection: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		detail, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("ollama api error: status %d model=%s body=%s", resp.StatusCode, currentModel, strings.TrimSpace(string(detail)))
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

// salvarPromptLocal gera um log do prompt em prompts/ com timestamp para depuração.
func ChamarAgente(nomeAgente, handoff string) (string, error) {
	modeloResolvido, promptResolvido, err := CarregarConfigAgente(nomeAgente)
	if err != nil {
		return "", fmt.Errorf("carregar agente %s: %w", nomeAgente, err)
	}

	mu.RLock()
	currentURL := ollamaURL
	mu.RUnlock()
	if currentURL == "" {
		return "", fmt.Errorf("ollama url não inicializada")
	}

	userPrompt := fmt.Sprintf(`Você recebeu um handoff para execução.

Leia o handoff e responda exclusivamente com um JSON válido.

Quando precisar criar ou atualizar arquivos de código, use este formato:
{
  "action": "write_code",
  "files": [
    {"path": "caminho/relativo/ao/workspace.ext", "content": "conteúdo completo do arquivo"}
  ]
}

O content de write_code deve ser codigo puro, sem markdown, sem crases triplas e sem explicacao.

Quando precisar chamar outro agente, use um handoff estruturado com "header" e "payload" ou:
{
  "action": "handoff",
  "handoffs": [
    {"header": {"sender": "[%s]", "recipient": "[AGENTE]", "task_ref": "TASK-XXX", "intent": "..."}, "payload": {}}
  ]
}

Se não houver arquivo ou handoff a gerar, use:
{"action":"note","message":"resumo objetivo da execução"}

Quando a resposta for atualização operacional, não escreva DOC/*.md. Use uma ação de banco:
{"action":"update_context","document_type":"CONTEXT","section":"...","title":"...","content":"...","tags":["..."]}
{"action":"update_plan","item_type":"plan","reference":"DOC/PLAN.md#...","title":"...","status":"...","priority":"...","content":"..."}
{"action":"record_test","test_name":"...","status":"PASSED|FAILED|SKIPPED","command":"...","output":"..."}
{"action":"ask_human","questions":[{"question":"...","priority":"High|Medium|Low","blocking":true}]}

Para wiki/Obsidian, crie arquivos markdown reais em wiki/ usando write_code.

Handoff recebido:
%s`, strings.ToUpper(nomeAgente), handoff)

	return callChat(currentURL, modeloResolvido, promptResolvido, userPrompt)
}

func ChamarCEOComCodexCLI(handoff, cliCommand, workspaceRoot string, timeout time.Duration) (string, error) {
	return ChamarAgenteComCodexCLI("CEO", handoff, cliCommand, workspaceRoot, timeout)
}

func ChamarAgenteComCodexCLI(nomeAgente, handoff, cliCommand, workspaceRoot string, timeout time.Duration) (string, error) {
	agentName := strings.ToUpper(strings.TrimSpace(nomeAgente))
	if agentName == "" {
		return "", fmt.Errorf("nome do agente vazio para Codex CLI")
	}
	agentModel, err := ResolverModeloAgente(agentName)
	if err != nil {
		return "", err
	}
	cliCommand = expandCLICommand(cliCommand, agentModel)
	if strings.TrimSpace(cliCommand) == "" {
		return "", fmt.Errorf("comando Codex CLI vazio para agente %s", agentName)
	}
	if timeout <= 0 {
		timeout = 10 * time.Minute
	}

	agentPrompt, err := carregarArquivoAgente(agentName)
	if err != nil {
		return "", fmt.Errorf("carregar agente %s para Codex CLI: %w", agentName, err)
	}

	prompt := fmt.Sprintf(`Voce esta atuando como o agente %s deste orquestrador local.
Modelo configurado: %s

INSTRUCOES DO AGENTE %s:
%s

CONTRATO DE RESPOSTA:
- Responda exclusivamente com JSON valido.
- Para chamar outro agente, use {"action":"handoff","handoffs":[{"header":{...},"payload":{...}}]}.
- Para atualizar contexto/plano/estado/testes, use acoes CQRS: update_context, update_plan, update_state, record_test, record_decision ou record_retrospective.
- Para perguntas ao humano, gere uma acao ask_human com questions. O sistema criara um handoff _TO_HUMAN que nao passa pela pipeline ate o humano responder e renomear para o agente destinatario.
- Para wiki/Obsidian ou codigo, use write_code com paths relativos ao workspace.
- Para codigo, files[].content deve conter apenas o conteudo bruto do arquivo, sem markdown, sem crases triplas e sem explicacao.

HANDOFF RECEBIDO:
%s`, agentName, agentModel, agentName, agentPrompt, handoff)

	args, err := shellquote.Split(cliCommand)
	if err != nil {
		return "", fmt.Errorf("parsear comando Codex CLI do agente %s: %w", agentName, err)
	}
	if len(args) == 0 {
		return "", fmt.Errorf("comando Codex CLI invalido para agente %s", agentName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	if strings.TrimSpace(workspaceRoot) != "" {
		cmd.Dir = workspaceRoot
	}
	cmd.Stdin = strings.NewReader(prompt)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("codex cli timeout para agente %s apos %s", agentName, timeout)
		}
		return "", fmt.Errorf("codex cli falhou para agente %s: %w: %s", agentName, err, strings.TrimSpace(stderr.String()))
	}

	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return "", fmt.Errorf("codex cli retornou resposta vazia para agente %s: %s", agentName, strings.TrimSpace(stderr.String()))
	}
	return output, nil
}

func expandCLICommand(cliCommand, modelName string) string {
	if strings.Contains(cliCommand, "{{MODEL}}") {
		return strings.ReplaceAll(cliCommand, "{{MODEL}}", modelName)
	}
	return cliCommand
}

func callChat(url, modelName, prompt, userPrompt string) (string, error) {
	body, err := json.Marshal(chatRequest{
		Model: modelName,
		Messages: []message{
			{Role: "system", Content: prompt},
			{Role: "user", Content: userPrompt},
		},
		Stream: true,
	})
	if err != nil {
		return "", err
	}

	resp, err := client.Post(
		url+"/api/chat",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("ollama connection: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		detail, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("ollama api error: status %d model=%s body=%s", resp.StatusCode, modelName, strings.TrimSpace(string(detail)))
	}

	var sb strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 512*1024)

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
		if chunk.Done {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("ler resposta do modelo: %w", err)
	}

	return strings.TrimSpace(sb.String()), nil
}

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

// ────────────────────────────────────────────────────────────────────
// Leitura de configuração de agentes
// ────────────────────────────────────────────────────────────────────

// carregarVariavelDeAgente lê AGENTS/{nomeArquivo}.md e extrai o valor da chave informada.
// nomeArquivo e chave são separados para evitar nomes de arquivo com caracteres inválidos.
func carregarVariavelDeAgente(nomeArquivo, chave string) (string, error) {
	mu.RLock()
	dir := agentsDir
	mu.RUnlock()

	filePath := filepath.Join(dir, nomeArquivo+".md")
	file, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	linhas := strings.Split(string(file), "\n")
	var capturando bool
	var resultado []string

	upperChave := strings.ToUpper(chave) + ":"

	for _, linha := range linhas {
		linhaTrim := strings.TrimSpace(linha)

		// Inicia captura ao encontrar a Chave: (case-insensitive na detecção,
		// mas usa o índice do prefixo para cortar — evita TrimPrefix com case errado)
		if strings.HasPrefix(strings.ToUpper(linhaTrim), upperChave) {
			capturando = true
			conteudoMesmaLinha := strings.TrimSpace(linhaTrim[len(upperChave):])
			if conteudoMesmaLinha != "" {
				resultado = append(resultado, conteudoMesmaLinha)
			}
			continue
		}

		// Para a captura ao encontrar outra definição de variável (padrão Chave:)
		if capturando && strings.Contains(linhaTrim, ":") &&
			!strings.HasPrefix(linha, " ") && !strings.HasPrefix(linha, "\t") {
			parts := strings.SplitN(linhaTrim, ":", 2)
			if len(parts[0]) < 25 && !strings.Contains(parts[0], " ") {
				break
			}
		}

		if capturando {
			resultado = append(resultado, linha)
		}
	}

	if len(resultado) == 0 {
		return "", fmt.Errorf("chave %q não encontrada em AGENTS/%s.md", chave, nomeArquivo)
	}

	return strings.TrimSpace(strings.Join(resultado, "\n")), nil
}

func carregarArquivoAgente(nomeAgente string) (string, error) {
	mu.RLock()
	dir := agentsDir
	mu.RUnlock()

	filePath := filepath.Join(dir, strings.ToUpper(nomeAgente)+".md")
	file, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	content := strings.TrimSpace(string(file))
	if content == "" {
		return "", fmt.Errorf("AGENTS/%s.md vazio", nomeAgente)
	}
	return content, nil
}

// CarregarVariavel mantém compatibilidade com chamadas existentes.
// Usa o próprio nome da chave como nome do arquivo (ex: "SYSTEM_PROMPT" → AGENTS/SYSTEM_PROMPT.md).
func CarregarVariavel(chave string) (string, error) {
	return carregarVariavelDeAgente(chave, chave)
}

// ────────────────────────────────────────────────────────────────────
// Resolução de Modelos por Agente (baseado em ENV_SETUP.md)
// ────────────────────────────────────────────────────────────────────

// ResolverModeloAgente retorna o modelo LLM para o agente especificado.
// Mapeia agentes a seus respectivos modelos conforme ENV_SETUP.md seção 3.1.
func ResolverModeloAgente(nomeAgente string) (string, error) {
	agentModelMap := map[string]string{
		"CEO":           "deepseek-r1:7b", // TIER_3_EXPERT_MODEL
		"BA":            "qwen2.5-coder",  // TIER_2_DEVELOPMENT_MODEL
		"CTO":           "deepseek-r1:7b", // TIER_3_EXPERT_MODEL
		"DEV_FRONTEND":  "qwen2.5-coder",  // TIER_2_DEVELOPMENT_MODEL
		"DEV_BACKEND":   "qwen2.5-coder",  // TIER_2_DEVELOPMENT_MODEL
		"DBA":           "qwen2.5-coder",  // TIER_2_DEVELOPMENT_MODEL
		"DS_ML":         "qwen2.5-coder",  // TIER_2_DEVELOPMENT_MODEL
		"SECURITY":      "deepseek-r1:7b", // TIER_3_EXPERT_MODEL
		"QA":            "gemma4:latest",  // TIER_1_EFFICIENCY_MODEL
		"DATA_ENGINEER": "qwen2.5-coder",  // TIER_2_DEVELOPMENT_MODEL
		"PM":            "deepseek-r1:7b", // TIER_3_EXPERT_MODEL
		"UX_RESEARCHER": "qwen2.5-coder",  // TIER_2_DEVELOPMENT_MODEL
		"WRITER":        "gemma4:latest",  // TIER_1_EFFICIENCY_MODEL
		"DOCUMENTATION": "gemma4:latest",  // TIER_1_EFFICIENCY_MODEL
		"ARTIST":        "gemma4:latest",  // TIER_1_EFFICIENCY_MODEL
		"DEVOPS":        "qwen2.5-coder",  // TIER_2_DEVELOPMENT_MODEL
		"CODE_REVIEWER": "deepseek-r1:7b", // TIER_3_EXPERT_MODEL
		"CMO":           "gemma4:latest",  // TIER_1_EFFICIENCY_MODEL
	}

	agentUpper := strings.ToUpper(nomeAgente)
	if modelFromEnv := strings.TrimSpace(os.Getenv(agentEnvKey(agentUpper, "MODEL"))); modelFromEnv != "" {
		return modelFromEnv, nil
	}
	if m, exists := agentModelMap[agentUpper]; exists {
		return m, nil
	}
	return "", fmt.Errorf("agente %q não mapeado em ENV_SETUP.md", nomeAgente)
}

// ResolverProviderAgente retorna "cli" ou "ollama" para o agente informado.
// A configuracao vem de <AGENTE>_AGENT_PROVIDER; por padrao todos usam Ollama.
func ResolverProviderAgente(nomeAgente string) string {
	agentUpper := strings.ToUpper(strings.TrimSpace(nomeAgente))
	if provider := strings.TrimSpace(os.Getenv(agentEnvKey(agentUpper, "PROVIDER"))); provider != "" {
		return strings.ToLower(provider)
	}
	return "ollama"
}

func agentEnvKey(agentUpper, suffix string) string {
	envNames := map[string]string{
		"CEO":           "CEO",
		"BA":            "BA",
		"CTO":           "CTO",
		"DEV_FRONTEND":  "DEV_FRONT",
		"DEV_BACKEND":   "DEV_BACK",
		"DBA":           "DBA",
		"DS_ML":         "DS_ML",
		"SECURITY":      "SECURITY",
		"QA":            "QA",
		"DATA_ENGINEER": "DATA_ENGINEER",
		"PM":            "PM",
		"UX_RESEARCHER": "UX_RESEARCHER",
		"WRITER":        "WRITER",
		"DOCUMENTATION": "DOCUMENTATION",
		"ARTIST":        "ARTIST",
		"DEVOPS":        "DEVOPS",
		"CODE_REVIEWER": "CODE_REVIEWER",
		"CMO":           "CMO",
	}
	prefix, ok := envNames[agentUpper]
	if !ok {
		prefix = strings.ReplaceAll(agentUpper, "-", "_")
	}
	return prefix + "_AGENT_" + suffix
}

// CarregarConfigAgente carrega modelo e instrucoes para um agente.
// Primeiro tenta extrair SYSTEM_PROMPT; se nao houver chave, usa o arquivo inteiro.
func CarregarConfigAgente(nomeAgente string) (modeloResolvido, systemPromptResolvido string, err error) {
	modeloResolvido, err = ResolverModeloAgente(nomeAgente)
	if err != nil {
		return "", "", err
	}

	// Lê AGENTS/{AGENTE}.md e extrai a chave SYSTEM_PROMPT
	systemPromptResolvido, err = carregarVariavelDeAgente(nomeAgente, "SYSTEM_PROMPT")
	if err != nil {
		systemPromptResolvido, err = carregarArquivoAgente(nomeAgente)
		if err != nil {
			return "", "", err
		}
	}

	return modeloResolvido, systemPromptResolvido, nil
}

// InitComAgente configura o cliente LLM específico para um agente.
// Carrega modelo e system prompt baseado em AGENTS/{nomeAgente}.md.
func InitComAgente(url, nomeAgente string) error {
	modelResolvido, promptResolvido, err := CarregarConfigAgente(nomeAgente)
	if err != nil {
		return fmt.Errorf("erro ao carregar config do agente %s: %w", nomeAgente, err)
	}

	mu.Lock()
	defer mu.Unlock()

	ollamaURL = url
	model = modelResolvido
	systemPrompt = promptResolvido

	fmt.Printf("✅ Agente %s carregado: modelo=%s\n", nomeAgente, model)
	return nil
}
