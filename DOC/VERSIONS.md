# VERSIONS.md

> **Owner:** Agent who closes the phase (updates at Stage Closure Gate).
> Append only — never delete entries. Timestamps mandatory.

## 1. Objective
Register project versions, rollbacks, bugfixes, and significant incidents with full traceability.

## 2. Versioning Convention
- SemVer: `MAJOR.MINOR.PATCH`
- Pre-releases: `1.0.0-alpha.1`, `1.0.0-beta.1`
- Rollback entries: `[ROLLBACK] v1.2.0 → v1.1.0`
- Bugfix entries: `[BUGFIX] vX.X.X-fix.1`

## 3. Commit Convention
- Format: `type(scope): short description`
- Types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `rollback`, `bugfix`
- Reference plan/task when applicable.

## 4. Entry Structure

### [Version] — YYYY-MM-DD — [Closing Agent]
- **Type:** Release / Rollback / Bugfix / Incident
- **Main Commit:** hash or reference
- **Description:** Summary of changes or incident
- **Files Changed:** [List]
- **Tests / Validation:** Reference to `TESTS.md`
- **Security Clearance:** SECURITY audit status
- **Deploy Notes:** Environment variables, release instructions
- **Retrospective Reference:** `RETROSPECTIVE.md#Phase-X` (if applicable)

---

## 5. Release History

### 1.1.0 — 2026-04-29 — [CEO]
- **Type:** Release
- **Main Commit:** boilerplate-finalization
- **Description:** GSD Boilerplate v1.1.0. Mapping and closure of all pending items. Introduction of ONBOARDING.md, stabilization of GSD-RULES.md, wiki structure, frontend architecture.
- **Files Changed:** `TODO.md`, `GSD-RULES.md`, `ARCHITECTURE.md`, `ENV_SETUP.md`, `PLAYBOOK.md`, `ONBOARDING.md`, `WIKI.md`, `DESIGN.md`
- **Tests / Validation:** Structure validated through technical review sessions.
- **Deploy Notes:** Boilerplate consolidated for new projects.

### 1.0.0 — YYYY-MM-DD — [Agent]
- **Type:** Release
- **Main Commit:** [hash]
- **Description:** Initial boilerplate release.
- **Files Changed:** `README.md`, `GSD-RULES.md`, `PROJECT.md`, `ROADMAP.md`, `TESTS.md`, `STATE.md`, `ENV_SETUP.md`

---

## 6. Rollback Entry Template

### [ROLLBACK] vX.X.X → vY.Y.Y — YYYY-MM-DD — [DEVOPS:Level]
- **Severity:** Critical / High
- **Trigger:** [What caused the rollback]
- **Authorized by:** [HUMAN] — [YYYY-MM-DD HH:MM]
- **Previous stable version:** vY.Y.Y
- **Retrospective Reference:** `RETROSPECTIVE.md#Phase-X`

---

## 7. Bugfix Entry Template

### [BUGFIX] vX.X.X-fix.N — YYYY-MM-DD — [Closing Agent]
- **Severity:** Critical / High / Medium
- **Branch:** `bugfix/phase-X-description`
- **Root Cause:** [Description]
- **Next Phase Frozen:** Yes / No
- **Merged into:** main + next phase branch
- **Retrospective Reference:** `RETROSPECTIVE.md#Phase-X`

---

## 8. Update History (Immutable)
- **[YYYY-MM-DD HH:MM] — CEO:** Template updated. Phase-closing agent ownership added. Rollback and bugfix entry templates added.
- [Add new entries here without deleting.]

---

## 9. PF_ai Release Entries

### 1.2.0-alpha.1 - 2026-05-05 - [CEO]
- **Type:** Release
- **Main Commit:** pending
- **Description:** PF_ai project onboarding milestone. Product vision, Go architecture, TypeScript web direction, PostgreSQL persistence, Docker local runtime path, API/local provider strategy, file handoffs for MVP, and MCP final-stage direction documented.
- **Files Changed:** `DOC/PROJECT.md`, `README.md`, `DOC/DESIGN.md`, `DOC/ROADMAP.md`, `DOC/PLAN.md`, `DOC/ENV_SETUP.md`, `DOC/ARCHITECTURE.md`, `QUESTIONS.md`, `PLAYBOOK.md`, `DOC/CONTEXT.md`, `DOC/TESTS.md`, `DOC/STATE.md`, `DOC/TASKS.md`.
- **Tests / Validation:** `DOC/TESTS.md#8-phase-0-validation-results`.
- **Security Clearance:** Documentation secret scan passed; implementation threat model pending before Phase 1.
- **Deploy Notes:** No deployable code yet.
- **Retrospective Reference:** Pending Phase 0 retrospective.
