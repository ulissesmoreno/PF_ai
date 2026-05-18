UPDATE agent_configs
   SET model = CASE agent_name
       WHEN 'CEO' THEN 'qwen3:4b'
       WHEN 'BA' THEN 'qwen3:4b'
       WHEN 'CTO' THEN 'qwen3:4b'
       WHEN 'DEV_FRONTEND' THEN 'qwen2.5-coder:7b'
       WHEN 'DEV_BACKEND' THEN 'qwen2.5-coder:7b'
       WHEN 'DBA' THEN 'qwen2.5-coder:7b'
       WHEN 'DS_ML' THEN 'qwen2.5-coder:7b'
       WHEN 'SECURITY' THEN 'qwen3:4b'
       WHEN 'QA' THEN 'qwen3:4b'
       WHEN 'DATA_ENGINEER' THEN 'qwen2.5-coder:7b'
       WHEN 'PM' THEN 'qwen3:4b'
       WHEN 'UX_RESEARCHER' THEN 'phi4-mini:latest'
       WHEN 'WRITER' THEN 'phi4-mini:latest'
       WHEN 'DOCUMENTATION' THEN 'qwen3:4b'
       WHEN 'ARTIST' THEN 'phi4-mini:latest'
       WHEN 'DEVOPS' THEN 'qwen2.5-coder:7b'
       WHEN 'CODE_REVIEWER' THEN 'qwen3:4b'
       WHEN 'CMO' THEN 'phi4-mini:latest'
       ELSE model
       END,
       tier = CASE agent_name
       WHEN 'DEV_FRONTEND' THEN 2
       WHEN 'DEV_BACKEND' THEN 2
       WHEN 'DBA' THEN 2
       WHEN 'DS_ML' THEN 2
       WHEN 'DATA_ENGINEER' THEN 2
       WHEN 'DEVOPS' THEN 2
       WHEN 'UX_RESEARCHER' THEN 1
       WHEN 'WRITER' THEN 1
       WHEN 'ARTIST' THEN 1
       WHEN 'CMO' THEN 1
       ELSE 3
       END,
       active = 1,
       updated_at = CURRENT_TIMESTAMP,
       updated_by = 'migration-008-local-profile'
 WHERE agent_name IN (
       'CEO', 'BA', 'CTO', 'DEV_FRONTEND', 'DEV_BACKEND', 'DBA', 'DS_ML',
       'SECURITY', 'QA', 'DATA_ENGINEER', 'PM', 'UX_RESEARCHER', 'WRITER',
       'DOCUMENTATION', 'ARTIST', 'DEVOPS', 'CODE_REVIEWER', 'CMO'
 );
