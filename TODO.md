# TODO

- [x] O OK no roadmap é dado pelo usuário. Antes da próxima etapa a anterior deve ser marcada como concluída.
- [x] .env não deve ser commitado nunca, assim como nenhuma informação que comprometa a segurança da aplicação
- [x] Na arquitetura indicar possibilidade de outros frontends, como aplicativos e destacar suas tecnologias
- [x] Passo-a-passo do teste que o dev vai fazer além do que já foi feito
- [x] Criar um questionário para iniciar um projeto novo, com preenchimento automático antes começar 
- [x] Antes de iniciar uma nova etapa pedir o preenchimento dos arquivos que indicam etapa concluída
- [x] Criar um arquivo onde preferências são descritas, esse deve ser agnóstico de projeto
- [x] Criar uma documentação para manutenção do projeto. Para isso criei um arquivo WIKI.md que contém a susgestão de usar o Obsidian para fazer a wiki. Criar um diretório wiki dentro do projeto e seguir a estrtura dessa arquivo. Revisar o mesmo para garantir o padrão com o que estamos trabalhando aqui. Criado em Inglês em ambas as pastas, fazer a tradução. Como a sugestão é usar o Obsidian, adicionar em ENV_SETUP.md
- [x] Atualizei o arquivo GSD-RULES em PT, reflita em EN
- [x] Foi criado um aruqivo DESIGN.md, nele será contido os detalhes estéticos do projeto. Pergunte ideias ou projetos similares que gostaria que fosse mantidos, como cores, fontes, layouts, animações, etc... Em caso de pedir para imitar o projeto da empresa X, ou do link / aplicativo Y, esse arquivo deverá conter apenas os informações necessárias para a contrução do visual sem nenhuma referência a inspiração (que pode estar no PROJECT.md)
- [x] Criar uma sessão (README ou PROJECT) onde fontes, inspirações, referências, etc... devem ser adicionadas. 
- [x] Testes de vulnerabilidades devem ser realizados antes da conclusão de cada etapa do projeto, sugerir ferramentas, padrões e etc
- [x] O código deve ser otimizado para performance e segurança.
- [x] Como é possível economizar tokens do agente durante a execução do desenvolvimento. 

> **Status:** Todas as pendências iniciais do Boilerplate foram validadas, integradas e concluídas (ver `VERSIONS.md` v1.1.0).

- [x] Até segunda preencher apenas documentos em inglês
- [x] Baseado em https://github.com/paperclipai/paperclip criar agentes que faram parte da atividade e terão acesso limitado a documentos. Criar um CEO que fará a orquestração e distribuição de tarefas conhecendo a regra do projeto, um CTO que conhece a arquitetura e suas diretrizes, um dev que desenvolverá o código e um QA que fará os testes. Criar arquivos para cada um (CEO.md, CTO.md, DEV.md, QA.md)

- [x] ENV_SETUP.md: Na seção de ferramentas ou variáveis de ambiente, podem ser definidas as chaves de API e os modelos específicos para cada função
- [x] ARCHITECTURE.md: Na tabela de Stack Tecnológico, pode-se adicionar uma coluna para "AI Engine" por agente.
- [x] PROJECT.md: adiconar as regras para criação de novos papeis / agentes (como Dev Backend, Dev Frontend, níveis de senioridade ou CMO). Se for necessário criar um novo agente, deve-se criar um arquivo .md com as regras específicas dele. NEW-INSTRUCTIONS.md será usado caso seja necessário criar um novo papel / agente após o inicio do projeto.
- [x] Transoforme o agente DEV em Dev Frontend e Dev Backend. Em caso de novas tecnologias serem adicionadas ao projeto, o CTO deve ser consultado para definir como será o gerenciamento das tecnologias.  
- [x] Altere o promtp inicial para "[SISTEMA DE ATIVAÇÃO GSD - OPERAÇÃO EM MODO EXPERTO]

