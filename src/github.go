package src

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh/v2/pkg/api"
)

func (m *Model) loadGitHubConfig() {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".ai-cli-manager", "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return
	}

	var config struct {
		GitHubUser string `json:"github_user"`
		GitHubRepo string `json:"github_repo"`
	}

	if json.Unmarshal(data, &config) == nil {
		m.githubUser = config.GitHubUser
		m.githubRepo = config.GitHubRepo
	}
}

func (m *Model) saveGitHubConfig() error {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".ai-cli-manager")
	configPath := filepath.Join(configDir, "config.json")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	config := struct {
		GitHubUser string `json:"github_user"`
		GitHubRepo string `json:"github_repo"`
	}{
		GitHubUser: m.githubUser,
		GitHubRepo: m.githubRepo,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func (m Model) checkGitHubCLI() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("gh", "--version")
		if err := cmd.Run(); err != nil {
			// Try to install gh if not present
			installCmd := exec.Command("brew", "install", "gh")
			if installCmd.Run() != nil {
				return githubSyncMsg{
					success: false,
					err:     fmt.Errorf("GitHub CLI (gh) is not installed. Please install it first"),
				}
			}
		}
		return nil
	}
}

func (m Model) syncWithGitHub() tea.Cmd {
	return func() tea.Msg {
		if m.githubUser == "" || m.githubRepo == "" {
			return githubSyncMsg{
				success: false,
				err:     fmt.Errorf("GitHub configuration not set"),
			}
		}

		// Create or update the repository
		repoName := fmt.Sprintf("%s/%s", m.githubUser, m.githubRepo)

		// Check if repo exists
		client, err := api.DefaultRESTClient()
		if err != nil {
			return githubSyncMsg{success: false, err: err}
		}

		// Check if repo exists
		err = client.Get(fmt.Sprintf("repos/%s", repoName), nil)
		if err != nil {
			// Create private repository
			// Create repo using gh CLI
			cmd := exec.Command("gh", "repo", "create", repoName, "--private", "-d", "AI CLI Manager Configuration")
			err = cmd.Run()
			if err != nil {
				return githubSyncMsg{success: false, err: err}
			}
		}

		// Save tools configuration as JSON
		data, err := json.MarshalIndent(m.tools, "", "  ")
		if err != nil {
			return githubSyncMsg{success: false, err: err}
		}

		// Create temp file
		tmpFile := "/tmp/ai-tools-config.json"
		if err := os.WriteFile(tmpFile, data, 0644); err != nil {
			return githubSyncMsg{success: false, err: err}
		}

		// Create gist or update repository file
		gistCmd := fmt.Sprintf("echo '%s' | gh gist create -f ai-tools.json -d 'AI CLI Tools Configuration' -", string(data))
		cmd := exec.Command("sh", "-c", gistCmd)
		if err := cmd.Run(); err != nil {
			return githubSyncMsg{success: false, err: err}
		}

		return githubSyncMsg{success: true}
	}
}

func (m *Model) pullFromGitHub() tea.Cmd {
	return func() tea.Msg {
		if m.githubUser == "" || m.githubRepo == "" {
			return githubSyncMsg{
				success: false,
				err:     fmt.Errorf("GitHub configuration not set"),
			}
		}

		// List gists and find the configuration
		cmd := exec.Command("gh", "gist", "list", "--limit", "100")
		output, err := cmd.Output()
		if err != nil {
			return githubSyncMsg{success: false, err: err}
		}

		// Parse output to find our config gist
		lines := strings.Split(string(output), "\n")
		var gistID string
		for _, line := range lines {
			if strings.Contains(line, "AI CLI Tools Configuration") {
				fields := strings.Fields(line)
				if len(fields) > 0 {
					gistID = fields[0]
					break
				}
			}
		}

		if gistID == "" {
			return githubSyncMsg{
				success: false,
				err:     fmt.Errorf("configuration gist not found"),
			}
		}

		// Download the gist
		cmd = exec.Command("gh", "gist", "view", gistID, "-f", "ai-tools.json")
		output, err = cmd.Output()
		if err != nil {
			return githubSyncMsg{success: false, err: err}
		}

		// Parse the tools
		var tools []AITool
		if err := json.Unmarshal(output, &tools); err != nil {
			return githubSyncMsg{success: false, err: err}
		}

		// Update local tools
		m.tools = tools
		saveAITools(tools)

		return githubSyncMsg{success: true}
	}
}

func (m Model) installFromGitHub(tool AITool) error {
	if tool.GitHubRepo == "" {
		return fmt.Errorf("no GitHub repository specified")
	}

	// Clone the repository
	tempDir := filepath.Join("/tmp", "ai-cli-install", tool.Name)
	os.RemoveAll(tempDir)
	defer os.RemoveAll(tempDir)

	cmd := exec.Command("git", "clone", tool.GitHubRepo, tempDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Look for installation script or use standard commands
	installScripts := []string{
		filepath.Join(tempDir, "install.sh"),
		filepath.Join(tempDir, "scripts/install.sh"),
		filepath.Join(tempDir, "setup.sh"),
	}

	for _, script := range installScripts {
		if _, err := os.Stat(script); err == nil {
			cmd := exec.Command("sh", script)
			cmd.Dir = tempDir
			return cmd.Run()
		}
	}

	// Try standard installation methods
	if _, err := os.Stat(filepath.Join(tempDir, "package.json")); err == nil {
		// Node.js project
		cmd := exec.Command("npm", "install", "-g", ".")
		cmd.Dir = tempDir
		return cmd.Run()
	}

	if _, err := os.Stat(filepath.Join(tempDir, "setup.py")); err == nil {
		// Python project
		cmd := exec.Command("pip", "install", ".")
		cmd.Dir = tempDir
		return cmd.Run()
	}

	if _, err := os.Stat(filepath.Join(tempDir, "go.mod")); err == nil {
		// Go project
		cmd := exec.Command("go", "install", ".")
		cmd.Dir = tempDir
		return cmd.Run()
	}

	return fmt.Errorf("no installation method found")
}

func (m Model) viewConfig() string {
	status := "Not configured"
	if m.githubUser != "" && m.githubRepo != "" {
		status = fmt.Sprintf("User: %s\nRepo: %s", m.githubUser, m.githubRepo)
	}

	return fmt.Sprintf(`
%s

Current GitHub Configuration:
%s

Options:
%s S: Sync configuration to GitHub
%s P: Pull configuration from GitHub
%s Esc: Back to menu

Enter GitHub username and repository name to configure.
Example: username/ai-cli-config

%s
`,
		titleStyle.Render("GitHub Configuration"),
		status,
		selectedStyle.Render("→"),
		selectedStyle.Render("→"),
		selectedStyle.Render("→"),
		m.message,
	)
}