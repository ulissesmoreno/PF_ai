package agent

import "testing"

func TestResolverProviderAgenteDefaults(t *testing.T) {
	t.Setenv("CEO_AGENT_PROVIDER", "")
	t.Setenv("DEV_BACK_AGENT_PROVIDER", "")

	if got := ResolverProviderAgente("CEO"); got != "ollama" {
		t.Fatalf("CEO provider = %q, want ollama", got)
	}

	if got := ResolverProviderAgente("DEV_BACKEND"); got != "ollama" {
		t.Fatalf("DEV_BACKEND provider = %q, want ollama", got)
	}
}

func TestResolverProviderAgenteFromEnv(t *testing.T) {
	t.Setenv("DEV_BACK_AGENT_PROVIDER", "cli")

	if got := ResolverProviderAgente("DEV_BACKEND"); got != "cli" {
		t.Fatalf("DEV_BACKEND provider = %q, want cli", got)
	}
}

func TestResolverModeloAgenteFromEnv(t *testing.T) {
	t.Setenv("DEV_BACK_AGENT_MODEL", "gpt-test")

	got, err := ResolverModeloAgente("DEV_BACKEND")
	if err != nil {
		t.Fatalf("ResolverModeloAgente returned error: %v", err)
	}
	if got != "gpt-test" {
		t.Fatalf("DEV_BACKEND model = %q, want gpt-test", got)
	}
}

func TestExpandCLICommandModelPlaceholder(t *testing.T) {
	got := expandCLICommand("codex exec --model {{MODEL}} -", "gpt-test")
	want := "codex exec --model gpt-test -"
	if got != want {
		t.Fatalf("expandCLICommand = %q, want %q", got, want)
	}
}

func TestCarregarConfigAgenteDoesNotMutateGlobalModel(t *testing.T) {
	Init("http://localhost:11434", "base-model")
	SetConfigDir("../../AGENTS")
	t.Setenv("CEO_AGENT_MODEL", "cli-model")

	if _, _, err := CarregarConfigAgente("CEO"); err != nil {
		t.Fatalf("CarregarConfigAgente returned error: %v", err)
	}

	mu.RLock()
	got := model
	mu.RUnlock()

	if got != "base-model" {
		t.Fatalf("global model = %q, want base-model", got)
	}
}
