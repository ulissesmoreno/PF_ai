# PLAYBOOK.md — Developer Preferences & Working Style

> **Immutability Rule:** Append only. Never delete or overwrite existing entries.
> Most recent entry for any topic is the active preference.
> Format: `[YYYY-MM-DD HH:MM] — Entry`
>
> **When to Update (mandatory triggers):**
> 1. After every answer in `QUESTIONS.md` — if the answer reveals a preference or decision pattern.
> 2. After every new entry in `NEW-INSTRUCTIONS.md` — if instructions reflect how [HUMAN] likes to work.
> 3. After `PROJECT.md` is defined — review stack, architecture, and process choices.
> 4. When any agent handoff contains `playbook_update: true` — CEO appends the revealed preference.
> 5. After `RETROSPECTIVE.md` phase entry — if a lesson reveals a [HUMAN] work preference.
>
> **Owner:** CEO identifies the trigger and appends. [HUMAN] is the source of truth.

---

## 1. Working Philosophy

| Timestamp | Principle | Notes |
| :--- | :--- | :--- |
| [2026-04-25 21:25] | **GSD (Get Shit Done)** — delivery over perfection, but quality by design. | Functional code ships. Untested code doesn't. |
| [2026-04-25 21:25] | **Traceability above speed.** Every decision, question, and context must be documented. | Nothing lost. Nothing assumed. |
| [2026-04-25 21:25] | **Documentation is not optional.** If it's not documented, it doesn't exist. | All `.md` files complete before "Done". |
| [2026-04-25 21:25] | **Immutability of history.** Files are append-only. History is never erased. | Timestamps on all entries. |
| [2026-04-25 21:25] | **User controls the roadmap.** AI proposes; user approves before advancing. | No stage starts without explicit user OK. |
| [2026-05-05 00:00] | **Chat responses: only the essential.** No narration of what is being done — state it in one line. | Details go into `.md` files. |
| [2026-05-05 00:00] | **Token economy is a design goal.** Multi-agent exists to minimize token consumption, not to add process. | Always assign lowest sufficient seniority. |
| [2026-05-05 00:00] | **Security is a structural foundation.** Not a gate — integrated at kickoff and closure of every phase. | Threat modeling before implementation. |

---

## 2. Architecture Defaults

| Timestamp | Decision | Rationale |
| :--- | :--- | :--- |
| [2026-04-25 21:25] | **Hexagonal Architecture** as default. | Testability, isolation, framework independence. |
| [2026-04-25 21:25] | **Domain layer has zero external dependencies.** | Pure business logic, maximally testable. |
| [2026-04-25 21:25] | **API routes born deny-all.** | Security-first by default. |
| [2026-04-25 21:25] | **MVP planned before any implementation.** | Validates hypotheses before building. |
| [2026-04-25 21:25] | **Directory naming convention:** `<name>-web`, `<name>-app`, `<name>-desktop`, `<name>-service`. | Consistency across projects. |
| [2026-05-05 00:00] | **Phase-based context loading.** CEO + CTO + BA read files once per phase, not per task. | Reduces token consumption. |
| [2026-05-05 00:00] | **Handoff-first communication between agents.** Structured JSON in `.agent_handoff/` over file reading. | More token-efficient than full file reads. |

---

## 3. Tech Stack Defaults

### 3.1 Backend
| Timestamp | Choice | Context |
| :--- | :--- | :--- |
| [2026-04-25 21:25] | **Java 21 + Spring Boot 3** | Primary backend. |
| [2026-04-29 20:35] | **C# / .NET** | Replacing Java for full CRUD in specific projects. |
| [2026-04-25 21:25] | **Python 3.12** | Data processing, ML pipelines. |

### 3.2 Frontend
| Timestamp | Choice | Context |
| :--- | :--- | :--- |
| [2026-04-25 21:25] | **Angular 18 / TypeScript** | Default web frontend. |
| [2026-04-29 20:35] | **React Native (Web + Mobile)** | Replacing Angular when single codebase needed. |

### 3.3 Database
| Timestamp | Choice | Context |
| :--- | :--- | :--- |
| [2026-04-25 21:25] | **PostgreSQL 16** | Default relational database. |

### 3.4 Infrastructure
| Timestamp | Choice | Context |
| :--- | :--- | :--- |
| [2026-04-25 21:25] | **Docker + Docker Compose** | All dev environments containerized. |
| [2026-04-25 21:25] | **GitHub Actions** | Default CI/CD. |
| [2026-04-25 21:25] | **Apache Kafka** | Async messaging when needed. |
| [2026-04-25 21:25] | **Redis** | Cache and volatile data. |

### 3.5 Observability
| Timestamp | Choice | Context |
| :--- | :--- | :--- |
| [2026-04-25 21:25] | **Prometheus + Grafana** | Metrics and dashboards. |

### 3.6 Security
| Timestamp | Choice | Context |
| :--- | :--- | :--- |
| [2026-04-25 21:25] | **Spring Security + JWT** | Stateless auth. |
| [2026-04-25 21:25] | **HashiCorp Vault / env vars** | Secrets never in files. |

---

## 4. Development Process

