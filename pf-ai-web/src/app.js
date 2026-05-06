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
  selectedMemory: "DOC/STATE.md",
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
  };
}

function renderAgents(agents) {
  document.querySelector("#agentRows").innerHTML = agents.length
    ? agents
    .map(
      (agent) => `<tr>
        <td>${escapeHTML(agent.name)}</td>
        <td>${escapeHTML(agent.role)}</td>
        <td><span class="badge">${escapeHTML(agent.seniority)}</span></td>
        <td>${escapeHTML(agent.provider_id ?? agent.provider)}</td>
      </tr>`,
    )
    .join("")
    : `<tr><td colspan="4">No agents registered.</td></tr>`;
}

function renderProviders(providers) {
  document.querySelector("#providerList").innerHTML = providers.length
    ? providers
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

  const options = [
    `<option value="">No provider</option>`,
    ...providers.map((provider) => `<option value="${escapeHTML(provider.id)}">${escapeHTML(provider.name)} (${escapeHTML(provider.mode)})</option>`),
  ];
  document.querySelector("#agentProvider").innerHTML = options.join("");
}

function renderBoard(cards) {
  const labels = {
    todo: "To Do",
    in_progress: "In Progress",
    done: "Done",
  };
  document.querySelector("#planningBoard").innerHTML = ["todo", "in_progress", "done"]
    .map((status) => {
      const items = cards.filter((card) => card.status === status);
      return `<section class="board-column">
        <h3>${labels[status]}</h3>
        ${items.length ? items.map((card) => `<article class="card">
          <strong>${escapeHTML(card.title)}</strong>
          <div class="card-meta">
            <span class="badge">${escapeHTML(card.priority)}</span>
            <span class="badge">${escapeHTML(card.owner)}</span>
          </div>
          <small>${escapeHTML(card.phase)} / ${escapeHTML(card.task_ref)}</small>
        </article>`).join("") : `<div class="list-item"><strong>No cards.</strong></div>`}
      </section>`;
    })
    .join("");
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
  const [agents, providers, cards] = await Promise.all([
    loadJSON("/api/agents", defaultAgents, "pf_ai_agents"),
    loadJSON("/api/providers", [], "pf_ai_providers"),
    loadJSON("/api/planning-cards", defaultCards, "pf_ai_cards"),
  ]);

  renderAgents(agents);
  renderProviders(providers);
  renderBoard(cards.length ? cards : defaultCards);
}

document.querySelector("#apiBase").value = state.apiBase;
document.querySelector("#authToken").value = state.apiToken;
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
  const payload = formData(event.currentTarget);
  try {
    await requestJSON("/api/providers", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    event.currentTarget.reset();
    await refresh();
    setStatus("Provider configured", true);
  } catch (error) {
    const providers = localItems("pf_ai_providers", []);
    const nextProviders = [...providers.filter((provider) => provider.id !== payload.id), payload];
    setLocalItems("pf_ai_providers", nextProviders);
    event.currentTarget.reset();
    await refresh();
    setStatus("Provider saved locally", true);
  }
});

document.querySelector("#projectForm").addEventListener("submit", async (event) => {
  event.preventDefault();
  const values = formData(event.currentTarget);
  try {
    const payload = {
      id: values.id,
      name: values.name,
      onboarding: JSON.parse(values.onboarding),
    };
    const result = await requestJSON("/api/projects", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    document.querySelector("#projectPreview").textContent = JSON.stringify(result, null, 2);
    setStatus("Project registered", true);
  } catch (error) {
    document.querySelector("#projectPreview").textContent = error.message;
    setStatus("Project registration failed", false);
  }
});

document.querySelector("#cardForm").addEventListener("submit", async (event) => {
  event.preventDefault();
  const payload = formData(event.currentTarget);
  try {
    await requestJSON("/api/planning-cards", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    event.currentTarget.reset();
    await refresh();
    setStatus("Card created", true);
  } catch (error) {
    const cards = localItems("pf_ai_cards", defaultCards);
    const nextCards = [...cards.filter((card) => card.id !== payload.id), payload];
    setLocalItems("pf_ai_cards", nextCards);
    event.currentTarget.reset();
    await refresh();
    setStatus("Card saved locally", true);
  }
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
    const agents = localItems("pf_ai_agents", defaultAgents);
    const nextAgents = [...agents.filter((agent) => agent.id !== payload.id), payload];
    setLocalItems("pf_ai_agents", nextAgents);
    event.currentTarget.reset();
    await refresh();
    setStatus("Agent saved locally", true);
  }
});

document.querySelector("#handoffForm").addEventListener("input", () => renderHandoffPreview());
document.querySelector("#handoffForm").addEventListener("submit", async (event) => {
  event.preventDefault();
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
    renderHandoffPreview(result.path);
    setStatus("Handoff created", true);
  } catch (error) {
    const localPath = `.agent_handoff/local_${Date.now()}_${values.task_ref}.json`;
    const handoffs = localItems("pf_ai_handoffs", []);
    setLocalItems("pf_ai_handoffs", [...handoffs, { path: localPath, request: values }]);
    renderHandoffPreview(localPath);
    setStatus("Handoff staged locally", true);
  }
});

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
