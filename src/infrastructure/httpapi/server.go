package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"pf-ai/src/domain"
	"pf-ai/src/infrastructure/filememory"
	"pf-ai/src/infrastructure/handoff"
)

type Server struct {
	AuthSecret     string
	DocRoot        string
	HandoffDir     string
	Agents         []domain.AgentDefinition
	Providers      []domain.ModelProvider
	ProviderHealth map[string]domain.ProviderHealth
	PlanningCards  []domain.PlanningCard
	Projects       []domain.ProjectRegistration
	MemoryFiles    []domain.MemoryFile
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.Handle("/api/agents", s.protected(http.HandlerFunc(s.agents)))
	mux.Handle("/api/providers", s.protected(http.HandlerFunc(s.providers)))
	mux.Handle("/api/provider-health", s.protected(http.HandlerFunc(s.providerHealth)))
	mux.Handle("/api/planning-cards", s.protected(http.HandlerFunc(s.planningCards)))
	mux.Handle("/api/projects", s.protected(http.HandlerFunc(s.projects)))
	mux.Handle("/api/mcp-envelope", s.protected(http.HandlerFunc(s.mcpEnvelope)))
	mux.Handle("/api/memory", s.protected(http.HandlerFunc(s.memory)))
	mux.Handle("/api/handoffs", s.protected(http.HandlerFunc(s.handoffs)))
	return cors(mux)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) agents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.Agents)
	case http.MethodPost:
		var request agentRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		agent, err := domain.NewAgentDefinition(request.ID, request.Name, request.Role, domain.Seniority(request.Seniority), request.ProviderID, request.Description)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		s.Agents = append(s.Agents, agent)
		writeJSON(w, http.StatusCreated, agent)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (s *Server) providers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.Providers)
	case http.MethodPost:
		var request providerRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		provider, err := domain.NewModelProvider(request.ID, request.Name, domain.ProviderMode(request.Mode), request.Endpoint, request.Model, request.SecretRef, request.LocalRuntime)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		s.Providers = append(s.Providers, provider)
		writeJSON(w, http.StatusCreated, provider)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (s *Server) providerHealth(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		for _, provider := range s.Providers {
			if provider.ID == id {
				health := s.healthForProvider(r.Context(), provider)
				writeJSON(w, http.StatusOK, domain.DecideProviderRoute(provider, health))
				return
			}
		}
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "provider not found"})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (s *Server) healthForProvider(ctx context.Context, provider domain.ModelProvider) domain.ProviderHealth {
	if s.ProviderHealth != nil {
		if health, ok := s.ProviderHealth[provider.ID]; ok {
			return health
		}
	}
	localHealthy := false
	if provider.LocalRuntime != "" && provider.Endpoint != "" {
		localHealthy = checkEndpoint(ctx, provider.Endpoint)
	}
	return domain.ProviderHealth{
		LocalHealthy: localHealthy,
		APIHealthy:   provider.SecretRef != "",
	}
}

func checkEndpoint(ctx context.Context, endpoint string) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return false
	}
	response, err := client.Do(request)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	return response.StatusCode >= 200 && response.StatusCode < 300
}

func (s *Server) planningCards(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.PlanningCards)
	case http.MethodPost:
		var request planningCardRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		card, err := domain.NewPlanningCard(request.ID, request.Title, request.Owner, domain.CardStatus(request.Status), domain.CardPriority(request.Priority), request.Phase, request.TaskRef)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		s.PlanningCards = append(s.PlanningCards, card)
		writeJSON(w, http.StatusCreated, card)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (s *Server) projects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.Projects)
	case http.MethodPost:
		var request projectRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		project, err := domain.NewProjectRegistration(request.ID, request.Name, request.Description, request.Audience, request.TechnicalStack, request.Onboarding)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		project.ID = s.uniqueProjectID(project.ID)
		s.Projects = append(s.Projects, project)
		writeJSON(w, http.StatusCreated, project)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (s *Server) mcpEnvelope(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var request handoffRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		item, err := domain.NewHandoff(request.Sender, request.Recipient, request.TaskRef, domain.HandoffIntent(request.Intent), request.Payload, nil)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, domain.NewMCPEnvelope(item))
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (s *Server) memory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		path := r.URL.Query().Get("path")
		for _, candidate := range s.MemoryFiles {
			if candidate.Path == path {
				reader := filememory.NewReader(s.DocRoot, memoryPaths(s.MemoryFiles))
				content, err := reader.ReadMemory(r.Context(), candidate)
				if err != nil {
					writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "memory read failed"})
					return
				}
				writeJSON(w, http.StatusOK, map[string]string{"path": candidate.Path, "label": candidate.Label, "content": content})
				return
			}
		}
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "memory file not found"})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (s *Server) handoffs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var request handoffRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		item, err := domain.NewHandoff(request.Sender, request.Recipient, request.TaskRef, domain.HandoffIntent(request.Intent), request.Payload, nil)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writer := handoff.FileWriter{Dir: s.HandoffDir}
		path, err := writer.WriteHandoff(context.Background(), item)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "handoff write failed"})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"path": path})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (s *Server) protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.AuthSecret == "" || r.Header.Get("Authorization") != "Bearer "+s.AuthSecret {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

type agentRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	Seniority   string `json:"seniority"`
	ProviderID  string `json:"provider_id"`
	Description string `json:"description"`
}

type providerRequest struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Mode         string `json:"mode"`
	Endpoint     string `json:"endpoint"`
	Model        string `json:"model"`
	SecretRef    string `json:"secret_ref"`
	LocalRuntime string `json:"local_runtime"`
}

type handoffRequest struct {
	Sender    string         `json:"sender"`
	Recipient string         `json:"recipient"`
	TaskRef   string         `json:"task_ref"`
	Intent    string         `json:"intent"`
	Payload   map[string]any `json:"payload"`
}

type planningCardRequest struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Owner    string `json:"owner"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
	Phase    string `json:"phase"`
	TaskRef  string `json:"task_ref"`
}

type projectRequest struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Audience       string         `json:"audience"`
	TechnicalStack string         `json:"technical_stack"`
	Onboarding     map[string]any `json:"onboarding"`
}

func (s *Server) uniqueProjectID(base string) string {
	candidate := base
	suffix := 2
	for {
		used := false
		for _, project := range s.Projects {
			if project.ID == candidate {
				used = true
				break
			}
		}
		if !used {
			return candidate
		}
		candidate = base + "-" + strconv.Itoa(suffix)
		suffix++
	}
}

func memoryPaths(files []domain.MemoryFile) []string {
	paths := make([]string, 0, len(files))
	for _, file := range files {
		paths = append(paths, file.Path)
	}
	return paths
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
