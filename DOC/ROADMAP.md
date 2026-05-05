# ROADMAP.md - PF_ai GSD Execution Planning

This file tracks roadmap progress. [HUMAN:Ulisses] is the only authority for marking phases complete and authorizing advancement.

## Phase 0: Onboarding and Foundation (Status: Done)
- [x] **Project Onboarding**
  - Criteria: `DOC/PROJECT.md`, `README.md`, `DOC/DESIGN.md`, `DOC/ARCHITECTURE.md`, `DOC/ENV_SETUP.md`, `DOC/ROADMAP.md`, and `DOC/PLAN.md` filled from [HUMAN] answers.
  - Tests: Placeholder scan and documentation consistency review.
  - Responsible: [CEO].
- [x] **Architecture Baseline**
  - Criteria: Go backend, TypeScript web, PostgreSQL, Docker, API/local model provider strategy, and handoff evolution path documented.
  - Tests: Architecture review before implementation.
  - Responsible: [CTO].
- [x] **Security Baseline**
  - Criteria: Deny-by-default routes, secret handling, sanitized logs, and local runtime isolation documented.
  - Tests: Threat model before MVP implementation.
  - Responsible: [SECURITY].

## Phase 1: MVP - Agent, Model, Memory, and File Handoffs (Status: Doing)
- [x] **MVP Planning**
  - Criteria: Acceptance criteria, tests, business rules, and dependencies finalized in `DOC/PLAN.md`.
  - Tests: BA/CTO validation and SECURITY threat modeling.
  - Responsible: [BA] / [CTO].
- [x] **MVP Backend** (Status: Done)
  - Criteria: Go API supports agent registry, model/provider registry, GSD memory reads, and file handoff creation.
  - Tests: Unit and integration tests created before production code.
  - Responsible: [DEV_BACKEND].
- [x] **MVP Frontend** (Status: Done)
  - Criteria: Web UI shows active phase/status, agents, models, key GSD memory files, and handoff actions.
  - Tests: Component and acceptance tests.
  - Responsible: [DEV_FRONTEND].
- [x] **MVP Persistence** (Status: Done)
  - Criteria: PostgreSQL schema stores agents, providers, sessions, and audit metadata.
  - Tests: Migration and repository contract tests.
  - Responsible: [DBA].
- [ ] **MVP Validation**
  - Criteria: Human can operate at least one agent through file handoff flow and inspect project state.
  - Tests: QA, SECURITY, and CODE_REVIEWER closure review.
  - Responsible: [QA] / [SECURITY] / [CODE_REVIEWER].

## Phase 2: Local Runtime and Hybrid Provider Routing (Status: To Do)
- [ ] **Local Model Runtime**
  - Criteria: Docker service can run the selected local model endpoint.
  - Tests: Health check, timeout handling, and resource validation.
- [ ] **Hybrid Routing**
  - Criteria: Agent execution can choose API, local, or hybrid provider mode.
  - Tests: Provider adapter contract tests.
- [ ] **Operational Hardening**
  - Criteria: Logs, monitoring hooks, secrets handling, and fallback behavior validated.
  - Tests: Security and reliability tests.

## Phase 3: MCP Handoffs and Advanced Orchestration (Status: To Do)
- [ ] **MCP Handoff Transport**
  - Criteria: File handoff schema maps cleanly to MCP-based communication.
  - Tests: Transport compatibility and regression tests.
- [ ] **Advanced Agent Manager**
  - Criteria: Agent routing, memory selection, model tiering, and skill mapping are managed through the UI/API.
  - Tests: End-to-end orchestration scenarios.
- [ ] **Release Readiness**
  - Criteria: Documentation, tests, observability, and release notes complete.
  - Tests: Stage Closure Gate.

## General Notes
- [HUMAN:Ulisses] must approve every Stage Closure Gate.
- `CONTEXT.md` updates occur after onboarding completion or phase intake, per [HUMAN] instruction.
- `wiki/` updates occur only at phase end unless [HUMAN] explicitly requests otherwise.
- MVP uses file-based handoffs; MCP is final-stage scope.
- Chat responses should remain compact in the Caveman style; details belong in files.

## Version History (Immutable)
- **[2026-05-05 15:06] - [CEO]:** Roadmap filled from onboarding answers and MVP confirmation.
- **[2026-05-05 15:16] - [CEO]:** Phase 0 task items marked complete; phase status set to Closure Pending until [HUMAN] approval.
- **[2026-05-05 15:20] - [CEO]:** [HUMAN:Ulisses] approved Phase 0 roadmap closure; phase status set to Done.
- **[2026-05-05 15:23] - [CEO]:** Phase 1 planning opened after `[USER_DONE]`; implementation remains blocked until Phase 1 plan, threat model, and kickoff handoff are complete.
- **[2026-05-05 15:30] - [CEO]:** Phase 1 implementation authorized by [HUMAN]; planning marked complete; backend, frontend, and persistence implementation started.
- **[2026-05-05 18:50] - [DBA]:** Phase 1 MVP Persistence marked Done after migration validation and repository contract tests for agents, providers, sessions, and audit events passed.
- **[2026-05-05 18:56] - [DEV_FRONTEND]:** Phase 1 MVP Frontend marked Done after static checks and local runtime validation at `http://127.0.0.1:5173`.
- **[2026-05-05 19:02] - [DEV_BACKEND]:** Phase 1 MVP Backend marked Done after HTTP handler acceptance flow validated provider, agent, memory, and handoff behavior.
- **[2026-05-05 19:03] - [TECH_LEAD]:** Phase 1 parallel review consolidated; no blocking implementation defect found, but OS-level backend background runtime validation remains unresolved before Stage Closure Gate.
