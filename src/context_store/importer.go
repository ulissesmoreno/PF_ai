package context_store

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type DocumentSeed struct {
	Path         string
	DocumentType string
	Owner        string
}

var OperationalDocuments = []DocumentSeed{
	{Path: "DOC/ARCHITECTURE.md", DocumentType: "ARCHITECTURE", Owner: "CTO"},
	{Path: "DOC/CONTEXT.md", DocumentType: "CONTEXT", Owner: "Management agents"},
	{Path: "DOC/DESIGN.md", DocumentType: "DESIGN", Owner: "UX_RESEARCHER/DEV_FRONTEND"},
	{Path: "DOC/ENV_SETUP.md", DocumentType: "ENV_SETUP", Owner: "DEVOPS"},
	{Path: "DOC/GSD-RULES.md", DocumentType: "GSD_RULES", Owner: "CEO/CTO"},
	{Path: "DOC/PLAN.md", DocumentType: "PLAN", Owner: "BA/CTO"},
	{Path: "DOC/PROJECT.md", DocumentType: "PROJECT", Owner: "CEO"},
	{Path: "DOC/RETROSPECTIVE.md", DocumentType: "RETROSPECTIVE", Owner: "CEO"},
	{Path: "DOC/ROADMAP.md", DocumentType: "ROADMAP", Owner: "CEO/PM"},
	{Path: "DOC/STATE.md", DocumentType: "STATE", Owner: "Technical agents"},
	{Path: "DOC/TASKS.md", DocumentType: "TASKS", Owner: "Technical agents"},
	{Path: "DOC/TESTS.md", DocumentType: "TESTS", Owner: "QA/SECURITY/Technical agents"},
	{Path: "DOC/VERSIONS.md", DocumentType: "VERSIONS", Owner: "Phase-closing agent"},
	{Path: "DOC/WIKI.md", DocumentType: "WIKI_PROTOCOL", Owner: "DOCUMENTATION"},
	{Path: "PLAYBOOK.md", DocumentType: "PLAYBOOK", Owner: "CEO"},
}

type markdownSection struct {
	Heading string
	Level   int
	Content string
}

var headingRE = regexp.MustCompile(`^(#{1,6})\s+(.+?)\s*$`)

func (s *Store) ImportOperationalDocuments() error {
	for _, doc := range OperationalDocuments {
		if err := s.ImportDocument(doc); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
	}
	return nil
}

func (s *Store) ImportDocument(seed DocumentSeed) error {
	fullPath := filepath.Join(s.workspaceRoot, seed.Path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	content := string(data)
	title := titleFromPath(seed.Path)
	hash := sha(content)
	timestamp := now()

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar importação de %s: %w", seed.Path, err)
	}
	defer tx.Rollback() //nolint

	if _, err := tx.Exec(
		`INSERT INTO documents(path, title, document_type, owner, created_at)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(path) DO UPDATE SET
			title = excluded.title,
			document_type = excluded.document_type,
			owner = excluded.owner`,
		seed.Path,
		title,
		seed.DocumentType,
		seed.Owner,
		timestamp,
	); err != nil {
		return fmt.Errorf("upsert documento %s: %w", seed.Path, err)
	}

	documentID, err := documentID(tx, seed.Path)
	if err != nil {
		return err
	}

	res, err := tx.Exec(
		"INSERT INTO document_imports(document_id, content_hash, imported_at) VALUES (?, ?, ?)",
		documentID,
		hash,
		timestamp,
	)
	if err != nil {
		return fmt.Errorf("registrar importação %s: %w", seed.Path, err)
	}
	importID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obter id da importação: %w", err)
	}

	sections := splitMarkdownSections(content)
	for i, section := range sections {
		res, err := tx.Exec(
			`INSERT INTO document_sections(
				document_id, import_id, heading, level, ordinal, content, content_hash, created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			documentID,
			importID,
			section.Heading,
			section.Level,
			i,
			section.Content,
			sha(section.Content),
			timestamp,
		)
		if err != nil {
			return fmt.Errorf("inserir seção %s[%d]: %w", seed.Path, i, err)
		}

		sectionID, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("obter id da seção: %w", err)
		}
		if strings.TrimSpace(section.Content) != "" {
			if _, err := tx.Exec(
				`INSERT INTO knowledge_chunks(
					source_type, source_id, document_type, heading, content, content_hash, created_at
				) VALUES (?, ?, ?, ?, ?, ?, ?)`,
				"document_section",
				sectionID,
				seed.DocumentType,
				section.Heading,
				section.Content,
				sha(section.Content),
				timestamp,
			); err != nil {
				return fmt.Errorf("inserir chunk %s[%d]: %w", seed.Path, i, err)
			}
			if _, err := tx.Exec(
				`INSERT INTO read_context_items(
					source_table, source_id, document_type, entry_type, title, heading, content, content_hash, projected_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				"document_sections",
				sectionID,
				seed.DocumentType,
				"document_section",
				title,
				section.Heading,
				section.Content,
				sha(section.Content),
				timestamp,
			); err != nil {
				return fmt.Errorf("projetar seção %s[%d]: %w", seed.Path, i, err)
			}
		}
	}

	return tx.Commit()
}

func documentID(tx *sql.Tx, path string) (int64, error) {
	var id int64
	if err := tx.QueryRow("SELECT id FROM documents WHERE path = ?", path).Scan(&id); err != nil {
		return 0, fmt.Errorf("buscar documento %s: %w", path, err)
	}
	return id, nil
}

func splitMarkdownSections(content string) []markdownSection {
	lines := strings.Split(content, "\n")
	sections := []markdownSection{
		{Heading: "_root", Level: 0},
	}

	for _, line := range lines {
		if matches := headingRE.FindStringSubmatch(line); len(matches) == 3 {
			sections = append(sections, markdownSection{
				Heading: strings.TrimSpace(matches[2]),
				Level:   len(matches[1]),
			})
			continue
		}

		idx := len(sections) - 1
		if sections[idx].Content == "" {
			sections[idx].Content = line
		} else {
			sections[idx].Content += "\n" + line
		}
	}

	for i := range sections {
		sections[i].Content = strings.TrimSpace(sections[i].Content)
	}

	return sections
}

func titleFromPath(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func sha(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}
