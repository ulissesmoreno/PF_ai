# Phase 3 Orchestration Board User Manual

## Updated
- **Timestamp:** 2026-05-06 07:33
- **Responsible:** [DOCUMENTATION]
- **Scope:** Project registration, planning cards, and MCP-compatible handoff export.

## Purpose
Use the Phase 3 dashboard as the operational planning surface for PF_ai. Planning cards show what exists, who owns it, priority, status, phase, and task reference.

## Project Registration
1. Open `http://127.0.0.1:5173`.
2. Fill `project id` and `project name`.
3. Fill onboarding JSON.
4. Click `Register Project`.

Secret-looking fields such as `token`, `password`, `secret`, or `api_key` are rejected.

## Planning Board
Cards are grouped by:
- To Do
- In Progress
- Done

Each card shows:
- Title
- Priority
- Responsible agent
- Phase
- Task reference

If the backend API is unavailable, cards are saved locally in browser `localStorage` for dashboard validation.

## MCP Envelope Export
1. Fill the handoff form.
2. Click `Export MCP Envelope`.
3. Review the generated `mcp-compatible` envelope.

The envelope wraps the existing validated handoff schema. Payloads with secret-like fields are rejected.
