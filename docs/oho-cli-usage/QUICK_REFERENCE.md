# oho CLI 快速参考卡片

> 一页纸速查表 - 打印或保存为书签

---

## 🔧 环境配置

```bash
# 基本配置
export OPENCODE_SERVER_HOST=127.0.0.1
export OPENCODE_SERVER_PORT=4096
export OPENCODE_SERVER_PASSWORD=your-password

# 超时配置（默认 5 分钟，复杂任务可延长）
export OPENCODE_CLIENT_TIMEOUT=600  # 10 分钟

# 验证配置
oho config get
oho config providers
```

---

## 📋 常用命令

### 会话管理
| 命令 | 说明 |
|------|------|
| `oho session list` | 列出所有会话 |
| `oho session create` | 创建新会话 |
| `oho session get -s <id>` | 会话详情 |
| `oho session abort -s <id>` | 中止会话 |
| `oho session delete <id>` | 删除会话 |

### 消息管理
| 命令 | 说明 |
|------|------|
| `oho message add -s <id> "msg"` | 发送消息 |
| `oho message add -s <id> "msg" --file <path>` | 带文件附件 |
| `oho message add -s <id> "msg" --no-reply` | 不等待响应 |
| `oho message list -s <id>` | 查看消息历史 |
| `oho message get <msgId> -s <id>` | 消息详情 |

### 配置管理
| 命令 | 说明 |
|------|------|
| `oho config get` | 查看配置 |
| `oho config providers` | 可用模型列表 |
| `oho agent list` | 可用代理列表 |

---

## 🐛 问题排查

### 快速诊断（5 步）

```bash
# 1. 服务器是否运行？
curl http://localhost:4096/global/health

# 2. 密码是否配置？
echo $OPENCODE_SERVER_PASSWORD

# 3. 认证是否成功？
oho config get

# 4. 会话是否可用？
oho session list

# 5. 模型是否可用？
oho config providers
```

### 常见问题速查

| 症状 | 可能原因 | 解决方案 |
|------|----------|----------|
| `connection refused` | 服务器未运行 | `opencode serve` |
| `401 Unauthorized` | 密码错误 | 检查 `OPENCODE_SERVER_PASSWORD` |
| `must start with "ses"` | 使用了 Slug | 使用完整会话 ID |
| `context deadline exceeded` | 请求超时 | `OPENCODE_CLIENT_TIMEOUT=600` |
| 消息无响应 | 会话忙碌 | `oho session status` |
| 消息无响应 | 模型未配置 | `oho config providers` |
| 文件不存在 | 路径错误 | 使用绝对路径 |

---

## 📊 诊断脚本

```bash
# 运行完整诊断
export OPENCODE_SERVER_PASSWORD=xxx
./debug_message.sh

# 快速检查
curl http://localhost:4096/global/health && \
  oho config get > /dev/null && \
  echo "✓ 一切正常" || echo "✗ 有问题"
```

---

## 📁 文件位置

| 文件 | 位置 |
|------|------|
| 配置文件 | `~/.config/oho/config.json` |
| 会话数据 | `~/.opencode/sessions/` |
| OpenCode 配置 | `~/.opencode/config.json` |

---

## 🔗 文档链接

- [完整教程](./docs/oho-cli-usage/README.md)
- [问题排查指南](./docs/oho-cli-usage/09-troubleshooting.md)
- [API 映射](./oho/API_MAPPING.md)

---

*最后更新：2026-03-15*
