const agents = [
  { name: "CEO", role: "Orchestration", seniority: "Senior", provider: "hybrid-default" },
  { name: "DEV_BACKEND", role: "Go API", seniority: "Pleno", provider: "hybrid-default" },
  { name: "SECURITY", role: "Threat model", seniority: "Senior", provider: "api-default" },
];

const providers = [
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

document.querySelector("#agentRows").innerHTML = agents
  .map(
    (agent) => `<tr>
      <td>${agent.name}</td>
      <td>${agent.role}</td>
      <td><span class="badge">${agent.seniority}</span></td>
      <td>${agent.provider}</td>
    </tr>`,
  )
  .join("");

document.querySelector("#providerList").innerHTML = providers
  .map(
    (provider) => `<div class="list-item">
      <strong>${provider.name}</strong>
      <span>${provider.mode}</span>
      <small>${provider.model}</small>
    </div>`,
  )
  .join("");

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
