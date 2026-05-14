package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pf_ai/context_store"
)

func TestAgentConfigAPI(t *testing.T) {
	store := openAPITestStore(t)
	defer store.Close()
	handler := NewHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/api/agents", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d body=%s", rec.Code, rec.Body.String())
	}

	body := bytes.NewBufferString(`{"model":"qwen2.5-coder","tier":2,"active":true,"updated_by":"test"}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/agents/QA", body)
	updateRec := httptest.NewRecorder()
	handler.ServeHTTP(updateRec, updateReq)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("update status = %d body=%s", updateRec.Code, updateRec.Body.String())
	}
	var updated context_store.AgentConfig
	if err := json.NewDecoder(updateRec.Body).Decode(&updated); err != nil {
		t.Fatal(err)
	}
	if updated.AgentName != "QA" || updated.Model != "qwen2.5-coder" {
		t.Fatalf("updated = %#v", updated)
	}
}

func TestTokenUsageAPI(t *testing.T) {
	store := openAPITestStore(t)
	defer store.Close()
	if err := store.SaveTokenUsage(context_store.TokenUsage{
		AgentName:        "CEO",
		Model:            "deepseek-r1:7b",
		Tier:             3,
		PromptTokens:     20,
		CompletionTokens: 7,
	}); err != nil {
		t.Fatal(err)
	}

	handler := NewHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/api/tokens/summary?agent=CEO", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("summary status = %d body=%s", rec.Code, rec.Body.String())
	}
	var usage context_store.TokenUsage
	if err := json.NewDecoder(rec.Body).Decode(&usage); err != nil {
		t.Fatal(err)
	}
	if usage.TotalTokens != 27 {
		t.Fatalf("usage = %#v", usage)
	}
}

func TestDashboardRoute(t *testing.T) {
	store := openAPITestStore(t)
	defer store.Close()

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	rec := httptest.NewRecorder()
	NewHandler(store).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("PF AI Cards")) {
		t.Fatalf("dashboard body missing title")
	}
}
