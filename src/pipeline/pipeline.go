package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"pf_ai/agent"
	"pf_ai/code_writer"
	"pf_ai/context_store"
	"pf_ai/embeddings"
	"pf_ai/extractor"
	"pf_ai/hand_off"
)

// Source indica a origem do job.
type Source int

const (
	SourceRaw       Source = iota // arquivo novo em inbox/raw/
	SourceExtracted               // .chunks.json pronto do Python
	SourceHandoff                 // handoff para um agente
)

// Job representa um arquivo a processar.
type Job struct {
	Path      string
	Source    Source
	AgentName string
}

// Config agrupa parâmetros do pipeline.
type Config struct {
	WorkerCount      int
	EmbedWorkerCount int
	PendingPython    string
	Extracted        string
	Processing       string
	ProcessedPython  string
	Success          string
	Failed           string
	PythonTimeout    time.Duration
	HandoffDir       string
	AgentOutputDir   string
	WorkspaceRoot    string
	ContextStore     *context_store.Store
	CodexCLI         string
	CodexTimeout     time.Duration
}

// Pipeline gerencia o worker pool.
type Pipeline struct {
	cfg       Config
	jobs      chan Job
	wg        sync.WaitGroup
	fileMu    sync.Mutex
	pendingMu sync.Mutex
	pending   map[string]chan string // id → path do .chunks.json
}

// New cria um Pipeline com a configuração fornecida.
func New(cfg Config) *Pipeline {
	return &Pipeline{
		cfg:     cfg,
		jobs:    make(chan Job, 50),
		pending: make(map[string]chan string),
	}
}

