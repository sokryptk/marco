package widgets

import tea "github.com/charmbracelet/bubbletea"

// Dialog is a simple over the screen modal
type Dialog struct {
	tea.Model
}

func (dialog Dialog) Init() tea.Cmd {
	return nil
}

func (dialog Dialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (dialog Dialog) View() string {
	return ""
}
