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

type Handler func(path string)

type AgentHandler func(path string, agentName string)

type Watcher struct {
	dir     string
	handler Handler
}

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
	".md":   true,
	".txt":  true,
	".json": true,
}

var knownAgents = []string{
	"DATA_ENGINEER",
	"CODE_REVIEWER",
	"DEV_FRONTEND",
	"DEV_BACKEND",
	"DOCUMENTATION",
	"UX_RESEARCHER",
	"TECH_LEAD",
	"SECURITY",
	"ARTIST",
	"WRITER",
	"DEVOPS",
	"DS_ML",
	"CEO",
	"CTO",
	"DBA",
	"CMO",
	"PM",
	"BA",
	"QA",
}

func New(dir string, handler Handler) *Watcher {
	return &Watcher{dir: dir, handler: handler}
}

func NewAgent(dir string, agentHandler AgentHandler) *AgentWatcher {
	return &AgentWatcher{dir: dir, agentHandler: agentHandler}
}

func ExtractAgentName(filename string) string {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	upper := strings.ToUpper(base)

	idx := strings.Index(upper, "_TO_")
	if idx == -1 {
		idx = strings.Index(upper, "TO_")
		if idx == -1 {
			return ""
		}
	} else {
		idx++
	}
	recipientPart := upper[idx+len("TO_"):]
	if recipientPart == "HUMAN" || strings.HasPrefix(recipientPart, "HUMAN_") {
		return ""
	}

	for _, agent := range knownAgents {
		if recipientPart == agent || strings.HasPrefix(recipientPart, agent+"_") {
			return agent
		}
	}

	re := regexp.MustCompile(`(?i)^([a-zA-Z0-9]+(?:_[a-zA-Z0-9]+)?)`)
	matches := re.FindStringSubmatch(recipientPart)
	if len(matches) > 1 {
		return strings.ToUpper(matches[1])
	}

	return ""
}

func isInternalAgentDir(path string) bool {
	switch strings.ToLower(filepath.Base(path)) {
	case "raw", "processing", "pending_python", "extracted", "processed_python", "success", "failed":
		return true
	default:
		return false
	}
}

func (w *Watcher) ScanInitial() error {
	log.Printf("Varrendo pasta inicial: %s", w.dir)
	found := 0
	return filepath.Walk(w.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Erro ao varrer %s: %v", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		log.Printf("Arquivo encontrado: %s (ext: %s)", filepath.Base(path), ext)

		if !supportedExtensions[ext] {
			log.Printf("Extensao nao suportada: %s", ext)
			return nil
		}

		base := filepath.Base(path)
		if strings.HasPrefix(base, ".") || strings.HasSuffix(base, ".error.txt") {
			log.Printf("Ignorando temporario/erro: %s", base)
			return nil
		}

		found++
		log.Printf("%d arquivo(s) encontrado(s) em %s", found, w.dir)
		log.Printf("Processando inicial: %s", base)
		go w.handler(path)

		return nil
	})
}

func (w *Watcher) Start(ctx context.Context) error {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer fw.Close()

	if err := fw.Add(w.dir); err != nil {
		return err
	}

	log.Printf("Monitorando: %s", w.dir)

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

				base := filepath.Base(path)
				if strings.HasPrefix(base, ".") || strings.HasSuffix(base, ".error.txt") {
					continue
				}

				log.Printf("Detectado: %s", base)
				go w.handler(path)
			}

		case err, ok := <-fw.Errors:
			if !ok {
				return nil
			}
			log.Printf("Watcher erro (%s): %v", w.dir, err)

		case <-ctx.Done():
			return nil
		}
	}
}

func (aw *AgentWatcher) ScanInitial() error {
	log.Printf("Varrendo pasta de agentes inicial: %s", aw.dir)
	found := 0
	return filepath.Walk(aw.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Erro ao varrer %s: %v", path, err)
			return err
		}
		if info.IsDir() {
			if path != aw.dir && isInternalAgentDir(path) {
				return filepath.SkipDir
			}
			return nil
		}

		filename := filepath.Base(path)
		ext := strings.ToLower(filepath.Ext(path))

		if !agentExtensions[ext] {
			return nil
		}

		if strings.HasPrefix(filename, ".") {
			log.Printf("Ignorando arquivo temporario: %s", filename)
			return nil
		}

		agentName := ExtractAgentName(filename)
		if agentName == "" {
			log.Printf("Arquivo nao segue padrao to_AGENTNAME: %s", filename)
			return nil
		}

		found++
		log.Printf("Tarefa de agente encontrada: %s -> %s", filename, agentName)
		go aw.agentHandler(path, agentName)

		return nil
	})
}

func (aw *AgentWatcher) Start(ctx context.Context) error {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer fw.Close()

	if err := fw.Add(aw.dir); err != nil {
		return err
	}

	log.Printf("Monitorando agentes: %s", aw.dir)

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

				if !agentExtensions[ext] {
					continue
				}

				if strings.HasPrefix(filename, ".") {
					continue
				}

				agentName := ExtractAgentName(filename)
				if agentName == "" {
					continue
				}

				log.Printf("Tarefa detectada: %s -> %s", filename, agentName)
				go aw.agentHandler(path, agentName)
			}

		case err, ok := <-fw.Errors:
			if !ok {
				return nil
			}
			log.Printf("AgentWatcher erro (%s): %v", aw.dir, err)

		case <-ctx.Done():
			return nil
		}
	}
}