// Start inicializa os workers. Retorna imediatamente.
func (p *Pipeline) Start(ctx context.Context) {
	for i := range p.cfg.WorkerCount {
		p.wg.Add(1)
		go func(id int) {
			defer p.wg.Done()
			log.Printf("🔧 Worker %d iniciado", id)
			for {
				select {
				case job, ok := <-p.jobs:
					if !ok {
						return
					}
					p.handle(ctx, job)
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}
}

// Submit envia um job para o pool. Não bloqueia.
func (p *Pipeline) Submit(job Job) {
	select {
	case p.jobs <- job:
	default:
		log.Printf("⚠️  Fila cheia — descartando: %s", filepath.Base(job.Path))
	}
}

// Wait aguarda todos os workers encerrarem.
func (p *Pipeline) Wait() {
	close(p.jobs)
	p.wg.Wait()
}

// ── Roteamento ────────────────────────────────────────────────────────────────

func (p *Pipeline) handle(ctx context.Context, job Job) {
	switch job.Source {
	case SourceRaw:
		p.handleRaw(ctx, job.Path)
	case SourceExtracted:
		p.handleExtracted(job.Path)
	case SourceHandoff:
		p.handleHandoff(ctx, job.Path, job.AgentName)
	}
}

// handleExtracted notifica jobs aguardando pelo chunks.json do Python.
func (p *Pipeline) handleExtracted(path string) {
	base := filepath.Base(path)
	if !strings.HasSuffix(base, ".chunks.json") {
		return
	}
	id := strings.TrimSuffix(base, ".chunks.json")

	p.pendingMu.Lock()
	ch, ok := p.pending[id]
	p.pendingMu.Unlock()

	if ok {
		ch <- path
	}
	// Se não há job aguardando, o fsnotify disparou antes do registro —
	// o job vai encontrar o arquivo ao verificar no disco.
}

func (p *Pipeline) handleHandoff(ctx context.Context, path, agentName string) {
	if agentName == "" {
		log.Printf("Handoff sem agente: %s", filepath.Base(path))
		return
	}

	id := fileID(path)
	log.Printf("[%s] Chamando agente %s", id, agentName)

	procPath, err := p.moveFile(path, p.cfg.Processing)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Printf("Mover handoff para processing: %v", err)
		return
	}

	data, err := os.ReadFile(procPath)
	if err != nil {
		log.Printf("[%s] Ler handoff: %v", id, err)
		writeError(p.cfg.Failed, id, filepath.Base(path), err)
		p.moveFile(procPath, p.cfg.Failed) //nolint
		return
	}

	if p.cfg.ContextStore != nil {
		if err := p.recordHandoffEvent(procPath, data); err != nil {
			log.Printf("[%s] Registrar handoff: %v", id, err)
		}
	}

	select {
	case <-ctx.Done():
		p.moveFile(procPath, p.cfg.Failed) //nolint
		return
	default:
	}

	agentInput := string(data)
	if p.cfg.ContextStore != nil {
		if enriched, err := p.enrichHandoffWithContext(agentName, data); err == nil {
			agentInput = enriched
		} else {
			log.Printf("[%s] Contexto consultável indisponível: %v", id, err)
		}
	}

	response, err := p.callAgent(agentName, agentInput)
	if err != nil {
		log.Printf("[%s] Agente %s falhou: %v", id, agentName, err)
		writeError(p.cfg.Failed, id, filepath.Base(path), err)
		p.moveFile(procPath, p.cfg.Failed) //nolint
		return
	}

	if err := p.applyAgentResponse(agentName, id, response); err != nil {
		log.Printf("[%s] Aplicar resposta do agente %s: %v", id, agentName, err)
		writeError(p.cfg.Failed, id, filepath.Base(path), err)
		p.moveFile(procPath, p.cfg.Failed) //nolint
		return
	}

	p.moveFile(procPath, p.cfg.Success) //nolint
	log.Printf("[%s] Agente %s concluido", id, agentName)
}

func (p *Pipeline) callAgent(agentName, input string) (string, error) {
	provider := agent.ResolverProviderAgente(agentName)
	if strings.EqualFold(provider, "cli") ||
		strings.EqualFold(provider, "codex") ||
		strings.EqualFold(provider, "codex_cli") {
		return agent.ChamarAgenteComCodexCLI(
			agentName,
			input,
			p.cfg.CodexCLI,
			p.cfg.WorkspaceRoot,
			p.cfg.CodexTimeout,
		)
	}

	return agent.ChamarAgente(agentName, input)
}

// ── Processamento de arquivo novo ─────────────────────────────────────────────

func (p *Pipeline) handleRaw(ctx context.Context, path string) {
	ext := strings.ToLower(filepath.Ext(path))
	id := fileID(path)
	nome := filepath.Base(path)
	if ext == ".json" {
		if agentName, ok := detectHandoffRecipient(path); ok {
			log.Printf("[%s] JSON de handoff detectado em raw; redirecionando para agente %s", id, agentName)
			p.handleHandoff(ctx, path, agentName)
			return
		}
	}

	log.Printf("⚙️  [%s] Iniciando", id)

	// Mover para processing/
	procPath, err := p.moveFile(path, p.cfg.Processing)
	if err != nil {
		log.Printf("❌ [%s] Mover para processing: %v", id, err)
		return
	}

	var chunks []ChunksFile

	switch ext {
	case ".json":
		chunks, err = p.processNative(ctx, procPath, id, extractor.ExtractJSON)
	default:
		log.Printf("⚠️  [%s] Extensão não suportada: %s", id, ext)
		p.moveFile(procPath, p.cfg.Failed) //nolint
		return
	}

	if err != nil {
		log.Printf("❌ [%s] Extração: %v", id, err)
		writeError(p.cfg.Failed, id, nome, err)
		p.moveFile(procPath, p.cfg.Failed) //nolint
		return
	}

	// Continuar pipeline com os chunks
	if err := p.indexAndGenerate(ctx, id, nome, chunks); err != nil {
		log.Printf("❌ [%s] Pipeline: %v", id, err)
		writeError(p.cfg.Failed, id, nome, err)
		p.moveFile(procPath, p.cfg.Failed) //nolint
		return
	}

	p.moveFile(procPath, p.cfg.Success) //nolint
	log.Printf("✅ [%s] Concluído", id)
}

// ── Extração nativa (JSON, TXT) ───────────────────────────────────────────────

func (p *Pipeline) processNative(
	ctx context.Context,
	path, id string,
	extract func(string) (string, error),
) ([]ChunksFile, error) {
	texto, err := extract(path)
	if err != nil {
		return nil, err
	}
	rawChunks := chunkSemântico(texto, 1000, 100)

	var result []ChunksFile
	for i, c := range rawChunks {
		result = append(result, ChunksFile{
			Index:      i,
			Text:       c,
			Tokens:     estimarTokens(c),
			PageApprox: 0,
		})
	}
	return result, nil
}

// ── Delegação ao Python (PDF, vídeo) ─────────────────────────────────────────

// pythonMeta é o JSON gerado pelo Go para o Python consumir.
type pythonMeta struct {
	ID           string `json:"id"`
	OriginalFile string `json:"original_file"`
	Type         string `json:"type"`
	CreatedAt    string `json:"created_at"`
}

func (p *Pipeline) processPython(
	ctx context.Context,
	path, id, ext string,
) ([]ChunksFile, error) {
	tipo := "pdf"
	if ext != ".pdf" {
		tipo = "video"
	}

	// Salvar metadados para o Python
	meta := pythonMeta{
		ID:           id,
		OriginalFile: path,
		Type:         tipo,
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
	}
	metaBytes, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("serializar meta: %w", err)
	}
	metaPath := filepath.Join(p.cfg.PendingPython, id+".json")
	if err := os.WriteFile(metaPath, metaBytes, 0644); err != nil {
		return nil, fmt.Errorf("salvar meta: %w", err)
	}

	// Mover arquivo para pending_python/
	pendingPath := filepath.Join(p.cfg.PendingPython, filepath.Base(path))
	if err := os.Rename(path, pendingPath); err != nil {
		return nil, fmt.Errorf("mover para pending_python: %w", err)
	}

	log.Printf("🐍 [%s] Aguardando Python (timeout: %s)", id, p.cfg.PythonTimeout)

	// Registrar canal de notificação
	ch := make(chan string, 1)
	p.pendingMu.Lock()
	p.pending[id] = ch
	p.pendingMu.Unlock()
	defer func() {
		p.pendingMu.Lock()
		delete(p.pending, id)
		p.pendingMu.Unlock()
	}()

	// Verificar se o arquivo já chegou antes de registrar (race condition)
	chunksPath := filepath.Join(p.cfg.Extracted, id+".chunks.json")
	if _, err := os.Stat(chunksPath); err == nil {
		return loadChunksFile(chunksPath)
	}

	// Aguardar notificação ou timeout
	select {
	case path := <-ch:
		return loadChunksFile(path)
	case <-time.After(p.cfg.PythonTimeout):
		return nil, fmt.Errorf(
			"timeout aguardando Python após %s — rode extractor.py manualmente",
			p.cfg.PythonTimeout,
		)
	case <-ctx.Done():
		return nil, fmt.Errorf("cancelado")
	}
}

// ── Embeddings + indexação + nota ────────────────────────────────────────────

type embeddedChunk struct {
	ID         string
	Text       string
	PageApprox int
	Embedding  []float64
}

func (p *Pipeline) indexAndGenerate(
	ctx context.Context,
	id, nome string,
	chunks []ChunksFile,
) error {
	log.Printf("🧮 [%s] Gerando embeddings (%d chunks)...", id, len(chunks))

	type result struct {
		chunk embeddedChunk
		err   error
	}

	chunkCh := make(chan ChunksFile, len(chunks))
	for _, c := range chunks {
		chunkCh <- c
	}
	close(chunkCh)

	resultCh := make(chan result, len(chunks))
	var wg sync.WaitGroup

	for range p.cfg.EmbedWorkerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for c := range chunkCh {
				select {
				case <-ctx.Done():
					return
				default:
				}
				emb, err := embeddings.Generate(c.Text)
				resultCh <- result{
					chunk: embeddedChunk{
						ID:         fmt.Sprintf("%s_%d", id, c.Index),
						Text:       c.Text,
						PageApprox: c.PageApprox,
						Embedding:  emb,
					},
					err: err,
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var embedded []embeddedChunk
	for r := range resultCh {
		if r.err != nil {
			log.Printf("⚠️  [%s] Embedding falhou: %v", id, r.err)
			continue
		}
		embedded = append(embedded, r.chunk)
	}

	if len(embedded) == 0 {
		return fmt.Errorf("nenhum embedding gerado")
	}

	// TODO: Indexar no SQLite (vectorstore)
	log.Printf("💾 [%s] Indexando %d chunks...", id, len(embedded))

	// Montar slice de textos para a nota — limita aos 5 primeiros chunks
	// (substitui a busca por similaridade até o vectorstore ser reativado)
	textos := make([]string, 0, min(5, len(embedded)))
	for _, c := range embedded {
		textos = append(textos, c.Text)
		if len(textos) == 5 {
			break
		}
	}

	// Gerar nota
	log.Printf("✍️ [%s] Recebendo resposta...", id)
	nota, err := agent.GerarArquivo(id, textos)
	if err != nil {
		return fmt.Errorf("gerar nota: %w", err)
	}

	// Salvar no vault
	return code_writer.SaveFile(id, nome, nota)
}

// ── Helpers ───────────────────────────────────────────────────────────────────

type agentAction struct {
	Action    string                  `json:"action"`
	Files     []code_writer.CodeFile  `json:"files"`
	File      *code_writer.CodeFile   `json:"file"`
	Handoffs  []json.RawMessage       `json:"handoffs"`
	Handoff   json.RawMessage         `json:"handoff"`
	Header    *hand_off.HandoffHeader `json:"header"`
	Payload   json.RawMessage         `json:"payload"`
	Message   string                  `json:"message"`
	Questions []humanQuestion         `json:"questions"`
	Question  string                  `json:"question"`
	Blocking  bool                    `json:"blocking"`
	Priority  string                  `json:"priority"`
	context_store.ContextEntry
	context_store.PlanningItem
	context_store.TestRecord
}

type humanQuestion struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Priority string `json:"priority"`
	Blocking bool   `json:"blocking"`
	TaskRef  string `json:"task_ref"`
	AskedBy  string `json:"asked_by"`
	Response string `json:"response"`
}

type humanHandoffPayload struct {
	ResponseKey  string          `json:"response_key"`
	RespondTo    string          `json:"respond_to"`
	Instructions string          `json:"instructions"`
	Questions    []humanQuestion `json:"questions"`
	Response     string          `json:"response"`
}

func (p *Pipeline) applyAgentResponse(agentName, taskID, response string) error {
	raw, err := extractJSONPayload(response)
	if err != nil {
		return fmt.Errorf("resposta do agente sem JSON valido: %w", err)
	}

	items, err := splitAgentActions(raw)
	if err != nil {
		return fmt.Errorf("resposta do agente nao e uma acao JSON valida: %w", err)
	}

	for _, item := range items {
		var action agentAction
		if err := json.Unmarshal(item, &action); err != nil {
			return fmt.Errorf("parsear ação do agente: %w", err)
		}

		if action.Header != nil {
			if _, err := hand_off.SaveRawHandoff(p.cfg.HandoffDir, item); err != nil {
				return err
			}
			continue
		}

		if action.File != nil {
			action.Files = append(action.Files, *action.File)
		}

		switch strings.ToLower(action.Action) {
		case "write_code", "create_file", "write_files":
			files, err := normalizeCodeFiles(action.Files)
			if err != nil {
				return err
			}
			if _, err := code_writer.WriteCodeFiles(p.cfg.WorkspaceRoot, files); err != nil {
				return err
			}
		case "update_context", "record_context", "record_decision", "record_retrospective":
			if p.cfg.ContextStore == nil {
				return p.saveAgentOutput(agentName, taskID, response)
			}
			entry := action.ContextEntry
			entry.SourceAgent = defaultString(entry.SourceAgent, agentName)
			entry.TaskRef = defaultString(entry.TaskRef, taskID)
			if entry.EntryType == "" {
				entry.EntryType = strings.ToLower(action.Action)
			}
			if err := p.cfg.ContextStore.SaveContextEntry(entry); err != nil {
				return err
			}
		case "update_plan", "record_plan", "update_state", "record_task":
			if p.cfg.ContextStore == nil {
				return p.saveAgentOutput(agentName, taskID, response)
			}
			item := action.PlanningItem
			item.SourceAgent = defaultString(item.SourceAgent, agentName)
			item.TaskRef = defaultString(item.TaskRef, taskID)
			if item.ItemType == "" {
				item.ItemType = strings.ToLower(action.Action)
			}
			if err := p.cfg.ContextStore.SavePlanningItem(item); err != nil {
				return err
			}
		case "record_test":
			if p.cfg.ContextStore == nil {
				return p.saveAgentOutput(agentName, taskID, response)
			}
			record := action.TestRecord
			record.SourceAgent = defaultString(record.SourceAgent, agentName)
			record.TaskRef = defaultString(record.TaskRef, taskID)
			if err := p.cfg.ContextStore.SaveTestRecord(record); err != nil {
				return err
			}
		case "ask_human", "question", "clarification_for_human":
			if err := p.createHumanHandoff(agentName, taskID, action); err != nil {
				return err
			}
		case "handoff":
			if len(action.Handoff) > 0 {
				action.Handoffs = append(action.Handoffs, action.Handoff)
			}
			for _, handoff := range action.Handoffs {
				if _, err := hand_off.SaveRawHandoff(p.cfg.HandoffDir, handoff); err != nil {
					return err
				}
			}
		case "note", "":
			if len(action.Files) > 0 {
				files, err := normalizeCodeFiles(action.Files)
				if err != nil {
					return err
				}
				if _, err := code_writer.WriteCodeFiles(p.cfg.WorkspaceRoot, files); err != nil {
					return err
				}
				continue
			}
			if len(action.Handoffs) > 0 {
				for _, handoff := range action.Handoffs {
					if _, err := hand_off.SaveRawHandoff(p.cfg.HandoffDir, handoff); err != nil {
						return err
					}
				}
				continue
			}
			if action.Message != "" {
				if err := p.saveAgentOutput(agentName, taskID, action.Message); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("acao de agente desconhecida: %q", action.Action)
		}
	}

	return nil
}

func detectHandoffRecipient(path string) (string, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}

	var parsed struct {
		Header hand_off.HandoffHeader `json:"header"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return "", false
	}

	recipient := strings.TrimSpace(parsed.Header.Recipient)
	if recipient == "" {
		return "", false
	}
	recipient = strings.Trim(recipient, "[]")
	if idx := strings.Index(recipient, ":"); idx >= 0 {
		recipient = recipient[:idx]
	}
	if strings.EqualFold(recipient, "HUMAN") {
		return "", false
	}
	return strings.ToUpper(recipient), true
}

func normalizeCodeFiles(files []code_writer.CodeFile) ([]code_writer.CodeFile, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("write_code sem arquivos")
	}

	normalized := make([]code_writer.CodeFile, 0, len(files))
	for _, file := range files {
		path := filepath.Clean(strings.TrimSpace(file.Path))
		if path == "." || path == "" {
			return nil, fmt.Errorf("write_code com arquivo sem path")
		}
		if filepath.Ext(path) == "" {
			return nil, fmt.Errorf("write_code path sem extensao: %s", file.Path)
		}

		slashPath := filepath.ToSlash(path)
		if strings.EqualFold(filepath.Ext(path), ".md") && !strings.HasPrefix(slashPath, "wiki/") {
			return nil, fmt.Errorf("write_code para .md fora de wiki/ nao permitido: %s", file.Path)
		}

		content := strings.TrimSpace(file.Content)
		content = stripSingleMarkdownFence(content)
		if strings.Contains(content, "```") {
			return nil, fmt.Errorf("write_code contem bloco markdown em %s", file.Path)
		}
		if content == "" {
			return nil, fmt.Errorf("write_code com conteudo vazio em %s", file.Path)
		}

		normalized = append(normalized, code_writer.CodeFile{
			Path:    path,
			Content: content,
		})
	}
	return normalized, nil
}

func stripSingleMarkdownFence(content string) string {
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) < 3 {
		return strings.TrimSpace(content)
	}
	if !strings.HasPrefix(strings.TrimSpace(lines[0]), "```") {
		return strings.TrimSpace(content)
	}
	if strings.TrimSpace(lines[len(lines)-1]) != "```" {
		return strings.TrimSpace(content)
	}
	return strings.TrimSpace(strings.Join(lines[1:len(lines)-1], "\n"))
}

