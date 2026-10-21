<div align="center">

# 🧩 mcp-prompt-server-go

### 把你的 prompt 库变成 MCP 工具——热重载、可版本化、完全属于你。

一个 Go 写的单二进制 [Model Context Protocol](https://modelcontextprotocol.io) 服务器：`prompts/` 目录里每个 YAML/JSON 文件都会变成一个可调用的 MCP 工具。内置 19 个即用模板 + 2 个管理工具，总共 1,021 行 Go、2 个直接依赖。改完文件即热重载——不用重启，不用重新部署。

[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg?style=flat-square)](go.mod)
[![Transport](https://img.shields.io/badge/MCP-stdio-blueviolet.svg?style=flat-square)](#-客户端接入)
[![Prompts](https://img.shields.io/badge/prompts-19-green.svg?style=flat-square)](prompts/)
[![Deps](https://img.shields.io/badge/直接依赖-2-orange.svg?style=flat-square)](go.mod)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](../../pulls)

[快速开始](#-快速开始) · [内置模板](#-内置-prompt-模板) · [客户端接入](#-客户端接入) · [English](README.md)

</div>

## 你最好的那条 prompt，正卡在某个聊天框里

- 你写过一条绝妙的代码审查 prompt——现在它埋在三月的某段对话里。
- 同一段 500 字指令，Cursor 贴一遍、Raycast 贴一遍，然后各改各的，慢慢漂移。
- 每改一次 prompt 就重启一次客户端，思路全断。

这个服务器把 prompt 变成**文件**：可 diff、可版本化、可分享——任何 MCP 客户端都能调用。

```text
$ echo '{"jsonrpc":"2.0","id":1,"method":"tools/call",
          "params":{"name":"get_prompt_names","arguments":{}}}' | ./bin/mcp-prompt-server

"可用的prompts (19):
- code_review
- gen_html_web_page
- writing_assistant
- ..."
```

> 首次启动即有 21 个工具：19 个 prompt 模板 + `reload_prompts` + `get_prompt_names`。

## ✨ 核心能力

| 你做 | 它做 |
|---|---|
| 把 `my_prompt.yaml` 丢进 `prompts/` | 立刻注册为带参数 schema 的 MCP 工具——fsnotify 热重载，零重启 |
| 在任意 MCP 客户端调用该工具 | 用你的参数渲染 `{{arg}}` 占位符，返回最终 prompt 文本 |
| 手滑保存了一个坏 YAML | 告警跳过，其余 18 个照常服务——不崩 |
| 多个客户端同时连 | RWMutex 保护的管理器，并发读安全 |
| 调用 `get_prompt_names` | 返回当前实际加载的清单，调试/发现两相宜 |

**刻意无聊的设计：** stdio 传输、JSON-RPC、无守护进程、无数据库、无网络。`ps` 里只有一个进程，`kill` 就是关闭按钮。

## 🚀 快速开始

**1. 构建**（需要 Go 1.21+）

```bash
git clone https://github.com/anthonyhann/mcp-prompt-server-go
cd mcp-prompt-server-go && make build
```

**2. 注册**：把二进制作为 stdio MCP server 加进客户端（配置见下）。

**3. 调用**：按名字用任意 prompt——例如 `code_review` 传 `{"code": "..."}`。

**给 AI Agent 的口令**——直接粘贴给你的 agent：

> 克隆 https://github.com/anthonyhann/mcp-prompt-server-go，运行 `make build`，然后把 `<仓库>/bin/mcp-prompt-server` 作为 stdio MCP server 注册进我的客户端配置，并用 `get_prompt_names` 工具验证。

## 🔌 客户端接入

**Cursor** — `~/.cursor/mcp_config.json`：

```json
{
  "mcpServers": {
    "prompt-server": {
      "command": "/绝对路径/mcp-prompt-server-go/bin/mcp-prompt-server"
    }
  }
}
```

**Windsurf** — `~/.codeium/windsurf/mcp_config.json`：同上格式。

**Raycast** — "Install Server (MCP)"：Name 填 `prompt`，Command 填二进制绝对路径，Arguments 留空。

验证任意配置：

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | ./bin/mcp-prompt-server
```

## 📚 内置 prompt 模板

5 大类 19 个模板——是起点，不是陈列品，直接原地改成你的：

| 类别 | 模板 |
|---|---|
| 内容创作 | `gen_title` · `gen_mimeng_headline_title` · `gen_summarize` · `writing_assistant` · `gen_podcast_script` |
| 网页生成 | `gen_html_web_page` · `gen_3d_webpage_html` · `gen_3d_edu_webpage_html` · `gen_bento_grid_html` · `gen_knowledge_card_html` · `gen_magazine_card_html` |
| 产品开发 | `gen_prd_prototype_html` · `project_architecture` · `api_documentation` |
| 代码相关 | `code_review` · `code_refactoring` · `test_case_generator` · `build_mcp_server` |
| 元工具 | `prompt_template_generator` |

加自己的模板只是一个 YAML 的事：

```yaml
name: my_prompt
description: 这个 prompt 做什么
arguments:
  - name: input_text
    description: 输入文本
    required: true
messages:
  - role: user
    content:
      type: text
      text: "请处理：{{input_text}}"
```

保存 → 文件监控自动重载 → 工具即刻可用。

## 📁 目录结构

```text
mcp-prompt-server-go/
├── main.go              # 入口：加载 prompts → 注册工具 → stdio 服务
├── internal/
│   ├── mcp/             # MCP 协议数据模型
│   ├── prompt/          # 模板模型 + 管理器（fsnotify 热重载、RWMutex）
│   └── server/          # stdio JSON-RPC 服务器
├── prompts/             # 19 个内置 YAML 模板（改的就是这里）
├── tools/test_mcp.go    # 手动测试工具
├── Makefile             # build / dev / test / package
└── GUIDE.md · UPGRADE.md
```

## ⚠️ 什么时候别用它

- **你想要 prompt 应用商店。** 它只服务*你自己的*文件；19 个内置模板是起步套件，不是内容流。
- **你要远程访问。** 只有 stdio——没有 HTTP/SSE、没有鉴权、没有多租户。一台机器，一个用户。
- **你的客户端只支持 MCP `prompts` 能力。** 本服务器把 prompt 暴露为 **tools**（`tools/call`）——主流客户端都支持，但请知晓这个区别。
- **你要 prompt 编排/工作流。** `{{var}}` 替换就是全部的模板引擎，故意的。

## 🙋 为什么重写它

Node.js 原版（[gdli6177/mcp-prompt-server](https://github.com/gdli6177/mcp-prompt-server)）验证了想法；我想要的是一个能随手丢到任何机器上的静态二进制，不需要装运行时。于是有了这 1,021 行 Go、2 个直接依赖，`make build` 出活。我自己每天在用——每个内置模板都经过真实使用淘汰。

⭐ 如果它每天帮你省一次复制粘贴，一个 star 告诉我继续补模板。

## License

[MIT](LICENSE) © 2026 anthonyhann — 协议：[Model Context Protocol](https://modelcontextprotocol.io)。致谢启发本次重写的 Node.js 原版作者。
