package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"pf_ai/context_store"
)

func TestListCards(t *testing.T) {
	store := openAPITestStore(t)
	defer store.Close()
	seedCard(t, store, "card-1")

	req := httptest.NewRequest(http.MethodGet, "/api/cards?status=open&task_ref=TASK-1", nil)
	rec := httptest.NewRecorder()
	NewHandler(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	var response struct {
		Cards []context_store.Card `json:"cards"`
		Count int                  `json:"count"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response.Count != 1 || response.Cards[0].ID != "card-1" {
		t.Fatalf("response = %#v", response)
	}
}

func TestGetCardThread(t *testing.T) {
	store := openAPITestStore(t)
	defer store.Close()
	seedCard(t, store, "card-2")
	if err := store.AddCardComment("card-2", "DEV_BACKEND", "response", "ok"); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/cards/card-2", nil)
	rec := httptest.NewRecorder()
	NewHandler(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	var thread context_store.CardThread
	if err := json.NewDecoder(rec.Body).Decode(&thread); err != nil {
		t.Fatal(err)
	}
	if thread.Card.ID != "card-2" || len(thread.Comments) != 2 {
		t.Fatalf("thread = %#v", thread)
	}
}

func openAPITestStore(t *testing.T) *context_store.Store {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	workspace := t.TempDir()
	migrationsDir := filepath.Clean(filepath.Join(wd, "..", "..", "db", "migrations"))
	store, err := context_store.Open(filepath.Join(workspace, "data", "test.db"), migrationsDir, workspace)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	return store
}

func seedCard(t *testing.T, store *context_store.Store, id string) {
	t.Helper()
	if _, err := store.SaveCardHandoff(context_store.CardHandoff{
		CardID:    id,
		Title:     "Implement API",
		TaskRef:   "TASK-1",
		Sender:    "[CEO]",
		Recipient: "[DEV_BACKEND]",
		Intent:    "IMPLEMENT",
		RawJSON:   `{}`,
	}); err != nil {
		t.Fatalf("SaveCardHandoff: %v", err)
	}
}