func extractJSONPayload(response string) ([]byte, error) {
	text := strings.TrimSpace(response)
	if strings.HasPrefix(text, "```") {
		lines := strings.Split(text, "\n")
		if len(lines) >= 3 {
			text = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}

	objectStart := strings.Index(text, "{")
	arrayStart := strings.Index(text, "[")
	start := -1
	if objectStart == -1 {
		start = arrayStart
	} else if arrayStart == -1 || objectStart < arrayStart {
		start = objectStart
	} else {
		start = arrayStart
	}
	if start == -1 {
		return nil, fmt.Errorf("resposta sem JSON")
	}

	objectEnd := strings.LastIndex(text, "}")
	arrayEnd := strings.LastIndex(text, "]")
	end := max(objectEnd, arrayEnd)
	if end < start {
		return nil, fmt.Errorf("json incompleto")
	}

	return []byte(strings.TrimSpace(text[start : end+1])), nil
}

func splitAgentActions(raw []byte) ([]json.RawMessage, error) {
	var list []json.RawMessage
	if err := json.Unmarshal(raw, &list); err == nil {
		return list, nil
	}

	var single json.RawMessage
	if err := json.Unmarshal(raw, &single); err != nil {
		return nil, err
	}
	return []json.RawMessage{single}, nil
}

func (p *Pipeline) saveAgentOutput(agentName, taskID, content string) error {
	if p.cfg.ContextStore != nil {
		return p.cfg.ContextStore.SaveAgentOutput(agentName, taskID, "note", content)
	}

	dir := filepath.Join(p.cfg.AgentOutputDir, strings.ToUpper(agentName))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	name := fmt.Sprintf("%s_%s.md", time.Now().Format("20060102_150405"), taskID)
	return os.WriteFile(filepath.Join(dir, name), []byte(content), 0644)
}

func (p *Pipeline) createHumanHandoff(agentName, taskID string, action agentAction) error {
	questions := action.Questions
	if strings.TrimSpace(action.Question) != "" {
		questions = append(questions, humanQuestion{
			Question: action.Question,
			Priority: action.Priority,
			Blocking: action.Blocking,
		})
	}
	if len(questions) == 0 {
		return fmt.Errorf("ask_human sem perguntas")
	}

	responseKey := fmt.Sprintf("response_%s_%d", sanitizeKey(taskID), time.Now().Unix())
	for i := range questions {
		if questions[i].ID == "" {
			questions[i].ID = fmt.Sprintf("%s_q%d", responseKey, i+1)
		}
		if questions[i].TaskRef == "" {
			questions[i].TaskRef = taskID
		}
		if questions[i].AskedBy == "" {
			questions[i].AskedBy = strings.ToUpper(agentName)
		}
		if questions[i].Priority == "" {
			questions[i].Priority = "Medium"
		}
	}

	payload := humanHandoffPayload{
		ResponseKey: responseKey,
		RespondTo:   fmt.Sprintf("[%s]", strings.ToUpper(agentName)),
		Instructions: "Preencha a chave response ou questions[].response. Depois renomeie este arquivo trocando _TO_HUMAN_ por _TO_" +
			strings.ToUpper(agentName) + "_ para retornar ao pipeline.",
		Questions: questions,
		Response:  "",
	}

	header := hand_off.HandoffHeader{
		Sender:    fmt.Sprintf("[%s]", strings.ToUpper(agentName)),
		Recipient: "[HUMAN]",
		TaskRef:   taskID,
		Intent:    "HUMAN_CLARIFICATION_REQUEST",
	}

	_, err := hand_off.CreateHandoff(p.cfg.HandoffDir, header, payload)
	return err
}

func sanitizeKey(value string) string {
	var sb strings.Builder
	for _, r := range strings.ToLower(value) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
			continue
		}
		sb.WriteRune('_')
	}
	result := strings.Trim(sb.String(), "_")
	if result == "" {
		return "question"
	}
	return result
}

