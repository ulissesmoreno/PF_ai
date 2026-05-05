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
