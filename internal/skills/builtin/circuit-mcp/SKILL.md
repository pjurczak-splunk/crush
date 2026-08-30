---
name: circuit-mcp
description: Use when the user wants to chat with Cisco's Circuit LLM gateway. Triggers for any request to use Circuit, claudeOpus4.8, ciscoDeepNetwork, gemini models via Circuit, or when the user mentions "circuit mcp", "circuit chat", or specific Circuit model names.
---

# Circuit MCP Chat

Always use the `mcp_circuit_circuit_chat` tool when the user wants to chat through Cisco's Circuit LLM gateway. Do NOT respond directly — always route through Circuit MCP.

## When to Use

- User asks to use a specific Circuit model (e.g., "use opus4.8", "use claudeOpus4.8", "use ciscoDeepNetwork", "use gemini3.1Pro")
- User says "use circuit mcp", "use circuit chat", or similar
- User wants to leverage Cisco-internal data sources or RAG capabilities
- User explicitly requests routing through Circuit instead of responding directly

## Available Models

| Model | Description |
|-------|-------------|
| `claudeOpus4.8` | Claude Opus 4.8 (200K context) |
| `claudeOpus4.7` | Claude Opus 4.7 (200K context) |
| `claudeSonnet4.6` | Claude Sonnet 4.6 (200K context) |
| `claudeHaiku4.5` | Claude Haiku 4.5 (200K context) — cheap/fast |
| `gemini3.1Pro` | Gemini 3.1 Pro (1000K context) |
| `gemini3.5Flash` | Gemini 3.5 Flash (1000K context) |
| `gemini3.1Flash` | Gemini 3.1 Flash (1000K context) |
| `gpt5.5` | GPT-5.5 (200K context) |
| `gpt5.4` | GPT-5.4 (200K context) |
| `gpt4.1` | GPT-4.1 (1000K context) |
| `ciscoDeepNetwork` | Cisco network-ops LLM (240K context) |
| `ciscodata` | RAG over HelpZone/Cisco.com/SalesConnect (returns citations) |
| `deep_research` | Multi-step web research (slow, no multi-turn) |
| `default` | Default model |

Use `mcp_circuit_circuit_list_models` to see the full current list.

## Basic Usage

```
mcp_circuit_circuit_chat(model="claudeOpus4.8", prompt="Your question here")
```

## Key Parameters

- `model`: Model ID (see table above). Default is `"default"`.
- `prompt`: The user's question or instruction.
- `new_session`: Set to `true` to start a fresh conversation (discard prior context). Default is `false`.
- `session_id`: Continue an existing session by providing its ID.
- `attach_files`: List of absolute file paths to upload and attach.
- `attach_file_ids`: Previously-uploaded file IDs to reuse without re-uploading.
- `preview_tokens`: Truncate response over this many tokens to disk; preview is returned. Set to `0` to disable. Default is `1000`.
- `auto_summarize`: If truncated, run a free Circuit summarization pass. Default is `false`.
- `deep_research`: Equivalent to `model='deep_research'`. Slow (1-10 min). Default is `false`.
- `agent_id`: Optional Circuit agent UUID. Use `mcp_circuit_circuit_list_agents` to discover agents.
- `agent_name`: Optional agent display name.

## File Uploads

For large files or files reused across multiple chats:

1. Upload once: `mcp_circuit_circuit_upload_file(path="/absolute/path/to/file.json")`
2. Reuse the returned ID: `mcp_circuit_circuit_chat(..., attach_file_ids=["returned_id"])`

Per-model token caps are enforced server-side. The upload tool reports exact token count upfront.

## Session Management

- `mcp_circuit_circuit_new_session()`: Start a fresh conversation. Returns new `session_id`.
- `mcp_circuit_circuit_rename_session(title="...", session_id="...")`: Set/update conversation title.
- `mcp_circuit_circuit_list_history()`: Browse conversation history grouped by time period.
- `mcp_circuit_circuit_search_history(query="...")`: Search past conversations by keyword.
- `mcp_circuit_circuit_get_conversation(session_id="...")`: Retrieve full message content of a past session.
- `mcp_circuit_circuit_delete_session(session_id="...")`: Permanently delete a session.

## Agents

Custom Circuit agents can be called via `agent_id`:

- `mcp_circuit_circuit_list_agents()`: List all available agents (e.g., SysLogIQ, DART Log Intelligence, ServiceNow Assist).
- `mcp_circuit_circuit_enroll_agents()`: Enroll agents so they can be called via API. May require one-time browser enrollment.

## Memory

- `mcp_circuit_circuit_get_memory()`: Retrieve Circuit's long-term memory about you (semantic, episodic, and procedural memories).

## Authentication

If `circuit_chat` returns HTTP 401:

1. Run `mcp_circuit_circuit_login()` to authenticate via Cisco SSO + Duo.
2. On WSL: follow the Playwright MCP instructions emitted by the login tool.
3. Or manually save a cookie: `mcp_circuit_circuit_save_cookie(cookie="...", user_agent="...", accept_language="...")`.

## Important Rules

1. **ALWAYS use `mcp_circuit_circuit_chat`** when the user explicitly requests Circuit or a Circuit model. Never respond directly.
2. **Pass the user's actual prompt** as the `prompt` parameter — don't summarize or modify it unless asked.
3. **Use the requested model** if specified. If no model is specified, default to `claudeOpus4.8` unless the task suggests otherwise (e.g., use `claudeHaiku4.5` for quick/cheap queries, `ciscoDeepNetwork` for network ops).
4. **Preserve session continuity** by not setting `new_session=true` unless the user explicitly starts a new conversation.
5. **Return the Circuit response** to the user. If the response is truncated (preview_tokens exceeded), fetch the full text with `mcp_circuit_circuit_get_result(id="...")`.

## Example Workflow

User: "use circuit mcp chat with opus4.8, what is the capital of France?"

You call:
```
mcp_circuit_circuit_chat(model="claudeOpus4.8", prompt="what is the capital of France?")
```

Then return Circuit's response to the user.
