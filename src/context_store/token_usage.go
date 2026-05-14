package context_store

import (
	"fmt"
	"strings"
)

type TokenUsage struct {
	ID               int64  `json:"id"`
	ProjectID        string `json:"project_id,omitempty"`
	CardID           string `json:"card_id,omitempty"`
	AgentName        string `json:"agent_name"`
	Model            string `json:"model"`
	Tier             int    `json:"tier"`
	PromptTokens     int    `json:"prompt_tokens"`
	CompletionTokens int    `json:"completion_tokens"`
	TotalTokens      int    `json:"total_tokens"`
	LatencyMS        int64  `json:"latency_ms"`
	CreatedAt        string `json:"created_at"`
}

type TokenUsageFilter struct {
	ProjectID string
	CardID    string
	AgentName string
	Model     string
	From      string
	To        string
}

func (s *Store) SaveTokenUsage(usage TokenUsage) error {
	if strings.TrimSpace(usage.AgentName) == "" {
		return fmt.Errorf("agent_name obrigatorio")
	}
	if strings.TrimSpace(usage.Model) == "" {
		usage.Model = "unknown"
	}
	if usage.TotalTokens == 0 {
		usage.TotalTokens = usage.PromptTokens + usage.CompletionTokens
	}
	_, err := s.db.Exec(
		`INSERT INTO token_usage(
			project_id, card_id, agent_name, model, tier, prompt_tokens, completion_tokens,
			total_tokens, latency_ms, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.projectIDOrNil(), nilIfEmpty(usage.CardID), strings.ToUpper(usage.AgentName), usage.Model,
		usage.Tier, usage.PromptTokens, usage.CompletionTokens, usage.TotalTokens, usage.LatencyMS, now(),
	)
	if err != nil {
		return fmt.Errorf("salvar token usage: %w", err)
	}
	if strings.TrimSpace(usage.CardID) != "" {
		summary := fmt.Sprintf("[TOKEN_USAGE] prompt=%d completion=%d total=%d model=%s tier=%d latency_ms=%d",
			usage.PromptTokens, usage.CompletionTokens, usage.TotalTokens, usage.Model, usage.Tier, usage.LatencyMS)
		if err := s.AddCardComment(usage.CardID, "SYSTEM", "token_usage", summary); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) SumTokenUsage(filter TokenUsageFilter) (TokenUsage, error) {
	row := s.db.QueryRow(
		`SELECT
			COALESCE(SUM(prompt_tokens), 0),
			COALESCE(SUM(completion_tokens), 0),
			COALESCE(SUM(total_tokens), 0),
			COALESCE(MAX(model), ''),
			COALESCE(MAX(agent_name), '')
		   FROM token_usage
		  WHERE (? = '' OR project_id = ?)
		    AND (? = '' OR card_id = ?)
		    AND (? = '' OR agent_name = ?)
		    AND (? = '' OR model = ?)
		    AND (? = '' OR created_at >= ?)
		    AND (? = '' OR created_at <= ?)`,
		filter.ProjectID, filter.ProjectID,
		filter.CardID, filter.CardID,
		strings.ToUpper(filter.AgentName), strings.ToUpper(filter.AgentName),
		filter.Model, filter.Model,
		filter.From, filter.From,
		filter.To, filter.To,
	)
	var usage TokenUsage
	if err := row.Scan(&usage.PromptTokens, &usage.CompletionTokens, &usage.TotalTokens, &usage.Model, &usage.AgentName); err != nil {
		return TokenUsage{}, fmt.Errorf("somar token usage: %w", err)
	}
	usage.ProjectID = filter.ProjectID
	usage.CardID = filter.CardID
	return usage, nil
}

func (s *Store) ListTokenUsageByCard(cardID string) ([]TokenUsage, error) {
	rows, err := s.db.Query(
		`SELECT id, COALESCE(project_id, ''), COALESCE(card_id, ''), agent_name, model, tier,
		        prompt_tokens, completion_tokens, total_tokens, COALESCE(latency_ms, 0), created_at
		   FROM token_usage
		  WHERE card_id = ?
		  ORDER BY created_at`,
		cardID,
	)
	if err != nil {
		return nil, fmt.Errorf("listar token usage por card: %w", err)
	}
	defer rows.Close()
	var usages []TokenUsage
	for rows.Next() {
		var usage TokenUsage
		if err := rows.Scan(&usage.ID, &usage.ProjectID, &usage.CardID, &usage.AgentName, &usage.Model, &usage.Tier, &usage.PromptTokens, &usage.CompletionTokens, &usage.TotalTokens, &usage.LatencyMS, &usage.CreatedAt); err != nil {
			return nil, err
		}
		usages = append(usages, usage)
	}
	return usages, rows.Err()
}

func nilIfEmpty(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}
