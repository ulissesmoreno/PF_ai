# PF_ai

## Portugues

### Descricao
PF_ai e uma camada operacional agentica para gerenciar modelos, memoria, agentes e handoffs GSD. O projeto cria uma interface humano-agentes para coordenar execucao, consultar estado, manter rastreabilidade e garantir que o boilerplate funcione com governanca tecnica.

### Objetivo
Centralizar a operacao de agentes de IA, permitindo execucao via APIs externas ou modelos locais em Docker, com memoria baseada em arquivos Markdown e handoffs estruturados.

### Publico-alvo
Ulisses como operador principal e, futuramente, usuarios tecnicos que precisam de um workspace agentico local-first, auditavel e orientado por fases.

### Instalacao
1. Prerequisitos: Go, Node.js, Docker, Docker Compose e PostgreSQL via Compose.
2. Clone o repositorio e entre na pasta do projeto.
3. Configure variaveis de ambiente conforme `DOC/ENV_SETUP.md`.
4. Copie `.env.example` para `.env` local e substitua placeholders fora do Git.
5. Suba PostgreSQL com `docker compose up -d postgres`.
6. Execute backend com `go run ./cmd/pf-ai-service` quando Go estiver instalado.
7. Execute frontend com `cd pf-ai-web` e `npm start`.
8. Acesse a UI local em `http://127.0.0.1:5173`.

### Uso
- Frontend web: interface humano-agentes para painel, agentes, modelos, memoria e handoffs.
- Backend Go: API local para orquestracao, leitura dos arquivos GSD e integracao com provedores.
- Handoffs MVP: arquivos JSON em `.agent_handoff/`.
- Handoffs finais: MCP planejado em fase posterior.
- Persistencia MVP: schema PostgreSQL e adapters Go compativeis com `database/sql` para agentes, provedores, sessoes e auditoria.
- Phase 1 MVP: o dashboard permite salvar providers, salvar agentes, criar agente sem provider, ler memoria GSD allowlisted e criar handoff. Quando a API local nao esta disponivel, saves de provider/agente/handoff usam fallback em `localStorage` para validacao humana no navegador.
- Todo novo projeto inicia com os agentes obrigatorios CEO, CTO e BA; demais agentes sao opcionais.

### Execucao local do MVP
1. Frontend: `cd pf-ai-web` e `npm start`; abrir `http://127.0.0.1:5173`.
2. Backend interativo: definir `PF_AI_AUTH_SECRET` e `PF_AI_HTTP_PORT=8081`, entao executar `go run ./cmd/pf-ai-service`.
3. Dashboard: informar API base `http://127.0.0.1:8081`, token local, e usar `Connect`.
4. Sem backend interativo: a UI continua testavel via fallback local para cadastro de providers, agentes e handoffs.

### Runtime local e roteamento hibrido
- Endpoint protegido: `GET /api/provider-health?id=<provider_id>`.
- Providers `local` e `hybrid` aceitam endpoint local apenas em loopback (`localhost`, `127.0.0.1`, `::1`).
- Provider `hybrid` tenta rota local quando saudavel e faz fallback para API com motivo auditavel quando local estiver indisponivel.
- Startup local recomendado: definir `PF_AI_AUTH_SECRET` no shell e executar `scripts\run-backend-local.ps1`.

### Orquestracao avancada
- Planning Board: cards por status, prioridade, agente responsavel, fase e task reference.
- Project Registration: cadastro do projeto com payload de onboarding; campos com aparencia de segredo sao rejeitados.
- MCP Envelope: handoffs podem ser exportados em envelope `mcp-compatible` para migracao de transporte.

### Desenvolvimento
- Arquitetura: `DOC/ARCHITECTURE.md`
- Regras GSD: `DOC/GSD-RULES.md`
- Setup: `DOC/ENV_SETUP.md`
- Roadmap: `DOC/ROADMAP.md`
- Duvidas: `QUESTIONS.md`

### Testes
Resultados obrigatorios devem ser registrados em `DOC/TESTS.md`.
- Backend Go: `go test ./...`
- Frontend shell: `node --check pf-ai-web/src/app.js` e `node --check pf-ai-web/server.mjs`

