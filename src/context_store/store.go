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

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type Store struct {
	db              *sql.DB
	workspaceRoot   string
	ActiveProjectID string
}

type Project struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Domain         string `json:"domain"`
	Description    string `json:"description"`
	TargetAudience string `json:"target_audience"`
	MainObjective  string `json:"main_objective"`
	StackBackend   string `json:"stack_backend"`
	StackFrontend  string `json:"stack_frontend"`
	StackDatabase  string `json:"stack_database"`
	StackInfra     string `json:"stack_infra"`
	StackNotes     string `json:"stack_notes"`
	Status         string `json:"status"`
	WorkspaceRoot  string `json:"workspace_root"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type ContextEntry struct {
	ProjectID    string   `json:"project_id,omitempty"`
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
	ProjectID   string `json:"project_id,omitempty"`
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
	ProjectID   string `json:"project_id,omitempty"`
	TestName    string `json:"test_name"`
	Status      string `json:"status"`
	Command     string `json:"command"`
	Output      string `json:"output"`
	SourceAgent string `json:"source_agent"`
	TaskRef     string `json:"task_ref"`
}

type HandoffEvent struct {
	ProjectID string
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
	if err := store.loadActiveProjectID(); err != nil {
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

func (s *Store) loadActiveProjectID() error {
	var id string
	err := s.db.QueryRow(
		`SELECT id FROM projects WHERE status = 'active' ORDER BY updated_at DESC, created_at DESC LIMIT 1`,
	).Scan(&id)
	if err == sql.ErrNoRows {
		s.ActiveProjectID = ""
		return nil
	}
	if err != nil {
		return fmt.Errorf("carregar projeto ativo: %w", err)
	}
	s.ActiveProjectID = id
	return nil
}

func (s *Store) CreateProject(project Project) (*Project, error) {
	project.ID = strings.TrimSpace(project.ID)
	if project.ID == "" {
		project.ID = uuid.NewString()
	}
	project.Name = strings.TrimSpace(project.Name)
	if project.Name == "" {
		return nil, fmt.Errorf("projeto sem nome")
	}
	project.Slug = slugify(defaultString(project.Slug, project.Name))
	if project.Slug == "" {
		return nil, fmt.Errorf("projeto sem slug")
	}
	if project.Status == "" {
		project.Status = "active"
	}
	if project.WorkspaceRoot == "" {
		project.WorkspaceRoot = filepath.Join(s.workspaceRoot, "projects", project.Slug)
	}
	if err := os.MkdirAll(project.WorkspaceRoot, 0755); err != nil {
		return nil, fmt.Errorf("criar diretorio do projeto: %w", err)
	}
	for _, dir := range []string{"src", "docs", "raw", "vault", "output"} {
		if err := os.MkdirAll(filepath.Join(project.WorkspaceRoot, dir), 0755); err != nil {
			return nil, fmt.Errorf("criar diretorio %s do projeto: %w", dir, err)
		}
	}
	timestamp := now()
	project.CreatedAt = timestamp
	project.UpdatedAt = timestamp

	_, err := s.db.Exec(
		`INSERT INTO projects(
			id, name, slug, domain, description, target_audience, main_objective,
			stack_backend, stack_frontend, stack_database, stack_infra, stack_notes,
			status, workspace_root, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		project.ID, project.Name, project.Slug, project.Domain, project.Description,
		project.TargetAudience, project.MainObjective, project.StackBackend,
		project.StackFrontend, project.StackDatabase, project.StackInfra, project.StackNotes,
		project.Status, project.WorkspaceRoot, project.CreatedAt, project.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("criar projeto: %w", err)
	}
	if project.Status == "active" {
		s.ActiveProjectID = project.ID
	}
	return &project, nil
}

func (s *Store) GetActiveProject() (*Project, error) {
	row := s.db.QueryRow(
		`SELECT id, name, slug, domain, description, target_audience, main_objective,
		        stack_backend, stack_frontend, stack_database, stack_infra, stack_notes,
		        status, workspace_root, created_at, updated_at
		   FROM projects
		  WHERE status = 'active'
		  ORDER BY updated_at DESC, created_at DESC
		  LIMIT 1`,
	)
	return scanProject(row)
}

func (s *Store) GetProject(id string) (*Project, error) {
	row := s.db.QueryRow(
		`SELECT id, name, slug, domain, description, target_audience, main_objective,
		        stack_backend, stack_frontend, stack_database, stack_infra, stack_notes,
		        status, workspace_root, created_at, updated_at
		   FROM projects
		  WHERE id = ?`,
		strings.TrimSpace(id),
	)
	return scanProject(row)
}

