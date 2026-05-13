package context_store

import (
	"database/sql"
	"fmt"
	"strings"
)

type Card struct {
	ID         string `json:"id"`
	ProjectID  string `json:"project_id,omitempty"`
	Title      string `json:"title"`
	TaskRef    string `json:"task_ref"`
	Sender     string `json:"sender"`
	Recipient  string `json:"recipient"`
	Intent     string `json:"intent"`
	Status     string `json:"status"`
	Priority   string `json:"priority"`
	RetryCount int    `json:"retry_count"`
	MaxRetries int    `json:"max_retries"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type CardComment struct {
	ID          int64  `json:"id"`
	CardID      string `json:"card_id"`
	Author      string `json:"author"`
	Content     string `json:"content"`
	CommentType string `json:"comment_type"`
	CreatedAt   string `json:"created_at"`
}

type CardHandoff struct {
	CardID    string
	Title     string
	TaskRef   string
	Sender    string
	Recipient string
	Intent    string
	Priority  string
	Payload   string
	RawJSON   string
}

type CardFilter struct {
	ProjectID string
	Status    string
	TaskRef   string
}

type CardThread struct {
	Card     Card          `json:"card"`
	Comments []CardComment `json:"comments"`
}

func (s *Store) SaveCardHandoff(event CardHandoff) (*Card, error) {
	cardID := strings.TrimSpace(event.CardID)
	if cardID == "" {
		return nil, fmt.Errorf("card_id vazio")
	}
	status := "open"
	if strings.EqualFold(event.Recipient, "[HUMAN]") || strings.EqualFold(strings.Trim(event.Recipient, "[]"), "HUMAN") {
		status = "blocked"
	}
	title := strings.TrimSpace(event.Title)
	if title == "" {
		title = strings.TrimSpace(event.Intent)
	}
	if title == "" {
		title = cardID
	}
	projectID := s.projectIDOrNil()
	timestamp := now()

	_, err := s.db.Exec(
		`INSERT INTO cards(
			id, project_id, title, task_ref, sender, recipient, intent, status, priority,
			retry_count, max_retries, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0, 3, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			recipient = excluded.recipient,
			status = excluded.status,
			updated_at = excluded.updated_at`,
		cardID, projectID, title, event.TaskRef, event.Sender, event.Recipient,
		event.Intent, status, event.Priority, timestamp, timestamp,
	)
	if err != nil {
		return nil, fmt.Errorf("salvar card: %w", err)
	}
	if err := s.AddCardComment(cardID, event.Sender, "handoff", defaultString(event.RawJSON, event.Payload)); err != nil {
		return nil, err
	}
	return s.GetCard(cardID)
}

func (s *Store) AddCardComment(cardID, author, commentType, content string) error {
	if strings.TrimSpace(cardID) == "" {
		return fmt.Errorf("card_id vazio")
	}
	if strings.TrimSpace(commentType) == "" {
		commentType = "note"
	}
	if strings.TrimSpace(content) == "" {
		content = "{}"
	}
	_, err := s.db.Exec(
		`INSERT INTO card_comments(card_id, author, content, comment_type, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		cardID, author, content, commentType, now(),
	)
	if err != nil {
		return fmt.Errorf("salvar comentario do card: %w", err)
	}
	return nil
}

func (s *Store) UpdateCardStatus(cardID, status string) error {
	switch status {
	case "open", "in_progress", "blocked", "done", "failed", "canceled":
	default:
		return fmt.Errorf("status de card invalido: %s", status)
	}
	_, err := s.db.Exec("UPDATE cards SET status = ?, updated_at = ? WHERE id = ?", status, now(), cardID)
	if err != nil {
		return fmt.Errorf("atualizar status do card: %w", err)
	}
	return nil
}

func (s *Store) CancelCard(cardID string) error {
	return s.UpdateCardStatus(cardID, "canceled")
}

func (s *Store) GetCard(id string) (*Card, error) {
	row := s.db.QueryRow(
		`SELECT id, project_id, title, task_ref, sender, recipient, intent, status, priority,
		        retry_count, max_retries, created_at, updated_at
		   FROM cards WHERE id = ?`,
		id,
	)
	return scanCard(row)
}

func (s *Store) GetCardThread(id string) (*CardThread, error) {
	card, err := s.GetCard(id)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(
		`SELECT id, card_id, author, content, comment_type, created_at
		   FROM card_comments WHERE card_id = ? ORDER BY id`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("listar comentarios do card: %w", err)
	}
	defer rows.Close()
	thread := &CardThread{Card: *card}
	for rows.Next() {
		var comment CardComment
		if err := rows.Scan(&comment.ID, &comment.CardID, &comment.Author, &comment.Content, &comment.CommentType, &comment.CreatedAt); err != nil {
			return nil, err
		}
		thread.Comments = append(thread.Comments, comment)
	}
	return thread, rows.Err()
}

func (s *Store) ListCards(filter CardFilter) ([]Card, error) {
	rows, err := s.db.Query(
		`SELECT id, project_id, title, task_ref, sender, recipient, intent, status, priority,
		        retry_count, max_retries, created_at, updated_at
		   FROM cards
		  WHERE (? = '' OR project_id = ?)
		    AND (? = '' OR status = ?)
		    AND (? = '' OR task_ref = ?)
		  ORDER BY updated_at DESC`,
		filter.ProjectID, filter.ProjectID, filter.Status, filter.Status, filter.TaskRef, filter.TaskRef,
	)
	if err != nil {
		return nil, fmt.Errorf("listar cards: %w", err)
	}
	defer rows.Close()
	var cards []Card
	for rows.Next() {
		card, err := scanCard(rows)
		if err != nil {
			return nil, err
		}
		cards = append(cards, *card)
	}
	return cards, rows.Err()
}

func scanCard(row scanner) (*Card, error) {
	var card Card
	var projectID, taskRef, priority sql.NullString
	if err := row.Scan(
		&card.ID, &projectID, &card.Title, &taskRef, &card.Sender, &card.Recipient,
		&card.Intent, &card.Status, &priority, &card.RetryCount, &card.MaxRetries,
		&card.CreatedAt, &card.UpdatedAt,
	); err != nil {
		return nil, err
	}
	card.ProjectID = projectID.String
	card.TaskRef = taskRef.String
	card.Priority = priority.String
	return &card, nil
}
