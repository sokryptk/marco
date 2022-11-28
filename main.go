package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"me.kryptk.marco/models/layout"
	"os"
)

func main() {
	_, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}

	home := layout.NewHome()
	defer home.Close()

	p := tea.NewProgram(home, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
