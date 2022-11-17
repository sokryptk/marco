package main

import (
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
    "me.kryptk.marco/models"
    "os"
)

func main() {
    p := tea.NewProgram(models.NewHome())
    if err := p.Start(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }

}
