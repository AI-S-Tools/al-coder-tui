package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lpm/ai-cli-manager/src"
)

func main() {
	p := tea.NewProgram(src.NewModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}