# ARCHITECTURE.md -

## 1. System Overview
This system is `[PROJECT_NAME]`. It uses **Hexagonal Architecture (Ports & Adapters)** to ensure that the business core is resilient, testable, and fully independent of external technologies.

## 2. Development Pillars (GSD Framework)
* **Get Shit Done (GSD):** Focus on functional delivery and high-impact executable code.
* **TDD:** Strict Red → Green → Refactor cycle before any commit.
* **Layered Security:** JWT implementation and rigorous sanitization in Adapters.
* **SOLID:** Extensible code for new models and technical indicators.
* **Token Economy:** Multi-agent architecture exists to minimize token consumption — specialized agents with limited scope over generalist agents.

## 3. Layer Organization (Hexagonal)

### A. Domain (The Heart)
Located in `src/domain`. **Golden Rule:** Zero dependency on external libraries.
- **Entities:** .
- **Value Objects:** .
- **Domain Services:** Complex business logic.

### B. Application (Use Cases)
Located in `src/application`. Orchestrates data flow between external world and domain.
- **Input Ports:** Command interfaces.
- **Output Ports:** Need interfaces.

### C. Infrastructure (Adapters)
Located in `src/infrastructure`. Technical implementation details.
- **Driving Adapters:** REST Controllers, Schedulers, CLI.
- **Driven Adapters:** PostgreSQL (JPA/SQL).

## 4. Technology Stack

| Component | Technology | Responsibility | Agent |
| :--- | :--- | :--- | :--- |
| **Backend** | Java 21 / Spring Boot 3 | Rule orchestration, API Gateway, Security | `[DEV_BACKEND]` |
| **Intelligence** | Python 3.12 / Scikit-Learn | Data capture, ML Modeling | `[DS/ML]` |
| **Data Pipelines** | Python / Spark (optional) | ETL, data quality, pipeline management | `[DATA_ENGINEER]` |
| **Frontend Web** | Angular 18 / TypeScript | Reactive visualization | `[DEV_FRONTEND]` |
| **Mobile** | React Native / Expo *(optional)* | iOS and Android | `[DEV_FRONTEND]` |
| **Desktop** | Electron / Tauri *(optional)* | Cross-platform desktop | `[DEV_FRONTEND]` |
| **Database** | PostgreSQL 16 | Central repository | `[DBA]` |
| **Security** | Spring Security + JWT | Stateless access control | `[SECURITY]` |
| **Container** | Docker / Compose | Environment orchestration | `[DEVOPS]` |
| **Messaging** | Apache Kafka | Async processing | `[DEV_BACKEND]` |
| **Cache** | Redis | Volatile data acceleration | `[DEV_BACKEND]` |
| **Monitoring** | Prometheus + Grafana | Metrics and dashboards | `[DEVOPS]` |
| **CI/CD** | GitHub Actions | Build, test, deploy | `[DEVOPS]` |

## 4.1 Alternative Frontends

| Type | Technology | When |
| :--- | :--- | :--- |
| Web | Angular 18 / React / Vue | Dashboards, admin panels |
| Mobile Hybrid | React Native + Expo | iOS + Android MVP |
| Mobile Native | Swift / Kotlin | Maximum hardware integration |
| Desktop | Electron / Tauri | Cross-platform with filesystem access |

> **Naming:** frontend → `<name>-web`; mobile → `<name>-app`; desktop → `<name>-desktop`; microservices → `<name>-service`.

## 4.2 Agent Communication Schemas (.agent_handoff/)

All handoff files must use structured JSON. Directory is git-ignored.

### 4.2.1 General Handoff Header (mandatory for all files)

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

### 4.2.2 Schema: Phase Kickoff (CEO → All Technical Agents)

```json
{
  "header": {
    "timestamp": "YYYY-MM-DD HH:MM",
    "sender": "[CEO]",
    "recipient": "[ALL_TECHNICAL]",
    "task_ref": "PHASE-X",
    "intent": "PHASE_KICKOFF"
  },
  "payload": {
    "phase_ref": "ROADMAP.md#Phase-X",
    "plan_ref": "PLAN.md#Stage-X",
    "seniority_level": "Junior | Pleno | Senior",
    "llm_tier": "1 | 2 | 3",
    "security_criteria": "Threat model summary from SECURITY kickoff",
    "applicable_skills": ["skill-1", "skill-2"],
    "constraints": ["TDD mandatory", "Hexagonal isolation"],
    "priority": "High | Medium | Low",
    "explicit_conclusion_authorized": false,
    "authorized_until_stage": null
  }
}
```

### 4.2.3 Schema: Clarification Request (Technical Agent → Management)

```json
{
  "header": {
    "timestamp": "YYYY-MM-DD HH:MM",
    "sender": "[DEV_BACKEND:Pleno]",
    "recipient": "[CTO]",
    "task_ref": "TASK-XXX",
    "intent": "CLARIFICATION_REQUEST"
  },
  "payload": {
    "doubt_type": "scope | architectural",
    "description": "Clear description of the doubt",
    "autonomous_decision_attempted": false,
    "blocking": true
  }
}
```

