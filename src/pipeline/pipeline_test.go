package pipeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"pf_ai/code_writer"
	"pf_ai/context_store"
	"pf_ai/hand_off"
)

func TestDetectHandoffRecipient(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "CEO_TO_DEV_BACKEND_TASK.json")
	data, err := json.Marshal(hand_off.HandoffSchema[map[string]string]{
		Header: hand_off.HandoffHeader{
			Sender:    "[CEO]",
			Recipient: "[DEV_BACKEND]",
			Intent:    "IMPLEMENTATION_REQUEST",
		},
		Payload: map[string]string{"task": "build"},
	})
	if err != nil {
		t.Fatalf("marshal handoff: %v", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write handoff: %v", err)
	}

	agentName, ok := detectHandoffRecipient(path)
	if !ok {
		t.Fatal("expected handoff recipient to be detected")
	}
	if agentName != "DEV_BACKEND" {
		t.Fatalf("agentName = %q, want DEV_BACKEND", agentName)
	}
}

func TestMoveFileConcurrentSameSource(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "handoff.json")
	dest := filepath.Join(dir, "processing")
	if err := os.WriteFile(src, []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}

	p := New(Config{})
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := p.moveFile(src, dest)
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)

	successes := 0
	notExist := 0
	for err := range errs {
		if err == nil {
			successes++
			continue
		}
		if os.IsNotExist(err) {
			notExist++
			continue
		}
		t.Fatalf("unexpected error: %v", err)
	}
	if successes != 1 || notExist != 1 {
		t.Fatalf("successes=%d notExist=%d, want 1/1", successes, notExist)
	}
}

func TestMoveFileConcurrentNameCollision(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "processing")
	srcA := filepath.Join(dir, "a", "handoff.json")
	srcB := filepath.Join(dir, "b", "handoff.json")
	if err := os.MkdirAll(filepath.Dir(srcA), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(srcB), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(srcA, []byte("a"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(srcB, []byte("b"), 0644); err != nil {
		t.Fatal(err)
	}

	p := New(Config{})
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for _, src := range []string{srcA, srcB} {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			_, err := p.moveFile(path, dest)
			errs <- err
		}(src)
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("moveFile returned error: %v", err)
		}
	}
	entries, err := os.ReadDir(dest)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("moved files = %d, want 2", len(entries))
	}
}

