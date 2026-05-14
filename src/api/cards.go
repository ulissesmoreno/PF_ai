package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"pf_ai/context_store"
	"pf_ai/hand_off"
)

type Handler struct {
	store      *context_store.Store
	handoffDir string
}

func NewHandler(store *context_store.Store, handoffDir ...string) http.Handler {
	h := &Handler{store: store}
	if len(handoffDir) > 0 {
		h.handoffDir = handoffDir[0]
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.dashboard)
	mux.HandleFunc("GET /dashboard", h.dashboard)
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /api/health", h.health)
	mux.HandleFunc("GET /api/projects", h.listProjects)
	mux.HandleFunc("POST /api/projects", h.createProject)
	mux.HandleFunc("POST /api/projects/", h.projectAction)
	mux.HandleFunc("PATCH /api/projects/", h.updateProject)
	mux.HandleFunc("GET /api/agents", h.listAgents)
	mux.HandleFunc("PUT /api/agents/", h.updateAgent)
	mux.HandleFunc("GET /api/tokens/summary", h.tokenSummary)
	mux.HandleFunc("GET /api/tokens/by-card/", h.tokensByCard)
	mux.HandleFunc("GET /api/cards", h.listCards)
	mux.HandleFunc("PUT /api/cards/", h.respondCard)
	mux.HandleFunc("GET /api/cards/", h.getCard)
	return mux
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	active, err := h.store.GetActiveProject()
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":            "ok",
			"active_project_id": "",
			"project_started":   false,
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":            "ok",
		"active_project_id": active.ID,
		"active_project":    active,
		"project_started":   true,
	})
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
	id := cardPathID(r.URL.Path, "")
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

type cardRespondRequest struct {
	Author    string          `json:"author"`
	Response  json.RawMessage `json:"response"`
	Questions json.RawMessage `json:"questions"`
}

func (h *Handler) respondCard(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "/respond") {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "card not found"})
		return
	}
	id := cardPathID(r.URL.Path, "/respond")
	if id == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "card not found"})
		return
	}
	card, err := h.store.GetCard(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "card not found"})
		return
	}
	var req cardRespondRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	content, err := json.Marshal(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	author := strings.TrimSpace(req.Author)
	if author == "" {
		author = "[HUMAN]"
	}
	if err := h.store.AddCardComment(id, author, "response", string(content)); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.store.UpdateCardStatus(id, "in_progress"); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if h.handoffDir != "" {
		payload := map[string]any{
			"reply_to_card_id": id,
			"response":         req.Response,
			"questions":        req.Questions,
		}
		path, err := hand_off.CreateHandoff(h.handoffDir, hand_off.HandoffHeader{
			CardID:    id,
			Sender:    "[HUMAN]",
			Recipient: card.Recipient,
			TaskRef:   card.TaskRef,
			Intent:    "HUMAN_RESPONSE",
		}, payload)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusAccepted, map[string]any{"card": card, "handoff_path": path})
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{"card": card})
}

func cardPathID(path, suffix string) string {
	id := strings.TrimPrefix(path, "/api/cards/")
	if suffix != "" {
		id = strings.TrimSuffix(id, suffix)
	}
	return strings.Trim(id, "/")
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
