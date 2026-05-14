package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"pf_ai/context_store"
)

func (h *Handler) listProjects(w http.ResponseWriter, _ *http.Request) {
	projects, err := h.store.ListProjects()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, struct {
		Projects []context_store.Project `json:"projects"`
		Count    int                     `json:"count"`
	}{Projects: projects, Count: len(projects)})
}

func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	var project context_store.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	created, err := h.store.CreateProject(project)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) updateProject(w http.ResponseWriter, r *http.Request) {
	id := projectPathID(r.URL.Path, "")
	if id == "" || strings.Contains(id, "/") {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "project not found"})
		return
	}
	var project context_store.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	updated, err := h.store.UpdateProject(id, project)
	if err != nil {
		writeProjectError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) projectAction(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "/activate") {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "project not found"})
		return
	}
	id := projectPathID(r.URL.Path, "/activate")
	if id == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "project not found"})
		return
	}
	if err := h.store.ActivateProject(id); err != nil {
		writeProjectError(w, err)
		return
	}
	project, err := h.store.GetProject(id)
	if err != nil {
		writeProjectError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, project)
}

func projectPathID(path, suffix string) string {
	id := strings.TrimPrefix(path, "/api/projects/")
	if suffix != "" {
		id = strings.TrimSuffix(id, suffix)
	}
	return strings.Trim(id, "/")
}

func writeProjectError(w http.ResponseWriter, err error) {
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "project not found"})
		return
	}
	writeError(w, http.StatusInternalServerError, err)
}
