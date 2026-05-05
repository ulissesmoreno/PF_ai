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
