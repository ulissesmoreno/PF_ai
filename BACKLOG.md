# BACKLOG — pf_ai Orchestration System
> Revisado em: 2026-05-13 | Base: análise do código-fonte + diretivas do [HUMAN]
> Prioridade: P1 (crítico/MVP) → P2 (ganho direto) → P3 (futuro/nice-to-have)
> Arquivos excluídos de CQRS: `README.md`, `GSD-RULES.md`

---

## Épico 0 — Projetos e Onboarding Estruturado

> **Contexto:** O sistema hoje assume um único projeto implícito — o `context_store` não tem isolamento, `ProjectStarted()` usa heurística frágil de placeholder em `.md`, e o onboarding é conduzido inteiramente via arquivos. A meta é persistir dados de projeto de forma estruturada, suportar múltiplos projetos ativos no mesmo banco, e conduzir o onboarding via card interativo respondido pelo humano.
>
> **Decisão arquitetural:** Isolamento por `project_id` (coluna UUID) em todas as tabelas de dados — não múltiplos bancos. SQLite não tem schemas; múltiplos `.db` complicam migrações, backups e consultas cross-project. Com `project_id` o dashboard exibe projetos lado a lado e consultas analíticas (tokens, cards, contexto) funcionam por projeto ou agregadas.

### [P1] E0-T1 — Tabela `projects` e propagação de `project_id`
**O que fazer:**
- Criar tabela `projects` na migration `003_projects.sql` (renumerar as migrations existentes a partir de `004`):
  ```sql
  CREATE TABLE projects (
    id TEXT PRIMARY KEY,              -- UUID v4 gerado no onboarding
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,        -- identificador curto p/ logs: "pf-ai", "crm-v2"
    domain TEXT,                      -- financeiro | saúde | logística | produtividade | outro
    description TEXT,
    target_audience TEXT,
    main_objective TEXT,
    stack_backend TEXT,
    stack_frontend TEXT,
    stack_database TEXT,
    stack_infra TEXT,
    stack_notes TEXT,                 -- campo livre para tecnologias adicionais
    status TEXT NOT NULL DEFAULT 'active',  -- active | paused | archived
    workspace_root TEXT,              -- caminho do workspace no filesystem
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
  );

  CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
  CREATE INDEX IF NOT EXISTS idx_projects_slug ON projects(slug);
  ```
- Adicionar coluna `project_id TEXT REFERENCES projects(id)` nas tabelas:
  - `cards`, `card_comments`, `token_usage`
  - `context_entries`, `planning_items`, `test_records`, `agent_outputs`, `handoff_events`
  - `documents`, `document_imports`, `document_sections`, `document_entries`
  - `knowledge_chunks`, `read_context_items`
- `project_id` é `NOT NULL` em tabelas novas (E2 em diante); `nullable` nas tabelas existentes com dados (retrocompatibilidade via `DEFAULT NULL`)
- Adicionar `project_id` como parâmetro obrigatório em `Store.Open()` e propagar para todos os métodos de escrita e leitura
- `QueryContextForHandoff` filtra por `project_id` — agentes de um projeto não recebem contexto de outro
- `Config` em `config.go`: adicionar `ActiveProjectID string` carregado do banco na inicialização (o projeto `active` mais recente); se nenhum existir, disparar onboarding

**Critério de aceite:** Dois projetos no banco; `QueryContextForHandoff` com `project_id=A` não retorna registros do projeto B; migration roda sem erro em banco existente (retrocompatibilidade).

---

