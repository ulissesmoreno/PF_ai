package context_store

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
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

type documentEntry struct {
	EntryType string
	Heading   string
	Content   string
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
	inserted := 0
	skipped := 0

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar importação de %s: %w", seed.Path, err)
	}
	defer tx.Rollback() //nolint

	if _, err := tx.Exec(
		`INSERT INTO documents(path, title, document_type, owner, project_id, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT(path) DO UPDATE SET
			title = excluded.title,
			document_type = excluded.document_type,
			owner = excluded.owner,
			project_id = excluded.project_id`,
		seed.Path,
		title,
		seed.DocumentType,
		seed.Owner,
		s.projectIDOrNil(),
		timestamp,
	); err != nil {
		return fmt.Errorf("upsert documento %s: %w", seed.Path, err)
	}

	documentID, err := documentID(tx, seed.Path)
	if err != nil {
		return err
	}

	res, err := tx.Exec(
		"INSERT INTO document_imports(document_id, project_id, content_hash, imported_at) VALUES (?, ?, ?, ?)",
		documentID,
		s.projectIDOrNil(),
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

	if isEntryDocument(seed.DocumentType) {
		return s.importDocumentEntries(tx, documentID, importID, seed, title, content, timestamp)
	}

	sections := splitMarkdownSections(content)
	for i, section := range sections {
		if s.sectionHashExists(tx, documentID, sha(section.Content)) {
			skipped++
			continue
		}
		res, err := tx.Exec(
			`INSERT INTO document_sections(
				document_id, import_id, project_id, heading, level, ordinal, content, content_hash, created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			documentID,
			importID,
			s.projectIDOrNil(),
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
		inserted++

		sectionID, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("obter id da seção: %w", err)
		}
		if strings.TrimSpace(section.Content) != "" {
			if _, err := tx.Exec(
				`INSERT INTO knowledge_chunks(
					source_type, source_id, project_id, document_type, heading, content, content_hash, created_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
				"document_section",
				sectionID,
				s.projectIDOrNil(),
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
					project_id, source_table, source_id, document_type, entry_type, title, heading, content, content_hash, projected_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				s.projectIDOrNil(),
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

	if err := s.importPlanningCQRS(tx, seed.DocumentType, content, timestamp); err != nil {
		return err
	}

	fmt.Printf("import %s: inserted=%d skipped=%d\n", seed.Path, inserted, skipped)
	return tx.Commit()
}

func (s *Store) importDocumentEntries(tx *sql.Tx, documentID, importID int64, seed DocumentSeed, title, content, timestamp string) error {
	entries := splitDocumentEntries(seed.DocumentType, content)
	inserted := 0
	skipped := 0
	for i, entry := range entries {
		if strings.TrimSpace(entry.Content) == "" {
			continue
		}
		if s.entryHashExists(tx, documentID, sha(entry.Content)) {
			skipped++
			continue
		}
		res, err := tx.Exec(
			`INSERT INTO document_entries(
				document_id, import_id, project_id, entry_type, heading, ordinal, content, content_hash, created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			documentID,
			importID,
			s.projectIDOrNil(),
			entry.EntryType,
			entry.Heading,
			i,
			entry.Content,
			sha(entry.Content),
			timestamp,
		)
		if err != nil {
			return fmt.Errorf("inserir entrada %s[%d]: %w", seed.Path, i, err)
		}
		inserted++
		entryID, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("obter id da entrada: %w", err)
		}
		if _, err := tx.Exec(
			`INSERT INTO read_context_items(
				project_id, source_table, source_id, document_type, entry_type, title, heading, content, content_hash, projected_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			s.projectIDOrNil(),
			"document_entries",
			entryID,
			seed.DocumentType,
			entry.EntryType,
			title,
			entry.Heading,
			entry.Content,
			sha(entry.Content),
			timestamp,
		); err != nil {
			return fmt.Errorf("projetar entrada %s[%d]: %w", seed.Path, i, err)
		}
	}

	if err := s.importPlanningCQRS(tx, seed.DocumentType, content, timestamp); err != nil {
		return err
	}

	fmt.Printf("import %s: inserted=%d skipped=%d\n", seed.Path, inserted, skipped)
	return tx.Commit()
}

func documentID(tx *sql.Tx, path string) (int64, error) {
	var id int64
	if err := tx.QueryRow("SELECT id FROM documents WHERE path = ?", path).Scan(&id); err != nil {
		return 0, fmt.Errorf("buscar documento %s: %w", path, err)
	}
	return id, nil
}

func (s *Store) sectionHashExists(tx *sql.Tx, documentID int64, contentHash string) bool {
	var exists int
	err := tx.QueryRow(
		`SELECT COUNT(1)
		   FROM document_sections
		  WHERE document_id = ?
		    AND content_hash = ?
		    AND (? IS NULL OR project_id = ?)`,
		documentID, contentHash, s.projectIDOrNil(), s.projectIDOrNil(),
	).Scan(&exists)
	return err == nil && exists > 0
}

func (s *Store) entryHashExists(tx *sql.Tx, documentID int64, contentHash string) bool {
	var exists int
	err := tx.QueryRow(
		`SELECT COUNT(1)
		   FROM document_entries
		  WHERE document_id = ?
		    AND content_hash = ?
		    AND (? IS NULL OR project_id = ?)`,
		documentID, contentHash, s.projectIDOrNil(), s.projectIDOrNil(),
	).Scan(&exists)
	return err == nil && exists > 0
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

func splitDocumentEntries(documentType, content string) []documentEntry {
	if documentType == "PLAYBOOK" {
		return splitPlaybookEntries(content)
	}
	return splitMarkdownEntryBlocks(documentType, content)
}

func splitPlaybookEntries(content string) []documentEntry {
	var entries []documentEntry
	currentHeading := ""
	for _, line := range strings.Split(content, "\n") {
		if matches := headingRE.FindStringSubmatch(line); len(matches) == 3 {
			currentHeading = strings.TrimSpace(matches[2])
			continue
		}

		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") || strings.Contains(trimmed, ":---") || strings.HasPrefix(trimmed, "| Timestamp ") {
			continue
		}
		cols := strings.Split(trimmed, "|")
		if len(cols) < 4 {
			continue
		}
		first := strings.TrimSpace(cols[1])
		if !strings.HasPrefix(first, "[") {
			continue
		}
		entries = append(entries, documentEntry{
			EntryType: "playbook_entry",
			Heading:   currentHeading,
			Content:   trimmed,
		})
	}
	return entries
}

func splitMarkdownEntryBlocks(documentType, content string) []documentEntry {
	var entries []documentEntry
	var currentSection string
	var current *documentEntry

	flush := func() {
		if current == nil {
			return
		}
		current.Content = strings.TrimSpace(current.Content)
		if current.Content != "" {
			entries = append(entries, *current)
		}
		current = nil
	}

	for _, line := range strings.Split(content, "\n") {
		if matches := headingRE.FindStringSubmatch(line); len(matches) == 3 {
			level := len(matches[1])
			heading := strings.TrimSpace(matches[2])
			if level <= 2 {
				flush()
				currentSection = heading
				continue
			}
			flush()
			current = &documentEntry{
				EntryType: strings.ToLower(documentType) + "_entry",
				Heading:   heading,
				Content:   line,
			}
			continue
		}

		if startsBulletEntry(line) {
			flush()
			current = &documentEntry{
				EntryType: strings.ToLower(documentType) + "_entry",
				Heading:   currentSection,
				Content:   line,
			}
			continue
		}

		if current != nil {
			current.Content += "\n" + line
		}
	}
	flush()

	return entries
}

func (s *Store) importPlanningCQRS(tx *sql.Tx, documentType, content, timestamp string) error {
	switch documentType {
	case "TASKS":
		for _, item := range parseTaskPlanningItems(content) {
			if err := s.insertPlanningItemTx(tx, item, timestamp); err != nil {
				return err
			}
		}
	case "ROADMAP":
		for _, item := range parseRoadmapPlanningItems(content) {
			if err := s.insertPlanningItemTx(tx, item, timestamp); err != nil {
				return err
			}
		}
	case "PLAN":
		for _, entry := range parsePlanContextEntries(content) {
			if err := s.insertContextEntryTx(tx, entry, timestamp); err != nil {
				return err
			}
		}
	case "STATE":
		for _, entry := range parseStateContextEntries(content) {
			if err := s.insertContextEntryTx(tx, entry, timestamp); err != nil {
				return err
			}
		}
	case "CONTEXT":
		for _, entry := range parseContextDecisionEntries(content) {
			if err := s.insertContextEntryTx(tx, entry, timestamp); err != nil {
				return err
			}
		}
	case "RETROSPECTIVE":
		for _, entry := range parseRetrospectiveContextEntries(content) {
			if err := s.insertContextEntryTx(tx, entry, timestamp); err != nil {
				return err
			}
		}
	case "VERSIONS":
		for _, item := range parseVersionPlanningItems(content) {
			if err := s.insertPlanningItemTx(tx, item, timestamp); err != nil {
				return err
			}
		}
	case "TESTS":
		for _, record := range parseTestRecords(content) {
			if err := s.insertTestRecordTx(tx, record, timestamp); err != nil {
				return err
			}
		}
	case "PLAYBOOK":
		for _, entry := range parsePlaybookContextEntries(content) {
			if err := s.insertContextEntryTx(tx, entry, timestamp); err != nil {
				return err
			}
		}
	}
	return nil
}

func parseTaskPlanningItems(content string) []PlanningItem {
	blocks := splitHeadingBlocks(content, 3)
	items := make([]PlanningItem, 0, len(blocks))
	for _, block := range blocks {
		if !strings.Contains(block.heading, "Task [") && !strings.Contains(block.heading, "Task ") {
			continue
		}
		item := PlanningItem{
			ItemType:  "task",
			Reference: block.heading,
			Title:     block.heading,
			Status:    fieldValue(block.content, "Status"),
			Priority:  fieldValue(block.content, "Priority"),
			Content:   strings.TrimSpace(block.content),
			TaskRef:   taskRefFromHeading(block.heading),
		}
		if assigned := fieldValue(block.content, "Assigned to"); assigned != "" {
			item.SourceAgent = assigned
		}
		items = append(items, item)
	}
	return items
}

func parseRoadmapPlanningItems(content string) []PlanningItem {
	var items []PlanningItem
	currentSection := ""
	for _, line := range strings.Split(content, "\n") {
		if matches := headingRE.FindStringSubmatch(line); len(matches) == 3 {
			currentSection = strings.TrimSpace(matches[2])
			continue
		}
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "- [ ]") && !strings.HasPrefix(trimmed, "- [x]") {
			continue
		}
		status := "todo"
		if strings.HasPrefix(trimmed, "- [x]") {
			status = "done"
		}
		title := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(trimmed, "- [ ]"), "- [x]"))
		items = append(items, PlanningItem{
			ItemType:  "roadmap_item",
			Reference: currentSection,
			Title:     stripMarkdownTitle(title),
			Status:    status,
			Content:   trimmed,
			TaskRef:   currentSection,
		})
	}
	return items
}

func parsePlanContextEntries(content string) []ContextEntry {
	blocks := splitHeadingBlocks(content, 2)
	entries := make([]ContextEntry, 0, len(blocks))
	for _, block := range blocks {
		if strings.TrimSpace(block.content) == "" {
			continue
		}
		entries = append(entries, ContextEntry{
			EntryType:    "plan_section",
			DocumentType: "PLAN",
			Section:      block.heading,
			Title:        block.heading,
			Content:      strings.TrimSpace(block.content),
			SourceAgent:  "IMPORTER",
			TaskRef:      "PLAN",
		})
	}
	return entries
}

func parseStateContextEntries(content string) []ContextEntry {
	var entries []ContextEntry
	currentSection := ""
	for _, block := range splitMarkdownEntryBlocks("STATE", content) {
		if block.Heading != "" {
			currentSection = block.Heading
		}
		entryType := "delivery"
		sectionLower := strings.ToLower(currentSection)
		contentLower := strings.ToLower(block.Content)
		switch {
		case strings.Contains(sectionLower, "blocker"):
			entryType = "blocker"
		case strings.Contains(sectionLower, "metric") || strings.Contains(contentLower, "metrics:"):
			entryType = "metric"
		}
		entries = append(entries, ContextEntry{
			EntryType:    entryType,
			DocumentType: "STATE",
			Section:      currentSection,
			Title:        firstLine(block.Content),
			Content:      block.Content,
			SourceAgent:  agentFromText(block.Content, "IMPORTER"),
			TaskRef:      "STATE",
		})
	}
	return entries
}

func parseContextDecisionEntries(content string) []ContextEntry {
	var entries []ContextEntry
	for _, block := range splitHeadingBlocks(content, 3) {
		if !strings.Contains(block.heading, "—") && !strings.Contains(block.heading, "-") {
			continue
		}
		entries = append(entries, ContextEntry{
			EntryType:    "decision",
			DocumentType: "CONTEXT",
			Section:      "Decisions Made",
			Title:        block.heading,
			Content:      strings.TrimSpace(block.content),
			SourceAgent:  agentFromText(block.heading, "IMPORTER"),
			TaskRef:      "CONTEXT",
		})
	}
	return entries
}

func parseRetrospectiveContextEntries(content string) []ContextEntry {
	var entries []ContextEntry
	for _, block := range splitHeadingBlocks(content, 3) {
		if strings.TrimSpace(block.content) == "" {
			continue
		}
		entries = append(entries, ContextEntry{
			EntryType:    "retrospective",
			DocumentType: "RETROSPECTIVE",
			Section:      block.heading,
			Title:        block.heading,
			Content:      strings.TrimSpace(block.content),
			SourceAgent:  agentFromText(block.content, "CEO"),
			TaskRef:      "RETROSPECTIVE",
		})
	}
	return entries
}

func parseVersionPlanningItems(content string) []PlanningItem {
	blocks := splitHeadingBlocks(content, 3)
	items := make([]PlanningItem, 0, len(blocks))
	for _, block := range blocks {
		if strings.TrimSpace(block.content) == "" {
			continue
		}
		itemType := strings.ToLower(fieldValue(block.content, "Type"))
		if itemType == "" {
			switch {
			case strings.Contains(strings.ToUpper(block.heading), "ROLLBACK"):
				itemType = "rollback"
			case strings.Contains(strings.ToUpper(block.heading), "BUGFIX"):
				itemType = "bugfix"
			default:
				itemType = "release"
			}
		}
		items = append(items, PlanningItem{
			ItemType:    itemType,
			Reference:   block.heading,
			Title:       block.heading,
			Status:      "recorded",
			Content:     strings.TrimSpace(block.content),
			SourceAgent: agentFromText(block.heading, "IMPORTER"),
			TaskRef:     "VERSIONS",
		})
	}
	return items
}

func parseTestRecords(content string) []TestRecord {
	var records []TestRecord
	for _, block := range splitHeadingBlocks(content, 3) {
		if !strings.Contains(strings.ToLower(block.heading), "test") {
			continue
		}
		status := fieldValuePlain(block.content, "Status")
		records = append(records, TestRecord{
			TestName:    block.heading,
			Status:      normalizeTestStatus(status),
			Command:     fieldValuePlain(block.content, "Command"),
			Output:      strings.TrimSpace(block.content),
			SourceAgent: "IMPORTER",
			TaskRef:     testRefFromHeading(block.heading),
		})
	}
	return records
}

func parsePlaybookContextEntries(content string) []ContextEntry {
	var entries []ContextEntry
	for _, entry := range splitPlaybookEntries(content) {
		entries = append(entries, ContextEntry{
			EntryType:    "playbook_entry",
			DocumentType: "PLAYBOOK",
			Section:      entry.Heading,
			Title:        firstTableColumn(entry.Content, 2),
			Content:      entry.Content,
			SourceAgent:  "IMPORTER",
			TaskRef:      "PLAYBOOK",
			Tags:         []string{firstTableColumn(entry.Content, 3)},
		})
	}
	return entries
}

type headingBlock struct {
	heading string
	content string
}

func splitHeadingBlocks(content string, level int) []headingBlock {
	var blocks []headingBlock
	var current *headingBlock
	prefix := strings.Repeat("#", level) + " "
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, prefix) {
			if current != nil && strings.TrimSpace(current.content) != "" {
				blocks = append(blocks, *current)
			}
			current = &headingBlock{heading: strings.TrimSpace(strings.TrimPrefix(line, prefix))}
			continue
		}
		if current != nil {
			current.content += line + "\n"
		}
	}
	if current != nil && strings.TrimSpace(current.content) != "" {
		blocks = append(blocks, *current)
	}
	return blocks
}

