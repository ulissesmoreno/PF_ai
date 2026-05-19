package pipeline

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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
	data = hand_off.StripUTF8BOM(data)

	currentCardID := cardIDFromHandoff(data)
	currentProjectID := projectIDFromHandoff(data)
	currentIntent := intentFromHandoff(data)
	if p.cfg.ContextStore != nil {
		if err := p.ensureProjectFromHandoff(data); err != nil {
			log.Printf("[%s] Restaurar projeto do handoff: %v", id, err)
		}
		cardID, projectID, err := p.recordCardHandoff(procPath, data)
		if cardID != "" {
			currentCardID = cardID
		}
		if projectID != "" {
			currentProjectID = projectID
		}
		if err != nil {
			log.Printf("[%s] Registrar card: %v", id, err)
		}
		if err := p.recordHandoffEvent(procPath, data, currentProjectID); err != nil {
			log.Printf("[%s] Registrar handoff: %v", id, err)
		}
		handled, err := p.applyProjectOnboardingResponse(id, data)
		if err != nil {
			log.Printf("[%s] Aplicar onboarding de projeto: %v", id, err)
			writeError(p.cfg.Failed, id, filepath.Base(path), err)
			p.moveFile(procPath, p.cfg.Failed) //nolint
			return
		}
		if handled {
			p.updateCardStatus(currentCardID, "done")
			p.moveFile(procPath, p.cfg.Success) //nolint
			log.Printf("[%s] Projeto criado via onboarding", id)
			return
		}
		p.updateCardStatus(currentCardID, "in_progress")
	}

	select {
	case <-ctx.Done():
		p.moveFile(procPath, p.cfg.Failed) //nolint
		return
	default:
	}

	agentInput := string(data)
	if p.cfg.ContextStore != nil {
		if enriched, err := p.enrichHandoffWithContext(agentName, data, currentProjectID); err == nil {
			agentInput = enriched
		} else {
			log.Printf("[%s] Contexto consultável indisponível: %v", id, err)
		}
	}

	response, err := p.callAgent(agentName, agentInput, currentCardID, currentProjectID)
	if err != nil {
		log.Printf("[%s] Agente %s falhou: %v", id, agentName, err)
		if p.retryCardHandoff(procPath, currentCardID, err) {
			return
		}
		p.moveFile(procPath, p.cfg.Failed) //nolint
		return
	}

	result, err := p.applyAgentResponseForProject(currentProjectID, agentName, id, currentCardID, response)
	if err != nil {
		log.Printf("[%s] Aplicar resposta do agente %s: %v", id, agentName, err)
		if p.retryCardHandoff(procPath, currentCardID, err) {
			return
		}
		p.moveFile(procPath, p.cfg.Failed) //nolint
		return
	}
	if shouldAutoDelegateKickoff(agentName, currentIntent, result) {
		usedCardID, err := p.createLeanImplementationHandoff(currentProjectID, id, currentCardID, data, response)
		if err != nil {
			log.Printf("[%s] Criar handoff lean apos kickoff: %v", id, err)
			if p.retryCardHandoff(procPath, currentCardID, err) {
				return
			}
			p.moveFile(procPath, p.cfg.Failed) //nolint
			return
		}
		result.Delegated = true
		if usedCardID != "" {
			result.CardID = usedCardID
		}
	}
	if shouldAutoCompleteTask(agentName, currentIntent, result) {
		usedCardID, err := p.createTaskCompleteHandoff(currentProjectID, id, currentCardID, agentName, currentIntent, response)
		if err != nil {
			log.Printf("[%s] Criar handoff de conclusao para CEO: %v", id, err)
			if p.retryCardHandoff(procPath, currentCardID, err) {
				return
			}
			p.moveFile(procPath, p.cfg.Failed) //nolint
			return
		}
		result.Delegated = true
		if usedCardID != "" {
			result.CardID = usedCardID
		}
	}

	p.recordCardResponse(result.CardID, agentName, response)
	if result.Delegated {
		p.updateCardStatus(result.CardID, "in_progress")
	} else if !result.Blocked {
		p.updateCardStatus(result.CardID, "done")
	}
	p.moveFile(procPath, p.cfg.Success) //nolint
	log.Printf("[%s] Agente %s concluido", id, agentName)
}

func (p *Pipeline) callAgent(agentName, input, cardID, projectID string) (string, error) {
	provider := agent.ResolverProviderAgente(agentName)
	model := mustAgentModel(agentName)
	tier := agent.ResolverTierAgente(agentName)
	if strings.TrimSpace(projectID) == "" {
		projectID = projectIDForToken(p.cfg.ContextStore, cardID)
	}
	log.Printf("[AGENT_CALL] status=calling agent=%s model=%s tier=%d card_id=%s project_id=%s", strings.ToUpper(agentName), model, tier, cardID, projectID)
	if strings.EqualFold(provider, "cli") ||
		strings.EqualFold(provider, "codex") ||
		strings.EqualFold(provider, "codex_cli") {
		start := time.Now()
		response, err := agent.ChamarAgenteComCodexCLI(
			agentName,
			input,
			p.cfg.CodexCLI,
			p.cfg.WorkspaceRoot,
			p.cfg.CodexTimeout,
		)
		if err == nil {
			p.recordTokenUsage(cardID, agentName, agent.TokenUsage{
				AgentName:        agentName,
				Model:            model,
				Tier:             tier,
				PromptTokens:     estimarTokens(input),
				CompletionTokens: estimarTokens(response),
				LatencyMS:        time.Since(start).Milliseconds(),
			})
		}
		p.logAgentResult(agentName, model, tier, cardID, projectID, start, err)
		return response, err
	}

	start := time.Now()
	response, usage, err := agent.ChamarAgenteComUso(agentName, input)
	if err == nil {
		p.recordTokenUsage(cardID, agentName, usage)
	}
	p.logAgentResult(agentName, model, tier, cardID, projectID, start, err)
	return response, err
}