### [P1] E0-T2 — Onboarding via card estruturado (CEO conduz)
**O que fazer:**
- Substituir a heurística `ProjectStarted()` (verifica placeholder em `PROJECT.md`) por: consulta à tabela `projects` — se `COUNT(*) WHERE status='active' = 0`, disparar onboarding
- No `createCEOKickoff` em `main.go`: quando sem projeto ativo, gerar handoff com `intent: PROJECT_ONBOARDING` em vez de `PHASE_KICKOFF`
- O CEO, ao receber `PROJECT_ONBOARDING`, gera um card com `status=blocked` e `comment_type=form` contendo as perguntas estruturadas como JSON:
  ```json
  {
    "action": "ask_human",
    "card_title": "Configuração do novo projeto",
    "questions": [
      { "id": "name",            "label": "Nome do projeto",             "type": "text",   "required": true },
      { "id": "slug",            "label": "Identificador curto (slug)",  "type": "text",   "required": true, "hint": "ex: pf-ai, crm-v2" },
      { "id": "domain",          "label": "Domínio",                     "type": "select", "options": ["financeiro","saúde","logística","produtividade","educação","outro"] },
      { "id": "description",     "label": "Descrição em uma frase",      "type": "text",   "required": true },
      { "id": "target_audience", "label": "Público-alvo",                "type": "text",   "required": true },
      { "id": "main_objective",  "label": "Objetivo principal",          "type": "textarea" },
      { "id": "stack_backend",   "label": "Backend",                     "type": "text",   "hint": "ex: Go 1.22, Java 21" },
      { "id": "stack_frontend",  "label": "Frontend",                    "type": "text",   "hint": "ex: React 18, Angular 18" },
      { "id": "stack_database",  "label": "Banco de dados",              "type": "text",   "hint": "ex: SQLite, PostgreSQL 16" },
      { "id": "stack_infra",     "label": "Infraestrutura",              "type": "text",   "hint": "ex: Docker, k8s, serverless" },
      { "id": "stack_notes",     "label": "Outras tecnologias",          "type": "textarea","required": false }
    ]
  }
  ```
- Quando o humano responde (renomeia o handoff `_TO_HUMAN_` → `_TO_CEO_`), o pipeline chama `Store.CreateProject(data)` que:
  - Gera UUID v4 para `id`
  - Deriva `slug` do campo informado (lowercase, hifens, sem espaços)
  - Persiste na tabela `projects`
  - Define `ActiveProjectID` no runtime
  - Registra comentário no card: `[PROJECT_CREATED] id=<uuid> slug=<slug>`
  - Dispara o kickoff normal (`PHASE_KICKOFF`) com o novo `project_id`
- Método `Store.GetActiveProject() (*Project, error)` retorna o projeto `active` mais recente
- Método `Store.ListProjects() ([]Project, error)` retorna todos (para o dashboard futuro)
- Método `Store.ArchiveProject(id string) error` muda status para `archived`

**Critério de aceite:** Sistema iniciado sem projeto → card de onboarding criado com `status=blocked`; humano responde → projeto persistido na tabela `projects` com UUID; segundo início usa `project_id` do projeto criado sem novo onboarding.

---

### [P2] E0-T3 — Seletor de projeto ativo na API e no dashboard
**O que fazer:**
- Endpoint `GET /api/projects` → lista todos os projetos com status, slug, nome e contagem de cards
- Endpoint `POST /api/projects/:id/activate` → muda projeto ativo no runtime (sem reiniciar o servidor)
- Endpoint `POST /api/projects` → inicia onboarding de novo projeto (cria card de formulário)
- Endpoint `PATCH /api/projects/:id` → atualiza campos do projeto (name, description, stack, status)
- No dashboard (E4-T2): seletor de projeto no header — dropdown com projetos ativos; trocar de projeto filtra todos os cards, tokens e contexto exibidos
- Projeto `archived` aparece em seção colapsável "Arquivados", somente leitura

**Critério de aceite:** Dois projetos no banco; trocar projeto no dashboard filtra todos os cards corretamente; `POST /api/projects` com projeto B iniciado enquanto A está ativo cria card de onboarding para B sem interromper A.

---

### [P2] E0-T4 — Formulário de onboarding no dashboard (E4 integração)
**O que fazer:**
- Dependência de E4-T2 (dashboard) e E0-T2 (card de onboarding)
- Quando o dashboard detectar card com `comment_type=form` e `status=blocked`, renderizar o JSON de perguntas como formulário visual em vez de texto bruto
- Campos `type=text` → input; `type=textarea` → textarea; `type=select` → dropdown com as opções
- Campos `required=true` validados antes de enviar
- Ao submeter, gera resposta no formato esperado pelo CEO e faz `PUT /api/cards/:id/respond` com o payload preenchido
- O backend processa a resposta como se o humano tivesse renomeado o handoff manualmente

**Critério de aceite:** Card de onboarding exibido como formulário preenchível no dashboard; submissão persiste projeto e fecha o card com `status=done`.

---

## Épico 1 — Persistência CQRS (Quebra de .md em Registros)

> **Contexto:** Hoje o `importer.go` persiste documentos como blocos de texto em `document_sections`. A meta é que cada item significativo de cada arquivo `.md` vire um registro independente no banco, com rastreabilidade por timestamp, agente e task_ref — seguindo o padrão append-only já adotado em `context_entries`.

