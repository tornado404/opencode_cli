# oho CLI 快速上手指南

## ⚡ 5 分钟快速开始

### 步骤 1: 设置环境

```bash
# 设置服务器密码（替换为你的实际密码）
export OPENCODE_SERVER_PASSWORD=your-password

# 设置超时时间（推荐 10 分钟，避免长时间任务超时）
export OPENCODE_CLIENT_TIMEOUT=600
```

### 步骤 2: 验证连接

```bash
# 检查服务器连接
oho config get

# 查看可用模型
oho config providers
```

### 步骤 3: 创建会话

```bash
# 创建新会话
SESSION_ID=$(oho session create 2>&1 | grep "ID:" | awk '{print $2}')
echo "创建会话：$SESSION_ID"
```

### 步骤 4: 发送消息

```bash
# 方法 A: 等待响应（适合简单任务）
oho message add -s "$SESSION_ID" "你好，请帮我分析项目结构"

# 方法 B: 不等待响应（适合长时间任务）⭐推荐
oho message add -s "$SESSION_ID" "分析项目结构" --no-reply

# 方法 C: 异步提交（后台处理）
oho message prompt-async -s "$SESSION_ID" "分析项目结构"
```

### 步骤 5: 查看结果

```bash
# 查看消息历史
oho message list -s "$SESSION_ID" --limit 5

# 查看会话状态
oho session status
```

---

## 🎯 避免超时错误的最佳实践

### ✅ 推荐做法

```bash
# 1. 设置足够的超时时间
export OPENCODE_CLIENT_TIMEOUT=600  # 10 分钟

# 2. 使用 --no-reply 发送长时间任务
oho message add -s ses_xxx "复杂任务" --no-reply

# 3. 稍后查看结果
sleep 60
oho message list -s ses_xxx --limit 3

# 4. 或使用异步提交
oho message prompt-async -s ses_xxx "任务"
```

### ❌ 避免做法

```bash
# 不要：不设置超时时间（可能使用默认值）
oho message add -s ses_xxx "复杂分析任务"

# 不要：对长时间任务使用同步模式
oho message add -s ses_xxx "重构整个项目"  # 可能超时！
```

---

## 📋 常用命令速查

### 会话操作

```bash
oho session create                    # 创建会话
oho session list                      # 列出所有会话
oho session get -s <id>               # 查看会话详情
oho session abort -s <id>             # 中止会话
oho session delete <id>               # 删除会话
```

### 消息操作

```bash
oho message add -s <id> "内容"                   # 发送消息
oho message add -s <id> "内容" --no-reply        # 不等待响应
oho message prompt-async -s <id> "内容"          # 异步提交
oho message list -s <id>                         # 查看历史
oho message get <msg_id> -s <id>                 # 消息详情
```

### 文件附件

```bash
# 附加单个文件
oho message add -s <id> "分析这个文件" --file /path/to/file.go

# 附加多个文件
oho message add -s <id> "对比这些文件" \
  --file file1.go \
  --file file2.go
```

---

## 🔧 环境配置（永久）

将以下内容添加到 `~/.bashrc` 或 `~/.zshrc`：

```bash
# OpenCode oho CLI 配置
export OPENCODE_SERVER_HOST=127.0.0.1
export OPENCODE_SERVER_PORT=4096
export OPENCODE_SERVER_PASSWORD=your-password
export OPENCODE_CLIENT_TIMEOUT=600

# 可选：添加别名
alias oho-session='oho session list'
alias oho-config='oho config get'
```

然后执行：
```bash
source ~/.bashrc  # 或 source ~/.zshrc
```

---

## 🐛 遇到问题？

### 超时错误

```bash
# 症状：context deadline exceeded
# 解决：增加超时时间或使用 --no-reply
export OPENCODE_CLIENT_TIMEOUT=600
oho message add -s ses_xxx "任务" --no-reply
```

### 认证失败

```bash
# 症状：API 错误 [401]
# 解决：检查密码
export OPENCODE_SERVER_PASSWORD=correct-password
```

### 连接被拒绝

```bash
# 症状：connection refused
# 解决：启动服务器
opencode serve --port 4096
```

### 运行诊断

```bash
# 完整诊断脚本
export OPENCODE_SERVER_PASSWORD=xxx
./debug_message.sh

# 超时测试
./test_timeout.sh
```

---

## 📚 更多文档

- [完整教程](docs/oho-cli-usage/README.md)
- [问题排查](docs/oho-cli-usage/09-troubleshooting.md)
- [快速参考](docs/oho-cli-usage/QUICK_REFERENCE.md)
- [安装说明](INSTALL_AND_USAGE.md)

---

*最后更新：2026-03-15*
