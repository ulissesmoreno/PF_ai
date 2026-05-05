# CONTEXT.md

> **Owner:** Management agents (`[CEO]`, `[CTO]`, `[BA]`).
> This file records decisions, justifications, and architectural choices. Technical execution memory lives in `STATE.md`.
> **Immutable:** append only. Never delete entries. Timestamps mandatory (YYYY-MM-DD HH:MM).

---

## 1. Objective

Record decisions and context before execution. Maintain transparency between phases.

---

## 2. Current Phase Context

- **Phase:** [Phase name and number]
- **Scope:** [What is being built this phase]
- **Components involved:** [E.g.: Backend API, PostgreSQL schema]
- **Known constraints:** [E.g.: Local hardware for ML]

---

## 3. Decisions Made

### [YYYY-MM-DD HH:MM] — [CEO/CTO/BA]: Decision title
- **Decision:** [What was decided]
- **Why:** [Justification]
- **Impact:** [Architectural / Business / Security impact]
- **References:** [Related files or handoffs]

---

## 4. Hypotheses and Assumptions

- **[YYYY-MM-DD HH:MM] — [Agent]:** [Hypothesis and basis]

---

## 5. Skill Mapping (per phase — CTO responsibility)

### Phase X — [Phase Name]
- **Applicable skills:** [List from GSD-RULES §14]
- **Recorded by:** [CTO]
- **Timestamp:** [YYYY-MM-DD HH:MM]

---

## 6. Integration with Other Files

- **STATE.md:** Technical execution memory (owned by technical agents).
- **QUESTIONS.md:** Open questions and decisions pending [HUMAN] input.
- **ROADMAP.md:** Phase progress tracking.
- **RETROSPECTIVE.md:** Phase learnings (owned by CEO).

---

## 7. Version History (Immutable)

- **[YYYY-MM-DD HH:MM] — CEO:** Initial template. Ownership clarified: management agents only.
- [Add new entries here without deleting.]
---

## 8. Activation Log

### [2026-05-05 14:51] - [CEO]: GSD activation stop-gate
- **Context:** GSD activation requested by [HUMAN]. CEO mandatory reading sequence executed using available files: `DOC/GSD-RULES.md`, `PLAYBOOK.md`, `NEW-INSTRUCTIONS.md`, `DOC/PLAN.md`, and `DOC/ONBOARDING.md`.
- **Finding:** Root `GSD-RULES.md` and root `PROJECT.md` are absent; canonical files exist under `DOC/`. `DOC/PROJECT.md`, `NEW-INSTRUCTIONS.md`, and `DOC/PLAN.md` still contain placeholders/templates.
- **Decision:** Phase execution is blocked. `DOC/ONBOARDING.md` is active and the 5-block onboarding conversation must be completed before any implementation or roadmap phase.
- **Impact:** Business / Architectural / Security.
- **References:** `QUESTIONS.md` entry `[2026-05-05 14:51:40] Project onboarding required`.

### [2026-05-05 15:00] - [CEO]: Onboarding answer reviewed
- **Context:** [HUMAN] answered the onboarding request in `QUESTIONS.md`.
- **Finding:** Product identity, visual inspiration, team model, and high-level objective are present. Stack and MVP remain decision points because [HUMAN] explicitly asked whether Go or Ruby on Rails can be used and requested MVP suggestions.
- **Decision:** Proposed Go as the recommended backend direction and a web-first MVP focused on agent/model/memory management, GSD file visibility, and structured handoffs. Auto-fill remains blocked until [HUMAN] confirms the proposal.
- **Impact:** Business / Architectural.
- **References:** `QUESTIONS.md` entry `[2026-05-05 15:00:16] Confirm onboarding auto-fill proposal`.

### [2026-05-05 15:06] - [CEO]: Onboarding completed
- **Context:** [HUMAN] confirmed Go, the suggested MVP, API/local agent execution, Docker local model runtime, file handoffs for MVP, MCP handoffs for the final stage, and the rule that `CONTEXT.md` is updated only after onboarding.
- **Decision:** PF_ai is defined as a Go-backed, TypeScript web, PostgreSQL, Docker-enabled agent/model/memory manager with file-based MVP handoffs and MCP planned for later.
- **Impact:** Business / Architectural / Security.
- **References:** `DOC/PROJECT.md`, `README.md`, `DOC/DESIGN.md`, `DOC/ROADMAP.md`, `DOC/PLAN.md`, `DOC/ENV_SETUP.md`, `DOC/ARCHITECTURE.md`, `PLAYBOOK.md`, `QUESTIONS.md`.

### [2026-05-05 15:14] - [CEO]: New instruction intake
- **Context:** [HUMAN] sent `[USER_DONE]`; `NEW-INSTRUCTIONS.md` now includes wiki timing preference and a Go/Hexagonal Architecture question.
- **Decision:** Wiki updates happen only at phase end unless explicitly requested. Go remains compatible with Hexagonal Architecture; PF_ai keeps domain/application/infrastructure boundaries.
- **Impact:** Architectural / Process.
- **References:** `NEW-INSTRUCTIONS.md`, `QUESTIONS.md` entry `[2026-05-05 15:14:57] Go and Hexagonal Architecture compatibility`, `PLAYBOOK.md`.
