const memory = [
  "DOC/PLAN.md",
  "DOC/ROADMAP.md",
  "DOC/CONTEXT.md",
  "DOC/STATE.md",
  "QUESTIONS.md",
];

const defaultAgents = [
  { id: "ceo", name: "CEO", role: "CEO", seniority: "Senior", provider_id: "", description: "Orchestration" },
  { id: "cto", name: "CTO", role: "CTO", seniority: "Senior", provider_id: "", description: "Architecture" },
  { id: "ba", name: "BA", role: "BA", seniority: "Senior", provider_id: "", description: "Business analysis" },
];

const state = {
  apiBase: globalThis.localStorage?.getItem("pf_ai_api_base") ?? "http://127.0.0.1:8081",
  apiToken: globalThis.localStorage?.getItem("pf_ai_token") ?? "",
  activeProjectID: globalThis.localStorage?.getItem("pf_ai_active_project") ?? "",
  selectedMemory: "DOC/STATE.md",
  selectedCard: "",
};

const defaultCards = [
  { id: "phase3-mcp", title: "MCP handoff envelope", owner: "DEV_BACKEND", status: "todo", priority: "high", phase: "Phase 3", task_ref: "PHASE3-BACKEND-001" },
  { id: "phase3-board", title: "Planning board dashboard", owner: "DEV_FRONTEND", status: "todo", priority: "high", phase: "Phase 3", task_ref: "PHASE3-FRONTEND-001" },
  { id: "phase3-release", title: "Release readiness review", owner: "QA", status: "todo", priority: "medium", phase: "Phase 3", task_ref: "PHASE3-QA-001" },
];

function localItems(key, fallback = []) {
  try {
    return JSON.parse(globalThis.localStorage?.getItem(key) ?? "null") ?? fallback;
  } catch {
    return fallback;
  }
}

function setLocalItems(key, items) {
  globalThis.localStorage?.setItem(key, JSON.stringify(items));
}

function activeProject(projects = localItems("pf_ai_projects", [])) {
  return projects.find((project) => project.id === state.activeProjectID) ?? null;
}

function scopedKey(key) {
  return state.activeProjectID ? `${key}_${state.activeProjectID}` : key;
}

function scopedDefaults(items) {
  if (!state.activeProjectID) {
    return [];
  }
  return items.map((item) => ({ ...item, project_id: state.activeProjectID }));
}

function token() {
  return state.apiToken;
}

function apiURL(path) {
  return new URL(path, state.apiBase).toString();
}