func (s *Store) ListProjects() ([]Project, error) {
	rows, err := s.db.Query(
		`SELECT id, name, slug, domain, description, target_audience, main_objective,
		        stack_backend, stack_frontend, stack_database, stack_infra, stack_notes,
		        status, workspace_root, created_at, updated_at
		   FROM projects
		  ORDER BY status, updated_at DESC, created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listar projetos: %w", err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		project, err := scanProject(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, *project)
	}
	return projects, rows.Err()
}

func (s *Store) UpdateProject(id string, project Project) (*Project, error) {
	current, err := s.GetProject(id)
	if err != nil {
		return nil, err
	}
	project.ID = current.ID
	project.Name = defaultString(project.Name, current.Name)
	project.Slug = slugify(defaultString(project.Slug, current.Slug))
	project.Domain = defaultString(project.Domain, current.Domain)
	project.Description = defaultString(project.Description, current.Description)
	project.TargetAudience = defaultString(project.TargetAudience, current.TargetAudience)
	project.MainObjective = defaultString(project.MainObjective, current.MainObjective)
	project.StackBackend = defaultString(project.StackBackend, current.StackBackend)
	project.StackFrontend = defaultString(project.StackFrontend, current.StackFrontend)
	project.StackDatabase = defaultString(project.StackDatabase, current.StackDatabase)
	project.StackInfra = defaultString(project.StackInfra, current.StackInfra)
	project.StackNotes = defaultString(project.StackNotes, current.StackNotes)
	project.Status = defaultString(project.Status, current.Status)
	project.WorkspaceRoot = defaultString(project.WorkspaceRoot, current.WorkspaceRoot)

	_, err = s.db.Exec(
		`UPDATE projects
		    SET name = ?, slug = ?, domain = ?, description = ?, target_audience = ?,
		        main_objective = ?, stack_backend = ?, stack_frontend = ?, stack_database = ?,
		        stack_infra = ?, stack_notes = ?, status = ?, workspace_root = ?, updated_at = ?
		  WHERE id = ?`,
		project.Name, project.Slug, project.Domain, project.Description, project.TargetAudience,
		project.MainObjective, project.StackBackend, project.StackFrontend, project.StackDatabase,
		project.StackInfra, project.StackNotes, project.Status, project.WorkspaceRoot, now(), project.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("atualizar projeto: %w", err)
	}
	if project.Status == "active" {
		s.ActiveProjectID = project.ID
	} else if s.ActiveProjectID == project.ID {
		s.ActiveProjectID = ""
		_ = s.loadActiveProjectID()
	}
	return s.GetProject(project.ID)
}

func (s *Store) ActivateProject(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id de projeto vazio")
	}
	res, err := s.db.Exec("UPDATE projects SET status = 'active', updated_at = ? WHERE id = ?", now(), id)
	if err != nil {
		return fmt.Errorf("ativar projeto: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("validar projeto ativado: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	s.ActiveProjectID = id
	return nil
}

func (s *Store) ArchiveProject(id string) error {
	id = strings.TrimSpace(id)
	res, err := s.db.Exec("UPDATE projects SET status = 'archived', updated_at = ? WHERE id = ?", now(), id)
	if err != nil {
		return fmt.Errorf("arquivar projeto: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("validar projeto arquivado: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	if s.ActiveProjectID == id {
		s.ActiveProjectID = ""
		return s.loadActiveProjectID()
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
			project_id, entry_type, document_type, section, title, content, source_agent, task_ref, tags, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.projectIDValueOrNil(entry.ProjectID),
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
	return s.projectContextItem(entry.ProjectID, "context_entries", id, strings.ToUpper(entry.DocumentType), entry.EntryType, entry.TaskRef, entry.SourceAgent, entry.Title, entry.Section, entry.Content)
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
			project_id, item_type, reference, title, status, priority, content, source_agent, task_ref, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.projectIDValueOrNil(item.ProjectID),
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
	return s.projectContextItem(item.ProjectID, "planning_items", id, "PLAN", item.ItemType, item.TaskRef, item.SourceAgent, item.Title, item.Reference, item.Content)
}

func (s *Store) SaveTestRecord(record TestRecord) error {
	if strings.TrimSpace(record.Status) == "" {
		return fmt.Errorf("test record sem status")
	}

	res, err := s.db.Exec(
		`INSERT INTO test_records(
			project_id, test_name, status, command, output, source_agent, task_ref, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		s.projectIDValueOrNil(record.ProjectID),
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
	return s.projectContextItem(record.ProjectID, "test_records", id, "TESTS", "record_test", record.TaskRef, record.SourceAgent, record.TestName, "", content)
}

func (s *Store) SaveAgentOutput(agentName, taskID, action, content string) error {
	return s.SaveAgentOutputForProject("", agentName, taskID, action, content)
}

func (s *Store) SaveAgentOutputForProject(projectID, agentName, taskID, action, content string) error {
	if strings.TrimSpace(content) == "" {
		return nil
	}
	if action == "" {
		action = "note"
	}

	res, err := s.db.Exec(
		`INSERT INTO agent_outputs(project_id, agent_name, task_id, action, content, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		s.projectIDValueOrNil(projectID),
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
	return s.projectContextItem(projectID, "agent_outputs", id, "", action, taskID, agentName, action, "", content)
}

func (s *Store) SaveHandoffEvent(event HandoffEvent) error {
	_, err := s.db.Exec(
		`INSERT INTO handoff_events(
			project_id, path, sender, recipient, task_ref, intent, payload, raw_json, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.projectIDValueOrNil(event.ProjectID),
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
	return s.QueryContextForHandoffForProject("", agentName, taskRef, limit)
}

func (s *Store) QueryContextForHandoffForProject(projectID, agentName, taskRef string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 12
	}
	agentName = strings.ToUpper(strings.TrimSpace(agentName))
	taskRef = strings.TrimSpace(taskRef)

	rows, err := s.db.Query(
		`SELECT content
		   FROM read_context_items
		  WHERE (? IS NULL OR project_id = ?)
		    AND NOT (
				source_table = 'document_sections'
				AND document_type IN ('TASKS', 'STATE', 'CONTEXT', 'PLAYBOOK', 'TESTS', 'VERSIONS', 'RETROSPECTIVE')
			)
		  ORDER BY
			CASE
				WHEN ? <> '' AND task_ref = ? THEN 100
				WHEN agent_name = ? THEN 70
				WHEN document_type IN ('TESTS', 'PLAN') AND ? <> '' THEN 50
				WHEN task_ref = '' OR task_ref IS NULL THEN 10
				ELSE 0
			END DESC,
			projected_at DESC
		  LIMIT ?`,
		s.projectIDValueOrNil(projectID),
		s.projectIDValueOrNil(projectID),
		taskRef,
		taskRef,
		agentName,
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
	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM projects WHERE status = 'active'").Scan(&count); err != nil {
		return false, fmt.Errorf("consultar projetos ativos: %w", err)
	}
	return count > 0, nil
}

func (s *Store) projectContextItem(projectID string, sourceTable string, sourceID int64, documentType, entryType, taskRef, agentName, title, heading, content string) error {
	if strings.TrimSpace(content) == "" {
		return nil
	}
	_, err := s.db.Exec(
		`INSERT INTO read_context_items(
			project_id, source_table, source_id, document_type, entry_type, task_ref, agent_name,
			title, heading, content, content_hash, projected_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.projectIDValueOrNil(projectID),
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

func (s *Store) projectIDOrNil() any {
	if strings.TrimSpace(s.ActiveProjectID) == "" {
		return nil
	}
	return s.ActiveProjectID
}

func (s *Store) projectIDValueOrNil(projectID string) any {
	if strings.TrimSpace(projectID) == "" {
		return s.projectIDOrNil()
	}
	return strings.TrimSpace(projectID)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanProject(row scanner) (*Project, error) {
	var project Project
	var domain, description, targetAudience, mainObjective sql.NullString
	var stackBackend, stackFrontend, stackDatabase, stackInfra, stackNotes sql.NullString
	var workspaceRoot sql.NullString
	if err := row.Scan(
		&project.ID,
		&project.Name,
		&project.Slug,
		&domain,
		&description,
		&targetAudience,
		&mainObjective,
		&stackBackend,
		&stackFrontend,
		&stackDatabase,
		&stackInfra,
		&stackNotes,
		&project.Status,
		&workspaceRoot,
		&project.CreatedAt,
		&project.UpdatedAt,
	); err != nil {
		return nil, err
	}
	project.Domain = domain.String
	project.Description = description.String
	project.TargetAudience = targetAudience.String
	project.MainObjective = mainObjective.String
	project.StackBackend = stackBackend.String
	project.StackFrontend = stackFrontend.String
	project.StackDatabase = stackDatabase.String
	project.StackInfra = stackInfra.String
	project.StackNotes = stackNotes.String
	project.WorkspaceRoot = workspaceRoot.String
	return &project, nil
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}

func slugify(value string) string {
	var sb strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			sb.WriteRune(r)
			lastDash = false
		default:
			if !lastDash && sb.Len() > 0 {
				sb.WriteRune('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(sb.String(), "-")
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}
