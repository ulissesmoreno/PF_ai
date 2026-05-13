package pipeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"pf_ai/code_writer"
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

	if err := p.applyAgentResponse("DEV_BACKEND", "TASK-1", response); err != nil {
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

	err := p.applyAgentResponse("DEV_BACKEND", "TASK-1", "Claro, aqui esta o codigo:\n```python\nprint(30)\n```")
	if err == nil {
		t.Fatal("expected text response to be rejected")
	}
	if !strings.Contains(err.Error(), "sem JSON") {
		t.Fatalf("unexpected error: %v", err)
	}
}
