# ENV_SETUP.md - PF_ai

All credentials must be stored securely through environment variables or vaults. Never commit real secrets.

## 1. Prerequisites
- **Operating System:** Windows with WSL 2 recommended, Linux, or macOS.
- **Tools:** Go, Node.js, Docker, Docker Compose, Git.
- **Database:** PostgreSQL 16 through Docker Compose.
- **Hardware:** Enough RAM for Docker plus optional local model runtime. Local models may require elevated memory depending on model size.

## 2. Stack Installation
- **Backend (Go):** Install Go and validate with `go version`.
- **Frontend (TypeScript Web):** Install Node.js and validate with `node --version` and `npm --version`.
- **Database (PostgreSQL):** Run through Docker Compose; no local native install required for MVP.
- **Container (Docker):** Install Docker Desktop or Docker Engine with Compose support.
- **Local Model Runtime:** Configure a Docker service for the selected local model provider in a later implementation phase.
- **External Providers:** Configure provider API keys through environment variables only.

## 3. Environment Variables
Configure these variables locally. Do not commit real values.

| Variable | Purpose |
| :--- | :--- |
| `PF_AI_ENV` | `local`, `test`, or `production` |
| `PF_AI_HTTP_PORT` | Backend HTTP port |
| `PF_AI_DB_URL` | PostgreSQL connection string |
| `PF_AI_DB_USER` | PostgreSQL user |
| `PF_AI_DB_PASSWORD` | PostgreSQL password |
| `PF_AI_AUTH_SECRET` | JWT or signed-session secret |
| `PF_AI_HANDOFF_DIR` | Path to `.agent_handoff/` |
| `PF_AI_DOC_ROOT` | Path to project root for Markdown memory |
| `PF_AI_PROVIDER_MODE` | `api`, `local`, or `hybrid` |
| `PF_AI_LOCAL_MODEL_URL` | Local model endpoint when enabled |
| `PF_AI_API_PROVIDER` | Selected external provider name |
| `PF_AI_API_KEY` | External provider key, stored securely |
| `LOG_LEVEL` | `INFO`, `WARN`, or `ERROR` |

## 3.1 AI Agents Configuration
- `CEO_AGENT_MODEL`: Tier 3 model for orchestration.
- `BA_AGENT_MODEL`: Tier 2 or Tier 3 model for requirements.
- `CTO_AGENT_MODEL`: Tier 3 model for architecture.
- `DEV_BACK_AGENT_MODEL`: Tier 2 model for Go implementation.
- `DEV_FRONT_AGENT_MODEL`: Tier 2 model for TypeScript implementation.
- `SECURITY_AGENT_MODEL`: Tier 3 model for threat modeling and audits.
- `QA_REVIEWER_AGENT_MODEL`: Tier 1 or Tier 2 model for repeatable verification.

## 4. Keys and Credentials
- Use `.env` locally only if it is git-ignored.
- Production secrets must use a secret manager.
- Never write provider keys, local auth secrets, database passwords, or tokens into Markdown files.
- If a required secret is missing, register the blocker in `QUESTIONS.md`.

## 5. Local Environment Configuration
1. Clone the repository.
2. Configure `.env` locally using the variable list above.
3. Start PostgreSQL with Docker Compose when the compose file exists.
4. Start the backend Go service when implemented.
5. Start the TypeScript web frontend when implemented.
6. Enable the local model container only after Phase 2 scope is approved.

## 6. Production Configuration
- Deploy as container images after CI/CD is defined.
- Store secrets in the selected production platform.
- Add monitoring with Prometheus/Grafana after MVP stabilization.

## 7. Validations
- Validate Docker is active.
- Validate PostgreSQL connection.
- Validate protected routes reject unauthenticated requests.
- Validate logs contain no secrets.
- Validate file handoffs are created in `.agent_handoff/` without credentials.
- Record all results in `DOC/TESTS.md`.

## 8. Troubleshooting
- Docker unavailable: verify Docker Desktop or daemon is running.
- PostgreSQL unavailable: verify Compose service health and credentials.
- Local model unavailable: verify container status, memory limits, and endpoint URL.
- Provider API failure: verify environment variable presence without logging secret values.

## 9. Obsidian - Project Maintenance Wiki
The project uses `wiki/` and `raw/` directories for local documentation. See `DOC/WIKI.md`.

## 10. Version History (Immutable)
- **[2026-05-05 15:06] - [CEO]:** Environment setup filled during onboarding for Go, TypeScript web, PostgreSQL, Docker, API providers, and local model runtime.
