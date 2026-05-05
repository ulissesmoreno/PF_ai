# GSD-RULES.md - Technical Execution Protocol (v1.5)

This document defines the inviolable development rules. Any deviation invalidates the "Done" state.

---

## 👑 §0. Agent Orchestration Protocol (Chain of Command)

### 0.1 CEO Precedence — Inviolable

The **[CEO]** is the first to act in every new phase. Reading sequence (strict order, once per phase):
1. `DOC/GSD-RULES.md`
2. `PLAYBOOK.md`
3. `NEW-INSTRUCTIONS.md`
4. `DOC/PLAN.md`
5. `DOC/ONBOARDING.md` — only on new projects with placeholder `PROJECT.md`

> Re-reading only occurs after direct [HUMAN] intervention mid-phase.

### 0.2 Human Interface

- **[HUMAN]** communicates exclusively with **[CEO]**.
- Every chat response is signed by the **active agent** (e.g., `[CTO]`, `[DEV_BACKEND:Senior]`).
- Active agent responds directly — CEO does not paraphrase.
- **Chat responses: state only what is being done. No narration, no process description.**
  - ✅ `[DEV_BACKEND:Pleno] Implementing JWT filter.`
  - ❌ `[DEV_BACKEND:Pleno] I will now proceed to implement the JWT filter as described in PLAN.md...`
- Detailed content goes into `.md` files, never into chat.
- **[HUMAN] acts only at phase closure** — approves or rejects the deliverable.
- Mid-phase interruptions only for: Critical bug, credential leak, impasse after 3 review rounds.

### 0.3 Management Autonomy

- CEO, CTO, and BA have full autonomy to execute phases until the deliverable is ready.
- No human approval required mid-phase unless one of the three exceptions above occurs.
- **Explicit Stage Conclusion Authorization:** [HUMAN] may extend autonomy explicitly in `NEW-INSTRUCTIONS.md`, specifying the stage limit. Expires automatically. Active agent notes: `[CEO] Concluding Stage X — authorized via NEW-INSTRUCTIONS.md [YYYY-MM-DD].`

### 0.4 Phase Kickoff — Single Intake

At the start of each phase:
1. CEO + CTO + BA perform full intake — **once per phase only**.
2. CTO assigns seniority based on BA business complexity input — recorded in `.agent_handoff/`, never in chat.
3. SECURITY performs **threat modeling** — maps attack surfaces, defines phase security criteria.
4. CEO generates `phase_context` snapshot and dispatches single `PHASE_KICKOFF` handoff to all technical agents.
5. Technical agents work from the handoff snapshot — no file reading for Junior/Pleno tasks.

### 0.5 Seniority Assignment (CTO-led)

- **BA** assesses business complexity → informs CTO via handoff.
- **CTO** decides seniority — BA does not co-decide.

| Complexity | Seniority | File Reading |
| :--- | :--- | :--- |
| Low — mechanical, well-defined | Junior | Handoff only |
| Medium — standard with some design | Pleno | Handoff only |
| High — architectural, complex domain | Senior | Handoff + `ARCHITECTURE.md` if needed |

> **Token Economy Rule:** Always assign the lowest sufficient seniority. The multi-agent objective is token economy — specialized agents with limited scope consume less context than a single generalist. When an external orchestrator is available, it **should** route to specific models per tier.

### 0.6 Technical Agents

| Agent | Scope |
| :--- | :--- |
| `[DEV_BACKEND]` | Backend implementation (TDD) |
| `[DEV_FRONTEND]` | Frontend implementation (TDD) |
| `[DBA]` | Schema, migrations, persistence |
| `[DS/ML]` | ML models, data pipelines |
| `[DEVOPS]` | CI/CD, containers, infrastructure |
| `[DATA_ENGINEER]` | ETL, data pipelines, data quality |

### 0.7 Escalation Protocol

- Technical agents escalate **once per phase** via `CLARIFICATION_REQUEST` handoff.
- Management responds via handoff if possible — no [HUMAN] interruption.
- If management cannot resolve → `QUESTIONS.md` → [HUMAN] → development freezes.
- **Autonomous decision rule:** Agent decides independently on *how to implement*. Escalates only for *what to implement* or *architectural scope*.

### 0.8 Quality & Review Gates

**Parallel review at phase closure:**

