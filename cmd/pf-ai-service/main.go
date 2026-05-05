package main

import (
	"log/slog"
	"net/http"
	"os"

	"pf-ai/src/infrastructure/httpapi"
)

func main() {
	port := getenv("PF_AI_HTTP_PORT", "8080")
	authSecret := os.Getenv("PF_AI_AUTH_SECRET")

	server := &httpapi.Server{AuthSecret: authSecret}

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
