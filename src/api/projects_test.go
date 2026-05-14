package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pf_ai/context_store"
)

func TestProjectAPI(t *testing.T) {
	store := openAPITestStore(t)
	defer store.Close()
	handler := NewHandler(store)

	createBody := bytes.NewBufferString(`{"name":"PF AI","slug":"PF AI","description":"Agents"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/api/projects", createBody)
	createRec := httptest.NewRecorder()
	handler.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create status = %d body=%s", createRec.Code, createRec.Body.String())
	}
	var created context_store.Project
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	if created.ID == "" || created.Slug != "pf-ai" {
		t.Fatalf("created = %#v", created)
	}

	patchReq := httptest.NewRequest(http.MethodPatch, "/api/projects/"+created.ID, bytes.NewBufferString(`{"description":"Updated"}`))
	patchRec := httptest.NewRecorder()
	handler.ServeHTTP(patchRec, patchReq)
	if patchRec.Code != http.StatusOK {
		t.Fatalf("patch status = %d body=%s", patchRec.Code, patchRec.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	listRec := httptest.NewRecorder()
	handler.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("list status = %d body=%s", listRec.Code, listRec.Body.String())
	}
	var list struct {
		Projects []context_store.Project `json:"projects"`
		Count    int                     `json:"count"`
	}
	if err := json.NewDecoder(listRec.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}
	if list.Count != 1 || list.Projects[0].Description != "Updated" {
		t.Fatalf("list = %#v", list)
	}
}

func TestHealthIncludesActiveProject(t *testing.T) {
	store := openAPITestStore(t)
	defer store.Close()
	if _, err := store.CreateProject(context_store.Project{Name: "PF AI", Slug: "pf-ai"}); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	NewHandler(store).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	var response map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response["project_started"] != true {
		t.Fatalf("health = %#v", response)
	}
}