func (p *Pipeline) recordHandoffEvent(path string, data []byte) error {
	var parsed struct {
		Header  hand_off.HandoffHeader `json:"header"`
		Payload json.RawMessage        `json:"payload"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}

	return p.cfg.ContextStore.SaveHandoffEvent(context_store.HandoffEvent{
		Path:      path,
		Sender:    parsed.Header.Sender,
		Recipient: parsed.Header.Recipient,
		TaskRef:   parsed.Header.TaskRef,
		Intent:    parsed.Header.Intent,
		Payload:   string(parsed.Payload),
		RawJSON:   string(data),
	})
}

func (p *Pipeline) enrichHandoffWithContext(agentName string, data []byte) (string, error) {
	var parsed struct {
		Header hand_off.HandoffHeader `json:"header"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return string(data), err
	}

	contextItems, err := p.cfg.ContextStore.QueryContextForHandoff(agentName, parsed.Header.TaskRef, 8)
	if err != nil {
		return string(data), err
	}
	if len(contextItems) == 0 {
		return string(data), nil
	}

	var sb strings.Builder
	sb.WriteString("HANDOFF:\n")
	sb.Write(data)
	sb.WriteString("\n\nCONTEXTO_CONSULTAVEL:\n")
	for i, item := range contextItems {
		sb.WriteString(fmt.Sprintf("\n--- contexto %d ---\n%s\n", i+1, item))
	}
	return sb.String(), nil
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}

