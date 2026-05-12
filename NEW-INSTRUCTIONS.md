# NEW-INSTRUCTIONS

> **Status:** Do not write new instructions directly in this file.
> New instructions must be sent as a JSON handoff to `[CEO]`.
> Copy the template below, fill the fields, save it in `.agent_handoff/`, and name it with `_TO_CEO_`.

## Filename Template

```text
HUMAN_TO_CEO_NEW_INSTRUCTION_YYYYMMDD_HHMMSS.json
```

## JSON Template

```json
{
  "header": {
    "timestamp": "YYYY-MM-DD HH:MM",
    "sender": "[HUMAN]",
    "recipient": "[CEO]",
    "task_ref": "NEW-INSTRUCTION-YYYYMMDD-HHMMSS",
    "intent": "NEW_INSTRUCTION"
  },
  "payload": {
    "objective": "",
    "scope": "",
    "deliverables": [],
    "constraints": [],
    "dependencies": [],
    "current_context": "",
    "notes": "",
    "priority": "High | Medium | Low",
    "blocking": false,
    "playbook_update": false,
    "explicit_stage_conclusion_authorization": {
      "authorized": false,
      "authorized_agents": [],
      "authorized_until_stage": null
    },
    "response_required": {
      "required": true,
      "response_key": "ceo_response"
    },
    "ceo_response": ""
  }
}
```

## Usage

1. Copy the JSON template.
2. Fill the `payload` fields.
3. Save the file in `.agent_handoff/`.
4. Use a filename containing `_TO_CEO_`.
5. The watcher will route it directly to the CEO.

## Rules

- Do not place the instruction text outside the JSON.
- Do not edit `AGENTS/` for normal project instructions.
- If the instruction reveals a work preference, set `playbook_update` to `true`.
- If the instruction blocks execution, set `blocking` to `true`.