### [P1] E1-T1 — CQRS: Parser granular para arquivos de planejamento
**Arquivos-alvo:** `PLAN.md`, `ROADMAP.md`, `TASKS.md`
**O que fazer:**
- Estender `splitDocumentEntries` em `importer.go` para extrair cada item como registro próprio:
  - `TASKS.md`: cada task (`### Task [ID]`) vira um `planning_item` com `status`, `assigned_to`, `priority`, `deadline`
  - `ROADMAP.md`: cada checkbox (`- [ ]` / `- [x]`) vira um `planning_item` com `status=todo|done`
  - `PLAN.md`: cada seção numerada (critérios, regras, riscos) vira um `context_entry` com `section` e `entry_type`
- Campos obrigatórios em cada registro: `timestamp`, `source_agent`, `task_ref`, `document_type`
- Nenhum registro é deletado; atualizações geram novo registro (append-only)

**Critério de aceite:** `SELECT COUNT(*) FROM planning_items` reflete exatamente o número de tasks/checkboxes dos arquivos; consulta por `task_ref` retorna histórico completo de mudanças de status.

---

### [P1] E1-T2 — CQRS: Parser granular para arquivos de estado e contexto
**Arquivos-alvo:** `STATE.md`, `CONTEXT.md`, `RETROSPECTIVE.md`, `VERSIONS.md`
**O que fazer:**
- `STATE.md`: cada bloco `- **[timestamp] — [AGENT]:**` vira um `context_entry` com `entry_type=delivery|blocker|metric`
- `CONTEXT.md`: cada decisão `### [timestamp] — [Agent]: título` vira um `context_entry` com `entry_type=decision`
- `RETROSPECTIVE.md`: cada bloco de fase vira um `context_entry` com `entry_type=retrospective`
- `VERSIONS.md`: cada entrada `### [Version]` vira um `planning_item` com `item_type=release|rollback|bugfix`
- Reutilizar a função `startsBulletEntry` já existente como base, estendendo os padrões reconhecidos

**Critério de aceite:** Consulta `QueryContextForHandoff("CEO", "PHASE-1", 10)` retorna decisões e entregas da fase sem precisar reler os `.md`.

---

### [P1] E1-T3 — CQRS: Parser granular para arquivos de qualidade
**Arquivos-alvo:** `TESTS.md`, `QUESTIONS.md`, `PLAYBOOK.md`
**O que fazer:**
- `TESTS.md`: cada bloco `### Test [ID]` vira um `test_record` com `status=passed|failed|pending`
- `QUESTIONS.md`: cada entrada `### [timestamp] Question:` vira um `context_entry` com `entry_type=question|answer`, status `open|answered|blocked`
- `PLAYBOOK.md`: já funciona via `splitPlaybookEntries`; revisar para capturar o campo `Rationale` como metadata
- Garantir que `SaveTestRecord` e `SaveContextEntry` são usados (não criar novas tabelas)

**Critério de aceite:** `record_test` via handoff de agente e importação de `TESTS.md` produzem registros na mesma tabela `test_records`; sem duplicatas por hash.

---

### [P2] E1-T4 — Deduplicação por hash no importer
**O que fazer:**
- Antes de inserir qualquer seção ou entrada, verificar se `content_hash` já existe para o `document_id` + `import_id` mais recente
- Se o conteúdo não mudou, pular a inserção (não criar registro duplicado)
- Logar quantos registros foram pulados vs inseridos em cada importação
- Manter a regra: se mudou, registrar novo (append); se não mudou, não duplicar

**Critério de aceite:** Segunda importação de um arquivo inalterado produz 0 novos registros em `document_sections` e `document_entries`.

---

### [P2] E1-T5 — Cadastro de modelos LLM por agente em persistência
**O que fazer:**
- Criar tabela `agent_configs` no SQLite:
  ```sql
  CREATE TABLE agent_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    agent_name TEXT NOT NULL,
    model TEXT NOT NULL,
    tier INTEGER NOT NULL,
    active INTEGER NOT NULL DEFAULT 1,
    updated_at TEXT NOT NULL,
    updated_by TEXT
  );
  ```
- Migrar o mapa hardcoded `agentModelMap` em `agent.go` para leitura desta tabela via `Store`
- `ResolverModeloAgente` passa a consultar banco; fallback para o mapa hardcoded se tabela vazia
- Seed inicial: inserir os mapeamentos atuais como registros na migration `004_agent_configs.sql`
- Expor métodos no `Store`: `GetAgentConfig(name)`, `UpsertAgentConfig(config)`, `ListAgentConfigs()`

