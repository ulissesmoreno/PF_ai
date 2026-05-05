# PLAN - PF_ai

This file defines the active roadmap stage before code. Tests are mandatory before implementation.

## Updated on: 2026-05-05 15:06:38

## 1. Stage Objective
Complete Phase 0 onboarding and prepare Phase 1 MVP execution for PF_ai. The project must establish the product vision, stack, architecture, design direction, environment setup, roadmap, and MVP acceptance criteria before any code is written.

## 2. Roadmap Stage
- **Stage name:** Phase 0 - Onboarding and Foundation
- **Related roadmap item:** `DOC/ROADMAP.md#phase-0-onboarding-and-foundation-status-done`

## 3. Detailed Description
- **What will be implemented:** No production code in Phase 0. This stage fills project documents and prepares the management handoff for Phase 1.
- **Expected behavior:** All required onboarding files contain real PF_ai-specific content and no onboarding placeholders in required sections.
- **Impacted components:** Documentation, architecture baseline, environment setup, design tokens, roadmap, and MVP plan.

## 4. Acceptance Criteria
- Criterion 1: `DOC/PROJECT.md` defines PF_ai vision, objectives, scope, stack, team, risks, and KPIs.
- Criterion 2: `README.md` describes PF_ai in Portuguese and English without onboarding placeholders.
- Criterion 3: `DOC/ARCHITECTURE.md` documents Go, TypeScript web, PostgreSQL, Docker, API/local providers, file handoffs, and MCP future scope.
- Criterion 4: `DOC/DESIGN.md` contains concrete design tokens and layout guidance.
- Criterion 5: `DOC/ENV_SETUP.md` defines local prerequisites and environment variables without real secrets.
- Criterion 6: `DOC/ROADMAP.md` defines Phase 0 through Phase 3.
- Criterion 7: `QUESTIONS.md` records the onboarding answers and decisions.

## 5. Business Rules
- Rule 1: [HUMAN:Ulisses] is the final authority for roadmap advancement.
- Rule 2: PF_ai must support agents through API providers and local runtime paths.
- Rule 3: MVP handoffs use files; MCP handoffs are final-stage scope.
- Rule 4: Chat responses must stay compact; detailed state belongs in Markdown files.
- Rule 5: `CONTEXT.md` must be updated only after onboarding completion or later phase intake.

## 6. Recommended Tests
- Documentation test 1: Scan required onboarding files for placeholder markers.
- Documentation test 2: Verify stack consistency across `PROJECT.md`, `ARCHITECTURE.md`, `ENV_SETUP.md`, `ROADMAP.md`, and `README.md`.
- Security test 1: Confirm no real secrets were written to documentation.
- Process test 1: Confirm `QUESTIONS.md` contains [HUMAN] answers and decisions.
- Regression test 1: Confirm GSD chain of command and Stage Closure Gate rules remain documented.

## 7. Security Validations
- Validation 1: Secrets are represented only as environment variable names.
- Validation 2: API routes are planned as deny-by-default.
- Validation 3: Logs are planned as structured and sanitized.
- Validation 4: Local model runtime is isolated through Docker Compose.
- Security notes: Threat modeling is mandatory before Phase 1 implementation.

## 8. Developer Validation Steps
1. Review onboarding files listed in acceptance criteria.
2. Confirm the stack and MVP scope.
3. Confirm no required file still has onboarding placeholders.
4. Approve or correct Phase 0 closure.
5. Authorize Phase 1 kickoff only after Stage Closure Gate.
6. Record any corrections in `QUESTIONS.md`.
7. Record test results in `DOC/TESTS.md` during closure validation.

## 9. Responsible and Dependencies
- **Responsible:** [CEO] for onboarding documents; [CTO]/[BA] for Phase 1 kickoff validation.
- **Dependencies:** [HUMAN] confirmation in `QUESTIONS.md`.
- **Time Estimate:** Phase 0 documentation pass: same work cycle.
- **Extra-Code Prerequisites:** Docker availability and provider/model decisions before implementation.

## 10. Risks and Mitigations
- Risk 1: Local model runtime may exceed machine resources.
  - Mitigation: Keep provider mode configurable and validate local runtime in Phase 2.
