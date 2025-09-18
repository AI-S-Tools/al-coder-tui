package src

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AITool struct {
	Name        string            `json:"name"`
	CLICommand  string            `json:"cli_command"`
	InstallCmd  string            `json:"install_cmd"`
	CheckCmd    string            `json:"check_cmd"`
	Description string            `json:"description"`
	GitHubRepo  string            `json:"github_repo,omitempty"`
	MCPServers  []MCPServerConfig `json:"mcp_servers,omitempty"`
	Config      map[string]string `json:"config,omitempty"`
	Installed   bool              `json:"-"`
}

type MCPServerConfig struct {
	Name        string            `json:"name"`
	Command     string            `json:"command"`
	Args        []string          `json:"args,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Description string            `json:"description"`
}

type Model struct {
	tools          []AITool
	table          table.Model
	selected       int
	mode           string // "menu", "table", "installing", "config", "mcp"
	message        string
	installing     bool
	installAllMode bool
	githubUser     string
	githubRepo     string
	configSynced   bool
	mcpConfigPath  string
}

type installMsg struct {
	tool    AITool
	success bool
	err     error
}

type checkCompleteMsg struct{}

type githubSyncMsg struct {
	success bool
	err     error
}

type mcpInstallMsg struct {
	tool    string
	success bool
	err     error
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(0, 1)

	menuStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true)

	installedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46"))

	notInstalledStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("196"))
)

func NewModel() Model {
	tools := loadAITools()

	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Name", Width: 20},
		{Title: "CLI Command", Width: 15},
		{Title: "Status", Width: 12},
		{Title: "MCP", Width: 8},
		{Title: "Description", Width: 35},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	// Detect MCP config path
	homeDir, _ := os.UserHomeDir()
	mcpPath := filepath.Join(homeDir, "Library", "Application Support", "Claude", "claude_desktop_config.json")

	m := Model{
		tools:         tools,
		table:         t,
		selected:      0,
		mode:          "table",
		message:       "Welcome to AI CLI Manager! Press Esc for menu.",
		mcpConfigPath: mcpPath,
	}

	// Load GitHub config
	m.loadGitHubConfig()

	// Initialize table data
	m.updateTable()

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		checkInstallations(m.tools),
		m.checkGitHubCLI(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case "menu":
			return m.handleMenuInput(msg)
		case "table":
			return m.handleTableInput(msg)
		case "config":
			return m.handleConfigInput(msg)
		case "mcp":
			return m.handleMCPInput(msg)
		case "installing":
			if msg.String() == "q" {
				return m, tea.Quit
			}
		}

	case checkCompleteMsg:
		m.updateTable()
		return m, nil

	case installMsg:
		m.installing = false
		if msg.success {
			m.message = successStyle.Render(fmt.Sprintf("✓ %s installed successfully!", msg.tool.Name))
			for i := range m.tools {
				if m.tools[i].Name == msg.tool.Name {
					m.tools[i].Installed = true
					break
				}
			}
			// After installing, configure MCP if available
			if len(msg.tool.MCPServers) > 0 {
				return m, m.configureMCPServers(msg.tool)
			}
		} else {
			m.message = errorStyle.Render(fmt.Sprintf("✗ Failed to install %s: %v", msg.tool.Name, msg.err))
		}
		m.updateTable()
		return m, nil

	case githubSyncMsg:
		if msg.success {
			m.message = successStyle.Render("✓ Configuration synced with GitHub!")
			m.configSynced = true
		} else {
			m.message = errorStyle.Render(fmt.Sprintf("✗ GitHub sync failed: %v", msg.err))
		}
		return m, nil

	case mcpInstallMsg:
		if msg.success {
			m.message = successStyle.Render(fmt.Sprintf("✓ MCP servers configured for %s!", msg.tool))
		} else {
			m.message = errorStyle.Render(fmt.Sprintf("✗ MCP configuration failed: %v", msg.err))
		}
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	if m.mode == "installing" {
		return fmt.Sprintf(
			"\n%s\n\n%s\n\n%s\n",
			titleStyle.Render("AI CLI Manager - Installing"),
			m.message,
			"Press 'q' to quit",
		)
	}

	if m.mode == "config" {
		return m.viewConfig()
	}

	if m.mode == "mcp" {
		return m.viewMCP()
	}

	if m.mode == "table" {
		help := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Render("↑/↓: Navigate • Enter: Install selected • M: Configure MCP • R: Refresh status • Esc: Main menu • Q: Quit")

		statusInfo := ""
		installedCount := 0
		for _, tool := range m.tools {
			if tool.Installed {
				installedCount++
			}
		}
		statusInfo = fmt.Sprintf("Status: %d/%d tools installed", installedCount, len(m.tools))

		return fmt.Sprintf(
			"\n%s\n\n%s\n\n%s\n%s\n\n%s\n",
			titleStyle.Render("AI CLI Tools Manager"),
			statusInfo,
			m.table.View(),
			m.message,
			help,
		)
	}

	// Menu view
	installedCount := 0
	for _, tool := range m.tools {
		if tool.Installed {
			installedCount++
		}
	}

	statusText := fmt.Sprintf("Status: %d/%d tools installed", installedCount, len(m.tools))

	githubStatus := "Not configured"
	if m.githubUser != "" && m.githubRepo != "" {
		githubStatus = fmt.Sprintf("Synced to %s/%s", m.githubUser, m.githubRepo)
	}

	menu := fmt.Sprintf(`
%s

%s
GitHub: %s

Choose an option:

%s 1. Return to tools table (main view)
%s 2. Install all missing tools
%s 3. Configure GitHub sync
%s 4. Configure MCP servers
%s 5. Refresh installation status

%s Q. Quit

Press 1 or Esc to return to the tools table.

%s
`,
		titleStyle.Render("🤖 AI CLI Manager - Main Menu"),
		successStyle.Render(statusText),
		githubStatus,
		selectedStyle.Render("→"),
		selectedStyle.Render("→"),
		selectedStyle.Render("→"),
		selectedStyle.Render("→"),
		selectedStyle.Render("→"),
		selectedStyle.Render("→"),
		m.message,
	)

	return menuStyle.Render(menu)
}