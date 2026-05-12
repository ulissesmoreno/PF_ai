CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT PRIMARY KEY,
    applied_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS documents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    document_type TEXT NOT NULL,
    owner TEXT,
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS document_imports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_id INTEGER NOT NULL,
    content_hash TEXT NOT NULL,
    imported_at TEXT NOT NULL,
    FOREIGN KEY (document_id) REFERENCES documents(id)
);

CREATE TABLE IF NOT EXISTS document_sections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_id INTEGER NOT NULL,
    import_id INTEGER NOT NULL,
    heading TEXT NOT NULL,
    level INTEGER NOT NULL,
    ordinal INTEGER NOT NULL,
    content TEXT NOT NULL,
    content_hash TEXT NOT NULL,
    created_at TEXT NOT NULL,
    FOREIGN KEY (document_id) REFERENCES documents(id),
    FOREIGN KEY (import_id) REFERENCES document_imports(id)
);

CREATE TABLE IF NOT EXISTS context_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    entry_type TEXT NOT NULL,
    document_type TEXT NOT NULL,
    section TEXT,
    title TEXT,
    content TEXT NOT NULL,
    source_agent TEXT,
    task_ref TEXT,
    tags TEXT,
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS planning_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item_type TEXT NOT NULL,
    reference TEXT,
    title TEXT,
    status TEXT,
    priority TEXT,
    content TEXT NOT NULL,
    source_agent TEXT,
    task_ref TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS handoff_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT,
    sender TEXT,
    recipient TEXT,
    task_ref TEXT,
    intent TEXT,
    payload TEXT,
    raw_json TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS agent_outputs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    agent_name TEXT NOT NULL,
    task_id TEXT NOT NULL,
    action TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS test_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    test_name TEXT,
    status TEXT NOT NULL,
    command TEXT,
    output TEXT,
    source_agent TEXT,
    task_ref TEXT,
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS knowledge_chunks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_type TEXT NOT NULL,
    source_id INTEGER NOT NULL,
    document_type TEXT,
    heading TEXT,
    content TEXT NOT NULL,
    content_hash TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS read_context_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_table TEXT NOT NULL,
    source_id INTEGER NOT NULL,
    document_type TEXT,
    entry_type TEXT,
    task_ref TEXT,
    agent_name TEXT,
    title TEXT,
    heading TEXT,
    content TEXT NOT NULL,
    content_hash TEXT NOT NULL,
    projected_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_document_sections_document_id ON document_sections(document_id);
CREATE INDEX IF NOT EXISTS idx_document_imports_document_id ON document_imports(document_id);
CREATE INDEX IF NOT EXISTS idx_context_entries_lookup ON context_entries(document_type, entry_type, task_ref);
CREATE INDEX IF NOT EXISTS idx_planning_items_lookup ON planning_items(item_type, status, task_ref);
CREATE INDEX IF NOT EXISTS idx_handoff_events_lookup ON handoff_events(recipient, task_ref, intent);
CREATE INDEX IF NOT EXISTS idx_agent_outputs_lookup ON agent_outputs(agent_name, task_id, action);
CREATE INDEX IF NOT EXISTS idx_knowledge_chunks_lookup ON knowledge_chunks(document_type, source_type);
CREATE INDEX IF NOT EXISTS idx_read_context_items_lookup ON read_context_items(document_type, entry_type, task_ref, agent_name);