**Critério de aceite:** Alterar o modelo do agente CEO via `UpsertAgentConfig` sem recompilar o binário; mudança refletida na próxima chamada a `ChamarAgente`.

---

## Épico 2 — Cards e Rastreabilidade de Handoffs

> **Contexto:** Hoje handoffs são arquivos JSON em `.agent_handoff/`. A meta é que cada handoff funcione como um **card** — com ID único, histórico de comentários (respostas encadeadas), status e rastreabilidade. O filesystem continua sendo o transporte; o banco persiste o estado do card.

### [P1] E2-T1 — Modelo de Card no banco
**O que fazer:**
- Criar tabela `cards` na migration `005_cards.sql`:
  ```sql
  CREATE TABLE cards (
    id TEXT PRIMARY KEY,           -- UUID gerado no handoff
    project_id TEXT NOT NULL REFERENCES projects(id),
    title TEXT NOT NULL,
    task_ref TEXT,
    sender TEXT NOT NULL,
    recipient TEXT NOT NULL,
    intent TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'open',  -- open | in_progress | blocked | done | failed
    priority TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    max_retries INTEGER NOT NULL DEFAULT 3,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
  );

  CREATE TABLE card_comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    card_id TEXT NOT NULL REFERENCES cards(id),
    author TEXT NOT NULL,          -- agente ou [HUMAN]
    content TEXT NOT NULL,         -- JSON do payload ou texto livre
    comment_type TEXT NOT NULL,    -- handoff | response | note | error | token_usage | form
    created_at TEXT NOT NULL
  );

  CREATE INDEX IF NOT EXISTS idx_cards_project ON cards(project_id, status);
  CREATE INDEX IF NOT EXISTS idx_cards_task_ref ON cards(task_ref);
  CREATE INDEX IF NOT EXISTS idx_card_comments_card ON card_comments(card_id);
  ```
- Adicionar campo `card_id` (UUID v4) no `HandoffHeader` em `hand_off.go`
- `SaveRawHandoff` deve criar ou atualizar o card correspondente e inserir um comentário com o payload
- `retry_count` e `max_retries` movidos para cá (removidos de E2-T3 que apenas implementa a lógica de uso)

**Critério de aceite:** Todo handoff processado tem registro em `cards` com `project_id` correto + ao menos um registro em `card_comments`; consulta por `card_id` retorna thread completa; cards de projetos diferentes não se misturam.

---

### [P1] E2-T2 — Encadeamento de respostas como comentários
**O que fazer:**
- Quando um agente responde a um handoff existente, o JSON de resposta deve incluir `"reply_to_card_id": "..."` 
- `applyAgentResponse` em `pipeline.go`: se `reply_to_card_id` presente, inserir comentário no card existente (não criar card novo) e atualizar `status`
- Handoffs `_TO_HUMAN_` criam card com `status=blocked`; quando o humano renomeia para `_TO_AGENTE_`, o pipeline insere comentário com a resposta e muda status para `in_progress`
- Fluxo de status: `open → in_progress → blocked (ask_human) → in_progress → done | failed`

**Critério de aceite:** Thread de 3 trocas (CEO → DEV → ask_human → CEO) aparece como 4 comentários em um único card, não 4 cards separados.

---

### [P2] E2-T3 — Retry automático com contador por card
**O que fazer:**
- As colunas `retry_count` e `max_retries` já existem na tabela `cards` (criadas em E2-T1)
- Em `handleHandoff`: antes de mover para `failed/`, consultar `retry_count < max_retries` no banco
  - Se sim: incrementar `retry_count`, mover arquivo de volta para `.agent_handoff/` (re-enfileirar), aguardar backoff exponencial (1s, 2s, 4s)
  - Se não: mover para `failed/`, atualizar `status=failed`, registrar comentário com o erro detalhado
- Backoff implementado com `time.Sleep` na goroutine do worker (não bloqueia outros workers)
- Registrar cada tentativa como comentário `comment_type=error` no card com stack trace resumido

**Critério de aceite:** Falha simulada de agente resulta em 3 tentativas logadas em `card_comments` antes de `status=failed`; sem falha silenciosa; backoff visível no timestamp dos comentários.

---

### [P2] E2-T4 — Contexto comprimido de gerenciais para técnicos
**O que fazer:**
- Agentes gerenciais (CEO, CTO, BA) ao gerar handoffs para agentes técnicos devem incluir apenas o contexto necessário para aquela task, não o contexto completo da sessão
- Em `enrichHandoffWithContext`: criar filtro por `document_type` e `task_ref` — agentes técnicos recebem máximo 5 itens de contexto diretamente relacionados à task, não os 8 globais atuais
- CEO/CTO/BA recebem contexto mais amplo (até 20 itens) para decisões arquiteturais
- Adicionar campo `context_scope` no `PHASE_KICKOFF` payload: `"task"` (técnicos) | `"phase"` (gerenciais)

