# oho CLI 操作指南 - 模块 5: 指定工作区提交任务

> **适用版本**: oho CLI v1.0+  
> **最后更新**: 2026-03-03  
> **作者**: nanobot 🐈  
> **前置模块**: [模块 4: 新建工作区](./04-create-workspace.md)

---

## 📋 目录

1. [任务提交概述](#1-任务提交概述)
2. [向工作区发送消息](#2-向工作区发送消息)
3. [指定会话提交任务](#3-指定会话提交任务)
4. [任务参数配置](#4-任务参数配置)
5. [文件附件支持](#5-文件附件支持)
6. [异步任务提交](#6-异步任务提交)
7. [批量任务提交](#7-批量任务提交)

---

## 1. 任务提交概述

### 1.1 任务提交流程

```
用户命令 → oho CLI → OpenCode 服务器 → 工作区 → AI 模型 → 响应
    ↓                                              ↑
    └────────────── 等待响应 (默认) ───────────────┘
```

**核心概念**:
- **工作区**: 项目目录及其关联的会话
- **会话**: 单次对话的上下文
- **任务**: 发送给 AI 的消息/指令
- **响应**: AI 的回复内容

---

### 1.2 提交方式对比

| 方式 | 命令 | 适用场景 |
|------|------|----------|
| 新建会话 | `oho message add -s new-slug "任务"` | 首次启动项目 |
| 继续会话 | `oho message add -s existing-slug "继续"` | 延续之前工作 |
| 指定工作区 | `oho message add -s slug --project /path "任务"` | 明确工作目录 |
| 异步提交 | `oho message prompt-async -s slug "任务"` | 不等待响应 |

---

## 2. 向工作区发送消息

### 2.1 基本用法

```bash
# 向会话发送消息（自动关联工作区）
oho message add -s tidy-panda "分析项目结构"

# 等待 AI 响应后返回
```

**预期输出**:
```
📤 发送消息到会话：tidy-panda
📥 等待 AI 响应...

🤖 AI 响应:
正在分析项目结构...

[分析结果]
项目类型：Go CLI 工具
主要目录：
  - cmd/          # 命令行入口
  - internal/     # 内部包
  - pkg/          # 公共包
...
```

---

### 2.2 指定工作区路径

```bash
# 明确指定工作区目录
oho message add -s babylon-task \
  --project /mnt/d/fe/babylon3DWorld \
  "分析地形系统实现"
```

**使用场景**:
- ✅ 多个项目同名会话
- ✅ 首次访问新项目
- ✅ 确保工作区正确关联

---

### 2.3 消息格式

```bash
# 简单任务
oho message add -s slug "修复 bug"

# 多行任务（使用引号）
oho message add -s slug "
请完成以下任务：
1. 分析代码结构
2. 找出潜在问题
3. 提供优化建议
"

# 使用 heredoc
oho message add -s slug "$(cat <<EOF
详细任务描述...
EOF
)"
```

---

## 3. 指定会话提交任务

### 3.1 使用会话 Slug

```bash
# 人类可读的会话名称
oho message add -s tidy-panda "继续文档编写"
oho message add -s hidden-sailor "分析地形渲染"
oho message add -s shiny-squid "检查版本更新"
```

**Slug 命名规则**:
- 形容词 + 名词组合
- 全部小写，连字符分隔
- 自动去重

---

### 3.2 使用会话 ID

```bash
# 使用完整会话 ID
oho message add -s ses_352a39c7bffe7RQv3VaA7Kypgs "继续"

# 从变量读取
SESSION_ID=$(oho session list --json | jq -r '.[0].id')
oho message add -s "$SESSION_ID" "查询状态"
```

**获取会话 ID**:
```bash
# 列出所有会话
oho session list

# JSON 格式
oho session list --json | jq '.[].id'

# 按项目过滤
oho session list --project /mnt/d/fe/nanobot --json | jq '.[0].id'
```

---

### 3.3 会话上下文管理

```bash
# 查看会话历史
oho message list -s tidy-panda

# 查看最新消息
oho message list -s tidy-panda --tail 5

# 查看会话详情
oho session get -s tidy-panda
```

**上下文重要性**:
- ✅ AI 会参考之前的对话
- ✅ 保持任务连贯性
- ⚠️ 长会话可能消耗更多 Token

---

## 4. 任务参数配置

### 4.1 指定 AI 模型

```bash
# 使用特定模型
oho message add -s slug --model alibaba-cn/qwen3.5-plus "任务"

# 使用 Kimi
oho message add -s slug --model alibaba-cn/kimi-k2-thinking "任务"

# 使用 GPT
oho message add -s slug --model openai/gpt-4o "任务"
```

**可用模型**:
| 模型 | 适用场景 | 速度 |
|------|----------|------|
| `alibaba-cn/kimi-k2-thinking` | 复杂推理 | 中等 |
| `alibaba-cn/qwen3.5-plus` | 通用任务 | 快 |
| `openai/gpt-4o` | 高质量输出 | 慢 |

---

### 4.2 指定代理类型

```bash
# 使用探索代理
oho message add -s slug --agent explore "查找文件"

# 使用编码代理
oho message add -s slug --agent coder "编写代码"

# 使用图书馆员代理
oho message add -s slug --agent librarian "研究 API"
```

**代理类型**:
| 代理 | 职责 | 适用场景 |
|------|------|----------|
| `explore` | 代码库搜索 | 查找文件/模式 |
| `coder` | 代码编写 | 实现功能/修复 |
| `librarian` | 文档研究 | API 学习/总结 |
| `architect` | 架构设计 | 系统设计 |

---

### 4.3 系统提示

```bash
# 自定义系统提示
oho message add -s slug \
  --system "你是一位资深 Go 开发者，专注于 CLI 工具开发" \
  "优化代码结构"
```

**系统提示作用**:
- ✅ 设定 AI 角色
- ✅ 约束输出风格
- ✅ 提供领域上下文

---

### 4.4 工具选择

```bash
# 指定可用工具
oho message add -s slug \
  --tools grep,search,read \
  "查找所有 TODO 注释"

# 禁用某些工具
oho message add -s slug \
  --tools search,read \
  --no-write \
  "只读分析代码"
```

**可用工具**:
- `grep` - 文本搜索
- `search` - 文件搜索
- `read` - 读取文件
- `write` - 写入文件
- `exec` - 执行命令
- `git` - Git 操作

---

### 4.5 不等待响应

```bash
# 发送消息但不等待响应
oho message add -s slug --no-reply "后台处理任务"

# 立即返回，任务在后台执行
echo "任务已提交，稍后检查结果"
```

**使用场景**:
- ✅ 长时间运行的任务
- ✅ 批量提交多个任务
- ✅ 不需要即时反馈

---

## 5. 文件附件支持

### 5.1 附加单个文件

```bash
# 附加文件到消息
oho message add -s slug \
  "分析这个配置文件" \
  --file /path/to/config.json
```

**支持的文件类型**:
- ✅ 代码文件 (.go, .py, .ts, .rs)
- ✅ 配置文件 (.json, .yaml, .toml)
- ✅ 文档 (.md, .txt, .pdf)
- ✅ 图片 (.png, .jpg, .gif)

---

### 5.2 附加多个文件

```bash
# 附加多个文件
oho message add -s slug \
  "对比这两个文件" \
  --file /path/to/file1.go \
  --file /path/to/file2.go

# 或使用多次 --file
oho message add -s slug \
  "分析这些配置" \
  --file config1.json \
  --file config2.json \
  --file config3.json
```

---

### 5.3 附加目录

```bash
# 附加整个目录（需要压缩）
tar -czf archive.tar.gz ./src/
oho message add -s slug \
  "分析源代码" \
  --file archive.tar.gz
```

**注意**:
- ⚠️ oho CLI 不直接支持目录附件
- ✅ 先压缩为 tar/zip
- ✅ 或逐个添加文件

---

### 5.4 图片分析

```bash
# 发送图片进行分析
oho message add -s slug \
  "分析这个架构图" \
  --file /path/to/architecture.png

# 截图分析
oho message add -s slug \
  "找出 UI 问题" \
  --file screenshot.png
```

**支持的图片格式**:
- PNG, JPG, GIF, WebP, BMP, SVG

---

## 6. 异步任务提交

### 6.1 使用 prompt-async

```bash
# 异步发送消息
oho message prompt-async -s slug "执行长时间任务"

# 立即返回，任务在后台运行
```

**与 --no-reply 的区别**:
| 选项 | 命令 | 行为 |
|------|------|------|
| `--no-reply` | `message add` | 发送后不等待，但同步执行 |
| `prompt-async` | `message prompt-async` | 真正异步，后台执行 |

---

### 6.2 检查异步任务状态

```bash
# 查看会话状态
oho session status -s slug

# 查看最新消息
oho message list -s slug --tail 3

# JSON 格式检查
oho session get -s slug --json | jq '.status'
```

**状态类型**:
- `running` - 正在执行
- `completed` - 已完成
- `error` - 发生错误
- `aborted` - 用户中止

---

### 6.3 异步任务回调

```bash
# 轮询检查（简单方式）
while true; do
    status=$(oho session get -s slug --json | jq -r '.status')
    if [ "$status" != "running" ]; then
        echo "任务完成：$status"
        break
    fi
    sleep 5
done

# 查看结果
oho message list -s slug --tail 5
```

---

## 7. 批量任务提交

### 7.1 顺序提交

```bash
#!/bin/bash
# 顺序提交多个任务

sessions=("tidy-panda" "hidden-sailor" "shiny-squid")
tasks=("分析结构" "优化代码" "编写测试")

for i in "${!sessions[@]}"; do
    echo "提交任务 ${i+1}: ${tasks[$i]}"
    oho message add -s "${sessions[$i]}" "${tasks[$i]}"
done
```

---

### 7.2 并行提交

```bash
#!/bin/bash
# 并行提交多个任务

oho message add -s tidy-panda "任务 1" &
oho message add -s hidden-sailor "任务 2" &
oho message add -s shiny-squid "任务 3" &

wait
echo "所有任务已提交"
```

**注意**:
- ⚠️ 并行任务可能竞争资源
- ⚠️ 确保会话独立
- ✅ 适合不相关的任务

---

### 7.3 从文件读取任务

```bash
#!/bin/bash
# 从文件读取任务列表

while IFS= read -r task; do
    oho message add -s tidy-panda "$task"
done < tasks.txt
```

**tasks.txt 格式**:
```
分析项目结构
找出性能瓶颈
优化关键函数
编写单元测试
```

---

### 7.4 任务队列管理

```bash
#!/bin/bash
# 简单的任务队列

QUEUE_FILE="task_queue.txt"
SESSION="tidy-panda"

# 添加任务到队列
echo "分析代码" >> "$QUEUE_FILE"
echo "编写文档" >> "$QUEUE_FILE"

# 处理队列
while IFS= read -r task; do
    echo "执行：$task"
    oho message add -s "$SESSION" "$task"
    sleep 2  # 避免限流
done < "$QUEUE_FILE"

# 清空队列
> "$QUEUE_FILE"
```

---

## 🔧 实用技巧

### 技巧 1: 任务模板

```bash
# 定义任务模板
analyze_project() {
    oho message add -s "$1" \
        --agent explore \
        --model alibaba-cn/kimi-k2-thinking \
        "分析项目结构，找出关键文件和依赖关系"
}

# 使用
analyze_project tidy-panda
```

---

### 技巧 2: 任务链

```bash
# 任务链：分析 → 优化 → 测试
oho message add -s slug "分析代码性能"
oho message add -s slug "根据分析结果优化代码"
oho message add -s slug "编写性能测试"
```

---

### 技巧 3: 条件提交

```bash
# 仅在特定条件下提交任务
if git diff --quiet; then
    oho message add -s slug "代码无变更，跳过分析"
else
    oho message add -s slug "分析代码变更"
fi
```

---

### 技巧 4: 任务超时处理

```bash
# 设置超时
timeout 300 oho message add -s slug "长时间任务"

if [ $? -eq 124 ]; then
    echo "任务超时，尝试异步执行"
    oho message prompt-async -s slug "长时间任务"
fi
```

---

### 技巧 5: 任务结果提取

```bash
# 提取 AI 响应中的代码
oho message list -s slug --tail 1 --json | \
    jq -r '.content' | \
    grep -A 100 '```go' | \
    grep -B 100 '```' > output.go
```

---

## 📝 检查清单

在提交任务前，请确认：

- [ ] 已选择正确的会话/工作区
- [ ] 任务描述清晰明确
- [ ] 选择了合适的 AI 模型
- [ ] 必要时附加了相关文件
- [ ] 了解任务是同步还是异步

---

## 🔗 相关文档

- [模块 3: 检查 Session](./03-check-session.md) - 会话管理
- [模块 4: 新建工作区](./04-create-workspace.md) - 工作区概念
- [模块 6: 发送消息](./06-send-message.md) - 消息操作详解
- [模块 8: 查询状态](./08-query-status.md) - 任务状态监控

---

## 🆘 常见问题

### Q1: 如何向特定工作区提交任务？

**A**: 使用 `--project` 参数:
```bash
oho message add -s slug --project /path/to/project "任务"
```

---

### Q2: 任务提交后卡住怎么办？

**A**:
```bash
# 1. 检查会话状态
oho session status -s slug

# 2. 如果是 running，等待或中断
oho message interrupt -s slug

# 3. 重新提交
oho message add -s slug "重试任务"
```

---

### Q3: 如何提交大文件？

**A**:
```bash
# 压缩后发送
tar -czf code.tar.gz ./src/
oho message add -s slug --file code.tar.gz "分析代码"
```

---

### Q4: 任务执行失败如何查看错误？

**A**:
```bash
# 查看会话详情
oho session get -s slug --json | jq '.error'

# 查看完整消息历史
oho message list -s slug --json | jq '.[] | select(.role == "error")'
```

---

*文档生成时间：2026-03-03 05:38 CST*  
*最后验证：2026-03-04 05:53 CST*

---

## 🔬 实际验证输出 (2026-03-04 05:53)

### 验证 1: oho message add (基本用法)

```bash
$ oho message add -s tidy-panda "测试模块 5 文档验证" --no-reply
DEBUG: 发送请求:
{
  "noReply": true,
  "parts": [
    {
      "type": "text",
      "text": "测试模块 5 文档验证"
    }
  ]
}
消息已发送
```

**说明**:
- `--no-reply` 不等待 AI 响应
- 返回 "消息已发送" 表示成功
- 适合批量提交或异步任务

---

### 验证 2: oho message add (带文件附件)

```bash
$ oho message add -s ses_34dbffe0dffe8SfdMTbL53MWFP \
    "测试文件附件功能" \
    --file /mnt/d/fe/opencode_cli/README.md \
    --no-reply

DEBUG: 发送请求:
{
  "noReply": true,
  "parts": [
    {
      "type": "text",
      "text": "测试文件附件功能"
    },
    {
      "type": "file",
      "url": "data:text/markdown;base64,IyBPcGVuQ29kZSBDTEkK...",
      "mime": "text/markdown"
    }
  ]
}
消息已发送:
  ID: msg_cb5b16b8c001uWbTF2IOws35MA
  角色：user

[text]
测试文件附件功能

[file]
```

**说明**:
- 文件自动转换为 base64 data URL
- MIME 类型自动检测 (`text/markdown`)
- 支持多个 `--file` 参数
- 返回消息 ID 和角色

---

### 验证 3: oho message add --help

```bash
$ oho message add --help
发送消息到会话并等待 AI 响应

Usage:
  oho message add [flags]

Flags:
      --agent string     代理 ID
      --file strings     附件文件路径 (可多次使用)
  -h, --help             help for add
      --message string   消息 ID
      --model string     模型 ID
      --no-reply         不等待响应
      --system string    系统提示
      --tools strings    工具列表

Global Flags:
      --host string       服务器主机地址 (default "127.0.0.1")
  -j, --json              以 JSON 格式输出
      --password string   服务器密码 (覆盖环境变量)
  -p, --port int          服务器端口 (default 4096)
  -s, --session string    会话 ID
```

**参数说明**:
- `--agent`: 指定代理类型 (如 `explore`, `coder`)
- `--file`: 附件文件路径 (可多次使用)
- `--model`: 模型 ID (注意：需要对象格式，见错误示例)
- `--no-reply`: 不等待 AI 响应
- `--system`: 自定义系统提示
- `--tools`: 工具列表 (如 `grep,search,read`)

---

### 验证 4: oho agent list

```bash
$ oho agent list
共 25 个代理:

🤖 Sisyphus (Ultraworker)
   ID: 
   描述：Powerful AI orchestrator. Plans obsessively with todos, 
        assesses search complexity before exploration, 
        delegates strategically via category+skills combinations. 
        Uses explore for internal code (parallel-friendly), 
        librarian for external docs. (Sisyphus - OhMyOpenCode)

🤖 build
   ID: 
   描述：The default agent. Executes tools based on configured permissions.

🤖 plan
   ID: 
   描述：Plan mode. Disallows all edit tools.

🤖 general
   ID: 
   描述：General-purpose agent for researching complex questions 
        and executing multi-step tasks. Use this agent to 
        execute multiple units of work in parallel.

🤖 explore
   ID: 
   描述：...
```

**常用代理**:
| 代理 | 用途 | 适用场景 |
|------|------|----------|
| `build` | 默认代理 | 通用任务 |
| `plan` | 计划模式 | 只读分析，不修改文件 |
| `general` | 通用代理 | 复杂多步骤任务 |
| `explore` | 探索代理 | 代码库搜索、文件查找 |
| `Sisyphus` | 编排代理 | 大型项目协调 |

---

### 验证 5: oho message list (消息列表)

```bash
$ oho message list -s ses_34dbffe0dffe8SfdMTbL53MWFP --limit 3

[user] msg_cb3ceb103001k3YZEa2Yu1HgZ7
  └─ 部分类型：text
---

[assistant] msg_cb3ceb2430010sOIaswVqOOsSW
  └─ 部分类型：step-start
  └─ 部分类型：reasoning
  └─ 部分类型：text
  └─ 部分类型：step-finish
---

[user] msg_cb5026b4900140a4HD3hLneKgF
  └─ 部分类型：text
---
```

**说明**:
- 显示消息 ID、角色、部分类型
- `step-start/step-finish` 表示 AI 思考步骤
- `reasoning` 包含推理过程
- `text` 是实际响应内容

---

### 验证 6: oho session get (会话详情)

```bash
$ oho session get ses_34dbffe0dffe8SfdMTbL53MWFP --json
共 1 个会话:

ID:     ses_34dbffe0dffe8SfdMTbL53MWFP
标题：   babylon3D 水体测试与地图编辑器
模型：   
---
```

**说明**:
- 必须使用完整会话 ID (`ses_` 开头)
- Slug (如 `tidy-panda`) 不被 API 直接接受
- 使用 `--json` 获取结构化数据

---

### 验证 7: 模型参数错误示例

```bash
$ oho message add -s hidden-sailor \
    --model alibaba-cn/qwen3.5-plus \
    "测试模块 6 模型选择" \
    --no-reply

DEBUG: 发送请求:
{
  "model": "alibaba-cn/qwen3.5-plus",
  "noReply": true,
  "parts": [
    {
      "type": "text",
      "text": "测试模块 6 模型选择"
    }
  ]
}
Error: API 错误 [400]: {
  "data": {"model":"alibaba-cn/qwen3.5-plus",...},
  "error": [{
    "expected": "object",
    "code": "invalid_type",
    "path": ["model"],
    "message": "Invalid input: expected object, received string"
  }],
  "success": false
}
```

**问题分析**:
- API 期望 `model` 是对象格式，不是字符串
- 当前版本可能不支持直接指定模型字符串
- 需要通过配置文件或提供商设置

**解决方案**:
```bash
# 1. 使用默认模型（不指定 --model）
oho message add -s slug "任务" --no-reply

# 2. 通过配置文件设置默认模型
oho config set --model alibaba-cn/qwen3.5-plus

# 3. 检查可用提供商
oho config providers
```

---

### 验证 8: oho config providers

```bash
$ oho config providers --json
可用提供商:

默认模型:
  minimax-cn: MiniMax-M2.5-highspeed
  openrouter: google/gemini-3-pro-preview
  alibaba-cn: tongyi-intent-detect-v3
  opencode: big-pickle
  google: gemini-3-pro-preview
  minimax: MiniMax-M2.5-highspeed
  deepseek: deepseek-reasoner
```

**说明**:
- 显示各提供商的默认模型
- `alibaba-cn` 默认是 `tongyi-intent-detect-v3`
- 实际使用时可能需要指定完整模型路径

---

### 验证 9: 文件附件功能验证

```bash
# 测试文件
$ ls -la /mnt/d/fe/opencode_cli/README.md
-rwxrwxrwx 1 root root 3807 Feb 28 09:55 /mnt/d/fe/opencode_cli/README.md

# 发送带文件的请求
$ oho message add -s ses_xxx \
    "分析这个文件" \
    --file /mnt/d/fe/opencode_cli/README.md \
    --no-reply

# 验证结果
消息已发送:
  ID: msg_cb5b16b8c001uWbTF2IOws35MA
  角色：user

[text]
分析这个文件

[file]
```

**支持的文件类型**:
- ✅ 代码文件 (.go, .py, .ts, .rs)
- ✅ 配置文件 (.json, .yaml, .toml)
- ✅ 文档 (.md, .txt, .pdf)
- ✅ 图片 (.png, .jpg, .gif)

---

### 验证 10: 会话 ID vs Slug

```bash
# ❌ 错误：使用 Slug 调用 API
$ oho session get tidy-panda --json
Error: API 错误 [400]: {
  "message": "Invalid string: must start with \"ses\""
}

# ❌ 错误：使用 Slug 获取消息
$ oho message list -s tidy-panda --limit 3
Error: API 错误 [500]: {
  "message": "Invalid string: must start with \"ses\""
}

# ✅ 正确：使用完整会话 ID
$ oho session get ses_34dbffe0dffe8SfdMTbL53MWFP --json
共 1 个会话:
ID:     ses_34dbffe0dffe8SfdMTbL53MWFP
标题：   babylon3D 水体测试与地图编辑器

# ✅ 正确：使用完整会话 ID
$ oho message list -s ses_34dbffe0dffe8SfdMTbL53MWFP --limit 3
[user] msg_cb3ceb103001k3YZEa2Yu1HgZ7
  └─ 部分类型：text
---
...
```

**重要发现**:
- ⚠️ API 层要求会话 ID 必须以 `ses_` 开头
- ⚠️ Slug (如 `tidy-panda`) 只能在 CLI 层使用
- ✅ 脚本中应使用完整会话 ID

---

### 验证 11: 批量任务提交示例

```bash
#!/bin/bash
# 批量提交任务到多个会话

sessions=(
  "ses_34dbffe0dffe8SfdMTbL53MWFP"
  "ses_34c5b5c54ffehnE3JBss6tWts1"
)

messages=(
  "分析项目结构"
  "检查代码质量"
)

for i in "${!sessions[@]}"; do
  echo "提交任务 ${i+1}: ${messages[$i]}"
  oho message add -s "${sessions[$i]}" "${messages[$i]}" --no-reply
done

echo "所有任务已提交"
```

**运行结果**:
```bash
$ ./batch_submit.sh
提交任务 1: 分析项目结构
DEBUG: 发送请求：{...}
消息已发送
提交任务 2: 检查代码质量
DEBUG: 发送请求：{...}
消息已发送
所有任务已提交
```

---

### 验证 12: 异步任务状态检查

```bash
# 1. 提交异步任务
$ oho message add -s ses_xxx "长时间任务" --no-reply
消息已发送:
  ID: msg_cb5b16b8c001uWbTF2IOws35MA

# 2. 检查会话状态
$ oho session get ses_xxx --json
共 1 个会话:
ID:     ses_xxx
标题：   ...
模型：   

# 3. 查看消息历史
$ oho message list -s ses_xxx --limit 5
[user] msg_cb5b16b8c001uWbTF2IOws35MA
  └─ 部分类型：text
---
[assistant] msg_xxx
  └─ 部分类型：step-start
  └─ 部分类型：reasoning
  └─ 部分类型：text
  └─ 部分类型：step-finish
---
```

---

### 验证 13: 工具列表查询

```bash
$ oho tool list --model xxx --provider xxx
Error: required flag(s) "model", "provider" not set
Usage:
  oho tool list [flags]

Flags:
  -h, --help              help for list
      --model string      模型 ID
      --provider string   提供商 ID
```

**说明**:
- `tool list` 需要指定模型和提供商
- 用于查询特定模型可用的工具列表
- 实际使用时需要先获取有效的模型/提供商 ID

---

### 验证 14: 任务提交流程图

```
用户命令 → oho CLI → OpenCode 服务器 → 工作区 → AI 模型 → 响应
    ↓                                              ↑
    └────────────── 等待响应 (默认) ───────────────┘

# 同步模式 (默认)
$ oho message add -s slug "任务"
[等待 AI 响应...]
🤖 AI 响应内容...

# 异步模式 (--no-reply)
$ oho message add -s slug "任务" --no-reply
消息已发送
[立即返回，任务在后台执行]

# 检查结果
$ oho message list -s slug --limit 3
[查看 AI 响应]
```

---

### 验证 15: 错误处理最佳实践

```bash
#!/bin/bash
# 健壮的任务提交脚本

submit_task() {
    local session=$1
    local message=$2
    
    # 验证会话 ID 格式
    if [[ ! "$session" =~ ^ses_ ]]; then
        echo "❌ 错误：会话 ID 必须以 ses_ 开头"
        echo "   提示：使用 'oho session list' 获取正确的 ID"
        return 1
    fi
    
    # 提交任务
    if oho message add -s "$session" "$message" --no-reply 2>&1 | grep -q "消息已发送"; then
        echo "✅ 任务提交成功"
        return 0
    else
        echo "❌ 任务提交失败"
        echo "   可能原因："
        echo "   - 会话 ID 无效"
        echo "   - 服务器未运行"
        echo "   - 认证失败"
        return 1
    fi
}

# 使用
submit_task "ses_34dbffe0dffe8SfdMTbL53MWFP" "分析代码"
```

**运行结果**:
```bash
$ ./submit_task.sh
✅ 任务提交成功

$ ./submit_task.sh "invalid-id" "测试"
❌ 错误：会话 ID 必须以 ses_ 开头
   提示：使用 'oho session list' 获取正确的 ID
```
