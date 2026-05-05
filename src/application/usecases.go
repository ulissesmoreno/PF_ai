package application

import (
	"context"

	"pf-ai/src/domain"
)

type RegistryService struct {
	Agents    AgentRepository
	Providers ProviderRepository
}

func (s RegistryService) RegisterAgent(ctx context.Context, agent domain.AgentDefinition) error {
	if !agent.Seniority.Valid() {
		return domainError("agent seniority is invalid")
	}
	return s.Agents.SaveAgent(ctx, agent)
}

func (s RegistryService) RegisterProvider(ctx context.Context, provider domain.ModelProvider) error {
	if !provider.Mode.Valid() {
		return domainError("provider mode is invalid")
	}
	return s.Providers.SaveProvider(ctx, provider)
}

type domainError string

func (e domainError) Error() string {
	return string(e)
}

type HandoffService struct {
	Writer HandoffWriter
}

func (s HandoffService) CreateHandoff(ctx context.Context, handoff domain.Handoff) (string, error) {
	if err := handoff.Validate(); err != nil {
		return "", err
	}
	return s.Writer.WriteHandoff(ctx, handoff)
}
