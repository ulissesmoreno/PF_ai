# GSD-RULES.md - Technical Execution Protocol (v1.3)

This document defines the inviolable development rules. Any deviation invalidates the "Done" state.

---

## 👑 §0. Agent Orchestration Protocol (Chain of Command)

This section is the **highest-priority rule in the entire framework**. It governs how agents activate, read, and delegate before any execution begins.

### 0.1 CEO Precedence — Inviolable

The **[CEO] Agent is invariably the first to act** in every new work cycle or session. No other agent may begin execution without prior CEO intake and delegation.

**CEO Mandatory Reading Sequence (strict order):**
1. `DOC/GSD-RULES.md` — Align with inviolable rules. This is non-negotiable before any action.
2. `PLAYBOOK.md` *(root)* — Review developer preferences and recurring patterns.
3. `NEW-INSTRUCTIONS.md` *(root)* — Understand the current cycle's user directives and scope.
4. `DOC/PLAN.md` — Validate business context and acceptance criteria before delegating.
5. `DOC/ONBOARDING.md` — **Trigger only at the start of a brand-new project** (when `PROJECT.md` still contains placeholder values). Run the full 5-block protocol before any other action.

> **Rule:** If `NEW-INSTRUCTIONS.md` contains pending questions or ambiguities, the CEO MUST register them in `QUESTIONS.md` and halt. Nothing proceeds with open questions.

### 0.2 Management Layer — CEO → BA / CTO

After CEO intake, execution is delegated to the **Management Layer**:

| Agent | Layer | Primary Responsibility |
| :--- | :--- | :--- |
| **[BA]** | Management | Translates user instructions into structured BDD requirements in `PLAN.md` |
| **[CTO]** | Management | Validates architectural feasibility and assigns technical tasks |

**Management Layer Mandatory Rule:** Both `[CTO]` and `[BA]` MUST read `DOC/GSD-RULES.md` **before executing or validating any task**, to confirm that technical orders comply with the protocol and no constraint is violated.

### 0.3 Technical Agents — Management → Technical

Technical agents receive work **only after** the Management Layer has processed and validated the CEO's delegation:

| Agents | Layer | Scope |
| :--- | :--- | :--- |
| `[DEV_BACKEND]`, `[DEV_FRONTEND]` | Technical | Feature implementation (TDD) |
| `[DBA]` | Technical | Schema, migrations, persistence optimization |
| `[DS/ML]` | Technical | ML models, data pipelines, intelligence layer |
| `[DEVOPS]` | Technical | CI/CD, containers, infrastructure, observability |

### 0.4 Quality & Review Gates

Quality agents act as **mandatory gates** after technical deliveries:

| Agent | When | Authority |
| :--- | :--- | :--- |
| `[SECURITY]` | After every stage; anytime a vulnerability is suspected | **Veto power** — can block Stage Closure Gate |
| `[REVIEWER]` | After every atomic DEV delivery, before QA | **Reject authority** — can block QA handoff |
| `[QA]` | After REVIEWER approval | **Stage approval** — validates acceptance criteria |

### 0.5 Creative & Market Agents

Activated on demand by the CEO or CMO when marketing or design output is required:

| Agent | Activated By | Output |
| :--- | :--- | :--- |
| `[CMO]` | CEO (market validation cycles) | KPI review, roadmap alignment |
| `[WRITER]` | CMO | Marketing copy |
| `[ARTIST]` | CMO / WRITER | Visual marketing assets |

> **Dormant Agent Rule:** Agents are called only when their domain is active. An agent with no relevant task in the current stage remains dormant and does not consume context.

---

## Fundamental Principles
Development must focus on **consistency**, **traceability**, **documentation**, and **context** over speed. All documents must be filled in **without fail** so that nothing is lost along the way. Documentation is the foundation of execution: without it, progress is invalid.

