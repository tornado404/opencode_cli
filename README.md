# OpenCode CLI
<p align="center">
  <a href="README.md">English</a> |
  <a href="README_zh.md">简体中文</a>
</p>

> Make OpenCode a command-line tool that can be invoked and supervised by other AI

This project is an extension of the OpenCode ecosystem, providing a complete Bash-based command-line client that enables OpenCode to be invoked and supervised by other AI systems.

## Vision

**oho** is designed to make OpenCode more accessible for invocation and supervision by other AI. Within the [OpenCode Ecosystem](https://opencode.ai/docs/ecosystem/), there are many similar applications, but **oho is the only solution implemented entirely in Bash**.

> "oho is callable from Bash" represents powerful extensibility and compatibility — this is the project's unique positioning.

## Core Features

### Complete API Coverage

oho is built on the OpenCode REST API, providing a complete command-line interface:

- ✅ Session management (create, delete, continue, terminate)
- ✅ Message sending and receiving
- ✅ Project and file operations
- ✅ Configuration and provider management
- ✅ MCP/LSP/Formatter status management

### Unique Linux Capabilities

In Linux environments, oho can provide capabilities that OpenCode CLI doesn't currently support:

- 📁 **Create Session in Specified Directory**: Start AI programming sessions in any directory
- 💬 **Continue Sending Messages Based on Session**: Resume previous session context
- 🗑️ **Destroy Session**: Complete lifecycle management for sessions
- 🔄 **Session Fork and Revert**: Easy switching for experimental development

### Bash Callable

As a pure Bash implementation, oho can be:

- Invoked by any AI Agent
- Integrated into automated workflows
- Run in CI/CD pipelines
- Seamlessly combined with other shell tools

## Quick Start

### Installation

```bash
cd oho
make build
```

### Basic Usage

```bash
# Configure server connection
export OPENCODE_SERVER_HOST=127.0.0.1
export OPENCODE_SERVER_PORT=4096
export OPENCODE_SERVER_PASSWORD=your-password

# List all sessions
oho session list

# Create a new session
oho session create

# Create session in specified directory
oho session create --path /your/project

# Send a message
oho message add -s <session-id> "Help me analyze this code"

# Continue existing session
oho message add -s <session-id> "Continue the previous task"

# Destroy session
oho session delete <session-id>
```

## Comparison with Other Ecosystem Projects

| Feature | oho | Other Ecosystem Projects |
|---------|-----|-------------------------|
| Implementation Language | Bash | TypeScript/Python/Go |
| AI Callable | ✅ Native support | Requires additional adaptation |
| Integration Difficulty | ⭐⭐⭐⭐⭐ Extremely Low | ⭐⭐⭐ Medium |

## Project Structure

```
.
├── oho/                    # OpenCode Bash Client
│   ├── cmd/                # Command-line entry
│   ├── internal/           # Internal packages
│   ├── go.mod              # Go module definition
│   └── README.md           # Client detailed documentation
├── docs/                   # Project documentation
│   └── plans/              # Design plans
├── assets/                 # Resource files
│   └── oho_cli.png         # CLI screenshot
├── AGENTS.md               # AI Development Guide
└── LICENSE                 # GPL v3 License
```

## Command Reference

For complete command list, see [oho/README.md](oho/README.md)

## License

This project is open source under the GPL v3 license. See [LICENSE](LICENSE) for details.

## References

- [OpenCode Official Documentation](https://opencode.ai/docs/)
- [OpenCode Ecosystem](https://opencode.ai/docs/ecosystem/)
- [OpenCode GitHub](https://github.com/anomalyco/opencode)
