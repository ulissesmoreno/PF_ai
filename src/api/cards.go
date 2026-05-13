package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"pf_ai/context_store"
)

type Handler struct {
	store *context_store.Store
}

func NewHandler(store *context_store.Store) http.Handler {
	h := &Handler{store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /api/cards", h.listCards)
	mux.HandleFunc("GET /api/cards/", h.getCard)
	return mux
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) listCards(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	cards, err := h.store.ListCards(context_store.CardFilter{
		ProjectID: q.Get("project_id"),
		Status:    q.Get("status"),
		TaskRef:   q.Get("task_ref"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, struct {
		Cards []context_store.Card `json:"cards"`
		Count int                  `json:"count"`
	}{Cards: cards, Count: len(cards)})
}

func (h *Handler) getCard(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/cards/")
	id = strings.TrimSpace(id)
	if id == "" || strings.Contains(id, "/") {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "card not found"})
		return
	}
	thread, err := h.store.GetCardThread(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "card not found"})
		return
	}
	writeJSON(w, http.StatusOK, thread)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
