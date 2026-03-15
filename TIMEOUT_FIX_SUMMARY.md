# 超时问题修复总结

## 问题描述

用户在使用 `oho message add` 提交消息时遇到超时错误：

```
Error: 请求失败：Post "http://127.0.0.1:4096/session/.../message": 
       context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```

## 问题分析

### 根本原因

客户端 HTTP 超时硬编码为 **30 秒**，对于需要长时间思考的 AI 任务来说太短了。

### 错误解读

这个错误**不是 Bug**，而是：
- ✅ 请求已成功发送到服务器
- ✅ 服务器正在处理（AI 正在思考）
- ❌ 客户端等不及响应就断开了（30 秒超时）

### 实际情况

消息实际上已经提交成功，AI 可能正在处理，但 CLI 客户端在收到响应前就超时断开了。

---

## 解决方案

### 1. 代码修改

**文件**: `oho/internal/client/client.go`

**修改内容**:
```go
// 之前：硬编码 30 秒
httpClient: &http.Client{
    Timeout: 30 * time.Second,
}

// 之后：默认 5 分钟，支持环境变量配置
timeoutSec := 300 // 5 分钟
if envTimeout := os.Getenv("OPENCODE_CLIENT_TIMEOUT"); envTimeout != "" {
    if parsed, err := strconv.Atoi(envTimeout); err == nil && parsed > 0 {
        timeoutSec = parsed
    }
}
httpClient: &http.Client{
    Timeout: time.Duration(timeoutSec) * time.Second,
}
```

### 2. 新增环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `OPENCODE_CLIENT_TIMEOUT` | 300 秒（5 分钟） | HTTP 请求超时时间 |

### 3. 使用方法

```bash
# 使用默认超时（5 分钟）
oho message add -s ses_xxx "任务"

# 设置更长超时（10 分钟）
export OPENCODE_CLIENT_TIMEOUT=600
oho message add -s ses_xxx "复杂任务"

# 不等待响应
oho message add -s ses_xxx "任务" --no-reply

# 异步提交
oho message prompt-async -s ses_xxx "任务"
```

---

## 文档更新

### 新增/更新的文件

1. **docs/oho-cli-usage/09-troubleshooting.md**
   - 添加 "4.0 超时错误" 章节
   - 详细说明超时原因和解决方案

2. **docs/oho-cli-usage/QUICK_REFERENCE.md**
   - 添加超时配置说明
   - 添加超时错误速查项

3. **oho/internal/client/client.go**
   - 修改超时配置逻辑
   - 添加环境变量支持

---

## 验证步骤

```bash
# 1. 重新编译
cd /mnt/d/fe/opencode_cli/oho
go build -o /usr/local/bin/oho ./cmd

# 2. 验证版本
oho --help

# 3. 测试超时配置
export OPENCODE_CLIENT_TIMEOUT=60
oho message add -s ses_xxx "测试" --no-reply

# 4. 检查超时是否生效
# 查看代码或运行诊断脚本
./debug_message.sh
```

---

## 最佳实践建议

### 超时配置建议

| 任务类型 | 推荐超时 | 说明 |
|----------|----------|------|
| 简单问答 | 60 秒 | 快速响应任务 |
| 代码分析 | 300 秒（默认） | 中等复杂度任务 |
| 复杂重构 | 600 秒 | 大型代码库分析 |
| 批量任务 | 900 秒 | 多文件处理 |

### 避免超时的方法

1. **使用 `--no-reply`**: 发送消息但不等待响应
2. **使用 `prompt-async`**: 异步提交任务
3. **增加超时时间**: 设置 `OPENCODE_CLIENT_TIMEOUT`
4. **简化任务**: 将大任务分解为小任务

---

## 技术细节

### HTTP 超时类型

| 超时类型 | 说明 | 当前配置 |
|----------|------|----------|
| Connection Timeout | 建立连接超时 | 由系统决定 |
| Response Timeout | 等待响应超时 | 300 秒（可配置） |
| Read Timeout | 读取数据超时 | 包含在 Response Timeout 中 |

### 超时的作用

- 防止客户端无限期等待
- 释放卡住的连接
- 避免资源浪费

### 为什么选择 5 分钟默认值

- 30 秒：太短，复杂任务会超时
- 5 分钟：平衡大多数任务需求
- 30 分钟：太长，可能掩盖服务器问题

---

*修复日期*: 2026-03-15  
*修复版本*: oho CLI v1.1+  
*相关文档*: docs/oho-cli-usage/09-troubleshooting.md
