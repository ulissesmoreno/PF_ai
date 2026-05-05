package main

import (
	"log/slog"
	"net/http"
	"os"

	"pf-ai/src/domain"
	"pf-ai/src/infrastructure/httpapi"
)

func main() {
	port := getenv("PF_AI_HTTP_PORT", "8080")
	authSecret := os.Getenv("PF_AI_AUTH_SECRET")
	docRoot := getenv("PF_AI_DOC_ROOT", ".")
	handoffDir := getenv("PF_AI_HANDOFF_DIR", ".agent_handoff")

	server := &httpapi.Server{
		AuthSecret: authSecret,
		DocRoot:    docRoot,
		HandoffDir: handoffDir,
		MemoryFiles: []domain.MemoryFile{
			{Path: "DOC/PLAN.md", Label: "Plan"},
			{Path: "DOC/ROADMAP.md", Label: "Roadmap"},
			{Path: "DOC/CONTEXT.md", Label: "Context"},
			{Path: "DOC/STATE.md", Label: "State"},
			{Path: "QUESTIONS.md", Label: "Questions"},
		},
	}

	slog.Info("starting pf-ai service", "port", port)
	if err := http.ListenAndServe(":"+port, server.Routes()); err != nil {
		slog.Error("pf-ai service stopped", "error", err)
		os.Exit(1)
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