func fieldValue(content, field string) string {
	re := regexp.MustCompile(`(?m)^-\s+\*\*` + regexp.QuoteMeta(field) + `:\*\*\s*(.+?)\s*$`)
	matches := re.FindStringSubmatch(content)
	if len(matches) < 2 {
		return ""
	}
	return strings.TrimSpace(matches[1])
}

func fieldValuePlain(content, field string) string {
	if value := fieldValue(content, field); value != "" {
		return value
	}
	re := regexp.MustCompile(`(?m)^-\s+` + regexp.QuoteMeta(field) + `:\s*(.+?)\s*$`)
	matches := re.FindStringSubmatch(content)
	if len(matches) < 2 {
		return ""
	}
	return strings.TrimSpace(matches[1])
}

func normalizeTestStatus(status string) string {
	status = strings.ToLower(strings.TrimSpace(status))
	switch status {
	case "passed", "pass":
		return "passed"
	case "failed", "fail":
		return "failed"
	case "pending":
		return "pending"
	default:
		return "pending"
	}
}

func testRefFromHeading(heading string) string {
	re := regexp.MustCompile(`(?i)test\s+([A-Za-z0-9._-]+|\[[^\]]+\])`)
	matches := re.FindStringSubmatch(heading)
	if len(matches) < 2 {
		return ""
	}
	return strings.Trim(matches[1], "[]")
}