**Critério de aceite:** Handoff para `DEV_BACKEND` contém ≤ 5 itens de contexto; handoff para `CEO` contém até 20; logs mostram contagem por handoff.

---

### [P2] E2-T5 — Contador de tokens usados por chamada LLM
**O que fazer:**
- Criar tabela `token_usage` na migration `006_token_usage.sql`:
  ```sql
  CREATE TABLE token_usage (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id TEXT NOT NULL REFERENCES projects(id),  -- isolamento por projeto
    card_id TEXT REFERENCES cards(id),     -- card que originou a chamada
    agent_name TEXT NOT NULL,              -- agente que chamou
    model TEXT NOT NULL,                   -- modelo usado
    tier INTEGER NOT NULL,                 -- 1, 2 ou 3
    prompt_tokens INTEGER NOT NULL,        -- tokens de entrada (prompt_eval_count)
    completion_tokens INTEGER NOT NULL,    -- tokens de saída (eval_count)
    total_tokens INTEGER NOT NULL,         -- soma dos dois
    latency_ms INTEGER,                    -- duração da chamada em ms
    created_at TEXT NOT NULL
  );

  CREATE INDEX IF NOT EXISTS idx_token_usage_project ON token_usage(project_id, created_at);
  CREATE INDEX IF NOT EXISTS idx_token_usage_card ON token_usage(card_id);
  CREATE INDEX IF NOT EXISTS idx_token_usage_agent ON token_usage(agent_name, created_at);
  CREATE INDEX IF NOT EXISTS idx_token_usage_model ON token_usage(model, created_at);
  ```
- Capturar `prompt_eval_count` e `eval_count` do chunk final (`done: true`) do streaming em `callChat` no `agent.go`:
  - Estender `streamChunk` para incluir os campos: `PromptEvalCount int`, `EvalCount int`
  - Após o loop de streaming, chamar `Store.SaveTokenUsage(usage)` com os valores coletados
- Registrar `latency_ms` medindo `time.Since(inicio)` antes e depois da chamada HTTP
- Expor método `Store.SumTokenUsage(filters)` para consultas analíticas: por agente, por modelo, por período, por card
- Após salvar no banco, inserir comentário no card associado com resumo legível:
  ```
  [TOKEN_USAGE] prompt=1240 completion=380 total=1620 model=deepseek-r1:7b tier=3 latency=4.2s
  ```
  Comentário do tipo `token_usage` (novo `comment_type` em `card_comments`)

**Endpoints a adicionar em E4-T1:**
```
GET /api/tokens/summary?agent=&model=&from=&to=   → totais agregados
GET /api/tokens/by-card/:card_id                  → uso detalhado por card
```

**Consultas analíticas habilitadas após esta task:**
- Custo total por agente no período
- Modelo com maior consumo
- Custo estimado por card/feature (base para priorização)
- Comparação de eficiência antes/depois de otimizações de prompt (E3-T3)

**Critério de aceite:** Toda chamada a `callChat` ou `ChamarAgenteComCodexCLI` gera um registro em `token_usage`; `GET /api/tokens/summary?agent=CEO` retorna total de tokens consumidos pelo CEO; comentário `[TOKEN_USAGE]` visível na thread do card correspondente.

---

## Épico 3 — Economia de Tokens

> **Contexto:** As tasks deste épico são economizadores — algumas já existem parcialmente, outras precisam ser implementadas do zero. Todas têm como métrica principal: redução de tokens por ciclo de agente.

### [P1] E3-T1 — LLM Tiering formalizado (Caveman Mode)
**O que fazer:**
- O mapa `agentModelMap` em `agent.go` já implementa tiering; formalizar com os 3 níveis explícitos:
  - **Tier 1 (Efficiency):** `gemma4:latest` — tarefas atômicas: formatação de logs, linting, seed de dados
  - **Tier 2 (Development):** `qwen2.5-coder` — implementação padrão TDD, features, migrações
  - **Tier 3 (Expert):** `deepseek-r1:7b` — decisões arquiteturais, resolução de ambiguidades, orchestration
