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
