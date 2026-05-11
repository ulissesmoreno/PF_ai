package main

import (
	"log"
	"os"
	"strconv"
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

	// ── Timeouts ───────────────────────────────────────────────
	PythonTimeout  time.Duration
	AgentTimeout   time.Duration
}

func loadConfig() Config {
	return Config{
		InboxRaw:        env("INBOX_RAW_DIR",        "../inbox/raw"),
		Processing:      env("PROCESSING_DIR",       "../inbox/processing"),
		PendingPython:   env("PENDING_PYTHON_DIR",   "../inbox/pending_python"),
		Extracted:       env("EXTRACTED_DIR",        "../inbox/extracted"),
		ProcessedPython: env("PROCESSED_PYTHON_DIR", "../inbox/processed_python"),
		Success:         env("SUCCESS_DIR",          "../inbox/success"),
		Failed:          env("FAILED_DIR",           "../inbox/failed"),
		VaultPath:       env("VAULT_PATH",           "vault"),
		DBPath:          env("DB_PATH",              "data/wiki.db"),

		AgentHandoffDir: env("AGENT_HANDOFF_DIR",    ".agent_handoff"),
		AgentsConfigDir: env("AGENTS_CONFIG_DIR",    "AGENTS"),
		AgentOutputDir:  env("AGENT_OUTPUT_DIR",     "output"),

		WorkerCount:      envInt("WORKER_COUNT",       3),
		EmbedWorkerCount: envInt("EMBED_WORKER_COUNT", 5),

		OllamaURL:        env("OLLAMA_URL",           "http://localhost:11434"),
		OllamaLLMModel:   env("OLLAMA_LLM_MODEL",    "deepseek-r1:7b"),
		OllamaEmbedModel: env("OLLAMA_EMBED_MODEL",   "nomic-embed-text"),
		OllamaGemmaModel: env("OLLAMA_GEMMA_MODEL",   "gemma4:latest"),

		PythonTimeout: envDuration("PYTHON_TIMEOUT_MINUTES", 60) * time.Minute,
		AgentTimeout:  envDuration("AGENT_TIMEOUT_SECONDS", 60) * time.Minute,
	}
}

func (c Config) log() {
	log.Println("── Configuração ──────────────────────────────")
	log.Printf("  inbox/raw       : %s", c.InboxRaw)
	log.Printf("  vault           : %s", c.VaultPath)
	log.Printf("  banco           : %s", c.DBPath)
	log.Println("──")
	log.Printf("  agent handoff   : %s", c.AgentHandoffDir)
	log.Printf("  agents config   : %s", c.AgentsConfigDir)
	log.Printf("  agent output    : %s", c.AgentOutputDir)
	log.Println("──")
	log.Printf("  ollama url      : %s", c.OllamaURL)
	log.Printf("  modelo llm      : %s", c.OllamaLLMModel)
	log.Printf("  modelo embed    : %s", c.OllamaEmbedModel)
	log.Printf("  modelo gemma4   : %s", c.OllamaGemmaModel)
	log.Println("──")
	log.Printf("  workers         : %d arquivo(s) paralelos", c.WorkerCount)
	log.Printf("  embed workers   : %d chunk(s) paralelos", c.EmbedWorkerCount)
	log.Printf("  python timeout  : %s", c.PythonTimeout)
	log.Printf("  agent timeout   : %s", c.AgentTimeout)
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
