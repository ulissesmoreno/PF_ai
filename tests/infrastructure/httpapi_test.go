package infrastructure_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"pf-ai/src/domain"
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

func TestOptionsPreflightDoesNotRequireToken(t *testing.T) {
	server := (&httpapi.Server{AuthSecret: "local-test-secret"}).Routes()

	request := httptest.NewRequest(http.MethodOptions, "/api/agents", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", response.Code)
	}
	if response.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("expected CORS allow origin header")
	}
}

func TestProtectedRoutesIncludeCORSHeaders(t *testing.T) {
	server := (&httpapi.Server{AuthSecret: "local-test-secret"}).Routes()

	request := httptest.NewRequest(http.MethodGet, "/api/agents", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	if response.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Fatal("expected CORS allow headers")
	}
}

func TestProviderHealthRouteRequiresToken(t *testing.T) {
	server := (&httpapi.Server{AuthSecret: "local-test-secret"}).Routes()

	request := httptest.NewRequest(http.MethodGet, "/api/provider-health?id=hybrid", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
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

func TestPostAgentAllowsNoProvider(t *testing.T) {
	server := (&httpapi.Server{AuthSecret: "local-test-secret"}).Routes()

	body := bytes.NewBufferString(`{"id":"codex","name":"Codex","role":"DEV_BACKEND","seniority":"Senior","description":"embedded model runtime"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/agents", body)
	request.Header.Set("Authorization", "Bearer local-test-secret")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", response.Code)
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

func TestProviderHealthReturnsHybridFallbackDecision(t *testing.T) {
	provider, err := domain.NewModelProvider("hybrid", "Hybrid", domain.ProviderModeHybrid, "http://127.0.0.1:11434", "llama", "PF_AI_API_KEY", "ollama")
	if err != nil {
		t.Fatal(err)
	}
	server := (&httpapi.Server{
		AuthSecret: "local-test-secret",
		Providers:  []domain.ModelProvider{provider},
		ProviderHealth: map[string]domain.ProviderHealth{
			"hybrid": {LocalHealthy: false, APIHealthy: true},
		},
	}).Routes()

	request := httptest.NewRequest(http.MethodGet, "/api/provider-health?id=hybrid", nil)
	request.Header.Set("Authorization", "Bearer local-test-secret")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"route":"api"`) {
		t.Fatalf("expected api fallback route, got %s", response.Body.String())
	}
	if strings.Contains(response.Body.String(), "local-test-secret") {
		t.Fatalf("provider health leaked auth secret: %s", response.Body.String())
	}
}

func TestProviderHealthChecksLocalRuntimeEndpoint(t *testing.T) {
	runtime := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer runtime.Close()
	provider, err := domain.NewModelProvider("local", "Local", domain.ProviderModeLocal, runtime.URL, "llama", "", "ollama")
	if err != nil {
		t.Fatal(err)
	}
	server := (&httpapi.Server{
		AuthSecret: "local-test-secret",
		Providers:  []domain.ModelProvider{provider},
	}).Routes()

	request := httptest.NewRequest(http.MethodGet, "/api/provider-health?id=local", nil)
	request.Header.Set("Authorization", "Bearer local-test-secret")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"route":"local"`) {
		t.Fatalf("expected local route, got %s", response.Body.String())
	}
}

func TestMVPFlowCreatesProviderAgentReadsMemoryAndWritesHandoff(t *testing.T) {
	root := t.TempDir()
	docDir := filepath.Join(root, "DOC")
	if err := os.MkdirAll(docDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(docDir, "STATE.md"), []byte("phase 1 state"), 0o600); err != nil {
		t.Fatal(err)
	}
	handoffDir := filepath.Join(root, ".agent_handoff")
	api := (&httpapi.Server{
		AuthSecret: "local-test-secret",
		DocRoot:    root,
		HandoffDir: handoffDir,
		MemoryFiles: []domain.MemoryFile{
			{Path: "DOC/STATE.md", Label: "State"},
		},
	}).Routes()

	postJSON(t, api, "/api/providers", `{"id":"api-default","name":"API Default","mode":"api","endpoint":"https://api.example.test","model":"gpt","secret_ref":"PF_AI_API_KEY"}`, http.StatusCreated)
	postJSON(t, api, "/api/agents", `{"id":"ceo","name":"CEO","role":"orchestration","seniority":"Senior","provider_id":"api-default","description":"orchestrates"}`, http.StatusCreated)

	memoryRequest := httptest.NewRequest(http.MethodGet, "/api/memory?path=DOC/STATE.md", nil)
	memoryRequest.Header.Set("Authorization", "Bearer local-test-secret")
	memoryResponse := httptest.NewRecorder()
	api.ServeHTTP(memoryResponse, memoryRequest)
	if memoryResponse.Code != http.StatusOK {
		t.Fatalf("expected memory 200, got %d: %s", memoryResponse.Code, memoryResponse.Body.String())
	}
	if !strings.Contains(memoryResponse.Body.String(), "phase 1 state") {
		t.Fatalf("expected memory content, got %s", memoryResponse.Body.String())
	}

	postJSON(t, api, "/api/handoffs", `{"sender":"CEO","recipient":"DEV_BACKEND","task_ref":"PHASE1-BACKEND-001","intent":"PHASE_KICKOFF","payload":{"summary":"validate flow"}}`, http.StatusCreated)

	files, err := os.ReadDir(handoffDir)
	if err != nil {
		t.Fatalf("expected handoff dir: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected one handoff file, got %d", len(files))
	}
	body, err := os.ReadFile(filepath.Join(handoffDir, files[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(strings.ToLower(string(body)), "local-test-secret") {
		t.Fatalf("handoff leaked auth secret: %s", string(body))
	}
}

func postJSON(t *testing.T, api http.Handler, path string, payload string, expectedStatus int) {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(payload))
	request.Header.Set("Authorization", "Bearer local-test-secret")
	response := httptest.NewRecorder()
	api.ServeHTTP(response, request)
	if response.Code != expectedStatus {
		var parsed map[string]any
		_ = json.Unmarshal(response.Body.Bytes(), &parsed)
		t.Fatalf("expected %d for %s, got %d: %v", expectedStatus, path, response.Code, parsed)
	}
}
