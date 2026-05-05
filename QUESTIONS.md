# QUESTIONS

> Immutable log of questions and answers. Append only — never delete entries.
> Timestamps mandatory (YYYY-MM-DD HH:MM). Each question has an inline response line for [HUMAN].

## Entry Index (For Navigation)
- [Add entries here as questions are registered]

---

### [YYYY-MM-DD HH:MM] Question: Short title
- **Context:** [Context description]
- **Author:** [Agent tag]
- **Priority:** High / Medium / Low
- **Impact:** Architectural / Functional / Security / Business
- **Status:** Open / Answered / Blocked
- **Question:** [Full text]
  > **Response:** [HUMAN fills here]
- **Decision / Action:** [What was decided or action taken — filled after response]
- **References:** [Links to related files]

---

## Structured Markers
- [ ] Timestamp present
- [ ] Status filled
- [ ] Response line present for [HUMAN]
- [ ] Decision documented after answer
- [ ] Priority and Impact defined
- [ ] References added

---

### [2026-05-05 14:51:40] Question: Project onboarding required
- **Context:** GSD activation detected `DOC/PROJECT.md`, `NEW-INSTRUCTIONS.md`, and `DOC/PLAN.md` still containing placeholders/templates. Per `DOC/GSD-RULES.md` and `DOC/ONBOARDING.md`, project onboarding must run before any roadmap phase or implementation.
- **Author:** [CEO]
- **Priority:** High
- **Impact:** Business / Architectural / Security
- **Status:** Answered
- **Question:** Please provide the onboarding answers for the 5 required blocks: project identity, technology stack, visual identity, team/process, and roadmap/MVP. Execution remains blocked until these are confirmed.
  > **Response:** O projeto manterá o nome PF_ai e terá como objetivo criar um gerenciador de modelos, memória e agentes para garantir a funcionalidade ideal do boiler plate. Criará uma interface para interação humano-agentes, inspirado no paperclip ai (ver NEW-INSTRUCTIONS). É possível trabalhar com Go ou Runy on Rails para esse projeto? Identidade visual se basear em https://getdesign.md/composio/design-md. Times será apenas Eu e os agentes. Aceito sugestoe para MVP. Usar CAVEMAN para repsosta em Chat. Open Claude pode ser usada para inspiração de gerenciador de agentes.
- **Decision / Action:** [2026-05-05 15:06:38] [HUMAN] provided onboarding direction. CEO created a stack/MVP proposal for confirmation before auto-fill.
- **References:** `DOC/GSD-RULES.md`, `DOC/ONBOARDING.md`, `DOC/PROJECT.md`, `NEW-INSTRUCTIONS.md`, `DOC/PLAN.md`

---

### [2026-05-05 15:00:16] Question: Confirm onboarding auto-fill proposal
- **Context:** [HUMAN] answered the onboarding prompt in `QUESTIONS.md`. The answer defines the product direction but asks whether Go or Ruby on Rails can be used and requests MVP suggestions.
- **Author:** [CEO]
- **Priority:** High
- **Impact:** Business / Architectural
- **Status:** Answered
- **Question:** Confirm whether the project files should be auto-filled using this proposal: PF_ai as an agent/model/memory manager and human-agent interface; web-first product; Go backend recommended over Ruby on Rails for local agent orchestration, process control, and efficient services; lightweight web frontend inspired by `getdesign.md/composio/design-md`; PostgreSQL for durable project state plus file-based `.md` traceability; MVP focused on registering agents, reading/writing GSD memory files, showing active phase/status, and creating structured handoffs for human-agent interaction.
  > **Response:** Use Go. MVP aprovado, apenas quero que os agentes possam funcionar por API e/ou localmente (nessa caso será criado um docker dentro do projeto rodando o modelo localmente). handoffs via arquivo no MVP e via MCP no final
- **Decision / Action:** [2026-05-05 15:06:38] Approved with Go, file handoffs for MVP, MCP handoffs for final stage, and API/local agent execution. CEO proceeded with onboarding auto-fill.
- **References:** `DOC/ONBOARDING.md`, `DOC/PROJECT.md`, `QUESTIONS.md`, `PLAYBOOK.md`

---

### [2026-05-05 15:14:57] Question: Go and Hexagonal Architecture compatibility
- **Context:** `NEW-INSTRUCTIONS.md` asks whether Hexagonal Architecture remains viable when using Go.
- **Author:** [CEO]
- **Priority:** High
- **Impact:** Architectural
- **Status:** Answered
- **Question:** Using Go, is Hexagonal Architecture still viable for PF_ai?
  > **Response:** Yes. Go is compatible with Hexagonal Architecture. PF_ai should implement the domain as pure Go packages, application use cases as ports/interfaces, and infrastructure adapters for HTTP, PostgreSQL, Markdown files, `.agent_handoff/`, Docker/local runtime, and external model APIs.
- **Decision / Action:** Keep Hexagonal Architecture as the architectural standard for the Go backend. CTO must enforce domain isolation and adapter boundaries before Phase 1 implementation.
- **References:** `NEW-INSTRUCTIONS.md`, `DOC/ARCHITECTURE.md`, `DOC/PROJECT.md`, `DOC/PLAN.md`