- Após E1-T5: o tiering passa a ser lido do banco `agent_configs`
- Adicionar log estruturado por chamada: `[AGENT:CEO] model=deepseek-r1:7b tier=3 tokens_estimate=~2400`
- Regra: agentes gerenciais nunca chamam Tier 1; agentes técnicos de tarefas mecânicas nunca chamam Tier 3

**Critério de aceite:** Log de cada chamada de agente registra model e tier; `agent_configs` no banco reflete os 3 tiers.

---

### [P2] E3-T2 — Memória seletiva: injeção de contexto por relevância
**O que fazer:**
- `enrichHandoffWithContext` hoje busca top-N por `projected_at DESC` (mais recente) — mudar para ranking por relevância:
  - Prioridade 1: itens com `task_ref` exatamente igual ao do handoff atual
  - Prioridade 2: itens com mesmo `document_type` do handoff
  - Prioridade 3: itens com `agent_name` igual ao `recipient` do handoff
  - Prioridade 4: demais, ordenados por recência
- Implementar como query SQL com `CASE WHEN` para score de relevância
- Limite por tier: Tier 1 recebe 0 itens de contexto (tarefa mecânica, não precisa); Tier 2 recebe até 5; Tier 3 recebe até 20

**Critério de aceite:** Agente `QA` recebendo handoff para `TASK-42` recebe contexto de `test_records` e `planning_items` de `TASK-42` antes de receber contexto genérico.

---

### [P2] E3-T3 — Compressão de prompts de sistema por agente
**O que fazer:**
- Hoje `carregarArquivoAgente` carrega o `.md` completo do agente como system prompt (~2-4KB por agente)
- Criar função `comprimirPromptAgente(content string) string` que:
  - Remove seções `## Version History` e `## Allowed Documents` do prompt (metadata, não instrução)
  - Remove linhas em branco consecutivas (mais de 2)
  - Substitui exemplos de JSON longos por referência: `[ver ARCHITECTURE.md §4.2.X]`
- Redução esperada: 30-40% do tamanho do system prompt
- Aplicar apenas em Tier 1 e Tier 2; Tier 3 recebe prompt completo

**Critério de aceite:** Prompt do agente `QA` (Tier 1) comprimido ≤ 1.5KB; log registra tamanho antes/depois.

---

### [P3] E3-T4 — Índice de símbolos de código (Token Savior)
**O que fazer:**
- Quando `write_code` persiste um arquivo, extrair símbolos públicos (funções, tipos, constantes) via parsing simples de regex
- Armazenar em nova tabela `code_symbols`: `(file_path, symbol_name, symbol_type, signature, created_at)`
- Agentes técnicos recebem índice de símbolos dos arquivos relevantes à task (não o arquivo completo)
- Evita que o agente peça para "reler" arquivos que já foram escritos nesta sessão

**Critério de aceite:** Handoff para `DEV_BACKEND` referenciando `src/agent/agent.go` recebe lista de símbolos exportados, não o conteúdo completo do arquivo.

---

### [P3] E3-T5 — MCP como protocolo de transporte (migração de filesystem)
**O que fazer:**
- Hoje o transporte é filesystem JSON em `.agent_handoff/`; o objetivo é MCP (Model Context Protocol)
- Fase preparatória (não quebra o atual):
  - Implementar servidor MCP em Go expondo as ferramentas existentes: `submit_handoff`, `get_card_status`, `query_context`, `list_agent_configs`
  - O filesystem continua funcionando em paralelo (fallback)
- Ganho esperado: latência de comunicação entre agentes cai de polling de arquivo para chamada direta; sem race condition de `moveFile`
- Dependência: E2-T1 (cards) deve estar completo antes — MCP opera sobre cards, não arquivos

**Critério de aceite:** Handoff CEO → DEV_BACKEND via MCP processado sem criação de arquivo intermediário; tempo de ciclo comparado em log.

---

## Épico 4 — API REST e Frontend

> **Contexto:** Não há servidor HTTP hoje. O dashboard pressupõe uma API — esta épica implementa a API primeiro, o frontend depois.

