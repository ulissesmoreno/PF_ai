package domain_test

import (
	"testing"

	"pf-ai/src/domain"
)

func TestHandoffRejectsSecretLikePayload(t *testing.T) {
	_, err := domain.NewHandoff("[CEO]", "[DEV_BACKEND:Pleno]", "TASK-1", domain.IntentPhaseKickoff, map[string]any{
		"api_key": "do-not-store",
	}, nil)
	if err == nil {
		t.Fatal("expected secret-like field rejection")
	}
}

func TestHandoffAcceptsValidPayload(t *testing.T) {
	item, err := domain.NewHandoff("[CEO]", "[DEV_BACKEND:Pleno]", "TASK-1", domain.IntentPhaseKickoff, map[string]any{
		"phase_ref": "DOC/ROADMAP.md#phase-1",
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Header.Intent != domain.IntentPhaseKickoff {
		t.Fatalf("expected PHASE_KICKOFF, got %s", item.Header.Intent)
	}
}
