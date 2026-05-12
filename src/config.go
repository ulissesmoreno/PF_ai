package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	// ── Pastas do pipeline legado ─────────────────────────────
	InboxRaw        string
	Processing      string
	PendingPython   string
	Extracted       string
	ProcessedPython string
	Success         string
	Failed          string
	VaultPath       string
	DBPath          string
	WorkspaceRoot   string
	MigrationsDir   string

	// ── Pastas do sistema de agentes ───────────────────────────
	AgentHandoffDir string // Entrada para processamento de agentes
	AgentsConfigDir string // Diretório com configurações de agentes
	AgentOutputDir  string // Diretório base para saída de agentes

	// ── Workers ────────────────────────────────────────────────
	WorkerCount      int
	EmbedWorkerCount int

	// ── Ollama ─────────────────────────────────────────────────
	OllamaURL        string
	OllamaLLMModel   string
	OllamaEmbedModel string
	OllamaGemmaModel string
	CEOProvider      string
	CodexCEOCLI      string

	// ── Timeouts ───────────────────────────────────────────────
	PythonTimeout   time.Duration
	AgentTimeout    time.Duration
	CodexCEOTimeout time.Duration
}

func loadConfig() Config {
	workspaceRoot := env("WORKSPACE_ROOT", detectWorkspaceRoot())
	handoffDir := env("AGENT_HANDOFF_DIR", filepathJoin(workspaceRoot, ".agent_handoff"))
	return Config{
		InboxRaw:        env("INBOX_RAW_DIR", filepathJoin(handoffDir, "raw")),
		Processing:      env("PROCESSING_DIR", filepathJoin(handoffDir, "processing")),
		PendingPython:   env("PENDING_PYTHON_DIR", filepathJoin(handoffDir, "pending_python")),
		Extracted:       env("EXTRACTED_DIR", filepathJoin(handoffDir, "extracted")),
		ProcessedPython: env("PROCESSED_PYTHON_DIR", filepathJoin(handoffDir, "processed_python")),
		Success:         env("SUCCESS_DIR", filepathJoin(handoffDir, "success")),
		Failed:          env("FAILED_DIR", filepathJoin(handoffDir, "failed")),
		VaultPath:       workspacePath(workspaceRoot, env("VAULT_PATH", "vault")),
		DBPath:          workspacePath(workspaceRoot, env("DB_PATH", filepathJoin("data", "wiki.db"))),
		WorkspaceRoot:   workspaceRoot,
		MigrationsDir:   workspacePath(workspaceRoot, env("MIGRATIONS_DIR", filepathJoin("db", "migrations"))),

		AgentHandoffDir: handoffDir,
		AgentsConfigDir: workspacePath(workspaceRoot, env("AGENTS_CONFIG_DIR", "AGENTS")),
		AgentOutputDir:  workspacePath(workspaceRoot, env("AGENT_OUTPUT_DIR", "output")),

		WorkerCount:      envInt("WORKER_COUNT", 3),
		EmbedWorkerCount: envInt("EMBED_WORKER_COUNT", 5),

		OllamaURL:        env("OLLAMA_URL", "http://localhost:11434"),
		OllamaLLMModel:   env("OLLAMA_LLM_MODEL", "deepseek-r1:7b"),
		OllamaEmbedModel: env("OLLAMA_EMBED_MODEL", "nomic-embed-text"),
		OllamaGemmaModel: env("OLLAMA_GEMMA_MODEL", "gemma4:latest"),
		CEOProvider:      env("CEO_PROVIDER", "ollama"),
		CodexCEOCLI:      env("CODEX_CEO_CLI", "codex exec -"),

		PythonTimeout:   envDuration("PYTHON_TIMEOUT_MINUTES", 60) * time.Minute,
		AgentTimeout:    envDuration("AGENT_TIMEOUT_SECONDS", 60) * time.Second,
		CodexCEOTimeout: envDuration("CODEX_CEO_TIMEOUT_SECONDS", 600) * time.Second,
	}
}

func (c Config) log() {
	log.Println("── Configuração ──────────────────────────────")
	log.Printf("  .agent_handoff/raw       : %s", c.InboxRaw)
	log.Printf("  vault           : %s", c.VaultPath)
	log.Printf("  banco           : %s", c.DBPath)
	log.Printf("  workspace       : %s", c.WorkspaceRoot)
	log.Printf("  migrations      : %s", c.MigrationsDir)
	log.Println("──")
	log.Printf("  agent handoff   : %s", c.AgentHandoffDir)
	log.Printf("  agents config   : %s", c.AgentsConfigDir)
	log.Printf("  agent output    : %s", c.AgentOutputDir)
	log.Println("──")
	log.Printf("  ollama url      : %s", c.OllamaURL)
	log.Printf("  modelo llm      : %s", c.OllamaLLMModel)
	log.Printf("  modelo embed    : %s", c.OllamaEmbedModel)
	log.Printf("  modelo gemma4   : %s", c.OllamaGemmaModel)
	log.Printf("  ceo provider    : %s", c.CEOProvider)
	if strings.EqualFold(c.CEOProvider, "codex_cli") {
		log.Printf("  codex ceo cli   : %s", c.CodexCEOCLI)
	}
	log.Println("──")
	log.Printf("  workers         : %d arquivo(s) paralelos", c.WorkerCount)
	log.Printf("  embed workers   : %d chunk(s) paralelos", c.EmbedWorkerCount)
	log.Printf("  python timeout  : %s", c.PythonTimeout)
	log.Printf("  agent timeout   : %s", c.AgentTimeout)
	log.Printf("  codex ceo timeout: %s", c.CodexCEOTimeout)
	log.Println("──────────────────────────────────────────────")
}

// Helpers

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
		log.Printf("⚠️  %s inválido, usando padrão %d", key, fallback)
	}
	return fallback
}

func envDuration(key string, fallbackMinutes int) time.Duration {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return time.Duration(n)
		}
		log.Printf("⚠️  %s inválido, usando padrão %d minutos", key, fallbackMinutes)
	}
	return time.Duration(fallbackMinutes)
}

func detectWorkspaceRoot() string {
	if dirExists("AGENTS") && dirExists("DOC") {
		return "."
	}
	if dirExists("../AGENTS") && dirExists("../DOC") {
		return ".."
	}
	return "."
}

func firstExisting(paths ...string) string {
	for _, path := range paths {
		if dirExists(path) {
			return path
		}
	}
	if len(paths) == 0 {
		return "."
	}
	return paths[0]
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func filepathJoin(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for _, part := range parts[1:] {
		if result == "" || result == "." {
			result = part
			continue
		}
		result += string(os.PathSeparator) + part
	}
	return result
}

func workspacePath(workspaceRoot, value string) string {
	if value == "" || filepath.IsAbs(value) {
		return value
	}
	clean := filepath.Clean(value)
	if strings.HasPrefix(clean, ".."+string(os.PathSeparator)) || clean == ".." {
		return clean
	}
	return filepath.Join(workspaceRoot, clean)
}
