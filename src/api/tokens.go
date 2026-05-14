package api

import (
	"net/http"
	"strings"

	"pf_ai/context_store"
)

func (h *Handler) tokenSummary(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	usage, err := h.store.SumTokenUsage(context_store.TokenUsageFilter{
		ProjectID: q.Get("project_id"),
		AgentName: q.Get("agent"),
		Model:     q.Get("model"),
		From:      q.Get("from"),
		To:        q.Get("to"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, usage)
}

func (h *Handler) tokensByCard(w http.ResponseWriter, r *http.Request) {
	cardID := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/tokens/by-card/"), "/")
	if cardID == "" || strings.Contains(cardID, "/") {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "card not found"})
		return
	}
	usages, err := h.store.ListTokenUsageByCard(cardID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, struct {
		Usages []context_store.TokenUsage `json:"usages"`
		Count  int                        `json:"count"`
	}{Usages: usages, Count: len(usages)})
}
