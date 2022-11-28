package widgets

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"me.kryptk.marco/utils"
)

type BarMsg struct {
	Selected *string
	Next     *Bar
}

type Bar struct {
	Message  string
	Triggers []string
}

func (b Bar) Init() tea.Cmd {
	return nil
}

func (b Bar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, r := range b.Triggers {
			if r == msg.String() {
				return b, func() tea.Msg {
					return BarMsg{
						Selected: utils.Ptr(msg.String()),
					}
				}
			}
		}
	}

	return b, nil
}

func (b Bar) View() string {
	return lipgloss.NewStyle().Faint(true).PaddingBottom(1).Render(fmt.Sprintf("%s", b.Message))
}
