package application

import (
	"context"

	"pf-ai/src/domain"
)

type AgentRepository interface {
	SaveAgent(ctx context.Context, agent domain.AgentDefinition) error
	ListAgents(ctx context.Context) ([]domain.AgentDefinition, error)
}

type ProviderRepository interface {
	SaveProvider(ctx context.Context, provider domain.ModelProvider) error
	ListProviders(ctx context.Context) ([]domain.ModelProvider, error)
}

type MemoryReader interface {
	ReadMemory(ctx context.Context, memory domain.MemoryFile) (string, error)
}

type HandoffWriter interface {
	WriteHandoff(ctx context.Context, handoff domain.Handoff) (string, error)
}