| Timestamp | Rule | Detail |
| :--- | :--- | :--- |
| [2026-04-25 21:25] | **TDD — strict mode.** | RED → GREEN → REFACTOR → SECURITY. |
| [2026-04-25 21:25] | **SOLID + Clean Code** in every refactor. | Semantic names, single responsibility. |
| [2026-04-25 21:25] | **Structured logs always.** INFO / WARN / ERROR. No sensitive data. | Auditable, not a security liability. |
| [2026-04-25 21:25] | **Stage Closure Gate** before advancing. | Checklist verified and confirmed by [HUMAN]. |
| [2026-04-25 21:25] | **`.env` never committed.** | All credentials in vault or local env vars. |
| [2026-04-25 21:25] | **Atomic commits.** Conventional Commits format. | `feat`, `fix`, `docs`, `refactor`, `test`, `chore`. |
| [2026-04-29 20:35] | **Developer as Manager.** Code, QA, and architecture fully delegated to AI agents. | |
| [2026-04-29 21:07] | **Branching per ROADMAP Phase.** Each phase in isolated branch. | Phase 0 on `main`. Next phases create PRs. |
| [2026-05-05 00:00] | **Parallel review at phase closure.** SECURITY + CODE_REVIEWER + QA run simultaneously. | Not sequential — reduces closure time. |
| [2026-05-05 00:00] | **Mutation testing over coverage percentage.** | Coverage % is outdated. Mutation score ≥ 80% for domain layer. |
| [2026-05-05 00:00] | **Rollback is [HUMAN]-authorized** except Critical severity (DEVOPS may act preventively). | |
| [2026-05-05 00:00] | **Bugfix branch has priority over next phase** for Critical/High severity. | Next phase freezes until bugfix merges. |

---

## 5. AI Agent Preferences

| Timestamp | Preference | Detail |
| :--- | :--- | :--- |
| [2026-04-25 21:25] | **Read all docs before acting.** | No assumptions. No undocumented actions. |
| [2026-04-25 21:25] | **Register all ambiguities in QUESTIONS.md.** | Pending questions block execution. |
| [2026-04-25 21:25] | **Never skip Stage Closure Checklist.** | [HUMAN] approval mandatory. |
| [2026-04-25 21:25] | **Skills first.** Check Antigravity skills before generic code. | Skills take precedence. |
| [2026-04-29 20:59] | **No autonomous execution before approval.** | Fill planning docs and await explicit [HUMAN] authorization. |
| [2026-04-29 20:59] | **Native files over artifacts.** | Plans go in project files, not ephemeral IDE artifacts. |
| [2026-05-05 00:00] | **[HUMAN] communicates only with CEO.** | Active agent signs chat responses. CEO does not paraphrase. |
| [2026-05-05 00:00] | **Chat: only the essential.** | "Creating repository Y" — not "I will now proceed to create repository Y based on Task Z..." |
| [2026-05-05 00:00] | **Phase-based file reading.** | Management reads files once per phase. Re-reads only after [HUMAN] intervention. |
| [2026-05-05 00:00] | **Seniority assigned by CTO** based on BA business complexity input. | BA informs, CTO decides. Never exposed in chat. |
| [2026-05-05 00:00] | **Escalation via handoff, once per phase.** | Technical agents escalate once. Management resolves via handoff if possible. |
| [2026-05-05 00:00] | **NEW-INSTRUCTIONS.md is [HUMAN]-only.** | CEO formats and timestamps — never modifies content. |
| [2026-05-05 00:00] | **PLAYBOOK.md updated on understanding change.** | Not on schedule — triggered when any interaction reveals or changes a [HUMAN] work preference. |

---

## 6. Documentation Style

| Timestamp | Preference | Detail |
| :--- | :--- | :--- |
| [2026-04-25 21:25] | **Markdown everywhere.** | Versionable, portable, readable. |
| [2026-04-25 21:25] | **Obsidian for maintenance wiki.** | `wiki/` and `raw/` directories. |
| [2026-04-25 21:25] | **Immutable history in all files.** | Append-only. Timestamps on every entry. |
| [2026-04-25 21:25] | **DESIGN.md for visual decisions only.** | No inspiration references. |
| [2026-05-05 00:00] | **README.md is a live project document.** | Updated post-onboarding. Not a boilerplate manual. Maintained per project. |
| [2026-05-05 00:00] | **STATE.md owned by technical agents.** | CONTEXT.md owned by management agents. Separate concerns. |
| [2026-05-05 00:00] | **RETROSPECTIVE.md per phase.** | CEO consolidates from handoff notes. Feeds PLAYBOOK.md when preference revealed. |
| [2026-05-05 00:00] | **wiki/user-manual/ for end-user docs.** | Non-technical, task-oriented. Owned by DOCUMENTATION agent. |

---

## 7. Code Quality Gates

| Timestamp | Gate | Threshold |
| :--- | :--- | :--- |
| [2026-04-25 21:25] | Unit test coverage — domain layer | Mutation score ≥ 80% |
| [2026-04-25 21:25] | Application layer | 100% use case coverage (happy + unhappy path) |
| [2026-04-25 21:25] | Infrastructure layer | Contract tests passing |
| [2026-04-25 21:25] | Integration tests | 100% passing before stage closure |
| [2026-04-25 21:25] | Routes without JWT return 401 | 100% |
| [2026-04-25 21:25] | Logs free of sensitive data | 100% |
| [2026-04-25 21:25] | Regressions on previous stages | 0 new failures |
| [2026-05-05 00:00] | Review rounds before escalation | Max 3 rounds |
| [2026-05-05 00:00] | Bugfix: Critical/High next phase freeze | Mandatory |

---

## 8. Entry Log (Immutable)

| Timestamp | Author | Change |
| :--- | :--- | :--- |
| [2026-04-25 21:25] | Developer | Initial PLAYBOOK created. Sections 1–7 populated. |
| [2026-04-30 16:03] | AI Agent | Branching, file verification, agent/task identification preferences added. |
| [2026-05-05 00:00] | CEO | Major update: multi-agent orchestration preferences, phase-based reading, chat conciseness, seniority model, security as foundation, rollback/bugfix protocols, memory ownership, retrospective, PLAYBOOK trigger rules. |

---

*To add preferences: append a new row with current timestamp. Never edit or remove existing rows.*
