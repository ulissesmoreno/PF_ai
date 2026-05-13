CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    domain TEXT,
    description TEXT,
    target_audience TEXT,
    main_objective TEXT,
    stack_backend TEXT,
    stack_frontend TEXT,
    stack_database TEXT,
    stack_infra TEXT,
    stack_notes TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    workspace_root TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_slug ON projects(slug);

ALTER TABLE context_entries ADD COLUMN project_id TEXT REFERENCES projects(id);
ALTER TABLE planning_items ADD COLUMN project_id TEXT REFERENCES projects(id);
ALTER TABLE test_records ADD COLUMN project_id TEXT REFERENCES projects(id);
ALTER TABLE agent_outputs ADD COLUMN project_id TEXT REFERENCES projects(id);
ALTER TABLE handoff_events ADD COLUMN project_id TEXT REFERENCES projects(id);
ALTER TABLE documents ADD COLUMN project_id TEXT REFERENCES projects(id);
ALTER TABLE document_imports ADD COLUMN project_id TEXT REFERENCES projects(id);
ALTER TABLE document_sections ADD COLUMN project_id TEXT REFERENCES projects(id);
ALTER TABLE document_entries ADD COLUMN project_id TEXT REFERENCES projects(id);
ALTER TABLE knowledge_chunks ADD COLUMN project_id TEXT REFERENCES projects(id);
ALTER TABLE read_context_items ADD COLUMN project_id TEXT REFERENCES projects(id);

CREATE INDEX IF NOT EXISTS idx_context_entries_project ON context_entries(project_id, document_type, entry_type, task_ref);
CREATE INDEX IF NOT EXISTS idx_planning_items_project ON planning_items(project_id, item_type, status, task_ref);
CREATE INDEX IF NOT EXISTS idx_test_records_project ON test_records(project_id, task_ref, status);
CREATE INDEX IF NOT EXISTS idx_agent_outputs_project ON agent_outputs(project_id, agent_name, task_id, action);
CREATE INDEX IF NOT EXISTS idx_handoff_events_project ON handoff_events(project_id, recipient, task_ref, intent);
CREATE INDEX IF NOT EXISTS idx_documents_project ON documents(project_id, document_type);
CREATE INDEX IF NOT EXISTS idx_document_imports_project ON document_imports(project_id, document_id);
CREATE INDEX IF NOT EXISTS idx_document_sections_project ON document_sections(project_id, document_id);
CREATE INDEX IF NOT EXISTS idx_document_entries_project ON document_entries(project_id, document_id);
CREATE INDEX IF NOT EXISTS idx_knowledge_chunks_project ON knowledge_chunks(project_id, document_type, source_type);
CREATE INDEX IF NOT EXISTS idx_read_context_items_project ON read_context_items(project_id, document_type, entry_type, task_ref, agent_name);
