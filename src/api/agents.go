package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"pf_ai/context_store"
)

func (h *Handler) listAgents(w http.ResponseWriter, _ *http.Request) {
	configs, err := h.store.ListAgentConfigs()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, struct {
		Agents []context_store.AgentConfig `json:"agents"`
		Count  int                         `json:"count"`
	}{Agents: configs, Count: len(configs)})
}

func (h *Handler) updateAgent(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/agents/"), "/")
	if name == "" || strings.Contains(name, "/") {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "agent not found"})
		return
	}
	var config context_store.AgentConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	config.AgentName = name
	updated, err := h.store.UpsertAgentConfig(config)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}
