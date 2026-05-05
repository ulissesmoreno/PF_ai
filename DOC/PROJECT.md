# PROJECT.md -

This file defines the product vision and must be filled **completely and without fail** to ensure project context and consistency.

## 1. Project Overview
- **Project Name:** [Fill with specific name]
- **Description:** [Detailed system description]
- **Domain:** [E.g.: Financial analysis, etc.]
- **Target Audience:** [E.g.: Traders, analysts, etc.]

## 2. Objectives and Scope
- **Main Objectives:** [List of goals]
- **Initial Scope:** [What will be included]
- **Limitations:** [What will not be covered]

## 3. Technology Stack (Based on Boilerplate)
- **Backend:** Java 21 / Spring Boot 3 - [Specific description]
- **Intelligence:** Python 3.12 / Scikit-Learn - [Specific description]
- **Frontend:** Angular 18 / TypeScript - [Specific description]
- **Database:** PostgreSQL 16 - [Specific description]
- **Security:** Spring Security + JWT - [Specific description]
- **Container:** Docker / Compose - [Specific description]
- **Others:** [Add as needed, e.g.: Redis, Kafka]

## 4. Roadmap and Stages
- **Reference:** Consult ROADMAP.md for details.
- **Initial Stages:** [List the first phases]

## 5. Team and Responsibilities
- **Developers:** [Names and roles]
- **Reviewer:** [Who approves stages]

### 5.1 Human Role — [HUMAN]
The `[HUMAN]` role represents the **final authority** in the GSD ecosystem. This is the real person (or persons) responsible for the project.

**[HUMAN] Responsibilities:**
- Provide instructions via `NEW-INSTRUCTIONS.md` at the start of each work cycle.
- Send the `[USER_DONE]` signal in chat to trigger the CEO Orchestration Cycle.
- Provide the **sole and exclusive approval** for all Stage Closure Gates — no agent may authorize stage advancement unilaterally.
- Mark stages as `[x]` in `ROADMAP.md` to confirm completion.
- Review and approve `PLAN.md`, `TESTS.md`, and `STATE.md` summaries before sign-off.

**Filling this section:** When a project starts, replace `[HUMAN]` with the actual person's name or identifier:
- `[HUMAN:Ulisses]` — for single-person projects
- `[HUMAN:CEO_Name]`, `[HUMAN:PO_Name]` — for teams with multiple approvers

### 5.2 AI Agent Roles
- **Active Agents by Default:** CEO, CTO, BA, QA, DEV_FRONTEND, DEV_BACKEND, DBA.
- **Available on Demand (Dormant):** DS_ML, SECURITY, CODE_REVIEWER, DEVOPS, CMO, WRITER, ARTIST.
- **Dormant Rule:** Agents are activated only when their domain is relevant to the current stage. A dormant agent does not participate and does not consume context.
- **Creating New Agents:** If it is necessary to create a new agent (e.g., seniority levels, niche specialists), you must create a new `.md` file in the `AGENTS/` directory with their specific rules, document access limits, and responsibilities.
- **Late Addition:** If a new role/agent is needed after the project has started, use `NEW-INSTRUCTIONS.md` to request and define its creation.

## 6. Risks and Mitigations
- **Technical Risks:** [E.g.: Hardware dependencies]
- **Mitigations:** [E.g.: Optimizations]

## 7. Success Metrics
- **KPIs:** [E.g.: Test coverage, performance]

---
*Reminder for AI: All questions must be registered in the QUESTIONS.md file.*
