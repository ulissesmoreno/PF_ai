# Phase 2 Runtime Routing User Manual

## Updated
- **Timestamp:** 2026-05-05 20:21
- **Responsible:** [DOCUMENTATION]
- **Scope:** Provider health and hybrid routing baseline.

## Purpose
Use Phase 2 controls to inspect whether a provider should route through API, local runtime, or hybrid fallback.

## Provider Health
The backend exposes:

```text
GET /api/provider-health?id=<provider_id>
```

The route is protected and requires the same bearer token used by the dashboard.

## Local Runtime Endpoint Contract
Local and hybrid providers must use loopback endpoints only:
- `localhost`
- `127.0.0.1`
- `::1`

Allowed schemes:
- `http`
- `https`

The backend performs a short HTTP GET to the configured endpoint and uses the result as local runtime health.

## Hybrid Routing
Hybrid provider behavior:
1. Use local route when the local runtime endpoint is healthy.
2. Fall back to API when local runtime is unhealthy and API secret reference exists.
3. Return an auditable fallback reason.
4. Return unavailable when no route is healthy.

## Dashboard
1. Open `http://127.0.0.1:5173`.
2. Connect to the backend API with the local token.
3. Register or select a provider.
4. Click `Check Status`.
5. Review the redacted route/status output.

## Backend Startup
Set `PF_AI_AUTH_SECRET` in the current shell, then run:

```powershell
scripts\run-backend-local.ps1
```

The script does not store secrets. It reads `PF_AI_AUTH_SECRET` from the current shell.