// ChunksFile representa um chunk no .chunks.json do Python.
type ChunksFile struct {
	Index      int    `json:"index"`
	Text       string `json:"text"`
	Tokens     int    `json:"tokens"`
	PageApprox int    `json:"page_approx"`
}

type chunksJSON struct {
	Chunks []ChunksFile `json:"chunks"`
}

func loadChunksFile(path string) ([]ChunksFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ler chunks.json: %w", err)
	}
	var cf chunksJSON
	if err := json.Unmarshal(data, &cf); err != nil {
		return nil, fmt.Errorf("parsear chunks.json: %w", err)
	}
	if len(cf.Chunks) == 0 {
		return nil, fmt.Errorf("chunks.json vazio")
	}
	return cf.Chunks, nil
}

func fileID(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func (p *Pipeline) moveFile(src, destDir string) (string, error) {
	p.fileMu.Lock()
	defer p.fileMu.Unlock()
	return moveFile(src, destDir)
}

func moveFile(src, destDir string) (string, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}
	base := filepath.Base(src)
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	for attempt := 0; ; attempt++ {
		dest := filepath.Join(destDir, base)
		if attempt > 0 {
			dest = filepath.Join(destDir, fmt.Sprintf("%s_%s_%02d%s", stem, time.Now().UTC().Format("20060102_150405.000000000"), attempt, ext))
		}
		if _, err := os.Stat(dest); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return "", err
		}
		if err := os.Rename(src, dest); err != nil {
			if os.IsExist(err) {
				continue
			}
			return "", err
		}
		return dest, nil
	}
}

