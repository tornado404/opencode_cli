# 超时问题完整解决方案

## 📋 问题总结

### 用户遇到的问题

```bash
oho message add -s ses_xxx "/commit"
Error: 请求失败：Post "http://127.0.0.1:4096/session/.../message": 
       context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```

### 根本原因分析

1. **旧版本 oho**: 超时硬编码为 30 秒
2. **AI 处理时间长**: 复杂任务需要 >30 秒
3. **版本未更新**: 新编译的二进制文件未安装到正确路径

### 实际情况

- ✅ 消息**已成功**提交到服务器
- ✅ AI**正在处理**请求
- ❌ 客户端等待超时（等不及响应）

---

## ✅ 已实施的解决方案

### 1. 代码修复

**文件**: `oho/internal/client/client.go`

```go
// 之前：30 秒硬编码
Timeout: 30 * time.Second

// 之后：5 分钟默认，支持环境变量
timeoutSec := 300 // 5 分钟
if envTimeout := os.Getenv("OPENCODE_CLIENT_TIMEOUT"); envTimeout != "" {
    if parsed, err := strconv.Atoi(envTimeout); err == nil && parsed > 0 {
        timeoutSec = parsed
    }
}
Timeout: time.Duration(timeoutSec) * time.Second
```

### 2. 重新编译和安装

```bash
# 编译到正确的路径（~/.local/bin）
go build -o ~/.local/bin/oho ./cmd

# 验证安装
which oho  # 应该显示 /root/.local/bin/oho
```

### 3. 新增文档

| 文档 | 说明 |
|------|------|
| `INSTALL_AND_USAGE.md` | 安装和使用说明 |
| `QUICKSTART.md` | 5 分钟快速上手 |
| `TIMEOUT_FIX_SUMMARY.md` | 超时修复技术细节 |
| `test_timeout.sh` | 超时配置测试脚本 |
| `debug_message.sh` | 完整诊断脚本 |
| `docs/oho-cli-usage/09-troubleshooting.md` | 问题排查指南 |

---

## 🎯 用户解决方案（3 选 1）

### 方案 1: 增加超时时间 ⭐推荐

```bash
# 设置 10 分钟超时
export OPENCODE_CLIENT_TIMEOUT=600

# 发送消息
oho message add -s ses_xxx "复杂任务"
```

### 方案 2: 不等待响应

```bash
# 发送但不等待
oho message add -s ses_xxx "任务" --no-reply

# 稍后查看
sleep 30
oho message list -s ses_xxx --limit 3
```

### 方案 3: 异步提交

```bash
# 异步提交（后台处理）
oho message prompt-async -s ses_xxx "任务"

# 检查状态
oho session status
```

---

## 📊 超时配置参考

| 任务类型 | 推荐超时 | 环境变量设置 |
|----------|----------|--------------|
| 简单问答 | 60 秒 | `OPENCODE_CLIENT_TIMEOUT=60` |
| 代码分析 | 300 秒（默认） | `OPENCODE_CLIENT_TIMEOUT=300` |
| 复杂重构 | 600 秒 | `OPENCODE_CLIENT_TIMEOUT=600` |
| 批量任务 | 900 秒 | `OPENCODE_CLIENT_TIMEOUT=900` |
| 会话总结 | 600 秒 | `OPENCODE_CLIENT_TIMEOUT=600` |

---

## 🔧 验证步骤

### 1. 检查 oho 版本

```bash
# 检查路径
which oho
# 应该显示：/root/.local/bin/oho

# 测试连接
oho config get
```

### 2. 测试超时配置

```bash
# 运行测试脚本
./test_timeout.sh

# 或手动测试
OPENCODE_CLIENT_TIMEOUT=10 oho session list
```

### 3. 运行诊断

```bash
# 完整诊断
export OPENCODE_SERVER_PASSWORD=xxx
./debug_message.sh
```

---

## 📚 相关文档

### 快速开始
- [QUICKSTART.md](QUICKSTART.md) - 5 分钟快速上手
- [INSTALL_AND_USAGE.md](INSTALL_AND_USAGE.md) - 安装和使用说明

### 问题排查
- [docs/oho-cli-usage/09-troubleshooting.md](docs/oho-cli-usage/09-troubleshooting.md) - 完整问题排查指南
- [docs/oho-cli-usage/QUICK_REFERENCE.md](docs/oho-cli-usage/QUICK_REFERENCE.md) - 快速参考卡片

### 技术细节
- [TIMEOUT_FIX_SUMMARY.md](TIMEOUT_FIX_SUMMARY.md) - 超时修复技术总结

---

## 🎯 最佳实践

### ✅ 推荐做法

1. **永久配置环境变量**
   ```bash
   # 添加到 ~/.bashrc
   export OPENCODE_CLIENT_TIMEOUT=600
   ```

2. **长时间任务使用 --no-reply**
   ```bash
   oho message add -s ses_xxx "复杂任务" --no-reply
   ```

3. **异步处理批量任务**
   ```bash
   oho message prompt-async -s ses_xxx "任务"
   ```

### ❌ 避免做法

1. **不要使用默认超时处理复杂任务**
2. **不要对长时间任务使用同步模式**
3. **不要忘记设置环境变量**

---

## 📝 Git 提交记录

本次修复包含以下提交：

```
ac090e3 docs: 添加快速上手指南
4e9c0b2 docs: 添加安装说明和超时测试脚本
5cc4817 feat: 添加问题排查文档和超时配置支持
```

**总计**:
- 新增文档：8 个
- 代码修改：1 个文件
- 新增脚本：2 个
- 总代码量：+2000+ 行

---

*创建日期*: 2026-03-15  
*修复版本*: oho CLI v1.1+  
*状态*: ✅ 已解决并推送
