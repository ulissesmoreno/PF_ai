CREATE TABLE IF NOT EXISTS agent_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    agent_name TEXT NOT NULL UNIQUE,
    model TEXT NOT NULL,
    tier INTEGER NOT NULL,
    active INTEGER NOT NULL DEFAULT 1,
    updated_at TEXT NOT NULL,
    updated_by TEXT
);

CREATE INDEX IF NOT EXISTS idx_agent_configs_active ON agent_configs(active, agent_name);

INSERT INTO agent_configs(agent_name, model, tier, active, updated_at, updated_by) VALUES
('CEO', 'qwen3:4b', 3, 1, CURRENT_TIMESTAMP, 'migration'),
('BA', 'qwen3:4b', 3, 1, CURRENT_TIMESTAMP, 'migration'),
('CTO', 'qwen3:4b', 3, 1, CURRENT_TIMESTAMP, 'migration'),
('DEV_FRONTEND', 'qwen2.5-coder:7b', 2, 1, CURRENT_TIMESTAMP, 'migration'),
('DEV_BACKEND', 'qwen2.5-coder:7b', 2, 1, CURRENT_TIMESTAMP, 'migration'),
('DBA', 'qwen2.5-coder:7b', 2, 1, CURRENT_TIMESTAMP, 'migration'),
('DS_ML', 'qwen2.5-coder:7b', 2, 1, CURRENT_TIMESTAMP, 'migration'),
('SECURITY', 'qwen3:4b', 3, 1, CURRENT_TIMESTAMP, 'migration'),
('QA', 'qwen3:4b', 3, 1, CURRENT_TIMESTAMP, 'migration'),
('DATA_ENGINEER', 'qwen2.5-coder:7b', 2, 1, CURRENT_TIMESTAMP, 'migration'),
('PM', 'qwen3:4b', 3, 1, CURRENT_TIMESTAMP, 'migration'),
('UX_RESEARCHER', 'phi4-mini:latest', 1, 1, CURRENT_TIMESTAMP, 'migration'),
('WRITER', 'phi4-mini:latest', 1, 1, CURRENT_TIMESTAMP, 'migration'),
('DOCUMENTATION', 'qwen3:4b', 3, 1, CURRENT_TIMESTAMP, 'migration'),
('ARTIST', 'phi4-mini:latest', 1, 1, CURRENT_TIMESTAMP, 'migration'),
('DEVOPS', 'qwen2.5-coder:7b', 2, 1, CURRENT_TIMESTAMP, 'migration'),
('CODE_REVIEWER', 'qwen3:4b', 3, 1, CURRENT_TIMESTAMP, 'migration'),
('CMO', 'phi4-mini:latest', 1, 1, CURRENT_TIMESTAMP, 'migration')
ON CONFLICT(agent_name) DO NOTHING;
