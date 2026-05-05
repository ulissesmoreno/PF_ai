const fallbackAgents = [
  { name: "CEO", role: "Orchestration", seniority: "Senior", provider: "hybrid-default" },
  { name: "DEV_BACKEND", role: "Go API", seniority: "Pleno", provider: "hybrid-default" },
  { name: "SECURITY", role: "Threat model", seniority: "Senior", provider: "api-default" },
];

const fallbackProviders = [
  { name: "API Provider", mode: "api", model: "configured by env" },
  { name: "Local Runtime", mode: "local", model: "Phase 2 container path" },
];

const memory = [
  "DOC/PLAN.md",
  "DOC/ROADMAP.md",
  "DOC/CONTEXT.md",
  "DOC/STATE.md",
  "QUESTIONS.md",
];

const apiToken = globalThis.localStorage?.getItem("pf_ai_token") ?? "";

async function loadJSON(path, fallback) {
  if (!apiToken) {
    return fallback;
  }

  try {
    const response = await fetch(path, {
      headers: {
        Authorization: `Bearer ${apiToken}`,
      },
    });
    if (!response.ok) {
      return fallback;
    }
    return await response.json();
  } catch {
    return fallback;
  }
}

function renderAgents(agents) {
  document.querySelector("#agentRows").innerHTML = agents
    .map(
      (agent) => `<tr>
        <td>${agent.name}</td>
        <td>${agent.role}</td>
        <td><span class="badge">${agent.seniority}</span></td>
        <td>${agent.provider_id ?? agent.provider}</td>
      </tr>`,
    )
    .join("");
}

function renderProviders(providers) {
  document.querySelector("#providerList").innerHTML = providers
    .map(
      (provider) => `<div class="list-item">
        <strong>${provider.name}</strong>
        <span>${provider.mode}</span>
        <small>${provider.model}</small>
      </div>`,
    )
    .join("");
}

document.querySelector("#memoryList").innerHTML = memory
  .map((path) => `<li><code>${path}</code><span>allowlisted</span></li>`)
  .join("");

document.querySelector("#handoffPreview").textContent = JSON.stringify(
  {
    header: {
      sender: "[CEO]",
      recipient: "[DEV_BACKEND:Pleno]",
      task_ref: "PHASE-1",
      intent: "PHASE_KICKOFF",
    },
    payload: {
      phase_ref: "DOC/ROADMAP.md#phase-1",
      constraints: ["TDD mandatory", "No secrets in handoffs"],
    },
  },
  null,
  2,
);

const [agents, providers] = await Promise.all([
  loadJSON("/api/agents", fallbackAgents),
  loadJSON("/api/providers", fallbackProviders),
]);

renderAgents(agents);
renderProviders(providers);
