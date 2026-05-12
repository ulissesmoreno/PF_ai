package code_writer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteCodeFilesWritesInsideWorkspace(t *testing.T) {
	baseDir := t.TempDir()

	written, err := WriteCodeFiles(baseDir, []CodeFile{
		{Path: "src/example.txt", Content: "ok"},
	})
	if err != nil {
		t.Fatalf("WriteCodeFiles returned error: %v", err)
	}
	if len(written) != 1 {
		t.Fatalf("written count = %d, want 1", len(written))
	}

	data, err := os.ReadFile(filepath.Join(baseDir, "src", "example.txt"))
	if err != nil {
		t.Fatalf("expected written file: %v", err)
	}
	if string(data) != "ok" {
		t.Fatalf("content = %q, want %q", string(data), "ok")
	}
}

func TestWriteCodeFilesRejectsPathTraversal(t *testing.T) {
	baseDir := t.TempDir()

	if _, err := WriteCodeFiles(baseDir, []CodeFile{
		{Path: "../outside.txt", Content: "no"},
	}); err == nil {
		t.Fatal("expected path traversal to be rejected")
	}
}
