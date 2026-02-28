# oho 功能完成情况

> 本文档记录 oho CLI 对 OpenCode Server API 的实现覆盖情况

## 完成指标

| 指标 | 数值 |
|------|------|
| 总 API 端点 | 66 |
| 已实现 | 64 |
| 未实现 | 2 |
| **覆盖率** | **97.0%** |

## 模块完成情况

### ✅ 已完成模块

| 模块 | API 端点数 | 已实现 | 状态 |
|------|-----------|--------|------|
| [Global](#global-全局命令) | 2 | 2 | ✅ |
| [Project](#project-项目管理) | 4 | 4 | ✅ |
| [Config](#config-配置管理) | 3 | 3 | ✅ |
| [Provider](#provider-提供商管理) | 4 | 4 | ✅ |
| [Session](#session-会话管理) | 18 | 18 | ✅ |
| [Message](#message-消息管理) | 6 | 6 | ✅ |
| [Command](#command-命令列表) | 1 | 1 | ✅ |
| [File](#file-文件操作) | 6 | 6 | ✅ |
| [Find](#find-查找功能) | 3 | 3 | ✅ |
| [Tool](#tool-工具列表) | 2 | 2 | ✅ |
| [LSP](#lsp-语言服务器) | 1 | 1 | ✅ |
| [Formatter](#formatter-格式化器) | 1 | 1 | ✅ |
| [MCP](#mcp-模型上下文协议) | 2 | 2 | ✅ |
| [Agent](#agent-代理列表) | 1 | 1 | ✅ |
| [TUI](#tui-终端界面控制) | 10 | 10 | ✅ |
| [Auth](#auth-认证管理) | 1 | 1 | ✅ |

### ❌ 未完成模块

| 模块 | API 端点 | 说明 |
|------|---------|------|
| Logging | `POST /log` | 日志写入功能（调试辅助） |
| Docs | `GET /doc` | OpenAPI 规范查看（可直接访问服务器 `/doc` 端点） |

---

## 详细 API 映射

### Global (全局命令)

| 命令 | API 端点 | 方法 | 状态 |
|------|----------|------|------|
| `oho global health` | `/global/health` | GET | ✅ |
| `oho global event` | `/global/event` | GET (SSE) | ✅ |

### Project (项目管理)

| 命令 | API 端点 | 方法 | 状态 |
|------|----------|------|------|
| `oho project list` | `/project` | GET | ✅ |
| `oho project current` | `/project/current` | GET | ✅ |
| `oho path` | `/path` | GET | ✅ |
| `oho vcs` | `/vcs` | GET | ✅ |
| `oho instance dispose` | `/instance/dispose` | POST | ✅ |

### Config (配置管理)

| 命令 | API 端点 | 方法 | 状态 |
|------|----------|------|------|
| `oho config get` | `/config` | GET | ✅ |
| `oho config set` | `/config` | PATCH | ✅ |
| `oho config providers` | `/config/providers` | GET | ✅ |

### Provider (提供商管理)

| 命令 | API 端点 | 方法 | 状态 |
|------|----------|------|------|
| `oho provider list` | `/provider` | GET | ✅ |
| `oho provider auth` | `/provider/auth` | GET | ✅ |
| `oho provider oauth authorize <id>` | `/provider/{id}/oauth/authorize` | POST | ✅ |
| `oho provider oauth callback <id>` | `/provider/{id}/oauth/callback` | POST | ✅ |

### Session (会话管理)

| 命令 | API 端点 | 方法 | 状态 |
|------|----------|------|------|
| `oho session list` | `/session` | GET | ✅ |
| `oho session create` | `/session` | POST | ✅ |
| `oho session status` | `/session/status` | GET | ✅ |
| `oho session get <id>` | `/session/:id` | GET | ✅ |
| `oho session delete <id>` | `/session/:id` | DELETE | ✅ |
| `oho session update <id>` | `/session/:id` | PATCH | ✅ |
| `oho session children <id>` | `/session/:id/children` | GET | ✅ |
| `oho session todo <id>` | `/session/:id/todo` | GET | ✅ |
| `oho session init <id>` | `/session/:id/init` | POST | ✅ |
| `oho session fork <id>` | `/session/:id/fork` | POST | ✅ |
| `oho session abort <id>` | `/session/:id/abort` | POST | ✅ |
| `oho session share <id>` | `/session/:id/share` | POST | ✅ |
| `oho session unshare <id>` | `/session/:id/share` | DELETE | ✅ |
| `oho session diff <id>` | `/session/:id/diff` | GET | ✅ |
| `oho session summarize <id>` | `/session/:id/summarize` | POST | ✅ |
| `oho session revert <id>` | `/session/:id/revert` | POST | ✅ |
| `oho session unrevert <id>` | `/session/:id/unrevert` | POST | ✅ |
| `oho session permissions <id>` | `/session/:id/permissions/:permissionID` | POST | ✅ |

### Message (消息管理)

| 命令 | API 端点 | 方法 | 状态 |
|------|----------|------|------|
| `oho message list -s <session>` | `/session/:id/message` | GET | ✅ |
| `oho message add -s <session>` | `/session/:id/message` | POST | ✅ |
| `oho message get -s <session> <msgId>` | `/session/:id/message/:messageID` | GET | ✅ |
| `oho message prompt-async -s <session>` | `/session/:id/prompt_async` | POST | ✅ |
| `oho message command -s <session>` | `/session/:id/command` | POST | ✅ |
| `oho message shell -s <session>` | `/session/:id/shell` | POST | ✅ |

### Command (命令列表)

| 命令 | API 端点 | 方法 | 状态 |
|------|----------|------|------|
| `oho command list` | `/command` | GET | ✅ |

### File (文件操作)

| 命令 | API 端点 | 方法 | 状态 |
|------|----------|------|------|
| `oho file list [path]` | `/file` | GET | ✅ |
| `oho file content <path>` | `/file/content` | GET | ✅ |
| `oho file status` | `/file/status` | GET | ✅ |

### Find (查找功能)

| 命令 | API 端点 | 方法 | 状态 |
|------|----------|------|------|
| `oho find text <pattern>` | `/find` | GET | ✅ |
| `oho find file <query>` | `/find/file` | GET | ✅ |
| `oho find symbol <query>` | `/find/symbol` | GET | ✅ |

### Tool (工具列表)

| 命令 | API 端点 | 方法 | 状态 |
|------|----------|------|------|
| `oho tool ids` | `/experimental/tool/ids` | GET | ✅ |
| `oho tool list` | `/experimental/tool` | GET | ✅ |

### LSP (语言服务器)

| 命令 | API 端点 | 方法 | 状态 |
|------|----------|------|------|
| `oho lsp status` | `/lsp` | GET | ✅ |

### Formatter (格式化器)

| 命令 | API 端点 | 方法 | 状态 |
|------|----------|------|------|
| `oho formatter status` | `/formatter` | GET | ✅ |

### MCP (模型上下文协议)

| 命令 | API 端点 | 方法 | 状态 |
|------|----------|------|------|
| `oho mcp list` | `/mcp` | GET | ✅ |
| `oho mcp add` | `/mcp` | POST | ✅ |

### Agent (代理列表)

| 命令 | API 端点 | 方法 | 状态 |
|------|----------|------|------|
| `oho agent list` | `/agent` | GET | ✅ |

### TUI (终端界面控制)

| 命令 | API 端点 | 方法 | 状态 |
|------|----------|------|------|
| `oho tui append-prompt` | `/tui/append-prompt` | POST | ✅ |
| `oho tui open-help` | `/tui/open-help` | POST | ✅ |
| `oho tui open-sessions` | `/tui/open-sessions` | POST | ✅ |
| `oho tui open-themes` | `/tui/open-themes` | POST | ✅ |
| `oho tui open-models` | `/tui/open-models` | POST | ✅ |
| `oho tui submit-prompt` | `/tui/submit-prompt` | POST | ✅ |
| `oho tui clear-prompt` | `/tui/clear-prompt` | POST | ✅ |
| `oho tui execute-command` | `/tui/execute-command` | POST | ✅ |
| `oho tui show-toast` | `/tui/show-toast` | POST | ✅ |
| `oho tui control-next` | `/tui/control/next` | GET | ✅ |
| `oho tui control-response` | `/tui/control/response` | POST | ✅ |

### Auth (认证管理)

| 命令 | API 端点 | 方法 | 状态 |
|------|----------|------|------|
| `oho auth set <provider>` | `/auth/:id` | PUT | ✅ |

---

## 未实现功能说明

### POST /log (日志写入)

```bash
# 未实现 - 可通过服务器直接访问
curl -X POST http://127.0.0.1:4096/log \
  -H "Content-Type: application/json" \
  -d '{"service": "app", "level": "info", "message": "test"}'
```

**说明**: 此功能为调试辅助功能，客户端可通过直接调用 API 使用。

### GET /doc (OpenAPI 规范)

```bash
# 未实现 - 可通过浏览器或 curl 直接访问
curl http://127.0.0.1:4096/doc
# 或浏览器访问 http://127.0.0.1:4096/doc
```

**说明**: OpenAPI 规范文档可直接通过浏览器访问查看，无需通过 CLI 调用。

---

## 更新日志

- **2026-02-28**: 初始版本，覆盖率 97.0% (64/66)
