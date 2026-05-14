package context_store

import (
	"database/sql"
	"fmt"
	"strings"
)

type AgentConfig struct {
	ID        int64  `json:"id"`
	AgentName string `json:"agent_name"`
	Model     string `json:"model"`
	Tier      int    `json:"tier"`
	Active    bool   `json:"active"`
	UpdatedAt string `json:"updated_at"`
	UpdatedBy string `json:"updated_by"`
}

func (s *Store) GetAgentConfig(name string) (*AgentConfig, error) {
	row := s.db.QueryRow(
		`SELECT id, agent_name, model, tier, active, updated_at, updated_by
		   FROM agent_configs
		  WHERE agent_name = ? AND active = 1`,
		strings.ToUpper(strings.TrimSpace(name)),
	)
	return scanAgentConfig(row)
}

func (s *Store) ListAgentConfigs() ([]AgentConfig, error) {
	rows, err := s.db.Query(
		`SELECT id, agent_name, model, tier, active, updated_at, updated_by
		   FROM agent_configs
		  ORDER BY agent_name`,
	)
	if err != nil {
		return nil, fmt.Errorf("listar agent configs: %w", err)
	}
	defer rows.Close()
	var configs []AgentConfig
	for rows.Next() {
		config, err := scanAgentConfig(rows)
		if err != nil {
			return nil, err
		}
		configs = append(configs, *config)
	}
	return configs, rows.Err()
}

func (s *Store) UpsertAgentConfig(config AgentConfig) (*AgentConfig, error) {
	config.AgentName = strings.ToUpper(strings.TrimSpace(config.AgentName))
	config.Model = strings.TrimSpace(config.Model)
	if config.AgentName == "" || config.Model == "" {
		return nil, fmt.Errorf("agent_name e model sao obrigatorios")
	}
	if config.Tier <= 0 {
		config.Tier = 2
	}
	if config.UpdatedBy == "" {
		config.UpdatedBy = "api"
	}
	active := 0
	if config.Active {
		active = 1
	}
	_, err := s.db.Exec(
		`INSERT INTO agent_configs(agent_name, model, tier, active, updated_at, updated_by)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT(agent_name) DO UPDATE SET
			model = excluded.model,
			tier = excluded.tier,
			active = excluded.active,
			updated_at = excluded.updated_at,
			updated_by = excluded.updated_by`,
		config.AgentName, config.Model, config.Tier, active, now(), config.UpdatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("salvar agent config: %w", err)
	}
	return s.GetAgentConfig(config.AgentName)
}

func scanAgentConfig(row scanner) (*AgentConfig, error) {
	var config AgentConfig
	var active int
	var updatedBy sql.NullString
	if err := row.Scan(&config.ID, &config.AgentName, &config.Model, &config.Tier, &active, &config.UpdatedAt, &updatedBy); err != nil {
		return nil, err
	}
	config.Active = active == 1
	config.UpdatedBy = updatedBy.String
	return &config, nil
}