func writeError(failedDir, id, nome string, err error) {
	os.MkdirAll(failedDir, 0755) //nolint
	content := fmt.Sprintf(
		"Arquivo : %s\nID      : %s\nData    : %s\nErro    :\n%v\n",
		nome, id, time.Now().UTC().Format(time.RFC3339), err,
	)
	errPath := filepath.Join(failedDir, id+".error.txt")
	os.WriteFile(errPath, []byte(content), 0644) //nolint
}

// estimarTokens — aproximação simples sem tokenizador externo.
func estimarTokens(s string) int {
	return max(1, len(s)/4)
}

// chunkSemântico — Estratégia D: agrupa parágrafos naturais.
func chunkSemântico(texto string, maxTokens, minTokens int) []string {
	blocos := splitBlocos(texto)
	var chunks []string
	var buffer []string
	bufTokens := 0

	flush := func() {
		if len(buffer) == 0 {
			return
		}
		chunks = append(chunks, strings.Join(buffer, "\n\n"))
		buffer = nil
		bufTokens = 0
	}

	for _, bloco := range blocos {
		t := estimarTokens(bloco)
		isTitulo := strings.HasPrefix(bloco, "#")

		if isTitulo && len(buffer) > 0 {
			flush()
		}

		if t > maxTokens {
			flush()
			for _, sub := range splitSentencas(bloco, maxTokens) {
				chunks = append(chunks, sub)
			}
			continue
		}

		if bufTokens+t > maxTokens {
			flush()
		}

		buffer = append(buffer, bloco)
		bufTokens += t
	}
	flush()
	return chunks
}

func splitBlocos(texto string) []string {
	raw := strings.Split(texto, "\n\n")
	var blocos []string
	for _, b := range raw {
		b = strings.TrimSpace(b)
		if b != "" {
			blocos = append(blocos, b)
		}
	}
	return blocos
}

func splitSentencas(texto string, maxTokens int) []string {
	var sentencas []string
	current := ""
	for _, r := range texto {
		current += string(r)
		if r == '.' || r == '!' || r == '?' {
			if s := strings.TrimSpace(current); s != "" {
				sentencas = append(sentencas, s)
			}
			current = ""
		}
	}
	if s := strings.TrimSpace(current); s != "" {
		sentencas = append(sentencas, s)
	}

	var chunks []string
	var buf []string
	bufT := 0
	for _, s := range sentencas {
		t := estimarTokens(s)
		if bufT+t > maxTokens && len(buf) > 0 {
			chunks = append(chunks, strings.Join(buf, " "))
			buf = nil
			bufT = 0
		}
		buf = append(buf, s)
		bufT += t
	}
	if len(buf) > 0 {
		chunks = append(chunks, strings.Join(buf, " "))
	}
	return chunks
}
