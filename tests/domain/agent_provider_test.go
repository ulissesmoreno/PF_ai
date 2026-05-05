package domain_test

import (
	"testing"

	"pf-ai/src/domain"
)

func TestNewAgentDefinitionRequiresValidSeniority(t *testing.T) {
	_, err := domain.NewAgentDefinition("ceo", "CEO", "orchestration", domain.Seniority("Lead"), "openai", "")
	if err == nil {
		t.Fatal("expected invalid seniority error")
	}
}

func TestNewAgentDefinitionAcceptsCompleteAgent(t *testing.T) {
	agent, err := domain.NewAgentDefinition("cto", "CTO", "architecture", domain.SenioritySenior, "local", "guards architecture")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if agent.ID != "cto" {
		t.Fatalf("expected id cto, got %s", agent.ID)
	}
}

func TestNewAgentDefinitionAllowsNoProvider(t *testing.T) {
	agent, err := domain.NewAgentDefinition("codex", "Codex", "DEV_BACKEND", domain.SenioritySenior, "", "embedded model runtime")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if agent.ProviderID != "" {
		t.Fatalf("expected no provider, got %s", agent.ProviderID)
	}
}

func TestAPIProviderRequiresSecretReference(t *testing.T) {
	_, err := domain.NewModelProvider("openai", "OpenAI", domain.ProviderModeAPI, "https://api.example.test", "gpt", "", "")
	if err == nil {
		t.Fatal("expected missing secret reference error")
	}
}

func TestLocalProviderRequiresRuntime(t *testing.T) {
	_, err := domain.NewModelProvider("ollama", "Ollama", domain.ProviderModeLocal, "http://localhost:11434", "llama", "", "")
	if err == nil {
		t.Fatal("expected missing local runtime error")
	}
}
