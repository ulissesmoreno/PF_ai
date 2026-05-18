# PF AI Orchestration System

Sistema local de orquestracao de agentes com persistencia em SQLite, cards rastreaveis e execucao via Ollama.

## Como Funciona

- Projetos ficam na tabela `projects`.
- Handoffs viram cards em `cards` e historico em `card_comments`.
- Contexto, plano, estado, testes e decisoes sao gravados por acoes CQRS no banco.
- Perguntas ao humano viram handoffs/cards `blocked` e podem ser respondidas pela API.
- Arquivos Markdown operacionais antigos sao apenas legado/importacao, nao sao mais o fluxo principal.

## Rodar

```powershell
cd src
go run .
```

API padrao:

- `GET http://127.0.0.1:8080/api/health`
- `GET http://127.0.0.1:8080/api/projects`
- `POST http://127.0.0.1:8080/api/projects`
- `POST http://127.0.0.1:8080/api/projects/{id}/activate`
- `GET http://127.0.0.1:8080/api/cards`
- `GET http://127.0.0.1:8080/api/cards/{id}`
- `PUT http://127.0.0.1:8080/api/cards/{id}/respond`

## Regras Atuais Para Agentes

- Nao escrever contexto/plano/estado/testes em arquivos operacionais.
- Usar `update_context`, `update_plan`, `update_state`, `record_test`, `record_decision` e `record_retrospective`.
- Usar `ask_human` quando precisar de decisao humana; nao preencher `QUESTIONS.md`.
- Usar `write_code` somente para codigo, wiki ou artefatos reais do workspace.

## Testes

```powershell
cd src
go test ./...
```