Você é o Agente de Orquestração deste ecossistema. Sua missão é garantir a execução do framework Get Shit Done (GSD), priorizando rastreabilidade, consistência e qualidade técnica sobre a velocidade.

PROTOCOLOS OBRIGATÓRIOS DE INÍCIO:
1. LEITURA DE DIRETRIZES: Antes de qualquer ação ou sugestão, leia obrigatoriamente GSD-RULES.md (regras invioláveis) e PLAYBOOK.md (preferências do desenvolvedor).
2. VERIFICAÇÃO DE ESTADO: Consulte o arquivo PROJECT.md. Se ele contiver placeholders como "[Fill with...]", você deve pausar qualquer desenvolvimento e disparar imediatamente o protocolo de ONBOARDING.md (Blocos 1 a 5) para definir o projeto com o usuário.
3. IDENTIDADE DE AGENTE: Identifique qual papel você está assumindo para a tarefa atual (CEO, CTO, DEV ou QA) com base nas responsabilidades definidas em seus respectivos arquivos .md. Assine cada interação ou alteração de arquivo com sua tag (ex: [CEO], [DEV]).
4. ORQUESTRAÇÃO E HANDOFF: Para transições de tarefas entre agentes, utilize estritamente o diretório .agent_handoff/ seguindo o padrão de comunicação programática (JSON/YAML) definido.
5. FLUXO OPERACIONAL: 
   - Leia NEW-INSTRUCTIONS.md e CONTEXT.md para entender a demanda e o histórico.
   - Registre qualquer ambiguidade em QUESTIONS.md; não prossiga com dúvidas pendentes.
   - Siga rigorosamente o ciclo TDD e a Arquitetura Hexagonal descritos em ARCHITECTURE.md.

DIRETRIZES DE SAÍDA:
- Mantenha a documentação bilíngue (PT/EN) conforme o padrão do boilerplate.
- Nunca delete o histórico; utilize apenas seções de histórico imutável com timestamps (YYYY-MM-DD HH:MM).
- Garanta que nenhum arquivo permaneça com placeholders após a conclusão de uma etapa (Definition of Done).

AGUARDANDO COMANDO INICIAL OU DEFINIÇÃO DE PROJETO."
- [x] Adicione o TODO.md ao .gitignore
- [x] Atualize links quebrados no README.md e outros .md. Deletei e movi arquivos.
- Primeiramente será listada todas as atividades aqui, organizadas em grupos semelhantes. Decidirei a ordem de realização.
- [x] Precedência do CEO: Formalizado em GSD-RULES.md §0.1 e CEO.md — o CEO é invariavelmente o primeiro a agir em qualquer ciclo, com leitura obrigatória de GSD-RULES.md → PLAYBOOK.md → NEW-INSTRUCTIONS.md → PLAN.md antes de qualquer delegação.
- [x] CEO lê ONBOARDING.md no inicio do projeto — formalizado em GSD-RULES.md §0.1 (passo 5) e no CEO.md (Responsibilities + Allowed Documents).
- [x] Handoff Gerencial: Hierarquia CEO → Camada Gerencial (CTO/BA) → Agentes Técnicos formalizada em GSD-RULES.md §0 (seções 0.2, 0.3, 0.4, 0.5). GSD-RULES.md adicionado como leitura obrigatória em CTO.md e BA.md.
- Qualquer agente pode interromper uma atividade para registrar dúvidas, essa diretriz deve ser consolidada no "coração" das normas do projeto. Esta regra reforça o princípio de que "nada começa com perguntas pendentes" e garante a rastreabilidade total exigida pelo framework GSD.
- [x] Movimentação de arquivos:
   - [x] NEW-INSTRUCTIONS e QUESTIONS mover para raiz.
   - [x] ONBOARDING mover para DOC
   - [x] Atualizar todos os arquivos que fazem referencia aos arquivos movimentados.
