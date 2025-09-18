# AI CLI Manager

A sophisticated Terminal User Interface (TUI) application for managing AI-powered CLI tools, their configurations, and MCP (Model Context Protocol) server setups.

## Features

- **🤖 Tool Management**: Install and manage popular AI CLI tools
  - Claude Code
  - Aider
  - Continue
  - GitHub Copilot CLI
  - Cursor
  - Codeium
  - Qodo
  - Windsurf

- **☁️ GitHub Sync**: Store and sync configurations using GitHub gists
- **🔌 MCP Integration**: Automatically configure MCP servers for Claude Desktop
- **📦 Smart Installation**: Install tools from package managers or GitHub repositories
- **🎨 Beautiful TUI**: Interactive interface built with Bubble Tea framework

## Installation

### Prerequisites

- Go 1.21 or higher
- GitHub CLI (`gh`) for configuration sync
- macOS (for MCP server configuration)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/lpm/ai-cli-manager
cd ai-cli-manager

# Install dependencies
make deps

# Build the application
make build

# Install globally (optional)
make install
```

## Usage

### Run the Application

```bash
# Run directly
make run

# Or if installed globally
ai-cli-manager
```

### Navigation

#### Main Menu
- **1**: View tools table
- **2**: Install selected tools
- **3**: Install all tools
- **4**: Configure GitHub sync
- **5**: Configure MCP servers
- **6**: Update installation status
- **Q**: Quit

#### Table View
- **↑/↓**: Navigate through tools
- **Enter**: Install selected tool
- **M**: Configure MCP for selected tool
- **R**: Refresh installation status
- **Esc**: Return to menu
- **Q**: Quit

## Configuration

### Tool Configuration
Tools are defined in `~/.ai-cli-manager/tools.json`. You can customize this file to add or modify tools.

Example tool configuration:
```json
{
  "name": "Claude Code",
  "cli_command": "claude",
  "install_cmd": "npm install -g @anthropic/claude-cli",
  "check_cmd": "claude --version",
  "description": "Anthropic's Claude AI coding assistant",
  "github_repo": "https://github.com/anthropics/claude-cli",
  "mcp_servers": [
    {
      "name": "filesystem",
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem"],
      "description": "File system access for Claude"
    }
  ]
}
```

### GitHub Synchronization
The app can sync your tool configurations to a private GitHub gist for backup and cross-machine synchronization.

1. Install and authenticate GitHub CLI: `gh auth login`
2. Configure sync in the app (option 4 from main menu)
3. Use sync options to push/pull configurations

### MCP Server Configuration
MCP servers are automatically configured in Claude Desktop's configuration file:
`~/Library/Application Support/Claude/claude_desktop_config.json`

The app will:
- Detect existing MCP configurations
- Merge new configurations without overwriting
- Support environment variables and custom arguments

## Development

### Project Structure
```
.
├── cmd/ai-cli-manager/    # Application entry point
├── src/                    # Source code
│   ├── model.go           # Core data structures
│   ├── handlers.go        # Input handling
│   ├── github.go          # GitHub integration
│   ├── mcp.go            # MCP configuration
│   └── config.go         # Tool definitions
├── Makefile              # Build commands
└── CLAUDE.md            # Documentation for Claude Code
```

### Available Make Commands
```bash
make build    # Build the application
make run      # Run the application
make install  # Install to /usr/local/bin
make deps     # Install dependencies
make clean    # Clean build artifacts
make test     # Run tests
make dev      # Development mode with auto-reload
make fmt      # Format code
make lint     # Run linter
make vuln     # Check for vulnerabilities
```

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

MIT License - See LICENSE file for details

## Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework
- Uses [GitHub CLI](https://cli.github.com/) for configuration sync
- Supports [MCP (Model Context Protocol)](https://modelcontextprotocol.io/) for Claude Desktop integration