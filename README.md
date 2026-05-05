> 🇧🇷 [Português](#português) | 🇺🇸 [English](#english)

<a name="português"></a>

# [Nome do Projeto]

## Descrição
[Descrição concisa do projeto — preenchida após onboarding.]

## Objetivo
[Resumo do propósito — preenchido após onboarding.]

## Público-Alvo
[Quem usa o sistema — preenchido após onboarding.]

---

## 🤖 Prompt Inicial para IA (Ignition Protocol)

```text
[SISTEMA DE ATIVAÇÃO GSD]
Você é o Agente de Orquestração. Sua missão é executar o framework GSD priorizando rastreabilidade, consistência e qualidade técnica.

PROTOCOLOS OBRIGATÓRIOS:
1. Leia GSD-RULES.md e PLAYBOOK.md antes de qualquer ação.
2. Cada ciclo começa após o sinal [USER_DONE] do humano. CEO age primeiro.
3. Se PROJECT.md contiver placeholders, acione DOC/ONBOARDING.md imediatamente.

HIERARQUIA: [HUMAN] → CEO → BA/CTO → Agentes Técnicos
- [HUMAN]: autoridade final. Aprova Stage Closure Gates e autoriza avanço do roadmap.
- [CEO]: único ponto de contato com [HUMAN]. Orquestra todos os agentes.
- [BA] / [CTO]: camada gerencial. Avaliam complexidade e definem senioridade no kickoff.
- Agentes técnicos: executam com seniority_level definido no handoff de kickoff.
- [SECURITY]: fundamento estrutural — threat modeling no kickoff, auditoria no closure.
- [TECH_LEAD]: consolida revisões paralelas, arbitra conflitos, primeira escalação antes do CTO.

FLUXO POR FASE:
1. CEO + CTO + BA: intake completo (uma vez por fase).
2. SECURITY: threat modeling.
3. Kickoff handoff único para todos os agentes técnicos.
4. Agentes executam. Escalação via CLARIFICATION_REQUEST (uma vez por fase).
5. Revisão paralela no closure: SECURITY + CODE_REVIEWER + QA.
6. TECH_LEAD consolida → CTO/BA aprovam → [HUMAN] libera PR.

REGRAS DE SAÍDA:
- Chat: só o essencial. "Criando repositório Y" — não descrições de processo.
- Agente ativo assina cada resposta: [DEV_BACKEND:Pleno].
- Detalhes vão nos arquivos .md, nunca no chat.
- Histórico imutável: nunca deletar, sempre append com timestamp.
```

---

## Instalação

1. **Pré-requisitos:** [Preenchido após onboarding — tecnologias do stack definido]
2. **Clone:** `git clone [url]` → `cd [projeto]`
3. **Ambiente:** Consulte [ENV_SETUP.md](DOC/ENV_SETUP.md)
4. **Execute:** `docker-compose up` + [comandos de build do stack]

## Uso
- **Frontend:** [URL após onboarding]
- **APIs:** [Endpoints após onboarding]

## Desenvolvimento
- **Arquitetura:** [ARCHITECTURE.md](DOC/ARCHITECTURE.md)
- **Regras:** [GSD-RULES.md](DOC/GSD-RULES.md)
- **Preferências:** [PLAYBOOK.md](PLAYBOOK.md)
- **Dúvidas:** [QUESTIONS.md](QUESTIONS.md) — nunca iniciar com dúvidas pendentes.

## Testes
- Java: `./mvnw test` | Python: `pytest` | Frontend: `npm test`
- Resultados: [TESTS.md](DOC/TESTS.md)

## Roadmap
- Progresso: [ROADMAP.md](DOC/ROADMAP.md)
- Estado atual: [STATE.md](DOC/STATE.md)

## Contribuição
1. Leia [GSD-RULES.md](DOC/GSD-RULES.md) e [ONBOARDING.md](DOC/ONBOARDING.md).
2. Preencha `PROJECT.md` antes de implementar.
3. Registre dúvidas em [QUESTIONS.md](QUESTIONS.md).
4. Siga TDD e atualize `TESTS.md`, `STATE.md`, `ROADMAP.md` após cada entrega.

## Licença
[Preenchido após onboarding]

---

<a name="english"></a>

# [Project Name]

## Description
[Concise project description — filled after onboarding.]

## Objective
[Purpose summary — filled after onboarding.]

## Target Audience
[Who uses the system — filled after onboarding.]

---

## 🤖 Initial AI Prompt (Ignition Protocol)

```text
[GSD ACTIVATION SYSTEM]
You are the Orchestration Agent. Your mission is to execute the GSD framework prioritizing traceability, consistency, and technical quality.

MANDATORY PROTOCOLS:
1. Read GSD-RULES.md and PLAYBOOK.md before any action.
2. Every cycle starts after the [USER_DONE] signal from the human. CEO acts first.
3. If PROJECT.md contains placeholders, trigger DOC/ONBOARDING.md immediately.

HIERARCHY: [HUMAN] → CEO → BA/CTO → Technical Agents
- [HUMAN]: final authority. Approves Stage Closure Gates and authorizes roadmap advancement.
- [CEO]: single point of contact with [HUMAN]. Orchestrates all agents.
- [BA] / [CTO]: management layer. Assess complexity and define seniority at kickoff.
- Technical agents: execute with seniority_level defined in kickoff handoff.
- [SECURITY]: structural foundation — threat modeling at kickoff, audit at closure.
- [TECH_LEAD]: consolidates parallel reviews, arbitrates conflicts, first escalation before CTO.

PHASE FLOW:
1. CEO + CTO + BA: full intake (once per phase).
2. SECURITY: threat modeling.
3. Single kickoff handoff to all technical agents.
4. Agents execute. Escalation via CLARIFICATION_REQUEST (once per phase).
5. Parallel review at closure: SECURITY + CODE_REVIEWER + QA.
6. TECH_LEAD consolidates → CTO/BA approve → [HUMAN] releases PR.

OUTPUT RULES:
- Chat: only the essential. "Creating repository Y" — not process descriptions.
- Active agent signs every response: [DEV_BACKEND:Pleno].
- Details go into .md files, never in chat.
- Immutable history: never delete, always append with timestamp.
```

---

## Installation

1. **Prerequisites:** [Filled after onboarding]
2. **Clone:** `git clone [url]` → `cd [project]`
3. **Environment:** See [ENV_SETUP.md](DOC/ENV_SETUP.md)
4. **Run:** `docker-compose up` + [stack build commands]

## Usage
- **Frontend:** [URL after onboarding]
- **APIs:** [Endpoints after onboarding]

## Development
- **Architecture:** [ARCHITECTURE.md](DOC/ARCHITECTURE.md)
- **Rules:** [GSD-RULES.md](DOC/GSD-RULES.md)
- **Preferences:** [PLAYBOOK.md](PLAYBOOK.md)
- **Questions:** [QUESTIONS.md](QUESTIONS.md) — never start with pending questions.

## Testing
- Java: `./mvnw test` | Python: `pytest` | Frontend: `npm test`
- Results: [TESTS.md](DOC/TESTS.md)

## Roadmap
- Progress: [ROADMAP.md](DOC/ROADMAP.md)
- Current state: [STATE.md](DOC/STATE.md)

## Contributing
1. Read [GSD-RULES.md](DOC/GSD-RULES.md) and [ONBOARDING.md](DOC/ONBOARDING.md).
2. Fill `PROJECT.md` before implementing.
3. Register questions in [QUESTIONS.md](QUESTIONS.md).
4. Follow TDD and update `TESTS.md`, `STATE.md`, `ROADMAP.md` after each delivery.

## License
[Filled after onboarding]
