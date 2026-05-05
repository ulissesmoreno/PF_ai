package domain_test

import (
	"testing"

	"pf-ai/src/domain"
)

func TestMemoryFileRejectsTraversal(t *testing.T) {
	_, err := domain.NewMemoryFile("../secrets.env", "secrets")
	if err == nil {
		t.Fatal("expected traversal rejection")
	}
}

func TestMemoryFileAcceptsRelativeProjectPath(t *testing.T) {
	memory, err := domain.NewMemoryFile("DOC/PLAN.md", "Plan")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if memory.Path != "DOC/PLAN.md" {
		t.Fatalf("expected normalized path DOC/PLAN.md, got %s", memory.Path)
	}
}
