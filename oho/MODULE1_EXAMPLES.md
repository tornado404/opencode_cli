# oho CLI 模块 1：客户端初始化 - 命令示例与预期输出

本文档提供 oho CLI 客户端初始化模块的完整命令示例和真实输出示例。

---

## 目录

1. [`oho auth set`](#1-oho-auth-set)
2. [`oho config get`](#2-oho-config-get)
3. [`oho config set`](#3-oho-config-set)
4. [`oho config providers`](#4-oho-config-providers)
5. [`oho session list`](#5-oho-session-list)

---

## 1. oho auth set

设置 AI 提供商的认证凭据。

### 命令语法

```bash
oho auth set <provider> --credentials <key=value> [--credentials <key=value>...]
```

### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `provider` | 位置参数 | 是 | 提供商 ID（如 `openai`, `anthropic`, `azure`） |
| `--credentials` | 字符串数组 | 是 | 认证凭据，格式为 `key=value`，可多次指定 |

### 使用示例

#### 示例 1：设置 OpenAI API Key

```bash
oho auth set openai --credentials apiKey=sk-abc123xyz789
```

**预期输出：**
```
提供商 openai 的认证凭据已设置
```

#### 示例 2：设置 Azure OpenAI 多凭据

```bash
oho auth set azure \
  --credentials apiKey=azure-key-12345 \
  --credentials endpoint=https://my-resource.openai.azure.com \
  --credentials deployment=my-deployment
```

**预期输出：**
```
提供商 azure 的认证凭据已设置
```

#### 示例 3：设置 Anthropic API Key

```bash
oho auth set anthropic --credentials apiKey=sk-ant-abc123
```

**预期输出：**
```
提供商 anthropic 的认证凭据已设置
```

#### 示例 4：JSON 输出模式

```bash
oho auth set openai --credentials apiKey=sk-test123 --json
```

**预期输出：**
```json
true
```

#### 示例 5：错误情况 - 缺少凭据

```bash
oho auth set openai
```

**预期输出：**
```
Error: required flag(s) "credentials" not set
```

#### 示例 6：错误情况 - 无效的凭据格式

```bash
oho auth set openai --credentials invalid-format
```

**预期输出：**
```
Error: 请提供至少一个凭据 (--credentials key=value)
```

---

## 2. oho config get

获取当前 OpenCode 配置。

### 命令语法

```bash
oho config get [--json]
```

### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `--json` | 布尔标志 | 否 | 以 JSON 格式输出配置 |

### 使用示例

#### 示例 1：默认文本输出

```bash
oho config get
```

**预期输出：**
```
当前配置:
  默认模型：gpt-4
  主题：dark
  语言：zh-CN
  最大 Token：4096
  温度：0.7
  自动批准：read, ls, cat
```

#### 示例 2：JSON 输出模式

```bash
oho config get --json
```

**预期输出：**
```json
{
  "defaultModel": "gpt-4",
  "theme": "dark",
  "language": "zh-CN",
  "maxTokens": 4096,
  "temperature": 0.7,
  "autoApprove": ["read", "ls", "cat"]
}
```

#### 示例 3：空配置情况

```bash
oho config get
```

**预期输出：**
```
当前配置:
  默认模型：
  主题：
  语言：
  最大 Token：0
  温度：0.00
```

#### 示例 4：带自动批准列表的配置

```bash
oho config get
```

**预期输出：**
```
当前配置:
  默认模型：claude-3-5-sonnet
  主题：light
  语言：en-US
  最大 Token：8192
  温度：0.5
  自动批准：read, ls, cat, grep, find
```

---

## 3. oho config set

更新 OpenCode 配置项。

### 命令语法

```bash
oho config set [flags]
```

### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `--theme` | 字符串 | 否 | 主题名称（如 `dark`, `light`） |
| `--language` | 字符串 | 否 | 语言设置（如 `zh-CN`, `en-US`） |
| `--model` | 字符串 | 否 | 默认模型名称 |
| `--max-tokens` | 整数 | 否 | 最大 Token 数 |
| `--temperature` | 浮点数 | 否 | 温度参数（0.0-2.0） |
| `--auto-approve` | 字符串数组 | 否 | 自动批准的工具列表 |

### 使用示例

#### 示例 1：更新主题

```bash
oho config set --theme dark
```

**预期输出：**
```
配置已更新
```

#### 示例 2：更新多个配置项

```bash
oho config set \
  --theme dark \
  --language zh-CN \
  --model gpt-4 \
  --max-tokens 8192 \
  --temperature 0.7
```

**预期输出：**
```
配置已更新
```

#### 示例 3：设置自动批准工具

```bash
oho config set --auto-approve read,ls,cat,grep
```

**预期输出：**
```
配置已更新
```

#### 示例 4：仅更新温度参数

```bash
oho config set --temperature 0.9
```

**预期输出：**
```
配置已更新
```

#### 示例 5：错误情况 - 未提供任何配置项

```bash
oho config set
```

**预期输出：**
```
Error: 请提供至少一个要更新的配置项
```

#### 示例 6：JSON 输出模式验证更新

```bash
oho config set --theme dark --json
oho config get --json
```

**预期输出：**
```
配置已更新
{
  "defaultModel": "gpt-4",
  "theme": "dark",
  "language": "zh-CN",
  "maxTokens": 4096,
  "temperature": 0.7,
  "autoApprove": []
}
```

---

## 4. oho config providers

列出所有可用的 AI 提供商和默认模型配置。

### 命令语法

```bash
oho config providers [--json]
```

### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `--json` | 布尔标志 | 否 | 以 JSON 格式输出 |

### 使用示例

#### 示例 1：默认文本输出

```bash
oho config providers
```

**预期输出：**
```
可用提供商:
  - OpenAI (openai)
  - Anthropic (anthropic)
  - Azure OpenAI (azure)
  - Google Gemini (google)
  - Ollama (ollama)

默认模型:
  openai: gpt-4
  anthropic: claude-3-5-sonnet
  azure: gpt-35-turbo
  google: gemini-pro
  ollama: llama2
```

#### 示例 2：JSON 输出模式

```bash
oho config providers --json
```

**预期输出：**
```json
{
  "providers": [
    {
      "id": "openai",
      "name": "OpenAI",
      "authType": "apiKey"
    },
    {
      "id": "anthropic",
      "name": "Anthropic",
      "authType": "apiKey"
    },
    {
      "id": "azure",
      "name": "Azure OpenAI",
      "authType": "multi"
    },
    {
      "id": "google",
      "name": "Google Gemini",
      "authType": "apiKey"
    },
    {
      "id": "ollama",
      "name": "Ollama",
      "authType": "none"
    }
  ],
  "default": {
    "openai": "gpt-4",
    "anthropic": "claude-3-5-sonnet",
    "azure": "gpt-35-turbo",
    "google": "gemini-pro",
    "ollama": "llama2"
  }
}
```

#### 示例 3：仅部分提供商配置

```bash
oho config providers
```

**预期输出：**
```
可用提供商:
  - OpenAI (openai)
  - Ollama (ollama)

默认模型:
  openai: gpt-3.5-turbo
  ollama: codellama
```

#### 示例 4：无默认模型配置

```bash
oho config providers
```

**预期输出：**
```
可用提供商:
  - OpenAI (openai)
  - Anthropic (anthropic)
```

---

## 5. oho session list

列出所有 OpenCode 会话。

### 命令语法

```bash
oho session list [--json] [--running]
```

### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `--json` | 布尔标志 | 否 | 以 JSON 格式输出 |
| `--running` | 布尔标志 | 否 | 只显示正在运行的会话 |

### 使用示例

#### 示例 1：默认文本输出

```bash
oho session list
```

**预期输出：**
```
共 3 个会话:

ID:     sess_abc123
标题：   项目重构
模型：   gpt-4
代理：   default
---
ID:     sess_def456
标题：   Bug 修复
模型：   claude-3-5-sonnet
---
ID:     sess_ghi789
模型：   gpt-3.5-turbo
代理：   code-reviewer
---
```

#### 示例 2：JSON 输出模式

```bash
oho session list --json
```

**预期输出：**
```json
[
  {
    "id": "sess_abc123",
    "title": "项目重构",
    "model": "gpt-4",
    "agent": "default",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T14:22:00Z"
  },
  {
    "id": "sess_def456",
    "title": "Bug 修复",
    "model": "claude-3-5-sonnet",
    "agent": "",
    "createdAt": "2024-01-14T09:15:00Z",
    "updatedAt": "2024-01-14T11:45:00Z"
  },
  {
    "id": "sess_ghi789",
    "title": "",
    "model": "gpt-3.5-turbo",
    "agent": "code-reviewer",
    "createdAt": "2024-01-13T16:00:00Z",
    "updatedAt": "2024-01-13T16:30:00Z"
  }
]
```

#### 示例 3：只显示正在运行的会话

```bash
oho session list --running
```

**预期输出：**
```
共 1 个会话:

ID:     sess_abc123
标题：   项目重构
模型：   gpt-4
代理：   default
---
```

#### 示例 4：空会话列表

```bash
oho session list
```

**预期输出：**
```
没有会话
```

#### 示例 5：JSON 格式空列表

```bash
oho session list --json
```

**预期输出：**
```json
[]
```

#### 示例 6：带运行状态的 JSON 输出

```bash
oho session list --json --running
```

**预期输出：**
```json
[
  {
    "id": "sess_abc123",
    "title": "项目重构",
    "model": "gpt-4",
    "agent": "default",
    "status": {
      "isReady": true,
      "isWorking": true,
      "status": "thinking"
    }
  }
]
```

---

## 快速参考表

### 命令速查

| 命令 | 功能 | 常用参数 |
|------|------|----------|
| `oho auth set <provider>` | 设置认证凭据 | `--credentials key=value` |
| `oho config get` | 获取配置 | `--json` |
| `oho config set` | 更新配置 | `--theme`, `--model`, `--temperature` |
| `oho config providers` | 列出提供商 | `--json` |
| `oho session list` | 列出会话 | `--json`, `--running` |

### 常见工作流

#### 工作流 1：初始化新环境

```bash
# 1. 设置认证凭据
oho auth set openai --credentials apiKey=sk-xxx

# 2. 查看可用提供商
oho config providers

# 3. 配置默认模型
oho config set --model gpt-4 --theme dark

# 4. 验证配置
oho config get
```

#### 工作流 2：查看和管理会话

```bash
# 1. 查看所有会话
oho session list

# 2. 查看正在运行的会话
oho session list --running

# 3. 以 JSON 格式导出会话列表
oho session list --json > sessions.json
```

#### 工作流 3：调整配置

```bash
# 1. 查看当前配置
oho config get

# 2. 更新温度和最大 token
oho config set --temperature 0.9 --max-tokens 8192

# 3. 设置自动批准的工具
oho config set --auto-approve read,ls,cat,grep

# 4. 验证更新
oho config get --json
```

---

## 错误处理示例

### 常见错误及解决方案

| 错误信息 | 原因 | 解决方案 |
|----------|------|----------|
| `Error: required flag(s) "credentials" not set` | 缺少 `--credentials` 参数 | 添加 `--credentials key=value` |
| `Error: 请提供至少一个凭据` | 凭据格式不正确 | 确保使用 `key=value` 格式 |
| `Error: 请提供至少一个要更新的配置项` | `config set` 未指定任何标志 | 添加如 `--theme`, `--model` 等标志 |
| `connection refused` | 服务器未启动 | 启动 OpenCode 服务器 |
| `401 Unauthorized` | 认证失败 | 检查 `auth set` 的凭据是否正确 |

---

## 环境配置

### 环境变量

在运行命令前，可设置以下环境变量：

```bash
export OPENCODE_SERVER_HOST=127.0.0.1
export OPENCODE_SERVER_PORT=4096
export OPENCODE_SERVER_PASSWORD=your-password
```

### 配置文件

配置文件位置：`~/.config/oho/config.json`

```json
{
  "host": "127.0.0.1",
  "port": 4096,
  "username": "opencode",
  "password": "",
  "json": false
}
```

---

## 附录：完整测试脚本

```bash
#!/bin/bash
# 测试 oho CLI 模块 1 的完整脚本

set -e

echo "=== oho CLI 模块 1 测试 ==="

# 1. 测试 auth set
echo -e "\n1. 测试 auth set..."
oho auth set openai --credentials apiKey=sk-test123
echo "✓ auth set 测试通过"

# 2. 测试 config get
echo -e "\n2. 测试 config get..."
oho config get
oho config get --json
echo "✓ config get 测试通过"

# 3. 测试 config set
echo -e "\n3. 测试 config set..."
oho config set --theme dark --model gpt-4
echo "✓ config set 测试通过"

# 4. 测试 config providers
echo -e "\n4. 测试 config providers..."
oho config providers
oho config providers --json
echo "✓ config providers 测试通过"

# 5. 测试 session list
echo -e "\n5. 测试 session list..."
oho session list
oho session list --json
oho session list --running
echo "✓ session list 测试通过"

echo -e "\n=== 所有测试完成 ==="
```

---

**文档版本**: 1.0  
**最后更新**: 2024-01-15  
**适用版本**: oho CLI v0.1.0+