### Inviolable Rule: Complete File Population
All project files (especially .md files such as PROJECT.md, STATE.md, TASKS.md, TESTS.md, CONTEXT.md, QUESTIONS.md, PLAN.md, ROADMAP.md, ARCHITECTURE.md, README.md, and NEW-INSTRUCTIONS.md) must be filled in **completely and without fail** as required by each stage. Placeholders must be replaced with real content, and no file may remain incomplete or with gaps. Any deviation invalidates the "Done" state and halts progress until corrected.

### Immutability Rule: History with Timestamps
No file may have information deleted; everything must maintain cumulative history with timestamps (YYYY-MM-DD HH:MM) and responsible parties. Add "History" sections to relevant files (e.g.: STATE.md, TESTS.md, CONTEXT.md) to track versions without data loss.

## 1. Pre-Execution Phase (The Stop-Gate)
1.1 **Mandatory Context:** Before any code, the scope, stack, security, and Antigravity environment must be declared in `.md` files.
1.2 **QUESTIONS.md:** Sole clarification channel. If there is ambiguity, register the question and **halt execution**. Nothing starts with pending questions.
1.3 **Immutable History:** Answers in `QUESTIONS.md` must never be deleted; they serve as an architectural decision log.
1.4 **Decisions and Context:** All decisions, hypotheses, contexts, and intentions must be documented before proceeding. Use `QUESTIONS.md` for questions and `CONTEXT.md` for structured context.
1.5 **Direct Order:** Development only begins after an explicit user command and clearing of questions.

## 1.8 Stage Closure Gate
Before starting any new ROADMAP stage, the agent **must mandatorily**:

1. **Verify and present to the user** the closure checklist below, with the real status of each item.
2. **Block advancement** while any item is incomplete.
3. **Wait for explicit confirmation** from the user (verbal OK or `[x]` marking in ROADMAP.md) before continuing.

### Stage Closure Checklist (mandatory)

| # | File | What to Check | Status |
| :- | :--- | :--- | :--- |
| 1 | `TESTS.md` | All stage tests recorded with real results (Passed/Failed), timestamp, and corrective action if applicable. | `[ ]` |
| 2 | `STATE.md` | "What Was Completed" section updated with delivery description, metrics, and timestamp. | `[ ]` |
| 3 | `TASKS.md` | All stage tasks marked as `Done`. | `[ ]` |
| 4 | `CONTEXT.md` | Context summary updated with decisions made in the stage. | `[ ]` |
| 5 | `README.md` | Reflecting the reality of the delivered code (no placeholders relating to the stage). | `[ ]` |
| 6 | `ROADMAP.md` | Previous stage marked as `[x]` **by the user**. | `[ ]` |
| 7 | `PLAN.md` | "Stage Closure Gate" section (§14) filled in and approved. | `[ ]` |
| 8 | `VERSIONS.md` | New entry added if there was a version delivery or significant milestone. | `[ ]` |

> **Rule:** The agent must not start the next stage until the **[HUMAN]** explicitly confirms the previous stage is complete. A simple instruction "start stage X" **is not sufficient** — the agent must first present the checklist above and await `[HUMAN]` validation.
>
> **[HUMAN] Authority:** The human user is the **final approval authority** for all Stage Closure Gates and roadmap advancement. No agent — including the CEO — may authorize stage advancement unilaterally. The human signals completion by sending `[USER_DONE]` in the chat or marking `[x]` in `ROADMAP.md`.

## 1.6 MVP and Feasibility
An MVP (Minimum Viable Product) must be planned at the beginning of development, ensuring project feasibility testing. Even with all other mappings completed (architecture, stack, roadmap), the MVP is mandatory to validate initial hypotheses. Post-MVP adjustments will be controlled exclusively by the NEW-INSTRUCTIONS.md file, which will guide refinements based on feedback and feasibility tests.