func TestApplyAgentResponseWritesCodeFile(t *testing.T) {
	workspace := t.TempDir()
	p := New(Config{WorkspaceRoot: workspace})

	response := `{
		"action": "write_code",
		"files": [
			{"path": "src/sum_code.py", "content": "number_a = 10\nnumber_b = 20\nprint(number_a + number_b)\n"}
		]
	}`

	if _, err := p.applyAgentResponse("DEV_BACKEND", "TASK-1", "", response); err != nil {
		t.Fatalf("applyAgentResponse returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(workspace, "src", "sum_code.py"))
	if err != nil {
		t.Fatalf("expected code file: %v", err)
	}
	if string(data) != "number_a = 10\nnumber_b = 20\nprint(number_a + number_b)" {
		t.Fatalf("unexpected content: %q", string(data))
	}
}

func TestApplyAgentResponseWritesCodeFileInsideProjectWorkspace(t *testing.T) {
	workspace := t.TempDir()
	store := openPipelineTestStore(t, workspace)
	defer store.Close()
	project, err := store.CreateProject(context_store.Project{Name: "Mini CRM", Slug: "mini-crm"})
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}

	p := New(Config{WorkspaceRoot: workspace, ContextStore: store})
	response := `{
		"action": "write_code",
		"files": [
			{"path": "src/app.go", "content": "package main\n"}
		]
	}`

	if _, err := p.applyAgentResponseForProject(project.ID, "DEV_BACKEND", "TASK-1", "", response); err != nil {
		t.Fatalf("applyAgentResponseForProject returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(project.WorkspaceRoot, "src", "app.go")); err != nil {
		t.Fatalf("expected code file in project workspace: %v", err)
	}
	if _, err := os.Stat(filepath.Join(workspace, "src", "app.go")); !os.IsNotExist(err) {
		t.Fatalf("unexpected file in root workspace, stat err=%v", err)
	}
}

func TestApplyAgentResponseUsesReplyToCardIDForHumanHandoff(t *testing.T) {
	workspace := t.TempDir()
	store := openPipelineTestStore(t, workspace)
	defer store.Close()
	if _, err := store.SaveCardHandoff(context_store.CardHandoff{
		CardID:    "card-1",
		Title:     "Need clarification",
		TaskRef:   "TASK-1",
		Sender:    "[CEO]",
		Recipient: "[DEV_BACKEND]",
		Intent:    "IMPLEMENT",
		RawJSON:   `{}`,
	}); err != nil {
		t.Fatal(err)
	}

	handoffDir := filepath.Join(workspace, ".agent_handoff")
	p := New(Config{WorkspaceRoot: workspace, HandoffDir: handoffDir, ContextStore: store})
	response := `{
		"action": "ask_human",
		"reply_to_card_id": "card-1",
		"question": "Qual endpoint?",
		"blocking": true
	}`

	result, err := p.applyAgentResponse("DEV_BACKEND", "TASK-1", "", response)
	if err != nil {
		t.Fatalf("applyAgentResponse: %v", err)
	}
	if !result.Blocked || result.CardID != "card-1" {
		t.Fatalf("result = %#v", result)
	}

	thread, err := store.GetCardThread("card-1")
	if err != nil {
		t.Fatal(err)
	}
	if thread.Card.Status != "blocked" || len(thread.Comments) != 2 {
		t.Fatalf("thread = %#v", thread)
	}
	handoff := readOnlyPipelineHandoff(t, handoffDir)
	if handoff.Header.CardID != "card-1" {
		t.Fatalf("card_id = %q, want card-1", handoff.Header.CardID)
	}
}

func TestApplyAgentResponseInjectsReplyToCardIDIntoOutgoingHandoff(t *testing.T) {
	workspace := t.TempDir()
	handoffDir := filepath.Join(workspace, ".agent_handoff")
	p := New(Config{WorkspaceRoot: workspace, HandoffDir: handoffDir})
	response := `{
		"action": "handoff",
		"reply_to_card_id": "card-1",
		"handoff": {
			"header": {
				"sender": "[DEV_BACKEND]",
				"recipient": "[QA]",
				"task_ref": "TASK-1",
				"intent": "TEST_REQUEST"
			},
			"payload": {"target": "api"}
		}
	}`

	result, err := p.applyAgentResponse("DEV_BACKEND", "TASK-1", "", response)
	if err != nil {
		t.Fatalf("applyAgentResponse: %v", err)
	}
	if result.CardID != "card-1" {
		t.Fatalf("cardID = %q, want card-1", result.CardID)
	}
	handoff := readOnlyPipelineHandoff(t, handoffDir)
	if handoff.Header.CardID != "card-1" {
		t.Fatalf("card_id = %q, want card-1", handoff.Header.CardID)
	}
}

func TestApplyAgentResponseCarriesProjectIDIntoOutgoingHandoff(t *testing.T) {
	workspace := t.TempDir()
	handoffDir := filepath.Join(workspace, ".agent_handoff")
	p := New(Config{WorkspaceRoot: workspace, HandoffDir: handoffDir})
	response := `{
		"action": "handoff",
		"reply_to_card_id": "card-1",
		"handoff": {
			"header": {
				"sender": "[DEV_BACKEND]",
				"recipient": "[QA]",
				"task_ref": "TASK-1",
				"intent": "TEST_REQUEST"
			},
			"payload": {"target": "api"}
		}
	}`

	if _, err := p.applyAgentResponseForProject("project-a", "DEV_BACKEND", "TASK-1", "", response); err != nil {
		t.Fatalf("applyAgentResponseForProject: %v", err)
	}
	handoff := readOnlyPipelineHandoff(t, handoffDir)
	if handoff.Header.ProjectID != "project-a" {
		t.Fatalf("project_id = %q, want project-a", handoff.Header.ProjectID)
	}
}

func TestRecordCardHandoffUsesHeaderProjectIDOverActiveProject(t *testing.T) {
	workspace := t.TempDir()
	store := openPipelineTestStore(t, workspace)
	defer store.Close()

	projectA, err := store.CreateProject(context_store.Project{Name: "Project A", Slug: "project-a"})
	if err != nil {
		t.Fatalf("CreateProject A: %v", err)
	}
	projectB, err := store.CreateProject(context_store.Project{Name: "Project B", Slug: "project-b"})
	if err != nil {
		t.Fatalf("CreateProject B: %v", err)
	}
	store.ActiveProjectID = projectB.ID

	p := New(Config{ContextStore: store, HandoffDir: filepath.Join(workspace, ".agent_handoff")})
	data, err := json.Marshal(hand_off.HandoffSchema[map[string]string]{
		Header: hand_off.HandoffHeader{
			CardID:    "card-a",
			ProjectID: projectA.ID,
			Sender:    "[CEO]",
			Recipient: "[DEV_BACKEND]",
			TaskRef:   "TASK-1",
			Intent:    "IMPLEMENT",
		},
		Payload: map[string]string{"task": "build"},
	})
	if err != nil {
		t.Fatalf("marshal handoff: %v", err)
	}

	cardID, projectID, err := p.recordCardHandoff(filepath.Join(workspace, "handoff.json"), data)
	if err != nil {
		t.Fatalf("recordCardHandoff: %v", err)
	}
	if cardID != "card-a" || projectID != projectA.ID {
		t.Fatalf("cardID=%q projectID=%q", cardID, projectID)
	}
	card, err := store.GetCard("card-a")
	if err != nil {
		t.Fatalf("GetCard: %v", err)
	}
	if card.ProjectID != projectA.ID {
		t.Fatalf("card project_id = %q, want %s", card.ProjectID, projectA.ID)
	}
}

func TestCEOCodeActionIsDelegatedToBackend(t *testing.T) {
	workspace := t.TempDir()
	handoffDir := filepath.Join(workspace, ".agent_handoff")
	p := New(Config{WorkspaceRoot: workspace, HandoffDir: handoffDir})
	response := `{
		"action": "write_code",
		"reply_to_card_id": "card-1",
		"files": [{"path": "src/ceo_should_not_write.py", "content": "print(5)"}]
	}`

	result, err := p.applyAgentResponse("CEO", "TASK-1", "", response)
	if err != nil {
		t.Fatalf("applyAgentResponse: %v", err)
	}
	if result.CardID != "card-1" {
		t.Fatalf("cardID = %q, want card-1", result.CardID)
	}
	if _, err := os.Stat(filepath.Join(workspace, "src", "ceo_should_not_write.py")); !os.IsNotExist(err) {
		t.Fatalf("CEO wrote code directly, stat err=%v", err)
	}
	handoff := readOnlyPipelineHandoff(t, handoffDir)
	if handoff.Header.Sender != "[CEO]" || handoff.Header.Recipient != "[DEV_BACKEND]" {
		t.Fatalf("handoff = %#v", handoff.Header)
	}
}

func TestCEODocumentActionIsDelegatedToDocumentation(t *testing.T) {
	workspace := t.TempDir()
	handoffDir := filepath.Join(workspace, ".agent_handoff")
	p := New(Config{WorkspaceRoot: workspace, HandoffDir: handoffDir})
	response := `{
		"action": "write_code",
		"reply_to_card_id": "card-1",
		"files": [{"path": "wiki/project-overview.md", "content": "# Overview"}]
	}`

	result, err := p.applyAgentResponse("CEO", "TASK-1", "", response)
	if err != nil {
		t.Fatalf("applyAgentResponse: %v", err)
	}
	if !result.Delegated || result.CardID != "card-1" {
		t.Fatalf("result = %#v", result)
	}
	handoff := readOnlyPipelineHandoff(t, handoffDir)
	if handoff.Header.Recipient != "[DOCUMENTATION]" {
		t.Fatalf("recipient = %q, want [DOCUMENTATION]", handoff.Header.Recipient)
	}
}

func TestCEONoteWithFilesIsDelegated(t *testing.T) {
	workspace := t.TempDir()
	handoffDir := filepath.Join(workspace, ".agent_handoff")
	p := New(Config{WorkspaceRoot: workspace, HandoffDir: handoffDir})
	response := `{
		"action": "note",
		"reply_to_card_id": "card-1",
		"files": [{"path": "frontend/App.tsx", "content": "export default function App(){return null}"}]
	}`

	result, err := p.applyAgentResponse("CEO", "TASK-1", "", response)
	if err != nil {
		t.Fatalf("applyAgentResponse: %v", err)
	}
	if !result.Delegated {
		t.Fatalf("result = %#v", result)
	}
	if _, err := os.Stat(filepath.Join(workspace, "frontend", "App.tsx")); !os.IsNotExist(err) {
		t.Fatalf("CEO wrote frontend directly, stat err=%v", err)
	}
	handoff := readOnlyPipelineHandoff(t, handoffDir)
	if handoff.Header.Recipient != "[DEV_FRONTEND]" {
		t.Fatalf("recipient = %q, want [DEV_FRONTEND]", handoff.Header.Recipient)
	}
}

func TestDetectHandoffRecipientAcceptsUTF8BOM(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "CEO_TO_DEV_BACKEND_TEST.json")
	data := append([]byte{0xEF, 0xBB, 0xBF}, []byte(`{"header":{"sender":"[CEO]","recipient":"[DEV_BACKEND]","intent":"TEST"},"payload":{}}`)...)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	recipient, ok := detectHandoffRecipient(path)
	if !ok || recipient != "DEV_BACKEND" {
		t.Fatalf("recipient=%q ok=%v, want DEV_BACKEND/true", recipient, ok)
	}
}

func TestApplyProjectOnboardingResponseCreatesProject(t *testing.T) {
	workspace := t.TempDir()
	store := openPipelineTestStore(t, workspace)
	defer store.Close()

	p := New(Config{ContextStore: store, HandoffDir: filepath.Join(workspace, ".agent_handoff")})
	data, err := json.Marshal(hand_off.HandoffSchema[map[string]string]{
		Header: hand_off.HandoffHeader{
			Sender:    "[HUMAN]",
			Recipient: "[CEO]",
			TaskRef:   "ONBOARDING-1",
			Intent:    "PROJECT_ONBOARDING_RESPONSE",
		},
		Payload: map[string]string{
			"name":            "PF AI",
			"slug":            "PF AI",
			"description":     "Agent orchestration",
			"target_audience": "developers",
		},
	})
	if err != nil {
		t.Fatalf("marshal handoff: %v", err)
	}

	handled, err := p.applyProjectOnboardingResponse("ONBOARDING-1", data)
	if err != nil {
		t.Fatalf("applyProjectOnboardingResponse: %v", err)
	}
	if !handled {
		t.Fatal("handled = false, want true")
	}

	project, err := store.GetActiveProject()
	if err != nil {
		t.Fatalf("GetActiveProject: %v", err)
	}
	if project.Slug != "pf-ai" {
		t.Fatalf("slug = %q, want pf-ai", project.Slug)
	}
	handoff := readOnlyPipelineHandoff(t, filepath.Join(workspace, ".agent_handoff"))
	if handoff.Header.ProjectID != project.ID || handoff.Header.Intent != "PHASE_KICKOFF" {
		t.Fatalf("kickoff header = %#v, want project %s PHASE_KICKOFF", handoff.Header, project.ID)
	}
	if !strings.Contains(string(handoff.Payload), "complexity_gate") {
		t.Fatalf("kickoff payload missing complexity gate: %s", string(handoff.Payload))
	}
}

func TestApplyProjectOnboardingResponseIgnoresPhaseKickoffProjectString(t *testing.T) {
	workspace := t.TempDir()
	store := openPipelineTestStore(t, workspace)
	defer store.Close()

	p := New(Config{ContextStore: store})
	data := []byte(`{
		"header": {
			"sender": "[SYSTEM_INIT]",
			"recipient": "[CEO]",
			"task_ref": "CRM-001",
			"intent": "PHASE_KICKOFF"
		},
		"payload": {
			"project": "Mini CRM pessoal",
			"goal": "Criar uma API local em Go com SQLite para contatos e oportunidades."
		}
	}`)

	handled, err := p.applyProjectOnboardingResponse("CRM-001", data)
	if err != nil {
		t.Fatalf("applyProjectOnboardingResponse: %v", err)
	}
	if handled {
		t.Fatal("handled = true, want false for PHASE_KICKOFF")
	}
}

func TestNormalizeCodeFilesStripsSingleFence(t *testing.T) {
	files, err := normalizeCodeFiles([]code_writer.CodeFile{
		{Path: "src/sum_code.py", Content: "```python\nprint(30)\n```"},
	})
	if err != nil {
		t.Fatalf("normalizeCodeFiles returned error: %v", err)
	}
	if files[0].Content != "print(30)" {
		t.Fatalf("content = %q, want print(30)", files[0].Content)
	}
}

func TestNormalizeCodeFilesRejectsMarkdownOutsideWiki(t *testing.T) {
	_, err := normalizeCodeFiles([]code_writer.CodeFile{
		{Path: "DOC/result.md", Content: "# result"},
	})
	if err == nil {
		t.Fatal("expected markdown outside wiki to be rejected")
	}
}

func TestNormalizeCodeFilesAllowsMarkdownInVault(t *testing.T) {
	files, err := normalizeCodeFiles([]code_writer.CodeFile{
		{Path: "vault/raw-summary.md", Content: "# Raw summary"},
	})
	if err != nil {
		t.Fatalf("normalizeCodeFiles returned error: %v", err)
	}
	if filepath.ToSlash(files[0].Path) != "vault/raw-summary.md" {
		t.Fatalf("path = %q, want vault/raw-summary.md", files[0].Path)
	}
}

func TestApplyAgentResponseRejectsTextResponse(t *testing.T) {
	p := New(Config{WorkspaceRoot: t.TempDir()})

	_, err := p.applyAgentResponse("DEV_BACKEND", "TASK-1", "", "Claro, aqui esta o codigo:\n```python\nprint(30)\n```")
	if err == nil {
		t.Fatal("expected text response to be rejected")
	}
	if !strings.Contains(err.Error(), "sem JSON") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApplyAgentResponseAcceptsJSONAfterTextPrefix(t *testing.T) {
	p := New(Config{WorkspaceRoot: t.TempDir()})

	response := `Resposta final:
{
  "action": "note",
  "message": "onboarding needs human input"
}`

	if _, err := p.applyAgentResponse("CEO", "ONBOARDING-1", "", response); err != nil {
		t.Fatalf("applyAgentResponse returned error: %v", err)
	}
}

func TestApplyAgentResponseIgnoresBracketedTextBeforeJSON(t *testing.T) {
	p := New(Config{WorkspaceRoot: t.TempDir()})

	response := `[RASCUNHO]
{
  "action": "note",
  "message": "onboarding needs human input"
}`

	if _, err := p.applyAgentResponse("CEO", "ONBOARDING-1", "", response); err != nil {
		t.Fatalf("applyAgentResponse returned error: %v", err)
	}
}

func TestApplyAgentResponseRejectsNoOpJSONResponse(t *testing.T) {
	p := New(Config{WorkspaceRoot: t.TempDir()})

	_, err := p.applyAgentResponse("CEO", "ONBOARDING-1", "", `{
		"name": "Sistema de Gestao Financeira",
		"slug": "finance-ai"
	}`)
	if err == nil {
		t.Fatal("expected no-op JSON response to be rejected")
	}
	if !strings.Contains(err.Error(), "sem action executavel") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApplyAgentResponseUpdatePlanCanDelegateHandoffs(t *testing.T) {
	workspace := t.TempDir()
	handoffDir := filepath.Join(workspace, ".agent_handoff")
	p := New(Config{WorkspaceRoot: workspace, HandoffDir: handoffDir})
	response := `{
		"action": "update_plan",
		"item_type": "plan",
		"title": "Complexity decision",
		"content": "complexity=simple execution_mode=lean",
		"complexity": "simple",
		"execution_mode": "lean",
		"recommended_agents": ["DEV_BACKEND", "QA"],
		"handoffs": [
			{
				"header": {
					"sender": "[CEO]",
					"recipient": "[DEV_BACKEND]",
					"task_ref": "TASK-1",
					"intent": "IMPLEMENT_LEAN_BACKEND"
				},
				"payload": {"task": "build minimal API"}
			}
		]
	}`

	result, err := p.applyAgentResponseForProject("project-1", "CEO", "PHASE-1", "card-1", response)
	if err != nil {
		t.Fatalf("applyAgentResponseForProject: %v", err)
	}
	if !result.Delegated {
		t.Fatalf("result = %#v, want delegated", result)
	}
	handoff := readOnlyPipelineHandoff(t, handoffDir)
	if handoff.Header.ProjectID != "project-1" || handoff.Header.Recipient != "[DEV_BACKEND]" {
		t.Fatalf("handoff header = %#v", handoff.Header)
	}
}

func TestApplyAgentResponseInfersUpdatePlanFromComplexity(t *testing.T) {
	workspace := t.TempDir()
	store := openPipelineTestStore(t, workspace)
	defer store.Close()
	project, err := store.CreateProject(context_store.Project{Name: "FIFO", Slug: "fifo"})
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	p := New(Config{WorkspaceRoot: workspace, ContextStore: store})

	response := `{
		"complexity": "simple",
		"execution_mode": "lean",
		"recommended_agents": ["DEV_BACKEND", "QA"]
	}`

	if _, err := p.applyAgentResponseForProject(project.ID, "CEO", "PHASE-1", "card-1", response); err != nil {
		t.Fatalf("applyAgentResponseForProject: %v", err)
	}
	items, err := store.QueryContextForHandoffForProject(project.ID, "CEO", "PHASE-1", 5)
	if err != nil {
		t.Fatalf("QueryContextForHandoffForProject: %v", err)
	}
	if len(items) != 1 || !strings.Contains(items[0], "complexity=simple") {
		t.Fatalf("planning items = %#v", items)
	}
}

func TestEnsureProjectFromHandoffRestoresMissingProject(t *testing.T) {
	workspace := t.TempDir()
	store := openPipelineTestStore(t, workspace)
	defer store.Close()
	p := New(Config{WorkspaceRoot: workspace, ContextStore: store})

	data := []byte(`{
		"header": {
			"card_id": "card-1",
			"project_id": "project-restore",
			"sender": "[SYSTEM_INIT]",
			"recipient": "[CEO]",
			"task_ref": "PHASE-1",
			"intent": "PHASE_KICKOFF"
		},
		"payload": {
			"project": {
				"id": "project-restore",
				"name": "FIFO de Produtos",
				"slug": "fifo-produtos-python",
				"status": "active"
			}
		}
	}`)

	if err := p.ensureProjectFromHandoff(data); err != nil {
		t.Fatalf("ensureProjectFromHandoff: %v", err)
	}
	if _, _, err := p.recordCardHandoff(filepath.Join(workspace, "handoff.json"), data); err != nil {
		t.Fatalf("recordCardHandoff: %v", err)
	}
	card, err := store.GetCard("card-1")
	if err != nil {
		t.Fatalf("GetCard: %v", err)
	}
	if card.ProjectID != "project-restore" {
		t.Fatalf("card project_id = %q", card.ProjectID)
	}
}

func TestCreateLeanImplementationHandoffFromKickoffNote(t *testing.T) {
	workspace := t.TempDir()
	handoffDir := filepath.Join(workspace, ".agent_handoff")
	p := New(Config{WorkspaceRoot: workspace, HandoffDir: handoffDir})
	source := []byte(`{
		"header": {
			"card_id": "card-1",
			"project_id": "project-1",
			"sender": "[SYSTEM_INIT]",
			"recipient": "[CEO]",
			"task_ref": "PHASE-1",
			"intent": "PHASE_KICKOFF"
		},
		"payload": {
			"project": {"name": "FIFO de Produtos", "slug": "fifo-produtos-python"}
		}
	}`)

	usedCardID, err := p.createLeanImplementationHandoff("project-1", "PHASE-1", "card-1", source, "You are assigned to set up the backend environment.")
	if err != nil {
		t.Fatalf("createLeanImplementationHandoff: %v", err)
	}
	if usedCardID != "card-1" {
		t.Fatalf("usedCardID = %q", usedCardID)
	}
	handoff := readOnlyPipelineHandoff(t, handoffDir)
	if handoff.Header.ProjectID != "project-1" || handoff.Header.Recipient != "[DEV_BACKEND]" || handoff.Header.Intent != "IMPLEMENT_LEAN_BACKEND" {
		t.Fatalf("handoff header = %#v", handoff.Header)
	}
	if !strings.Contains(string(handoff.Payload), "write_code") {
		t.Fatalf("handoff payload missing write_code instruction: %s", string(handoff.Payload))
	}
}

func TestCreateTaskCompleteHandoffIncludesEvidence(t *testing.T) {
	workspace := t.TempDir()
	handoffDir := filepath.Join(workspace, ".agent_handoff")
	store := openPipelineTestStore(t, workspace)
	defer store.Close()
	project, err := store.CreateProject(context_store.Project{Name: "FIFO", Slug: "fifo"})
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	if err := os.WriteFile(filepath.Join(project.WorkspaceRoot, "src", "main.py"), []byte("print('ok')\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	p := New(Config{WorkspaceRoot: workspace, HandoffDir: handoffDir, ContextStore: store})
	response := `{
		"action": "write_code",
		"files": [{"path": "src/main.py", "content": "print('ok')\n"}]
	}`

	usedCardID, err := p.createTaskCompleteHandoff(project.ID, "TASK-1", "card-1", "DEV_BACKEND", "IMPLEMENT_LEAN_BACKEND", response)
	if err != nil {
		t.Fatalf("createTaskCompleteHandoff: %v", err)
	}
	if usedCardID != "card-1" {
		t.Fatalf("usedCardID = %q", usedCardID)
	}
	handoff := readOnlyPipelineHandoff(t, handoffDir)
	if handoff.Header.Sender != "[DEV_BACKEND]" || handoff.Header.Recipient != "[CEO]" || handoff.Header.Intent != "TASK_COMPLETE" {
		t.Fatalf("handoff header = %#v", handoff.Header)
	}
	if !strings.Contains(string(handoff.Payload), `"write_code_detected": true`) {
		t.Fatalf("payload missing write_code evidence: %s", string(handoff.Payload))
	}
	if !strings.Contains(string(handoff.Payload), "src/main.py") {
		t.Fatalf("payload missing physical file evidence: %s", string(handoff.Payload))
	}
}

func TestShouldAutoCompleteTaskAvoidsLoops(t *testing.T) {
	if !shouldAutoCompleteTask("DEV_BACKEND", "IMPLEMENT_LEAN_BACKEND", agentResponseResult{}) {
		t.Fatal("DEV_BACKEND completion should auto-complete")
	}
	if shouldAutoCompleteTask("CEO", "TASK_COMPLETE", agentResponseResult{}) {
		t.Fatal("CEO should not auto-complete")
	}
	if shouldAutoCompleteTask("DEV_BACKEND", "TASK_COMPLETE", agentResponseResult{}) {
		t.Fatal("TASK_COMPLETE should not generate another TASK_COMPLETE")
	}
	if shouldAutoCompleteTask("DEV_BACKEND", "IMPLEMENT_LEAN_BACKEND", agentResponseResult{Delegated: true}) {
		t.Fatal("delegated result should not auto-complete")
	}
	if shouldAutoCompleteTask("DEV_BACKEND", "IMPLEMENT_LEAN_BACKEND", agentResponseResult{Blocked: true}) {
		t.Fatal("blocked result should not auto-complete")
	}
}

func openPipelineTestStore(t *testing.T, workspace string) *context_store.Store {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	migrationsDir := filepath.Clean(filepath.Join(wd, "..", "..", "db", "migrations"))
	store, err := context_store.Open(filepath.Join(workspace, "data", "test.db"), migrationsDir, workspace)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	return store
}

func readOnlyPipelineHandoff(t *testing.T, dir string) hand_off.HandoffSchema[json.RawMessage] {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("handoffs = %d, want 1", len(entries))
	}
	data, err := os.ReadFile(filepath.Join(dir, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	var handoff hand_off.HandoffSchema[json.RawMessage]
	if err := json.Unmarshal(data, &handoff); err != nil {
		t.Fatal(err)
	}
	return handoff
}
