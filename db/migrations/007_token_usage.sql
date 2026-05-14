CREATE TABLE IF NOT EXISTS token_usage (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id TEXT REFERENCES projects(id),
    card_id TEXT REFERENCES cards(id),
    agent_name TEXT NOT NULL,
    model TEXT NOT NULL,
    tier INTEGER NOT NULL,
    prompt_tokens INTEGER NOT NULL,
    completion_tokens INTEGER NOT NULL,
    total_tokens INTEGER NOT NULL,
    latency_ms INTEGER,
    created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_token_usage_project ON token_usage(project_id, created_at);
CREATE INDEX IF NOT EXISTS idx_token_usage_card ON token_usage(card_id);
CREATE INDEX IF NOT EXISTS idx_token_usage_agent ON token_usage(agent_name, created_at);
CREATE INDEX IF NOT EXISTS idx_token_usage_model ON token_usage(model, created_at);