## 1.7 File Population
1.6.1 **NEW-INSTRUCTIONS.md:** Must be filled with a clear and structured pattern whenever there are new instructions.
- Title and update date.
- Cycle objective.
- Scope and deliverables.
- Constraints, dependencies, and context.
- Structured markers to ensure readability and traceability.
1.6.2 **QUESTIONS.md:** Each question and answer must have a timestamp and status, ensuring complete history.
- Use headers per entry: `### [YYYY-MM-DD HH:MM:SS] Question`
- Record `Question`, `Context`, `Status`, `Author`, `Answer`, and `Decision / Action`.
- Never delete old entries; always add new entries to maintain an immutable log.
1.6.3 **CONTEXT.md:** After each reading of `NEW-INSTRUCTIONS.md` or new answer in `QUESTIONS.md`, update the context and decisions summary.
1.6.4a **PLAYBOOK.md:** After each new answer in `QUESTIONS.md`, after each new entry in `NEW-INSTRUCTIONS.md`, and after the project definition (`PROJECT.md`) is filled in, the agent must **review PLAYBOOK.md** and append new entries (with timestamp) if any preference, recurring decision, or working pattern has been revealed or confirmed. This step is mandatory and happens *before* the Atomic Commit.
1.6.4 **Recommended markers:**
- `#` for main title
- `##` for objective, scope, constraints, context, and deliverables sections
- `###` for individual question entries
- Lists with `-` for items and decisions
- Standardized fields to facilitate automated reading and traceability

## 2. Planning and Traceability
2.1 **PLAN.md:** Each roadmap stage must be described completely and operationally.
- Stage description and clear objective.
- Detailed and measurable acceptance criteria.
- Applicable business rules.
- Tests required to validate functionality and business rules.
- Mandatory security validations for production.
- Validation steps the developer must execute.
- Dependencies and preconditions.
2.2 **TASKS.md:** Registry of all atomic project activities. Each task must have a status (To Do, Doing, Done).
2.3 **TESTS.md:** Test results performed by the AI agent must be saved in `tests.md`.
- Record each test with step-by-step, environment, expected result, obtained result, and status.
- Document corrective actions when a test fails.
- Use `tests.md` as a concrete guide for the developer to replicate and validate manual or automated tests.

## 3. TDD Development Cycle (Strict Mode)
No line of production code without a prior test.
1. **RED:** Create a unit/integration test that fails.
2. **GREEN:** Minimum code to pass the test.
3. **REFACTOR:** Apply SOLID and Clean Code without breaking the test.
4. **SECURITY:** Validate sanitization, injection protection, and JWT compliance.

## 4. Hexagonal Architecture & Isolation
- **DOMAIN:** Pure logic. Importing frameworks (Spring, Pandas, etc.) is forbidden.
- **APPLICATION (Ports):** Interfaces that define the input and output contract.
- **INFRASTRUCTURE (Adapters):** Volatile technical implementations (DB, APIs, UI).

## 4.1 Naming and Directory Structure
Before creating directories for backend, frontend, or services, validate the structure and names.
- If the solution is Web, the frontend should be named `<abbreviated-name>-web`.
- If the solution is a native/hybrid app, apply appropriate naming logic to the project.
- If there is more than one microservice, each service should be created as `<name>-service`.
- If it is a monolith, the directory can be created with just the abbreviated project name.

## 5. Code Standards & SOLID
- **S:** Classes with single responsibility.
- **O:** Extensibility via interfaces (new models/indicators).
- **L/I/D:** Strict dependency injection and specific interfaces.
- **Clean Code:** Semantic names and readability as priority.

## 6. GSD File Workflow
For each delivery, the order is mandatory and **all documents must be filled without fail** to ensure consistency and complete traceability:
1. **Check:** Read `NEW-INSTRUCTIONS.md` and validate `PLAN.md`.
1.a **Roadmap:** Start or update `ROADMAP.md` with each reading of `NEW-INSTRUCTIONS.md`, to ensure the plan reflects the most recent context.
2. **Test:** Create/update the test file.
3. **Impl:** Create/update the implementation code.
4. **Migration:** If there is a schema change, create/update the SQL migration file.
5. **Docs:** Update `README.md`, `ROADMAP.md`, `STATE.md`, `TASKS.md`, `CONTEXT.md`, and `tests.md`. **Mandatory:** No document may remain incomplete or with failures.
5.a **PLAYBOOK:** Review `PLAYBOOK.md` and add new entries (with timestamp) if the delivery revealed developer preferences, decisions, or patterns relevant to future projects.
6. **Atomic Commit:** Simulate/perform the activity commit.