- Risk 2: Handoff schemas may change before MCP migration.
  - Mitigation: Version file handoff schemas from the MVP.
- Risk 3: Documentation and application state may drift.
  - Mitigation: Treat Markdown memory as an explicit adapter with tests.

## 11. Integration with Other Files
- **ROADMAP.md:** Defines Phase 0 through Phase 3.
- **STATE.md:** Updated by technical agents after implementation begins.
- **TESTS.md:** Receives validation results at closure.
- **CONTEXT.md:** Updated after onboarding completion per [HUMAN] instruction.

## 12. Version History (Immutable)
- **[2026-05-05 15:06] - [CEO]:** Phase 0 onboarding plan filled from [HUMAN] answers.

## 13. Plan Validation
- **Technical Feasibility:** Approved for planning. Go, TypeScript web, PostgreSQL, Docker, and file handoffs are feasible for the MVP.
- **Project Alignment:** Aligned with PF_ai objective as agent/model/memory manager.
- **Risks Assessed:** Local model runtime and MCP migration risks documented.
- **Approval:** Approved for Phase 0 documentation completion. Phase 1 implementation still requires Stage Closure Gate.
- **ROADMAP.md Reference:** This plan is validated here; `ROADMAP.md` tracks phase progress.

## 14. Stage Closure Gate

### Completed Stage: Phase 0 - Onboarding and Foundation
### Closure Date: 2026-05-05 15:16
### Responsible: [CEO]

### Closure Checklist

| # | File | Status | Note |
| :- | :--- | :---: | :--- |
| 1 | `TESTS.md` | `[x]` | Phase 0 validation results recorded. |
| 2 | `STATE.md` | `[x]` | Phase 0 onboarding completion recorded. |
| 3 | `TASKS.md` | `[x]` | Phase 0 documentation and architecture tasks recorded as Done. |
| 4 | `CONTEXT.md` | `[x]` | Updated after onboarding completion and new instruction intake. |
| 5 | `README.md` | `[x]` | Filled during onboarding. |
| 6 | `ROADMAP.md` | `[x]` | Phase 0 approved by [HUMAN:Ulisses] and marked Done. |
| 7 | `VERSIONS.md` | `[x]` | `1.2.0-alpha.1` onboarding milestone recorded. |

### Pending Items / Blockers
- None.

### Delivery Summary
- Phase 0 onboarding documentation was filled from [HUMAN] answers and validated.
- Metrics: 0 onboarding placeholders found; 0 real secrets found; Go/Hexagonal compatibility confirmed.

### User Approval
- [x] User confirmed stage closure and authorizes roadmap advancement on 2026-05-05 15:20.

---

## 15. Phase 1 Plan - MVP Agent, Model, Memory, and File Handoffs

## Updated on: 2026-05-05 15:23:09

## 15.1 Stage Objective
Deliver the PF_ai MVP planning baseline and prepare implementation for a web-first agent/model/memory manager. Phase 1 must prove that [HUMAN:Ulisses] can operate PF_ai through a UI/API that registers agents and model providers, reads GSD memory files, and creates structured file handoffs.

## 15.2 Roadmap Stage
- **Stage name:** Phase 1 - MVP Agent, Model, Memory, and File Handoffs.
- **Related roadmap item:** `DOC/ROADMAP.md#phase-1-mvp---agent-model-memory-and-file-handoffs-status-doing`

## 15.3 Detailed Description
- **What will be implemented after kickoff approval:** Go backend, TypeScript web frontend, PostgreSQL schema, file-memory adapters, `.agent_handoff/` JSON creation, and protected local API routes.
- **Expected behavior:** [HUMAN:Ulisses] can inspect current phase/status, register agents, register model/provider configurations, view key GSD files, and create a structured handoff file.
- **Impacted components:** `pf-ai-service`, `pf-ai-web`, PostgreSQL, Docker Compose, `.agent_handoff/`, Markdown memory adapters, authentication/session handling.

