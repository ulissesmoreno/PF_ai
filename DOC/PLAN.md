# PLAN.md

This file must be filled **completely and without fail** for each roadmap phase, ensuring consistency, traceability, and total documentation. Fill at phase closure — not before. Reminder: Tests are mandatory (TDD).

## Updated on: YYYY-MM-DD HH:MM:SS

---

## 1. Phase Objective
- Describe the purpose of this roadmap phase.
- What must be delivered and why.

## 2. Value Delivered
> **What the user/client can do after this phase that they could not before.**
- [Clear, non-technical description of the value delivered.]
- If this phase has no direct user-facing value (e.g., infrastructure): justify here why it enables future value.

## 3. Roadmap Reference
- Phase name: [E.g.: MVP — Phase 1]
- Related ROADMAP item: [Link to ROADMAP.md#Phase-X]

## 4. Detailed Description
- What was implemented: [Details]
- Expected behavior: [Description]
- Impacted components: [List]

## 5. Acceptance Criteria
- Criterion 1: [Measurable]
- Criterion 2: [Measurable]
- Criterion 3: [Measurable]

## 6. Business Rules
- Rule 1: [Description]
- Rule 2: [Description]

## 7. Recommended Tests
- Unit test 1: description and coverage.
- Integration test 1: description and coverage.
- Acceptance test 1: description and scenario.
- Security test 1: mandatory validations.
- Regression test 1: critical scenario.

## 8. Security Validations
- Validation 1: [E.g.: JWT active]
- Validation 2: [E.g.: Data sanitized]
- Validation 3: [E.g.: Zero leak]
- Security notes: [Notes]

## 9. Responsible and Dependencies
- Responsible: [Agent:Level]
- Dependencies: [Prerequisites]
- Time Estimate: [Hours/days]
- **Extra-Code Prerequisites:** [API keys, environment variables, external services needed]

## 10. Risks and Mitigations
- Risk 1: [Description]
  - Mitigation: [Action]
- Risk 2: [Description]
  - Mitigation: [Action]

## 11. Integration with Other Files
- **ROADMAP.md:** Related phase.
- **STATE.md:** Current progress.
- **TESTS.md:** Validations performed.

## 12. Version History (Immutable)
- **[YYYY-MM-DD HH:MM] — [Responsible]:** Initial phase version.
- [Add new entries here without deleting.]

## 13. Plan Validation
- **Technical Feasibility:** [Confirm if components and stack support the phase.]
- **Value Alignment:** [Verify if "Value Delivered" meets PROJECT.md objectives and client expectations.]
- **Risks Assessed:** [Review mitigations; register in QUESTIONS.md if open questions.]
- **Approval:** [Status: Approved / Rejected — Reason.]

---

## 14. Phase Closure Gate ⛔

> Filled by the agent at phase closure. Presented to [HUMAN] for approval. Advancement blocked until [HUMAN] confirms.

### Closed Phase: [Phase name]
### Closure Date: [YYYY-MM-DD HH:MM]
### Responsible: [Agent]

### Value Delivered Confirmation
- [ ] Value Delivered as described in §2 was achieved and is usable.
- [ ] [HUMAN] can validate the delivered value without additional setup.

### Closure Checklist

| # | File | Status | Note |
| :- | :--- | :---: | :--- |
| 1 | `TESTS.md` | `[ ]` | All tests recorded with real results and timestamp |
| 2 | `STATE.md` | `[ ]` | Updated with delivery description and metrics |
| 3 | `TASKS.md` | `[ ]` | All phase tasks marked Done *(Medium+)* |
| 4 | `CONTEXT.md` | `[ ]` | Decisions documented *(Medium+)* |
| 5 | `README.md` | `[ ]` | Reflects delivered code reality |
| 6 | `ROADMAP.md` | `[ ]` | Stage marked `[x]` by [HUMAN] *(Medium+)* |
| 7 | `VERSIONS.md` | `[ ]` | Entry added *(Medium+, Phase 1+)* |
| 8 | `WIKI` | `[ ]` | Updated with phase deliverables *(Medium+)* |
| 9 | `RETROSPECTIVE.md` | `[ ]` | Phase entry consolidated by CEO *(Medium+)* |

### Pending Items / Blockers
- [Describe any incomplete item and reason]

### Delivery Summary
- [Brief description of what was delivered]
- Metrics: [E.g.: mutation score 82%, 0 security failures]

### ✅ User Approval
> **Awaiting [HUMAN] confirmation to advance to the next phase.**
- [ ] [HUMAN] confirmed phase closure and authorizes start of next phase.