- [x] Análise de Skills (CTO): Formalizado em GSD-RULES.md §9.3.6 — [CTO] é o único responsável pelo Antigravity Skill Mapping antes de cada fase do roadmap, comunicando via .agent_handoff/ e registrando em CONTEXT.md. Também refletido em CTO.md (Responsibilities).
- Cobertura de Arquivos: Implementar uma regra no STATE.md que garanta que todos os arquivos do repositório sejam lidos e processados por pelo menos um agente em momentos específicos do roadmap, eliminando "zonas cegas" de contexto.
- [x] Remover citação a NEW_PROJECT_PROMPT.md
- [x] Papel Humano (User): [HUMAN] criado formalmente em PROJECT.md §5.1 como autoridade final de aprovação de Stage Closure Gates. Referenciado em GSD-RULES.md §1.8 (Gate Rule) e README.md (Ignition Prompt PT+EN). Inclui nota de personalização [HUMAN:Nome].
- Wiki Colaborativa: Definir na WIKI.md e no ENV_SETUP.md que BA e DEVs preenchem a Wiki em páginas distintas:
   - BA: Foca no contexto de negócio, manuais de usuário e tradução de valor.
   - DEVs: Registram a documentação técnica de manutenção e detalhes de implementação.
- To signal that human intervention is finished and the GSD Orchestration Cycle should resume, the user must send the following command strictly in the chat:
[USER_DONE]
   - Trigger: This prompt informs the [CEO] Agent that all manual updates to NEW-INSTRUCTIONS.md, QUESTIONS.md, or reviews of the ROADMAP.md are complete
- [x] Criação de Novos Perfis de Agentes (.md)
   - [x] Data Scientist / ML Engineer: `AGENTS/DS_ML.md` — Focado no componente de Intelligence (Python/Scikit-Learn). Responsável por modelagem de ML e garantia de processamento paralelo no hardware local (n_jobs=-1).
   - [x] Security Auditor: `AGENTS/SECURITY.md` — Focado em auditorias de segurança, conformidade com JWT, sanitização nos adapters e execução de ferramentas de SAST (SonarQube/Snyk) antes da conclusão de cada etapa.
   - [x] CMO (Chief Marketing Officer): `AGENTS/CMO.md` — Auxilia o CEO na definição de visão de mercado, público-alvo e proposição de valor no PROJECT.md.
   - [x] Redator: `AGENTS/WRITER.md` — escreve os textos de Marketing.
   - [x] Artista: `AGENTS/ARTIST.md` — cria artes para divulgação.
   - [x] Revisor de Código: `AGENTS/CODE_REVIEWER.md` — Especialista em qualidade técnica. Analisa a aderência aos padrões SOLID, Clean Code e resultados de ferramentas como SonarQube após cada entrega atômica e antes do teste final do usuário.
   - [x] DevOps / SRE Specialist: `AGENTS/DEVOPS.md` — Responsável por configurar os pipelines de CI/CD no GitHub Actions, orquestrar os containers via Docker Compose e garantir a resiliência do ambiente produtivo.
   - [x] Agente BA (Business Analyst): `AGENTS/BA.md` — Ponte entre as instruções do usuário e a tradução técnica, responsável por critérios BDD e manuais na WIKI.

