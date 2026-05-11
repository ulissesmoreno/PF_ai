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
	"pf_ai/embeddings"
	"pf_ai/extractor"
)

// Source indica a origem do job.
type Source int

const (
	SourceRaw       Source = iota // arquivo novo em inbox/raw/
	SourceExtracted                // .chunks.json pronto do Python
)

// Job representa um arquivo a processar.
type Job struct {
	Path   string
	Source Source
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
}

// Pipeline gerencia o worker pool.
type Pipeline struct {
	cfg     Config
	jobs    chan Job
	wg      sync.WaitGroup
	// pending mapeia id → canal para notificar quando chunks chegarem
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
		log.Printf("⚠️  Fila cheia — esperando: %s", filepath.Base(job.Path))
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

// ── Processamento de arquivo novo ─────────────────────────────────────────────

func (p *Pipeline) handleRaw(ctx context.Context, path string) {
	ext := strings.ToLower(filepath.Ext(path)) // alterar para padrão de nome
	id := fileID(path)
	nome := filepath.Base(path)

	log.Printf("⚙️  [%s] Iniciando", id)

	// Mover para processing/
	procPath, err := moveFile(path, p.cfg.Processing)
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
		moveFile(procPath, p.cfg.Failed) //nolint
		return
	}

	if err != nil {
		log.Printf("❌ [%s] Extração: %v", id, err)
		writeError(p.cfg.Failed, id, nome, err)
		moveFile(procPath, p.cfg.Failed) //nolint
		return
	}

	// Continuar pipeline com os chunks
	if err := p.indexAndGenerate(ctx, id, nome, chunks); err != nil {
		log.Printf("❌ [%s] Pipeline: %v", id, err)
		writeError(p.cfg.Failed, id, nome, err)
		moveFile(procPath, p.cfg.Failed) //nolint
		return
	}

	moveFile(procPath, p.cfg.Success) //nolint
	log.Printf("✅ [%s] Concluído", id)
}

// ── Extração nativa (EPUB, TXT) ───────────────────────────────────────────────

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
			PageApprox: 0, // TXT e EPUB não têm páginas
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
	metaPath := filepath.Join(p.cfg.PendingPython, id+".json")
	metaBytes, _ := json.MarshalIndent(meta, "", "  ")
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

	// Fan-out: embeddings em paralelo
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
	// var vsChunks []vectorstore.Chunk
	// for _, e := range embedded {
	// 	vsChunks = append(vsChunks, vectorstore.Chunk{
	// 		ID:         e.ID,
	// 		Source:     nome,
	// 		Text:       e.Text,
	// 		PageApprox: e.PageApprox,
	// 		Embedding:  e.Embedding,
	// 	})
	// }
	// if err := vectorstore.Add(vsChunks); err != nil {
	// 	return fmt.Errorf("vectorstore: %w", err)
	// }

	// TODO: Buscar contexto para a nota
	// queryEmb, err := embeddings.Generate("resumo conceitos principais " + id)
	// if err != nil {
	// 	return fmt.Errorf("embedding query: %w", err)
	// }
	// topChunks, err := vectorstore.Query(queryEmb, 5)
	// if err != nil {
	// 	return fmt.Errorf("query vectorstore: %w", err)
	// }

	var textos []string
	// for _, c := range topChunks {
	// 	textos = append(textos, c.Text)
	// }

	// Gerar nota
	log.Printf("✍️ [%s] Recebendo resposta...", id)
	nota, err := agent.GerarArquivo(id, textos)
	if err != nil {
		return fmt.Errorf("gerar nota: %w", err)
	}

	// Salvar no vault
	return code_writer.SalvarNota(id, nome, nota)
}

// ── Helpers ───────────────────────────────────────────────────────────────────

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

func moveFile(src, destDir string) (string, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}
	dest := filepath.Join(destDir, filepath.Base(src))
	// Evitar colisão
	if _, err := os.Stat(dest); err == nil {
		ts := time.Now().Format("20060102_150405")
		ext := filepath.Ext(dest)
		stem := strings.TrimSuffix(filepath.Base(dest), ext)
		dest = filepath.Join(destDir, fmt.Sprintf("%s_%s%s", stem, ts, ext))
	}
	return dest, os.Rename(src, dest)
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
			// Subdividir por sentença
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
	// Divide em frases por . ! ?
	var sentencas []string
	current := ""
	for _, r := range texto {
		current += string(r)
		if r == '.' || r == '!' || r == '?' {
			s := strings.TrimSpace(current)
			if s != "" {
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
