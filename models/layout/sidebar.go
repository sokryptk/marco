package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var tabBody = lipgloss.NewStyle().MarginRight(3)
var tab = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 5).Width(20).AlignHorizontal(0.5)
var activeTab = tab.Copy().Border(lipgloss.DoubleBorder())

type Sidebar struct {
	items    []teaModelWithName
	selected int
}

func (s Sidebar) Init() tea.Cmd {
	return nil
}

func (s Sidebar) Update(msg tea.Msg) (Sidebar, tea.Cmd) {
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