- [x] Dormant Agent Rule: Formalizado em GSD-RULES.md §0.5, PROJECT.md §5.2 e README.md. Agentes sem domínio ativo permanecem adormecidos e não consomem contexto.
- Estratégias para cada agente:
   1. Data Scientist / ML Engineer (Camada de Inteligência)
      - Method: Vectorized ML Pipeline. Evitar loops nativos do Python em favor de operações vetorizadas com Pandas e NumPy, utilizando tipos de dados otimizados (como float32) para reduzir a pegada de memória em hardware local.
      - Strategy: Parallel Resource Exhaustion. Implementar obrigatoriamente n_jobs=-1 em todos os estimadores compatíveis para exaurir a capacidade da CPU local, garantindo que o processamento massivo não se torne um gargalo.
      - Traceability: Registrar métricas de precisão e recall no TESTS.md para cada versão de modelo, vinculando-as ao snapshot correspondente no STATE.md.
   2. Security Auditor (Guardião da Integridade)
      - Method: Shift-Left Security Auditing. Integrar a análise de segurança no início de cada estágio. Realizar scans de SAST (Static Application Security Testing) usando ferramentas como SonarQube ou Snyk antes do código ser enviado para o QA.
      - Strategy: Zero Trust & Sanitization. Auditar todos os Adapters de infraestrutura para garantir que os dados sejam sanitizados antes de chegarem ao Domínio. Verificar se todas as rotas Java nascem com "deny-all" e se o JWT está corretamente implementado
   3. CMO - Chief Marketing Officer (Visão de Mercado)
      - Method: Customer-Centric Validation. Validar o PROJECT.md e o README.md sob a perspectiva do Público-Alvo, garantindo que a proposta de valor e o objetivo do sistema estejam claros e atendam a uma dor real.
      - Strategy: KPI-Driven Roadmap Alignment. Monitorar as métricas de sucesso definidas na Seção 7 do PROJECT.md e sugerir ajustes no ROADMAP.md caso o MVP não esteja atingindo os indicadores de tração ou engajamento esperados.
   4. Revisor de Código (Especialista em Qualidade)
      - Method: Static & Formal Analysis. Revisar o código após cada entrega atômica do DEV, verificando a aderência aos padrões SOLID, Clean Code e, principalmente, o isolamento da Arquitetura Hexagonal (garantindo que o Domínio tenha zero dependências externas).
      - Strategy: Debt Control & Refactoring. Identificar e documentar "débitos técnicos" no QUESTIONS.md ou STATE.md. Sua estratégia é garantir que o ciclo Red-Green-Refactor do TDD seja seguido rigorosamente, impedindo que "código de gambiarra" avance para o Stage Closure Gate.
   5. Agente BA (Business Analyst - Requisitos)
      - Método: BDD (Behavior Driven Development). Escrever os Critérios de Aceitação no PLAN.md §4 em um formato que o QA possa transformar diretamente em testes (ex: "Dado que... Quando... Então...").
      - Estratégia de Documentação: User Value Translation. Garantir que cada tarefa no TASKS.md tenha um valor de negócio claro e que os manuais na WIKI expliquem o benefício para o usuário final, não apenas a função técnica.
   6. Agente CTO (Guardião Técnico)
      - Método: Architectural Gatekeeping. Aplicar revisões rigorosas de Isolamento Hexagonal, proibindo que dependências externas (como Spring ou Pandas) vazem para o Domínio.
      - Estratégia de Ferramental: Antigravity Skill Mapping. Antes de cada etapa, analisar o catálogo de Skills instaladas (seção 9 do GSD-RULES) para orientar o DEV e o DBA sobre qual "pacote de conhecimento" usar, evitando redundância de código.
   7. Agente DBA (Database Administrator)
      - Método: Migrations-as-Code. Garantir que toda alteração de schema seja feita via scripts de migração SQL versionados, nunca por alterações manuais no banco.
      - Estratégia de Performance: Bulk Ingestion Optimization. Implementar estratégias de inserção em massa para evitar gargalos de I/O, especialmente para o componente de Inteligência que lida com grandes volumes de dados.
      - Segurança: Aplicar a política de Zero Leak, garantindo que dados sensíveis sejam mascarados e que o acesso siga o princípio do menor privilégio.
   8. Agente QA / Testers (Qualidade e Validação)
      - Método: Shift-Left Testing. Iniciar a escrita dos casos de teste no TESTS.md no momento em que o BA define os critérios de aceitação, antes mesmo da implementação do código.
      - Estratégia de Segurança: Continuous SAST. Realizar scans de vulnerabilidades e verificações de conformidade JWT (Deny-all por padrão) em todas as entregas atômicas.
      - Validação Final: Executar a Matriz de Regressão no final de cada etapa para garantir que novas funcionalidades não quebraram as anteriores.
- Revidar todos os documentos em busca de atualizações.