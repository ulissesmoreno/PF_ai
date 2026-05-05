package domain

import (
	"errors"
	"strings"
	"time"
)

type Session struct {
	ID        string
	Subject   string
	ExpiresAt time.Time
}

func NewSession(id, subject string, expiresAt time.Time) (Session, error) {
	session := Session{
		ID:        strings.TrimSpace(id),
		Subject:   strings.TrimSpace(subject),
		ExpiresAt: expiresAt,
	}

	if session.ID == "" {
		return Session{}, errors.New("session id is required")
	}
	if session.Subject == "" {
		return Session{}, errors.New("session subject is required")
	}
	if session.ExpiresAt.IsZero() {
		return Session{}, errors.New("session expiration is required")
	}

	return session, nil
}
