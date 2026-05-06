package domain_test

import (
	"testing"

	"pf-ai/src/domain"
)

func TestPlanningCardRequiresOwnerAndStatus(t *testing.T) {
	_, err := domain.NewPlanningCard("card-1", "Phase 3 board", "", domain.CardStatusTodo, domain.CardPriorityHigh, "CEO", "PHASE3-CARD-001")
	if err == nil {
		t.Fatal("expected missing owner error")
	}
}

func TestProjectRegistrationRejectsSecretLikeOnboarding(t *testing.T) {
	_, err := domain.NewProjectRegistration("pf-ai", "PF_ai", map[string]any{"api_key": "raw"})
	if err == nil {
		t.Fatal("expected secret-like onboarding rejection")
	}
}

func TestHandoffExportsMCPEnvelope(t *testing.T) {
	handoff, err := domain.NewHandoff("CEO", "DEV_BACKEND", "PHASE3", domain.IntentPhaseKickoff, map[string]any{"summary": "build"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	envelope := domain.NewMCPEnvelope(handoff)

	if envelope.Protocol != "mcp-compatible" {
		t.Fatalf("expected mcp-compatible protocol, got %s", envelope.Protocol)
	}
	if envelope.Message.Header.TaskRef != "PHASE3" {
		t.Fatalf("expected task ref, got %s", envelope.Message.Header.TaskRef)
	}
}
