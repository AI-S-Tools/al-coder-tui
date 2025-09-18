# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AI CLI Manager is a sophisticated Terminal User Interface (TUI) application written in Go that manages AI-powered CLI tools, their configurations, and MCP (Model Context Protocol) server setups. It provides centralized management with GitHub synchronization and automatic MCP server configuration for Claude Desktop.

## Architecture

The application follows a modular architecture using the Bubble Tea framework:

### Core Structure
- `cmd/ai-cli-manager/main.go` - Application entry point
- `src/model.go` - Core data structures and main model
- `src/handlers.go` - Input handling and state management
- `src/github.go` - GitHub integration for configuration sync
- `src/mcp.go` - MCP server configuration management
- `src/config.go` - Tool definitions and configuration loading

### Key Features
1. **Tool Management**: Install and manage AI CLI tools (Claude, Aider, Continue, Cursor, etc.)
2. **GitHub Sync**: Store and sync configurations using GitHub gists
3. **MCP Integration**: Automatically configure MCP servers for Claude Desktop
4. **Repository-based Installation**: Install tools directly from GitHub repositories

### Data Flow
- Tools configuration stored in `~/.ai-cli-manager/tools.json`
- GitHub configuration in `~/.ai-cli-manager/config.json`
- MCP servers configured in `~/Library/Application Support/Claude/claude_desktop_config.json`

## Commands

```bash
# Build the application
make build

# Run directly
make run

# Install globally
make install

# Development mode with auto-reload
make dev

# Install dependencies
make deps

# Run tests
make test

# Format code
make fmt
```

## Key Components

### AITool Structure
- Supports standard installation commands
- GitHub repository-based installation
- MCP server configurations per tool
- Custom configuration parameters

### Modes
- **menu**: Main menu with options
- **table**: Interactive tool list with installation status
- **config**: GitHub synchronization settings
- **mcp**: MCP server configuration management
- **installing**: Installation progress display

### Navigation
- **Default Start**: Application starts in Table View (main interface)
- **Table Mode**: Arrow keys to navigate, Enter to install, M for MCP config, R to refresh, Esc for menu
- **Menu Mode**: Number keys (1-5) to select options, 1 or Esc to return to table
- **All Modes**: Q to quit

## MCP Server Integration

The app automatically configures MCP servers for tools that support them:
- Detects Claude Desktop configuration location
- Merges MCP server configurations
- Preserves existing configurations
- Supports environment variables and arguments

## GitHub Integration

Uses GitHub CLI (`gh`) for configuration management:
- Creates private gists for configuration storage
- Syncs tool configurations across machines
- Supports pull/push operations
- Requires `gh` CLI to be installed and authenticated