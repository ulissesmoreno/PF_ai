package domain

import (
	"errors"
	"net/url"
	"strings"
)

type ProviderMode string

const (
	ProviderModeAPI    ProviderMode = "api"
	ProviderModeLocal  ProviderMode = "local"
	ProviderModeHybrid ProviderMode = "hybrid"
)

type ModelProvider struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Mode         ProviderMode `json:"mode"`
	Endpoint     string       `json:"endpoint"`
	Model        string       `json:"model"`
	SecretRef    string       `json:"secret_ref,omitempty"`
	LocalRuntime string       `json:"local_runtime,omitempty"`
}

func NewModelProvider(id, name string, mode ProviderMode, endpoint, model, secretRef, localRuntime string) (ModelProvider, error) {
	provider := ModelProvider{
		ID:           strings.TrimSpace(id),
		Name:         strings.TrimSpace(name),
		Mode:         mode,
		Endpoint:     strings.TrimSpace(endpoint),
		Model:        strings.TrimSpace(model),
		SecretRef:    strings.TrimSpace(secretRef),
		LocalRuntime: strings.TrimSpace(localRuntime),
	}

	if provider.ID == "" {
		return ModelProvider{}, errors.New("provider id is required")
	}
	if provider.Name == "" {
		return ModelProvider{}, errors.New("provider name is required")
	}
	if !provider.Mode.Valid() {
		return ModelProvider{}, errors.New("provider mode is invalid")
	}
	if provider.Model == "" {
		return ModelProvider{}, errors.New("provider model is required")
	}
	if provider.Mode == ProviderModeAPI && provider.SecretRef == "" {
		return ModelProvider{}, errors.New("api provider requires a secret reference")
	}
	if provider.Mode == ProviderModeLocal && provider.LocalRuntime == "" {
		return ModelProvider{}, errors.New("local provider requires a local runtime")
	}
	if provider.Mode == ProviderModeHybrid && (provider.SecretRef == "" || provider.LocalRuntime == "") {
		return ModelProvider{}, errors.New("hybrid provider requires a secret reference and local runtime")
	}
	if provider.Endpoint != "" && !safeProviderEndpoint(provider.Mode, provider.Endpoint) {
		return ModelProvider{}, errors.New("provider endpoint scheme is not allowed")
	}

	return provider, nil
}

func (m ProviderMode) Valid() bool {
	switch m {
	case ProviderModeAPI, ProviderModeLocal, ProviderModeHybrid:
		return true
	default:
		return false
	}
}

type ProviderRouteStatus string

const (
	ProviderRouteAvailable   ProviderRouteStatus = "available"
	ProviderRouteUnavailable ProviderRouteStatus = "unavailable"
)

type ProviderHealth struct {
	LocalHealthy bool `json:"local_healthy"`
	APIHealthy   bool `json:"api_healthy"`
}

type ProviderRouteDecision struct {
	ProviderID     string              `json:"provider_id"`
	Mode           ProviderMode        `json:"mode"`
	Route          string              `json:"route"`
	Status         ProviderRouteStatus `json:"status"`
	FallbackReason string              `json:"fallback_reason,omitempty"`
}

func DecideProviderRoute(provider ModelProvider, health ProviderHealth) ProviderRouteDecision {
	decision := ProviderRouteDecision{
		ProviderID: provider.ID,
		Mode:       provider.Mode,
		Status:     ProviderRouteAvailable,
	}

	switch provider.Mode {
	case ProviderModeAPI:
		if health.APIHealthy {
			decision.Route = "api"
			return decision
		}
	case ProviderModeLocal:
		if health.LocalHealthy {
			decision.Route = "local"
			return decision
		}
	case ProviderModeHybrid:
		if health.LocalHealthy {
			decision.Route = "local"
			return decision
		}
		if health.APIHealthy {
			decision.Route = "api"
			decision.FallbackReason = "local runtime unhealthy"
			return decision
		}
	}

	decision.Status = ProviderRouteUnavailable
	decision.Route = "none"
	return decision
}

func safeProviderEndpoint(mode ProviderMode, raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	if mode == ProviderModeAPI {
		return true
	}
	switch parsed.Hostname() {
	case "localhost", "127.0.0.1", "::1":
		return true
	default:
		return false
	}
}
