package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var tabBody = lipgloss.NewStyle().Margin(0, 1)
var tab = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 3)
var activeTab = tab.Copy().Border(lipgloss.ThickBorder())

type Sidebar struct {
	items    []teaModelWithName
	selected int
}

func (s Sidebar) Init() tea.Cmd {
	return nil
}

func (s Sidebar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}

func (s Sidebar) View() string {
	return tabBody.Render(lipgloss.JoinVertical(lipgloss.Center, s.toTabs()...))
}

func (s Sidebar) toTabs() []string {
	names := make([]string, len(s.items))

	for i, m := range s.items {
		if i == s.selected {
			names[i] = activeTab.Render(m.Name())
			continue
		}

		names[i] = tab.Render(m.Name())
	}

	return names
}
