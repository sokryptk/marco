package pages

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Padding(1, 0)
	tabStyle       = lipgloss.NewStyle().PaddingLeft(1).Bold(true)
	activeTabStyle = tabStyle.Copy().Foreground(lipgloss.Color("5")).Border(bo, false, false, false, true).BorderForeground(lipgloss.Color("5"))
	bo             = lipgloss.Border{
		Left: "▐",
	}
)
