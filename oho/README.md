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
- 📊 **[API Completion Status](./COMPLETION.md)** - View implementation coverage

## Installation

### Quick Install (Recommended)

```bash
# Download and run the installer
curl -sSL https://raw.githubusercontent.com/tornado404/opencode_cli/master/oho/install.sh | bash

# Or with wget
wget -qO- https://raw.githubusercontent.com/tornado404/opencode_cli/master/oho/install.sh | bash
```

### Build from Source

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
oho session list --id ses_abc         # Filter by ID (fuzzy match)
oho session list --title "测试"        # Filter by title (fuzzy)
oho session list --project-id proj1   # Filter by project ID
oho session list --directory babylon  # Filter by directory
oho session list --created 1773537883643  # Filter by created timestamp
oho session list --updated 1773538142930  # Filter by updated timestamp
oho session list --sort updated --order desc  # Sort by updated (desc)
oho session list --limit 10 --offset 0  # Pagination
oho session create                    # Create new session
oho session create --title "名称"       # Create with custom title
oho session create --parent ses_xxx   # Create child session
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

**List Command Flags**:

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--id` | string | Filter by ID (fuzzy, case-insensitive) | - |
| `--title` | string | Filter by title (fuzzy, case-insensitive) | - |
| `--created` | int64 | Filter by created timestamp (exact) | - |
| `--updated` | int64 | Filter by updated timestamp (exact) | - |
| `--project-id` | string | Filter by project ID (fuzzy) | - |
| `--directory` | string | Filter by directory (fuzzy) | - |
| `--status` | string | Filter by status (running/completed/error/aborted/idle) | - |
| `--running` | bool | Show only running sessions | false |
| `--sort` | string | Sort field (created/updated) | updated |
| `--order` | string | Sort order (asc/desc) | desc |
| `--limit` | int | Limit results count | - |
| `--offset` | int | Pagination offset | 0 |
| `-j, --json` | bool | JSON output format | false |

### Message Management

```bash
oho message list -s <session>         # List messages
oho message add -s <session> "content"   # Send message
oho message get -s <session> <msg-id> # Get message details
oho message prompt-async -s <session> "content"  # Send async
oho message command -s <session> /help  # Execute command
oho message shell -s <session> --agent default "ls -la"  # Run shell
```

### Quick Start (Session + Message)

```bash
oho add "帮我分析这个项目"                    # Create session and send message
oho add "修复登录 bug" --title "Bug 修复"       # Create session with custom title
oho add "测试功能" --no-reply --agent default  # Don't wait for AI response
oho add "分析日志" --file /var/log/app.log     # Attach file to message
oho add "任务描述" --directory /path/to/project # Specify working directory
oho add "消息内容" --json                      # Output in JSON format
```

### ⚠️ Timeout Considerations

The `oho add` command waits for the AI response by default. For complex tasks, the AI may need extended time to think, which could result in a timeout.

**Methods to Avoid Timeouts**:

1. **Use `--no-reply` flag** (Recommended):
   ```bash
   # Send message and return immediately without waiting for AI response
   oho add "Analyze project structure" --no-reply
   
   # Check results later
   oho message list -s <session-id>
   ```

2. **Increase timeout duration**:
   ```bash
   # Set timeout to 10 minutes (600 seconds)
   export OPENCODE_CLIENT_TIMEOUT=600
   oho add "Complex task"
   
   # Or set temporarily
   OPENCODE_CLIENT_TIMEOUT=600 oho add "Complex task"
   ```

3. **Use async command**:
   ```bash
   # Create session first
   oho session create --title "Task"
   
   # Send message asynchronously
   oho message prompt-async -s <session-id> "Task description"
   ```

**Timeout Configuration**:
| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `OPENCODE_CLIENT_TIMEOUT` | 300 seconds | HTTP request timeout (seconds) |

**`oho add` Flags**:

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--title` | string | Session title (auto-generated if not provided) | auto-generated |
| `--parent` | string | Parent session ID (for creating sub-session) | - |
| `--directory` | string | Working directory for the session | current directory |
| `--agent` | string | Agent ID for message | - |
| `--model` | string | Model ID for message (e.g., `provider:model`) | default model |
| `--no-reply` | bool | Don't wait for AI response | false |
| `--system` | string | System prompt | - |
| `--tools` | string[] | Tools list (can be specified multiple times) | - |
| `--file` | string[] | File attachments (can be specified multiple times) | - |
| `-j, --json` | bool | Output in JSON format | false |

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
```

## MCP Server

oho can be used as an MCP (Model Context Protocol) server, allowing external MCP clients (like Claude Desktop, Cursor, etc.) to call OpenCode API through MCP protocol.

### Start MCP Server

```bash
# Start the MCP server (stdio mode)
oho mcpserver
```

The MCP server uses stdio transport, which is the standard mode for local MCP clients.

### Available MCP Tools

The following tools are available via MCP:

| Tool | Description |
|------|-------------|
| `session_list` | List all sessions |
| `session_create` | Create a new session |
| `session_get` | Get session details |
| `session_delete` | Delete a session |
| `session_status` | Get all session statuses |
| `message_list` | List messages in a session |
| `message_add` | Send a message to a session |
| `config_get` | Get OpenCode configuration |
| `project_list` | List all projects |
| `project_current` | Get current project |
| `provider_list` | List AI providers |
| `file_list` | List files in a directory |
| `file_content` | Read file content |
| `find_text` | Search text in project |
| `find_file` | Find files by name |
| `global_health` | Check server health |

### MCP Client Configuration

#### Claude Desktop (macOS/Windows)

Add to `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "oho": {
      "command": "oho",
      "args": ["mcpserver"],
      "env": {
        "OPENCODE_SERVER_HOST": "127.0.0.1",
        "OPENCODE_SERVER_PORT": "4096",
        "OPENCODE_SERVER_PASSWORD": "your-password"
      }
    }
  }
}
```

#### Cursor

Add to Cursor settings (JSON mode):

```json
{
  "mcpServers": {
    "oho": {
      "command": "oho",
      "args": ["mcpserver"],
      "env": {
        "OPENCODE_SERVER_HOST": "127.0.0.1",
        "OPENCODE_SERVER_PORT": "4096",
        "OPENCODE_SERVER_PASSWORD": "your-password"
      }
    }
  }
}
```

#### VS Code (with Copilot Free)

VS Code doesn't have native MCP support. Use the [MCP VS Code extension](https://github.com/modelcontextprotocol/servers):

```json
{
  "mcpServers": {
    "oho": {
      "command": "oho",
      "args": ["mcpserver"]
    }
  }
}
```

## License
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
