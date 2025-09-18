package src

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func loadAITools() []AITool {
	// Default AI tools with MCP server configurations - embedded in binary
	defaultTools := []AITool{
		{
			Name:        "Claude Code",
			CLICommand:  "claude",
			InstallCmd:  "npm install -g @anthropic/claude-cli",
			CheckCmd:    "claude --version",
			Description: "Anthropic's Claude AI coding assistant",
			GitHubRepo:  "https://github.com/anthropics/claude-cli",
			MCPServers: []MCPServerConfig{
				{
					Name:        "filesystem",
					Command:     "npx",
					Args:        []string{"-y", "@modelcontextprotocol/server-filesystem"},
					Description: "File system access for Claude",
				},
				{
					Name:        "github",
					Command:     "npx",
					Args:        []string{"-y", "@modelcontextprotocol/server-github"},
					Env:         map[string]string{"GITHUB_PERSONAL_ACCESS_TOKEN": "${GITHUB_TOKEN}"},
					Description: "GitHub integration for Claude",
				},
			},
		},
		{
			Name:        "Aider",
			CLICommand:  "aider",
			InstallCmd:  "pip install aider-chat",
			CheckCmd:    "aider --version",
			Description: "AI pair programming in your terminal",
			GitHubRepo:  "https://github.com/paul-gauthier/aider",
			MCPServers: []MCPServerConfig{
				{
					Name:        "git",
					Command:     "npx",
					Args:        []string{"-y", "@modelcontextprotocol/server-git"},
					Description: "Git integration for AI tools",
				},
			},
		},
		{
			Name:        "Continue",
			CLICommand:  "continue",
			InstallCmd:  "npm install -g continue",
			CheckCmd:    "continue --version",
			Description: "Open-source AI code assistant",
			GitHubRepo:  "https://github.com/continuedev/continue",
			MCPServers: []MCPServerConfig{
				{
					Name:        "sqlite",
					Command:     "npx",
					Args:        []string{"-y", "@modelcontextprotocol/server-sqlite"},
					Description: "SQLite database access",
				},
			},
		},
		{
			Name:        "GitHub Copilot CLI",
			CLICommand:  "gh-copilot",
			InstallCmd:  "gh extension install github/gh-copilot",
			CheckCmd:    "gh copilot --version",
			Description: "GitHub Copilot command-line interface",
		},
		{
			Name:        "Codeium",
			CLICommand:  "codeium",
			InstallCmd:  "curl -Ls https://github.com/Exafunction/codeium/releases/latest/download/install.sh | bash",
			CheckCmd:    "codeium --version",
			Description: "Free AI code completion",
			GitHubRepo:  "https://github.com/Exafunction/codeium-cli",
		},
		{
			Name:        "Cursor",
			CLICommand:  "cursor",
			InstallCmd:  "brew install --cask cursor",
			CheckCmd:    "cursor --version",
			Description: "AI-powered code editor",
			MCPServers: []MCPServerConfig{
				{
					Name:        "typescript",
					Command:     "npx",
					Args:        []string{"-y", "@modelcontextprotocol/server-typescript"},
					Description: "TypeScript language server",
				},
				{
					Name:        "python",
					Command:     "uvx",
					Args:        []string{"mcp-server-python"},
					Description: "Python language server",
				},
			},
		},
		{
			Name:        "Qodo",
			CLICommand:  "qodo",
			InstallCmd:  "npm install -g @qodo/cli",
			CheckCmd:    "qodo --version",
			Description: "AI test generation and code quality",
			GitHubRepo:  "https://github.com/qodo-ai/qodo-cli",
		},
		{
			Name:        "Windsurf",
			CLICommand:  "windsurf",
			InstallCmd:  "brew install --cask windsurf",
			CheckCmd:    "windsurf --version",
			Description: "Codeium's AI-powered IDE",
		},
		{
			Name:        "Ollama",
			CLICommand:  "ollama",
			InstallCmd:  "brew install ollama",
			CheckCmd:    "ollama --version",
			Description: "Run large language models locally",
			MCPServers: []MCPServerConfig{
				{
					Name:        "ollama",
					Command:     "npx",
					Args:        []string{"-y", "mcp-server-ollama"},
					Description: "Ollama model server",
				},
			},
		},
		{
			Name:        "LM Studio CLI",
			CLICommand:  "lms",
			InstallCmd:  "brew install --cask lm-studio",
			CheckCmd:    "lms --version",
			Description: "Local LLM management",
		},
		{
			Name:        "Sourcegraph Cody",
			CLICommand:  "cody",
			InstallCmd:  "brew install sourcegraph/cody/cody-cli",
			CheckCmd:    "cody --version",
			Description: "Sourcegraph's AI coding assistant",
			GitHubRepo:  "https://github.com/sourcegraph/cody-cli",
		},
		{
			Name:        "Amazon Q",
			CLICommand:  "q",
			InstallCmd:  "brew install --cask amazon-q",
			CheckCmd:    "q --version",
			Description: "Amazon's AI developer assistant",
		},
		{
			Name:        "Tabnine CLI",
			CLICommand:  "tabnine",
			InstallCmd:  "curl -fsSL https://raw.githubusercontent.com/codota/tabnine-cli/master/install.sh | bash",
			CheckCmd:    "tabnine --version",
			Description: "AI code completion",
			GitHubRepo:  "https://github.com/codota/tabnine-cli",
		},
		{
			Name:        "Pieces CLI",
			CLICommand:  "pieces",
			InstallCmd:  "brew install pieces-cli",
			CheckCmd:    "pieces --version",
			Description: "AI-powered code snippet manager",
		},
		{
			Name:        "Mentat",
			CLICommand:  "mentat",
			InstallCmd:  "pip install mentat",
			CheckCmd:    "mentat --version",
			Description: "AI coding assistant with context awareness",
			GitHubRepo:  "https://github.com/AbanteAI/mentat",
		},
		{
			Name:        "GPT Engineer",
			CLICommand:  "gpt-engineer",
			InstallCmd:  "pip install gpt-engineer",
			CheckCmd:    "gpt-engineer --version",
			Description: "AI engineer that builds entire codebases",
			GitHubRepo:  "https://github.com/gpt-engineer-org/gpt-engineer",
		},
		{
			Name:        "Smol Developer",
			CLICommand:  "smol-dev",
			InstallCmd:  "pip install smol-developer",
			CheckCmd:    "smol-dev --version",
			Description: "Smallest AI developer",
			GitHubRepo:  "https://github.com/smol-ai/developer",
		},
		{
			Name:        "Auto-GPT",
			CLICommand:  "autogpt",
			InstallCmd:  "pip install auto-gpt",
			CheckCmd:    "autogpt --version",
			Description: "Autonomous GPT-4 agent",
			GitHubRepo:  "https://github.com/Significant-Gravitas/AutoGPT",
		},
		{
			Name:        "Open Interpreter",
			CLICommand:  "interpreter",
			InstallCmd:  "pip install open-interpreter",
			CheckCmd:    "interpreter --version",
			Description: "Natural language interface for computers",
			GitHubRepo:  "https://github.com/OpenInterpreter/open-interpreter",
			MCPServers: []MCPServerConfig{
				{
					Name:        "code-execution",
					Command:     "npx",
					Args:        []string{"-y", "mcp-server-code-execution"},
					Description: "Safe code execution environment",
				},
			},
		},
		{
			Name:        "Sweep AI",
			CLICommand:  "sweep",
			InstallCmd:  "pip install sweep-ai",
			CheckCmd:    "sweep --version",
			Description: "AI-powered code reviewer",
			GitHubRepo:  "https://github.com/sweepai/sweep",
		},
	}

	// Try to load from user's custom JSON file first
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".ai-cli-manager", "tools.json")

	if data, err := os.ReadFile(configPath); err == nil {
		var tools []AITool
		if json.Unmarshal(data, &tools) == nil && len(tools) > 0 {
			return tools
		}
	}

	// Save default tools to user config on first run
	saveAITools(defaultTools)
	return defaultTools
}

func saveAITools(tools []AITool) error {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".ai-cli-manager")
	configPath := filepath.Join(configDir, "tools.json")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(tools, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// ExportToolsConfig exports the current tools configuration to a JSON file
func ExportToolsConfig(tools []AITool, filename string) error {
	data, err := json.MarshalIndent(tools, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// ImportToolsConfig imports tools configuration from a JSON file
func ImportToolsConfig(filename string) ([]AITool, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var tools []AITool
	if err := json.Unmarshal(data, &tools); err != nil {
		return nil, err
	}

	return tools, nil
}