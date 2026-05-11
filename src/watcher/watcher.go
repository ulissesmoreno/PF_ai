package watcher

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// Handler é chamado com o caminho completo do arquivo detectado.
type Handler func(path string)

// AgentHandler é chamado quando um arquivo de agente é detectado (extrai agente do nome).
type AgentHandler func(path string, agentName string)

// Watcher monitora uma pasta e chama Handler para arquivos novos.
type Watcher struct {
	dir    string
	handler Handler
}

// AgentWatcher monitora pasta de handoff de agentes.
type AgentWatcher struct {
	dir          string
	agentHandler AgentHandler
}

var supportedExtensions = map[string]bool{
	".pdf":  true,
	".epub": true,
	".txt":  true,
	".mp4":  true,
	".mkv":  true,
	".mov":  true,
	".m4a":  true,
	".mp3":  true,
	".json": true, 
}

var agentExtensions = map[string]bool{
	".md":    true,
	".txt":   true,
	".json":  true,
}

// New cria um Watcher para o diretório informado.
func New(dir string, handler Handler) *Watcher {
	return &Watcher{dir: dir, handler: handler}
}

// NewAgent cria um AgentWatcher para monitorar tarefas de agentes.
func NewAgent(dir string, agentHandler AgentHandler) *AgentWatcher {
	return &AgentWatcher{dir: dir, agentHandler: agentHandler}
}

// ExtractAgentName extrai o nome do agente do padrão "to_AGENTNAME" no nome do arquivo.
// Exemplo: "documento_to_CODE_REVIEWER.md" → "CODE_REVIEWER"
// Retorna vazio string se o padrão não for encontrado.
func ExtractAgentName(filename string) string {
	// Remove extensão para análise
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	
	// Procura por padrão "to_NOMEDO" (case-insensitive)
	re := regexp.MustCompile(`(?i)to_([a-zA-Z0-9_]+)`)
	matches := re.FindStringSubmatch(base)
	if len(matches) > 1 {
		// Retorna em UPPERCASE para padronização (ex: CODE_REVIEWER)
		return strings.ToUpper(matches[1])
	}
	
	return ""
}

// ScanInitial varre a pasta recursivamente e processa arquivos já existentes.
func (w *Watcher) ScanInitial() error {
	log.Printf("🔍 Varrendo pasta inicial: %s", w.dir)
	found := 0
	return filepath.Walk(w.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("⚠️  Erro ao varrer %s: %v", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		log.Printf("📋 Arquivo encontrado: %s (ext: %s)", filepath.Base(path), ext)

		if !supportedExtensions[ext] {
			log.Printf("⏭️  Extensão não suportada: %s", ext)
			return nil
		}

		// Ignorar arquivos temporários e .error.txt
		base := filepath.Base(path)
		if strings.HasPrefix(base, ".") || strings.HasSuffix(base, ".error.txt") {
			log.Printf("⏭️  Ignorando temporário/erro: %s", base)
			return nil
		}

		found++
		log.Printf("🔍 %d arquivo(s) encontrado(s) em %s", found, w.dir)
		log.Printf("📥 Processando inicial: %s", base)
		go w.handler(path) // não bloqueia o scan

		return nil
	})
}

// Start inicia o monitoramento. Bloqueia até ctx ser cancelado.
func (w *Watcher) Start(ctx context.Context) error {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer fw.Close()

	if err := fw.Add(w.dir); err != nil {
		return err
	}

	log.Printf("👁️  Monitorando: %s", w.dir)

	for {
		select {
		case event, ok := <-fw.Events:
			if !ok {
				return nil
			}
			if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {
				path := event.Name
				ext := strings.ToLower(filepath.Ext(path))

				if !supportedExtensions[ext] {
					continue
				}

				// Ignorar arquivos temporários e .error.txt
				base := filepath.Base(path)
				if strings.HasPrefix(base, ".") || strings.HasSuffix(base, ".error.txt") {
					continue
				}

				log.Printf("📥 Detectado: %s", base)
				go w.handler(path) // não bloqueia o watcher

			}

		case err, ok := <-fw.Errors:
			if !ok {
				return nil
			}
			log.Printf("⚠️  Watcher erro (%s): %v", w.dir, err)

		case <-ctx.Done():
			return nil
		}
	}
}

// ────────────────────────────────────────────────────────────────────
// AgentWatcher — Monitoramento de tarefas de agentes
// ────────────────────────────────────────────────────────────────────

// ScanInitial varre a pasta de agentes e processa arquivos já existentes.
func (aw *AgentWatcher) ScanInitial() error {
	log.Printf("🔍 Varrendo pasta de agentes inicial: %s", aw.dir)
	found := 0
	return filepath.Walk(aw.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("⚠️  Erro ao varrer %s: %v", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}

		filename := filepath.Base(path)
		ext := strings.ToLower(filepath.Ext(path))

		// Verificar extensão
		if !agentExtensions[ext] {
			return nil
		}

		// Ignorar arquivos temporários e hidden files
		if strings.HasPrefix(filename, ".") {
			log.Printf("⏭️  Ignorando arquivo temporário: %s", filename)
			return nil
		}

		// Extrair nome do agente do padrão to_AGENTNAME
		agentName := ExtractAgentName(filename)
		if agentName == "" {
			log.Printf("⏭️  Arquivo não segue padrão to_AGENTNAME: %s", filename)
			return nil
		}

		found++
		log.Printf("📋 Tarefa de agente encontrada: %s → %s", filename, agentName)
		go aw.agentHandler(path, agentName) // não bloqueia o scan

		return nil
	})
}

// Start inicia o monitoramento de tarefas de agentes. Bloqueia até ctx ser cancelado.
func (aw *AgentWatcher) Start(ctx context.Context) error {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer fw.Close()

	if err := fw.Add(aw.dir); err != nil {
		return err
	}

	log.Printf("👁️  Monitorando agentes: %s", aw.dir)

	for {
		select {
		case event, ok := <-fw.Events:
			if !ok {
				return nil
			}
			if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {
				path := event.Name
				filename := filepath.Base(path)
				ext := strings.ToLower(filepath.Ext(path))

				// Verificar extensão
				if !agentExtensions[ext] {
					continue
				}

				// Ignorar arquivos temporários
				if strings.HasPrefix(filename, ".") {
					continue
				}

				// Extrair nome do agente
				agentName := ExtractAgentName(filename)
				if agentName == "" {
					continue
				}

				log.Printf("📥 Tarefa detectada: %s → %s", filename, agentName)
				go aw.agentHandler(path, agentName) // não bloqueia o watcher

			}

		case err, ok := <-fw.Errors:
			if !ok {
				return nil
			}
			log.Printf("⚠️  AgentWatcher erro (%s): %v", aw.dir, err)

		case <-ctx.Done():
			return nil
		}
	}
}