| Reviewer | Scope | Always included |
| :--- | :--- | :--- |
| `[SECURITY]` | Vulnerabilities, JWT, sanitization, credential scan | ✅ |
| `[CODE_REVIEWER]` | SOLID, Clean Code, Hexagonal isolation | ✅ |
| `[QA]` | Acceptance criteria, test coverage | ✅ |
| `[DBA]` | Schema/migration changes | Only if applicable |
| `[DEVOPS]` | Infrastructure/pipeline changes | Only if applicable |

**Conflict resolution:**
- *How code was written* → `[CODE_REVIEWER]` decides.
- *Whether behavior is correct* → `[QA]` decides.
- Cross-scope → `[TECH_LEAD]` arbitrates → `[CTO]` last resort.

**Review rounds:** Max 3. After 3rd without resolution → `QUESTIONS.md` → [HUMAN].

**Approval flow:**
```
Reviewers (parallel) → TECH_LEAD consolidates → CTO/BA approve → [HUMAN] releases PR → Merge
```

> Every phase must produce a functional, validated, and tested codebase ready for merge.

### 0.9 SECURITY as Structural Foundation

- **Phase Kickoff:** Threat modeling — attack surfaces, security criteria for the phase.
- **Phase Closure:** Conformance audit — verifies implementation against kickoff threat model.
- **Pre-production:** Full audit before any merge to production branch. Veto power active.
- These are distinct activities — no overlap.

### 0.10 Creative & Market Agents

Activated on demand by CEO or CMO:

| Agent | Output |
| :--- | :--- |
| `[CMO]` | KPI review, roadmap alignment |
| `[WRITER]` | Marketing copy |
| `[ARTIST]` | Visual marketing assets |

> **Dormant Agent Rule:** Agents with no active domain remain dormant and do not consume context.

---

## §1. Document Ownership & Write Rules

### 1.1 When to Write

| File | Written | By |
| :--- | :--- | :--- |
| `STATE.md` | During and at end of phase | Technical agents |
| `NEW-INSTRUCTIONS.md` | Only if [HUMAN] intervenes mid-phase | [HUMAN] exclusively |
| `CONTEXT.md` | Phase closure only | Management agents (CEO, CTO, BA) |
| `PLAN.md` | Phase closure only | BA |
| `TASKS.md` | Phase closure only | All agents |
| `TESTS.md` | Phase closure only | QA / Technical agents |
| `WIKI` | Phase closure only | DOCUMENTATION agent |
| `VERSIONS.md` | Phase closure only, from Phase 1 onward | Phase-closing agent |
| `RETROSPECTIVE.md` | Phase closure only | CEO |
| `QUESTIONS.md` | Only when agent needs to ask [HUMAN] | Any agent |
| `ROADMAP.md` | Marked by [HUMAN] at closure | [HUMAN] |
| `PLAYBOOK.md` | When any interaction changes understanding of [HUMAN] work preferences | CEO |

### 1.2 File Purpose (Inviolable Definitions)

- **PLAYBOOK.md** — [HUMAN]'s work preferences and working style. Updated only when understanding changes.
- **CONTEXT.md** — Architectural and project decisions with justifications. Owned by management agents.
- **STATE.md** — Agent execution memory: what is done, pending, and blocked. Owned by technical agents.
- **QUESTIONS.md** — Exclusively for agents to ask [HUMAN]. No decisions, no logs, no architecture notes.
- **WIKI** — Project knowledge base + end-user manual (`wiki/user-manual/`). Filled at phase closure.
- **VERSIONS.md** — Release history. Written at phase closure, starting from Phase 1.

### 1.3 Reading Rules

- **Technical agents (Junior/Pleno):** Handoff only — no file reading.
- **Technical agents (Senior):** Handoff + `ARCHITECTURE.md` only if task requires it.
- **Management agents:** Full intake once per phase. No re-reading unless [HUMAN] intervenes.
- **QUESTIONS.md:** Written only when registering a question for [HUMAN]. Never read at kickoff.
- **RETROSPECTIVE.md / WIKI:** Written at closure only.

---

## §2. Phase Lifecycle

### 2.1 Phase Start
1. CEO + CTO + BA: full intake (once).
2. SECURITY: threat modeling.
3. CTO: seniority assignment via handoff.
4. CEO: generates `phase_context` snapshot + dispatches single `PHASE_KICKOFF` handoff.

### 2.2 During Phase
- Technical agents execute from handoff context.
- `STATE.md` updated as tasks complete.
- `NEW-INSTRUCTIONS.md` only if [HUMAN] intervenes.
- `playbook_update: true` flag in handoff when a decision reveals a work preference.
- `retrospective_note` field populated in handoffs during the phase — consolidated at closure.

