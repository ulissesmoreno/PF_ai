package infrastructure_test

import (
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
