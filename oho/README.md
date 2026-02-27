# oho - OpenCode CLI

> Make OpenCode a command-line tool that can be invoked and supervised by other AI

[![GitHub Stars](https://img.shields.io/github/stars/tornado404/opencode_cli?style=flat-square)](https://github.com/tornado404/opencode_cli/stargazers)
[![License](https://img.shields.io/badge/license-GPLv3-blue?style=flat-square)](LICENSE)

oho is the command-line client tool for OpenCode Server, providing complete access to the OpenCode Server API.

## Project Positioning

### Unique Value

**oho** is the **only command-line client implemented entirely in Bash** within the [OpenCode Ecosystem](https://opencode.ai/docs/ecosystem/).

> "oho is callable from Bash" represents powerful extensibility and compatibility — this is the project's unique positioning.

### Design Goals

Make OpenCode more accessible for invocation and supervision by other AI:

- 🤖 Natively callable by any AI Agent
- 🔄 Integrated into automated workflows
- 📦 Run in CI/CD pipelines
- 🔗 Seamlessly combined with other shell tools

### Unique Linux Capabilities

In Linux environments, oho can provide capabilities that OpenCode CLI doesn't currently support:

| Feature | Description |
|---------|-------------|
| 📁 Create Session in Specified Directory | Start AI programming sessions in any directory |
| 💬 Continue Sending Messages Based on Session | Resume previous session context |
| 🗑️ Destroy Session | Complete lifecycle management for sessions |
| 🔄 Session Fork and Revert | Easy switching for experimental development |

## Interface Preview

![oho CLI](assets/oho_cli.png)

## Features

- ✅ Complete API mapping and封装
- ✅ HTTP Basic Auth authentication support
- ✅ JSON/Text dual output mode
- ✅ Configuration file and environment variable support
- ✅ All session management operations
- ✅ Message sending and management
- ✅ File and symbol lookup
- ✅ TUI interface control
- ✅ MCP/LSP/Formatter status management

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/tornado404/opencode_cli.git
cd opencode_cli/oho

# Build
make build

# Or build Linux version
make build-linux
```

### Dependencies

- Go 1.21+
- Cobra CLI framework
- Standard library net/http

## Quick Start

### 1. Configure Server Connection

```bash
# Using environment variables
export OPENCODE_SERVER_HOST=127.0.0.1
export OPENCODE_SERVER_PORT=4096
export OPENCODE_SERVER_PASSWORD=your-password

# Or use command-line flags
oho --host 127.0.0.1 --port 4096 --password your-password session list
```

### 2. Basic Usage

```bash
# Check server health
oho global health

# List all sessions
oho session list

# Create a new session
oho session create

# Create session in specified directory
oho session create --path /your/project

# Send a message
oho message add -s <session-id> "Hello, please help me analyze this project"

# Continue existing session
oho message add -s <session-id> "Continue the previous task"

# View message list
oho message list -s <session-id>

# Destroy session
oho session delete <session-id>

# Get configuration
oho config get

# List providers
oho provider list
```

## Comparison with Other Ecosystem Projects

| Feature | oho | Other Ecosystem Projects |
|---------|-----|-------------------------|
| Implementation Language | Bash | TypeScript/Python/Go |
| AI Callable | ✅ Native support | Requires additional adaptation |
| Cross-platform | Linux/Mac/Windows | Runtime dependent |
| Integration Difficulty | ⭐⭐⭐⭐⭐ Extremely Low | ⭐⭐⭐ Medium |

Reference: [Other Projects in OpenCode Ecosystem](https://opencode.ai/docs/ecosystem/)

## Command Reference

### Global Commands

```bash
oho global health          # Check server health status
oho global event           # Listen to global event stream (SSE)
```

### Project Management

```bash
oho project list           # List all projects
oho project current        # Get current project
oho path                   # Get current path
oho vcs                    # Get VCS information
oho instance dispose       # Dispose current instance
```

### Session Management

```bash
oho session list                      # List all sessions
oho session create                    # Create new session
oho session create --path /path        # Create session in specified directory
oho session status                    # Get all session statuses
oho session get <id>                  # Get session details
oho session delete <id>               # Delete session
oho session update <id> --title "New Title"  # Update session
oho session children <id>             # Get child sessions
oho session todo <id>                 # Get todo items
oho session fork <id>                 # Fork session
oho session abort <id>                # Abort session
oho session share <id>                # Share session
oho session unshare <id>              # Unshare session
oho session diff <id>                 # Get file diff
oho session summarize <id>            # Summarize session
oho session revert <id> --message <msg-id>  # Revert message
oho session unrevert <id>             # Undo revert
oho session permissions <id> <perm-id> --response allow  # Respond to permission
```

### Message Management

```bash
oho message list -s <session>         # List messages
oho message add -s <session> "content"   # Send message
oho message get -s <session> <msg-id> # Get message details
oho message prompt-async -s <session> "content"  # Send async
oho message command -s <session> /help  # Execute command
oho message shell -s <session> --agent default "ls -la"  # Run shell
```

### Configuration Management

```bash
oho config get                      # Get configuration
oho config set --theme dark         # Update configuration
oho config providers                # List providers and default models
```

### Provider Management

```bash
oho provider list                   # List all providers
oho provider auth                   # Get authentication methods
oho provider oauth authorize <id>   # OAuth authorize
oho provider oauth callback <id>    # Handle callback
```

### File Operations

```bash
oho file list [path]                # List files
oho file content <path>             # Read file content
oho file status                     # Get file status
```

### Find Features

```bash
oho find text "pattern"             # Search text
oho find file "query"               # Find files
oho find symbol "query"             # Find symbols
```

### Other Commands

```bash
oho agent list                      # List agents
oho command list                    # List commands
oho tool ids                        # List tool IDs
oho tool list --provider xxx --model xxx  # List tools
oho lsp status                      # LSP status
oho formatter status                # Formatter status
oho mcp list                        # List MCP servers
oho mcp add <name> --config '{}'    # Add MCP server
oho tui open-help                   # Open help
oho tui show-toast --message "message"  # Show toast
oho auth set <provider> --credentials key=value  # Set authentication
```

## Output Format

Use `-j` or `--json` flags for JSON output:

```bash
oho session list -j
oho config get --json
```

## Configuration File

Configuration file is located at `~/.config/oho/config.json`:

```json
{
  "host": "127.0.0.1",
  "port": 4096,
  "username": "opencode",
  "password": "",
  "json": false
}
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `OPENCODE_SERVER_HOST` | Server host | `127.0.0.1` |
| `OPENCODE_SERVER_PORT` | Server port | `4096` |
| `OPENCODE_SERVER_USERNAME` | Username | `opencode` |
| `OPENCODE_SERVER_PASSWORD` | Password | empty |

## Development

```bash
# Run
go run ./cmd/oho --help

# Test
go test ./...

# Format
go fmt ./...

# Clean
make clean
```

## Project Structure

```
oho/
├── cmd/
│   └── oho/
│       ├── main.go           # Entry file
│       ├── root.go           # Root command
│       ├── cmd/              # Subcommands
│       │   ├── global/
│       │   ├── project/
│       │   ├── session/
│       │   ├── message/
│       │   ├── configcmd/
│       │   ├── provider/
│       │   ├── file/
│       │   ├── find/
│       │   ├── tool/
│       │   ├── agent/
│       │   ├── command/
│       │   ├── lsp/
│       │   ├── formatter/
│       │   ├── mcp/
│       │   ├── tui/
│       │   └── auth/
│       └── internal/
│           ├── client/       # HTTP client
│           ├── config/       # Configuration management
│           ├── types/        # Type definitions
│           └── util/         # Utility functions
├── Makefile
├── build.sh
└── README.md
```

## License

GPL v3 License - See project root [LICENSE](../LICENSE)

## References

- [OpenCode Official Documentation](https://opencode.ai/docs/)
- [OpenCode Ecosystem](https://opencode.ai/docs/ecosystem/)
- [OpenCode GitHub](https://github.com/anomalyco/opencode)

## Contributing

Issues and Pull Requests are welcome!
