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

### [2026-05-05 15:23] - [CEO]: Phase 1 intake opened
- **Context:** Phase 0 is Done and [HUMAN] sent `[USER_DONE]` to start the next cycle.
- **Decision:** Phase 1 opens in planning mode only. Implementation remains blocked until Phase 1 kickoff handoff, SECURITY threat model validation, and explicit [HUMAN] implementation authorization.
- **Architecture:** Go Hexagonal backend, TypeScript web frontend, PostgreSQL persistence, Markdown read adapters, `.agent_handoff/` write adapter, Docker Compose for dependencies.
- **Impact:** Architectural / Security / Process.
- **References:** `DOC/ROADMAP.md`, `DOC/PLAN.md#15-phase-1-plan---mvp-agent-model-memory-and-file-handoffs`, `DOC/ARCHITECTURE.md`.

### Phase 1 - Skill Mapping
- **Applicable skills:** Backend Go/Hexagonal patterns, frontend design, database design, PostgreSQL best practices, security audit, TDD workflow, webapp testing, Docker expert, documentation.
- **Recorded by:** [CTO]
- **Timestamp:** 2026-05-05 15:23
- **Decision:** Use the listed domains as required expertise for kickoff handoff. No external skill installation is required before planning; implementation agents must prefer installed/local skills when available.

### [2026-05-05 15:23] - [SECURITY]: Phase 1 preliminary threat model
- **Threats:** Path traversal in Markdown memory reads; secret leakage into logs/handoffs; unauthenticated local API access; prompt/handoff injection from memory files; provider-mode confusion.
- **Controls:** Canonical path validation; approved file allowlist; redaction and denylisted fields; deny-by-default protected routes; treating Markdown as data; explicit provider mode contracts.
- **Decision:** Threat model is sufficient for planning. SECURITY validation is mandatory before production code begins.
- **Impact:** Security / Architectural.

### [2026-05-05 15:29] - [CEO]: Phase 1 implementation authorized
- **Context:** [HUMAN] answered the direct authorization question in `QUESTIONS.md`.
- **Decision:** Phase 1 implementation may start. CEO/CTO/BA may approve planning-level work; [HUMAN] gives final opinion at phase end. `[USER_DONE]` is treated as explicit approval of the latest pending step when context is unambiguous.
- **Impact:** Process / Functional.
- **References:** `QUESTIONS.md` entry `[2026-05-05 15:29:40] Phase 1 implementation authorization`.

### [2026-05-05 15:30] - [DEVOPS]: Environment verification
- **Context:** Phase 1 implementation started with backend Go files, tests, frontend workbench, database migration, and Docker Compose baseline.
- **Decision:** Node.js is available and JavaScript syntax validation passed. Go is not available on PATH, so Go test execution is blocked until Go is installed or a valid Go runtime is made available. Docker daemon access is inconsistent from the sandbox for image inspection.
- **Impact:** Technical / DevOps.
- **References:** `DOC/TESTS.md#11-phase-1-executed-results`.
