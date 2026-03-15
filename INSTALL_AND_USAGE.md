# oho CLI 安装和使用说明

## 🚨 重要：超时问题解决方案

如果你遇到 `context deadline exceeded` 错误，请按以下步骤操作：

### 问题原因

- ✅ 消息已成功提交到服务器
- ✅ AI 正在处理你的请求
- ❌ 但客户端等待超时（默认 5 分钟）

### 解决方案（3 选 1）

#### 方案 1: 增加超时时间（推荐）

```bash
# 设置 10 分钟超时
export OPENCODE_CLIENT_TIMEOUT=600

# 然后发送消息
oho message add -s ses_xxx "复杂任务"
```

#### 方案 2: 不等待响应

```bash
# 发送消息但不等待 AI 响应
oho message add -s ses_xxx "任务" --no-reply

# 稍后查看结果
sleep 30
oho message list -s ses_xxx --limit 3
```

#### 方案 3: 异步提交

```bash
# 异步提交任务（后台处理）
oho message prompt-async -s ses_xxx "任务"

# 检查状态
oho session status
```

---

## 📦 安装说明

### 从源码编译

```bash
cd /path/to/opencode_cli/oho

# 编译并安装到 ~/.local/bin（推荐）
go build -o ~/.local/bin/oho ./cmd

# 或安装到 /usr/local/bin
sudo go build -o /usr/local/bin/oho ./cmd
```

### 验证安装

```bash
# 检查 oho 路径
which oho

# 应该看到：/root/.local/bin/oho 或 /usr/local/bin/oho

# 测试连接
oho config get
```

---

## 🔧 环境配置

### 基本配置

```bash
# 服务器连接
export OPENCODE_SERVER_HOST=127.0.0.1
export OPENCODE_SERVER_PORT=4096
export OPENCODE_SERVER_PASSWORD=your-password

# 超时配置（可选）
export OPENCODE_CLIENT_TIMEOUT=600  # 10 分钟
```

### 永久配置（推荐）

将以下内容添加到 `~/.bashrc` 或 `~/.zshrc`：

```bash
# OpenCode oho CLI 配置
export OPENCODE_SERVER_HOST=127.0.0.1
export OPENCODE_SERVER_PORT=4096
export OPENCODE_SERVER_PASSWORD=your-password
export OPENCODE_CLIENT_TIMEOUT=600
```

然后执行：
```bash
source ~/.bashrc  # 或 source ~/.zshrc
```

---

## 📋 常用命令

### 会话管理

```bash
# 创建会话
oho session create

# 列出会话
oho session list

# 查看会话详情
oho session get -s ses_xxx

# 中止会话
oho session abort -s ses_xxx
```

### 消息管理

```bash
# 发送消息（等待响应）
oho message add -s ses_xxx "你好"

# 发送消息（不等待响应）
oho message add -s ses_xxx "分析项目" --no-reply

# 异步发送消息
oho message prompt-async -s ses_xxx "任务"

# 查看消息历史
oho message list -s ses_xxx

# 带文件附件
oho message add -s ses_xxx "分析这个文件" --file /path/to/file.go
```

### 配置管理

```bash
# 查看配置
oho config get

# 查看可用模型
oho config providers
```

---

## 🐛 故障排查

### 1. 超时错误

```
Error: context deadline exceeded
```

**解决**:
```bash
export OPENCODE_CLIENT_TIMEOUT=600
oho message add -s ses_xxx "任务" --no-reply
```

### 2. 认证失败

```
Error: API 错误 [401]: 认证失败
```

**解决**:
```bash
export OPENCODE_SERVER_PASSWORD=your-password
```

### 3. 连接被拒绝

```
Error: connection refused
```

**解决**:
```bash
# 启动 OpenCode 服务器
opencode serve --port 4096
```

### 4. 运行诊断脚本

```bash
# 完整诊断
export OPENCODE_SERVER_PASSWORD=xxx
./debug_message.sh

# 快速测试
./test_timeout.sh
```

---

## 📊 超时配置参考

| 任务类型 | 推荐超时 | 命令示例 |
|----------|----------|----------|
| 简单问答 | 60 秒 | `export OPENCODE_CLIENT_TIMEOUT=60` |
| 代码分析 | 300 秒（默认） | `export OPENCODE_CLIENT_TIMEOUT=300` |
| 复杂重构 | 600 秒 | `export OPENCODE_CLIENT_TIMEOUT=600` |
| 批量任务 | 900 秒 | `export OPENCODE_CLIENT_TIMEOUT=900` |

---

## 🔗 相关文档

- [问题排查指南](docs/oho-cli-usage/09-troubleshooting.md)
- [快速参考卡片](docs/oho-cli-usage/QUICK_REFERENCE.md)
- [超时修复总结](TIMEOUT_FIX_SUMMARY.md)
- [完整教程](docs/oho-cli-usage/README.md)

---

*最后更新：2026-03-15*