### 2.3 Phase Closure
1. Parallel review: SECURITY + CODE_REVIEWER + QA ± DBA/DEVOPS.
2. TECH_LEAD consolidates → CTO/BA approve.
3. All documentation written: CONTEXT, PLAN, TASKS, TESTS, WIKI, VERSIONS (Phase 1+), RETROSPECTIVE.
4. CEO updates PLAYBOOK if `playbook_update: true` flags received.
5. Stage Closure Checklist presented to [HUMAN].
6. [HUMAN] approves → marks ROADMAP → releases PR → Merge.

### 2.4 Stage Closure Checklist

| # | File | What to Check |
| :- | :--- | :--- |
| 1 | `TESTS.md` | All tests recorded with results and timestamp |
| 2 | `STATE.md` | Updated with delivery and metrics |
| 3 | `TASKS.md` | All tasks Done; estimated vs. actual time recorded |
| 4 | `CONTEXT.md` | Architectural decisions documented |
| 5 | `README.md` | Reflects delivered code reality |
| 6 | `ROADMAP.md` | Stage marked `[x]` by [HUMAN] |
| 7 | `PLAN.md` | Closure Gate section filled |
| 8 | `VERSIONS.md` | Entry added (Phase 1 onward) |
| 9 | `RETROSPECTIVE.md` | Phase entry consolidated by CEO |
| 10 | `WIKI` | Updated with phase deliverables |

> **[HUMAN] is the sole approval authority for stage advancement.**

---

## §3. Rollback Protocol

### 3.1 Severity Levels

| Level | Criteria | Rollback |
| :--- | :--- | :--- |
| **Critical** | Security vulnerability, data loss, system down | Mandatory — immediate |
| **High** | Core feature regression, all-users impact | Mandatory — [HUMAN] authorized |
| **Medium** | Secondary flow bug, performance degradation | Optional — [HUMAN] decides |
| **Low** | Cosmetic, rare edge case | No rollback — fix in next phase |

> **Critical only:** DEVOPS may initiate preventive rollback while awaiting [HUMAN].

### 3.2 Rollback Flow
```
Bug identified + severity defined
→ Critical/High: next phase FREEZES
→ SECURITY/QA registers in QUESTIONS.md
→ [HUMAN] authorizes rollback
→ DEVOPS executes rollback to last stable version
→ TECH_LEAD opens priority task in TASKS.md
→ Records: VERSIONS.md, RETROSPECTIVE.md, CONTEXT.md (if architectural)
```

---

## §4. Bugfix Branch Protocol

```
Bug identified
→ Critical/High: next phase FREEZES
→ DEVOPS creates: bugfix/phase-X-<description>
→ TECH_LEAD opens priority task in TASKS.md
→ CTO/BA focused intake (no full file reading)
→ DEV fixes in bugfix branch
→ Full review (SECURITY mandatory for Critical/High)
→ [HUMAN] authorizes merge → main + next phase branch
→ Next phase resumes
```

| Severity | Next Phase | Action |
| :--- | :--- | :--- |
| Critical | Freezes immediately | Absolute priority |
| High | Freezes | Starts within 1 cycle |
| Medium | Continues | Parallel isolated branch |
| Low | Continues | Normal backlog |

---

## §5. Credential Leak Protocol

Always **Critical** — no exceptions. No agent resolves alone.

```
Leak identified
→ SECURITY registers in QUESTIONS.md immediately
→ Development FREEZES
→ [HUMAN] revokes credential at source
→ DEVOPS removes from git history
→ SECURITY validates cleanup
→ New credential via vault
→ SECURITY audits all phase files
→ Development resumes after SECURITY clearance
```

---

## §6. Quality Standards

### 6.1 Test Quality

| Layer | Criteria |
| :--- | :--- |
| Domain | Mutation score ≥ 80% |
| Application | 100% use case coverage (happy + unhappy path) |
| Infrastructure | Contract tests |
| Frontend | Behavior tests, not implementation tests |

**By seniority:**
- Junior: happy path mandatory.
- Pleno: happy + unhappy path.
- Senior: mutation testing + edge cases.

### 6.2 Conflict Resolution

- *How code was written* → `[CODE_REVIEWER]`.
- *Whether behavior is correct* → `[QA]`.
- Cross-scope → `[TECH_LEAD]` → `[CTO]`.

---

## §7. TDD Development Cycle (Strict Mode)

