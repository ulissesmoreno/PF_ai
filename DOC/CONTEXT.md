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

### [2026-05-05 20:06] - [CEO] / [CTO] / [BA]: Phase 1 closure-ready MVP scope
- **Context:** [HUMAN] tested the dashboard and added notes in `NEW-INSTRUCTIONS.md`. Phase 1 implementation was adjusted until provider save, agent save, optional provider, automatic agent ID, role list with `Outro`, local-save fallback, and default CEO/CTO/BA agents were present.
- **Decision:** Phase 1 MVP is closure-ready for [HUMAN] review. The MVP includes a Go Hexagonal API, PostgreSQL migration and repository contracts, Markdown memory read adapter, `.agent_handoff/` writer, and an operational web workbench with API-backed paths plus localStorage fallback for human browser testing when the backend process runner is unavailable.
- **Why:** The proposed MVP is to prove agent/model/provider/memory/handoff management, not Phase 2 local model execution or full persistent service orchestration.
- **Impact:** Business / Architectural / Security.
- **References:** `NEW-INSTRUCTIONS.md`, `DOC/PLAN.md#15-phase-1-plan---mvp-agent-model-memory-and-file-handoffs`, `DOC/TESTS.md#11-phase-1-executed-results`, `pf-ai-web/src/app.js`, `src/infrastructure/httpapi/server.go`.

### [2026-05-05 20:06] - [CTO]: Provider ownership and agent-provider decoupling
- **Decision:** Local/online/hybrid runtime configuration belongs to provider registration. Agent registration selects a provider when needed, or selects no provider for embedded-model environments such as Codex and Claude Code.
- **Why:** This matches the UI model identified by [HUMAN] during testing and prevents agent records from duplicating provider credentials/runtime data.
- **Impact:** Architectural / Security / UX.
- **References:** `NEW-INSTRUCTIONS.md`, `src/domain/agent.go`, `src/domain/provider.go`, `pf-ai-web/index.html`.

### [2026-05-05 20:06] - [SECURITY]: Phase 1 closure audit
- **Decision:** SECURITY approves Phase 1 for closure review. No real secrets were found in scanned implementation/documentation hits. Handoff creation denies secret-looking payloads and API tests verify generated handoffs do not contain the local auth sentinel.
- **Known limitation:** Backend OS-level background runtime validation remains an environment/process-runner issue. MVP API behavior is covered through HTTP handler tests; manual API-backed browser validation still requires starting the backend interactively.
- **Impact:** Security / DevOps.
- **References:** `DOC/TESTS.md#12-phase-1-parallel-review-snapshot`, `TEST-PHASE1-CLOSURE-001`.

### [2026-05-05 20:21] - [CEO] / [CTO] / [BA]: Phase 2 runtime routing baseline
- **Context:** Phase 2 opened after [HUMAN] approved Phase 1 and sent `[USER_DONE]` for the next cycle.
- **Decision:** Phase 2 delivers a runtime/routing baseline: provider health endpoint, deterministic route decision, hybrid fallback reason, loopback-only local runtime endpoint validation, dashboard provider status action, and repeatable foreground backend startup script.
- **Why:** This satisfies the roadmap intent for local runtime and hybrid provider routing without requiring a large local model image selection during this cycle.
- **Impact:** Architectural / Security / DevOps / UX.
- **References:** `DOC/PLAN.md#16-phase-2-plan---local-runtime-and-hybrid-provider-routing`, `.agent_handoff/2026-05-05_2015_phase2_kickoff.json`, `src/domain/provider.go`, `src/infrastructure/httpapi/server.go`, `pf-ai-web/src/app.js`.

### [2026-05-05 20:21] - [SECURITY]: Phase 2 closure audit
- **Decision:** SECURITY approves Phase 2 baseline for Stage Closure Gate review.
- **Controls validated:** Protected health route; no auth sentinel leakage in response; local/hybrid endpoints restricted to loopback hosts; local runtime field is treated as configuration, not executed shell text.
- **Known limitation:** Real local model image selection remains configurable and should be chosen explicitly when [HUMAN] wants concrete model execution beyond endpoint health.
- **Impact:** Security / Runtime.
- **References:** `DOC/TESTS.md#14-phase-2-executed-results`.

### [2026-05-06 07:33] - [CEO] / [CTO] / [BA]: Phase 3 user-visible orchestration
- **Context:** [HUMAN] questioned whether Phase 2 was only technical. Phase 3 therefore prioritizes visible value over transport-only work.
- **Decision:** Phase 3 delivers project registration, onboarding payload validation, planning cards as the dashboard, and MCP-compatible handoff envelope export.
- **Why:** This converts PF_ai from a technical runtime manager into a usable orchestration workbench where planning, ownership, and handoff readiness are visible.
- **Impact:** Business / UX / Architectural.
- **References:** `DOC/PLAN.md#17-phase-3-plan---mcp-handoffs-and-advanced-orchestration`, `src/domain/orchestration.go`, `pf-ai-web/index.html`.

### [2026-05-06 07:33] - [SECURITY]: Phase 3 closure audit
- **Decision:** SECURITY approves Phase 3 for Stage Closure Gate review.
- **Controls validated:** Secret-like onboarding payload fields rejected; MCP envelope export reuses handoff validation; frontend escapes card/project text; protected APIs remain bearer-token gated.
- **Impact:** Security / Product.
- **References:** `DOC/TESTS.md#15-phase-3-executed-results`.

### [2026-05-06 08:10] - [CEO] / [CTO] / [BA]: R1 project-centric workspace root
- **Context:** [HUMAN] requested a full project refactor beginning with project registration and approved the R1 plan via `[USER_DONE]`.
- **Decision:** Project becomes the root workspace entity for the UI. The sidebar lists projects only, while feature navigation lives inside the selected project workspace. R1 stores project summaries through the API and scopes local frontend records by active project.
- **Why:** This makes PF_ai usable as a multi-project orchestration tool and prevents the dashboard from mixing unrelated project data.
- **Impact:** Business / UX / Architectural.
- **References:** `DOC/PLAN.md#18-refactor-plan---r1-project-registration`, `pf-ai-web/index.html`, `pf-ai-web/src/app.js`, `src/domain/orchestration.go`.

### [2026-05-06 08:10] - [SECURITY]: R1 project field validation
- **Decision:** Project registration rejects secret-like keys in onboarding payloads and secret-like text in project summary fields. Browser-created project question signals remain local UI state and do not write into `QUESTIONS.md`.
- **Why:** Project descriptions and stack notes are likely to receive human-entered text; refusing obvious credential markers reduces accidental secret capture during onboarding.
- **Impact:** Security / Documentation.
- **References:** `DOC/TESTS.md#17-refactor-r1-executed-results`, `src/domain/orchestration.go`, `tests/domain/orchestration_test.go`.
