package context_store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	db            *sql.DB
	workspaceRoot string
}

type ContextEntry struct {
	EntryType    string   `json:"entry_type"`
	DocumentType string   `json:"document_type"`
	Section      string   `json:"section"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	SourceAgent  string   `json:"source_agent"`
	TaskRef      string   `json:"task_ref"`
	Tags         []string `json:"tags"`
}

type PlanningItem struct {
	ItemType    string `json:"item_type"`
	Reference   string `json:"reference"`
	Title       string `json:"title"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	Content     string `json:"content"`
	SourceAgent string `json:"source_agent"`
	TaskRef     string `json:"task_ref"`
}

type TestRecord struct {
	TestName    string `json:"test_name"`
	Status      string `json:"status"`
	Command     string `json:"command"`
	Output      string `json:"output"`
	SourceAgent string `json:"source_agent"`
	TaskRef     string `json:"task_ref"`
}

type HandoffEvent struct {
	Path      string
	Sender    string
	Recipient string
	TaskRef   string
	Intent    string
	Payload   string
	RawJSON   string
}

func Open(dbPath, migrationsDir, workspaceRoot string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("criar diretório do banco: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("abrir sqlite: %w", err)
	}

	store := &Store{db: db, workspaceRoot: workspaceRoot}
	if err := store.applyMigrations(migrationsDir); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) applyMigrations(migrationsDir string) error {
	if strings.TrimSpace(migrationsDir) == "" {
		return nil
	}

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("listar migrations: %w", err)
	}
	sort.Strings(files)

	for _, file := range files {
		version := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))

		var exists int
		err := s.db.QueryRow("SELECT COUNT(1) FROM schema_migrations WHERE version = ?", version).Scan(&exists)
		if err == nil && exists > 0 {
			continue
		}

		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("ler migration %s: %w", file, err)
		}

		tx, err := s.db.Begin()
		if err != nil {
			return fmt.Errorf("iniciar migration %s: %w", version, err)
		}
		if _, err := tx.Exec(string(data)); err != nil {
			tx.Rollback() //nolint
			return fmt.Errorf("executar migration %s: %w", version, err)
		}
		if _, err := tx.Exec(
			"INSERT INTO schema_migrations(version, applied_at) VALUES (?, ?)",
			version,
			now(),
		); err != nil {
			tx.Rollback() //nolint
			return fmt.Errorf("registrar migration %s: %w", version, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", version, err)
		}
	}

	return nil
}

