const memory = [
  "DOC/PLAN.md",
  "DOC/ROADMAP.md",
  "DOC/CONTEXT.md",
  "DOC/STATE.md",
  "QUESTIONS.md",
];

const state = {
  apiBase: globalThis.localStorage?.getItem("pf_ai_api_base") ?? "http://127.0.0.1:8081",
  apiToken: globalThis.localStorage?.getItem("pf_ai_token") ?? "",
  selectedMemory: "DOC/STATE.md",
};

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

async function loadJSON(path, fallback) {
  try {
    return await requestJSON(path);
  } catch (error) {
    setStatus(error.message === "token required" ? "Token required" : "API unavailable", false);
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
      </div>`,
    )
    .join("")
    : `<div class="list-item"><strong>No providers configured.</strong><span>Use Configure.</span></div>`;
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
  const [agents, providers] = await Promise.all([
    loadJSON("/api/agents", []),
    loadJSON("/api/providers", []),
  ]);

  renderAgents(agents);
  renderProviders(providers);
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
  try {
    await requestJSON("/api/providers", {
      method: "POST",
      body: JSON.stringify(formData(event.currentTarget)),
    });
    event.currentTarget.reset();
    await refresh();
    setStatus("Provider configured", true);
  } catch (error) {
    setStatus(error.message, false);
  }
});

document.querySelector("#agentForm").addEventListener("submit", async (event) => {
  event.preventDefault();
  try {
    await requestJSON("/api/agents", {
      method: "POST",
      body: JSON.stringify(formData(event.currentTarget)),
    });
    event.currentTarget.reset();
    await refresh();
    setStatus("Agent created", true);
  } catch (error) {
    setStatus(error.message, false);
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
    setStatus(error.message, false);
  }
});

renderHandoffPreview();
await refresh();
