package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"pf_ai/agent"
	"pf_ai/code_writer"
	"pf_ai/context_store"
	"pf_ai/embeddings"
	"pf_ai/hand_off"
	"pf_ai/pipeline"
	"pf_ai/watcher"

	"github.com/joho/godotenv"
)

type ceoKickoffPayload struct {
	Mode                         string   `json:"mode"`
	PhaseRef                     string   `json:"phase_ref"`
	PlanRef                      string   `json:"plan_ref"`
	SeniorityLevel               string   `json:"seniority_level"`
	LLMTier                      int      `json:"llm_tier"`
	SecurityCriteria             string   `json:"security_criteria"`
	ApplicableSkills             []string `json:"applicable_skills"`
	Constraints                  []string `json:"constraints"`
	Priority                     string   `json:"priority"`
	ExplicitConclusionAuthorized bool     `json:"explicit_conclusion_authorized"`
	AuthorizedUntilStage         *string  `json:"authorized_until_stage"`
	Instructions                 string   `json:"instructions"`
	PersistedContext             []string `json:"persisted_context"`
}

func createCEOKickoff(cfg Config, store *context_store.Store) error {
	projectStarted, err := store.ProjectStarted()
	if err != nil {
		return err
	}

	header := hand_off.HandoffHeader{
		Sender:    "[SYSTEM_INIT]",
		Recipient: "[CEO]",
		TaskRef:   "PHASE-1",
		Intent:    "PHASE_KICKOFF",
	}

	payload := ceoKickoffPayload{
		Mode:             "CONTINUE_PLAN",
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

	if projectStarted {
		payload.Instructions = "Projeto iniciado detectado na persistência. Continue o plano usando o contexto persistido. Gere handoffs estruturados para os agentes apropriados. Atualizações operacionais devem usar ações CQRS como update_context, update_plan e record_test. Wiki deve ser criada/atualizada como arquivo em wiki/. QUESTIONS.md continua append-only para perguntas ao humano."
		payload.PersistedContext, err = store.QueryContextForHandoff("CEO", "", 20)
		if err != nil {
			return err
		}
	} else {
		payload.Mode = "PROJECT_IGNITION"
		header.TaskRef = "ONBOARDING-1"
		payload.PhaseRef = "PERSISTENCE://PROJECT"
		payload.PlanRef = "PERSISTENCE://PLAN"
		payload.Instructions, err = buildInitialPromptFromREADME(cfg.WorkspaceRoot)
		if err != nil {
			return err
		}
		payload.PersistedContext, err = store.QueryContextForHandoff("CEO", "", 12)
		if err != nil {
			return err
		}
	}

	path, err := hand_off.CreateHandoff(cfg.AgentHandoffDir, header, payload)
	if err != nil {
		return err
	}

	log.Printf("Handoff de kickoff gerado: %s", path)
	return nil
}

func buildInitialPromptFromREADME(workspaceRoot string) (string, error) {
	data, err := os.ReadFile(filepath.Join(workspaceRoot, "README.md"))
	if err != nil {
		return "", fmt.Errorf("ler README para ignition: %w", err)
	}

	readme := string(data)
	prompt := extractReadmeIgnitionPrompt(readme)
	if prompt == "" {
		prompt = readme
	}

	adaptation := `

[ADAPTAÇÃO PARA ORQUESTRAÇÃO COM PERSISTÊNCIA]
- Não edite DOC/*.md diretamente para contexto, plano, estado, testes, decisões ou retrospectiva.
- Use ações CQRS append-only: update_context, update_plan, update_state, record_test, record_decision e record_retrospective.
- Use write_code apenas para código e arquivos reais exigidos pelo filesystem, incluindo wiki/*.md para Obsidian.
- QUESTIONS.md permanece arquivo append-only para perguntas ao humano; NEW-INSTRUCTIONS.md permanece humano-only.
- Monte próximos handoffs usando o contexto persistido recebido em persisted_context.
- Se o PROJECT persistido ainda estiver em template/placeholders, conduza onboarding pelo CEO antes de delegar implementação.
`

	return strings.TrimSpace(prompt) + adaptation, nil
}

func extractReadmeIgnitionPrompt(readme string) string {
	marker := "Prompt Inicial para IA"
	idx := strings.Index(readme, marker)
	if idx == -1 {
		idx = strings.Index(readme, "Initial AI Prompt")
	}
	if idx == -1 {
		return ""
	}

	rest := readme[idx:]
	start := strings.Index(rest, "```text")
	if start == -1 {
		return ""
	}
	rest = rest[start+len("```text"):]
	end := strings.Index(rest, "```")
	if end == -1 {
		return ""
	}
	return strings.TrimSpace(rest[:end])
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
		filepathJoin(cfg.WorkspaceRoot, "data"),
		filepathJoin(cfg.WorkspaceRoot, "logs"),
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

	contextStore, err := context_store.Open(cfg.DBPath, cfg.MigrationsDir, cfg.WorkspaceRoot)
	if err != nil {
		log.Fatalf("Context store: %v", err)
	}
	defer contextStore.Close()

	if err := contextStore.ImportOperationalDocuments(); err != nil {
		log.Fatalf("Importar documentos operacionais: %v", err)
	}

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
		ContextStore:     contextStore,
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
	if err := createCEOKickoff(cfg, contextStore); err != nil {
		log.Printf("Erro ao criar CEO kickoff: %v", err)
	}

	<-ctx.Done()
	log.Println("Encerrando; aguardando jobs em andamento...")
	pl.Wait()
	log.Println("Agent encerrado")
}
