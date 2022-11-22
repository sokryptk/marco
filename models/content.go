package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var content = lipgloss.NewStyle().Border(lipgloss.NormalBorder())

type Content struct {
	width, height int
	model         tea.Model
}

func (c Content) Init() tea.Cmd {
	return nil
}

func (c Content) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (c Content) View() string {
	width := content.GetHorizontalBorderSize() * 2
	return content.Height(c.height).Width(c.width - width).Render(c.model.View())
}