## 7. Security & Persistence (Infra)
7.1 **API Lockdown:** Java routes are born `deny-all`. Release via JWT.
7.2 **Zero Leak:** Logging keys, payloads, or transactions in plain text is forbidden.
7.3 **Sanitization:** Adapters must clean data before sending it to the Domain.
7.4 **Docker & DB:** Before starting, verify Docker is active. Bring up the database and apply **Migrations** as needed.
7.5 **Credential Security:** No password, API key, or sensitive credential should be in code or documentation files. The `.env` file must never be committed, nor any information that compromises application security. Everything needed in the environment must be clearly informed in [ENV_SETUP.md or specific section], with a defined location for secure storage (e.g.: external vaults). Register instructions in QUESTIONS.md if there are questions about configuration.
7.6 **Mandatory Logging:** Every action performed by the solution must have logs for activity and/or error traceability. Configure structured logs (e.g.: INFO, WARN, ERROR levels) without exposing sensitive data. Record logging validations in TESTS.md.

## 8. Documentation as Code
The project is guided by the following files — note their locations:

| File | Location | Purpose |
| :--- | :--- | :--- |
| `README.md` | Root | Project overview and quick-start |
| `PLAYBOOK.md` | Root | Developer preferences (project-agnostic) |
| `NEW-INSTRUCTIONS.md` | **Root** | Current cycle instructions from the user |
| `QUESTIONS.md` | **Root** | Sole clarification channel; immutable decision log |
| `ONBOARDING.md` | `DOC/` | New project setup guide (run once per project) |
| `GSD-RULES.md` | `DOC/` | Inviolable execution rules |
| `ARCHITECTURE.md` | `DOC/` | Hexagonal architecture, stack, components |
| `PROJECT.md` | `DOC/` | Vision, objectives, team, and KPIs |
| `PLAN.md` | `DOC/` | Stage-level acceptance criteria and business rules |
| `ROADMAP.md` | `DOC/` | High-level stage tracking |
| `STATE.md` | `DOC/` | Current state, deliveries, and blockers |
| `TASKS.md` | `DOC/` | Atomic task registry with statuses |
| `TESTS.md` | `DOC/` | Test results, steps, and validations |
| `CONTEXT.md` | `DOC/` | Decisions and context summary |
| `ENV_SETUP.md` | `DOC/` | Environment variables and secure credential guidance |
| `DESIGN.md` | `DOC/` | Visual identity and design tokens |
| `VERSIONS.md` | `DOC/` | Release history and milestone log |

## 9. Installed Skills & Environment Tools

The Antigravity environment has **skills** (specialized instruction packages) and **external tools** that extend the AI agent's capabilities. The agent **must** consult the available skills before starting any development phase and record in `QUESTIONS.md` which ones can be used and in which context.

## 9.1. LLM Tiering Strategy
| Tier | Model | Purpose & Use Cases |
| :--- | :--- | :--- | 
| Tier 1: Efficiency | {{TIER_1_EFFICIENCY_MODEL}} | Atomic & Repetitive Tasks: Log generation, JSON/YAML formatting in .agent_handoff/, syntax linting, and rapid documentation research. |
| Tier 2: Development | {{TIER_2_DEVELOPMENT_MODEL}} | Standard Implementation: TDD cycles (Red-Green-Refactor), unit test creation in `TESTS.md`, and SQL/infrastructure optimization. |
| Tier 3: Expert | {{TIER_3_EXPERT_MODEL}} | Architecture & Decision Making: Resolving critical ambiguities in `QUESTIONS.md`, enforcing Hexagonal Architecture patterns, and strategic Roadmap planning. |


