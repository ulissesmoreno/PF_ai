package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"pf_ai/agent"
	"pf_ai/code_writer"
	"pf_ai/embeddings"
	"pf_ai/pipeline"
	"pf_ai/watcher"	
	"pf_ai/hand_off"

	"github.com/joho/godotenv"
)

func createCEOKickoff(cfg config) error {
	// 1. Definição do Header conforme §4.2.1
	header := hand_off.HandoffHeader{
		Sender:    "[SYSTEM_INIT]",
		Recipient: "[CEO]",
		TaskRef:   "PHASE-1",
		Intent:    "PHASE_KICKOFF",
	}

	// 2. Definição do Payload conforme §4.2.2
	// Nota: Certifique-se de que esta struct esteja visível ou definida no pacote agent
	payload := hand_off.PhaseKickoffPayload{
		PhaseRef:         "ROADMAP.md#Phase-1",
		PlanRef:          "PLAN.md#Stage-1",
		SeniorityLevel:   "Senior",
		LLMTier:          3,
		SecurityCriteria: "Initial threat model: JWT isolation and input sanitization.",
		ApplicableSkills: []string{"Orchestration", "Task_Delegation", "Project_Management"},
		Constraints: []string{
			"TDD mandatory",
			"Hexagonal isolation",
			"GSD-RULES strict adherence",
		},
		Priority:                     "High",
		ExplicitConclusionAuthorized: false,
		AuthorizedUntilStage:         nil,
	}

	// 3. Chamada ao criador genérico de hand_off.go
	path, err := hand_off.CreateHandoff(cfg.AgentHandoffDir, header, payload)
	if err != nil {
		return err
	}

	log.Printf("🚀 Handoff de Kickoff gerado com sucesso: %s", path)
	return nil
}

func main() {
	// Carregar .env (ignora se não existir — usa variáveis do sistema)
	// Tenta carregar da pasta atual, se falhar tenta na pasta anterior
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("⚠️  .env não encontrado nas pastas local ou superior.")
		}
	}

	cfg := loadConfig()
	cfg.log()

	// Garantir pastas
	dirs := []string{
		cfg.InboxRaw,
		cfg.Processing,
		cfg.PendingPython,
		cfg.Extracted,
		cfg.ProcessedPython,
		cfg.Success,
		cfg.Failed,
		cfg.VaultPath,
		cfg.AgentHandoffDir,
		cfg.AgentsConfigDir,
		cfg.AgentOutputDir,
		"data",
		"logs",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			log.Fatalf("❌ Não foi possível criar pasta %s: %v", d, err)
		}
	}

	// // Inicializar dependências
	// if err := vectorstore.Init(cfg.DBPath); err != nil {
	// 	log.Fatalf("❌ SQLite: %v", err)
	// }
	// defer vectorstore.Close()

	code_writer.Init(cfg.VaultPath)
	embeddings.Init(cfg.OllamaURL, cfg.OllamaEmbedModel)
	agent.Init(cfg.OllamaURL, cfg.OllamaLLMModel)

	// Inicializar agente CEO para orquestração
	log.Println("🎯 Inicializando agente CEO...")
	if err := agent.InitComAgente(cfg.OllamaURL, "CEO"); err != nil {
		log.Printf("⚠️  Aviso ao inicializar CEO: %v", err)
	}

	// Contexto com cancelamento via Ctrl+C
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Iniciar pipeline
	pl := pipeline.New(pipeline.Config{
		WorkerCount:      cfg.WorkerCount,
		EmbedWorkerCount: cfg.EmbedWorkerCount,
		PendingPython:    cfg.PendingPython,
		Extracted:        cfg.Extracted,
		Processing:       cfg.Processing,
		ProcessedPython:  cfg.ProcessedPython,
		Success:          cfg.Success,
		Failed:           cfg.Failed,
		PythonTimeout:    cfg.PythonTimeout,
	})
	pl.Start(ctx)

	// Watcher em inbox/raw/
	handOffWatcher := watcher.New(cfg.InboxRaw, func(path string) {
		pl.Submit(pipeline.Job{Path: path, Source: pipeline.SourceRaw})
	})	
	
	log.Println("🚀 Agent iniciado — aguardando arquivos em: ", cfg.InboxRaw)

	// Processar arquivos já existentes antes de iniciar monitoramento
	if err := handOffWatcher.ScanInitial(); err != nil {
		log.Printf("⚠️  Erro ao varrer inbox/raw: %v", err)
	}

	go func() {
		if err := handOffWatcher.Start(ctx); err != nil {
			log.Printf("❌ Watcher raw: %v", err)
			cancel()
		}
	}()

	// ────────────────────────────────────────────────────────────────────
	// Watcher para tarefas de agentes em .agent_handoff/
	// ────────────────────────────────────────────────────────────────────
	agentWatcher := watcher.NewAgent(cfg.AgentHandoffDir, func(path string, agentName string) {
		log.Printf("🤖 Tarefa recebida: agente=%s arquivo=%s", agentName, path)
		// TODO: Processar tarefa de agente através do pipeline de agentes
		// Por enquanto apenas registra a detecção
	})

	log.Printf("🚀 Monitorando agentes em: %s", cfg.AgentHandoffDir)

	// Processar tarefas de agentes já existentes
	if err := agentWatcher.ScanInitial(); err != nil {
		log.Printf("⚠️  Erro ao varrer agentes: %v", err)
	}

	go func() {
		if err := agentWatcher.Start(ctx); err != nil {
			log.Printf("❌ AgentWatcher: %v", err)
			cancel()
		}
	}()

	// ────────────────────────────────────────────────────────────────────
	// Kickoff: Enviar PHASE_KICKOFF para o CEO
	// ────────────────────────────────────────────────────────────────────
	log.Println("📋 Chamando CEO...")
	if err := createCEOKickoff(cfg); err != nil {
		log.Printf("⚠️  Erro ao criar CEO kickoff: %v", err)
	}

	<-ctx.Done()
	log.Println("🛑 Encerrando — aguardando jobs em andamento...")
	pl.Wait()
	log.Println("✅ Agent encerrado")
}
