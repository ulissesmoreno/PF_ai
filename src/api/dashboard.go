package api

import "net/http"

func (h *Handler) dashboard(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(dashboardHTML))
}

const dashboardHTML = `<!doctype html>
<html lang="pt-BR">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>PF AI Cards</title>
  <style>
    :root { color-scheme: light; --bg:#f6f7f9; --panel:#fff; --line:#d9dee7; --text:#141820; --muted:#667085; --accent:#0969da; --blocked:#b54708; --done:#067647; --failed:#b42318; }
    * { box-sizing: border-box; }
    body { margin:0; font-family: Segoe UI, Arial, sans-serif; background:var(--bg); color:var(--text); }
    header { height:56px; display:flex; align-items:center; gap:12px; padding:0 18px; background:var(--panel); border-bottom:1px solid var(--line); }
    h1 { font-size:18px; margin:0 10px 0 0; }
    select, input, button, textarea { font:inherit; border:1px solid var(--line); background:#fff; color:var(--text); border-radius:6px; padding:8px 10px; }
    button { cursor:pointer; background:var(--text); color:#fff; border-color:var(--text); }
    main { display:grid; grid-template-columns: 1fr 380px; gap:0; min-height:calc(100vh - 56px); }
    .board { padding:16px; display:grid; grid-template-columns: repeat(5, minmax(190px, 1fr)); gap:12px; overflow:auto; }
    .column { min-width:190px; }
    .column h2 { font-size:13px; text-transform:uppercase; letter-spacing:.04em; color:var(--muted); margin:0 0 8px; display:flex; justify-content:space-between; }
    .card { background:var(--panel); border:1px solid var(--line); border-radius:8px; padding:10px; margin-bottom:8px; cursor:pointer; }
    .card:hover { border-color:#9aa4b2; }
    .title { font-size:14px; font-weight:650; margin-bottom:8px; }
    .meta { font-size:12px; color:var(--muted); display:grid; gap:3px; }
    .badge { display:inline-block; width:max-content; padding:2px 7px; border-radius:999px; background:#eef2f6; color:#344054; font-size:12px; }
    .blocked { color:var(--blocked); } .done { color:var(--done); } .failed { color:var(--failed); }
    aside { background:var(--panel); border-left:1px solid var(--line); padding:16px; overflow:auto; max-height:calc(100vh - 56px); }
    .empty { color:var(--muted); padding:18px 0; }
    .thread h2 { margin:0 0 6px; font-size:18px; }
    .thread .sub { color:var(--muted); margin-bottom:12px; font-size:13px; }
    .comment { border-top:1px solid var(--line); padding:10px 0; }
    .comment strong { font-size:13px; }
    pre { white-space:pre-wrap; word-break:break-word; background:#f2f4f7; padding:8px; border-radius:6px; font-size:12px; max-height:220px; overflow:auto; }
    .formbox { border-top:1px solid var(--line); margin-top:12px; padding-top:12px; display:grid; gap:8px; }
    textarea { min-height:96px; resize:vertical; }
    @media (max-width: 900px) { main { grid-template-columns:1fr; } aside { border-left:0; border-top:1px solid var(--line); max-height:none; } .board { grid-template-columns:1fr; } }
  </style>
</head>
<body>
  <header>
    <h1>PF AI Cards</h1>
    <select id="project"></select>
    <select id="status">
      <option value="">Todos</option><option>open</option><option>in_progress</option><option>blocked</option><option>done</option><option>failed</option><option>canceled</option>
    </select>
    <input id="task" placeholder="task_ref">
    <button id="refresh">Atualizar</button>
  </header>
  <main>
    <section class="board" id="board"></section>
    <aside id="details"><div class="empty">Selecione um card.</div></aside>
  </main>
  <script>
    const statuses = ["open","in_progress","blocked","done","failed"];
    const board = document.getElementById("board");
    const details = document.getElementById("details");
    const project = document.getElementById("project");
    const statusFilter = document.getElementById("status");
    const task = document.getElementById("task");
    document.getElementById("refresh").onclick = loadCards;
    project.onchange = loadCards; statusFilter.onchange = loadCards; task.onchange = loadCards;

    async function api(path, opts) {
      const res = await fetch(path, opts);
      if (!res.ok) throw new Error(await res.text());
      return res.json();
    }
    async function loadProjects() {
      const data = await api("/api/projects");
      project.innerHTML = '<option value="">Todos os projetos</option>' + data.projects.map(p => '<option value="'+p.id+'">'+escapeHtml(p.slug || p.name)+'</option>').join("");
    }
    async function loadCards() {
      const params = new URLSearchParams();
      if (project.value) params.set("project_id", project.value);
      if (statusFilter.value) params.set("status", statusFilter.value);
      if (task.value) params.set("task_ref", task.value);
      const data = await api("/api/cards?" + params.toString());
      renderBoard(data.cards || []);
    }
    function renderBoard(cards) {
      board.innerHTML = "";
      const grouped = Object.fromEntries(statuses.map(s => [s, []]));
      for (const card of cards) (grouped[card.status] ||= []).push(card);
      for (const s of statuses) {
        const col = document.createElement("div"); col.className = "column";
        col.innerHTML = '<h2><span>'+s+'</span><span>'+((grouped[s]||[]).length)+'</span></h2>';
        for (const card of grouped[s] || []) col.appendChild(cardEl(card));
        board.appendChild(col);
      }
    }
    function cardEl(card) {
      const el = document.createElement("div"); el.className = "card";
      el.innerHTML = '<div class="title">'+escapeHtml(card.title || card.intent || card.id)+'</div>'+
        '<div class="meta"><span>'+escapeHtml(card.recipient || "")+'</span><span>'+escapeHtml(card.task_ref || "")+'</span>'+
        '<span class="badge '+card.status+'">'+card.status+'</span><span>retry '+card.retry_count+'/'+card.max_retries+'</span></div>';
      el.onclick = () => loadThread(card.id);
      return el;
    }
    async function loadThread(id) {
      const thread = await api("/api/cards/" + encodeURIComponent(id));
      const card = thread.card;
      details.innerHTML = '<div class="thread"><h2>'+escapeHtml(card.title || card.id)+'</h2><div class="sub">'+escapeHtml(card.id)+' · '+escapeHtml(card.status)+'</div><div id="comments"></div></div>';
      const comments = document.getElementById("comments");
      for (const c of thread.comments || []) {
        const div = document.createElement("div"); div.className = "comment";
        div.innerHTML = '<strong>'+escapeHtml(c.author)+' · '+escapeHtml(c.comment_type)+'</strong><pre>'+escapeHtml(formatContent(c.content))+'</pre>';
        comments.appendChild(div);
      }
      if (card.status === "blocked") renderResponseForm(card.id);
    }
    function renderResponseForm(id) {
      const form = document.createElement("div"); form.className = "formbox";
      form.innerHTML = '<strong>Responder card</strong><textarea id="answer" placeholder="Resposta para o agente"></textarea><button id="send">Enviar resposta</button>';
      details.appendChild(form);
      document.getElementById("send").onclick = async () => {
        const answer = document.getElementById("answer").value;
        await api("/api/cards/" + encodeURIComponent(id) + "/respond", {method:"PUT", headers:{"Content-Type":"application/json"}, body:JSON.stringify({author:"[HUMAN]", response:{answer}})});
        await loadThread(id); await loadCards();
      };
    }
    function formatContent(content) { try { return JSON.stringify(JSON.parse(content), null, 2); } catch { return content || ""; } }
    function escapeHtml(s) { return String(s ?? "").replace(/[&<>"']/g, m => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[m])); }
    loadProjects().then(loadCards).catch(e => details.innerHTML = '<pre>'+escapeHtml(e.message)+'</pre>');
    setInterval(loadCards, 5000);
  </script>
</body>
</html>`
