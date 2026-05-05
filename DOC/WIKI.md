# Project Wiki — Maintenance Documentation

A knowledge base maintained by the Agent and developers.
Based on Andrej Karpathy's LLM Wiki pattern.
**Recommended tool: Obsidian** — see `ENV_SETUP.md §9` for installation and setup.

---

## Purpose

This wiki is a structured, interlinked knowledge base for the project.
The Agent maintains the wiki. The developer curates sources, asks questions, and guides the analysis.
The goal is to ensure **maintenance traceability** — architectural decisions, known bugs, configurations, and history of significant changes.

---

## Folder Structure

```
wiki/             -- markdown pages maintained by Agent (and devs)
wiki/index.md     -- master table of contents for the entire wiki
wiki/log.md       -- append-only record of all operations
raw/              -- source documents (immutable — never modify these)
```

### How to create the initial structure
```bash
# From the project root
mkdir -p wiki raw
touch wiki/index.md wiki/log.md
```

> **Note:** The `wiki/` and `raw/` directories should be versioned in the repository.
> Add only Obsidian cache files to `.gitignore` (`.obsidian/workspace`).

---

## Ingest Workflow

When the developer adds a new source to `raw/` and asks you to ingest it:

1. Read the full source document.
2. Discuss key takeaways with the developer before writing anything.
3. Create a summary page in `wiki/` named after the source.
4. Create or update concept pages for each major idea or entity.
5. Add wiki-links (`[[page-name]]`) to connect related pages.
6. Update `wiki/index.md` with new pages and one-line descriptions.
7. Append an entry to `wiki/log.md` with the date, source name, and what changed.

> A single source may touch 10–15 wiki pages. That is normal.

---

## Page Format

Every wiki page should follow this structure:

```markdown
# Page Title

**Summary**: One to two sentences describing this page.

**Sources**: List of raw source files this page draws from.

**Last updated**: [YYYY-MM-DD HH:MM] — Responsible.

---

Main content goes here. Use clear headings and short paragraphs.

Link to related concepts using [[wiki-links]] throughout the text.

## Related Pages

- [[related-concept-1]]
- [[related-concept-2]]
```

---

## Citation Rules

- Every factual claim should reference its source file.
- Use the format `(source: filename.pdf)` after the claim.
- If two sources disagree, note the contradiction explicitly.
- If a claim has no source, mark it as `[needs verification]`.

---

## Question Answering

When the developer asks a question:

1. Read `wiki/index.md` first to find relevant pages.
2. Read those pages and synthesize an answer.
3. Cite specific wiki pages in your response.
4. If the answer is not in the wiki, say so clearly.
5. If the answer is valuable, offer to save it as a new wiki page.

> Good answers should be filed back into the wiki so they compound over time.

---

## Lint / Audit

When the developer asks you to lint or audit the wiki:

- Check for contradictions between pages.
- Find orphan pages (no inbound links from other pages).
- Identify concepts mentioned in pages that lack their own page.
- Flag claims that may be outdated based on newer sources.
- Check that all pages follow the page format above.
- Report findings as a numbered list with suggested fixes.

---

## Rules

- **Never modify** anything in the `raw/` folder.
- Always update `wiki/index.md` and `wiki/log.md` after changes.
- Keep page names lowercase with hyphens (e.g. `machine-learning.md`).
- Write in clear, plain language.
- When uncertain about how to categorize something, ask the developer.
- Every entry in `wiki/log.md` must have a **timestamp (YYYY-MM-DD HH:MM)** and responsible party (GSD immutability standard).

---

## Integration with the GSD Workflow

| GSD Event | Wiki Action |
| :--- | :--- |
| Architectural decision recorded in `CONTEXT.md` | Create/update a page in `wiki/` with the decision |
| Critical bug resolved, documented in `TESTS.md` | Add a known-bug page in `wiki/` |
| New ROADMAP stage started | Update `wiki/index.md` with stage context |
| Delivery recorded in `VERSIONS.md` | Add an entry to `wiki/log.md` |

---

## Version History (Immutable)
- **[YYYY-MM-DD HH:MM] - Responsible:** Initial wiki template version.
- [Add new entries here without deleting.]
