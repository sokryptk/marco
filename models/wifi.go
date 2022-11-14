package models

import (
    tea "github.com/charmbracelet/bubbletea"
    "me.kryptk.marco/repository"
    "me.kryptk.marco/services"
)

type WiFi struct {
    Service repository.WiFi
}

func NewWiFi() WiFi {
    return WiFi{
        Service: services.NewNMWiFi(),
    }
}
func (w WiFi) Init() tea.Cmd {
    return nil
}

func (w WiFi) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return w, nil
}

func (w WiFi) View() string {
    return "Wifi"
}