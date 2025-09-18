package src

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "1", "esc":
		m.mode = "table"
		m.updateTable()
		return m, nil
	case "2":
		return m, m.installAll()
	case "3":
		m.mode = "config"
		return m, nil
	case "4":
		m.mode = "mcp"
		return m, nil
	case "5":
		m.mode = "table"
		return m, checkInstallations(m.tools)
	}
	return m, nil
}

func (m Model) handleTableInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "q", "esc":
		m.mode = "menu"
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		selected := m.table.Cursor()
		if selected < len(m.tools) {
			if !m.tools[selected].Installed {
				return m, m.installTool(m.tools[selected])
			} else {
				m.message = "Tool already installed"
			}
		}
	case "m", "M":
		selected := m.table.Cursor()
		if selected < len(m.tools) && len(m.tools[selected].MCPServers) > 0 {
			return m, m.configureMCPServers(m.tools[selected])
		}
	case "r", "R":
		return m, checkInstallations(m.tools)
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) handleConfigInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.mode = "menu"
		return m, nil
	case "s", "S":
		return m, m.syncWithGitHub()
	case "p", "P":
		return m, m.pullFromGitHub()
	}
	return m, nil
}

func (m Model) handleMCPInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.mode = "menu"
		return m, nil
	case "a", "A":
		// Install all MCP servers
		return m, m.installAllMCPServers()
	}
	return m, nil
}

func checkInstallations(tools []AITool) tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		for i := range tools {
			tools[i].Installed = isInstalled(tools[i])
		}
		return checkCompleteMsg{}
	})
}

func isInstalled(tool AITool) bool {
	if tool.CheckCmd == "" {
		cmd := exec.Command("which", tool.CLICommand)
		return cmd.Run() == nil
	}

	parts := strings.Fields(tool.CheckCmd)
	if len(parts) == 0 {
		return false
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	return cmd.Run() == nil
}

func (m Model) installSelected() tea.Cmd {
	return func() tea.Msg {
		for _, tool := range m.tools {
			if !tool.Installed {
				return m.installTool(tool)()
			}
		}
		return installMsg{success: false, err: fmt.Errorf("all tools are already installed")}
	}
}

func (m Model) installAll() tea.Cmd {
	return func() tea.Msg {
		m.installing = true
		m.installAllMode = true

		var cmds []tea.Cmd
		for _, tool := range m.tools {
			if !tool.Installed {
				cmds = append(cmds, m.installTool(tool))
			}
		}

		if len(cmds) == 0 {
			return installMsg{success: false, err: fmt.Errorf("all tools are already installed")}
		}

		// Execute installations sequentially
		for _, cmd := range cmds {
			result := cmd()
			if msg, ok := result.(installMsg); ok && !msg.success {
				return msg
			}
		}

		return installMsg{success: true}
	}
}

func (m Model) installTool(tool AITool) tea.Cmd {
	return func() tea.Msg {
		m.installing = true
		m.message = fmt.Sprintf("Installing %s...", tool.Name)

		// If tool has a GitHub repo, clone and install from there
		if tool.GitHubRepo != "" {
			if err := m.installFromGitHub(tool); err == nil {
				return installMsg{
					tool:    tool,
					success: true,
					err:     nil,
				}
			}
		}

		// Fallback to standard install command
		parts := strings.Fields(tool.InstallCmd)
		if len(parts) == 0 {
			return installMsg{
				tool:    tool,
				success: false,
				err:     fmt.Errorf("no install command specified"),
			}
		}

		cmd := exec.Command(parts[0], parts[1:]...)
		err := cmd.Run()

		return installMsg{
			tool:    tool,
			success: err == nil,
			err:     err,
		}
	}
}

func (m *Model) updateTable() {
	rows := []table.Row{}
	for i, tool := range m.tools {
		status := "Not installed"
		if tool.Installed {
			status = "✅ Installed"
		} else {
			status = "❌ Missing"
		}

		mcpStatus := "-"
		if len(tool.MCPServers) > 0 {
			mcpStatus = fmt.Sprintf("✅ %d", len(tool.MCPServers))
		}

		// Truncate description if too long
		description := tool.Description
		if len(description) > 35 {
			description = description[:32] + "..."
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", i+1),
			tool.Name,
			tool.CLICommand,
			status,
			mcpStatus,
			description,
		})
	}
	m.table.SetRows(rows)
}