function headers() {
  return {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token()}`,
  };
}

async function requestJSON(path, options = {}) {
  if (!token()) {
    throw new Error("token required");
  }

  const response = await fetch(apiURL(path), {
    ...options,
    headers: {
      ...headers(),
      ...(options.headers ?? {}),
    },
  });
  const payload = await response.json().catch(() => ({}));
  if (!response.ok) {
    throw new Error(payload.error ?? `request failed: ${response.status}`);
  }
  return payload;
}

async function loadJSON(path, fallback, localKey = "") {
  try {
    return await requestJSON(path);
  } catch (error) {
    setStatus(error.message === "token required" ? "Token required" : "API unavailable", false);
    if (localKey) {
      return localItems(localKey, fallback);
    }
    return fallback;
  }
}

function setStatus(message, ok) {
  const status = document.querySelector("#apiStatus");
  status.textContent = message;
  status.dataset.state = ok ? "ok" : "warn";
}

function escapeHTML(value) {
  return String(value ?? "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

function formData(form) {
  return Object.fromEntries(new FormData(form).entries());
}

function slug(value) {
  return String(value ?? "")
    .trim()
    .toLowerCase()
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}

function projectPayload(form) {
  const values = formData(form);
  const id = slug(values.name);
  return {
    id,
    name: values.name,
    description: values.description,
    audience: values.audience,
    technical_stack: values.technical_stack,
    onboarding: {
      description: values.description,
      audience: values.audience,
      technical_stack: values.technical_stack,
    },
  };
}

function agentPayload(form) {
  const values = formData(form);
  const role = values.role === "OTHER" ? values.custom_role : values.role;
  return {
    id: slug(`${role}-${values.name}`),
    name: values.name,
    role,
    seniority: values.seniority,
    provider_id: values.provider_id,
    description: values.description,
    project_id: state.activeProjectID,
  };
}

function renderAgents(agents) {
  const scopedAgents = scoped(agents);
  const visibleAgents = scopedAgents.length ? scopedAgents : scopedDefaults(defaultAgents);
  document.querySelector("#agentRows").innerHTML = visibleAgents.length
    ? visibleAgents
    .map(
      (agent) => `<tr>
        <td>${escapeHTML(agent.name)}</td>
        <td>${escapeHTML(agent.role)}</td>
        <td><span class="badge">${escapeHTML(agent.seniority)}</span></td>
        <td>${escapeHTML(agent.provider_id ?? agent.provider)}</td>
      </tr>`,
    )
    .join("")
    : `<tr><td colspan="4">Select a project to view agents.</td></tr>`;
}

function renderProviders(providers) {
  const scopedProviders = scoped(providers);
  const html = scopedProviders.length
    ? scopedProviders
    .map(
      (provider) => `<div class="list-item">
        <strong>${escapeHTML(provider.name)}</strong>
        <span>${escapeHTML(provider.mode)}</span>
        <small>${escapeHTML(provider.model)}</small>
        <button type="button" data-provider-health="${escapeHTML(provider.id)}">Check Status</button>
      </div>`,
    )
    .join("")
    : `<div class="list-item"><strong>No providers configured.</strong><span>Use Configure.</span></div>`;
  document.querySelector("#providerList").innerHTML = html;
  document.querySelector("#providerListSummary").innerHTML = html;

  const options = [
    `<option value="">No provider</option>`,
    ...scopedProviders.map((provider) => `<option value="${escapeHTML(provider.id)}">${escapeHTML(provider.name)} (${escapeHTML(provider.mode)})</option>`),
  ];
  document.querySelector("#agentProvider").innerHTML = options.join("");
}

function scoped(items) {
  if (!state.activeProjectID) {
    return [];
  }
  return items.filter((item) => item.project_id === state.activeProjectID);
}

function cardComments(cardID) {
  return localItems(scopedKey("pf_ai_comments"), []).filter((comment) => comment.card_id === cardID);
}

function renderBoard(cards) {
  const scopedCards = scoped(cards);
  const labels = {
    todo: "To Do",
    in_progress: "In Progress",
    done: "Done",
  };
  document.querySelector("#planningBoard").innerHTML = ["todo", "in_progress", "done"]
    .map((status) => {
      const items = scopedCards.filter((card) => card.status === status);
      return `<section class="board-column">
        <h3>${labels[status]}</h3>
        ${items.length ? items.map((card) => `<article class="card ${state.selectedCard === card.id ? "selected" : ""}" data-card-id="${escapeHTML(card.id)}">
          <strong>${escapeHTML(card.title)}</strong>
          <div class="card-meta">
            <span class="badge">${escapeHTML(card.priority)}</span>
            <span class="badge">${escapeHTML(card.owner)}</span>
          </div>
          <small>${escapeHTML(card.phase)} / ${escapeHTML(card.task_ref)}</small>
          <div class="comments">
            ${cardComments(card.id).slice(-3).map((comment) => `<div class="comment"><strong>${escapeHTML(comment.author)}</strong><br />${escapeHTML(comment.body)}</div>`).join("")}
          </div>
        </article>`).join("") : `<div class="list-item"><strong>No cards.</strong></div>`}
      </section>`;
    })
    .join("");
}

function showPage(page) {
  document.querySelectorAll("[data-page]").forEach((section) => {
    section.classList.toggle("active-page", section.dataset.page === page);
  });
  document.querySelectorAll("[data-page-link]").forEach((link) => {
    link.classList.toggle("active", link.dataset.pageLink === page);
  });
}

function selectProject(id) {
  state.activeProjectID = id;
  globalThis.localStorage?.setItem("pf_ai_active_project", id);
  state.selectedCard = "";
}

function renderProjects(projects) {
  const current = activeProject(projects);
  document.querySelector("#activeProjectLabel").textContent = current ? current.name : "No project selected";
  document.querySelector("#projectNav").innerHTML = projects.length
    ? projects
      .map((project) => `<a href="#projects" class="${project.id === state.activeProjectID ? "active" : ""}" data-project-id="${escapeHTML(project.id)}">
        <strong>${escapeHTML(project.name)}</strong>
        <small>${escapeHTML(project.technical_stack)}</small>
      </a>`)
      .join("")
    : `<div class="empty-state">Register a project to start.</div>`;
  document.querySelector("#projectDetail").innerHTML = current
    ? `<div class="list-item">
        <strong>${escapeHTML(current.name)}</strong>
        <span>${escapeHTML(current.audience)}</span>
        <small>${escapeHTML(current.technical_stack)}</small>
      </div>
      <div class="list-item"><strong>Description</strong><span>${escapeHTML(current.description)}</span></div>`
    : `<div class="list-item"><strong>No project selected.</strong><span>Create or select a project.</span></div>`;
  renderQuestionSignal(current);
}

function renderQuestionSignal(project) {
  if (!project) {
    document.querySelector("#questionSignal").innerHTML = "";
    return;
  }
  const signals = localItems("pf_ai_project_question_signals", []).filter((signal) => signal.project_id === project.id);
  document.querySelector("#questionSignal").innerHTML = signals.length
    ? signals.map((signal) => `<div class="list-item"><strong>CEO question signal</strong><span>${escapeHTML(signal.body)}</span></div>`).join("")
    : `<div class="list-item"><strong>No pending question signal.</strong><span>Onboarding summary complete.</span></div>`;
}

function requireActiveProject() {
  if (state.activeProjectID) {
    return true;
  }
  setStatus("Select a project first", false);
  showPage("projects");
  return false;
}

function addProjectQuestionSignal(project) {
  const signals = localItems("pf_ai_project_question_signals", []);
  const body = `Novo projeto "${project.name}" registrado. CEO deve validar se as informacoes resumidas bastam para iniciar o onboarding.`;
  setLocalItems("pf_ai_project_question_signals", [
    ...signals.filter((signal) => signal.project_id !== project.id),
    { project_id: project.id, body, created_at: new Date().toISOString() },
  ]);
}

function redactProviderStatus(status) {
  const safe = {
    provider_id: status.provider_id,
    mode: status.mode,
    route: status.route,
    status: status.status,
    fallback_reason: status.fallback_reason,
  };
  return JSON.stringify(safe, null, 2);
}

document.querySelector("#memoryList").innerHTML = memory
  .map((path) => `<li data-path="${escapeHTML(path)}"><code>${escapeHTML(path)}</code><span>allowlisted</span></li>`)
  .join("");

function renderHandoffPreview(path = "") {
  document.querySelector("#handoffPreview").textContent = JSON.stringify(
    {
      path,
      request: formData(document.querySelector("#handoffForm")),
    },
    null,
    2,
  );
}

async function refresh() {
  const [projects, agents, providers, cards] = await Promise.all([
    loadJSON("/api/projects", [], "pf_ai_projects"),
    loadJSON("/api/agents", scopedDefaults(defaultAgents), scopedKey("pf_ai_agents")),
    loadJSON("/api/providers", [], scopedKey("pf_ai_providers")),
    loadJSON("/api/planning-cards", scopedDefaults(defaultCards), scopedKey("pf_ai_cards")),
  ]);

  if (!state.activeProjectID && projects.length) {
    selectProject(projects[0].id);
  }
  renderProjects(projects);
  renderAgents(agents);
  renderProviders(providers);
  renderBoard(cards.length ? cards : scopedDefaults(defaultCards));
}

document.querySelector("#apiBase").value = state.apiBase;
document.querySelector("#authToken").value = state.apiToken;
document.querySelectorAll("[data-page-link]").forEach((link) => {
  link.addEventListener("click", (event) => {
    event.preventDefault();
    showPage(link.dataset.pageLink);
  });
});

document.querySelector("#projectNav").addEventListener("click", async (event) => {
  const link = event.target.closest("[data-project-id]");
  if (!link) {
    return;
  }
  event.preventDefault();
  selectProject(link.dataset.projectId);
  showPage("projects");
  await refresh();
});
document.querySelector("#authForm").addEventListener("submit", async (event) => {
  event.preventDefault();
  const values = formData(event.currentTarget);
  state.apiBase = values.apiBase.trim();
  state.apiToken = values.authToken.trim();
  globalThis.localStorage?.setItem("pf_ai_api_base", state.apiBase);
  globalThis.localStorage?.setItem("pf_ai_token", state.apiToken);
  await refresh();
  if (state.apiToken) {
    setStatus("Connected", true);
  }
});

document.querySelectorAll("#memoryList li").forEach((item) => {
  item.addEventListener("click", () => {
    state.selectedMemory = item.dataset.path;
    document.querySelectorAll("#memoryList li").forEach((node) => node.classList.remove("selected"));
    item.classList.add("selected");
  });
});

document.querySelector("#memoryList li[data-path='DOC/STATE.md']")?.classList.add("selected");
document.querySelector("#agentRole").addEventListener("change", (event) => {
  const customRole = document.querySelector("#customRole");
  const isOther = event.currentTarget.value === "OTHER";
  customRole.hidden = !isOther;
  customRole.required = isOther;
  if (!isOther) {
    customRole.value = "";
  }
});

document.querySelector("#readMemory").addEventListener("click", async () => {
  try {
    const result = await requestJSON(`/api/memory?path=${encodeURIComponent(state.selectedMemory)}`);
    document.querySelector("#memoryContent").textContent = result.content;
    setStatus("Memory read", true);
  } catch (error) {
    document.querySelector("#memoryContent").textContent = error.message;
    setStatus("Memory read failed", false);
  }
});

document.querySelector("#providerForm").addEventListener("submit", async (event) => {
  event.preventDefault();
  if (!requireActiveProject()) {
    return;
  }
  const payload = { ...formData(event.currentTarget), project_id: state.activeProjectID };
  try {
    await requestJSON("/api/providers", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    event.currentTarget.reset();
    await refresh();
    setStatus("Provider configured", true);
  } catch (error) {
    const providers = localItems(scopedKey("pf_ai_providers"), []);
    const nextProviders = [...providers.filter((provider) => provider.id !== payload.id), payload];
    setLocalItems(scopedKey("pf_ai_providers"), nextProviders);
    event.currentTarget.reset();
    await refresh();
    setStatus("Provider saved locally", true);
  }
});

document.querySelector("#projectForm").addEventListener("submit", async (event) => {
  event.preventDefault();
  const payload = projectPayload(event.currentTarget);
  try {
    const result = await requestJSON("/api/projects", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    const projects = localItems("pf_ai_projects", []);
    setLocalItems("pf_ai_projects", [...projects.filter((project) => project.id !== result.id), result]);
    selectProject(result.id);
    addProjectQuestionSignal(result);
    document.querySelector("#projectPreview").textContent = JSON.stringify(result, null, 2);
    event.currentTarget.reset();
    await refresh();
    setStatus("Project registered", true);
  } catch (error) {
    const projects = localItems("pf_ai_projects", []);
    setLocalItems("pf_ai_projects", [...projects.filter((project) => project.id !== payload.id), payload]);
    selectProject(payload.id);
    addProjectQuestionSignal(payload);
    document.querySelector("#projectPreview").textContent = JSON.stringify(payload, null, 2);
    event.currentTarget.reset();
    await refresh();
    setStatus("Project saved locally", true);
  }
});

document.querySelector("#cardForm").addEventListener("submit", async (event) => {
  event.preventDefault();
  if (!requireActiveProject()) {
    return;
  }
  const payload = { ...formData(event.currentTarget), project_id: state.activeProjectID };
  try {
    await requestJSON("/api/planning-cards", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    event.currentTarget.reset();
    await refresh();
    setStatus("Card created", true);
  } catch (error) {
    const cards = localItems(scopedKey("pf_ai_cards"), scopedDefaults(defaultCards));
    const nextCards = [...cards.filter((card) => card.id !== payload.id), payload];
    setLocalItems(scopedKey("pf_ai_cards"), nextCards);
    event.currentTarget.reset();
    await refresh();
    setStatus("Card saved locally", true);
  }
});

document.querySelector("#planningBoard").addEventListener("click", async (event) => {
  const card = event.target.closest("[data-card-id]");
  if (!card) {
    return;
  }
  state.selectedCard = card.dataset.cardId;
  document.querySelector("#commentForm [name='card_id']").value = state.selectedCard;
  await refresh();
});

document.querySelector("#commentForm").addEventListener("submit", async (event) => {
  event.preventDefault();
  if (!requireActiveProject()) {
    return;
  }
  const payload = {
    ...formData(event.currentTarget),
    created_at: new Date().toISOString(),
  };
  const scopedComments = localItems(scopedKey("pf_ai_comments"), []);
  setLocalItems(scopedKey("pf_ai_comments"), [...scopedComments, payload]);
  document.querySelector("#commentPreview").textContent = JSON.stringify(payload, null, 2);
  event.currentTarget.reset();
  await refresh();
  setStatus("Comment added", true);
});

document.querySelector("#providerList").addEventListener("click", async (event) => {
  const button = event.target.closest("[data-provider-health]");
  if (!button) {
    return;
  }
  try {
    const providerID = button.dataset.providerHealth;
    const result = await requestJSON(`/api/provider-health?id=${encodeURIComponent(providerID)}`);
    document.querySelector("#providerStatus").textContent = redactProviderStatus(result);
    setStatus("Provider status checked", true);
  } catch (error) {
    document.querySelector("#providerStatus").textContent = error.message;
    setStatus("Provider status unavailable", false);
  }
});

document.querySelector("#agentForm").addEventListener("submit", async (event) => {
  event.preventDefault();
  if (!requireActiveProject()) {
    return;
  }
  const payload = agentPayload(event.currentTarget);
  try {
    await requestJSON("/api/agents", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    event.currentTarget.reset();
    await refresh();
    setStatus("Agent created", true);
  } catch (error) {
    const agents = localItems(scopedKey("pf_ai_agents"), scopedDefaults(defaultAgents));
    const nextAgents = [...agents.filter((agent) => agent.id !== payload.id), payload];
    setLocalItems(scopedKey("pf_ai_agents"), nextAgents);
    event.currentTarget.reset();
    await refresh();
    setStatus("Agent saved locally", true);
  }
});

document.querySelector("#handoffForm").addEventListener("input", () => renderHandoffPreview());
document.querySelector("#handoffForm").addEventListener("submit", async (event) => {
  event.preventDefault();
  if (!requireActiveProject()) {
    return;
  }
  const values = formData(event.currentTarget);
  try {
    const payload = JSON.parse(values.payload);
    const result = await requestJSON("/api/handoffs", {
      method: "POST",
      body: JSON.stringify({
        sender: values.sender,
        recipient: values.recipient,
        task_ref: values.task_ref,
        intent: values.intent,
        payload,
      }),
    });
    addHandoffComment(values, result.path);
    renderHandoffPreview(result.path);
    setStatus("Handoff created", true);
  } catch (error) {
    const localPath = `.agent_handoff/local_${Date.now()}_${values.task_ref}.json`;
    const handoffs = localItems(scopedKey("pf_ai_handoffs"), []);
    setLocalItems(scopedKey("pf_ai_handoffs"), [...handoffs, { path: localPath, request: values, project_id: state.activeProjectID }]);
    addHandoffComment(values, localPath);
    renderHandoffPreview(localPath);
    setStatus("Handoff staged locally", true);
  }
});

function addHandoffComment(values, path) {
  const comments = localItems(scopedKey("pf_ai_comments"), []);
  const cardID = state.selectedCard || values.task_ref;
  setLocalItems(scopedKey("pf_ai_comments"), [
    ...comments,
    {
      card_id: cardID,
      author: "HANDOFF",
      body: `${values.sender} -> ${values.recipient}: ${path}`,
      created_at: new Date().toISOString(),
    },
  ]);
}

document.querySelector("#exportMcp").addEventListener("click", async () => {
  const values = formData(document.querySelector("#handoffForm"));
  try {
    const result = await requestJSON("/api/mcp-envelope", {
      method: "POST",
      body: JSON.stringify({
        sender: values.sender,
        recipient: values.recipient,
        task_ref: values.task_ref,
        intent: values.intent,
        payload: JSON.parse(values.payload),
      }),
    });
    document.querySelector("#handoffPreview").textContent = JSON.stringify(result, null, 2);
    setStatus("MCP envelope exported", true);
  } catch (error) {
    document.querySelector("#handoffPreview").textContent = error.message;
    setStatus("MCP export failed", false);
  }
});

renderHandoffPreview();
await refresh();
