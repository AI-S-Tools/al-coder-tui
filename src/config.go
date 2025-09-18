package src

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func loadAITools() []AITool {
	// Load tools from ai_tools.json in the project directory
	data, err := os.ReadFile("ai_tools.json")
	if err != nil {
		// Fallback: try relative path
		data, err = os.ReadFile("../ai_tools.json")
		if err != nil {
			// Last fallback: try in current working directory
			if wd, err := os.Getwd(); err == nil {
				data, err = os.ReadFile(filepath.Join(wd, "ai_tools.json"))
			}
		}
	}

	if err != nil {
		// Return empty slice if no JSON file found
		return []AITool{}
	}

	var tools []AITool
	if err := json.Unmarshal(data, &tools); err != nil {
		// Return empty slice if JSON is invalid
		return []AITool{}
	}

	// Also save to user config for backup
	saveAITools(tools)
	return tools
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