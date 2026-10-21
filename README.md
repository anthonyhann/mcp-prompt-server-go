<div align="center">

# 🧩 mcp-prompt-server-go

### Your prompt library, served as MCP tools — hot-reloaded, version-controlled, yours.

A single-binary [Model Context Protocol](https://modelcontextprotocol.io) server written in Go: every YAML/JSON file in `prompts/` becomes a callable MCP tool. 19 ready-made templates + 2 management tools, 1,021 lines of Go, 2 direct dependencies. Edit a file, the server reloads — no restart, no redeploy.

[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg?style=flat-square)](go.mod)
[![Transport](https://img.shields.io/badge/MCP-stdio-blueviolet.svg?style=flat-square)](#-client-setup)
[![Prompts](https://img.shields.io/badge/prompts-19-green.svg?style=flat-square)](prompts/)
[![Deps](https://img.shields.io/badge/direct_deps-2-orange.svg?style=flat-square)](go.mod)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](../../pulls)

[Quickstart](#-quickstart) · [Prompts](#-built-in-prompts) · [Client setup](#-client-setup) · [中文](README_zh.md)

</div>

## Your best prompt is stuck in a chat box

- You wrote a killer code-review prompt once. Now it's buried in a conversation from March.
- You paste the same 500-word instruction into Cursor, then Raycast, then wherever — and they drift apart.
- You tweak a prompt, restart the client, lose your place, repeat.

This server makes prompts **files**: diffable, versionable, shareable — and callable from any MCP client.

```text
$ echo '{"jsonrpc":"2.0","id":1,"method":"tools/call",
          "params":{"name":"get_prompt_names","arguments":{}}}' | ./bin/mcp-prompt-server

"可用的prompts (19):
- code_review
- gen_html_web_page
- writing_assistant
- ..."
```

> 21 tools on first boot: 19 prompt templates from `prompts/` + `reload_prompts` + `get_prompt_names`.

## ✨ What it does

| You do | It does |
|---|---|
| Drop `my_prompt.yaml` into `prompts/` | Registers `my_prompt` as an MCP tool with typed arguments — instantly, via fsnotify hot reload |
| Call the tool from any MCP client | Renders `{{arg}}` placeholders with your arguments, returns the final prompt text |
| Save a broken YAML by accident | Skips it with a warning, keeps serving the other 18 — no crash |
| Run many clients at once | RWMutex-guarded manager; concurrent reads are safe |
| Ask `get_prompt_names` | Live list of what's loaded, for debugging and discovery |

**Deliberately boring:** stdio transport, JSON-RPC, no daemon, no database, no network. `ps` shows one process; `kill` is the off switch.

## 🚀 Quickstart

**1. Build** (Go 1.21+)

```bash
git clone https://github.com/anthonyhann/mcp-prompt-server-go
cd mcp-prompt-server-go && make build
```

**2. Register** the binary as a stdio MCP server in your client (full configs below).

**3. Call** any prompt by name — e.g. `code_review` with `{"code": "..."}`.

**For AI Agents** — paste this to your agent:

> Clone https://github.com/anthonyhann/mcp-prompt-server-go, run `make build`, then register `<repo>/bin/mcp-prompt-server` as a stdio MCP server in my client config, and verify with the `get_prompt_names` tool.

## 🔌 Client setup

**Cursor** — `~/.cursor/mcp_config.json`:

```json
{
  "mcpServers": {
    "prompt-server": {
      "command": "/absolute/path/to/mcp-prompt-server-go/bin/mcp-prompt-server"
    }
  }
}
```

**Windsurf** — `~/.codeium/windsurf/mcp_config.json`: same shape.

**Raycast** — "Install Server (MCP)": Name `prompt`, Command = absolute binary path, no arguments.

Verify any setup:

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | ./bin/mcp-prompt-server
```

## 📚 Built-in prompts

19 templates in 5 families — starting points, meant to be edited in place:

| Family | Prompts |
|---|---|
| Content | `gen_title` · `gen_mimeng_headline_title` · `gen_summarize` · `writing_assistant` · `gen_podcast_script` |
| Web pages | `gen_html_web_page` · `gen_3d_webpage_html` · `gen_3d_edu_webpage_html` · `gen_bento_grid_html` · `gen_knowledge_card_html` · `gen_magazine_card_html` |
| Product | `gen_prd_prototype_html` · `project_architecture` · `api_documentation` |
| Code | `code_review` · `code_refactoring` · `test_case_generator` · `build_mcp_server` |
| Meta | `prompt_template_generator` |

Your own template is a YAML file away:

```yaml
name: my_prompt
description: What it does
arguments:
  - name: input_text
    description: The input
    required: true
messages:
  - role: user
    content:
      type: text
      text: "Process this: {{input_text}}"
```

Save → the file watcher reloads → the tool is live.

## 📁 Repository layout

```text
mcp-prompt-server-go/
├── main.go              # Entry: load prompts → register tools → serve stdio
├── internal/
│   ├── mcp/             # MCP protocol models
│   ├── prompt/          # Template models + manager (fsnotify hot reload, RWMutex)
│   └── server/          # stdio JSON-RPC server
├── prompts/             # 19 built-in YAML templates (edit these!)
├── tools/test_mcp.go    # Manual test harness
├── Makefile             # build / dev / test / package
└── GUIDE.md · UPGRADE.md
```

## ⚠️ When NOT to use it

- **You want a prompt marketplace.** This serves *your* files; the 19 built-ins are a starter kit, not a feed.
- **You need remote access.** stdio only — no HTTP/SSE, no auth, no multi-tenancy. One user, one machine.
- **Your client only supports the MCP `prompts` capability.** This server exposes prompts as **tools** (`tools/call`), which every major client supports — but know the difference.
- **You need prompt chaining/workflows.** `{{var}}` substitution is the whole templating engine, on purpose.

## 🙋 Why it exists

The Node.js original ([gdli6177/mcp-prompt-server](https://github.com/gdli6177/mcp-prompt-server)) proved the idea; I wanted the same thing as one static binary I could drop anywhere without a runtime. So: 1,021 lines of Go, 2 direct deps, `make build` → done. I run it daily; every built-in prompt here survived real use.

⭐ If it saves you one copy-paste a day, a star tells me to keep the templates coming.

## License

[MIT](LICENSE) © 2026 anthonyhann — protocol: [Model Context Protocol](https://modelcontextprotocol.io). Hat tip to the Node.js originals that inspired this rewrite.
