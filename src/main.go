package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pf_ai/agent"
	"pf_ai/code_writer"
	"pf_ai/embeddings"
	"pf_ai/hand_off"
	"pf_ai/pipeline"
	"pf_ai/watcher"

	"github.com/joho/godotenv"
)

func createCEOKickoff(cfg Config) error {
	header := hand_off.HandoffHeader{
		Sender:    "[SYSTEM_INIT]",
		Recipient: "[CEO]",
		TaskRef:   "PHASE-1",
		Intent:    "PHASE_KICKOFF",
	}

	payload := hand_off.PhaseKickoffPayload{
		PhaseRef:         "DOC/ROADMAP.md#Phase-1",
		PlanRef:          "DOC/PLAN.md#Stage-1",
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

	path, err := hand_off.CreateHandoff(cfg.AgentHandoffDir, header, payload)
	if err != nil {
		return err
	}

	log.Printf("Handoff de kickoff gerado: %s", path)
	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Println(".env nao encontrado; usando variaveis do sistema.")
		}
	}

	cfg := loadConfig()
	cfg.log()

	dirs := []string{
		cfg.InboxRaw,
		cfg.AgentHandoffDir,
		cfg.Processing,
		cfg.PendingPython,
		cfg.Extracted,
		cfg.ProcessedPython,
		cfg.Success,
		cfg.Failed,
		cfg.VaultPath,
		cfg.AgentsConfigDir,
		cfg.AgentOutputDir,
		"data",
		"logs",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			log.Fatalf("Nao foi possivel criar pasta %s: %v", d, err)
		}
	}

	code_writer.Init(cfg.VaultPath)
	agent.SetConfigDir(cfg.AgentsConfigDir)
	embeddings.Init(cfg.OllamaURL, cfg.OllamaEmbedModel)
	agent.Init(cfg.OllamaURL, cfg.OllamaLLMModel)

	log.Println("Inicializando agente CEO...")
	if err := agent.InitComAgente(cfg.OllamaURL, "CEO"); err != nil {
		log.Printf("Aviso ao inicializar CEO: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

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
		HandoffDir:       cfg.AgentHandoffDir,
		AgentOutputDir:   cfg.AgentOutputDir,
		WorkspaceRoot:    cfg.WorkspaceRoot,
	})
	pl.Start(ctx)

	rawWatcher := watcher.New(cfg.InboxRaw, func(path string) {
		pl.Submit(pipeline.Job{Path: path, Source: pipeline.SourceRaw})
	})

	log.Println("Agent iniciado; aguardando arquivos em:", cfg.InboxRaw)
	if err := rawWatcher.ScanInitial(); err != nil {
		log.Printf("Erro ao varrer inbox/raw: %v", err)
	}

	go func() {
		if err := rawWatcher.Start(ctx); err != nil {
			log.Printf("Watcher raw: %v", err)
			cancel()
		}
	}()

	agentWatcher := watcher.NewAgent(cfg.AgentHandoffDir, func(path string, agentName string) {
		log.Printf("Tarefa recebida: agente=%s arquivo=%s", agentName, path)
		pl.Submit(pipeline.Job{Path: path, Source: pipeline.SourceHandoff, AgentName: agentName})
	})

	log.Printf("Monitorando agentes em: %s", cfg.AgentHandoffDir)

	if err := agentWatcher.ScanInitial(); err != nil {
		log.Printf("Erro ao varrer agentes: %v", err)
	}

	go func() {
		if err := agentWatcher.Start(ctx); err != nil {
			log.Printf("AgentWatcher: %v", err)
			cancel()
		}
	}()

	time.Sleep(250 * time.Millisecond)
	log.Println("Chamando CEO...")
	if err := createCEOKickoff(cfg); err != nil {
		log.Printf("Erro ao criar CEO kickoff: %v", err)
	}

	<-ctx.Done()
	log.Println("Encerrando; aguardando jobs em andamento...")
	pl.Wait()
	log.Println("Agent encerrado")
}