### 4.2.4 Schema: Review Result (Reviewer → TECH_LEAD)

```json
{
  "header": {
    "timestamp": "YYYY-MM-DD HH:MM",
    "sender": "[SECURITY:Senior]",
    "recipient": "[TECH_LEAD]",
    "task_ref": "TASK-XXX",
    "intent": "REVIEW_RESULT"
  },
  "payload": {
    "status": "APPROVED | REJECTED | APPROVED_WITH_NOTES",
    "round": 1,
    "findings": [],
    "corrective_actions": [],
    "retrospective_note": "Optional phase observation",
    "playbook_update": false
  }
}
```

### 4.2.5 Schema: Stage Approval (TECH_LEAD → CTO/BA)

```json
{
  "header": {
    "timestamp": "YYYY-MM-DD HH:MM",
    "sender": "[TECH_LEAD]",
    "recipient": "[CTO]",
    "task_ref": "PHASE-X",
    "intent": "STAGE_APPROVAL"
  },
  "payload": {
    "status": "APPROVED | REJECTED | ESCALATED_TO_HUMAN",
    "review_rounds_used": 1,
    "consolidated_findings": [],
    "security_cleared": true,
    "checklist_validated": true,
    "retrospective_notes": []
  }
}
```

### 4.2.6 Schema: Delivery (Technical Agent → TECH_LEAD)

```json
{
  "header": {
    "timestamp": "YYYY-MM-DD HH:MM",
    "sender": "[DEV_BACKEND:Pleno]",
    "recipient": "[TECH_LEAD]",
    "task_ref": "TASK-XXX",
    "intent": "REVIEW_RESULT"
  },
  "payload": {
    "artifacts": ["src/domain/entity/X.java"],
    "test_file": "TESTS.md#Test-X",
    "coverage_report": "target/site/jacoco/index.html",
    "status": "Green",
    "retrospective_note": "Optional observation",
    "playbook_update": false
  }
}
```

### 4.2.7 Schema: Bugfix Request (TECH_LEAD → DEVOPS + DEV)

```json
{
  "header": {
    "timestamp": "YYYY-MM-DD HH:MM",
    "sender": "[TECH_LEAD]",
    "recipient": "[DEVOPS]",
    "task_ref": "BUG-XXX",
    "intent": "BUGFIX_REQUEST"
  },
  "payload": {
    "severity": "Critical | High | Medium | Low",
    "branch_name": "bugfix/phase-X-description",
    "next_phase_frozen": true,
    "root_cause": "Description",
    "affected_phase": "Phase-X"
  }
}
```

## 5. Document Ownership

| File | Owner | Scope |
| :--- | :--- | :--- |
| `STATE.md` | Technical agents | Execution memory |
| `CONTEXT.md` | Management agents (CEO, CTO, BA) | Decisions and justifications |
| `PLAYBOOK.md` | CEO | [HUMAN] work preferences |
| `RETROSPECTIVE.md` | CEO | Phase learnings |
| `VERSIONS.md` | Phase-closing agent | Release history |
| `TESTS.md` | QA + SECURITY + Technical | Test results |
| `DESIGN.md` | UX_RESEARCHER + DEV_FRONTEND | Visual identity |
| `wiki/` | DOCUMENTATION | Project knowledge base |
| `wiki/user-manual/` | DOCUMENTATION | End-user documentation |

## 6. Machine Learning Optimization

- `n_jobs=-1` mandatory on all compatible estimators.
- Pandas with optimized data types (`float32`), vectorized operations only.
- Bulk Insert persistence to avoid I/O bottlenecks.

## 7. Definition of "Done" (DoD)

1. **Tests:** Unit and integration coverage validated per quality standards (GSD-RULES §8).
2. **Security:** SECURITY conformance audit cleared.
3. **Documentation:** All operational `.md` files updated. DOCUMENTATION wiki updated.
4. **Code:** Clean Code and SOLID patterns applied.
5. **STATE:** `STATE.md` updated by technical agents.
6. **Review:** TECH_LEAD consolidated approval delivered to CTO/BA.

## 8. LLM Tiering

Seniority and LLM Tier are separate concepts. When an external orchestrator is available, route tasks to specific models per tier:

| Tier | Model Placeholder | Purpose |
| :--- | :--- | :--- |
| 1 — Efficiency | `{{TIER_1_EFFICIENCY_MODEL}}` | Logs, formatting, linting |
| 2 — Development | `{{TIER_2_DEVELOPMENT_MODEL}}` | Standard TDD, feature implementation |
| 3 — Expert | `{{TIER_3_EXPERT_MODEL}}` | Architectural decisions, complex logic |

> Model routing is a strong recommendation — depends on project infrastructure.
