package context_store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImportDocumentIsAppendOnly(t *testing.T) {
	workspace := t.TempDir()
	docPath := filepath.Join(workspace, "DOC")
	if err := os.MkdirAll(docPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(docPath, "PLAN.md"), []byte("# Plan\n\nfirst"), 0644); err != nil {
		t.Fatal(err)
	}

	store := openTestStore(t, workspace)
	defer store.Close()

	seed := DocumentSeed{Path: "DOC/PLAN.md", DocumentType: "PLAN", Owner: "BA/CTO"}
	if err := store.ImportDocument(seed); err != nil {
		t.Fatalf("first import: %v", err)
	}
	if err := os.WriteFile(filepath.Join(docPath, "PLAN.md"), []byte("# Plan\n\nsecond"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := store.ImportDocument(seed); err != nil {
		t.Fatalf("second import: %v", err)
	}

	var imports int
	if err := store.db.QueryRow("SELECT COUNT(*) FROM document_imports").Scan(&imports); err != nil {
		t.Fatal(err)
	}
	if imports != 2 {
		t.Fatalf("imports = %d, want 2", imports)
	}

	var sections int
	if err := store.db.QueryRow("SELECT COUNT(*) FROM document_sections").Scan(&sections); err != nil {
		t.Fatal(err)
	}
	if sections < 4 {
		t.Fatalf("sections = %d, want append-only history from both imports", sections)
	}
}

func TestImportEntryDocumentStoresEntriesOnly(t *testing.T) {
	workspace := t.TempDir()
	docPath := filepath.Join(workspace, "DOC")
	if err := os.MkdirAll(docPath, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# STATE.md\n\n## Done\n\n- **[2026-05-12 10:00] - [CEO]:** first entry\n  - Ref: A\n\n- **[2026-05-12 10:01] - [CTO]:** second entry"
	if err := os.WriteFile(filepath.Join(docPath, "STATE.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	store := openTestStore(t, workspace)
	defer store.Close()

	if err := store.ImportDocument(DocumentSeed{Path: "DOC/STATE.md", DocumentType: "STATE", Owner: "Technical agents"}); err != nil {
		t.Fatalf("ImportDocument: %v", err)
	}

	var sections int
	if err := store.db.QueryRow("SELECT COUNT(*) FROM document_sections").Scan(&sections); err != nil {
		t.Fatal(err)
	}
	if sections != 0 {
		t.Fatalf("sections = %d, want 0 for entry document", sections)
	}

	var entries int
	if err := store.db.QueryRow("SELECT COUNT(*) FROM document_entries").Scan(&entries); err != nil {
		t.Fatal(err)
	}
	if entries != 2 {
		t.Fatalf("entries = %d, want 2", entries)
	}
}

func TestImportTasksCreatesPlanningItems(t *testing.T) {
	workspace := t.TempDir()
	docPath := filepath.Join(workspace, "DOC")
	if err := os.MkdirAll(docPath, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# TASKS\n\n## Active\n\n### Task [TASK-42]: Build API\n- **Status:** Doing\n- **Assigned to:** DEV_BACKEND:Senior\n- **Priority:** High\n- **Deadline:** 2026-05-14\n"
	if err := os.WriteFile(filepath.Join(docPath, "TASKS.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	store := openTestStore(t, workspace)
	defer store.Close()

	if err := store.ImportDocument(DocumentSeed{Path: "DOC/TASKS.md", DocumentType: "TASKS", Owner: "Technical agents"}); err != nil {
		t.Fatalf("ImportDocument: %v", err)
	}

	var count int
	var status string
	if err := store.db.QueryRow("SELECT COUNT(*), MAX(status) FROM planning_items WHERE task_ref = 'TASK-42'").Scan(&count, &status); err != nil {
		t.Fatal(err)
	}
	if count != 1 || status != "Doing" {
		t.Fatalf("count=%d status=%q, want 1/Doing", count, status)
	}
}

func TestImportRoadmapCreatesPlanningItems(t *testing.T) {
	workspace := t.TempDir()
	docPath := filepath.Join(workspace, "DOC")
	if err := os.MkdirAll(docPath, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# ROADMAP\n\n## Phase 1\n- [ ] **Build API:** ship endpoints\n- [x] **Setup:** done\n"
	if err := os.WriteFile(filepath.Join(docPath, "ROADMAP.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	store := openTestStore(t, workspace)
	defer store.Close()

	if err := store.ImportDocument(DocumentSeed{Path: "DOC/ROADMAP.md", DocumentType: "ROADMAP", Owner: "CEO/PM"}); err != nil {
		t.Fatalf("ImportDocument: %v", err)
	}

	var count int
	if err := store.db.QueryRow("SELECT COUNT(*) FROM planning_items WHERE item_type = 'roadmap_item'").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("roadmap planning items = %d, want 2", count)
	}
}

func TestImportPlanCreatesContextEntries(t *testing.T) {
	workspace := t.TempDir()
	docPath := filepath.Join(workspace, "DOC")
	if err := os.MkdirAll(docPath, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# PLAN\n\n## 1. Stage Objective\n- Deliver value\n\n## 2. Risks\n- Risk one\n"
	if err := os.WriteFile(filepath.Join(docPath, "PLAN.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	store := openTestStore(t, workspace)
	defer store.Close()

	if err := store.ImportDocument(DocumentSeed{Path: "DOC/PLAN.md", DocumentType: "PLAN", Owner: "BA/CTO"}); err != nil {
		t.Fatalf("ImportDocument: %v", err)
	}

	var count int
	if err := store.db.QueryRow("SELECT COUNT(*) FROM context_entries WHERE document_type = 'PLAN' AND entry_type = 'plan_section'").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("plan context entries = %d, want 2", count)
	}
}

func TestImportStateCreatesContextEntries(t *testing.T) {
	workspace := t.TempDir()
	docPath := filepath.Join(workspace, "DOC")
	if err := os.MkdirAll(docPath, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# STATE\n\n## What Was Completed\n\n- **[2026-05-12 10:00] — [DEV_BACKEND:Senior]:** API delivered\n  - Tests Performed: go test ./...\n\n## Blockers\n\n- **[2026-05-12 11:00] — [QA:Senior]:** Need fixture\n"
	if err := os.WriteFile(filepath.Join(docPath, "STATE.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	store := openTestStore(t, workspace)
	defer store.Close()

	if err := store.ImportDocument(DocumentSeed{Path: "DOC/STATE.md", DocumentType: "STATE", Owner: "Technical agents"}); err != nil {
		t.Fatalf("ImportDocument: %v", err)
	}

	var deliveries, blockers int
	if err := store.db.QueryRow("SELECT COUNT(*) FROM context_entries WHERE document_type = 'STATE' AND entry_type = 'delivery'").Scan(&deliveries); err != nil {
		t.Fatal(err)
	}
	if err := store.db.QueryRow("SELECT COUNT(*) FROM context_entries WHERE document_type = 'STATE' AND entry_type = 'blocker'").Scan(&blockers); err != nil {
		t.Fatal(err)
	}
	if deliveries != 1 || blockers != 1 {
		t.Fatalf("deliveries=%d blockers=%d, want 1/1", deliveries, blockers)
	}
}

func TestImportContextCreatesDecisionEntries(t *testing.T) {
	workspace := t.TempDir()
	docPath := filepath.Join(workspace, "DOC")
	if err := os.MkdirAll(docPath, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# CONTEXT\n\n## Decisions\n\n### [2026-05-12 10:00] — [CTO]: Use SQLite\n- **Decision:** Keep local DB\n- **Why:** Simple MVP\n"
	if err := os.WriteFile(filepath.Join(docPath, "CONTEXT.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	store := openTestStore(t, workspace)
	defer store.Close()

	if err := store.ImportDocument(DocumentSeed{Path: "DOC/CONTEXT.md", DocumentType: "CONTEXT", Owner: "Management agents"}); err != nil {
		t.Fatalf("ImportDocument: %v", err)
	}

	var count int
	if err := store.db.QueryRow("SELECT COUNT(*) FROM context_entries WHERE document_type = 'CONTEXT' AND entry_type = 'decision'").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("decisions = %d, want 1", count)
	}
}

func TestImportRetrospectiveCreatesContextEntries(t *testing.T) {
	workspace := t.TempDir()
	docPath := filepath.Join(workspace, "DOC")
	if err := os.MkdirAll(docPath, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# RETRO\n\n### [2026-05-12 10:00] — Phase 1: MVP\n- **Consolidated by:** [CEO]\n#### What Worked\n- Fast loop\n"
	if err := os.WriteFile(filepath.Join(docPath, "RETROSPECTIVE.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	store := openTestStore(t, workspace)
	defer store.Close()

	if err := store.ImportDocument(DocumentSeed{Path: "DOC/RETROSPECTIVE.md", DocumentType: "RETROSPECTIVE", Owner: "CEO"}); err != nil {
		t.Fatalf("ImportDocument: %v", err)
	}

	var count int
	if err := store.db.QueryRow("SELECT COUNT(*) FROM context_entries WHERE document_type = 'RETROSPECTIVE' AND entry_type = 'retrospective'").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("retrospectives = %d, want 1", count)
	}
}

func TestImportVersionsCreatesPlanningItems(t *testing.T) {
	workspace := t.TempDir()
	docPath := filepath.Join(workspace, "DOC")
	if err := os.MkdirAll(docPath, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# VERSIONS\n\n### 1.1.0 — 2026-05-12 — [CEO]\n- **Type:** Release\n- **Description:** shipped\n\n### [BUGFIX] v1.1.1-fix.1 — 2026-05-13 — [QA]\n- **Severity:** High\n"
	if err := os.WriteFile(filepath.Join(docPath, "VERSIONS.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	store := openTestStore(t, workspace)
	defer store.Close()

	if err := store.ImportDocument(DocumentSeed{Path: "DOC/VERSIONS.md", DocumentType: "VERSIONS", Owner: "Phase-closing agent"}); err != nil {
		t.Fatalf("ImportDocument: %v", err)
	}

	var count int
	if err := store.db.QueryRow("SELECT COUNT(*) FROM planning_items WHERE task_ref = 'VERSIONS'").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("version planning items = %d, want 2", count)
	}
}

func TestImportTestsCreatesTestRecords(t *testing.T) {
	workspace := t.TempDir()
	docPath := filepath.Join(workspace, "DOC")
	if err := os.MkdirAll(docPath, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# TESTS\n\n### Test T-1 - Unit\n- Status: Passed\n- Command: go test ./...\n- Obtained result: ok\n"
	if err := os.WriteFile(filepath.Join(docPath, "TESTS.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	store := openTestStore(t, workspace)
	defer store.Close()

	if err := store.ImportDocument(DocumentSeed{Path: "DOC/TESTS.md", DocumentType: "TESTS", Owner: "QA"}); err != nil {
		t.Fatalf("ImportDocument: %v", err)
	}

	var count int
	var status string
	if err := store.db.QueryRow("SELECT COUNT(*), MAX(status) FROM test_records WHERE task_ref = 'T-1'").Scan(&count, &status); err != nil {
		t.Fatal(err)
	}
	if count != 1 || status != "passed" {
		t.Fatalf("count=%d status=%q, want 1/passed", count, status)
	}
}

func TestImportQuestionsCreatesContextEntries(t *testing.T) {
	workspace := t.TempDir()
	if err := os.WriteFile(filepath.Join(workspace, "QUESTIONS.md"), []byte("# QUESTIONS\n\n### [2026-05-12 10:00] Question: Choose DB\n- **Status:** Open\n- **Author:** [CTO]\n- **Question:** SQLite or Postgres?\n  > **Response:** [HUMAN fills here]\n"), 0644); err != nil {
		t.Fatal(err)
	}
	store := openTestStore(t, workspace)
	defer store.Close()

	if err := store.ImportDocument(DocumentSeed{Path: "QUESTIONS.md", DocumentType: "QUESTIONS", Owner: "Management agents"}); err != nil {
		t.Fatalf("ImportDocument: %v", err)
	}

	var count int
	if err := store.db.QueryRow("SELECT COUNT(*) FROM context_entries WHERE document_type = 'QUESTIONS' AND entry_type = 'question'").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("questions = %d, want 1", count)
	}
}

func TestImportPlaybookCreatesContextEntries(t *testing.T) {
	workspace := t.TempDir()
	content := "# PLAYBOOK\n\n## Preferences\n\n| Timestamp | Principle | Rationale |\n| :--- | :--- | :--- |\n| [2026-05-12 10:00] | Concise chat | Saves tokens |\n"
	if err := os.WriteFile(filepath.Join(workspace, "PLAYBOOK.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	store := openTestStore(t, workspace)
	defer store.Close()

	if err := store.ImportDocument(DocumentSeed{Path: "PLAYBOOK.md", DocumentType: "PLAYBOOK", Owner: "CEO"}); err != nil {
		t.Fatalf("ImportDocument: %v", err)
	}

	var count int
	if err := store.db.QueryRow("SELECT COUNT(*) FROM context_entries WHERE document_type = 'PLAYBOOK' AND entry_type = 'playbook_entry'").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("playbook entries = %d, want 1", count)
	}
}

func TestSaveContextEntryProjectsReadModel(t *testing.T) {
	workspace := t.TempDir()
	store := openTestStore(t, workspace)
	defer store.Close()

	if err := store.SaveContextEntry(ContextEntry{
		EntryType:    "decision",
		DocumentType: "CONTEXT",
		Title:        "Decision",
		Content:      "Use CQRS append-only storage.",
		SourceAgent:  "CTO",
		TaskRef:      "TASK-1",
	}); err != nil {
		t.Fatalf("SaveContextEntry: %v", err)
	}

	items, err := store.QueryContextForHandoff("CEO", "TASK-1", 3)
	if err != nil {
		t.Fatalf("QueryContextForHandoff: %v", err)
	}
	if len(items) != 1 || items[0] != "Use CQRS append-only storage." {
		t.Fatalf("items = %#v", items)
	}
}

func TestProjectLifecycleAndStartedFlag(t *testing.T) {
	workspace := t.TempDir()
	store := openTestStore(t, workspace)
	defer store.Close()

	started, err := store.ProjectStarted()
	if err != nil {
		t.Fatalf("ProjectStarted before project: %v", err)
	}
	if started {
		t.Fatal("ProjectStarted = true, want false without active project")
	}

	project, err := store.CreateProject(Project{Name: "PF AI", Slug: "PF AI"})
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	if project.ID == "" || project.Slug != "pf-ai" {
		t.Fatalf("project = %#v", project)
	}

	started, err = store.ProjectStarted()
	if err != nil {
		t.Fatalf("ProjectStarted after project: %v", err)
	}
	if !started {
		t.Fatal("ProjectStarted = false, want true with active project")
	}

	active, err := store.GetActiveProject()
	if err != nil {
		t.Fatalf("GetActiveProject: %v", err)
	}
	if active.ID != project.ID {
		t.Fatalf("active project = %s, want %s", active.ID, project.ID)
	}
}

func TestQueryContextForHandoffFiltersByActiveProject(t *testing.T) {
	workspace := t.TempDir()
	store := openTestStore(t, workspace)
	defer store.Close()

	projectA, err := store.CreateProject(Project{Name: "Project A", Slug: "project-a"})
	if err != nil {
		t.Fatalf("CreateProject A: %v", err)
	}
	if err := store.SaveContextEntry(ContextEntry{
		EntryType:    "decision",
		DocumentType: "CONTEXT",
		Title:        "A",
		Content:      "context from project A",
		SourceAgent:  "CTO",
		TaskRef:      "TASK-1",
	}); err != nil {
		t.Fatalf("SaveContextEntry A: %v", err)
	}

	projectB, err := store.CreateProject(Project{Name: "Project B", Slug: "project-b"})
	if err != nil {
		t.Fatalf("CreateProject B: %v", err)
	}
	if err := store.SaveContextEntry(ContextEntry{
		EntryType:    "decision",
		DocumentType: "CONTEXT",
		Title:        "B",
		Content:      "context from project B",
		SourceAgent:  "CTO",
		TaskRef:      "TASK-1",
	}); err != nil {
		t.Fatalf("SaveContextEntry B: %v", err)
	}

	store.ActiveProjectID = projectA.ID
	items, err := store.QueryContextForHandoff("CEO", "TASK-1", 10)
	if err != nil {
		t.Fatalf("QueryContextForHandoff A: %v", err)
	}
	if len(items) != 1 || items[0] != "context from project A" {
		t.Fatalf("items for A = %#v", items)
	}

	store.ActiveProjectID = projectB.ID
	items, err = store.QueryContextForHandoff("CEO", "TASK-1", 10)
	if err != nil {
		t.Fatalf("QueryContextForHandoff B: %v", err)
	}
	if len(items) != 1 || items[0] != "context from project B" {
		t.Fatalf("items for B = %#v", items)
	}
}

func openTestStore(t *testing.T, workspace string) *Store {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	migrationsDir := filepath.Clean(filepath.Join(wd, "..", "..", "db", "migrations"))
	store, err := Open(filepath.Join(workspace, "data", "test.db"), migrationsDir, workspace)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	return store
}
