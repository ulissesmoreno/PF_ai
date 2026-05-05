# PROJECT.md - PF_ai

This file defines the product vision and must be filled completely and without fail to ensure project context and consistency.

## 1. Project Overview
- **Project Name:** PF_ai
- **Description:** PF_ai is an agentic operating layer for managing AI models, persistent memory, agent roles, and GSD handoffs. It provides a human-agent interface so the user can coordinate agents, inspect project state, and keep boilerplate execution consistent across local and API-based model providers.
- **Domain:** Agent orchestration, AI developer tooling, project governance, and local-first automation.
- **Target Audience:** Ulisses as the primary operator, plus future technical users who need a governed AI-agent workspace with traceable memory, models, and phase execution.

## 2. Objectives and Scope
- **Main Objectives:**
  - Manage agent definitions, responsibilities, seniority, and routing.
  - Manage AI model providers that can run through external APIs or local containers.
  - Maintain project memory through GSD documents and structured handoff files.
  - Provide a human-agent interface inspired by agent-company orchestration patterns.
  - Preserve token economy by using compact chat responses and file-based detail.
- **Initial Scope:**
  - Web-first interface for observing and operating the GSD workspace.
  - Go backend for orchestration APIs, file adapters, provider adapters, and local process control.
  - PostgreSQL for durable application state.
  - Docker Compose for PostgreSQL and optional local model runtime.
  - MVP handoffs through files in `.agent_handoff/`.
- **Limitations:**
  - MCP-based handoffs are planned for the final stage, not the MVP.
  - Full autonomous implementation loops are outside the first MVP unless explicitly approved.
  - Production multi-tenant security is not part of the first MVP.

## 3. Technology Stack
- **Backend:** Go - Primary orchestration backend, REST API, local process control, provider adapters, and GSD file adapters.
- **Frontend:** TypeScript web application - Human-agent interface, dashboard, agent registry, memory view, and handoff controls.
- **Database:** PostgreSQL 16 - Durable state for agents, models, providers, sessions, and audit metadata.
- **Local Model Runtime:** Docker Compose service - Optional local model container for offline or private execution.
- **External Model Access:** Provider adapters - API-based access for selected LLM providers; credentials via environment variables only.
- **Security:** JWT or local signed session tokens for protected routes; deny-by-default API policy.
- **Container:** Docker / Docker Compose - Local development services and local model runtime.
- **Handoffs:** File-based JSON in `.agent_handoff/` for MVP; MCP transport planned for the final stage.
- **Documentation Memory:** Markdown files under `DOC/`, root operational files, and `wiki/`.
- **Inspirations:** Paperclip AI for agent-company orchestration, Caveman/Cavemem/Cavekit for token economy and persistent memory, Antigravity skill catalogs for runtime capability mapping, Open Claude as inspiration for agent-manager UX patterns.

## 4. Roadmap and Stages
- **Reference:** Consult `DOC/ROADMAP.md` for details.
- **Initial Stages:**
  - Phase 0: Onboarding and foundation.
  - Phase 1: MVP planning and implementation.
  - Phase 2: Local model runtime and provider routing hardening.
  - Phase 3: MCP handoffs and advanced orchestration.

## 5. Team and Responsibilities
- **Developers:** Ulisses as [HUMAN:Ulisses], with AI agents executing under the GSD chain of command.
- **Reviewer:** [HUMAN:Ulisses] is the final approver for Stage Closure Gates and roadmap advancement.

### 5.1 Human Role - [HUMAN:Ulisses]
The `[HUMAN:Ulisses]` role represents the final authority in the GSD ecosystem.

**[HUMAN:Ulisses] Responsibilities:**
- Provide instructions via chat and `NEW-INSTRUCTIONS.md` at the start of each work cycle.
- Send the `[USER_DONE]` signal in chat to trigger the CEO Orchestration Cycle.
- Provide the sole and exclusive approval for all Stage Closure Gates.
- Mark stages as `[x]` in `ROADMAP.md` to confirm completion.
- Review and approve `PLAN.md`, `TESTS.md`, and `STATE.md` summaries before sign-off.

### 5.2 AI Agent Roles
- **Active Agents by Default:** CEO, CTO, BA, SECURITY, TECH_LEAD, QA, CODE_REVIEWER, DEV_BACKEND, DEV_FRONTEND, DBA, DEVOPS.
- **Available on Demand (Dormant):** DS_ML, CMO, WRITER, ARTIST, DOCUMENTATION, DATA_ENGINEER.
- **Dormant Rule:** Agents are activated only when their domain is relevant to the current stage.
- **Creating New Agents:** New role files under `AGENTS/` must define rules, document access limits, and responsibilities.
- **Late Addition:** If a new role/agent is needed after the project has started, use `NEW-INSTRUCTIONS.md` to request and define it.

## 6. Risks and Mitigations
- **Technical Risks:**
  - Local model containers may be heavy for the development machine.
  - File-based memory can drift from database state if adapters are not designed carefully.
  - API provider credentials create leakage risk if logs are not sanitized.
  - MCP handoff migration can become complex if MVP handoff schemas are unstable.
- **Mitigations:**
  - Keep the MVP provider abstraction small and testable.
  - Treat Markdown and `.agent_handoff/` as first-class ports/adapters.
  - Use deny-by-default API routes and sanitized structured logs.
  - Freeze handoff JSON schemas early and version them.
  - Store all credentials in environment variables or secure vaults only.

## 7. Success Metrics
- **KPIs:**
  - Project onboarding has no remaining placeholders in required files.
  - MVP can register agents and models through the UI/API.
  - MVP can read key GSD files and create structured file handoffs.
  - MVP supports at least one API-based provider and one local-provider configuration path.
  - All protected routes reject unauthenticated access.
  - All tests recorded in `DOC/TESTS.md` with real results before stage closure.

## 8. Version History (Immutable)
- **[2026-05-05 15:06] - [CEO]:** Onboarding auto-fill completed from [HUMAN] answers in `QUESTIONS.md`.

---
*Reminder for AI: All questions must be registered in `QUESTIONS.md`.*
