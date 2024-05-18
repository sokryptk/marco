package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"me.kryptk.marco/models/messages"
)

var content = lipgloss.NewStyle().Border(lipgloss.NormalBorder())

type Content struct {
	width, height int
	model         tea.Model
	dialog        tea.Model
}

func (c Content) Init() tea.Cmd {
	return nil
}

func (c Content) Update(msg tea.Msg) (Content, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width, c.height = msg.Width, msg.Height

		msg.Width -= content.GetHorizontalFrameSize()
		msg.Height -= content.GetVerticalFrameSize()

		c.model, cmd = c.model.Update(msg)
	case PageSwitchMsg:
		c.model = msg.new
	case messages.ShowDialogMsg:
		c.model = msg.Dialog
	case tea.KeyMsg:
		c.model, cmd = c.model.Update(msg)
	default:
		c.model, cmd = c.model.Update(msg)
	}

	return c, cmd
}

func (c Content) View() string {
	width, height := content.GetHorizontalFrameSize(), content.GetVerticalFrameSize()
	return content.Height(c.height - height).Width(c.width - width).AlignHorizontal(0.5).Render(c.model.View())
}
