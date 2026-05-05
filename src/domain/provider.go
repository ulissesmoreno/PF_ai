package domain

import (
	"errors"
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
