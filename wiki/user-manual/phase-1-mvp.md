# Phase 1 MVP User Manual

## Updated
- **Timestamp:** 2026-05-05 20:06
- **Responsible:** [DOCUMENTATION]
- **Scope:** Dashboard operation for Phase 1 MVP.

## Purpose
Use the PF_ai dashboard to manage providers, agents, GSD memory reads, and file handoffs during the Phase 1 MVP.

## Open The Dashboard
1. Start the frontend with `cd pf-ai-web` and `npm start`.
2. Open `http://127.0.0.1:5173`.
3. Optional API-backed mode: start the backend interactively with `PF_AI_AUTH_SECRET` and `PF_AI_HTTP_PORT=8081`, then connect the dashboard to `http://127.0.0.1:8081`.

## Providers
1. Open the provider form.
2. Choose provider mode:
   - `online` for API providers. Use a secret reference, never a raw key.
   - `local` for local runtime paths.
   - `hybrid` for mixed routing.
3. Click `Save Provider`.
4. Saved providers appear in the providers list and in the agent provider selector.

## Agents
1. Open the agent form.
2. Select a role from the list or choose `Outro` for a custom role.
3. Enter the agent name; the ID is generated automatically.
4. Choose a provider, or select `No provider` for embedded-model environments such as Codex or Claude Code.
5. Click `Save Agent`.

## Mandatory Core Agents
Every new project starts with:
- CEO
- CTO
- BA

Other agents are optional and can be added as the phase requires.

## Memory
Use the memory list to read allowlisted GSD files. The MVP reads approved files only and does not write to memory docs through the dashboard.

## Handoffs
Use the handoff form to create structured JSON handoffs. MVP handoffs are file-based and written under `.agent_handoff/`.

## Local Fallback
If the API is unavailable, provider, agent, and handoff saves use browser `localStorage`. This keeps human dashboard validation possible while backend persistent service startup is hardened in a later runtime task.
