# OpenCode CLI

OpenCode 命令行工具集，包含 OpenCode Server 的 Go 语言客户端实现。

## 项目结构

```
.
├── oho/                    # OpenCode Server Go 客户端
│   ├── cmd/                # 命令行入口
│   ├── internal/           # 内部包
│   ├── go.mod              # Go 模块定义
│   └── README.md           # 客户端详细文档
├── docs/                   # 项目文档
│   └── plans/              # 设计计划
├── AGENTS.md               # AI 开发指南
└── LICENSE                 # GPL v3 许可证
```

## 快速开始

### 安装 oho (OpenCode Server 客户端)

```bash
cd oho
make build
```

详细使用说明请参考 [oho/README.md](oho/README.md)

## 许可证

本项目基于 GPL v3 许可证开源，详见 [LICENSE](LICENSE) 文件。

## 贡献

欢迎提交 Issue 和 Pull Request！
