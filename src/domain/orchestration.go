package domain

import (
	"errors"
	"regexp"
	"strings"
)

type CardStatus string

const (
	CardStatusTodo       CardStatus = "todo"
	CardStatusInProgress CardStatus = "in_progress"
	CardStatusDone       CardStatus = "done"
)

type CardPriority string

const (
	CardPriorityLow    CardPriority = "low"
	CardPriorityMedium CardPriority = "medium"
	CardPriorityHigh   CardPriority = "high"
)

type PlanningCard struct {
	ID       string       `json:"id"`
	Title    string       `json:"title"`
	Owner    string       `json:"owner"`
	Status   CardStatus   `json:"status"`
	Priority CardPriority `json:"priority"`
	Phase    string       `json:"phase"`
	TaskRef  string       `json:"task_ref"`
}

func NewPlanningCard(id, title, owner string, status CardStatus, priority CardPriority, phase, taskRef string) (PlanningCard, error) {
	card := PlanningCard{
		ID:       strings.TrimSpace(id),
		Title:    strings.TrimSpace(title),
		Owner:    strings.TrimSpace(owner),
		Status:   status,
		Priority: priority,
		Phase:    strings.TrimSpace(phase),
		TaskRef:  strings.TrimSpace(taskRef),
	}
	if card.ID == "" {
		return PlanningCard{}, errors.New("planning card id is required")
	}
	if card.Title == "" {
		return PlanningCard{}, errors.New("planning card title is required")
	}
	if card.Owner == "" {
		return PlanningCard{}, errors.New("planning card owner is required")
	}
	if !card.Status.Valid() {
		return PlanningCard{}, errors.New("planning card status is invalid")
	}
	if !card.Priority.Valid() {
		return PlanningCard{}, errors.New("planning card priority is invalid")
	}
	return card, nil
}

func (s CardStatus) Valid() bool {
	return s == CardStatusTodo || s == CardStatusInProgress || s == CardStatusDone
}

func (p CardPriority) Valid() bool {
	return p == CardPriorityLow || p == CardPriorityMedium || p == CardPriorityHigh
}

type ProjectRegistration struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Audience       string         `json:"audience"`
	TechnicalStack string         `json:"technical_stack"`
	Onboarding     map[string]any `json:"onboarding,omitempty"`
}

func NewProjectRegistration(id, name, description, audience, technicalStack string, onboarding map[string]any) (ProjectRegistration, error) {
	name = strings.TrimSpace(name)
	id = strings.TrimSpace(id)
	if id == "" {
		id = ProjectIDFromName(name)
	}
	project := ProjectRegistration{
		ID:             id,
		Name:           name,
		Description:    strings.TrimSpace(description),
		Audience:       strings.TrimSpace(audience),
		TechnicalStack: strings.TrimSpace(technicalStack),
		Onboarding:     onboarding,
	}
	if project.ID == "" {
		return ProjectRegistration{}, errors.New("project id is required")
	}
	if project.Name == "" {
		return ProjectRegistration{}, errors.New("project name is required")
	}
	if project.Description == "" {
		return ProjectRegistration{}, errors.New("project description is required")
	}
	if project.Audience == "" {
		return ProjectRegistration{}, errors.New("project audience is required")
	}
	if project.TechnicalStack == "" {
		return ProjectRegistration{}, errors.New("project technical stack is required")
	}
	if containsSecretLikeText(project.Name, project.Description, project.Audience, project.TechnicalStack) {
		return ProjectRegistration{}, errors.New("project summary contains secret-like value")
	}
	if containsSecretLikeField(project.Onboarding) {
		return ProjectRegistration{}, errors.New("project onboarding contains secret-like field")
	}
	return project, nil
}

func ProjectIDFromName(name string) string {
	normalized := strings.ToLower(strings.TrimSpace(name))
	normalized = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(normalized, "-")
	return strings.Trim(normalized, "-")
}

func containsSecretLikeText(values ...string) bool {
	for _, value := range values {
		normalized := strings.ToLower(value)
		if strings.Contains(normalized, "secret") || strings.Contains(normalized, "password") || strings.Contains(normalized, "token") || strings.Contains(normalized, "api_key") || strings.Contains(normalized, "apikey") {
			return true
		}
	}
	return false
}

type MCPEnvelope struct {
	Protocol string  `json:"protocol"`
	Version  string  `json:"version"`
	Message  Handoff `json:"message"`
}

func NewMCPEnvelope(handoff Handoff) MCPEnvelope {
	return MCPEnvelope{
		Protocol: "mcp-compatible",
		Version:  "2026-05-06",
		Message:  handoff,
	}
}
