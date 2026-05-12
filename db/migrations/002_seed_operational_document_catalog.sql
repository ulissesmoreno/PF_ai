INSERT OR IGNORE INTO documents (
    path,
    title,
    document_type,
    owner,
    created_at
) VALUES
    ('DOC/ARCHITECTURE.md', 'ARCHITECTURE', 'ARCHITECTURE', 'CTO', datetime('now')),
    ('DOC/CONTEXT.md', 'CONTEXT', 'CONTEXT', 'Management agents', datetime('now')),
    ('DOC/DESIGN.md', 'DESIGN', 'DESIGN', 'UX_RESEARCHER/DEV_FRONTEND', datetime('now')),
    ('DOC/ENV_SETUP.md', 'ENV_SETUP', 'ENV_SETUP', 'DEVOPS', datetime('now')),
    ('DOC/GSD-RULES.md', 'GSD-RULES', 'GSD_RULES', 'CEO/CTO', datetime('now')),
    ('DOC/PLAN.md', 'PLAN', 'PLAN', 'BA/CTO', datetime('now')),
    ('DOC/PROJECT.md', 'PROJECT', 'PROJECT', 'CEO', datetime('now')),
    ('DOC/RETROSPECTIVE.md', 'RETROSPECTIVE', 'RETROSPECTIVE', 'CEO', datetime('now')),
    ('DOC/ROADMAP.md', 'ROADMAP', 'ROADMAP', 'CEO/PM', datetime('now')),
    ('DOC/STATE.md', 'STATE', 'STATE', 'Technical agents', datetime('now')),
    ('DOC/TASKS.md', 'TASKS', 'TASKS', 'Technical agents', datetime('now')),
    ('DOC/TESTS.md', 'TESTS', 'TESTS', 'QA/SECURITY/Technical agents', datetime('now')),
    ('DOC/VERSIONS.md', 'VERSIONS', 'VERSIONS', 'Phase-closing agent', datetime('now')),
    ('DOC/WIKI.md', 'WIKI', 'WIKI_PROTOCOL', 'DOCUMENTATION', datetime('now')),
    ('PLAYBOOK.md', 'PLAYBOOK', 'PLAYBOOK', 'CEO', datetime('now'));
