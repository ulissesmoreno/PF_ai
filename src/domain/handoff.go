package domain

import (
	"errors"
	"strings"
	"time"
)

type HandoffIntent string

const (
	IntentPhaseKickoff          HandoffIntent = "PHASE_KICKOFF"
	IntentClarificationRequest  HandoffIntent = "CLARIFICATION_REQUEST"
	IntentClarificationResponse HandoffIntent = "CLARIFICATION_RESPONSE"
	IntentReviewResult          HandoffIntent = "REVIEW_RESULT"
	IntentStageApproval         HandoffIntent = "STAGE_APPROVAL"
	IntentBugfixRequest         HandoffIntent = "BUGFIX_REQUEST"
	IntentRetrospective         HandoffIntent = "RETROSPECTIVE"
)

type HandoffHeader struct {
	Timestamp string        `json:"timestamp"`
	Sender    string        `json:"sender"`
	Recipient string        `json:"recipient"`
	TaskRef   string        `json:"task_ref"`
	Intent    HandoffIntent `json:"intent"`
}

type Handoff struct {
	Header  HandoffHeader         `json:"header"`
	Payload map[string]any        `json:"payload"`
	Meta    map[string]string     `json:"meta,omitempty"`
	Sources []MemoryFile          `json:"sources,omitempty"`
}

func NewHandoff(sender, recipient, taskRef string, intent HandoffIntent, payload map[string]any, sources []MemoryFile) (Handoff, error) {
	header := HandoffHeader{
		Timestamp: time.Now().Format("2006-01-02 15:04"),
		Sender:    strings.TrimSpace(sender),
		Recipient: strings.TrimSpace(recipient),
		TaskRef:   strings.TrimSpace(taskRef),
		Intent:    intent,
	}

	handoff := Handoff{
		Header:  header,
		Payload: payload,
		Sources: sources,
	}

	if err := handoff.Validate(); err != nil {
		return Handoff{}, err
	}

	return handoff, nil
}

func (h Handoff) Validate() error {
	if h.Header.Sender == "" {
		return errors.New("handoff sender is required")
	}
	if h.Header.Recipient == "" {
		return errors.New("handoff recipient is required")
	}
	if h.Header.TaskRef == "" {
		return errors.New("handoff task ref is required")
	}
	if !h.Header.Intent.Valid() {
		return errors.New("handoff intent is invalid")
	}
	if h.Payload == nil {
		return errors.New("handoff payload is required")
	}
	if containsSecretLikeField(h.Payload) {
		return errors.New("handoff payload contains secret-like field")
	}

	return nil
}

func (i HandoffIntent) Valid() bool {
	switch i {
	case IntentPhaseKickoff, IntentClarificationRequest, IntentClarificationResponse, IntentReviewResult, IntentStageApproval, IntentBugfixRequest, IntentRetrospective:
		return true
	default:
		return false
	}
}

func containsSecretLikeField(payload map[string]any) bool {
	for key, value := range payload {
		normalized := strings.ToLower(key)
		if strings.Contains(normalized, "secret") || strings.Contains(normalized, "password") || strings.Contains(normalized, "token") || strings.Contains(normalized, "api_key") || strings.Contains(normalized, "apikey") {
			return true
		}
		nested, ok := value.(map[string]any)
		if ok && containsSecretLikeField(nested) {
			return true
		}
	}

	return false
}
