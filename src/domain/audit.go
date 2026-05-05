package domain

import (
	"errors"
	"strings"
)

type AuditEvent struct {
	Actor        string
	Action       string
	Target       string
	MetadataJSON string
}

func NewAuditEvent(actor, action, target, metadataJSON string) (AuditEvent, error) {
	event := AuditEvent{
		Actor:        strings.TrimSpace(actor),
		Action:       strings.TrimSpace(action),
		Target:       strings.TrimSpace(target),
		MetadataJSON: strings.TrimSpace(metadataJSON),
	}

	if event.Actor == "" {
		return AuditEvent{}, errors.New("audit actor is required")
	}
	if event.Action == "" {
		return AuditEvent{}, errors.New("audit action is required")
	}
	if event.Target == "" {
		return AuditEvent{}, errors.New("audit target is required")
	}
	if event.MetadataJSON == "" {
		event.MetadataJSON = "{}"
	}

	return event, nil
}
