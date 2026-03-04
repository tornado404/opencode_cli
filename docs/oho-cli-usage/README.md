# oho CLI 操作指南

> **完整教程**: 从入门到精通  
> **适用版本**: oho CLI v1.0+  
> **最后更新**: 2026-03-03  
> **状态**: 🟡 进行中（持续完善中）

---

## 📚 教程目录

本教程共 8 个模块，涵盖 oho CLI 的核心功能：

| 模块 | 主题 | 文件 | 状态 |
|------|------|------|------|
| 1 | 客户端初始化 | [01-client-initialization.md](./01-client-initialization.md) | ✅ 已完成 |
| 2 | 验证 | [02-validation.md](./02-validation.md) | ✅ 已完成 |
| 3 | 检查 Session | [03-check-session.md](./03-check-session.md) | ✅ 已完成 |
| 4 | 新建工作区 | [04-create-workspace.md](./04-create-workspace.md) | ✅ 已完成 |
| 5 | 指定工作区提交任务 | [05-submit-task.md](./05-submit-task.md) | ✅ 已完成 |
| 6 | 指定 session_id 和模型发消息 | [06-send-message.md](./06-send-message.md) | ✅ 已完成 |
| 7 | 中断任务 | [07-interrupt-task.md](./07-interrupt-task.md) | ✅ 已完成 |
| 8 | 查询任务执行情况及状态 | [08-query-status.md](./08-query-status.md) | ✅ 已完成 |

**总体进度**: 8/8 模块创建完成 (100%) 🎉  
**完善进度**: 🟡 持续完善中（补充实际输出示例、命令验证）

---

## 🔧 持续完善计划

**当前任务**: 补充完整命令示例、预期输出、文档结构优化

**完善内容**:
- [🔧] 完整命令示例 (持续补充)
- [🔧] 预期输出示例 (持续补充)
- [🔧] 文档结构优化 (持续改进)
- [🔧] 实际运行验证 (持续进行)

**完善计划文档**: [IMPROVEMENT_PLAN.md](./IMPROVEMENT_PLAN.md)

**下次执行**: 补充模块 1-2 的实际输出验证

---

## 🚀 快速开始

### 前置要求

- ✅ 已安装 oho CLI (`/usr/local/bin/oho`)
- ✅ OpenCode Server 正在运行
- ✅ 已配置服务器密码

### 5 分钟上手

```bash
# 1. 设置认证
export OPENCODE_SERVER_PASSWORD="your_password"

# 2. 验证连接
oho config get

# 3. 创建会话
SESSION=$(oho session create --json | jq -r '.id')

# 4. 发送消息
oho message add -s "$SESSION" "Hello, World!"

# 5. 查看响应
oho message list -s "$SESSION"
```

---

## 📖 模块详情

### 模块 1: 客户端初始化

**内容**:
- 认证配置 (`oho auth set`)
- 服务器连接参数 (`--host`, `--port`, `--password`)
- 配置管理 (`oho config get/set`)
- 连接验证方法
- 常见问题排查

**适用场景**:
- 首次使用 oho CLI
- 切换服务器环境
- 配置问题诊断

**阅读时间**: 10 分钟

👉 [开始学习](./01-client-initialization.md)

---

### 模块 2: 验证

**计划内容**:
- 身份验证流程
- 权限检查
- Token 管理
- 错误处理

**状态**: ⏳ 待编写

---

### 模块 3: 检查 Session

**计划内容**:
- 查询会话列表 (`oho session list`)
- 获取会话详情 (`oho session get`)
- 查看会话状态 (`oho session status`)
- 子会话管理 (`oho session children`)

**状态**: ⏳ 待编写

---

### 模块 4: 新建工作区

**计划内容**:
- 创建工作区
- 工作区配置
- 工作区切换
- 工作区清理

**状态**: ⏳ 待编写

---

### 模块 5: 指定工作区提交任务

**计划内容**:
- 任务提交方式
- 任务参数配置
- 任务优先级
- 批量任务处理

**状态**: ⏳ 待编写

---

### 模块 6: 指定 session_id 和模型发消息

**计划内容**:
- 发送消息 (`oho message add`)
- 指定模型参数
- 文件上传 (`--file`)
- 异步消息 (`oho message prompt-async`)

**状态**: ⏳ 待编写

---

### 模块 7: 中断任务

**计划内容**:
- 中止任务 (`oho session abort`)
- 任务恢复
- 超时处理
- 优雅中断

**状态**: ⏳ 待编写

---

### 模块 8: 查询任务执行情况及状态

**计划内容**:
- 任务状态查询
- 执行日志查看
- 结果导出
- 性能分析

**状态**: ⏳ 待编写

---

## 🛠️ 常用命令速查

### 会话管理
```bash
oho session create              # 创建会话
oho session list                # 列出会话
oho session get -s <id>         # 获取详情
oho session delete -s <id>      # 删除会话
```

### 消息管理
```bash
oho message add -s <id> "msg"   # 发送消息
oho message list -s <id>        # 查看消息
oho message command -s <id> /help  # 执行命令
```

### 配置管理
```bash
oho config get                  # 查看配置
oho config set --model <name>   # 设置模型
oho config providers            # 列出模型
```

### 认证管理
```bash
oho auth set                    # 设置密码
```

---

## 📝 更新日志

### 2026-03-02
- ✅ 创建教程目录结构
- ✅ 完成模块 1: 客户端初始化
- ⏳ 计划编写模块 2-8

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 改进文档！

**反馈渠道**:
- GitHub Issues
- Telegram 群组
- 邮件列表

---

*文档由 nanobot 🐈 生成和维护*
