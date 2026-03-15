# oho CLI 操作指南 - 模块 3: 检查 Session

> **适用版本**: oho CLI v1.1+  
> **最后更新**: 2026-03-04  
> **作者**: nanobot 🐈  
> **前置模块**: [模块 2: 验证](./02-validation.md)

---

## 📋 目录

1. [查询会话列表](#1-查询会话列表)
2. [获取会话详情](#2-获取会话详情)
3. [查看会话状态](#3-查看会话状态)
4. [子会话管理](#4-子会话管理)
5. [会话过滤与搜索](#5-会话过滤与搜索)
6. [会话导出与分享](#6-会话导出与分享)

---

## 1. 查询会话列表

### 1.1 基本用法

```bash
# 列出所有会话
oho session list

# JSON 格式输出
oho session list --json
```

**预期输出**:
```
会话列表:
  ses_352a39c7bffe7RQv3VaA7Kypgs — shiny-squid (nanobot)
  ses_35261d371ffePsrv1yCkitEceB — quick-lagoon (hidden-falcon)
  ses_352536943ffet6eENMjHdaXW0Z — crisp-nebula (hidden-sailor)
  ...
```

**JSON 输出**:
```json
[
  {
    "id": "ses_352a39c7bffe7RQv3VaA7Kypgs",
    "slug": "shiny-squid",
    "title": "nanobot 版本检查",
    "project": "nanobot",
    "created": "2026-03-02T07:03:59Z",
    "updated": "2026-03-02T07:10:15Z",
    "messageCount": 12
  }
]
```

---

### 1.2 列表选项

```bash
# 限制显示数量
oho session list --limit 10

# 按时间排序
oho session list --sort created     # 按创建时间
oho session list --sort updated     # 按更新时间
oho session list --sort name        # 按名称

# 反向排序
oho session list --sort updated --reverse
```

**可用选项**:
| 选项 | 说明 | 默认值 |
|------|------|--------|
| `--limit` | 限制结果数量 | 全部 |
| `--sort` | 排序字段 | `updated` |
| `--reverse` | 反向排序 | `false` |
| `--json` | JSON 格式输出 | `false` |

---

### 1.3 按项目过滤

```bash
# 指定项目目录
oho session list --project /mnt/d/fe/nanobot

# 使用项目 ID
oho session list --project-id 086d65ace5213c9435dd217c4b5f3869c990e714
```

**输出示例**:
```
nanobot 项目会话:
  ses_352a39c7bffe7RQv3VaA7Kypgs — shiny-squid (07:03)
  ses_352a37f6dffeJqHEfobsUfCK4s — jolly-engine (07:04)
```

---

## 2. 获取会话详情

### 2.1 基本用法

```bash
# 获取会话详情
oho session get ses_352a39c7bffe7RQv3VaA7Kypgs

# 使用 -s 参数
oho session get -s ses_352a39c7bffe7RQv3VaA7Kypgs

# JSON 格式
oho session get ses_352a39c7bffe7RQv3VaA7Kypgs --json
```

**预期输出**:
```
会话详情:
  ID: ses_352a39c7bffe7RQv3VaA7Kypgs
  名称：shiny-squid
  标题：nanobot 版本检查
  项目：nanobot
  创建：2026-03-02 07:03:59
  更新：2026-03-02 07:10:15
  消息数：12
  状态：completed
  模型：alibaba-cn/qwen3.5-plus
```

---

### 2.2 详情字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | string | 会话唯一标识 |
| `slug` | string | 人类可读的名称 |
| `title` | string | 会话标题/描述 |
| `project` | string | 所属项目 |
| `directory` | string | 工作目录 |
| `created` | datetime | 创建时间 |
| `updated` | datetime | 最后更新时间 |
| `messageCount` | int | 消息数量 |
| `status` | string | 会话状态 |
| `model` | string | 使用的模型 |
| `parentID` | string | 父会话 ID (子会话) |

---

### 2.3 查看会话消息

```bash
# 列出会话中的消息
oho message list -s ses_352a39c7bffe7RQv3VaA7Kypgs

# JSON 格式
oho message list -s ses_352a39c7bffe7RQv3VaA7Kypgs --json

# 限制消息数量
oho message list -s ses_352a39c7bffe7RQv3VaA7Kypgs --limit 5

# 查看最新消息
oho message list -s ses_352a39c7bffe7RQv3VaA7Kypgs --tail 3
```

**消息输出**:
```
消息列表 (ses_352a39c7bffe7RQv3VaA7Kypgs):
  [1] 用户 (07:03:59) — "检查 nanobot 最新版本"
  [2] 助手 (07:04:15) — "正在检查 GitHub releases..."
  [3] 助手 (07:05:30) — "找到最新版本 v1.2.15"
  ...
```

---

## 3. 查看会话状态

### 3.1 基本用法

```bash
# 获取所有会话状态
oho session status

# JSON 格式
oho session status --json
```

**预期输出**:
```
会话状态:
  ✅ ses_352a39c7bffe7RQv3VaA7Kypgs — shiny-squid (completed)
  🔄 ses_35261d371ffePsrv1yCkitEceB — quick-lagoon (running)
  ⏸️  ses_352536943ffet6eENMjHdaXW0Z — crisp-nebula (paused)
  ❌ ses_35252fd63fferFlPqYpg0eb3Sv — clever-tiger (error)
```

---

### 3.2 状态类型

| 状态 | 标识 | 说明 |
|------|------|------|
| `running` | 🔄 | 正在执行任务 |
| `completed` | ✅ | 任务已完成 |
| `paused` | ⏸️ | 已暂停 |
| `error` | ❌ | 发生错误 |
| `aborted` | 🛑 | 用户中止 |
| `idle` | ⏳ | 等待输入 |

---

### 3.3 按状态过滤

```bash
# 查看运行中的会话
oho session status --filter running

# 查看错误的会话
oho session status --filter error

# 查看未完成的会话
oho session status --filter running,error,pending
```

**使用场景**:
- 🔄 `running` — 查看当前活跃任务
- ❌ `error` — 诊断失败任务
- ⏳ `idle` — 等待响应的会话
- ✅ `completed` — 查看已完成工作

---

### 3.4 实时监控

```bash
# 持续监控会话状态 (每秒刷新)
watch -n 1 "oho session status"

# 仅监控特定会话
watch -n 2 "oho session get ses_xxxxx --json | jq '.status'"
```

---

## 4. 子会话管理

### 4.1 查看子会话

```bash
# 获取子会话列表
oho session children ses_352a39c7bffe7RQv3VaA7Kypgs

# JSON 格式
oho session children ses_352a39c7bffe7RQv3VaA7Kypgs --json
```

**预期输出**:
```
子会话 (ses_352a39c7bffe7RQv3VaA7Kypgs):
  ses_352a37f6dffeJqHEfobsUfCK4s — jolly-engine (Fetch release notes)
  ses_352a37f5cffefLIwB1jqELAvJg — stellar-star (Git workflows 分析)
```

---

### 4.2 子会话层级

```bash
# 查看完整层级树
oho session tree ses_352a39c7bffe7RQv3VaA7Kypgs

# 限制深度
oho session tree ses_352a39c7bffe7RQv3VaA7Kypgs --depth 2
```

**层级示例**:
```
ses_352a39c7bffe7RQv3VaA7Kypgs (主会话)
├─ ses_352a37f6dffeJqHEfobsUfCK4s (子代理：GitHub release)
│  └─ ses_352a37f6dffeJqHEfobsUfCK4s-sub1 (嵌套子代理)
└─ ses_352a37f5cffefLIwB1jqELAvJg (子代理：Git workflows)
```

---

### 4.3 子会话详情

```bash
# 获取子会话详情
oho session get ses_352a37f6dffeJqHEfobsUfCK4s

# 查看子会话消息
oho message list -s ses_352a37f6dffeJqHEfobsUfCK4s

# 查看子会话与父会话的关系
oho session get ses_352a37f6dffeJqHEfobsUfCK4s --json | jq '.parentID'
```

---

## 5. 会话过滤与搜索

### 5.1 按时间过滤

```bash
# 今天创建的会话
oho session list --since today

# 最近 24 小时
oho session list --since 24h

# 特定日期之后
oho session list --since 2026-03-01

# 特定日期范围
oho session list --since 2026-03-01 --until 2026-03-02
```

**时间格式**:
- `today` — 今天 00:00
- `yesterday` — 昨天 00:00
- `24h` — 最近 24 小时
- `7d` — 最近 7 天
- `2026-03-01` — 特定日期
- `2026-03-01T10:00:00` — 精确时间

---

### 5.2 按关键词搜索

```bash
# 搜索标题包含关键词的会话
oho session list --search "nanobot"

# 搜索项目名
oho session list --search "hidden-sailor"

# 组合搜索
oho session list --search "tilemap" --since 24h
```

**搜索范围**:
- ✅ 会话标题
- ✅ 项目名称
- ✅ 会话 slug
- ❌ 消息内容 (使用 `oho message search`)

---

### 5.3 高级过滤

```bash
# 组合条件
oho session list --project nanobot --since today --status completed

# 排除条件
oho session list --exclude-status error,aborted

# 复杂查询 (JSON 输出后处理)
oho session list --json | jq '.[] | select(.messageCount > 10)'
```

---

### 5.4 多字段精确过滤（新增）

**oho CLI v1.1+** 支持按会话的每个字段进行精确过滤或模糊查询：

```bash
# 按会话 ID 过滤（支持模糊查询）
oho session list --id ses_abc123
oho session list --id "ses_34db"

# 按标题过滤（支持模糊查询，不区分大小写）
oho session list --title "babylon3D"
oho session list --title "水体测试"

# 按创建时间过滤（时间戳，精确匹配）
oho session list --created 1773537883643

# 按更新时间过滤（时间戳，精确匹配）
oho session list --updated 1773538142930

# 按项目 ID 过滤（支持模糊查询）
oho session list --project-id "1f01524d641bdfc5f4e43134de956b66c0b1332b"
oho session list --project-id "proj1"

# 按目录过滤（支持模糊查询）
oho session list --directory "rl_mockgame"
oho session list --directory "/mnt/d/code"
```

**过滤参数说明**:

| 参数 | 类型 | 匹配方式 | 示例 |
|------|------|---------|------|
| `--id` | string | 模糊匹配（不区分大小写） | `--id "ses_abc"` |
| `--title` | string | 模糊匹配（不区分大小写） | `--title "测试"` |
| `--created` | int64 | 精确匹配（Unix 时间戳毫秒） | `--created 1773537883643` |
| `--updated` | int64 | 精确匹配（Unix 时间戳毫秒） | `--updated 1773538142930` |
| `--project-id` | string | 模糊匹配（不区分大小写） | `--project-id "proj1"` |
| `--directory` | string | 模糊匹配（不区分大小写） | `--directory "babylon"` |

---

### 5.5 组合过滤示例

```bash
# 组合多个字段过滤
oho session list --title "babylon3D" --directory "babylon3DWorld"

# 组合状态 + 字段过滤
oho session list --status running --title "测试"

# 组合排序 + 过滤 + 分页
oho session list --title "项目" --sort updated --order desc --limit 10

# JSON 输出 + 字段过滤
oho session list --project-id "proj1" --json | jq '.[] | {id, title, directory}'
```

**实际输出示例**:

```bash
$ oho session list --title "babylon3D" --directory "babylon3DWorld"
共 2 个会话:

ID:     ses_34dbffe0dffe8SfdMTbL53MWFP
标题：   babylon3D 水体测试与地图编辑器
模型：   
---
ID:     ses_35720ca6cffetpjG9PEV9bIcKZ
标题：   探索 babylon3DWorld 湖泊渲染代码 (@explore subagent)
模型：   
---
```

---

### 5.6 过滤与排序、分页组合

```bash
# 过滤后排序
oho session list --title "测试" --sort created --order asc

# 过滤后分页
oho session list --directory "project1" --limit 5 --offset 10

# 完整组合
oho session list \
  --title "babylon3D" \
  --directory "babylon3DWorld" \
  --sort updated \
  --order desc \
  --limit 10 \
  --offset 0
```

**排序参数**:

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--sort` | `updated` | 排序字段：`created` 或 `updated` |
| `--order` | `desc` | 排序顺序：`asc`（升序）或 `desc`（降序） |
| `--limit` | 无限制 | 返回结果数量上限 |
| `--offset` | 0 | 分页偏移量 |

---

## 6. 会话导出与分享

### 6.1 导出会话

```bash
# 导出会话为 Markdown
oho session export ses_352a39c7bffe7RQv3VaA7Kypgs

# 导出为 JSON
oho session export ses_352a39c7bffe7RQv3VaA7Kypgs --format json

# 导出到文件
oho session export ses_352a39c7bffe7RQv3VaA7Kypgs --output session.md
```

**导出内容**:
- ✅ 会话元数据
- ✅ 所有消息
- ✅ 子会话引用
- ✅ 时间戳

---

### 6.2 分享会话

```bash
# 生成分享链接
oho session share ses_352a39c7bffe7RQv3VaA7Kypgs

# 设置过期时间
oho session share ses_352a39c7bffe7RQv3VaA7Kypgs --expires 24h

# 设置密码
oho session share ses_352a39c7bffe7RQv3VaA7Kypgs --password "share123"
```

**分享选项**:
| 选项 | 说明 | 默认值 |
|------|------|--------|
| `--expires` | 过期时间 | 永久 |
| `--password` | 访问密码 | 无 |
| `--public` | 公开访问 | `false` |

---

### 6.3 取消分享

```bash
# 取消分享
oho session unshare ses_352a39c7bffe7RQv3VaA7Kypgs

# 查看所有分享的会话
oho session list --shared
```

---

## 🔧 实用技巧

### 技巧 1: 快速定位活跃会话

```bash
# 一键获取最近更新的会话
oho session list --sort updated --reverse --limit 1 --json | jq '.[0].id'

# 保存为变量
LATEST_SESSION=$(oho session list --sort updated --reverse --limit 1 --json | jq -r '.[0].id')
oho message list -s "$LATEST_SESSION"
```

---

### 技巧 2: 会话统计

```bash
# 统计总会话数
oho session list --json | jq 'length'

# 统计各状态会话数
oho session status --json | jq 'group_by(.status) | map({status: .[0].status, count: length})'

# 统计今日会话数
oho session list --since today --json | jq 'length'
```

---

### 技巧 3: 批量操作

```bash
# 删除所有错误会话
oho session list --filter error --json | jq -r '.[].id' | xargs -I {} oho session delete {}

# 导出所有完成的会话
oho session list --filter completed --json | jq -r '.[].id' | while read id; do
    oho session export "$id" --output "export-$id.md"
done
```

---

### 技巧 4: 会话对比

```bash
# 对比两个会话的差异
oho session diff ses_xxxxx ses_yyyyy

# 查看会话变更历史
oho session history ses_xxxxx
```

---

## 📝 检查清单

在开始会话检查前，请确认：

- [ ] 已认证并连接到服务器
- [ ] 知道目标会话 ID 或项目名称
- [ ] 了解需要的信息类型 (列表/详情/状态)
- [ ] 选择合适的输出格式 (文本/JSON)

---

## 🔗 相关文档

- [模块 1: 客户端初始化](./01-client-initialization.md) - 连接服务器
- [模块 2: 验证](./02-validation.md) - 身份验证
- [模块 6: 发送消息](./06-send-message.md) - 消息操作
- [模块 8: 查询状态](./08-query-status.md) - 任务状态监控

---

## 🆘 常见问题

### Q1: 如何找到特定项目的会话？

**A**: 使用项目过滤:
```bash
oho session list --project /mnt/d/fe/nanobot
```

---

### Q2: 如何查看子代理的执行结果？

**A**: 先获取子会话 ID，再查看消息:
```bash
# 获取子会话
oho session children ses_main

# 查看子会话消息
oho message list -s ses_child
```

---

### Q3: 会话列表为空怎么办？

**A**: 检查以下几点:
```bash
# 1. 确认服务器连接
oho config get

# 2. 检查项目目录
oho session list --project /correct/path

# 3. 查看时间范围
oho session list --since 30d
```

---

### Q4: 如何清理旧会话？

**A**: 批量删除:
```bash
# 删除 7 天前的会话
oho session list --until 7d --json | jq -r '.[].id' | xargs oho session delete
```

---

*文档生成时间：2026-03-02 23:55 CST*  
*最后验证：2026-03-04 02:42 CST*

---

## 🔬 实际验证输出 (2026-03-04 02:42)

### 验证 1: oho session list

```bash
$ oho session list
共 48 个会话:

ID:     ses_34c5b5c54ffehnE3JBss6tWts1
标题：   New session - 2026-03-03T12:20:37.425Z
模型：   
---
ID:     ses_34dbffe0dffe8SfdMTbL53MWFP
标题：   babylon3D 水体测试与地图编辑器
模型：   
---
ID:     ses_35725f2eeffecp7ZPxdGfCnPkO
标题：   New session - 2026-03-01T10:03:08.433Z
模型：   
---
ID:     ses_35720ca6cffetpjG9PEV9bIcKZ
标题：   探索 babylon3DWorld 湖泊渲染代码 (@explore subagent)
模型：   
---
ID:     ses_357212bd2ffeF92okRPBuhHVlp
标题：   搜索 wujimanager 中湖泊相关代码 (@explore subagent)
模型：   
---
... (共 48 个会话)
```

**说明**:
- 默认按更新时间排序
- 显示会话 ID、标题、模型
- 子代理会话标记为 `(@explore subagent)`

---

### 验证 2: oho session get

```bash
$ oho session get ses_34dbffe0dffe8SfdMTbL53MWFP
共 1 个会话:

ID:     ses_34dbffe0dffe8SfdMTbL53MWFP
标题：   babylon3D 水体测试与地图编辑器
模型：   
---
```

**说明**:
- 支持完整会话 ID 或 slug
- 返回会话基本信息
- 使用 `--json` 获取详细信息

---

### 验证 3: oho message list

```bash
$ oho message list -s ses_34dbffe0dffe8SfdMTbL53MWFP --limit 5

[assistant] msg_cb243748e0014GeHtPFdynAR6l
  └─ 部分类型：step-start
  └─ 部分类型：reasoning
  └─ 部分类型：text
  └─ 部分类型：step-finish
---

[user] msg_cb244d0f5001YWkvi0wcW5x6bk
  └─ 部分类型：text
---

[assistant] msg_cb244d5d1001B1s4qzqOZRoV75
  └─ 部分类型：step-start
  └─ 部分类型：reasoning
  └─ 部分类型：text
  └─ 部分类型：step-finish
---

[user] msg_cb3ceb103001k3YZEa2Yu1HgZ7
  └─ 部分类型：text
---

[assistant] msg_cb3ceb2430010sOIaswVqOOsSW
  └─ 部分类型：step-start
  └─ 部分类型：reasoning
  └─ 部分类型：text
  └─ 部分类型：step-finish
---
```

**说明**:
- 显示消息 ID、角色、部分类型
- `step-start/step-finish` 表示 AI 思考步骤
- `reasoning` 包含推理过程
- `text` 是实际响应内容

---

### 验证 4: oho session create

```bash
$ oho session create
会话创建成功:
  ID: ses_34afd94f6ffe4IWeoe4rpzHidB
  标题：New session - 2026-03-03T18:42:40.522Z
  模型：
```

**说明**:
- 自动分配唯一会话 ID
- 默认标题为 "New session - 时间戳"
- 模型字段为空，等待首次消息时确定

---

### 验证 5: oho message add (带 --no-reply)

```bash
$ oho message add -s ses_34dbffe0dffe8SfdMTbL53MWFP "测试模块 3 文档验证" --no-reply
DEBUG: 发送请求:
{
  "noReply": true,
  "parts": [
    {
      "type": "text",
      "text": "测试模块 3 文档验证"
    }
  ]
}
消息已发送:
  ID: msg_cb5026b4900140a4HD3hLneKgF
  角色：user

[text]
测试模块 3 文档验证
```

**说明**:
- `--no-reply` 不等待 AI 响应
- 返回消息 ID 和角色
- 适合批量发送或异步任务

---

### 验证 6: oho --help (完整命令列表)

```bash
$ oho --help
Available Commands:
  agent       代理命令
  auth        认证管理
  command     命令管理
  completion  Generate the autocompletion script for the specified shell
  config      配置管理命令
  file        文件管理命令
  find        查找命令
  formatter   格式化器状态
  global      全局命令
  help        Help about any command
  lsp         LSP 服务器状态
  mcp         MCP 服务器管理
  mcpserver   启动 MCP 服务器
  message     消息管理命令
  project     项目管理命令
  provider    提供商管理命令
  session     会话管理命令
  tool        工具命令
  tui         TUI 控制命令
```

**常用命令分类**:
- **会话管理**: `session list`, `session get`, `session create`
- **消息操作**: `message add`, `message list`
- **项目管理**: `project list`, `project current`, `project path`
- **配置管理**: `config get`, `config providers`
- **认证管理**: `auth set`

---

### 验证 7: 会话过滤与搜索

```bash
# 查找包含关键词的会话
$ oho session list | grep "babylon3D"
ID:     ses_34dbffe0dffe8SfdMTbL53MWFP
标题：   babylon3D 水体测试与地图编辑器

# 查找子代理会话
$ oho session list | grep "@explore"
ID:     ses_35720ca6cffetpjG9PEV9bIcKZ
标题：   探索 babylon3DWorld 湖泊渲染代码 (@explore subagent)
ID:     ses_357212bd2ffeF92okRPBuhHVlp
标题：   搜索 wujimanager 中湖泊相关代码 (@explore subagent)
```

---

### 验证 8: 会话状态监控

```bash
# 查看会话数量统计
$ oho session list | grep -c "^ID:"
48

# 查看最新消息
$ oho message list -s ses_34dbffe0dffe8SfdMTbL53MWFP --limit 1

# 查看会话创建时间
$ oho session list | grep "2026-03-04"
```

---

### 验证 9: 错误处理示例

```bash
# 错误：无效的会话 ID
$ oho session get invalid-id
Error: API 错误 [400]: {"error": "invalid session id"}

# 错误：会话不存在
$ oho message list -s ses_nonexistent
Error: 会话不存在

# 错误：认证失败
$ unset OPENCODE_SERVER_PASSWORD
$ oho session list
Error: API 错误 [401]: {"error": "unauthorized"}
```

**解决方案**:
```bash
# 1. 使用正确的会话 ID
oho session list  # 查看有效 ID

# 2. 重新认证
export OPENCODE_SERVER_PASSWORD="your_password"

# 3. 检查服务器状态
ps aux | grep opencode
```