func firstTableColumn(row string, index int) string {
	cols := strings.Split(row, "|")
	if index >= len(cols) {
		return ""
	}
	return strings.TrimSpace(cols[index])
}

func taskRefFromHeading(heading string) string {
	re := regexp.MustCompile(`Task\s+\[?([A-Za-z0-9._-]+)\]?`)
	matches := re.FindStringSubmatch(heading)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func stripMarkdownTitle(value string) string {
	value = strings.ReplaceAll(value, "**", "")
	if idx := strings.Index(value, ":"); idx >= 0 {
		return strings.TrimSpace(value[:idx])
	}
	return strings.TrimSpace(value)
}

func firstLine(value string) string {
	for _, line := range strings.Split(value, "\n") {
		if strings.TrimSpace(line) != "" {
			return strings.TrimSpace(line)
		}
	}
	return ""
}

func agentFromText(value, fallback string) string {
	re := regexp.MustCompile(`\[([A-Za-z_/-]+)(?::[^\]]+)?\]`)
	matches := re.FindStringSubmatch(value)
	if len(matches) < 2 {
		return fallback
	}
	return strings.ToUpper(strings.TrimSpace(matches[1]))
}

func (s *Store) insertPlanningItemTx(tx *sql.Tx, item PlanningItem, timestamp string) error {
	if strings.TrimSpace(item.Content) == "" {
		return nil
	}
	if item.ItemType == "" {
		item.ItemType = "plan"
	}
	res, err := tx.Exec(
		`INSERT INTO planning_items(
			project_id, item_type, reference, title, status, priority, content, source_agent, task_ref, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.projectIDOrNil(), item.ItemType, item.Reference, item.Title, item.Status, item.Priority,
		item.Content, item.SourceAgent, item.TaskRef, timestamp, timestamp,
	)
	if err != nil {
		return fmt.Errorf("salvar planejamento importado: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obter id de planejamento importado: %w", err)
	}
	return s.projectContextItemTx(tx, "planning_items", id, "PLAN", item.ItemType, item.TaskRef, item.SourceAgent, item.Title, item.Reference, item.Content, timestamp)
}

func (s *Store) insertContextEntryTx(tx *sql.Tx, entry ContextEntry, timestamp string) error {
	tags, err := json.Marshal(entry.Tags)
	if err != nil {
		return fmt.Errorf("serializar tags importadas: %w", err)
	}
	res, err := tx.Exec(
		`INSERT INTO context_entries(
			project_id, entry_type, document_type, section, title, content, source_agent, task_ref, tags, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.projectIDOrNil(), entry.EntryType, strings.ToUpper(entry.DocumentType), entry.Section,
		entry.Title, entry.Content, entry.SourceAgent, entry.TaskRef, string(tags), timestamp,
	)
	if err != nil {
		return fmt.Errorf("salvar contexto importado: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obter id de contexto importado: %w", err)
	}
	return s.projectContextItemTx(tx, "context_entries", id, strings.ToUpper(entry.DocumentType), entry.EntryType, entry.TaskRef, entry.SourceAgent, entry.Title, entry.Section, entry.Content, timestamp)
}

func (s *Store) insertTestRecordTx(tx *sql.Tx, record TestRecord, timestamp string) error {
	if strings.TrimSpace(record.Status) == "" {
		record.Status = "pending"
	}
	res, err := tx.Exec(
		`INSERT INTO test_records(
			project_id, test_name, status, command, output, source_agent, task_ref, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		s.projectIDOrNil(), record.TestName, record.Status, record.Command,
		record.Output, record.SourceAgent, record.TaskRef, timestamp,
	)
	if err != nil {
		return fmt.Errorf("salvar teste importado: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obter id de teste importado: %w", err)
	}
	content := strings.TrimSpace(record.Status + "\n" + record.Command + "\n" + record.Output)
	return s.projectContextItemTx(tx, "test_records", id, "TESTS", "record_test", record.TaskRef, record.SourceAgent, record.TestName, "", content, timestamp)
}

func (s *Store) projectContextItemTx(tx *sql.Tx, sourceTable string, sourceID int64, documentType, entryType, taskRef, agentName, title, heading, content, timestamp string) error {
	if strings.TrimSpace(content) == "" {
		return nil
	}
	_, err := tx.Exec(
		`INSERT INTO read_context_items(
			project_id, source_table, source_id, document_type, entry_type, task_ref, agent_name,
			title, heading, content, content_hash, projected_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.projectIDOrNil(), sourceTable, sourceID, documentType, entryType, taskRef,
		strings.ToUpper(agentName), title, heading, content, sha(content), timestamp,
	)
	if err != nil {
		return fmt.Errorf("projetar read model importado: %w", err)
	}
	return nil
}

func startsBulletEntry(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "- **[") ||
		strings.HasPrefix(trimmed, "- [") ||
		strings.HasPrefix(trimmed, "- **Test ") ||
		strings.HasPrefix(trimmed, "- **Phase ") ||
		strings.HasPrefix(trimmed, "- **Version ")
}

func isEntryDocument(documentType string) bool {
	switch documentType {
	case "TASKS", "STATE", "CONTEXT", "PLAYBOOK", "TESTS", "VERSIONS", "RETROSPECTIVE":
		return true
	default:
		return false
	}
}

func titleFromPath(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func sha(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}
