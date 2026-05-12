CREATE TABLE IF NOT EXISTS document_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_id INTEGER NOT NULL,
    import_id INTEGER NOT NULL,
    entry_type TEXT NOT NULL,
    heading TEXT,
    ordinal INTEGER NOT NULL,
    content TEXT NOT NULL,
    content_hash TEXT NOT NULL,
    created_at TEXT NOT NULL,
    FOREIGN KEY (document_id) REFERENCES documents(id),
    FOREIGN KEY (import_id) REFERENCES document_imports(id)
);

CREATE INDEX IF NOT EXISTS idx_document_entries_document_id ON document_entries(document_id);