Nota local: se o cache Go em `AppData` estiver bloqueado, defina `GOCACHE` para `.gocache` dentro do workspace antes de rodar testes.

## English

### Description
PF_ai is an agentic operating layer for managing AI models, memory, agents, and GSD handoffs. It provides a human-agent interface for coordinating execution, inspecting state, preserving traceability, and keeping the boilerplate technically governed.

### Objective
Centralize AI-agent operations with support for external API providers and local Docker-based models, using Markdown memory and structured handoff files.

### Target Audience
Ulisses as the primary operator, with future support for technical users who need a local-first, auditable, phase-driven agentic workspace.

### Installation
1. Prerequisites: Go, Node.js, Docker, Docker Compose, and PostgreSQL through Compose.
2. Clone the repository and enter the project directory.
3. Configure environment variables according to `DOC/ENV_SETUP.md`.
4. Copy `.env.example` to a local `.env` and replace placeholders outside Git.
5. Start PostgreSQL with `docker compose up -d postgres`.
6. Run the backend with `go run ./cmd/pf-ai-service` when Go is installed.
7. Run the frontend with `cd pf-ai-web` and `npm start`.
8. Open the local UI at `http://127.0.0.1:5173`.

Current local note: if the Go cache under `AppData` is blocked, set `GOCACHE` to `.gocache` inside the workspace before running tests.

### Usage
- Web frontend: human-agent interface for dashboard, agents, models, memory, and handoffs.
- Go backend: local orchestration API, GSD file reading, and provider integration.
- MVP handoffs: JSON files in `.agent_handoff/`.
- Final handoffs: MCP planned for a later stage.
- MVP persistence: PostgreSQL schema and `database/sql`-compatible Go adapters for agents, providers, sessions, and audit events.
- Phase 1 MVP: the dashboard can save providers, save agents, create no-provider agents, read allowlisted GSD memory, and create handoffs. When the local API is unavailable, provider/agent/handoff saves use `localStorage` fallback for human browser validation.
- Every new project starts with mandatory CEO, CTO, and BA agents; all other agents are optional.

### Local MVP Runtime
1. Frontend: `cd pf-ai-web` and `npm start`; open `http://127.0.0.1:5173`.
2. Interactive backend: set `PF_AI_AUTH_SECRET` and `PF_AI_HTTP_PORT=8081`, then run `go run ./cmd/pf-ai-service`.
3. Dashboard: use API base `http://127.0.0.1:8081`, local token, and `Connect`.
4. Without interactive backend: the UI remains testable through local fallback for provider, agent, and handoff saves.

### Local Runtime And Hybrid Routing
- Protected endpoint: `GET /api/provider-health?id=<provider_id>`.
- `local` and `hybrid` providers accept local endpoints only on loopback (`localhost`, `127.0.0.1`, `::1`).
- A `hybrid` provider uses the local route when healthy and falls back to API with an auditable reason when local is unavailable.
- Recommended local startup: set `PF_AI_AUTH_SECRET` in the shell and run `scripts\run-backend-local.ps1`.

### Advanced Orchestration
- Planning Board: cards by status, priority, responsible agent, phase, and task reference.
- Project Registration: project registration with onboarding payload; secret-looking fields are rejected.
- MCP Envelope: handoffs can be exported as an `mcp-compatible` envelope for transport migration.

### Development
- Architecture: `DOC/ARCHITECTURE.md`
- GSD rules: `DOC/GSD-RULES.md`
- Setup: `DOC/ENV_SETUP.md`
- Roadmap: `DOC/ROADMAP.md`
- Questions: `QUESTIONS.md`

### Testing
Mandatory results must be recorded in `DOC/TESTS.md`.
- Go backend: `go test ./...`
- Frontend shell: `node --check pf-ai-web/src/app.js` and `node --check pf-ai-web/server.mjs`

## GSD Ignition

Every cycle starts only after `[USER_DONE]`. CEO acts first, then delegates through the GSD chain of command. Chat stays compact; details live in Markdown files.

## License
To be defined by [HUMAN:Ulisses].
