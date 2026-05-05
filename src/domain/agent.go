package domain

import (
	"errors"
	"strings"
)

type Seniority string

const (
	SeniorityJunior Seniority = "Junior"
	SeniorityPleno  Seniority = "Pleno"
	SenioritySenior Seniority = "Senior"
)

type AgentDefinition struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Role        string    `json:"role"`
	Seniority   Seniority `json:"seniority"`
	ProviderID  string    `json:"provider_id,omitempty"`
	Description string    `json:"description"`
}

func NewAgentDefinition(id, name, role string, seniority Seniority, providerID, description string) (AgentDefinition, error) {
	agent := AgentDefinition{
		ID:          strings.TrimSpace(id),
		Name:        strings.TrimSpace(name),
		Role:        strings.TrimSpace(role),
		Seniority:   seniority,
		ProviderID:  strings.TrimSpace(providerID),
		Description: strings.TrimSpace(description),
	}

	if agent.ID == "" {
		return AgentDefinition{}, errors.New("agent id is required")
	}
	if agent.Name == "" {
		return AgentDefinition{}, errors.New("agent name is required")
	}
	if agent.Role == "" {
		return AgentDefinition{}, errors.New("agent role is required")
	}
	if !agent.Seniority.Valid() {
		return AgentDefinition{}, errors.New("agent seniority is invalid")
	}
	return agent, nil
}

func (s Seniority) Valid() bool {
	switch s {
	case SeniorityJunior, SeniorityPleno, SenioritySenior:
		return true
	default:
		return false
	}
}
