# ONBOARDING.md — New Project Setup Guide

> **When to run:** At the very start of every new project — before any code, any roadmap, and any plan.
> The agent runs this guide as a structured conversation with the developer, then automatically
> populates the relevant project files. Nothing is hardcoded here; values come from the developer.
>
> **Language:** Conduct the conversation in the developer's preferred language.
> All file outputs follow the bilingual standard (DOC/ and DOC/).

---

## Purpose

This guide ensures every project starts with a complete, consistent context.
The agent explores 5 blocks with the developer, then uses the answers to fill in:

| Block | Output Files |
| :--- | :--- |
| 1. Project Identity | `PROJECT.md`, `README.md` (name, description, domain, audience) |
| 2. Technology Stack | `PROJECT.md §3`, `ARCHITECTURE.md`, `ENV_SETUP.md` (selected technologies) |
| 3. Visual Identity | `DESIGN.md` (colors, typography, motion, layout) |
| 4. Team & Process | `PROJECT.md §5`, `GSD-RULES.md` context, `PLAYBOOK.md` (if new patterns emerge) |
| 5. Roadmap & MVP | `ROADMAP.md`, `PLAN.md §1-2` (first stage definition) |

---

## Agent Protocol

### Before Starting
- Read `PLAYBOOK.md` to understand the developer's recurring preferences and defaults.
- Inform the developer: *"I will ask about 5 topics to set up the project. After our conversation,
  I will automatically fill in all relevant files. You can answer briefly or in detail — I will
  translate your answers into the right structure."*
- Do not fill any file until the developer has confirmed all 5 blocks are complete.
- Register any uncertainty in `QUESTIONS.md` before proceeding.

---

### Block 1 — Project Identity
Explore: project name, what it does, who it serves, and what problem it solves.
Aim to understand the core value proposition in one or two sentences.

**Fills:** `PROJECT.md §1-2`, `README.md` (Description, Objective, Target Audience)

Key topics to cover:
- The name and its abbreviation (used for directory naming conventions)
- The domain (financial, health, logistics, productivity, etc.)
- The primary user and their main pain point
- What "success" looks like for this project
- Any inspiration or reference products (record in `PROJECT.md` — never in `DESIGN.md`)

---

### Block 2 — Technology Stack
Explore: which technologies the project actually needs, and where defaults should be overridden.

**Fills:** `PROJECT.md §3`, `ARCHITECTURE.md §4 Stack`, `ENV_SETUP.md §2`

Key topics to cover:
- Confirm or replace the boilerplate defaults (Java / Python / Angular / PostgreSQL)
  - Is this a web app, mobile, desktop, or hybrid? → Affects frontend choice
  - Is there a data/ML component? → Confirms Python layer
  - Is the backend purely API-based, or does it need scheduled jobs?
  - Any specific database requirement (relational vs. document vs. time-series)?
- Third-party integrations: authentication providers, payment, messaging, external APIs
- Infrastructure: cloud provider (if any), deployment target (Docker, k8s, serverless)
- Any constraints: existing tech debt, client mandates, team expertise

---

### Block 3 — Visual Identity
Explore: the look, feel, and personality of the product's UI.

**Fills:** `DESIGN.md` (all token sections and rationale prose)

Follow the agent protocol defined in `DESIGN.md` itself (refer to the "Agent protocol" section there).
Do not record inspiration sources in `DESIGN.md` — log them in `PROJECT.md §6 Risks / Notes`.

Key topics to cover:
- Overall visual mood (dark/light, minimal/rich, serious/playful)
- Brand or palette constraints (existing colors, hex codes, brand guide)
- Typography preferences (specific fonts, or general feel: modern/classic/geometric)
- Motion philosophy (subtle, energetic, minimal)
- Layout preference (card-heavy, list-based, sidebar navigation)
- Any design systems or components the developer wants to align with (e.g., Material, Ant Design)

If the developer has no visual preferences, apply sensible defaults from `PLAYBOOK.md §3`
and note that the full design system will be defined iteratively.

---

### Block 4 — Team & Process
Explore: who is working on the project and how they want to operate.

**Fills:** `PROJECT.md §5`, `PLAYBOOK.md` (append any new process decisions with timestamp)

Key topics to cover:
- Who is responsible for each layer (backend, frontend, ML, infra)
- Who approves roadmap stage transitions (the GSD gate)
- Estimated project timeline and delivery cadence
- Any process exceptions to the GSD defaults (e.g., no ML layer, no CI/CD yet)
- Environment readiness: are Docker, database credentials, and API keys available?
  If not, add entries to `QUESTIONS.md` and `ENV_SETUP.md` as blockers.

---

### Block 5 — Roadmap & MVP
Explore: what is the minimum version that validates the core hypothesis, and what are the stages after that.

**Fills:** `ROADMAP.md`, `PLAN.md §1-2` (first stage)

Key topics to cover:
- What is the single most important thing the MVP must prove?
- What are the minimum features for the MVP to be testable by a real user?
- How many stages are expected? (rough count — can be refined later)
- What does "stage 1 done" look like? (Acceptance criteria for the first delivery)
- Any hard deadlines or external dependencies between stages?

---

## Auto-Fill Execution

After the developer confirms all 5 blocks:

1. Fill `PROJECT.md` completely (English) — no placeholders remaining.
2. Fill `README.md` project description sections (English).
3. Fill `DESIGN.md` with collected visual tokens and rationale (English).
4. Create the first entries in `ROADMAP.md` based on Block 5.
5. Fill `PLAN.md §1-2` for the first stage.
6. Update `ENV_SETUP.md §2` with the confirmed stack.
7. Update `ARCHITECTURE.md §4` if the stack deviates from defaults.
8. Append to `PLAYBOOK.md` any new recurring preferences revealed (with timestamp).
9. Register any unresolved items in `QUESTIONS.md`.
10. Add an entry to `wiki/log.md` with the project setup timestamp.

> **Confirm with the developer before writing.** Present a brief summary of what will be
> filled in each file and ask: *"Should I proceed with the auto-fill?"*

---

## Completion Criteria

The onboarding is complete when:
- [ ] All 5 blocks have been explored with the developer
- [ ] All target files have been filled (no `[placeholder]` remaining in the files listed above)
- [ ] `QUESTIONS.md` has entries for any unresolved items
- [ ] Developer has confirmed the auto-filled content is correct
- [ ] `wiki/log.md` has a new entry with timestamp

---

## Version History (Immutable)
- **[2026-04-25 21:42] — Developer:** ONBOARDING.md created. Initial 5-block protocol defined.
- **[2026-05-02 21:55] — Orchestration Agent:** ONBOARDING.md moved from project root to `DOC/` as per TODO.md file reorganization. All references updated across README.md, ROADMAP.md, and agent files.
- [Add new entries here without deleting.]
