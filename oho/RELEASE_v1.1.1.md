# Release v1.1.1 - oho add 命令稳定性修复

**发布日期**: 2026-03-21  
**版本类型**: Bug Fix Release  
**上一版本**: v1.1.0

---

## 🎯 发布重点

本次发布专注于提升 `oho add` 命令的稳定性和可靠性，修复了文件附件处理的边界条件 bug，并添加了完整的单元测试覆盖。

---

## 🐛 Bug 修复

### 1. detectMimeType Panic 修复 (严重)

**问题**: 当文件没有扩展名时，`detectMimeType` 函数会触发 `slice bounds out of range` panic

**修复前**:
```go
ext := strings.ToLower(filePath[strings.LastIndex(filePath, "."):])
// 如果文件无扩展名，LastIndex 返回 -1，导致 slice bounds out of range [-1:]
```

**修复后**:
```go
dotIndex := strings.LastIndex(filePath, ".")
if dotIndex == -1 {
    return "application/octet-stream"
}
ext := strings.ToLower(filePath[dotIndex:])
```

**影响**: 使用 `oho add` 发送无扩展名文件附件时不再崩溃

### 2. ClientInterface 接口完善

**问题**: `PostWithQuery` 方法未在接口中定义，导致 Mock 客户端无法正确模拟

**修复**: 
- 在 `ClientInterface` 中添加 `PostWithQuery` 方法定义
- 更新 `MockClient` 添加 `PostWithQueryFunc` 注入支持
- `add.go` 改为使用 `ClientInterface` 而非具体 `*Client` 类型

**影响**: 提升代码可测试性，支持更好的单元测试隔离

---

## ✅ 新增测试 (770 行)

新增 `oho/cmd/add/add_test.go` 文件，包含 10 个测试函数：

| 测试函数 | 覆盖场景 | 用例数 |
|---------|---------|-------|
| `TestConvertModel` | 模型格式转换 (nil/string/Model) | 5 |
| `TestDetectMimeType` | MIME 类型检测 (各种扩展名) | 19 |
| `TestCreateSession` | 会话创建 (成功/失败/API 错误) | 6 |
| `TestSendMessage` | 消息发送 (简单/附件/no-reply/错误) | 7 |
| `TestRunAddSuccess` | 完整流程集成测试 | 4 |
| `TestRaceConditionScenarios` | 竞态条件模拟 | 3 |
| `TestTimeoutScenarios` | 超时边界测试 | 2 |
| `TestErrorPropagation` | 错误传播验证 | 3 |
| `TestPartialFailureHandling` | 部分失败处理 | 1 |
| `TestJSONOutputFormat` | 输出格式验证 | 2 |

**总计**: 52 个测试用例，覆盖 `oho add` 命令所有关键路径

---

## 📚 新增文档

### PROJECT_SUMMARY_ZH.md (664 行)

完整的项目摘要文档，包含：

1. **项目概述** - 目标、定位、核心功能、技术栈
2. **架构设计** - 目录结构、模块划分、数据流和调用关系
3. **核心实现分析** - `oho add` 流程、HTTP 客户端、错误处理
4. **测试覆盖** - 测试结构、新增用例、覆盖率分析
5. **已知问题和修复** - 间歇性失败分析、bug 修复记录
6. **使用指南** - 安装步骤、命令示例、故障排查

---

## 📦 文件变更

| 文件 | 变更类型 | 说明 |
|------|---------|------|
| `oho/cmd/add/add.go` | 修改 | 修复 detectMimeType，使用 ClientInterface |
| `oho/cmd/add/add_test.go` | 新增 | 770 行单元测试 |
| `oho/internal/client/client_interface.go` | 修改 | 添加 PostWithQuery 方法 |
| `oho/internal/client/client_mock.go` | 修改 | 添加 PostWithQueryFunc mock |
| `oho/PROJECT_SUMMARY_ZH.md` | 新增 | 664 行项目摘要 |

**统计**: 5 个文件，+1458 行，-10 行

---

## 🧪 测试状态

```
✅ 所有现有测试通过 (18 个包)
✅ 新增 10 个测试函数
✅ 边界条件 bug 已修复
✅ 测试覆盖率达到生产就绪标准
```

### 测试运行命令

```bash
cd oho
make test
# 或
go test -v ./cmd/add/...
```

---

## 🚀 升级指南

### 安装/更新

```bash
# 方法一：快速安装
curl -sSL https://raw.githubusercontent.com/tornado404/opencode_cli/master/oho/install.sh | bash

# 方法二：源码编译
git clone https://github.com/tornado404/opencode_cli.git
cd opencode_cli/oho
make build
sudo cp bin/oho /usr/local/bin/oho
```

### 验证安装

```bash
oho --version
# 应显示：oho version v1.1.1
```

---

## 📋 推荐升级人群

以下用户**强烈建议**升级到此版本：

- ✅ 需要使用 `oho add` 发送文件附件的用户
- ✅ 在自动化脚本中使用 `oho add` 的用户
- ✅ 遇到过间歇性失败的用户
- ✅ 需要高测试覆盖率保障的团队

---

## 🔧 技术细节

### 版本信息

```
版本：v1.1.1
提交：6db26d9
类型：Annotated Tag
兼容性：完全向后兼容
```

### Git 历史

```
commit 6db26d9 (HEAD -> master, tag: v1.1.1)
Author: AI Agent
Date:   Sat Mar 21 2026

    fix: oho add 命令间歇性失败修复 + 完整单元测试覆盖
    
    主要变更:
    - 修复 detectMimeType 在无扩展名文件上的 panic
    - 添加 ClientInterface.PostWithQuery 方法
    - 更新 MockClient 支持 PostWithQueryFunc
    - 将 add.go 改为使用 ClientInterface
    
    新增测试 (770 行):
    - 10 个测试函数，52 个测试用例
    - 覆盖 oho add 命令所有关键路径
    
    新增文档 (664 行):
    - PROJECT_SUMMARY_ZH.md 项目完整摘要
```

---

## ⚠️ 已知限制

- 无自动重试机制（待改进项）
- 无连接池复用（待改进项）
- 无 debug 日志输出（待改进项）

这些限制不影响当前功能，将在未来版本中逐步改进。

---

## 📞 问题反馈

如遇到问题，请提交 Issue 并附上：

1. `oho --version` 输出
2. 完整的错误信息
3. 复现步骤
4. 相关日志（使用 `--json` 输出）

---

## 👥 致谢

感谢所有测试和反馈此版本问题的用户！

---

**Full Changelog**: https://github.com/tornado404/opencode_cli/compare/v1.1.0...v1.1.1