func (s *Store) SaveContextEntry(entry ContextEntry) error {
	if strings.TrimSpace(entry.Content) == "" {
		return fmt.Errorf("context entry sem conteúdo")
	}
	if entry.EntryType == "" {
		entry.EntryType = "context"
	}
	if entry.DocumentType == "" {
		entry.DocumentType = "CONTEXT"
	}

	tags, err := json.Marshal(entry.Tags)
	if err != nil {
		return fmt.Errorf("serializar tags: %w", err)
	}

	res, err := s.db.Exec(
		`INSERT INTO context_entries(
			entry_type, document_type, section, title, content, source_agent, task_ref, tags, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entry.EntryType,
		strings.ToUpper(entry.DocumentType),
		entry.Section,
		entry.Title,
		entry.Content,
		entry.SourceAgent,
		entry.TaskRef,
		string(tags),
		now(),
	)
	if err != nil {
		return fmt.Errorf("salvar contexto: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obter id de contexto: %w", err)
	}
	return s.projectContextItem("context_entries", id, strings.ToUpper(entry.DocumentType), entry.EntryType, entry.TaskRef, entry.SourceAgent, entry.Title, entry.Section, entry.Content)
}

func (s *Store) SavePlanningItem(item PlanningItem) error {
	if strings.TrimSpace(item.Content) == "" {
		return fmt.Errorf("planning item sem conteúdo")
	}
	if item.ItemType == "" {
		item.ItemType = "plan"
	}

	res, err := s.db.Exec(
		`INSERT INTO planning_items(
			item_type, reference, title, status, priority, content, source_agent, task_ref, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.ItemType,
		item.Reference,
		item.Title,
		item.Status,
		item.Priority,
		item.Content,
		item.SourceAgent,
		item.TaskRef,
		now(),
		now(),
	)
	if err != nil {
		return fmt.Errorf("salvar planejamento: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obter id de planejamento: %w", err)
	}
	return s.projectContextItem("planning_items", id, "PLAN", item.ItemType, item.TaskRef, item.SourceAgent, item.Title, item.Reference, item.Content)
}

func (s *Store) SaveTestRecord(record TestRecord) error {
	if strings.TrimSpace(record.Status) == "" {
		return fmt.Errorf("test record sem status")
	}

	res, err := s.db.Exec(
		`INSERT INTO test_records(
			test_name, status, command, output, source_agent, task_ref, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		record.TestName,
		record.Status,
		record.Command,
		record.Output,
		record.SourceAgent,
		record.TaskRef,
		now(),
	)
	if err != nil {
		return fmt.Errorf("salvar teste: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obter id de teste: %w", err)
	}
	content := strings.TrimSpace(record.Status + "\n" + record.Command + "\n" + record.Output)
	return s.projectContextItem("test_records", id, "TESTS", "record_test", record.TaskRef, record.SourceAgent, record.TestName, "", content)
}

func (s *Store) SaveAgentOutput(agentName, taskID, action, content string) error {
	if strings.TrimSpace(content) == "" {
		return nil
	}
	if action == "" {
		action = "note"
	}

	res, err := s.db.Exec(
		`INSERT INTO agent_outputs(agent_name, task_id, action, content, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		strings.ToUpper(agentName),
		taskID,
		action,
		content,
		now(),
	)
	if err != nil {
		return fmt.Errorf("salvar saída do agente: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obter id de saída do agente: %w", err)
	}
	return s.projectContextItem("agent_outputs", id, "", action, taskID, agentName, action, "", content)
}

func (s *Store) SaveHandoffEvent(event HandoffEvent) error {
	_, err := s.db.Exec(
		`INSERT INTO handoff_events(
			path, sender, recipient, task_ref, intent, payload, raw_json, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		event.Path,
		event.Sender,
		event.Recipient,
		event.TaskRef,
		event.Intent,
		event.Payload,
		event.RawJSON,
		now(),
	)
	if err != nil {
		return fmt.Errorf("salvar handoff event: %w", err)
	}
	return nil
}

func (s *Store) QueryContextForHandoff(agentName, taskRef string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 12
	}

	rows, err := s.db.Query(
		`SELECT content
		   FROM read_context_items
		  WHERE (? = '' OR task_ref = ? OR task_ref = '' OR task_ref IS NULL)
		    AND NOT (
				source_table = 'document_sections'
				AND document_type IN ('TASKS', 'STATE', 'CONTEXT', 'PLAYBOOK', 'TESTS', 'VERSIONS', 'RETROSPECTIVE')
			)
		  ORDER BY projected_at DESC
		  LIMIT ?`,
		taskRef,
		taskRef,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("consultar contexto para %s: %w", agentName, err)
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return nil, err
		}
		result = append(result, content)
	}
	return result, rows.Err()
}

func (s *Store) ProjectStarted() (bool, error) {
	rows, err := s.db.Query(
		`SELECT ds.content
		   FROM document_sections ds
		   JOIN documents d ON d.id = ds.document_id
		   JOIN document_imports di ON di.id = ds.import_id
		  WHERE d.document_type = 'PROJECT'
		    AND di.id = (
				SELECT MAX(di2.id)
				  FROM document_imports di2
				  JOIN documents d2 ON d2.id = di2.document_id
				 WHERE d2.document_type = 'PROJECT'
			)
		  ORDER BY ds.ordinal`,
	)
	if err != nil {
		return false, fmt.Errorf("consultar PROJECT persistido: %w", err)
	}
	defer rows.Close()

	var parts []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return false, err
		}
		parts = append(parts, content)
	}
	if err := rows.Err(); err != nil {
		return false, err
	}
	if len(parts) == 0 {
		return false, nil
	}

	content := strings.ToLower(strings.Join(parts, "\n"))
	placeholderSignals := []string{
		"[project_name]",
		"[nome do projeto]",
		"[descrição",
		"[descricao",
		"[resumo",
		"[quem usa",
		"[placeholder",
		"[preenchido",
		"[filled",
	}
	for _, signal := range placeholderSignals {
		if strings.Contains(content, signal) {
			return false, nil
		}
	}

	return true, nil
}

func (s *Store) projectContextItem(sourceTable string, sourceID int64, documentType, entryType, taskRef, agentName, title, heading, content string) error {
	if strings.TrimSpace(content) == "" {
		return nil
	}
	_, err := s.db.Exec(
		`INSERT INTO read_context_items(
			source_table, source_id, document_type, entry_type, task_ref, agent_name,
			title, heading, content, content_hash, projected_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		sourceTable,
		sourceID,
		documentType,
		entryType,
		taskRef,
		strings.ToUpper(agentName),
		title,
		heading,
		content,
		sha(content),
		now(),
	)
	if err != nil {
		return fmt.Errorf("projetar read model: %w", err)
	}
	return nil
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}
