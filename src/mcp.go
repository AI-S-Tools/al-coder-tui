package src

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

type ClaudeConfig struct {
	MCPServers map[string]MCPServerEntry `json:"mcpServers,omitempty"`
}

type MCPServerEntry struct {
	Command string            `json:"command"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

func (m Model) configureMCPServers(tool AITool) tea.Cmd {
	return func() tea.Msg {
		if len(tool.MCPServers) == 0 {
			return mcpInstallMsg{
				tool:    tool.Name,
				success: false,
				err:     fmt.Errorf("no MCP servers configured for %s", tool.Name),
			}
		}

		// Read existing Claude config
		config, err := m.readClaudeConfig()
		if err != nil {
			config = &ClaudeConfig{
				MCPServers: make(map[string]MCPServerEntry),
			}
		}

		// Add MCP servers for this tool
		for _, server := range tool.MCPServers {
			serverKey := fmt.Sprintf("%s-%s", tool.Name, server.Name)
			config.MCPServers[serverKey] = MCPServerEntry{
				Command: server.Command,
				Args:    server.Args,
				Env:     server.Env,
			}
		}

		// Write updated config
		if err := m.writeClaudeConfig(config); err != nil {
			return mcpInstallMsg{
				tool:    tool.Name,
				success: false,
				err:     err,
			}
		}

		return mcpInstallMsg{
			tool:    tool.Name,
			success: true,
		}
	}
}

func (m Model) installAllMCPServers() tea.Cmd {
	return func() tea.Msg {
		config, err := m.readClaudeConfig()
		if err != nil {
			config = &ClaudeConfig{
				MCPServers: make(map[string]MCPServerEntry),
			}
		}

		count := 0
		for _, tool := range m.tools {
			if len(tool.MCPServers) > 0 {
				for _, server := range tool.MCPServers {
					serverKey := fmt.Sprintf("%s-%s", tool.Name, server.Name)
					config.MCPServers[serverKey] = MCPServerEntry{
						Command: server.Command,
						Args:    server.Args,
						Env:     server.Env,
					}
					count++
				}
			}
		}

		if count == 0 {
			return mcpInstallMsg{
				tool:    "all",
				success: false,
				err:     fmt.Errorf("no MCP servers to configure"),
			}
		}

		if err := m.writeClaudeConfig(config); err != nil {
			return mcpInstallMsg{
				tool:    "all",
				success: false,
				err:     err,
			}
		}

		return mcpInstallMsg{
			tool:    fmt.Sprintf("all (%d servers)", count),
			success: true,
		}
	}
}

func (m Model) readClaudeConfig() (*ClaudeConfig, error) {
	data, err := os.ReadFile(m.mcpConfigPath)
	if err != nil {
		return nil, err
	}

	var config ClaudeConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (m Model) writeClaudeConfig(config *ClaudeConfig) error {
	// Ensure directory exists
	dir := filepath.Dir(m.mcpConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.mcpConfigPath, data, 0644)
}

func (m Model) viewMCP() string {
	config, err := m.readClaudeConfig()

	serverCount := 0
	if err == nil && config != nil {
		serverCount = len(config.MCPServers)
	}

	// Count tools with MCP servers
	toolsWithMCP := 0
	totalServers := 0
	for _, tool := range m.tools {
		if len(tool.MCPServers) > 0 {
			toolsWithMCP++
			totalServers += len(tool.MCPServers)
		}
	}

	configStatus := fmt.Sprintf("Current MCP servers configured: %d", serverCount)
	availableStatus := fmt.Sprintf("Available MCP servers: %d (from %d tools)", totalServers, toolsWithMCP)

	return fmt.Sprintf(`
%s

MCP Server Configuration
%s
%s

Config location: %s

Options:
%s A: Configure all available MCP servers
%s Esc: Back to menu

MCP servers enable AI tools to interact with local services
and provide enhanced functionality within Claude Desktop.

%s
`,
		titleStyle.Render("MCP Configuration"),
		configStatus,
		availableStatus,
		m.mcpConfigPath,
		selectedStyle.Render("→"),
		selectedStyle.Render("→"),
		m.message,
	)
}