## 15.4 Acceptance Criteria
- Criterion 1: Backend exposes protected API endpoints for agents, model providers, GSD memory reads, and handoff creation.
- Criterion 2: Domain packages contain no framework, database, HTTP, Docker, or filesystem dependencies.
- Criterion 3: PostgreSQL schema persists agents, providers, sessions, and audit metadata.
- Criterion 4: File adapter reads approved GSD memory files without writing outside the project root.
- Criterion 5: Handoff adapter creates schema-valid JSON files under `.agent_handoff/` without secrets.
- Criterion 6: Web UI shows phase/status, agents, providers, selected memory files, and handoff actions.
- Criterion 7: Unauthenticated requests to protected routes return 401.
- Criterion 8: Tests are written before production code and results are recorded in `DOC/TESTS.md`.

## 15.5 Business Rules
- Rule 1: MVP handoffs use files; MCP remains later-stage scope.
- Rule 2: Agents can be configured for API, local, or hybrid provider mode, but local model execution itself belongs to Phase 2 unless explicitly pulled forward.
- Rule 3: Secrets never appear in `.md`, logs, database seed data, or handoff files.
- Rule 4: `QUESTIONS.md` is used only when an agent must ask [HUMAN] a direct question.
- Rule 5: `VERSIONS.md` and `wiki/` are updated at phase closure only.

## 15.6 Recommended Tests
- Unit: Domain validation for agent definitions, provider modes, memory paths, and handoff metadata.
- Unit: Handoff schema builder rejects missing required fields and secret-looking fields.
- Integration: PostgreSQL repositories persist and retrieve agents, providers, sessions, and audit metadata.
- Integration: File-memory adapter rejects path traversal and unauthorized files.
- Security: Protected routes return 401 without a valid session/token.
- Acceptance: UI flow creates one agent, one provider config, reads a GSD file, and generates one handoff file.
- Regression: Phase 0 documentation validation remains clean.

## 15.7 Security Validations
- Threat 1: Path traversal through memory-file reads.
  - Mitigation: Canonicalize paths and restrict reads to approved project-root files.
- Threat 2: Secret leakage into logs or handoff files.
  - Mitigation: Redaction tests and denylisted field checks.
- Threat 3: Unauthenticated local API access.
  - Mitigation: Deny-by-default protected routes and auth/session tests.
- Threat 4: Prompt or handoff injection through Markdown memory.
  - Mitigation: Treat memory text as data; preserve source labels; avoid executing embedded instructions automatically.
- Threat 5: Provider-mode confusion between local and API models.
  - Mitigation: Explicit provider mode and adapter contract tests.

## 15.8 Developer Validation Steps
1. Confirm Docker is active.
2. Confirm Go and Node.js are installed.
3. Create tests before production code.
4. Run backend unit/integration tests.
5. Run frontend component/acceptance tests.
6. Validate protected routes return 401 without auth.
7. Validate generated handoff files contain no secrets.
8. Record all results in `DOC/TESTS.md`.

## 15.9 Responsible and Dependencies
- **Responsible:** [BA] and [CTO] for planning/kickoff; [SECURITY] for threat model; [DEV_BACKEND], [DEV_FRONTEND], [DBA], [DEVOPS] for implementation after kickoff.
- **Dependencies:** Phase 0 closed; Docker available; Go available; Node.js available; PostgreSQL via Compose.
- **Time Estimate:** To be estimated by CTO during kickoff handoff.
- **Extra-Code Prerequisites:** Local `.env` must be created outside git before runtime tests.

## 15.10 Risks and Mitigations
- Risk 1: MVP scope expands into Phase 2 local runtime.
  - Mitigation: Keep local runtime as configurable provider path only in Phase 1.
- Risk 2: File-memory adapter can corrupt docs.
  - Mitigation: MVP reads approved GSD memory files; writes only to `.agent_handoff/` unless explicitly approved.
- Risk 3: UI becomes decorative instead of operational.
  - Mitigation: Use dense workbench design from `DOC/DESIGN.md`.

## 15.11 Plan Validation
- **Business Feasibility:** Approved for MVP planning.
- **Technical Feasibility:** Approved pending environment verification.
- **Security Readiness:** Threat model drafted; SECURITY must validate before implementation.
- **Approval:** Implementation authorized by [HUMAN] after kickoff handoff.

## 15.12 Stage Closure Gate Placeholder
- Phase 1 closure gate will be filled only after implementation, review, QA, SECURITY audit, and [HUMAN] approval.