func (p *Pipeline) logAgentResult(agentName, model string, tier int, cardID, projectID string, start time.Time, err error) {
	status := "done"
	if err != nil {
		status = "failed"
	}
	log.Printf("[AGENT_CALL] status=%s agent=%s model=%s tier=%d card_id=%s project_id=%s latency_ms=%d", status, strings.ToUpper(agentName), model, tier, cardID, projectID, time.Since(start).Milliseconds())
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
	Action            string                  `json:"action"`
	ReplyToCardID     string                  `json:"reply_to_card_id"`
	Complexity        string                  `json:"complexity"`
	ExecutionMode     string                  `json:"execution_mode"`
	RecommendedAgents []string                `json:"recommended_agents"`
	Files             []code_writer.CodeFile  `json:"files"`
	File              *code_writer.CodeFile   `json:"file"`
	Handoffs          []json.RawMessage       `json:"handoffs"`
	Handoff           json.RawMessage         `json:"handoff"`
	Header            *hand_off.HandoffHeader `json:"header"`
	Payload           json.RawMessage         `json:"payload"`
	Message           string                  `json:"message"`
	Questions         []humanQuestion         `json:"questions"`
	Question          string                  `json:"question"`
	Blocking          bool                    `json:"blocking"`
	Priority          string                  `json:"priority"`
	context_store.ContextEntry
	context_store.PlanningItem
	context_store.TestRecord
}

