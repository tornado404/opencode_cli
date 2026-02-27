# oho API 映射参考

本文档列出所有 oho 命令与 OpenCode Server API 端点的映射关系。

## 全局命令 (Global)

| 命令 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| `oho global health` | `/global/health` | GET | 检查服务器健康状态 |
| `oho global event` | `/global/event` | GET (SSE) | 监听全局事件流 |

## 项目命令 (Project)

| 命令 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| `oho project list` | `/project` | GET | 列出所有项目 |
| `oho project current` | `/project/current` | GET | 获取当前项目 |
| `oho path` | `/path` | GET | 获取当前路径 |
| `oho vcs` | `/vcs` | GET | 获取 VCS 信息 |
| `oho instance dispose` | `/instance/dispose` | POST | 销毁当前实例 |

## 配置命令 (Config)

| 命令 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| `oho config get` | `/config` | GET | 获取配置 |
| `oho config set` | `/config` | PATCH | 更新配置 |
| `oho config providers` | `/config/providers` | GET | 列出提供商和默认模型 |

## 提供商命令 (Provider)

| 命令 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| `oho provider list` | `/provider` | GET | 列出所有提供商 |
| `oho provider auth` | `/provider/auth` | GET | 获取提供商认证方式 |
| `oho provider oauth authorize <id>` | `/provider/{id}/oauth/authorize` | POST | OAuth 授权 |
| `oho provider oauth callback <id>` | `/provider/{id}/oauth/callback` | POST | 处理 OAuth 回调 |

## 会话命令 (Session)

| 命令 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| `oho session list` | `/session` | GET | 列出所有会话 |
| `oho session create` | `/session` | POST | 创建新会话 |
| `oho session status` | `/session/status` | GET | 获取所有会话状态 |
| `oho session get <id>` | `/session/:id` | GET | 获取会话详情 |
| `oho session delete <id>` | `/session/:id` | DELETE | 删除会话 |
| `oho session update <id>` | `/session/:id` | PATCH | 更新会话属性 |
| `oho session children <id>` | `/session/:id/children` | GET | 获取子会话 |
| `oho session todo <id>` | `/session/:id/todo` | GET | 获取待办事项 |
| `oho session init <id>` | `/session/:id/init` | POST | 创建 AGENTS.md |
| `oho session fork <id>` | `/session/:id/fork` | POST | 分叉会话 |
| `oho session abort <id>` | `/session/:id/abort` | POST | 中止会话 |
| `oho session share <id>` | `/session/:id/share` | POST | 分享会话 |
| `oho session unshare <id>` | `/session/:id/share` | DELETE | 取消分享 |
| `oho session diff <id>` | `/session/:id/diff` | GET | 获取文件差异 |
| `oho session summarize <id>` | `/session/:id/summarize` | POST | 总结会话 |
| `oho session revert <id>` | `/session/:id/revert` | POST | 回退消息 |
| `oho session unrevert <id>` | `/session/:id/unrevert` | POST | 恢复回退 |
| `oho session permissions <id> <permId>` | `/session/:id/permissions/:permissionID` | POST | 响应权限请求 |

## 消息命令 (Message)

| 命令 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| `oho message list -s <session>` | `/session/:id/message` | GET | 列出消息 |
| `oho message add -s <session> "内容"` | `/session/:id/message` | POST | 发送消息 |
| `oho message get -s <session> <msgId>` | `/session/:id/message/:messageID` | GET | 获取消息详情 |
| `oho message prompt-async -s <session>` | `/session/:id/prompt_async` | POST | 异步发送消息 |
| `oho message command -s <session>` | `/session/:id/command` | POST | 执行斜杠命令 |
| `oho message shell -s <session>` | `/session/:id/shell` | POST | 运行 shell 命令 |

## 文件命令 (File)

| 命令 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| `oho file list [path]` | `/file` | GET | 列出文件 |
| `oho file content <path>` | `/file/content` | GET | 读取文件内容 |
| `oho file status` | `/file/status` | GET | 获取文件状态 |

## 查找命令 (Find)

| 命令 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| `oho find text <pattern>` | `/find` | GET | 搜索文本 |
| `oho find file <query>` | `/find/file` | GET | 查找文件 |
| `oho find symbol <query>` | `/find/symbol` | GET | 查找符号 |

## 工具命令 (Tool)

| 命令 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| `oho tool ids` | `/experimental/tool/ids` | GET | 列出工具 ID |
| `oho tool list` | `/experimental/tool` | GET | 列出工具 |

## 其他命令

| 命令 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| `oho agent list` | `/agent` | GET | 列出代理 |
| `oho command list` | `/command` | GET | 列出命令 |
| `oho lsp status` | `/lsp` | GET | LSP 状态 |
| `oho formatter status` | `/formatter` | GET | 格式化器状态 |
| `oho mcp list` | `/mcp` | GET | MCP 服务器列表 |
| `oho mcp add` | `/mcp` | POST | 添加 MCP 服务器 |

## TUI 控制命令

| 命令 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| `oho tui append-prompt` | `/tui/append-prompt` | POST | 追加提示词 |
| `oho tui open-help` | `/tui/open-help` | POST | 打开帮助 |
| `oho tui open-sessions` | `/tui/open-sessions` | POST | 打开会话选择器 |
| `oho tui open-themes` | `/tui/open-themes` | POST | 打开主题选择器 |
| `oho tui open-models` | `/tui/open-models` | POST | 打开模型选择器 |
| `oho tui submit-prompt` | `/tui/submit-prompt` | POST | 提交提示词 |
| `oho tui clear-prompt` | `/tui/clear-prompt` | POST | 清除提示词 |
| `oho tui execute-command` | `/tui/execute-command` | POST | 执行命令 |
| `oho tui show-toast` | `/tui/show-toast` | POST | 显示提示 |
| `oho tui control-next` | `/tui/control/next` | GET | 等待控制请求 |
| `oho tui control-response` | `/tui/control/response` | POST | 响应控制请求 |

## 认证命令

| 命令 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| `oho auth set <provider>` | `/auth/:id` | PUT | 设置认证凭据 |

## 文档命令

| 命令 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| `oho doc` | `/doc` | GET | 获取 OpenAPI 规范 |

---

## 通用标志

| 标志 | 简写 | 描述 | 默认值 |
|------|------|------|--------|
| `--host` | | 服务器主机地址 | `127.0.0.1` |
| `--port` | `-p` | 服务器端口 | `4096` |
| `--password` | | 服务器密码 | 环境变量 |
| `--json` | `-j` | JSON 格式输出 | `false` |

## 会话标志

| 标志 | 简写 | 描述 |
|------|------|------|
| `--session` | `-s` | 会话 ID (用于 message 等命令) |

## 环境变量

| 变量 | 描述 | 默认值 |
|------|------|--------|
| `OPENCODE_SERVER_HOST` | 服务器主机 | `127.0.0.1` |
| `OPENCODE_SERVER_PORT` | 服务器端口 | `4096` |
| `OPENCODE_SERVER_USERNAME` | 用户名 | `opencode` |
| `OPENCODE_SERVER_PASSWORD` | 密码 | 空 |
