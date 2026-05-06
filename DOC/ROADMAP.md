# ROADMAP.md - GSD Execution Planning

This file must be **started or updated mandatorily** with each new phase, ensuring the plan reflects the most recent context without loss of traceability. Every phase represents a **value deliverable** — something functional and usable by the end user or client.

> **Rule:** Advancement to the next phase requires explicit [HUMAN] approval. Stages are marked `[x]` exclusively by [HUMAN].

---

## 🟢 PHASE 0: Foundation and Infra (Status: To Do)

- **Value Delivered:** Development environment ready; project structure validated. No user-facing value — justified as enabler for Phase 1.
- [ ] **Initial Configuration:** Fill PROJECT.md with vision, objectives, and stack.
  - Criteria: All placeholders replaced; stack validated.
  - Tests: Environment validation.
- [ ] **Environment Setup:** Install technologies from the defined stack.
  - Criteria: Functional environment; Docker active.
  - Tests: Security (JWT configured).
  - Dependencies: PROJECT.md complete.
- [ ] **Base Architecture:** Define components per ARCHITECTURE.md.
  - Criteria: Documented and validated.
  - Tests: Hexagonal isolation check.
  - DoD: Files filled without placeholders.

---

## 🟡 PHASE 1: MVP (Status: To Do)

- **Value Delivered:** [Describe what the user/client can do after this phase — the core hypothesis validated.]
- [ ] **MVP Planning:** Define minimum features to validate the core hypothesis.
  - Criteria: Clear scope; hypotheses defined.
  - Tests: Functionality and security before code.
- [ ] **MVP Implementation:** Develop core features with TDD.
  - Criteria: Executable code; tests passing.
  - Tests: Unit/integration coverage.
  - Dependencies: Phase 0 complete.
- [ ] **MVP Validation:** Feasibility tests; adjustments via NEW-INSTRUCTIONS.md if needed.
  - Criteria: Core hypothesis validated; metrics achieved (see PROJECT.md KPIs).
  - Tests: Security and performance.
  - DoD: STATE.md updated.

---

## 🔴 PHASE 2: [Phase Name] (Status: To Do)

- **Value Delivered:** [Describe the incremental value the user/client gains — what they can do now that they couldn't after Phase 1.]
- [ ] **Feature:** [Description]
  - Criteria: [Measurable]
  - Tests: [Required validations]
  - Dependencies: Phase 1 complete.

---

## 🔴 PHASE 3: [Phase Name] (Status: To Do)

- **Value Delivered:** [Describe incremental value.]
- [ ] **Feature:** [Description]
  - Criteria: [Measurable]
  - Tests: [Required validations]
  - Dependencies: Phase 2 complete.

---

## 📋 General Notes

- **Value First:** Every phase must deliver something usable. Purely technical phases must justify the absence of user-facing value in PLAN.md.
- **MVP is mandatory:** Phase 1 always validates the core hypothesis before expansion.
- **Approval:** Advancement authorized exclusively by [HUMAN]. Previous phase must be marked `[x]` before next begins.
- **Updates:** Updated at phase start and closure. Register in QUESTIONS.md if open questions exist.
- **Traceability:** Each phase with status, value delivered, responsible, and dependencies.
- **References:** PROJECT.md (objectives), TESTS.md (validations), STATE.md (progress), PLAN.md (detail).

## 📋 Version History (Immutable)
- **[YYYY-MM-DD HH:MM] — [Responsible]:** Initial ROADMAP created.
- [Add new entries here without deleting.]
