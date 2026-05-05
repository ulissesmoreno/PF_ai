import { createServer } from "node:http";
import { readFile } from "node:fs/promises";
import { extname, join, normalize, relative } from "node:path";
import { fileURLToPath } from "node:url";

const port = Number.parseInt(process.env.PF_AI_WEB_PORT ?? "5173", 10);
const apiBase = process.env.PF_AI_API_BASE ?? "http://127.0.0.1:8081";
const root = fileURLToPath(new URL(".", import.meta.url));

const contentTypes = {
  ".html": "text/html; charset=utf-8",
  ".css": "text/css; charset=utf-8",
  ".js": "text/javascript; charset=utf-8",
};

createServer(async (request, response) => {
  const url = new URL(request.url ?? "/", `http://${request.headers.host}`);
  if (url.pathname === "/health" || url.pathname.startsWith("/api/")) {
    await proxyAPI(request, response, url);
    return;
  }

  const requestedPath = url.pathname === "/" ? "/index.html" : url.pathname;
  const normalized = normalize(requestedPath).replace(/^[/\\]+/, "").replace(/^(\.\.[/\\])+/, "");
  const filePath = join(root, normalized);
  const relativePath = relative(root, filePath);

  if (relativePath.startsWith("..")) {
    response.writeHead(403, { "Content-Type": "text/plain; charset=utf-8" });
    response.end("Forbidden");
    return;
  }

  try {
    const body = await readFile(filePath);
    response.writeHead(200, {
      "Content-Type": contentTypes[extname(filePath)] ?? "application/octet-stream",
    });
    response.end(body);
  } catch {
    response.writeHead(404, { "Content-Type": "text/plain; charset=utf-8" });
    response.end("Not found");
  }
}).listen(port, () => {
  console.log(`PF_ai web listening on http://localhost:${port}`);
});

async function proxyAPI(request, response, url) {
  const target = new URL(url.pathname + url.search, apiBase);
  try {
    const upstream = await fetch(target, {
      method: request.method,
      headers: request.headers,
      body: ["GET", "HEAD"].includes(request.method ?? "GET") ? undefined : request,
      duplex: "half",
    });
    const body = await upstream.arrayBuffer();
    response.writeHead(upstream.status, {
      "Content-Type": upstream.headers.get("content-type") ?? "application/json",
    });
    response.end(Buffer.from(body));
  } catch {
    response.writeHead(502, { "Content-Type": "application/json" });
    response.end(JSON.stringify({ error: "backend unavailable" }));
  }
}
