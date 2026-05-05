package infrastructure_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"pf-ai/src/domain"
	"pf-ai/src/infrastructure/filememory"
)

func TestFileMemoryReaderRejectsPathOutsideAllowlist(t *testing.T) {
	root := t.TempDir()
	reader := filememory.NewReader(root, []string{"DOC/PLAN.md"})
	memory, err := domain.NewMemoryFile("README.md", "Readme")
	if err != nil {
		t.Fatalf("unexpected memory error: %v", err)
	}

	_, err = reader.ReadMemory(context.Background(), memory)
	if err == nil {
		t.Fatal("expected allowlist rejection")
	}
}

func TestFileMemoryReaderReadsAllowedFile(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "DOC"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "DOC", "PLAN.md"), []byte("phase plan"), 0o600); err != nil {
		t.Fatal(err)
	}

	reader := filememory.NewReader(root, []string{"DOC/PLAN.md"})
	memory, err := domain.NewMemoryFile("DOC/PLAN.md", "Plan")
	if err != nil {
		t.Fatalf("unexpected memory error: %v", err)
	}

	content, err := reader.ReadMemory(context.Background(), memory)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}
	if content != "phase plan" {
		t.Fatalf("expected phase plan, got %s", content)
	}
}
