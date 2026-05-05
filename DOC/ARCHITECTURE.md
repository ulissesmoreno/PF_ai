# ARCHITECTURE.md - PF_ai

## 1. System Overview
PF_ai uses Hexagonal Architecture to keep orchestration rules independent from web, database, file-system, Docker, and model-provider details.

## 2. Development Pillars
- **GSD:** Delivery with traceability and documented decisions.
- **TDD:** Red, Green, Refactor, Security before delivery.
- **Token Economy:** Compact chat, detailed Markdown memory, scoped agents.
- **Local-first:** The project must operate with local files and local model runtime when configured.
- **Provider-flexible:** Agents may run through API providers or local containers.

## 3. Layer Organization

### Domain
Located in `src/domain`.
- Agent definitions and role rules.
- Model/provider capabilities.
- Memory and handoff policies.
- Stage and phase state concepts.

### Application
Located in `src/application`.
- Agent registry use cases.
- Model provider registry use cases.
- GSD memory read/write use cases.
- Handoff creation and validation use cases.
- Runtime routing use cases for API or local model execution.

### Infrastructure
Located in `src/infrastructure`.
- HTTP controllers.
- PostgreSQL repositories.
- Markdown file adapters.
- `.agent_handoff/` JSON adapters.
- Docker/local model runtime adapters.
- External LLM API adapters.
- Authentication and structured logging adapters.

## 4. Technology Stack

| Component | Technology | Responsibility | Agent |
| :--- | :--- | :--- | :--- |
| Backend | Go | Orchestration API, ports/adapters, provider routing, file integration | `[DEV_BACKEND]` |
| Frontend Web | TypeScript | Human-agent dashboard and control surface | `[DEV_FRONTEND]` |
| Database | PostgreSQL 16 | Durable state for agents, models, providers, sessions, and audit metadata | `[DBA]` |
| Local Runtime | Docker / Docker Compose | PostgreSQL and optional local model container | `[DEVOPS]` |
| External Providers | API adapters | Remote LLM execution where configured | `[DEV_BACKEND]` |
| Handoffs MVP | JSON files in `.agent_handoff/` | Structured agent communication | `[TECH_LEAD]` |
| Handoffs Final | MCP | Future transport for agent communication | `[CTO]` |
| Security | JWT or signed local sessions | Deny-by-default API protection | `[SECURITY]` |
| Documentation Memory | Markdown | GSD state, context, tests, roadmap, decisions | `[CEO]` / `[DOCUMENTATION]` |

## 4.1 Naming
- Backend service: `pf-ai-service`
- Web frontend: `pf-ai-web`
- Local model service: `pf-ai-model-runtime`

## 4.2 Agent Communication Schemas
MVP handoffs are structured JSON files under `.agent_handoff/`. MCP is planned for the final stage after file schemas stabilize.

### General Handoff Header
```json
{
  "header": {
    "timestamp": "YYYY-MM-DD HH:MM",
    "sender": "[AGENT_TAG]",
    "recipient": "[AGENT_TAG]",
    "task_ref": "TASK-XXX",
    "intent": "PHASE_KICKOFF | CLARIFICATION_REQUEST | CLARIFICATION_RESPONSE | REVIEW_RESULT | STAGE_APPROVAL | BUGFIX_REQUEST | RETROSPECTIVE"
  }
}
```

## 5. Document Ownership

| File | Owner | Scope |
| :--- | :--- | :--- |
| `STATE.md` | Technical agents | Execution memory |
| `CONTEXT.md` | Management agents | Decisions after onboarding and phase intake |
| `PLAYBOOK.md` | CEO | [HUMAN] work preferences |
| `RETROSPECTIVE.md` | CEO | Phase learnings |
| `VERSIONS.md` | Phase-closing agent | Release history |
| `TESTS.md` | QA + SECURITY + Technical | Test results |
| `DESIGN.md` | UX_RESEARCHER + DEV_FRONTEND | Visual identity |
| `wiki/` | DOCUMENTATION | Project knowledge base |
| `QUESTIONS.md` | CEO | Only direct questions from agents to [HUMAN] |

## 6. Security Architecture
- All API routes are born deny-by-default.
- Provider API keys must be environment variables or external secrets.
- Logs must be structured and sanitized.
- Handoff files must not contain secrets or raw provider credentials.
- Local model runtime must be isolated through Docker Compose.

## 7. Definition of Done
1. Tests are created before production code.
2. Security validations are recorded.
3. Documentation files reflect the delivered state.
4. Handoff schemas stay versioned and validated.
5. [HUMAN:Ulisses] approves Stage Closure Gate before roadmap advancement.

## 8. Version History (Immutable)
- **[2026-05-05 15:06] - [CEO]:** Architecture filled during onboarding. Stack set to Go, TypeScript web, PostgreSQL, Docker, API/local model providers, file handoffs first, MCP later.
- **[2026-05-05 15:23] - [CEO]:** Document ownership clarified: `QUESTIONS.md` for agent-to-human questions only; `CONTEXT.md` for architectural decisions; `STATE.md` for technical memory; `VERSIONS.md` at phase closure from Phase 1 onward.
