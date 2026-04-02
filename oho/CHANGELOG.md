# Changelog

所有重要的版本更新都会在此记录。遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/) 规范。

---

## [v1.1.2] - 2026-04-02

### Added

- **config 优化**: 优化配置文件搜索路径和加载逻辑 (`e13d2f1`)
- **config 测试**: 添加配置初始化和基础URL的单元测试 (`a5b18dc`)
- **超时参数**: `oho add` 命令新增 `--timeout` 参数 (`d1b7a1c`)
- **no-reply 追踪**: 添加 no-reply 模式 sessionID 追踪任务 (`44a18f6`)
- **OpenCode API 技能手册**: 添加 OpenCode Server API 技能手册 (SKILL.md) (`7f0e9e8`)

### Fixed

- **Init() 返回值检查**: 检查 Init() 返回值修复 errcheck lint 错误 (`73d0699`)
- **no-reply 模式修复**: 使用 /prompt_async 端点修复 --no-reply 模式 (`bdad95e`)
- **CI ldflags 修复**: 修复 GitHub Actions workflow 中缺少 ldflags 的问题 (`76b0953`)
- **测试用例修复**: add missing wantMsgID assertions in TestSendMessage (`3ad01ce`)
- **测试用例修复**: correct TestSendMessage test cases for noReply and file not found (`0f44c33`)
- **lint 修复**: add error checks in add_test.go for golangci-lint (`96742f9`)
- **config 帮助文本**: 更新 config set 命令帮助文本和提示信息 (`9580e5f`)

### Documentation

- **构建指南**: 添加 oho 构建指南并修复 Makefile build-linux 目标 (`22a7ab4`)
- **发布说明**: 添加 v1.1.1 发布说明 (`99644ae`)

---

## [v1.1.1] - 2026-03-21

### Added

- **oho add 命令**: 一键创建会话并发送消息 (`aed8700`, `bcb63e2`)
- **session submit 命令**: 支持一键提交任务 (`cf21aa0`)
- **session 列表增强**: 列表输出增加目录和项目信息 (`77baaf9`)
- **完整版本信息**: 始终显示完整版本信息 (`fff40c6`)
- **超时文档**: 添加超时问题完整解决方案总结 (`a37a63b`)
- **快速上手指南**: 添加快速上手指南 (`ac090e3`)
- **安装说明**: 添加安装说明和超时测试脚本 (`4e9c0b2`)
- **问题排查文档**: 添加问题排查文档和超时配置支持 (`5cc4817`)
- **session create 增强**: 为 session create 添加 --title 和 --parent 参数 (`8c96383`)

### Fixed

- **submit directory 参数**: submit 命令使用 directory query 参数创建会话 (`e2d612c`)
- **submit 文档**: submit 命令添加 OpenCode Server path 参数限制说明 (`d43acc7`)
- **add 消息格式**: oho add 命令消息输出格式 (`1a8e56e`)
- **add 间歇性失败**: oho add 命令间歇性失败修复 + 完整单元测试覆盖 (`6db26d9`)
- **session create 文档**: 更新 session create 参数文档 (`0674e74`)

### Changed

- **构建流程**: 添加版本信息注入到构建流程 (`449d662`)

### Chore

- **任务文件清理**: 清理已提交的任务文件 (`4753ee5`)

---

## [v1.1.0] - 2026-03-17

### Added

- **初始功能集**: 基础的 session、message、config、provider 等命令
- **HTTP Basic Auth**: 支持 HTTP Basic Auth 认证
- **JSON/Text 双输出模式**: 支持 JSON 和文本格式输出
- **配置文件支持**: 支持配置文件和环境变量配置
- **TUI 控制**: 支持 Toast、Help 等 TUI 控制命令
- **MCP 服务器**: 支持 MCP (Model Context Protocol) 服务器模式

### Features

- 会话管理 (创建/删除/列出/更新/分叉/回退)
- 消息发送 (同步/异步/命令/shell)
- 项目和文件操作
- 提供商管理
- LSP/Formatter 状态查询
- MCP 服务器管理

---

## [v0.0.2] - 更早版本

### Added

- 初始版本基础功能

---

## [v0.0.1] - 首次发布

- 首次发布