1. **RED:** Failing test.
2. **GREEN:** Minimum code to pass.
3. **REFACTOR:** SOLID + Clean Code.
4. **SECURITY:** Sanitization, injection protection, JWT.

---

## §8. Hexagonal Architecture & Isolation

- **DOMAIN:** Pure logic. No framework imports.
- **APPLICATION (Ports):** Input/output contract interfaces.
- **INFRASTRUCTURE (Adapters):** DB, APIs, UI implementations.

**Naming:**
- Web → `<name>-web` | Mobile → `<name>-app` | Desktop → `<name>-desktop` | Service → `<name>-service`

---

## §9. Security & Persistence

- API routes born `deny-all`. Released via JWT only.
- No logging of keys, payloads, or sensitive data.
- Adapters sanitize data before Domain.
- Credentials never in code or docs. `.env` never committed.
- Structured logs mandatory: INFO / WARN / ERROR.

---

## §10. LLM Tiering & Model Routing

Seniority and LLM Tier are **separate concepts**.
- Seniority → autonomy scope.
- Tier → which model runs.

| Tier | Purpose |
| :--- | :--- |
| Tier 1 — Efficiency | Logs, formatting, linting, repetitive tasks |
| Tier 2 — Development | Standard TDD, feature implementation |
| Tier 3 — Expert | Architectural decisions, complex domain, security design |

> When an external orchestrator is available, it **should** route tasks to specific models per tier. The multi-agent objective is token economy.

---

## §11. Installed Skills & Environment Tools

CTO is the sole responsible agent for skill mapping per phase. Recorded in `CONTEXT.md` at closure.

| Domain | Skills | When |
| :--- | :--- | :--- |
| Architecture | `architect-review`, `senior-architect` | Design review |
| Backend Java | `api-patterns`, `backend-architect` | API design |
| Python/ML | `scikit-learn`, `ml-pipeline-workflow` | ML pipelines |
| Frontend | `frontend-design`, `react-patterns` | UI componentization |
| Database | `database-design`, `postgres-best-practices` | Schema, migrations |
| Security | `security-auditor`, `differential-review` | Audits |
| TDD | `tdd-workflow`, `webapp-testing` | Test cycles |
| DevOps | `docker-expert`, `github-actions-templates` | Pipelines |
| Documentation | `documentation`, `wiki-architect` | Doc generation |
| Debugging | `systematic-debugging` | Bug diagnosis |

**External tools:** ESLint, Prettier, Checkstyle, SonarQube, Semgrep, k6, Dependabot, Swagger, Mermaid, HashiCorp Vault, Docker Compose, Kafka, Redis.

---

## §12. Documentation as Code

| File | Location | Owner | Written |
| :--- | :--- | :--- | :--- |
| `README.md` | Root | DOCUMENTATION | Post-onboarding; updated at closure |
| `PLAYBOOK.md` | Root | CEO | On preference change |
| `NEW-INSTRUCTIONS.md` | Root | [HUMAN] exclusively | On new instructions |
| `QUESTIONS.md` | Root | Any agent | Only to ask [HUMAN] |
| `RETROSPECTIVE.md` | Root | CEO | Phase closure |
| `GSD-RULES.md` | `DOC/` | — | Inviolable |
| `ARCHITECTURE.md` | `DOC/` | CTO | Phase closure |
| `PROJECT.md` | `DOC/` | CEO / BA | Onboarding |
| `PLAN.md` | `DOC/` | BA | Phase closure |
| `ROADMAP.md` | `DOC/` | CEO / [HUMAN] | Phase closure |
| `STATE.md` | `DOC/` | Technical agents | During + closure |
| `CONTEXT.md` | `DOC/` | Management agents | Phase closure |
| `TASKS.md` | `DOC/` | All agents | Phase closure |
| `TESTS.md` | `DOC/` | QA / Technical | Phase closure |
| `ENV_SETUP.md` | `DOC/` | DEVOPS / SECURITY | Phase closure |
| `DESIGN.md` | `DOC/` | UX_RESEARCHER | Phase closure |
| `VERSIONS.md` | `DOC/` | Phase-closing agent | Phase closure, Phase 1+ |
| `ONBOARDING.md` | `DOC/` | CEO | Once per new project |
| `wiki/` | `wiki/` | DOCUMENTATION | Phase closure |
| `wiki/user-manual/` | `wiki/` | DOCUMENTATION | Phase closure |

---

**AGENT SIGNATURE:** Operating in GSD Mode — Execution on demand, quality by design. v1.5
