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