### [P1] E4-T1 — Servidor HTTP em Go (API REST)
**O que fazer:**
- Adicionar servidor `net/http` em `main.go` (sem framework externo — stdlib é suficiente)
- Endpoints mínimos para o dashboard funcionar:
  ```
  GET  /api/projects                    → lista projetos (active | archived)
  POST /api/projects                    → inicia onboarding de novo projeto
  POST /api/projects/:id/activate       → muda projeto ativo
  PATCH /api/projects/:id               → atualiza campos do projeto
  GET  /api/cards?project_id=&status=   → lista cards com filtros
  GET  /api/cards/:id                   → card + thread de comentários
  PUT  /api/cards/:id/respond           → humano responde card bloqueado
  GET  /api/agents                      → lista agent_configs
  PUT  /api/agents/:name                → atualiza model/tier de um agente
  GET  /api/tokens/summary?project_id=&agent=&model=&from=&to=  → totais agregados
  GET  /api/tokens/by-card/:card_id     → uso detalhado por card
  GET  /api/context?project_id=&task_ref=&limit=  → retorna read_context_items
  GET  /api/health                      → status do sistema (ollama, workers, projeto ativo)
  ```
- Autenticação: token estático via env var `API_TOKEN` (header `Authorization: Bearer`); sem OAuth no MVP
- CORS configurado para `localhost:3000` (React dev server)
- Porta via env var `HTTP_PORT` (default `8080`)

**Critério de aceite:** `curl http://localhost:8080/api/health` retorna JSON com status dos componentes; `curl /api/cards` retorna lista paginada.

---

### [P2] E4-T2 — Dashboard React: visão de cards
**O que fazer:**
- Header com seletor de projeto ativo (dropdown com projetos `active`); trocar projeto filtra tudo
- Interface para visualizar cards como quadro Kanban (colunas: open, in_progress, blocked, done, failed)
- Cada card mostra: título, agente responsável, task_ref, timestamp, retry_count
- Card com `comment_type=form` renderizado como formulário preenchível (E0-T4)
- Clicar no card abre thread de comentários (histórico completo de handoffs encadeados)
- Filtros: por agente, por status, por task_ref, por data
- Painel lateral "Tokens" com totais do projeto ativo: por agente, por modelo, por dia
- Atualização em tempo real via polling longo (`GET /api/cards` a cada 5s) ou SSE
- Design: minimalista técnico — inspiração Linear/Vercel

**Critério de aceite:** Dois projetos no banco; trocar no seletor filtra cards e tokens corretamente; card bloqueado com formulário exibe campos preenchíveis; thread de 4 comentários exibida corretamente.

---

### [P2] E4-T3 — Modal de cadastro de modelos por agente
**O que fazer:**
- No dashboard, botão "Agents" abre modal com tabela de todos os agentes e seus modelos atuais
- Cada linha: `agent_name | model | tier | active | updated_at | ação: editar`
- Editar abre inline form: dropdown de modelos disponíveis (buscado de `GET /api/ollama/models` → `ollama list`) + seleção de tier
- Salvar chama `PUT /api/agents/:name` e atualiza o banco via E1-T5
- Sem reload do servidor — mudança refletida na próxima chamada ao agente

**Critério de aceite:** Alterar modelo do `QA` de `gemma4:latest` para `qwen2.5-coder` via modal e confirmar que próximo handoff para QA usa o novo modelo (verificável via log).

---

### [P3] E4-T4 — Health check contínuo do Ollama no runtime
**O que fazer:**
- Goroutine de health check a cada 30s: `GET {OLLAMA_URL}/api/tags`
- Se falhar: marcar flag global `ollamaHealthy=false`, logar `WARN`, não processar novos handoffs de agentes (enfileirar)
- Se recuperar: marcar `ollamaHealthy=true`, logar `INFO`, retomar fila
- Endpoint `/api/health` expõe o status atual
- Hoje o `run.bat` checa na inicialização, mas não durante a execução

**Critério de aceite:** Parar o Ollama durante execução; novos handoffs ficam em fila (não vão para `failed/`); reiniciar Ollama faz a fila ser processada automaticamente.

---

## Épico 5 — Estabilidade e Correções

> Tasks de baixo custo e alto valor operacional. Devem ser feitas antes dos épicos 2-4 para não acumular dívida técnica.

### [P1] E5-T1 — Lock em moveFile (race condition de workers)
**O que fazer:**
- Adicionar `sync.Mutex` no `Pipeline` específico para operações de filesystem
- `moveFile` e `handleHandoff` adquirem o lock antes de verificar existência + mover
- Alternativa mais simples: usar `os.Rename` atomicamente + verificar `os.IsNotExist` no erro (já feito parcialmente); garantir que o fallback de timestamp no nome seja suficiente
- Adicionar teste unitário simulando 2 workers movendo o mesmo arquivo simultaneamente

**Critério de aceite:** `go test -race ./...` passa sem data race detectada em `pipeline.go`.

---