### 9.1 Antigravity Environment Skills
Skills are thematic extensions installed in `<appDataDir>\skills\`. Before implementing any component, the agent must check if there is a relevant skill and follow it.

**Mandatory skills by domain (examples):**

| Domain | Recommended Skills | When to Use |
| :--- | :--- | :--- |
| **Architecture** | `architect-review`, `senior-architect`, `ddd-context-mapping` | Design review and context mapping |
| **Backend Java/Spring** | `api-patterns`, `backend-architect`, `dotnet-backend` (reference) | API design, service patterns |
| **Python / ML** | `scikit-learn`, `ml-pipeline-workflow`, `python-pro` | ML pipelines, modeling, optimization |
| **Frontend Angular** | `react-patterns` (component reference), `frontend-design` | UI design, componentization |
| **Database** | `database`, `database-design`, `postgres-best-practices` | Schema, migrations, SQL optimization |
| **Security** | `security-auditor`, `differential-review`, `gdpr-data-handling` | Audits, secure code review |
| **TDD / Testing** | `tdd-workflow`, `tdd-workflows-tdd-cycle`, `webapp-testing` | Red-Green-Refactor cycle, E2E tests |
| **DevOps / CI/CD** | `github-actions-templates`, `docker-expert`, `terraform-skill` | Pipelines, containers, IaC |
| **Documentation** | `documentation`, `wiki-architect`, `mermaid-expert` | Doc generation, diagrams |
| **Debugging** | `debugger`, `systematic-debugging`, `error-diagnostics-smart-debug` | Bug and failure diagnosis |
| **Observability** | `grafana-dashboards`, `prometheus-configuration`, `distributed-tracing` | Monitoring, metrics, tracing |
| **Git / PR** | `github`, `pr-writer`, `git-pr-workflows-pr-enhance` | Git workflows, PR creation |

### 9.2 Suggested External Tools
In addition to skills, the following external tools are recommended to complement the GSD workflow:

- **Linting & Formatting:** ESLint, Prettier (frontend), Checkstyle/SpotBugs (Java), Black/Ruff (Python).
- **Static Application Security Testing (SAST):** SonarQube, Semgrep, Snyk.
- **Performance Testing:** k6, JMeter, Locust.
- **Dependency Monitoring:** Dependabot, Renovate, OWASP Dependency-Check.
- **API Documentation:** Swagger/OpenAPI (Spring), Sphinx (Python).
- **Diagrams:** Mermaid (integrated via `mermaid-expert` skill), draw.io, PlantUML.
- **Secrets Management:** HashiCorp Vault, AWS Secrets Manager, Azure Key Vault (see ENV_SETUP.md).
- **Containerization & Orchestration:** Docker Compose (dev), Kubernetes/Helm (production).
- **Messaging & Cache:** Apache Kafka, Redis (as per ARCHITECTURE.md).

### 9.3 Usage Rules
9.3.1 **Mandatory Consultation:** Before starting any ROADMAP phase, the agent must list applicable skills and record them in `CONTEXT.md`.
9.3.2 **Skill Priority:** If a skill exists for the domain, it takes precedence over generic knowledge.
9.3.3 **Developer Suggestions:** The agent must suggest relevant external tools in `PLAN.md` (Dependencies section) when identifying needs not covered by skills.
9.3.4 **Updates:** When new skills are installed or removed, update this section and `CONTEXT.md`.
9.3.5 **Usage Documentation:** Record in `TESTS.md` validations that used specific skills or tools.
9.3.6 **CTO Ownership — Antigravity Skill Mapping:** The **[CTO] Agent is the sole responsible agent** for analyzing the installed skills catalog (this section) before each roadmap stage. The CTO must:
  - Identify which skills apply to the current stage's scope and tech stack.
  - Communicate the selected skills to the relevant technical agents ([DEV_BACKEND], [DEV_FRONTEND], [DBA], [DS/ML], [DEVOPS]) via `.agent_handoff/`.
  - Record the skill mapping decision in `CONTEXT.md` before any implementation begins.
  - Update this section if new skills are installed or deprecated.

---
**AGENT SIGNATURE:** Operating in GSD Mode - Execution on demand, quality by design.
