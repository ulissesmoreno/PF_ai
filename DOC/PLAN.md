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
