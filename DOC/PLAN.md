# PLAN

This file must be filled **completely and without fail** for each roadmap stage, ensuring consistency, traceability, and total documentation. Nothing can be lost along the way. Reminder: Fill before code; tests are mandatory (TDD).

## Updated on: YYYY-MM-DD HH:MM:SS

## 1. Stage Objective
- Describe the purpose of this roadmap phase.
- What must be delivered and why.

## 2. Roadmap Stage
- Stage name: [E.g.: MVP Implementation]
- Related to which roadmap item: [Link to ROADMAP.md]

## 3. Detailed Description
- What will be implemented: [Details]
- Expected behavior: [Description]
- Impacted components: [List]

## 4. Acceptance Criteria
- Criterion 1: [Measurable]
- Criterion 2: [Measurable]
- Criterion 3: [Measurable]

## 5. Business Rules
- Rule 1: [Description]
- Rule 2: [Description]
- Rule 3: [Description]

## 6. Recommended Tests
- Unit test 1: description and coverage.
- Integration test 1: description and coverage.
- Acceptance test 1: description and scenario.
- Security test 1: mandatory validations.
- Regression test 1: critical scenario.

## 7. Security Validations
- Validation 1: [E.g.: JWT active]
- Validation 2: [E.g.: Data sanitized]
- Validation 3: [E.g.: Zero leak]
- Security notes: [Notes]

## 8. Developer Validation Steps
1. Prepare environment: [Instructions]
2. Execute unit/module test: [Commands]
3. Execute integration test: [Commands]
4. Verify acceptance criteria: [How]
5. Validate business rules: [How]
6. Validate security: [How]
7. Document results in `tests.md`: [How]

## 9. Responsible and Dependencies
- Responsible: [Name]
- Dependencies: [Prerequisites, e.g.: Previous phase complete]
- Time Estimate: [Hours/days]
- **Extra-Code Prerequisites:** [Non-code aspects that must be ready, e.g.: API keys, configured environment variables. Notify developer to arrange before starting.]

## 10. Risks and Mitigations
- Risk 1: [Description]
  - Mitigation: [Action]
- Risk 2: [Description]
  - Mitigation: [Action]

## 11. Integration with Other Files
- **ROADMAP.md:** Related stage.
- **STATE.md:** Current progress.
- **TESTS.md:** Validations performed.

## 12. Version History (Immutable)
- **[YYYY-MM-DD HH:MM] - Responsible:** Initial stage version.
- [Add new entries here without deleting.]

## 13. Plan Validation
- **Technical Feasibility:** [Confirm if components and stack support the stage.]
- **Project Alignment:** [Verify if it meets PROJECT.md objectives.]
- **Risks Assessed:** [Review mitigations; register in QUESTIONS.md if questions.]
- **Approval:** [Status: Approved / Rejected - Reason.]
- **ROADMAP.md Reference:** This plan is validated here; ROADMAP.md serves only for development tracking (To Do/Doing/Done status).

## 9. Dependencies and Preconditions
- Dependency 1:
- Dependency 2:
- Precondition 1:
- Precondition 2:

## 10. Risks and Attention
- Risk 1:
- Risk 2:
- Mitigation:

## 11. References
- Related files:
- Technical documents:
- Additional notes:

---

## 14. Stage Closure Gate ⛔ (Fill before advancing)

> This section must be **filled in by the agent at the end of the stage** and presented to the user for approval. Advancement to the next stage is **blocked** until the user explicitly confirms.

### Completed Stage: [Stage name — e.g.: MVP Phase 1]
### Closure Date: [YYYY-MM-DD HH:MM]
### Responsible: [Name]

### Closure Checklist

| # | File | Status | Note |
| :- | :--- | :---: | :--- |
| 1 | `TESTS.md` | `[ ]` | [All tests recorded with real results and timestamp?] |
| 2 | `STATE.md` | `[ ]` | ["What Was Completed" section updated?] |
| 3 | `TASKS.md` | `[ ]` | [All stage tasks marked as Done?] |
| 4 | `CONTEXT.md` | `[ ]` | [Stage context and decisions documented?] |
| 5 | `README.md` | `[ ]` | [Reflects the reality of the delivered code?] |
| 6 | `ROADMAP.md` | `[ ]` | [Stage marked as `[x]` by the user?] |
| 7 | `VERSIONS.md` | `[ ]` | [New entry added if there was a version delivery?] |

### Pending Items / Blockers
- [Describe any incomplete item and the reason]

### Delivery Summary
- [Brief description of what was delivered in this stage]
- Metrics: [E.g.: 82% test coverage, 0 security failures]

### ✅ User Approval
> **Awaiting user confirmation to advance to the next stage.**
- [ ] User confirmed stage closure and authorizes start of the next stage.
