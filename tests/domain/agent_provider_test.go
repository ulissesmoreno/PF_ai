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

func TestHybridProviderRequiresSecretAndRuntime(t *testing.T) {
	_, err := domain.NewModelProvider("hybrid", "Hybrid", domain.ProviderModeHybrid, "http://127.0.0.1:11434", "llama", "PF_AI_API_KEY", "")
	if err == nil {
		t.Fatal("expected missing local runtime error")
	}
}

func TestLocalProviderRejectsUnsafeEndpointScheme(t *testing.T) {
	_, err := domain.NewModelProvider("local", "Local", domain.ProviderModeLocal, "file:///tmp/model", "llama", "", "ollama")
	if err == nil {
		t.Fatal("expected unsafe endpoint error")
	}
}

func TestLocalProviderRejectsNonLoopbackEndpoint(t *testing.T) {
	_, err := domain.NewModelProvider("local", "Local", domain.ProviderModeLocal, "http://example.com", "llama", "", "ollama")
	if err == nil {
		t.Fatal("expected non-loopback endpoint error")
	}
}

func TestHybridRouteFallsBackToAPIWhenLocalIsUnhealthy(t *testing.T) {
	provider, err := domain.NewModelProvider("hybrid", "Hybrid", domain.ProviderModeHybrid, "http://127.0.0.1:11434", "llama", "PF_AI_API_KEY", "ollama")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	decision := domain.DecideProviderRoute(provider, domain.ProviderHealth{LocalHealthy: false, APIHealthy: true})

	if decision.Route != "api" {
		t.Fatalf("expected api fallback, got %s", decision.Route)
	}
	if decision.FallbackReason == "" {
		t.Fatal("expected fallback reason")
	}
}

func TestLocalRouteReportsUnavailableWhenHealthFails(t *testing.T) {
	provider, err := domain.NewModelProvider("local", "Local", domain.ProviderModeLocal, "http://127.0.0.1:11434", "llama", "", "ollama")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	decision := domain.DecideProviderRoute(provider, domain.ProviderHealth{LocalHealthy: false})

	if decision.Status != domain.ProviderRouteUnavailable {
		t.Fatalf("expected unavailable route, got %s", decision.Status)
	}
}
