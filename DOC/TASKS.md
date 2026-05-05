# TASKS

> Immutable registry of all atomic project activities. Append only — never delete entries.
> Updated after every task change. Timestamps mandatory (YYYY-MM-DD HH:MM).

## Updated on: YYYY-MM-DD HH:MM:SS

---

## Active Tasks

### Task [ID]: [Description]
- **Status:** To Do / Doing / Done
- **Assigned to:** [Agent:Level]
- **Priority:** High / Medium / Low
- **Deadline:** [Date]
- **Estimated time:** [Hours]
- **Actual time:** [Hours — filled on completion]
- **Dependencies:** [List]
- **Phase reference:** [ROADMAP.md#Phase-X]
- **Business value:** [What problem this solves for the user]
- **Notes:** [Notes]

---

## Completed Tasks

### [YYYY-MM-DD] — Task [ID]: [Description]
- **Status:** Done
- **Assigned to:** [Agent:Level]
- **Completed on:** [YYYY-MM-DD HH:MM]
- **Estimated time:** [Hours]
- **Actual time:** [Hours]
- **Variance:** [+/- Hours — for future estimation calibration]
- **Notes:** [Result]

---

## Blocked Tasks

### Task [ID]: [Description]
- **Status:** Blocked
- **Assigned to:** [Agent:Level]
- **Reason:** [E.g.: External dependency / Awaiting CLARIFICATION_REQUEST response]
- **Action:** [Registered in QUESTIONS.md / Handoff sent]

---

## Bugfix Tasks (Priority Override)

### [BUGFIX] Task [ID]: [Description]
- **Status:** To Do / Doing / Done
- **Severity:** Critical / High / Medium
- **Assigned to:** [Agent:Level]
- **Branch:** `bugfix/phase-X-description`
- **Next phase frozen:** Yes / No
- **Root cause:** [Description]
- **Opened by:** [TECH_LEAD — YYYY-MM-DD HH:MM]

---

## Estimation Calibration Log
> Tracks estimated vs. actual time to improve future estimates.

| Task ID | Estimated | Actual | Variance | Agent Level |
| :--- | :--- | :--- | :--- | :--- |
| [ID] | [Xh] | [Yh] | [+/-Zh] | [Junior/Pleno/Senior] |

---

## Version History (Immutable)
- **[YYYY-MM-DD HH:MM] — CEO:** Template updated. Estimated vs. actual time field added. Bugfix task section added. Business value field added.
- [Add new entries here without deleting.]

---

## Phase 0 Task Entries

### [2026-05-05] - Task PHASE0-DOC-001: Complete PF_ai onboarding documentation
- **Status:** Done
- **Assigned to:** [CEO]
- **Completed on:** 2026-05-05 15:16
- **Estimated time:** Same work cycle
- **Actual time:** Same work cycle
- **Variance:** 0
- **Dependencies:** [HUMAN] onboarding answers in `QUESTIONS.md`.
- **Phase reference:** `DOC/ROADMAP.md#phase-0-onboarding-and-foundation-status-doing`
- **Business value:** Establishes validated project context before implementation.
- **Notes:** Required onboarding files filled and validated. Roadmap advancement remains blocked until [HUMAN] approves Stage Closure Gate.

### [2026-05-05] - Task PHASE0-ARCH-001: Confirm Go and Hexagonal Architecture
- **Status:** Done
- **Assigned to:** [CEO] / [CTO]
- **Completed on:** 2026-05-05 15:16
- **Estimated time:** Same work cycle
- **Actual time:** Same work cycle
- **Variance:** 0
- **Dependencies:** `NEW-INSTRUCTIONS.md` question.
- **Phase reference:** `DOC/ROADMAP.md#phase-0-onboarding-and-foundation-status-doing`
- **Business value:** Confirms the backend architecture remains testable, isolated, and aligned with PF_ai.
- **Notes:** Answer registered in `QUESTIONS.md`.

---

## Phase 1 Task Entries

### [2026-05-05] - Task PHASE1-PLAN-001: Prepare MVP plan and threat model
- **Status:** Done
- **Assigned to:** [CEO] / [BA] / [CTO] / [SECURITY]
- **Completed on:** 2026-05-05 15:23
- **Estimated time:** Same work cycle
- **Actual time:** Same work cycle
- **Variance:** 0
- **Dependencies:** Phase 0 Done.
- **Phase reference:** `DOC/ROADMAP.md#phase-1-mvp---agent-model-memory-and-file-handoffs-status-doing`
- **Business value:** Converts approved roadmap into implementation-ready acceptance criteria and security constraints.
- **Notes:** Implementation remains blocked until kickoff handoff and [HUMAN] authorization.

### [2026-05-05] - Task PHASE1-KICKOFF-001: Create technical kickoff handoff
- **Status:** Done
- **Assigned to:** [CEO] / [CTO]
- **Completed on:** 2026-05-05 15:23
- **Priority:** High
- **Deadline:** Before any Phase 1 code.
- **Estimated time:** Same work cycle
- **Actual time:** Same work cycle
- **Variance:** 0
- **Dependencies:** Phase 1 plan and threat model.
- **Phase reference:** `DOC/ROADMAP.md#phase-1-mvp---agent-model-memory-and-file-handoffs-status-doing`
- **Business value:** Provides one authoritative instruction package for technical agents.
- **Notes:** Handoff created at `.agent_handoff/2026-05-05_1523_phase1_kickoff.json`. Implementation remains blocked until [HUMAN] authorizes code.

### [2026-05-05] - Task PHASE1-ENV-001: Verify local development prerequisites
- **Status:** Done
- **Assigned to:** [DEVOPS]
- **Priority:** High
- **Deadline:** Before Phase 1 implementation.
- **Estimated time:** Same work cycle
- **Actual time:** Same work cycle
- **Dependencies:** Docker, Go, Node.js availability.
- **Phase reference:** `DOC/ROADMAP.md#phase-1-mvp---agent-model-memory-and-file-handoffs-status-doing`
- **Business value:** Confirms the project can run backend, frontend, and database locally.
- **Reason:** Go runtime is available through `C:\Program Files\Go\bin\go.exe`; Docker daemon access works with elevated permission.
- **Action:** Go tests passed with workspace-local `GOCACHE`; PostgreSQL container and migration validated.

### [2026-05-05] - Task PHASE1-BACKEND-001: Create Go Hexagonal MVP foundation
- **Status:** Doing
- **Assigned to:** [DEV_BACKEND:Pleno]
- **Priority:** High
- **Deadline:** Phase 1.
- **Estimated time:** TBD after Go runtime is available.
- **Dependencies:** Phase 1 authorization; Go runtime for test execution.
- **Phase reference:** `DOC/ROADMAP.md#phase-1-mvp---agent-model-memory-and-file-handoffs-status-doing`
- **Business value:** Provides the backend foundation for agent/provider/memory/handoff operations.
- **Notes:** Domain, application ports, file-memory adapter, handoff writer, HTTP API, and tests were created. Go test suite passed with workspace-local `GOCACHE`.

### [2026-05-05] - Task PHASE1-FRONTEND-001: Create MVP web workbench shell
- **Status:** Doing
- **Assigned to:** [DEV_FRONTEND:Pleno]
- **Priority:** Medium
- **Deadline:** Phase 1.
- **Estimated time:** Same work cycle
- **Dependencies:** Design tokens in `DOC/DESIGN.md`.
- **Phase reference:** `DOC/ROADMAP.md#phase-1-mvp---agent-model-memory-and-file-handoffs-status-doing`
- **Business value:** Gives [HUMAN] a first operational view of phase/status, agents, providers, memory, and handoff preview.
- **Notes:** Static web shell created; JavaScript syntax validation passed.

### [2026-05-05] - Task PHASE1-DB-001: Create MVP PostgreSQL migration
- **Status:** Done
- **Assigned to:** [DBA:Pleno]
- **Priority:** High
- **Deadline:** Phase 1.
- **Estimated time:** Same work cycle
- **Dependencies:** PostgreSQL via Docker Compose.
- **Phase reference:** `DOC/ROADMAP.md#phase-1-mvp---agent-model-memory-and-file-handoffs-status-doing`
- **Business value:** Establishes durable state for agents, providers, sessions, and audit events.
- **Notes:** Initial migration created and validated against `pf-ai-postgres`. Repository-level persistence remains a separate implementation task.

### [2026-05-05] - Task PHASE1-FRONTEND-002: Validate local web runtime
- **Status:** Blocked
- **Assigned to:** [DEV_FRONTEND:Pleno] / [DEVOPS:Pleno]
- **Priority:** Medium
- **Deadline:** Before Phase 1 closure.
- **Estimated time:** Same work cycle
- **Dependencies:** Durable local process runner or interactive terminal.
- **Phase reference:** `DOC/ROADMAP.md#phase-1-mvp---agent-model-memory-and-file-handoffs-status-doing`
- **Business value:** Allows [HUMAN] to inspect the MVP workbench in browser.
- **Reason:** Background launch did not keep the Node server active from the current shell.
- **Action:** Run `cd pf-ai-web && npm start` interactively or provide an approved persistent process runner.
