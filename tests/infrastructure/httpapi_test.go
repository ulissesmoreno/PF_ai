package infrastructure_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"pf-ai/src/infrastructure/httpapi"
)

func TestProtectedRoutesReturnUnauthorizedWithoutToken(t *testing.T) {
	server := (&httpapi.Server{AuthSecret: "local-test-secret"}).Routes()

	request := httptest.NewRequest(http.MethodGet, "/api/agents", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}
}

func TestHealthRouteDoesNotRequireToken(t *testing.T) {
	server := (&httpapi.Server{AuthSecret: "local-test-secret"}).Routes()

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
}

func TestPostAgentCreatesAgent(t *testing.T) {
	server := (&httpapi.Server{AuthSecret: "local-test-secret"}).Routes()

	body := bytes.NewBufferString(`{"id":"ceo","name":"CEO","role":"orchestration","seniority":"Senior","provider_id":"api-default","description":"orchestrates"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/agents", body)
	request.Header.Set("Authorization", "Bearer local-test-secret")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", response.Code, response.Body.String())
	}
}

func TestPostProviderRejectsMissingSecretReference(t *testing.T) {
	server := (&httpapi.Server{AuthSecret: "local-test-secret"}).Routes()

	body := bytes.NewBufferString(`{"id":"api","name":"API","mode":"api","model":"gpt"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/providers", body)
	request.Header.Set("Authorization", "Bearer local-test-secret")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}
}
