# Multi-Language README Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create English (default) and Chinese (`_zh` suffix) versions of both README files for the OpenCode CLI project.

**Architecture:** Translation-based approach preserving all formatting, tables, code blocks, and links. English becomes the default/main README, Chinese versions get `_zh` suffix.

**Tech Stack:** Manual translation with web search assistance for CLI terminology consistency.

---

## Task Overview

| File | Current | Target |
|------|---------|--------|
| `README.md` | Chinese | English (new default) |
| `README_zh.md` | - | Chinese (keep original) |
| `oho/README.md` | Chinese | English (new default) |
| `oho/README_zh.md` | - | Chinese (keep original) |

---

## Translation Strategy

### Approach
1. **Manual translation** with web search for CLI terminology
2. **Preserve all formatting**: tables, code blocks, links, badges
3. **Professional CLI documentation standards**
4. **Consistent terminology** across both files

### Key Terminology Mapping

| Chinese | English |
|---------|---------|
| 命令行工具 | CLI tool |
| 命令行客户端 | command-line client |
| 会话 | session |
| 消息 | message |
| 提供商 | provider |
| 配置文件 | config file |
| 环境变量 | environment variable |
| 源码编译 | build from source |
| 子命令 | subcommand |

---

## Execution Tasks

### Task 1: Translate Root README.md to English

**Files:**
- Read: `README.md` (existing Chinese)
- Create: `README.md` (new English)
- Create: `README_zh.md` (copy original Chinese)

**Step 1: Read existing Chinese README**

```bash
# Already read - 116 lines, contains:
# - Project vision
# - Core features (API coverage, Linux capabilities, Bash callable)
# - Quick start (install, basic usage)
# - Comparison table
# - Project structure
# - Command reference
# - License, resources
```

**Step 2: Create English version (new default)**

Translating these key sections:

| Section | Key Terms to Translate |
|---------|----------------------|
| 项目愿景 → Project Vision | 愿景→vision, 调用→invoke, 监督→monitor |
| 核心特性 → Core Features | API覆盖→API coverage, Linux能力→Linux capabilities |
| 完整的 API 覆盖 → Complete API Coverage | 会话管理→session management, 消息发送→message sending |
| 独特的 Linux 能力 → Unique Linux Capabilities | 指定目录→specify directory, 会话分叉→session fork |
| Bash 可调用性 → Bash Callability | AI Agent, 自动化工作流→automated workflow |
| 快速开始 → Quick Start | 安装→installation, 基本用法→basic usage |
| 与其他生态项目的对比 → Comparison | 实现语言→implementation language, 集成难度→integration difficulty |
| 项目结构 → Project Structure | 目录结构→directory structure |

**Step 3: Verify English README**

Checklist:
- [ ] All headings translated
- [ ] All code blocks preserved
- [ ] All links working (URLs unchanged)
- [ ] Tables format intact
- [ ] Emoji preserved where appropriate
- [ ] Chinese specific URLs changed to English versions where available

**Step 4: Create README_zh.md**

```bash
cp README.md README_zh.md
```

**Step 5: Commit**

```bash
git add README.md README_zh.md
git commit -m "docs: add English README, preserve Chinese as README_zh"
```

---

### Task 2: Translate oho/README.md to English

**Files:**
- Read: `oho/README.md` (existing Chinese - 332 lines)
- Create: `oho/README.md` (new English)
- Create: `oho/README_zh.md` (copy original Chinese)

**Step 1: Read existing Chinese oho README**

```bash
# Already read - 332 lines, contains:
# - Project positioning (unique value, design goals, Linux capabilities)
# - Interface preview
# - Features list
# - Installation (build from source)
# - Quick start (configuration, basic usage)
# - Comparison table
# - Comprehensive command reference (15+ command categories)
# - Output format
# - Config file
# - Environment variables
# - Development commands
# - Project structure
# - License, resources, contribution
```

**Step 2: Create English version (new default)**

Translating 15 command categories:

| Category | Chinese | English |
|----------|---------|---------|
| 全局命令 | 全局命令 | Global Commands |
| 项目管理 | 项目管理 | Project Management |
| 会话管理 | 会话管理 | Session Management |
| 消息管理 | 消息管理 | Message Management |
| 配置管理 | 配置管理 | Config Management |
| 提供商管理 | 提供商管理 | Provider Management |
| 文件操作 | 文件操作 | File Operations |
| 查找功能 | 查找功能 | Find Commands |
| 其他命令 | 其他命令 | Other Commands |

**Special attention sections:**

1. **Command examples** - preserve exact command syntax, translate comments only:
   ```bash
   # 中文: # 列出所有会话
   # English: # List all sessions
   ```

2. **Environment variables table** - preserve variable names, translate descriptions:
   | Variable | Description | Default |
   |----------|-------------|---------|
   | `OPENCODE_SERVER_HOST` | Server host | `127.0.0.1` |

3. **Config file** - preserve JSON structure, translate comments

4. **URLs** - change Chinese docs to English where available:
   - `https://opencode.ai/docs/zh-cn/` → `https://opencode.ai/docs/`
   - `https://opencode.ai/docs/zh-cn/ecosystem/` → `https://opencode.ai/docs/ecosystem/`

**Step 3: Verify English README**

Checklist:
- [ ] All 15+ command categories translated
- [ ] Command syntax preserved exactly
- [ ] Comments translated, commands unchanged
- [ ] Tables format intact
- [ ] Badges (GitHub Stars, License) preserved
- [ ] Image reference `assets/oho_cli.png` unchanged
- [ ] URLs converted to English versions

**Step 4: Create oho/README_zh.md**

```bash
cp oho/README.md oho/README_zh.md
```

**Step 5: Commit**

```bash
git add oho/README.md oho/README_zh.md
git commit -m "docs: add English oho README, preserve Chinese as README_zh"
```

---

## Verification Steps

### Post-Translation Verification

1. **File existence check**
   ```bash
   ls -la README.md README_zh.md
   ls -la oho/README.md oho/README_zh.md
   ```

2. **Format verification**
   ```bash
   # Check tables render correctly
   markdownlint README.md || true
   markdownlint oho/README.md || true
   ```

3. **Link verification**
   ```bash
   # Verify no broken internal links
   grep -o ']([^)]*)' README.md | head -10
   grep -o ']([^)]*)' oho/README.md | head -10
   ```

4. **Line count comparison**
   ```bash
   wc -l README.md README_zh.md
   wc -l oho/README.md oho/README_zh.md
   ```

5. **Content spot-check**
   - [ ] Root README: "Project Vision" section exists in English
   - [ ] Root README: "Core Features" section exists in English
   - [ ] oho README: All 15 command categories present
   - [ ] oho README: Environment variables table intact

---

## Summary

| Task | Steps |
|------|-------|
| Task 1: Root README | 5 steps (read, translate EN, verify, copy ZH, commit) |
| Task 2: oho README | 5 steps (read, translate EN, verify, copy ZH, commit) |
| Verification | 5 verification checks |

**Total:** 10 execution steps + 5 verification checks

---

## Execution Choice

**Plan complete and saved to `docs/plans/2026-02-27-multi-language-readme.md`. Two execution options:**

1. **Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks, fast iteration

2. **Parallel Session (separate)** - Open new session with executing-plans, batch execution with checkpoints

**Which approach?**