type agentResponseResult struct {
	Blocked   bool
	Delegated bool
	CardID    string
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

type projectOnboardingPayload struct {
	Action         string                 `json:"action"`
	Project        *context_store.Project `json:"project"`
	Name           string                 `json:"name"`
	Slug           string                 `json:"slug"`
	Domain         string                 `json:"domain"`
	Description    string                 `json:"description"`
	TargetAudience string                 `json:"target_audience"`
	MainObjective  string                 `json:"main_objective"`
	StackBackend   string                 `json:"stack_backend"`
	StackFrontend  string                 `json:"stack_frontend"`
	StackDatabase  string                 `json:"stack_database"`
	StackInfra     string                 `json:"stack_infra"`
	StackNotes     string                 `json:"stack_notes"`
}

func (p *Pipeline) applyAgentResponse(agentName, taskID, cardID, response string) (agentResponseResult, error) {
	return p.applyAgentResponseForProject("", agentName, taskID, cardID, response)
}

func (p *Pipeline) applyAgentResponseForProject(projectID, agentName, taskID, cardID, response string) (agentResponseResult, error) {
	result := agentResponseResult{CardID: cardID}
	raw, err := extractJSONPayload(response)
	if err != nil {
		return result, fmt.Errorf("resposta do agente sem JSON valido: %w", err)
	}

	items, err := splitAgentActions(raw)
	if err != nil {
		return result, fmt.Errorf("resposta do agente nao e uma acao JSON valida: %w", err)
	}

	for _, item := range items {
		var action agentAction
		if err := json.Unmarshal(item, &action); err != nil {
			return result, fmt.Errorf("parsear ação do agente: %w", err)
		}

		if action.ReplyToCardID != "" {
			merged, err := mergeReplyCardID(result.CardID, cardID, action.ReplyToCardID)
			if err != nil {
				return result, err
			}
			result.CardID = merged
		}

		if action.Header != nil {
			usedCardID, err := p.saveRawHandoffForCard(projectID, item, result.CardID)
			if err != nil {
				return result, err
			}
			result.Delegated = true
			if usedCardID != "" {
				merged, err := mergeReplyCardID(result.CardID, cardID, usedCardID)
				if err != nil {
					return result, err
				}
				result.CardID = merged
			}
			continue
		}

		if action.File != nil {
			action.Files = append(action.Files, *action.File)
		}
		inferAgentAction(&action)
		if strings.TrimSpace(action.Action) == "" &&
			len(action.Files) == 0 &&
			len(action.Handoffs) == 0 &&
			len(action.Handoff) == 0 &&
			action.Header == nil &&
			strings.TrimSpace(action.Message) == "" {
			return result, fmt.Errorf("acao de agente sem action executavel")
		}

		switch strings.ToLower(action.Action) {
		case "write_code", "create_file", "write_files":
			files, err := normalizeCodeFiles(action.Files)
			if err != nil {
				return result, err
			}
			if strings.EqualFold(agentName, "CEO") {
				usedCardID, err := p.createDeliveryHandoff(projectID, taskID, result.CardID, files)
				if err != nil {
					return result, err
				}
				result.Delegated = true
				if usedCardID != "" {
					merged, err := mergeReplyCardID(result.CardID, cardID, usedCardID)
					if err != nil {
						return result, err
					}
					result.CardID = merged
				}
				continue
			}
			if _, err := code_writer.WriteCodeFiles(p.workspaceRootForProject(projectID), files); err != nil {
				return result, err
			}
		case "update_context", "record_context", "record_decision", "record_retrospective":
			if p.cfg.ContextStore == nil {
				return result, p.saveAgentOutput(projectID, agentName, taskID, response)
			}
			entry := action.ContextEntry
			entry.ProjectID = defaultString(entry.ProjectID, projectID)
			entry.SourceAgent = defaultString(entry.SourceAgent, agentName)
			entry.TaskRef = defaultString(entry.TaskRef, taskID)
			if entry.EntryType == "" {
				entry.EntryType = strings.ToLower(action.Action)
			}
			if err := p.cfg.ContextStore.SaveContextEntry(entry); err != nil {
				return result, err
			}
		case "update_plan", "record_plan", "update_state", "record_task":
			if p.cfg.ContextStore == nil {
				if err := p.saveActionHandoffs(projectID, action, cardID, &result); err != nil {
					return result, err
				}
				return result, p.saveAgentOutput(projectID, agentName, taskID, response)
			}
			item := action.PlanningItem
			item.ProjectID = defaultString(item.ProjectID, projectID)
			item.SourceAgent = defaultString(item.SourceAgent, agentName)
			item.TaskRef = defaultString(item.TaskRef, taskID)
			if item.ItemType == "" {
				item.ItemType = strings.ToLower(action.Action)
			}
			if err := p.cfg.ContextStore.SavePlanningItem(item); err != nil {
				return result, err
			}
			if err := p.saveActionHandoffs(projectID, action, cardID, &result); err != nil {
				return result, err
			}
		case "record_test":
			if p.cfg.ContextStore == nil {
				return result, p.saveAgentOutput(projectID, agentName, taskID, response)
			}
			record := action.TestRecord
			record.ProjectID = defaultString(record.ProjectID, projectID)
			record.SourceAgent = defaultString(record.SourceAgent, agentName)
			record.TaskRef = defaultString(record.TaskRef, taskID)
			if err := p.cfg.ContextStore.SaveTestRecord(record); err != nil {
				return result, err
			}
		case "ask_human", "question", "clarification_for_human":
			if err := p.createHumanHandoff(projectID, agentName, taskID, result.CardID, action); err != nil {
				return result, err
			}
			p.updateCardStatus(result.CardID, "blocked")
			result.Blocked = true
		case "handoff":
			if err := p.saveActionHandoffs(projectID, action, cardID, &result); err != nil {
				return result, err
			}
		case "note", "":
			if len(action.Files) > 0 {
				files, err := normalizeCodeFiles(action.Files)
				if err != nil {
					return result, err
				}
				if strings.EqualFold(agentName, "CEO") {
					usedCardID, err := p.createDeliveryHandoff(projectID, taskID, result.CardID, files)
					if err != nil {
						return result, err
					}
					result.Delegated = true
					if usedCardID != "" {
						merged, err := mergeReplyCardID(result.CardID, cardID, usedCardID)
						if err != nil {
							return result, err
						}
						result.CardID = merged
					}
					continue
				}
				if _, err := code_writer.WriteCodeFiles(p.workspaceRootForProject(projectID), files); err != nil {
					return result, err
				}
				continue
			}
			if len(action.Handoffs) > 0 {
				for _, handoff := range action.Handoffs {
					usedCardID, err := p.saveRawHandoffForCard(projectID, handoff, result.CardID)
					if err != nil {
						return result, err
					}
					result.Delegated = true
					if usedCardID != "" {
						merged, err := mergeReplyCardID(result.CardID, cardID, usedCardID)
						if err != nil {
							return result, err
						}
						result.CardID = merged
					}
				}
				continue
			}
			if action.Message != "" {
				if err := p.saveAgentOutput(projectID, agentName, taskID, action.Message); err != nil {
					return result, err
				}
			}
		default:
			return result, fmt.Errorf("acao de agente desconhecida: %q", action.Action)
		}
	}

	return result, nil
}

func inferAgentAction(action *agentAction) {
	if action == nil || strings.TrimSpace(action.Action) != "" {
		return
	}
	if len(action.Files) > 0 {
		action.Action = "write_code"
		return
	}
	if len(action.Handoffs) > 0 || len(action.Handoff) > 0 {
		action.Action = "handoff"
		return
	}
	if strings.TrimSpace(action.Complexity) != "" ||
		strings.TrimSpace(action.ExecutionMode) != "" ||
		len(action.RecommendedAgents) > 0 ||
		strings.TrimSpace(action.PlanningItem.Title) != "" ||
		strings.TrimSpace(action.PlanningItem.Content) != "" {
		action.Action = "update_plan"
		if strings.TrimSpace(action.PlanningItem.ItemType) == "" {
			action.PlanningItem.ItemType = "plan"
		}
		if strings.TrimSpace(action.PlanningItem.Title) == "" {
			action.PlanningItem.Title = "Complexity decision"
		}
		if strings.TrimSpace(action.PlanningItem.Content) == "" {
			action.PlanningItem.Content = fmt.Sprintf(
				"complexity=%s execution_mode=%s recommended_agents=%s",
				action.Complexity,
				action.ExecutionMode,
				strings.Join(action.RecommendedAgents, ","),
			)
		}
	}
}

func (p *Pipeline) saveActionHandoffs(projectID string, action agentAction, fallbackCardID string, result *agentResponseResult) error {
	if len(action.Handoff) > 0 {
		action.Handoffs = append(action.Handoffs, action.Handoff)
	}
	for _, handoff := range action.Handoffs {
		usedCardID, err := p.saveRawHandoffForCard(projectID, handoff, result.CardID)
		if err != nil {
			return err
		}
		result.Delegated = true
		if usedCardID != "" {
			merged, err := mergeReplyCardID(result.CardID, fallbackCardID, usedCardID)
			if err != nil {
				return err
			}
			result.CardID = merged
		}
	}
	return nil
}

func detectHandoffRecipient(path string) (string, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	data = hand_off.StripUTF8BOM(data)

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
		if strings.EqualFold(filepath.Ext(path), ".md") &&
			!strings.HasPrefix(slashPath, "wiki/") &&
			!strings.HasPrefix(slashPath, "docs/") &&
			!strings.HasPrefix(slashPath, "vault/") {
			return nil, fmt.Errorf("write_code para .md fora de wiki/, docs/ ou vault/ nao permitido: %s", file.Path)
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

	for start, r := range text {
		if r != '{' && r != '[' {
			continue
		}
		fragment := text[start:]
		var raw json.RawMessage
		decoder := json.NewDecoder(strings.NewReader(fragment))
		if err := decoder.Decode(&raw); err != nil {
			continue
		}
		if len(raw) > 0 {
			return []byte(strings.TrimSpace(string(raw))), nil
		}
	}

	return nil, fmt.Errorf("resposta sem JSON valido")
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

func (p *Pipeline) saveRawHandoffForCard(projectID string, data []byte, cardID string) (string, error) {
	var handoff struct {
		ReplyToCardID string                 `json:"reply_to_card_id"`
		Header        hand_off.HandoffHeader `json:"header"`
		Payload       json.RawMessage        `json:"payload"`
	}
	if err := json.Unmarshal(data, &handoff); err != nil {
		return "", err
	}
	targetCardID := strings.TrimSpace(handoff.ReplyToCardID)
	if targetCardID == "" {
		targetCardID = strings.TrimSpace(cardID)
	}
	if targetCardID != "" {
		handoff.Header.CardID = targetCardID
	}
	if strings.TrimSpace(handoff.Header.ProjectID) == "" {
		handoff.Header.ProjectID = strings.TrimSpace(projectID)
	}

	normalized, err := json.Marshal(hand_off.HandoffSchema[json.RawMessage]{
		Header:  handoff.Header,
		Payload: handoff.Payload,
	})
	if err != nil {
		return "", err
	}
	if _, err := hand_off.SaveRawHandoff(p.cfg.HandoffDir, normalized); err != nil {
		return "", err
	}
	return strings.TrimSpace(handoff.Header.CardID), nil
}

func (p *Pipeline) workspaceRootForProject(projectID string) string {
	if p.cfg.ContextStore == nil || strings.TrimSpace(projectID) == "" {
		return p.cfg.WorkspaceRoot
	}
	project, err := p.cfg.ContextStore.GetProject(projectID)
	if err != nil || strings.TrimSpace(project.WorkspaceRoot) == "" {
		return p.cfg.WorkspaceRoot
	}
	return project.WorkspaceRoot
}

func mergeReplyCardID(current, fallback, next string) (string, error) {
	current = strings.TrimSpace(current)
	fallback = strings.TrimSpace(fallback)
	next = strings.TrimSpace(next)
	if next == "" || next == current {
		return current, nil
	}
	if current == "" || current == fallback {
		return next, nil
	}
	return current, fmt.Errorf("resposta referencia multiplos cards: %s e %s", current, next)
}

func (p *Pipeline) saveAgentOutput(projectID, agentName, taskID, content string) error {
	if p.cfg.ContextStore != nil {
		return p.cfg.ContextStore.SaveAgentOutputForProject(projectID, agentName, taskID, "note", content)
	}

	dir := filepath.Join(p.cfg.AgentOutputDir, strings.ToUpper(agentName))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	name := fmt.Sprintf("%s_%s.md", time.Now().Format("20060102_150405"), taskID)
	return os.WriteFile(filepath.Join(dir, name), []byte(content), 0644)
}

func (p *Pipeline) createHumanHandoff(projectID, agentName, taskID, cardID string, action agentAction) error {
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
		CardID:    cardID,
		ProjectID: projectID,
		Sender:    fmt.Sprintf("[%s]", strings.ToUpper(agentName)),
		Recipient: "[HUMAN]",
		TaskRef:   taskID,
		Intent:    "HUMAN_CLARIFICATION_REQUEST",
	}

	if _, err := hand_off.CreateHandoff(p.cfg.HandoffDir, header, payload); err != nil {
		return err
	}
	if p.cfg.ContextStore != nil && strings.TrimSpace(cardID) != "" {
		rawPayload, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		if err := p.cfg.ContextStore.AddCardComment(cardID, strings.ToUpper(agentName), "form", string(rawPayload)); err != nil {
			return err
		}
	}
	return nil
}

func (p *Pipeline) createDeliveryHandoff(projectID, taskID, cardID string, files []code_writer.CodeFile) (string, error) {
	recipient := routeDeliveryAgent(files)
	payload := map[string]any{
		"action":       "write_code",
		"instructions": "Executar a entrega solicitada pelo CEO. Responder com action write_code e conteudo bruto quando houver arquivo a persistir.",
		"files":        files,
	}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	handoff := hand_off.HandoffSchema[json.RawMessage]{
		Header: hand_off.HandoffHeader{
			CardID:    cardID,
			ProjectID: projectID,
			Sender:    "[CEO]",
			Recipient: fmt.Sprintf("[%s]", recipient),
			TaskRef:   taskID,
			Intent:    "DELIVERY_REQUEST",
		},
		Payload: rawPayload,
	}
	data, err := json.Marshal(handoff)
	if err != nil {
		return "", err
	}
	usedCardID, err := p.saveRawHandoffForCard(projectID, data, cardID)
	if err != nil {
		return "", err
	}
	p.recordDelegationComment(cardID, recipient, "DELIVERY_REQUEST")
	return usedCardID, nil
}

func shouldAutoDelegateKickoff(agentName, intent string, result agentResponseResult) bool {
	return strings.EqualFold(agentName, "CEO") &&
		strings.EqualFold(intent, "PHASE_KICKOFF") &&
		!result.Delegated &&
		!result.Blocked
}

func shouldAutoCompleteTask(agentName, intent string, result agentResponseResult) bool {
	agentName = strings.Trim(strings.ToUpper(strings.TrimSpace(agentName)), "[]")
	return agentName != "" &&
		agentName != "CEO" &&
		agentName != "HUMAN" &&
		!strings.EqualFold(intent, "TASK_COMPLETE") &&
		!result.Delegated &&
		!result.Blocked
}

func (p *Pipeline) createLeanImplementationHandoff(projectID, taskID, cardID string, sourceHandoff []byte, ceoResponse string) (string, error) {
	var source struct {
		Payload json.RawMessage `json:"payload"`
	}
	if err := json.Unmarshal(sourceHandoff, &source); err != nil {
		return "", err
	}
	payload := map[string]any{
		"action": "implement_project",
		"instructions": strings.Join([]string{
			"Implementar o projeto simples solicitado em Python.",
			"Responder obrigatoriamente com JSON valido usando action write_code ou write_files.",
			"Criar arquivos apenas dentro do diretorio do projeto.",
			"Manter solucao lean: CLI simples, persistencia JSON local, README curto e testes basicos se fizer sentido.",
			"Nao devolver apenas nota ou texto livre.",
		}, " "),
		"source_payload": source.Payload,
		"ceo_response":   ceoResponse,
	}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	handoff := hand_off.HandoffSchema[json.RawMessage]{
		Header: hand_off.HandoffHeader{
			CardID:    cardID,
			ProjectID: projectID,
			Sender:    "[CEO]",
			Recipient: "[DEV_BACKEND]",
			TaskRef:   taskID,
			Intent:    "IMPLEMENT_LEAN_BACKEND",
		},
		Payload: rawPayload,
	}
	data, err := json.Marshal(handoff)
	if err != nil {
		return "", err
	}
	usedCardID, err := p.saveRawHandoffForCard(projectID, data, cardID)
	if err != nil {
		return "", err
	}
	p.recordDelegationComment(cardID, "DEV_BACKEND", "IMPLEMENT_LEAN_BACKEND")
	return usedCardID, nil
}

func (p *Pipeline) createTaskCompleteHandoff(projectID, taskID, cardID, agentName, sourceIntent, response string) (string, error) {
	evidence := p.collectTaskEvidence(projectID, response)
	payload := map[string]any{
		"status":              "completed",
		"completed_by":        strings.ToUpper(agentName),
		"source_intent":       sourceIntent,
		"original_task_ref":   taskID,
		"agent_response":      truncateForPayload(response, 4000),
		"evidence":            evidence,
		"requires_ceo_review": true,
		"instructions": strings.Join([]string{
			"Verifique se a atividade foi realmente executada antes de aceitar.",
			"Compare agent_response com evidence.files_found e evidence.write_code_detected.",
			"Se nao houver evidencia fisica suficiente, reabra a tarefa ou gere novo handoff corretivo.",
			"Se houver entrega e faltar validacao, envie para QA.",
			"Se estiver completa, atualize plano/contexto e decida o proximo passo.",
		}, " "),
	}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	handoff := hand_off.HandoffSchema[json.RawMessage]{
		Header: hand_off.HandoffHeader{
			CardID:    cardID,
			ProjectID: projectID,
			Sender:    fmt.Sprintf("[%s]", strings.ToUpper(agentName)),
			Recipient: "[CEO]",
			TaskRef:   taskID,
			Intent:    "TASK_COMPLETE",
		},
		Payload: rawPayload,
	}
	data, err := json.Marshal(handoff)
	if err != nil {
		return "", err
	}
	usedCardID, err := p.saveRawHandoffForCard(projectID, data, cardID)
	if err != nil {
		return "", err
	}
	p.recordDelegationComment(cardID, "CEO", "TASK_COMPLETE")
	return usedCardID, nil
}

func (p *Pipeline) collectTaskEvidence(projectID, response string) map[string]any {
	projectRoot := p.workspaceRootForProject(projectID)
	files, totalFiles := listFilesForEvidence(projectRoot, 200)
	actions, declaredFiles := responseActionEvidence(response)
	return map[string]any{
		"project_root":          projectRoot,
		"files_found":           files,
		"files_found_count":     totalFiles,
		"files_found_truncated": totalFiles > len(files),
		"actions_detected":      actions,
		"declared_files":        declaredFiles,
		"write_code_detected":   containsString(actions, "write_code") || containsString(actions, "write_files") || containsString(actions, "create_file"),
	}
}

func listFilesForEvidence(root string, limit int) ([]string, int) {
	root = strings.TrimSpace(root)
	if root == "" || limit <= 0 {
		return nil, 0
	}
	var files []string
	total := 0
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry == nil {
			return nil
		}
		if entry.IsDir() {
			switch strings.ToLower(entry.Name()) {
			case ".git", ".gocache", "__pycache__", "node_modules", ".venv", "venv":
				if path != root {
					return filepath.SkipDir
				}
			}
			return nil
		}
		total++
		if len(files) >= limit {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			rel = path
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	return files, total
}

func responseActionEvidence(response string) ([]string, []string) {
	raw, err := extractJSONPayload(response)
	if err != nil {
		return nil, nil
	}
	items, err := splitAgentActions(raw)
	if err != nil {
		return nil, nil
	}
	var actions []string
	var files []string
	for _, item := range items {
		var action struct {
			Action string                 `json:"action"`
			File   *code_writer.CodeFile  `json:"file"`
			Files  []code_writer.CodeFile `json:"files"`
		}
		if err := json.Unmarshal(item, &action); err != nil {
			continue
		}
		if strings.TrimSpace(action.Action) != "" {
			actions = append(actions, strings.ToLower(strings.TrimSpace(action.Action)))
		}
		if action.File != nil && strings.TrimSpace(action.File.Path) != "" {
			files = append(files, filepath.ToSlash(action.File.Path))
		}
		for _, file := range action.Files {
			if strings.TrimSpace(file.Path) != "" {
				files = append(files, filepath.ToSlash(file.Path))
			}
		}
	}
	return uniqueStrings(actions), uniqueStrings(files)
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if strings.EqualFold(value, want) {
			return true
		}
	}
	return false
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	var unique []string
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		unique = append(unique, value)
	}
	return unique
}

func truncateForPayload(value string, limit int) string {
	value = strings.TrimSpace(value)
	if limit <= 0 || len(value) <= limit {
		return value
	}
	return value[:limit] + "\n[truncated]"
}

func routeDeliveryAgent(files []code_writer.CodeFile) string {
	for _, file := range files {
		path := strings.ToLower(filepath.ToSlash(file.Path))
		ext := strings.ToLower(filepath.Ext(path))
		switch {
		case strings.HasPrefix(path, "wiki/") || strings.HasPrefix(path, "docs/") || strings.HasPrefix(path, "doc/") || ext == ".md" || ext == ".txt":
			return "DOCUMENTATION"
		case strings.Contains(path, "frontend") || strings.Contains(path, "web/") || strings.Contains(path, "ui/") || ext == ".tsx" || ext == ".jsx" || ext == ".vue" || ext == ".css" || ext == ".scss" || ext == ".html":
			return "DEV_FRONTEND"
		case strings.Contains(path, "migration") || strings.Contains(path, "schema") || ext == ".sql":
			return "DBA"
		case strings.Contains(path, "docker") || strings.Contains(path, ".github/workflows") || strings.Contains(path, "deploy") || ext == ".yml" || ext == ".yaml" || ext == ".tf":
			return "DEVOPS"
		case strings.Contains(path, "data") || strings.Contains(path, "etl") || strings.Contains(path, "pipeline"):
			return "DATA_ENGINEER"
		case strings.Contains(path, "security") || strings.Contains(path, "auth"):
			return "SECURITY"
		}
	}
	return "DEV_BACKEND"
}

func (p *Pipeline) recordDelegationComment(cardID, recipient, intent string) {
	if p.cfg.ContextStore == nil || strings.TrimSpace(cardID) == "" {
		return
	}
	content := fmt.Sprintf(`{"recipient":"[%s]","intent":"%s"}`, recipient, intent)
	if err := p.cfg.ContextStore.AddCardComment(cardID, "SYSTEM", "delegation", content); err != nil {
		log.Printf("Registrar delegacao do card %s: %v", cardID, err)
	}
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

func (p *Pipeline) recordHandoffEvent(path string, data []byte, projectID string) error {
	var parsed struct {
		Header  hand_off.HandoffHeader `json:"header"`
		Payload json.RawMessage        `json:"payload"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}

	return p.cfg.ContextStore.SaveHandoffEvent(context_store.HandoffEvent{
		ProjectID: defaultString(parsed.Header.ProjectID, projectID),
		Path:      path,
		Sender:    parsed.Header.Sender,
		Recipient: parsed.Header.Recipient,
		TaskRef:   parsed.Header.TaskRef,
		Intent:    parsed.Header.Intent,
		Payload:   string(parsed.Payload),
		RawJSON:   string(data),
	})
}

func (p *Pipeline) ensureProjectFromHandoff(data []byte) error {
	if p.cfg.ContextStore == nil {
		return nil
	}
	var parsed struct {
		Header  hand_off.HandoffHeader `json:"header"`
		Payload struct {
			Project *context_store.Project `json:"project"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}
	projectID := strings.TrimSpace(parsed.Header.ProjectID)
	if projectID == "" || parsed.Payload.Project == nil {
		return nil
	}
	if _, err := p.cfg.ContextStore.GetProject(projectID); err == nil {
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	project := *parsed.Payload.Project
	if strings.TrimSpace(project.ID) == "" {
		project.ID = projectID
	}
	if strings.TrimSpace(project.Name) == "" {
		project.Name = project.Slug
	}
	if strings.TrimSpace(project.Name) == "" {
		return fmt.Errorf("handoff referencia projeto %s ausente, mas payload.project nao tem name/slug", projectID)
	}
	_, err := p.cfg.ContextStore.CreateProject(project)
	return err
}

func (p *Pipeline) recordCardHandoff(path string, data []byte) (string, string, error) {
	var parsed struct {
		ReplyToCardID string                 `json:"reply_to_card_id"`
		Header        hand_off.HandoffHeader `json:"header"`
		Payload       json.RawMessage        `json:"payload"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return "", "", err
	}
	if parsed.Header.CardID == "" {
		parsed.Header.CardID = strings.TrimSpace(parsed.ReplyToCardID)
	}
	hand_off.EnsureCardID(&parsed.Header)
	card, err := p.cfg.ContextStore.SaveCardHandoff(context_store.CardHandoff{
		ProjectID: parsed.Header.ProjectID,
		CardID:    parsed.Header.CardID,
		Title:     parsed.Header.Intent,
		TaskRef:   parsed.Header.TaskRef,
		Sender:    parsed.Header.Sender,
		Recipient: parsed.Header.Recipient,
		Intent:    parsed.Header.Intent,
		Payload:   string(parsed.Payload),
		RawJSON:   string(data),
	})
	if err != nil {
		return parsed.Header.CardID, parsed.Header.ProjectID, err
	}
	return card.ID, card.ProjectID, nil
}

func (p *Pipeline) recordCardResponse(cardID, agentName, response string) {
	if p.cfg.ContextStore == nil || strings.TrimSpace(cardID) == "" {
		return
	}
	if err := p.cfg.ContextStore.AddCardComment(cardID, strings.ToUpper(agentName), "response", response); err != nil {
		log.Printf("Registrar resposta do card %s: %v", cardID, err)
	}
}

func (p *Pipeline) recordCardError(cardID string, err error) {
	if p.cfg.ContextStore == nil {
		return
	}
	if strings.TrimSpace(cardID) == "" {
		return
	}
	if commentErr := p.cfg.ContextStore.AddCardComment(cardID, "SYSTEM", "error", err.Error()); commentErr != nil {
		log.Printf("Registrar erro do card %s: %v", cardID, commentErr)
	}
	p.updateCardStatus(cardID, "failed")
}

func (p *Pipeline) recordTokenUsage(cardID, agentName string, usage agent.TokenUsage) {
	if p.cfg.ContextStore == nil {
		return
	}
	if usage.AgentName == "" {
		usage.AgentName = agentName
	}
	if usage.Model == "" {
		usage.Model = mustAgentModel(agentName)
	}
	if usage.Tier == 0 {
		usage.Tier = agent.ResolverTierAgente(agentName)
	}
	if usage.TotalTokens == 0 {
		usage.TotalTokens = usage.PromptTokens + usage.CompletionTokens
	}
	err := p.cfg.ContextStore.SaveTokenUsage(context_store.TokenUsage{
		ProjectID:        projectIDForToken(p.cfg.ContextStore, cardID),
		CardID:           cardID,
		AgentName:        usage.AgentName,
		Model:            usage.Model,
		Tier:             usage.Tier,
		PromptTokens:     usage.PromptTokens,
		CompletionTokens: usage.CompletionTokens,
		TotalTokens:      usage.TotalTokens,
		LatencyMS:        usage.LatencyMS,
	})
	if err != nil {
		log.Printf("Registrar token usage do card %s: %v", cardID, err)
	}
}

func mustAgentModel(agentName string) string {
	model, err := agent.ResolverModeloAgente(agentName)
	if err != nil {
		return "unknown"
	}
	return model
}

func (p *Pipeline) retryCardHandoff(path, cardID string, err error) bool {
	if p.cfg.ContextStore == nil || strings.TrimSpace(cardID) == "" {
		writeError(p.cfg.Failed, fileID(path), filepath.Base(path), err)
		p.recordCardError(cardID, err)
		return false
	}
	retry, recordErr := p.cfg.ContextStore.RecordCardFailure(cardID, err.Error())
	if recordErr != nil {
		if errors.Is(recordErr, sql.ErrNoRows) {
			log.Printf("Card %s nao encontrado para retry; enviando handoff para failed: %v", cardID, err)
			writeError(p.cfg.Failed, fileID(path), filepath.Base(path), err)
			return false
		}
		log.Printf("Registrar retry do card %s: %v", cardID, recordErr)
		writeError(p.cfg.Failed, fileID(path), filepath.Base(path), err)
		p.recordCardError(cardID, err)
		return false
	}
	if !retry {
		writeError(p.cfg.Failed, fileID(path), filepath.Base(path), err)
		return false
	}
	card, getErr := p.cfg.ContextStore.GetCard(cardID)
	backoff := time.Second
	if getErr == nil && card.RetryCount > 0 {
		backoff = time.Duration(1<<(card.RetryCount-1)) * time.Second
	}
	time.Sleep(backoff)
	if _, moveErr := p.moveFile(path, p.cfg.HandoffDir); moveErr != nil {
		log.Printf("Reenfileirar card %s: %v", cardID, moveErr)
		writeError(p.cfg.Failed, fileID(path), filepath.Base(path), moveErr)
		return false
	}
	log.Printf("Card %s reenfileirado apos erro: %v", cardID, err)
	return true
}

func (p *Pipeline) updateCardStatus(cardID, status string) {
	if p.cfg.ContextStore == nil || strings.TrimSpace(cardID) == "" {
		return
	}
	if err := p.cfg.ContextStore.UpdateCardStatus(cardID, status); err != nil {
		log.Printf("Atualizar status do card %s: %v", cardID, err)
	}
}

func cardIDFromHandoff(data []byte) string {
	var parsed struct {
		ReplyToCardID string                 `json:"reply_to_card_id"`
		Header        hand_off.HandoffHeader `json:"header"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return ""
	}
	if strings.TrimSpace(parsed.Header.CardID) == "" {
		return strings.TrimSpace(parsed.ReplyToCardID)
	}
	return parsed.Header.CardID
}

func projectIDFromHandoff(data []byte) string {
	var parsed struct {
		Header hand_off.HandoffHeader `json:"header"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return ""
	}
	return strings.TrimSpace(parsed.Header.ProjectID)
}

func intentFromHandoff(data []byte) string {
	var parsed struct {
		Header hand_off.HandoffHeader `json:"header"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return ""
	}
	return strings.TrimSpace(parsed.Header.Intent)
}

func projectIDForToken(store *context_store.Store, cardID string) string {
	if store == nil || strings.TrimSpace(cardID) == "" {
		return ""
	}
	card, err := store.GetCard(cardID)
	if err != nil {
		return ""
	}
	return card.ProjectID
}

func (p *Pipeline) applyProjectOnboardingResponse(taskID string, data []byte) (bool, error) {
	var parsed struct {
		Header  hand_off.HandoffHeader `json:"header"`
		Payload json.RawMessage        `json:"payload"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return false, err
	}

	var probe struct {
		Action string `json:"action"`
	}
	if len(parsed.Payload) > 0 {
		if err := json.Unmarshal(parsed.Payload, &probe); err != nil {
			return false, err
		}
	}

	intent := strings.ToUpper(strings.TrimSpace(parsed.Header.Intent))
	action := strings.ToLower(strings.TrimSpace(probe.Action))
	if intent != "PROJECT_ONBOARDING_RESPONSE" && action != "create_project" {
		return false, nil
	}

	var payload projectOnboardingPayload
	if err := json.Unmarshal(parsed.Payload, &payload); err != nil {
		return true, err
	}

	project := context_store.Project{}
	if payload.Project != nil {
		project = *payload.Project
	} else {
		project = context_store.Project{
			Name:           payload.Name,
			Slug:           payload.Slug,
			Domain:         payload.Domain,
			Description:    payload.Description,
			TargetAudience: payload.TargetAudience,
			MainObjective:  payload.MainObjective,
			StackBackend:   payload.StackBackend,
			StackFrontend:  payload.StackFrontend,
			StackDatabase:  payload.StackDatabase,
			StackInfra:     payload.StackInfra,
			StackNotes:     payload.StackNotes,
		}
	}

	created, err := p.cfg.ContextStore.CreateProject(project)
	if err != nil {
		return true, err
	}
	msg := fmt.Sprintf("[PROJECT_CREATED] id=%s slug=%s", created.ID, created.Slug)
	if err := p.cfg.ContextStore.SaveAgentOutputForProject(created.ID, "SYSTEM", taskID, "project_created", msg); err != nil {
		return true, err
	}
	if err := p.createPhaseKickoff(created); err != nil {
		return true, err
	}
	return true, nil
}

func (p *Pipeline) createPhaseKickoff(project *context_store.Project) error {
	if project == nil {
		return fmt.Errorf("projeto vazio para phase kickoff")
	}
	payload := map[string]any{
		"mode":            "CONTINUE_PLAN",
		"project":         project,
		"phase_ref":       "PERSISTENCE://ROADMAP/Phase-1",
		"plan_ref":        "PERSISTENCE://PLAN/Stage-1",
		"priority":        "High",
		"complexity_gate": "Antes de delegar, classifique a complexidade como simple, medium ou complex. Para simple, use execucao lean com poucas etapas e somente agentes indispensaveis. Inclua complexity, execution_mode e recommended_agents na resposta. Nao inclua skipped_agents.",
		"instructions":    "Gere handoffs estruturados para os agentes apropriados. Para projetos simples, evite arquitetura pesada e excesso de etapas. Use update_plan para registrar a decisao de complexidade e inclua handoffs quando houver delegacao.",
	}
	header := hand_off.HandoffHeader{
		ProjectID: project.ID,
		Sender:    "[SYSTEM_INIT]",
		Recipient: "[CEO]",
		TaskRef:   "PHASE-1",
		Intent:    "PHASE_KICKOFF",
	}
	_, err := hand_off.CreateHandoff(p.cfg.HandoffDir, header, payload)
	return err
}

func (p *Pipeline) enrichHandoffWithContext(agentName string, data []byte, projectID string) (string, error) {
	var parsed struct {
		Header hand_off.HandoffHeader `json:"header"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return string(data), err
	}

	limit := contextLimitForAgent(agentName)
	queryProjectID := defaultString(parsed.Header.ProjectID, projectID)
	contextItems, err := p.cfg.ContextStore.QueryContextForHandoffForProject(queryProjectID, agentName, parsed.Header.TaskRef, limit)
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
	log.Printf("Contexto injetado: agente=%s task_ref=%s itens=%d limite=%d", agentName, parsed.Header.TaskRef, len(contextItems), limit)
	return sb.String(), nil
}

func contextLimitForAgent(agentName string) int {
	switch strings.ToUpper(strings.TrimSpace(agentName)) {
	case "CEO", "CTO", "BA", "PM", "TECH_LEAD", "CODE_REVIEWER", "SECURITY":
		return 20
	case "QA", "WRITER", "DOCUMENTATION", "ARTIST", "CMO":
		return 3
	default:
		return 5
	}
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
