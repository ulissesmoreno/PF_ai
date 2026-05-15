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

	p := New(Config{ContextStore: store})
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