### [P1] E5-T2 — Corrigir dependência `go-shellquote` no go.mod
**O que fazer:**
- `go-shellquote` é usado em `agent.go` mas não declarado como dependência direta em `go.mod`
- Executar `go mod tidy` e commitar `go.mod` e `go.sum` atualizados
- Verificar se há outras dependências indiretas usadas diretamente (rodar `go mod why` para cada import)

**Critério de aceite:** `go mod tidy` não altera `go.mod`; `go build ./...` passa sem warnings.

---

### [P2] E5-T3 — Cron scheduler para routines de agentes
**O que fazer:**
- Goroutine com ticker configurável via env vars: `CRON_ENABLED=true`, `CRON_INTERVAL_MINUTES=60`
- A cada tick, gerar automaticamente handoffs para agentes de manutenção: backup do banco, relatório de cards blocked há mais de N horas, health report do sistema
- Implementar somente quando houver caso de uso concreto definido (não implementar "por garantia")
- Pré-requisito: E4-T1 (API) pronto para expor status ao cron

**Critério de aceite:** Task definida, implementação adiada até E4 estar completo.

---

## Resumo de Prioridades

| Épica | Task | Prioridade | Depende de | Esforço est. |
|:------|:-----|:----------:|:----------:|:------------:|
| E0 | T1 — Tabela projects + project_id em todas as tabelas | P1 | — | M |
| E0 | T2 — Onboarding via card estruturado (CEO) | P1 | E0-T1 | M |
| E1 | T1 — Parser CQRS planejamento | P1 | E0-T1 | M |
| E1 | T2 — Parser CQRS estado/contexto | P1 | E0-T1 | M |
| E1 | T3 — Parser CQRS qualidade | P1 | E0-T1 | S |
| E5 | T1 — Lock race condition | P1 | — | S |
| E5 | T2 — Fix go.mod | P1 | — | XS |
| E1 | T4 — Deduplicação por hash | P2 | E1-T1,T2,T3 | S |
| E1 | T5 — Cadastro de modelos no banco | P2 | E0-T1 | M |
| E2 | T1 — Modelo de Card | P1 | E0-T1 | M |
| E2 | T2 — Encadeamento de respostas | P1 | E2-T1 | M |
| E2 | T3 — Retry com contador | P2 | E2-T1 | S |
| E2 | T4 — Contexto comprimido gerencial→técnico | P2 | E1-T2 | S |
| E2 | T5 — Contador de tokens por chamada LLM | P2 | E2-T1 | M |
| E3 | T1 — LLM Tiering formalizado | P1 | E1-T5 | S |
| E3 | T2 — Memória seletiva por relevância | P2 | E1-T1,T2 | M |
| E3 | T3 — Compressão de prompts | P2 | — | S |
| E4 | T1 — API REST Go | P1 | E0-T1, E2-T1 | M |
| E4 | T2 — Dashboard React: Kanban + projetos | P2 | E4-T1 | L |
| E4 | T3 — Modal de modelos | P2 | E4-T2, E1-T5 | M |
| E4 | T4 — Health check Ollama runtime | P2 | E4-T1 | S |
| E0 | T3 — Seletor de projeto ativo na API e dashboard | P2 | E4-T2, E0-T1 | M |
| E0 | T4 — Formulário de onboarding no dashboard | P2 | E4-T2, E0-T2 | M |
| E3 | T4 — Índice de símbolos (Token Savior) | P3 | E2-T1 | L |
| E3 | T5 — MCP como transporte | P3 | E2-T1, E4-T1 | XL |
| E5 | T3 — Cron scheduler | P3 | E4-T1 | M |

---

## Sequência de Implementação Sugerida

```
Semana 1: E5-T2 → E5-T1 → E0-T1 → E0-T2
Semana 2: E1-T1 → E1-T2 → E1-T3 → E1-T4
Semana 3: E1-T5 → E2-T1 → E2-T2 → E2-T3
Semana 4: E3-T1 → E3-T3 → E2-T4 → E2-T5
Semana 5: E4-T1 → E4-T4 → E3-T2
Semana 6: E4-T2 → E4-T3 → E0-T3 → E0-T4
Fase 2:   E3-T4 → E3-T5 → E5-T3
```

---

*Backlog gerado com base em: análise do código-fonte (Go, migrations SQL, agent configs), diretivas do [HUMAN] em 2026-05-13, e cruzamento com o documento de arquitetura do sistema. Última atualização: 2026-05-13 — adicionados Épico 0 (projetos e onboarding), E2-T5 (contador de tokens), propagação de `project_id` em todas as tabelas.*
